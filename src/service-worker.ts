/// <reference types="@sveltejs/kit" />
/// <reference no-default-lib="true"/>
/// <reference lib="esnext" />
/// <reference lib="webworker" />

// eslint-disable-next-line @typescript-eslint/no-explicit-any
(self as any).__WB_DISABLE_DEV_LOGS = true;

interface ManifestEntry {
	url: string;
	revision: string | null;
}

declare const self: ServiceWorkerGlobalScope & {
	__WB_MANIFEST: (string | ManifestEntry)[];
};

import { precacheAndRoute, cleanupOutdatedCaches } from 'workbox-precaching';
import { registerRoute } from 'workbox-routing';
import { NetworkFirst } from 'workbox-strategies';

const sw = self as unknown as ServiceWorkerGlobalScope;

// 1. 이전 버전의 오래된 캐시 자동 삭제
cleanupOutdatedCaches();

// 2. 프리캐싱 (빌드 시 생성된 자산들)
precacheAndRoute(self.__WB_MANIFEST);

// 3. 라이프사이클 제어
sw.addEventListener('install', () => {
	console.log('[SW] 설치 완료, 대기 중...');
});

// 클라이언트(App.svelte)에서 updateServiceWorker(true)를 호출하면 이 메시지가 옵니다.
sw.addEventListener('message', (event) => {
	if (event.data && event.data.type === 'SKIP_WAITING') {
		console.log('[SW] SKIP_WAITING');
		sw.skipWaiting();
	}
});
sw.addEventListener('activate', (event) => event.waitUntil(sw.clients.claim()));

// 4. 동적 캐싱 (기존 fetch 이벤트를 Workbox 식으로 변경)
// 직접 fetch 리스너를 짜는 것보다 Workbox 전략을 쓰는 게 훨씬 안전합니다.
registerRoute(
	({ request, url }) => request.method === 'GET' && !url.pathname.startsWith('/api/sse'), //  API SSE 요청은 워커가 가로채지 않음
	// ({ request }) => request.method === 'GET',
	new NetworkFirst({
		cacheName: 'dynamic-cache',
	}),
);

sw.addEventListener('push', (event) => {
	if (!event.data) return;

	// 기본 데이터 설정 (이전 답변에서 누락된 url, tag, data 등을 추가하여 견고하게 만듦)
	let data = {
		title: '새로운 알림',
		body: '내용이 없습니다.',
		icon: '/logo/icon-192x192.png',
		badge: '/logo/badge-96x96.png',
		url: '/app', // 기본 URL
		tag: 'default-tag',
		data: {},
	};

	// 데이터 파싱 (JSON 파싱 실패 대비)
	try {
		const jsonData = event.data.json();
		data = { ...data, ...jsonData };
		console.log('[SW] Push Data:', data);
	} catch (e) {
		console.log('Push 데이터가 JSON이 아닙니다. 텍스트로 처리합니다.', e);
		data.body = event.data.text();
	}

	const options: NotificationOptions = {
		body: data.body,
		icon: data.icon,
		badge: data.badge,
		data: { url: data.url, custom: data.data || {} }, // data 필드에 url 포함
		tag: `push-${Date.now()}`,
	};

	event.waitUntil(sw.registration.showNotification(data.title, options));
});

sw.addEventListener('notificationclick', (event) => {
	event.notification.close();
	const targetUrl = event.notification.data.url;

	event.waitUntil(
		sw.clients.matchAll({ type: 'window', includeUncontrolled: true }).then((clientList) => {
			for (const client of clientList) {
				if (client.url.includes(sw.registration.scope) && 'focus' in client) {
					return client.focus().then(() => client.navigate(targetUrl));
				}
			}
			if (sw.clients.openWindow) return sw.clients.openWindow(targetUrl);
		}),
	);
});
