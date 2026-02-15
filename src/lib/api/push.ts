import { api, catchError, type Result } from '$lib/pkg/fetch';

export interface CheckSubscriptionResponse {
	isOwner: boolean;
}

export async function checkSubscription(
	endpoint: string,
): Promise<Result<CheckSubscriptionResponse>> {
	return catchError(
		api<CheckSubscriptionResponse>(`/subscriptions/check`, {
			method: 'POST',
			body: { endpoint: endpoint },
		}),
	);
}

export async function subscribe(subscription: PushSubscription): Promise<void> {
	await api<void>(`/subscriptions`, {
		method: 'POST',
		body: subscription,
	});
}
