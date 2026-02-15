import { logout } from '$lib/client/auth/lifecycle';
import { getErrorMessage } from './error-message';
import { shouldShowToast, showToast, type ToastType } from './toast';

export type ApiError = {
	status: number;
	message?: string;
	code?: string;
};

type ApiOptions<TBody> = {
	method?: 'GET' | 'POST' | 'PUT' | 'PATCH' | 'DELETE';
	body?: TBody;
	headers?: HeadersInit;
	signal?: AbortSignal;
	toastType?: ToastType;
};

type ApiResponse<T> = {
	code: number;
	success: boolean;
	data: T;
	message?: string;
	error?: {
		code: string;
		message: string;
	};
};

// 토큰 갱신 중인지 확인하는 플래그 (중복 요청 방지)
let isRefreshing = false;
// 갱신을 기다리는 요청들의 큐 (동시 요청 처리용)
let refreshSubscribers: ((token: boolean) => void)[] = [];

const onRefreshed = (success: boolean) => {
	refreshSubscribers.forEach((callback) => callback(success));
	refreshSubscribers = [];
};

const showToastWrapper = (type: ToastType, message: string, code?: string) => {
	if (type === 'none') return;

	if (code === 'NETWORK_ERROR' && !shouldShowToast('NETWORK_ERROR')) {
		return;
	}

	if (type === 'error') {
		showToast.error(message);
	} else if (type === 'warning') {
		console.warn(`[Warning] ${message}`);
	}
};

const onAuthFailure: () => Promise<void> | void = async () => {
	showToast.message('안전한 이용을 위해 다시 한번 로그인해 주세요.');
	await logout(); //순환참조 괜찮을까?
};

export async function api<TResponse, TBody = unknown>(
	url: string,
	options: ApiOptions<TBody> = {},
): Promise<TResponse> {
	const {
		method = 'GET',
		body,
		headers,
		signal,
		toastType = 'error', // 토스트를 직접 컨트롤 하려면 'none' 으로 사용
	} = options;
	url = `/api${url}`;
	// 에러를 던지기 전에 토스트를 처리하는 내부 헬퍼
	const handleError = (status: number, message?: string, code?: string) => {
		showToastWrapper(toastType, getErrorMessage(code, message), code);

		const error: ApiError = {
			status,
			message,
			code,
		};
		throw error;
	};
	const executeRequest = async (): Promise<Response> => {
		return fetch(url, {
			method,
			credentials: 'include',
			headers: { 'Content-Type': 'application/json', ...headers },
			body: body ? JSON.stringify(body) : undefined,
			signal,
		});
	};

	try {
		const res = await executeRequest();
		// Access Token 만료 (401) 발생 시
		if (res.status === 401 && !url.includes('/auth/refresh')) {
			// 갱신 시도 로직
			if (!isRefreshing) {
				isRefreshing = true;

				(async () => {
					try {
						const refreshRes = await fetch(`/api/auth/refresh`, {
							method: 'POST',
							credentials: 'include',
						});

						if (refreshRes.ok) {
							onRefreshed(true);
						} else {
							// 리프레시 토큰 만료 -> 로그아웃
							onRefreshed(false);
							onAuthFailure();
						}
					} catch (_) {
						// 네트워크 에러 등
						onRefreshed(false);
						onAuthFailure();
					} finally {
						isRefreshing = false;
					}
				})();
			}

			// 여기서 new Promise를 리턴하여 현재 요청을 대기 시킵니다.
			return new Promise<TResponse>((resolve, reject) => {
				refreshSubscribers.push((success) => {
					//onRefreshed 가 호출되면 아래 로직 수행
					if (success) {
						// 재요청 성공 시 resolve
						resolve(api<TResponse, TBody>(url, options));
					} else {
						// 재요청 실패(로그아웃) 시 reject
						reject({ status: 401, message: 'Session Expired', code: 'SESSION_EXPIRED' });
					}
				});
			});
		}

		// HTTP 에러 처리 (400~500번대)
		if (!res.ok) {
			const { status, statusText } = res;
			let message = statusText;
			let code = String(status);

			const errData: ApiResponse<never> = await res.json().catch(() => ({}));

			if (errData.error) {
				message = errData.error.message;
				code = errData.error.code;
			}

			return handleError(status, message, code);
		}

		// 204 No Content 처리
		if (res.status === 204) {
			return undefined as TResponse;
		}

		// 응답 바디 파싱 (JSON이 아닐 경우 대비)
		let json: ApiResponse<TResponse>;
		try {
			json = await res.json();
		} catch (_) {
			// 200 OK인데 JSON이 아닌 경우
			return handleError(res.status, 'Invalid JSON response');
		}

		//  비즈니스 로직 에러 핸들링 (success: false)
		if (!json.success) {
			return handleError(
				json.code ?? res.status,
				json.message ?? 'Unknown API Error',
				String(json.code),
			);
		}

		return json.data;
	} catch (error) {
		// 이미 handleError에서 던져진 ApiError라면 그대로 통과
		if ((error as ApiError).status) {
			throw error;
		}
		// 진짜 네트워크 에러 등 예상치 못한 에러
		// 백엔드 서버 통신이 안되는 경우
		return handleError(0, undefined, 'NETWORK_ERROR');
	}
}

// Result 타입 개선: 에러 타입을 ApiError로 구체화 가능하게 함
export type Result<T, E = ApiError> = { ok: true; data: T } | { ok: false; error: E };

export async function catchError<T>(promise: Promise<T>): Promise<Result<T>> {
	try {
		const data = await promise;
		return { ok: true, data };
	} catch (error) {
		return { ok: false, error: error as ApiError };
	}
}
