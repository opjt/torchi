import { goto } from '$app/navigation';
import { resolve } from '$app/paths';

import { api } from '$lib/pkg/fetch';
import { auth } from '$lib/client/auth/auth';
import { push } from '../pushManager.svelte';

export async function logout() {
	const endpint = await push.getEndpoint();
	console.log(endpint);
	await api<void>(`/auth/logout`, {
		method: 'POST',
		body: { endpoint: endpint },
	});
	auth.logout();
}

export async function guestLogin() {
	const userId = localStorage.getItem('guest_user_id');

	const res = await api<{ user_id: string }>(`/auth/guest`, {
		method: 'POST',
		body: { user_id: userId ?? null },
	});

	localStorage.setItem('guest_user_id', res.user_id);
	await auth.init();
	goto(resolve('/app'));
}
