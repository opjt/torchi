import { api, catchError, type Result } from '$lib/pkg/fetch';
export interface UserInfo {
	user_id: string;
	email?: string;
	terms_agreed: boolean;
	is_guest: boolean;
}

export async function fetchWhoami(): Promise<UserInfo | null> {
	const res = await api<UserInfo>(`/users/whoami`);
	return res;
}

export async function agreeToTerms(): Promise<void> {
	await api<void>(`/users/terms-agree`, {
		method: 'POST',
	});
}

export async function withdraw(): Promise<Result<void>> {
	return catchError(
		api<void>(`/users`, {
			method: 'DELETE',
		}),
	);
}
