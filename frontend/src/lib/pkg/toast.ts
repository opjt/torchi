import { toast } from 'svelte-sonner';

export const defaultToastCooldown = 1400; // ms
export const errorToastDuration = 2300;

const lastToastMap = new Map<string, number>();

// 같은 의미의 토스트를 던지기 싫을 때 사용하는 확인 함수
export function shouldShowToast(key: string): boolean {
	const now = Date.now();
	const last = lastToastMap.get(key) ?? 0;

	if (now - last < defaultToastCooldown) {
		return false;
	}

	lastToastMap.set(key, now);
	return true;
}

export type ToastType = 'error' | 'warning' | 'info' | 'message' | 'none';

const TOAST_DURATION: Record<Exclude<ToastType, 'none'>, number> = {
	info: 1500,
	warning: 2500,
	error: 2300,
	message: 3000,
};

export const showToast = {
	info: (msg: string) => {
		toast.info(msg, { duration: TOAST_DURATION.info });
	},
	warning: (msg: string) => {
		toast.warning(msg, { duration: TOAST_DURATION.warning });
	},
	error: (msg: string) => {
		toast.error(msg, { duration: TOAST_DURATION.error });
	},
	message: (msg: string) => {
		toast.message(msg, { duration: TOAST_DURATION.message });
	},
};
