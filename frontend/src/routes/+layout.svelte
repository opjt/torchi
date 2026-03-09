<script lang="ts">
	import '$src/app.css';
	import favicon from '$lib/assets/favicon.svg';
	import { Toaster } from '$lib/components/ui/sonner/index';
	import { push } from '$lib/client/pushManager.svelte';
	import { onMount, type ComponentProps } from 'svelte';
	import { toast } from 'svelte-sonner';
	import { isDev } from '$lib/pkg/util';
	import { defaultToastCooldown } from '$lib/pkg/toast';

	type Position = ComponentProps<typeof Toaster>['position'];

	let position = $state<Position>('top-center');
	const swSessionStorageKey = 'sw-just-updated';
	function updatePosition() {
		position = window.innerWidth < 768 ? 'bottom-center' : 'top-center';
	}

	onMount(() => {
		push.initialize();
		updatePosition();
		window.addEventListener('resize', updatePosition);
		return () => {
			window.removeEventListener('resize', updatePosition);
		};
	});

	onMount(async () => {
		if (isDev()) {
			toast.info('dev mode', { duration: 800 });
			return;
		}
		if (typeof window !== 'undefined') {
			const { registerSW } = await import('virtual:pwa-register');

			const updateServiceWorker = registerSW({
				immediate: true,

				onNeedRefresh() {
					console.log('새 버전 발견');

					// 토스트로 알림 (자동 업데이트 예정)
					toast.info('새 버전을 다운로드 중입니다...', {
						duration: 2000,
					});

					// 3초 후 자동 업데이트 (사용자가 현재 작업 마무리할 시간)
					setTimeout(() => {
						toast.success('업데이트를 적용합니다...', {
							duration: 1200,
						});

						setTimeout(() => {
							sessionStorage.setItem(swSessionStorageKey, 'true');
							updateServiceWorker(true);
						}, 1400);
					}, 2500);
				},

				onOfflineReady() {
					console.log('오프라인 사용이 가능합니다!');
				},
			});
		}
	});

	// 업데이트 완료 후 토스트
	onMount(() => {
		const justUpdated = sessionStorage.getItem(swSessionStorageKey);

		if (justUpdated === 'true') {
			toast.success('앱이 최신 버전으로 업데이트되었습니다!', {
				duration: 3000,
			});
			sessionStorage.removeItem(swSessionStorageKey);
		}
	});

	let { children } = $props();
</script>

<svelte:head>
	<link rel="icon" href={favicon} />
</svelte:head>

<Toaster duration={defaultToastCooldown} {position} />
{@render children()}
