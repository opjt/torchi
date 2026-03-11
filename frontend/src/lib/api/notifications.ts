import { api } from '$lib/pkg/fetch';

export type NotificationStatus =
	| 'pending'
	| 'sent'
	| 'failed'
	| 'mute'
	| 'timeout_reply'
	| 'reacted';

// 백엔드 응답 구조와 일치
export interface NotificationApiResponse {
	id: string; // notification UUID
	endpoint_name: string;
	endpoint_id: string | null;
	body: string;
	is_read: boolean;
	created_at: string;
	status: NotificationStatus;
	mute: boolean;
	actions: string[] | null;
	reaction: string | null;
}

export interface PaginatedNotiResponse {
	items: NotificationApiResponse[];
	next_cursor: string | null;
	has_more: boolean;
}

export interface DisplayNotification {
	id: string;
	endpointName: string;
	endpointId: string | null;
	body: string;
	isRead: boolean;
	timestamp: string;
	createdAt: string;
	isMute: boolean;
	actions: string[] | null; // ex) ["승인", "거절"] or null
	reaction: string | null; // 반응한 리액션 값, 없으면 null
	isExpired: boolean; // timeout_reply면 true → 버튼 비활성화용
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
		endpointId: apiData.endpoint_id,
		body: apiData.body,
		isRead: apiData.is_read,
		timestamp: formatTimestamp(apiData.created_at),
		createdAt: apiData.created_at,
		isMute: apiData.mute,
		actions: apiData.actions,
		reaction: apiData.reaction,
		isExpired: apiData.status === 'timeout_reply',
	};
}

// 실제 API 호출 (Cursor 지원)
export async function getNotifications(
	cursor?: string,
	endpointID?: string,
	searchQuery?: string,
): Promise<PaginatedNotiResponse> {
	let path = `/notifications?limit=20`;

	if (cursor) path += `&cursor=${encodeURIComponent(cursor)}`;
	if (endpointID && endpointID !== 'ALL') path += `&endpoint_id=${encodeURIComponent(endpointID)}`;
	if (searchQuery) path += `&query=${encodeURIComponent(searchQuery)}`;
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

// 알림 리액션 api
export async function postReaction(id: string, action: string): Promise<void> {
	await api(`/v1/react/${id}`, {
		method: 'POST',
		body: { reaction: action },
	});
}
