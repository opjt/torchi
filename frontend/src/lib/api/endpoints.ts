import { api, catchError, type Result } from '$lib/pkg/fetch';
export interface Endpoint {
	id: string;
	name: string;
	token: string;
	active: boolean;
}

export async function addEndpoint(name: string): Promise<Result<void>> {
	return await catchError(
		api<void>(`/endpoints`, {
			method: 'POST',
			body: {
				serviceName: name.trim(),
			},
		}),
	);
}

export async function fetchEndpoints(): Promise<Endpoint[]> {
	const res = await api<Endpoint[]>('/endpoints');
	return res;
}

export async function deleteEndpoint(token: string): Promise<void> {
	await api<void>(`/endpoints/${token}`, {
		method: 'DELETE',
	});
}

export async function muteEndpoint(token: string): Promise<Result<void>> {
	return await catchError(
		api<void>(`/endpoints/${token}/mute`, {
			method: 'POST',
		}),
	);
}

export async function unmuteEndpoint(token: string): Promise<Result<void>> {
	return catchError(
		api<void>(`/endpoints/${token}/mute`, {
			method: 'DELETE',
		}),
	);
}
