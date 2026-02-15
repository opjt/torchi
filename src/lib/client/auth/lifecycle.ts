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

export async function fakeLogin() {
	await api<void>(`/auth/fake/login`);
	await auth.init();
	goto(resolve('/app'));
}
