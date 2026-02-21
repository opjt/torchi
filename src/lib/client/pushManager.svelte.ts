import { browser } from '$app/environment';
import { subscribe } from '$lib/api/push';
import { PUBLIC_VAPID_KEY } from '$lib/config';
import { api, catchError } from '$lib/pkg/fetch';
import { toast } from 'svelte-sonner';

const PushError = {
	PERMISSION_DENIED: {
		type: 'error',
		message: '브라우저 알림 권한이 차단되어 있어요. 권한을 허용해 주세요.',
	},
	SUBSCRIBE_FAILED: {
		type: 'error',
		message: '알림 신청을 완료하지 못했어요.',
	},
	SUBSCRIBE_SUCCESS: {
		type: 'success',
		message: '이제 실시간 알림을 받을 수 있어요.',
	},
	UNSUBSCRIBE_FAILED: {
		type: 'error',
		message: '알림 해제 처리에 실패했어요.',
	},
	UNSUBSCRIBE_SUCCESS: {
		type: 'success',
		message: '알림 구독을 해제했어요.',
	},
	UNSUPPORTED: {
		type: 'error',
		message: '현재 환경에서는 알림 기능을 지원하지 않아요.',
	},
	SW_REG_FAILED: {
		type: 'error',
		message: '서비스 워커 초기화에 실패했어요.',
	},
} as const;

class PushNotificationManager {
	isLoading = $state(true);
	isToggling = $state(false); // 푸시 알림 상태변환 loading...
	isSubscribed = $state(false);
	subscription = $state<PushSubscription | null>(null);
	permissionState = $state<NotificationPermission | null>(null);

	private VAPID_PUBLIC_KEY = PUBLIC_VAPID_KEY;

	private notifyError(localError: (typeof PushError)[keyof typeof PushError], e?: unknown) {
		// 1. 상세 메시지(description) 추출 로직
		let description: string | null = null;

		if (e instanceof Error) {
			description = e.message;
		} else if (typeof e === 'object' && e !== null && 'error' in e) {
			// api 유틸리티의 공통 에러 구조인 { error: { message: string } } 대응
			const errorBody = (e as { error: { message?: string } }).error;
			description = errorBody?.message || null;
		}

		// 2. 성공/실패 여부에 따른 토스트 함수 결정
		const toastFn = localError.type === 'success' ? toast.success : toast.error;

		// 3. 토스트 출력
		toastFn(localError.message, {
			...(description && { description: `${description}` }),
		});
	}

	get isSupported() {
		return browser && 'serviceWorker' in navigator && 'Notification' in window;
	}

	private checkSupport(): boolean {
		if (!this.isSupported) {
			console.warn('[PushManager] 이 브라우저는 푸시 알림을 지원하지 않습니다.');
			toast.error('지원되지 않는 브라우저입니다.');
			return false;
		}
		return true;
	}

	constructor() {
		if (browser) {
			// 브라우저가 Notification 기능이 있을 경우만.
			this.permissionState = 'Notification' in window ? Notification.permission : 'default';
		}
	}

	async initialize() {
		if (!browser) return;

		// 1. 지원 여부 확인
		if (!this.checkSupport()) {
			this.isLoading = false;
			return;
		}

		try {
			this.watchPermission();

			if (Notification.permission == 'denied') {
				// 차단된 상태라면 명시적으로 상태 업데이트
				this.isSubscribed = false;
				this.subscription = null;
			} else {
				await this.loadSubscription();
			}
		} catch (e) {
			console.error('[PushManager] Initialization failed:', e);
		} finally {
			this.isLoading = false;
		}
	}
	async getEndpoint(): Promise<string | null> {
		if (!this.checkSupport()) return null;

		try {
			// 1. 서비스 워커가 활성화될 때까지 대기
			const reg = await navigator.serviceWorker.ready;

			// 2. 기존 구독 정보 확인
			const sub = await reg.pushManager.getSubscription();

			if (!sub) {
				console.log('[PushManager] 등록된 구독 정보가 없습니다.');
				return null;
			}

			// 3. 엔드포인트 URL 반환
			return sub.endpoint;
		} catch (e) {
			console.error('[PushManager] 엔드포인트 추출 실패:', e);
			return null;
		}
	}
	private watchPermission() {
		if (!this.checkSupport()) return;

		navigator.permissions.query({ name: 'notifications' }).then((status) => {
			status.onchange = async () => {
				const newState = Notification.permission;
				this.permissionState = newState;
				console.log('[PushManager] Permission changed:', newState);

				// 권한이 취소(denied)되었다면 상태 초기화 및 서버 통보 시도
				if ((newState === 'denied' || newState === 'default') && this.subscription) {
					await this.handleUnsubscribe(); // 재사용
				} else if (newState === 'granted' && !this.subscription) {
					await this.loadSubscription();
				}
			};
		});
	}

	private urlBase64ToUint8Array(base64String: string) {
		const padding = '='.repeat((4 - (base64String.length % 4)) % 4);
		const base64 = (base64String + padding).replace(/-/g, '+').replace(/_/g, '/');
		const rawData = atob(base64);
		return Uint8Array.from([...rawData].map((c) => c.charCodeAt(0)));
	}

	// 구독 로직
	async handleSubscribe() {
		if (!this.checkSupport()) return;

		this.isToggling = true;
		let tempSub: PushSubscription | null = null;

		try {
			const reg = await navigator.serviceWorker.ready;

			// 1. 권한 요청
			const permission = await Notification.requestPermission();
			this.permissionState = permission;
			if (permission !== 'granted') throw new Error('알림 권한이 거부되었습니다.');

			// 2. 브라우저 구독 생성
			// 기존 구독이 있다면 그것을 사용, 없다면 새로 생성
			tempSub = await reg.pushManager.getSubscription();
			if (!tempSub) {
				tempSub = await reg.pushManager.subscribe({
					userVisibleOnly: true,
					applicationServerKey: this.urlBase64ToUint8Array(this.VAPID_PUBLIC_KEY),
				});
			}

			// 3. 서버 전송 (중요: this.subscription이 아닌 tempSub를 보냄)
			// PushSubscription 객체는 toJSON()을 가지고 있으므로 이를 명시적으로 호출하거나,
			// JSON.stringify가 내부적으로 호출하도록 객체 자체를 넘깁니다.
			await subscribe(tempSub);

			// 4. 모든 과정 성공 시 상태 업데이트
			this.subscription = tempSub;
			this.isSubscribed = true;
			this.notifyError(PushError.SUBSCRIBE_SUCCESS);
		} catch (_) {
			// 서버 등록 실패 시 브라우저 구독도 취소 (데이터 정합성 유지)
			if (tempSub) {
				await tempSub.unsubscribe().catch(() => {}); // 롤백 실패는 무시
			}

			this.subscription = null;
			this.isSubscribed = false;
		} finally {
			this.isToggling = false;
		}
	}

	// 구독 해제 로직
	async handleUnsubscribe() {
		// 이미 구독 정보가 없다면 종료
		if (!this.subscription) return;
		this.isToggling = true;
		const subToUnsubscribe = this.subscription; // 현재 구독 객체 캡처

		try {
			await catchError(
				// 실패해도 조용히 넘어갑니다.
				api<void>(`/subscriptions/unsubscribe`, {
					method: 'POST',
					body: subToUnsubscribe.toJSON(),
					toastType: 'none',
				}),
			);
			await subToUnsubscribe.unsubscribe(); //브라우저 구독 해제
			this.notifyError(PushError.UNSUBSCRIBE_SUCCESS);
		} catch (e) {
			this.notifyError(PushError.UNSUBSCRIBE_FAILED, e);
		} finally {
			// 성공하든 실패하든 클라이언트 상태는 초기화 (사용자 입장에서 해제됨)
			this.subscription = null;
			this.isSubscribed = false;

			this.isToggling = false;
		}
	}

	async loadSubscription(initialize = false) {
		if (!this.checkSupport()) return;
		try {
			const reg = await navigator.serviceWorker.ready;
			const sub = await reg.pushManager.getSubscription();

			this.subscription = sub;
			this.permissionState = Notification.permission;

			this.isSubscribed = !!sub;
			if (initialize && this.subscription) {
				await subscribe(this.subscription);
			}

			console.log('[PushManager] Loaded subscription:', !!sub);
		} catch (e) {
			this.isSubscribed = false;
			console.error('[PushManager] Load failed:', e);
		}
	}

	async testNotification() {
		if (!this.subscription) {
			return;
		}

		await catchError(
			api<void>(`/push-test`, {
				method: 'POST',
				body: this.subscription.toJSON(),
			}),
		);
	}

	async handleDemoPush(message: string) {
		if (!this.checkSupport()) return;

		this.isToggling = true;
		let tempSub: PushSubscription | null = null;
		let isNewSubscription = false; // 이번에 새로 만들었는지 여부 확인

		try {
			const reg = await navigator.serviceWorker.ready;

			// 1. 권한 요청
			const permission = await Notification.requestPermission();
			this.permissionState = permission;

			if (permission !== 'granted') {
				this.notifyError(PushError.PERMISSION_DENIED);
				return;
			}

			// 기존에 구독된 게 있는지 먼저 확인
			tempSub = await reg.pushManager.getSubscription();

			if (!tempSub) {
				// 기존 구독이 없을 때만 새로 생성하고, '새로 만들었다'고 표시
				tempSub = await reg.pushManager.subscribe({
					userVisibleOnly: true,
					applicationServerKey: this.urlBase64ToUint8Array(this.VAPID_PUBLIC_KEY),
				});
				isNewSubscription = true;
			}

			// 데모 전용 엔드포인트 호출
			const subJson = tempSub.toJSON();

			await api<void>(`/api/push-demo`, {
				method: 'POST',
				body: {
					endpoint: subJson.endpoint,
					auth: subJson.keys?.auth,
					p256dh: subJson.keys?.p256dh,
					message: message,
				},
			});
		} catch (e) {
			// 에러가 났는데, 만약 이번에 새로 만든 임시 구독이었다면 즉시 삭제
			if (tempSub && isNewSubscription) {
				await tempSub.unsubscribe().catch(() => {});
			}
			console.log(e);
		} finally {
			// 기존 유저(이미 구독한 유저)가 데모를 눌렀을 때는 구독을 유지해야 합니다.
			if (tempSub && isNewSubscription) {
				setTimeout(() => {
					tempSub?.unsubscribe().catch(() => {});
				}, 3000);
			}

			this.isToggling = false;
		}
	}
}

export const push = new PushNotificationManager();
