import { api } from '$lib/pkg/fetch';

// 백엔드 응답 구조와 일치
export interface NotificationApiResponse {
	id: string; // 실제 DB UUID
	endpoint_name: string;
	body: string;
	is_read: boolean;
	created_at: string;
	mute: boolean;
}

export interface PaginatedNotiResponse {
	items: NotificationApiResponse[];
	next_cursor: string | null;
	has_more: boolean;
}

export interface DisplayNotification {
	id: string;
	endpointName: string;
	body: string;
	isRead: boolean;
	timestamp: string;
	createdAt: string;
	isMute: boolean;
}

function formatTimestamp(dateString: string): string {
	const date = new Date(dateString);
	const now = new Date();
	const diffMs = now.getTime() - date.getTime();
	const diffMins = Math.floor(diffMs / 60000);
	const diffHours = Math.floor(diffMs / 3600000);

	if (diffMins < 1) return '방금 전';
	if (diffMins < 60) return `${diffMins}분 전`;
	if (diffHours < 24) return `${diffHours}시간 전`;
	return date.toLocaleDateString('ko-KR', { month: 'short', day: 'numeric' });
}

// 데이터 변환 헬퍼
export function transformNotification(apiData: NotificationApiResponse): DisplayNotification {
	return {
		id: apiData.id,
		endpointName: apiData.endpoint_name,
		body: apiData.body,
		isRead: apiData.is_read,
		timestamp: formatTimestamp(apiData.created_at),
		createdAt: apiData.created_at,
		isMute: apiData.mute,
	};
}

// 실제 API 호출 (Cursor 지원)
export async function getNotifications(
	cursor?: string,
	endpointID?: string,
): Promise<PaginatedNotiResponse> {
	let path = `/notifications?limit=20`;

	if (cursor) path += `&cursor=${encodeURIComponent(cursor)}`;
	if (endpointID && endpointID !== 'ALL') path += `&endpoint_id=${encodeURIComponent(endpointID)}`;

	return await api<PaginatedNotiResponse>(path);
}

// 알림 읽음 처리
export async function markAsReadUntil(lastId: string, endpointID?: string): Promise<void> {
	const actualEndpointID = endpointID && endpointID !== 'ALL' ? endpointID : undefined;
	await api(`/notifications/read-until`, {
		method: 'POST',
		body: { last_id: lastId, endpoint_id: actualEndpointID },
	});
}

// 알림 삭제 처리
export async function deleteNotification(id: string): Promise<void> {
	await api(`/notifications/${id}`, {
		method: 'DELETE',
	});
}
