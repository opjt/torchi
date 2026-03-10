<script lang="ts">
	import '$src/app.css';
	import { goto } from '$app/navigation';
	import { Plus, Trash2, Power, Copy, Settings, Activity, ChevronLeft } from 'lucide-svelte';
	import { slide } from 'svelte/transition';

	// UI 상태 관리를 위한 임시 데이터 구조
	type Service = {
		id: string;
		name: string;
		token: string;
		active: boolean;
		lastActive?: string;
	};

	// Svelte 5 $state
	let services = $state<Service[]>([
		{
			id: '1',
			name: 'Github Actions',
			token: 'pook_gh_v8a9s8',
			active: true,
			lastActive: '2m ago'
		},
		{
			id: '2',
			name: 'Server Monitor',
			token: 'pook_srv_b7x2c1',
			active: false,
			lastActive: '1d ago'
		}
	]);

	let isAdding = $state(false);
	let newServiceName = $state('');
	let copiedId: string | null = $state(null);

	// 활성화된 서비스 개수 계산 (Derived state)
	let activeCount = $derived(services.filter((s) => s.active).length);

	// 서비스 추가
	function addService() {
		if (!newServiceName.trim()) return;

		const newService = {
			id: crypto.randomUUID(),
			name: newServiceName,
			token: `pook_${Math.random().toString(36).substr(2, 9)}`,
			active: true,
			lastActive: 'Just now'
		};

		services = [newService, ...services]; // 최신순 정렬
		newServiceName = '';
		isAdding = false;
	}

	// 서비스 삭제
	function deleteService(id: string) {
		if (!confirm('정말 이 서비스를 삭제하시겠습니까?')) return;
		services = services.filter((s) => s.id !== id);
	}

	// 서비스 토글
	function toggleServiceActive(id: string) {
		const idx = services.findIndex((s) => s.id === id);
		if (idx !== -1) {
			services[idx].active = !services[idx].active;
		}
	}

	// 엔드포인트 복사
	async function copyEndpoint(token: string, id: string) {
		const url = `https://pook.io/api/push/${token}`;
		const curlCmd = `curl -X POST "${url}" -d "msg=Hello!"`;

		await navigator.clipboard.writeText(curlCmd);
		copiedId = id;
		setTimeout(() => (copiedId = null), 2000);
	}
</script>

<div class="bg-base-100 font-sans text-base-content flex min-h-screen flex-col">
	<header
		class="top-0 bg-base-100/80 px-6 py-6 backdrop-blur-md sticky z-20 flex items-center justify-between"
	>
		<div class="flex items-center">
			<button
				onclick={() => goto('/app')}
				class="mr-2 -ml-2 p-2 opacity-50 transition-opacity hover:opacity-100"
				title="home"
			>
				<ChevronLeft />
			</button>
			<h1 class="text-xl font-black tracking-tight">My Endpoints</h1>
		</div>
	</header>
	<!-- <header
		class="px-6 py-6 top-0 bg-base-100/80 backdrop-blur-md sticky z-20 flex items-center justify-between"
	>
		<h1 class="text-xl font-black tracking-tight gap-2 flex items-center">
			My Services
			{#if activeCount > 0}
				<div class="badge badge-success badge-xs gap-1 font-mono opacity-80">
					<div class="w-1 h-1 bg-white animate-pulse rounded-full"></div>
					{activeCount} RUNNING
				</div>
			{/if}
		</h1>
		<button
			onclick={() => goto('/settings')}
			class="p-2 -mr-2 opacity-50 transition-all hover:rotate-90 hover:opacity-100 active:scale-90"
			title="Settings"
		>
			<Settings size={22} />
		</button>
	</header> -->

	<main class="px-6 pt-2 pb-20 flex-1 overflow-x-hidden">
		<div class="mb-6 flex items-center justify-between">
			<h2 class="font-bold text-[11px] tracking-[0.2em] uppercase opacity-40">Manage Endpoints</h2>
			<button
				onclick={() => (isAdding = !isAdding)}
				class="btn btn-circle btn-ghost btn-xs hover:bg-base-200 opacity-60 transition-all hover:opacity-100"
				class:rotate-45={isAdding}
				class:bg-base-200={isAdding}
				class:opacity-100={isAdding}
			>
				<Plus size={16} class="transition-transform duration-300" />
			</button>
		</div>

		<div class="space-y-4">
			{#if isAdding}
				<div
					transition:slide
					class="mb-4 rounded-3xl border-primary/20 bg-base-200/80 p-5 shadow-lg shadow-primary/5 border"
				>
					<p class="mb-3 ml-1 text-xs font-bold">Create New Service</p>
					<div class="gap-2 flex">
						<input
							type="text"
							bind:value={newServiceName}
							placeholder="ex) 결제 서버, 비트코인 알림"
							class="input-bordered input input-sm h-10 rounded-2xl bg-base-100 focus:border-primary w-full focus:outline-none"
							onkeydown={(e) => e.key === 'Enter' && addService()}
						/>
						<button
							onclick={addService}
							class="btn h-10 rounded-2xl px-4 font-bold btn-soft btn-sm"
						>
							등록
						</button>
					</div>
				</div>
			{/if}

			{#if services.length === 0}
				<div
					class="gap-3 rounded-3xl border-base-content/10 bg-base-200/30 py-12 text-xs flex flex-col items-center border border-dashed text-center opacity-40"
				>
					<Activity size={24} strokeWidth={1.5} />
					<p>
						등록된 서비스가 없습니다.<br />
						우측 상단 + 버튼을 눌러 시작하세요.
					</p>
				</div>
			{/if}

			{#each services as service (service.id)}
				<div
					class="group rounded-3xl border-base-content/5 bg-base-200/40 p-5 hover:bg-base-200/70 hover:shadow-lg hover:shadow-base-content/5 relative overflow-hidden border transition-all duration-300"
				>
					<div class="mb-4 flex items-center justify-between">
						<div class="gap-3 flex items-center">
							<div class="relative">
								<div
									class="h-2.5 w-2.5 rounded-full transition-all duration-500 {service.active
										? 'bg-success shadow-[0_0_10px_rgba(34,197,94,0.6)]'
										: 'bg-base-content/20'}"
								></div>
								{#if service.active}
									<div
										class="inset-0 h-2.5 w-2.5 animate-ping bg-success absolute rounded-full opacity-20"
									></div>
								{/if}
							</div>

							<div>
								<h3
									class="leading-tight font-bold text-[15px] {service.active ? '' : 'opacity-40'}"
								>
									{service.name}
								</h3>
								{#if service.active && service.lastActive}
									<p class="mt-0.5 font-mono text-[9px] opacity-40">
										Last active: {service.lastActive}
									</p>
								{/if}
							</div>
						</div>

						<div class="gap-1 flex items-center">
							<button
								onclick={() => toggleServiceActive(service.id)}
								class="btn btn-square rounded-xl btn-ghost btn-sm {service.active
									? 'text-success hover:bg-success/10'
									: 'text-base-content/30 hover:bg-base-content/10'}"
								title={service.active ? 'Pause Service' : 'Resume Service'}
							>
								<Power size={16} strokeWidth={2.5} />
							</button>
							<button
								onclick={() => deleteService(service.id)}
								class="btn btn-square rounded-xl text-error/40 btn-ghost btn-sm hover:bg-error/10 hover:text-error"
								title="Delete Service"
							>
								<Trash2 size={16} />
							</button>
						</div>
					</div>

					<div class="group/input relative">
						<div class="inset-y-0 left-3 pointer-events-none absolute flex items-center opacity-30">
							<span class="font-bold text-[10px]">POST</span>
						</div>
						<input
							type="text"
							readonly
							value="https://pook.io/api/push/{service.token}"
							class="rounded-2xl border-base-content/5 bg-base-100 py-3 pr-12 pl-12 font-mono focus:border-primary/30 focus:bg-base-100/50 w-full truncate border text-[10px] opacity-60 transition-all focus:opacity-100 focus:outline-none"
						/>
						<button
							onclick={() => copyEndpoint(service.token, service.id)}
							class="top-1.5 right-1.5 bottom-1.5 rounded-xl text-base-content/40 hover:bg-primary/10 hover:text-primary absolute flex aspect-square items-center justify-center transition-colors"
						>
							{#if copiedId === service.id}
								<span class="animate-bounce font-bold text-success text-[10px]">OK</span>
							{:else}
								<Copy size={14} />
							{/if}
						</button>
					</div>
				</div>
			{/each}
		</div>
	</main>
</div>
