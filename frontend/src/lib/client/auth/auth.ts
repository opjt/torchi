import { get, writable, derived, type Readable } from 'svelte/store';
import { fetchWhoami, type UserInfo } from '$lib/api/user';
import { push } from '../pushManager.svelte';

type AuthState = UserInfo | null | undefined;

type AuthStore = Readable<AuthState> & {
	init: () => Promise<void>;
	ready: Readable<boolean>;
	logout: () => void;
	isAuthenticated: () => boolean;
	getUser: () => UserInfo | null;
	hasAgreedToTerms: () => boolean;
	whenReady: () => Promise<void>;
	getDisplayName: () => string;
};

function createAuthStore(): AuthStore {
	const { subscribe, set } = writable<AuthState>(undefined);
	const readyStore = writable(false);

	let readyResolve: (() => void) | null = null;
	const readyPromise = new Promise<void>((resolve) => {
		readyResolve = resolve;
	});

	const store: AuthStore = {
		subscribe,

		init: async () => {
			try {
				const userInfo = await fetchWhoami();
				set(userInfo);
				await push.loadSubscription(true);
			} catch (error) {
				console.error('Auth init failed:', error);
				set(null);
			} finally {
				readyStore.set(true);
				readyResolve?.();
			}
		},

		ready: {
			subscribe: readyStore.subscribe,
		},

		isAuthenticated: () => {
			const state = get({ subscribe });
			return state !== null && state !== undefined;
		},

		getUser: () => {
			const state = get({ subscribe });
			return state !== null && state !== undefined ? state : null;
		},
		// 약관 동의 여부 확인
		hasAgreedToTerms: () => {
			const state = get({ subscribe });
			return state?.terms_agreed ?? false;
		},

		whenReady: () => readyPromise, // init 함수 대기

		logout: () => {
			set(null);
		},
		getDisplayName: () => {
			const state = get({ subscribe });
			if (!state) return 'Torchi User';
			if (state.is_guest) return `${state.user_id.slice(0, 8)}`;
			return state.email?.split('@')[0] ?? 'Torchi User';
		},
	};

	return store;
}

export const auth = createAuthStore();

// 편의를 위한 derived stores
export const isAuthenticated = derived(auth, ($auth) => {
	return $auth !== null && $auth !== undefined;
});

export const hasAgreedToTerms = derived(auth, ($auth) => {
	return $auth?.terms_agreed ?? false;
});
