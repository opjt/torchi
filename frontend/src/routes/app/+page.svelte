<script lang="ts">
	import { goto } from '$app/navigation';
	import { fetchEndpoints, type Endpoint } from '$lib/api/endpoints';
	import {
		deleteNotification,
		getNotifications,
		markAsReadUntil,
		postReaction,
		transformNotification,
		type DisplayNotification,
	} from '$lib/api/notifications';
	import { auth } from '$lib/client/auth/auth';
	import { debugLog, linkify } from '$lib/pkg/util';
	import { BellOff, ChevronDown, ChevronLeft, Search, Settings, X } from 'lucide-svelte';
	import { onMount, tick } from 'svelte';
	import { cubicOut } from 'svelte/easing';
	import { slide } from 'svelte/transition';

	// --- 상태 관리 ---
	let notifications = $state<DisplayNotification[]>([]);
	let nextCursor = $state<string | null>(null);
	let hasMore = $state(false);
	let loading = $state(false);
	let isFilterOpen = $state(false);
	let selectedServiceId = $state<string | 'ALL'>('ALL');

	// --- 검색 상태 ---
	let searchInputValue = $state(''); // input에 바인딩되는 값 (타이핑 중인 값)
	let searchQuery = $state(''); // 실제 API 검색에 사용되는 확정 값
	let isSearchMode = $state(false);
	let searchInputEl = $state<HTMLInputElement | null>(null);

	// 관찰할 하단 요소
	let observerTarget = $state<HTMLElement | null>(null);

	// 엔드포인트 목록
	let endpoints = $state<Endpoint[]>([]);

	function connectSSE(): EventSource {
		const es = new EventSource('/api/sse/notifications', {
			withCredentials: true,
		});

		es.addEventListener('notification', async (e) => {
			const newNoti = transformNotification(JSON.parse(e.data));

			// 이미 있는 알림이면 스킵 (중복 방지)
			if (notifications.some((n) => n.id === newNoti.id)) return;

			// 검색 중이거나 필터가 다르면 스킵
			if (isSearchMode) return;
			if (selectedServiceId !== 'ALL' && newNoti.id !== selectedServiceId) return;

			notifications = [newNoti, ...notifications];
			await tick();

			await markAsReadUntil(newNoti.id, selectedServiceId);
		});

		es.addEventListener('connected', () => {
			debugLog('SSE connected');
		});

		es.onerror = () => {
			// 브라우저가 자동 재연결하므로 별도 처리 불필요
		};
		return es;
	}

	// 선택된 서비스의 "이름"을 찾기 위한 derived
	let currentFilterName = $derived.by(() => {
		if (selectedServiceId === 'ALL') return '모든 서비스';
		return endpoints.find((e) => e.id === selectedServiceId)?.name || '알 수 없는 서비스';
	});

	async function loadEndpoints() {
		try {
			endpoints = await fetchEndpoints();
		} catch (e) {
			console.error('Failed to fetch endpoints:', e);
		}
	}

	async function loadNotifications(isFirst = false) {
		if (loading || (!isFirst && !hasMore)) {
			debugLog('Skip loading', { loading, hasMore, isFirst });
			return;
		}

		loading = true;
		debugLog('Start loading', { isFirst, nextCursor, searchQuery });
		try {
			const res = await getNotifications(
				isFirst ? undefined : (nextCursor ?? undefined),
				selectedServiceId,
				searchQuery.trim() || undefined,
			);
			const newItems = res.items.map(transformNotification);

			if (isFirst) {
				notifications = newItems;
				await tick();

				window.scrollTo({ top: 0, behavior: 'instant' });
			} else {
				notifications = [...notifications, ...newItems];
			}

			nextCursor = res.next_cursor;
			hasMore = res.has_more;

			if (newItems.length > 0) {
				const lastIdOfBatch = newItems[newItems.length - 1].id;
				await markAsReadUntil(lastIdOfBatch, selectedServiceId);
			}
			debugLog('Load success', { itemsCount: newItems.length, hasMore });
		} catch (e) {
			console.error('Failed:', e);
		} finally {
			loading = false;
		}
	}

	// 무한스크롤 IntersectionObserver
	$effect(() => {
		if (!observerTarget) {
			debugLog('Target element not found yet');
			return;
		}

		const observer = new IntersectionObserver(
			(entries) => {
				const entry = entries[0];
				if (entry.isIntersecting && hasMore && !loading) {
					loadNotifications();
				}
			},
			{ threshold: 0.1, rootMargin: '100px' },
		);

		observer.observe(observerTarget);
		return () => observer.disconnect();
	});

	onMount(async () => {
		await Promise.all([loadNotifications(true), loadEndpoints()]);
	});

	$effect(() => {
		let es: EventSource | null = null;

		// SSE 연결 함수 호출
		debugLog('SSE try to connect...');
		es = connectSSE();

		// 이 리턴 함수가 컴포넌트가 파괴될 때 호출됩니다.
		return () => {
			if (es) {
				debugLog('SSE 연결 종료 (Cleanup)');
				es.close();
			}
		};
	});

	function toggleFilter() {
		isFilterOpen = !isFilterOpen;
	}

	function selectFilter(id: string | 'ALL') {
		selectedServiceId = id;
		isFilterOpen = false;
		nextCursor = null;
		hasMore = true;
		loadNotifications(true);
	}

	async function handleDelete(id: string) {
		try {
			await deleteNotification(id);
			notifications = notifications.filter((n) => n.id !== id);
		} catch (e) {
			console.error(e);
		}
	}

	async function enterSearchMode() {
		isSearchMode = true;
		notifications = []; // 목록 즉시 날림
		hasMore = false;
		await tick();
		searchInputEl?.focus();
	}

	function exitSearchMode() {
		isSearchMode = false;
		searchInputValue = '';
		searchQuery = '';
		nextCursor = null;
		hasMore = true;
		loadNotifications(true); // 일반 목록 복구
	}

	function clearSearch() {
		searchInputValue = '';
		searchInputEl?.focus();
	}

	// 엔터 시 검색 실행
	function handleSearchKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			searchQuery = searchInputValue; // 확정값에 반영
			nextCursor = null;
			hasMore = true;
			loadNotifications(true);
		}
	}

	function highlight(html: string, query: string): string {
		if (!query.trim()) return html;
		const escaped = query.replace(/[.*+?^${}()|[\]\\]/g, '\\$&');
		const regex = new RegExp(escaped, 'gi');
		const MERGE_GAP = 2; // 이 글자 수 이하 간격이면 합침

		return html.replace(/(<[^>]*>)|([^<]+)/g, (_, tag, text) => {
			if (tag) return tag;

			// 1. 모든 매치 위치 수집
			const matches: { start: number; end: number }[] = [];
			let m: RegExpExecArray | null;
			while ((m = regex.exec(text)) !== null) {
				matches.push({ start: m.index, end: m.index + m[0].length });
			}
			if (matches.length === 0) return text;

			// 2. 근접한 매치 병합
			const merged = [matches[0]];
			for (let i = 1; i < matches.length; i++) {
				const prev = merged[merged.length - 1];
				const curr = matches[i];
				if (curr.start - prev.end <= MERGE_GAP) {
					prev.end = curr.end; // 병합
				} else {
					merged.push(curr);
				}
			}

			// 3. 병합된 구간으로 하이라이트 재구성
			let result = '';
			let cursor = 0;
			for (const { start, end } of merged) {
				result += text.slice(cursor, start);
				result += `<mark class="bg-primary/30 text-inherit rounded-sm px-0.5">${text.slice(start, end)}</mark>`;
				cursor = end;
			}
			result += text.slice(cursor);
			return result;
		});
	}
	function slideX(node: Element, { duration = 200, delay = 0, from = 1 } = {}) {
		return {
			duration,
			delay,
			css: (t: number) => {
				const e = cubicOut(t);
				return `
        opacity: ${e};
        transform: translateX(${(1 - e) * from * 24}px);
      `;
			},
		};
	}
	async function handleReaction(id: string, action: string) {
		await postReaction(id, action);
		notifications = notifications.map((n) => (n.id === id ? { ...n, reaction: action } : n));
	}
</script>

<div class="bg-base-100 font-sans text-base-content flex min-h-screen flex-col">
	<header
		class="top-0 border-base-content/10 bg-base-100/90 backdrop-blur-md sticky z-30 overflow-hidden border-b"
		style="height: 73px;"
	>
		{#if !isSearchMode}
			<div
				class="px-5 flex h-full items-center justify-between"
				in:slideX={{ duration: 200, delay: 80, from: 1 }}
				out:slideX={{ duration: 120, from: 1 }}
			>
				<div>
					<h1 class="text-xl font-black tracking-tight gap-0.5 flex">
						Torchi<span class="text-primary">.</span>
					</h1>
					<p class="font-mono text-[10px] opacity-40">{$auth?.email || 'Guest'}</p>
				</div>
				<div class="gap-1.5 flex items-center">
					<button
						title="검색"
						onclick={enterSearchMode}
						class="btn btn-square rounded-xl btn-ghost btn-sm opacity-40 transition-all hover:opacity-100"
					>
						<Search />
					</button>
					<button
						title="설정"
						onclick={() => goto('/app/setting')}
						class="btn btn-square rounded-xl btn-ghost btn-sm opacity-40 transition-all hover:opacity-100"
					>
						<Settings />
					</button>
				</div>
			</div>
		{:else}
			<div
				class="px-2 gap-2 flex h-full items-center"
				in:slideX={{ duration: 200, delay: 80, from: -1 }}
				out:slideX={{ duration: 120, from: -1 }}
			>
				<button
					onclick={exitSearchMode}
					class="btn btn-square rounded-xl btn-ghost btn-sm shrink-0 opacity-50 transition-all hover:opacity-100 active:scale-90"
					title="뒤로"
				>
					<ChevronLeft size={20} />
				</button>
				<div
					class="border-base-content/10 bg-base-content/5 focus-within:border-primary/40 focus-within:bg-base-100 gap-2 rounded-xl px-3 py-2 min-w-0 flex flex-1 items-center border transition-all duration-200"
				>
					<Search size={13} class="shrink-0 opacity-30" />

					<!-- scope chip (ALL이 아닐 때만 표시) -->
					{#if selectedServiceId !== 'ALL'}
						<span
							class="bg-primary/15 text-primary rounded-md px-2 py-0.5 font-black max-w-20 shrink-0 truncate text-[11px]"
						>
							{currentFilterName}
						</span>
					{/if}

					<input
						bind:this={searchInputEl}
						bind:value={searchInputValue}
						onkeydown={handleSearchKeydown}
						type="text"
						placeholder={selectedServiceId === 'ALL' ? '알림 검색 후 엔터...' : '검색 후 엔터...'}
						class="placeholder:text-base-content/25 min-w-0 text-sm font-medium flex-1 bg-transparent outline-none"
					/>

					{#if searchInputValue}
						<button
							onclick={clearSearch}
							class="p-0.5 bg-base-300 shrink-0 rounded-full opacity-30 transition-all hover:opacity-70 active:scale-90"
						>
							<X size={13} />
						</button>
					{/if}
				</div>
				<button
					onclick={exitSearchMode}
					class="text-xs p-1 font-bold shrink-0 opacity-80 transition-all hover:opacity-100 active:scale-90"
				>
					닫기
				</button>
			</div>
		{/if}
	</header>

	<main class="px-4 relative flex-1">
		{#if isSearchMode}
			<div class="h-4"></div>
		{:else}
			<div class="py-2 top-15 sticky z-20">
				<div class="gap-2 flex items-center">
					<div class="relative inline-block">
						<button
							onclick={toggleFilter}
							class="btn h-10 gap-2 rounded-sm border-base-content/10 bg-base-100 pr-3 pl-4 shadow-xs btn-sm hover:border-primary hover:bg-base-100 flex items-center border transition-all"
						>
							<span class="text-xs font-bold opacity-60">Filter:</span>
							<span class="text-xs font-black text-primary max-w-30 truncate"
								>{currentFilterName}</span
							>

							<ChevronDown
								size={14}
								class="opacity-40 transition-transform {isFilterOpen ? 'rotate-180' : ''}"
							/>
						</button>
						{#if isFilterOpen}
							<div
								transition:slide={{ duration: 150 }}
								class="top-12 left-0 w-56 gap-1 rounded-sm border-white/10 bg-base-100/95 p-1 shadow-xl backdrop-blur-xl absolute flex flex-col overflow-hidden border"
							>
								<button
									onclick={() => selectFilter('ALL')}
									class="rounded-sm px-4 py-3 text-xs font-bold hover:bg-white/5 flex items-center justify-between text-left transition-colors {selectedServiceId ===
									'ALL'
										? 'bg-primary/5 text-primary'
										: 'opacity-60'}"
								>
									모든 서비스
									{#if selectedServiceId === 'ALL'}
										<span class="h-1.5 w-1.5 bg-primary rounded-full"></span>
									{/if}
								</button>
								<div class="mx-2 my-1 bg-white/5 h-px"></div>
								{#each endpoints as enp}
									<button
										onclick={() => selectFilter(enp.id)}
										class="rounded-sm px-4 py-3 text-xs font-bold hover:bg-base-content/5 flex items-center justify-between text-left transition-colors {selectedServiceId ===
										enp.id
											? 'bg-primary/5 text-primary'
											: 'opacity-60'}"
									>
										{enp.name}
										{#if selectedServiceId === enp.id}
											<span class="h-1.5 w-1.5 bg-primary rounded-full"></span>
										{/if}
									</button>
								{/each}
							</div>
							<button
								title="close"
								class="inset-0 fixed z-[-1]"
								onclick={() => (isFilterOpen = false)}
							></button>
						{/if}
					</div>
				</div>
			</div>
		{/if}

		<div class="space-y-3">
			{#each notifications as noti (noti.id)}
				<div transition:slide={{ duration: 200, axis: 'y' }}>
					<div
						class="group/card rounded-xl px-5 py-4 relative border transition-all
                    {isSearchMode
							? 'bg-base hover:bg-base-content/3 border-base-content/10'
							: noti.isRead
								? 'bg-base-content/2 hover:bg-base-content/3 border-transparent'
								: 'bg-primary/12 hover:border-primary/30 hover:bg-primary/9 shadow-sm border-base-content/10'}"
					>
						<div class="mb-2 flex items-center justify-between">
							<div class="gap-2 flex items-center">
								<div
									class="h-1.5 w-1.5 rounded-full {noti.isRead
										? 'bg-base-content/15'
										: 'animate-pulse bg-primary shadow-sm shadow-primary/50'}"
								></div>
								<span class="font-bold text-xs">
									{@html highlight(noti.endpointName ?? '', searchQuery)}
								</span>
							</div>

							<div class="gap-3 flex items-center">
								{#if noti.isMute}
									<span class="text-xs opacity-30">
										<BellOff size={14} />
									</span>
								{/if}
								<span class="font-mono text-xs opacity-35">{noti.timestamp}</span>
								<button
									onclick={() => handleDelete(noti.id)}
									class="p-1.5 transition-all hover:opacity-100 active:scale-90 {noti.isRead
										? 'hover:text-base-content opacity-30'
										: 'hover:text-primary opacity-50'}"
									title="삭제"
								>
									<X size={18} strokeWidth={2.5} />
								</button>
							</div>
						</div>

						<p
							class="pl-3.5 leading-relaxed font-medium border-l text-[14px] whitespace-pre-wrap
                    {noti.isRead
								? 'border-base-content/10 text-base-content/75'
								: 'border-primary/30 text-base-content/95'}"
						>
							{@html highlight(linkify(noti.body), searchQuery)}
						</p>
						<!-- 리액션 버튼 -->
						{#if noti.actions && noti.actions.length > 0}
							<div class="mt-3 ml-1.5 gap-2 flex">
								{#if noti.reaction}
									<span class="text-xs font-bold font-mono opacity-40">
										✓ {noti.reaction}
									</span>
								{:else if noti.isExpired}
									<span class="text-xs font-bold font-mono opacity-30"> ⏱ 시간 초과 </span>
								{:else}
									{#each noti.actions as action, i}
										<button
											onclick={() => handleReaction(noti.id, action)}
											class="px-3 py-1.5 rounded-lg text-xs font-bold transition-all active:scale-95
                    {i === 0
												? 'bg-primary/15 text-primary hover:bg-primary/25 border-primary/20 hover:border-primary/40 border'
												: 'bg-base-content/5 text-base-content/60 hover:bg-base-content/10 border-base-content/10 hover:border-base-content/20 border'}"
										>
											{action}
										</button>
									{/each}
								{/if}
							</div>
						{/if}
					</div>
				</div>
			{:else}
				<div class="py-24 text-center opacity-30 flex flex-col items-center">
					{#if loading}
						<!-- 로딩 중엔 빈 상태 메시지 숨김 -->
					{:else if isSearchMode}
						<Search size={20} class="mb-3 opacity-40" />
						<p class="text-xs font-mono">검색어 입력 후 엔터</p>
					{:else}
						<p class="text-xs font-mono">NO_NOTIFICATIONS</p>
					{/if}
				</div>
			{/each}

			<div class="py-4 flex items-center justify-center" bind:this={observerTarget}>
				{#if loading}
					<span class="loading loading-spinner loading-lg"></span>
				{/if}
			</div>
		</div>
	</main>
</div>
