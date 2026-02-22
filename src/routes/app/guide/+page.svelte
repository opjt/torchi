<script lang="ts">
	import { page } from '$app/stores';
	import { Copy, Check, Zap, MessageSquare, Terminal } from 'lucide-svelte';

	// 엔드포인트 토큰 (설정 페이지에서 넘어올 수 있음, 없으면 placeholder)
	let token = $derived($page.url.searchParams.get('token') ?? '{YOUR_TOKEN}');

	let copiedKey: string | null = $state(null);

	async function copy(text: string, key: string) {
		await navigator.clipboard.writeText(text);
		copiedKey = key;
		setTimeout(() => (copiedKey = null), 1500);
	}

	type Example = {
		key: string;
		label: string;
		description: string;
		command: string;
	};

	let examples: Example[] = $derived([
		{
			key: 'simple',
			label: '단순 알림',
			description: '메세지를 보내고 끝.',
			command: `curl "https://torchi.app/api/v1/push/${token}" \\\n  -d 'msg=Hello!'`,
		},
		{
			key: 'ask',
			label: '승인 요청',
			description: '모바일에서 응답할 때까지 대기.',
			command: `reaction=$(curl -s "https://torchi.app/api/v1/push/${token}/ask" \\\n  -d 'msg=프로덕션 배포할까요?' \\\n  -d 'actions=승인,거절')\n\nif [ "$reaction" = "승인" ]; then\n  ./deploy.sh\nfi`,
		},
		{
			key: 'timeout',
			label: '타임아웃 지정',
			description: '응답 대기 시간을 직접 지정.',
			command: `curl "https://torchi.app/api/v1/push/${token}/ask" \\\n  -d 'msg=DB 마이그레이션 실행할까요?' \\\n  -d 'actions=승인,거절' \\\n  -d 'timeout=60'`,
		},
		{
			key: 'pipe',
			label: '파이프 연결',
			description: '스트림으로 실시간 알림 수신.',
			command: `while read msg; do\n  notify-send "$msg"\ndone < <(stdbuf -i0 -o0 \\\n  curl -s "https://torchi.app/api/v1/push/${token}/raw")`,
		},
	]);

	let activeTab = $state('simple');
	let activeExample: Example = $derived(examples.find((e) => e.key === activeTab) ?? examples[0]);
</script>

<div class="bg-base-100 font-sans text-base-content flex min-h-screen flex-col">
	<header
		class="top-0 border-base-content/8 bg-base-100/80 px-5 py-4 backdrop-blur-md sticky z-20 border-b"
		style="height: 73px;"
	>
		<div class="flex h-full items-center justify-between">
			<h1 class="text-xl font-black tracking-tight gap-0.5 flex">
				Torchi<span class="text-primary">.</span>
				<span class="text-base-content/30 font-mono text-sm ml-2 font-normal mb-0.5 self-end"
					>docs</span
				>
			</h1>
			<a
				href="/app/setting"
				class="text-xs font-bold opacity-40 transition-opacity hover:opacity-100"
			>
				← 설정으로
			</a>
		</div>
	</header>

	<main class="px-5 pb-20 pt-6 space-y-10 max-w-2xl mx-auto w-full flex-1">
		<!-- 소개 -->
		<section>
			<p class="font-mono text-primary mb-2 tracking-widest text-[11px] uppercase">Quick Start</p>
			<h2 class="text-2xl font-black mb-3">두 가지 방법으로<br />알림을 보내세요.</h2>
			<p class="text-sm text-base-content/50 leading-relaxed">
				단순 알림부터 승인/거절 리액션까지.<br />
				curl 한 줄로 시작할 수 있어요.
			</p>
		</section>

		<!-- 탭 -->
		<section>
			<div class="gap-1 mb-4 border-base-content/8 flex overflow-x-auto border-b">
				{#each examples as ex}
					<button
						onclick={() => (activeTab = ex.key)}
						class="px-4 py-2.5 text-xs font-bold -mb-px border-b-2 whitespace-nowrap transition-all
							{activeTab === ex.key
							? 'border-primary text-primary'
							: 'text-base-content/40 hover:text-base-content/70 border-transparent'}"
					>
						{ex.label}
					</button>
				{/each}
			</div>

			<div class="rounded-2xl border-base-content/8 bg-base-200/50 overflow-hidden border">
				<!-- 설명 -->
				<div class="px-4 py-3 border-base-content/5 flex items-center justify-between border-b">
					<div class="gap-2 flex items-center">
						<Terminal size={13} class="text-primary opacity-70" />
						<span class="text-xs text-base-content/50">{activeExample.description}</span>
					</div>
					<button
						onclick={() => copy(activeExample.command.replace(/\\\n\s*/g, ' '), activeTab)}
						class="gap-1.5 font-bold flex items-center text-[11px] opacity-40 transition-opacity hover:opacity-100"
					>
						{#if copiedKey === activeTab}
							<Check size={12} class="text-success" />
							<span class="text-success">복사됨</span>
						{:else}
							<Copy size={12} />
							복사
						{/if}
					</button>
				</div>

				<!-- 코드 -->
				<pre
					class="px-4 py-4 font-mono leading-relaxed text-base-content/80 overflow-x-auto text-[12px] whitespace-pre">{activeExample.command}</pre>
			</div>
		</section>

		<!-- 응답 포맷 -->
		<section>
			<p class="font-mono text-primary mb-3 tracking-widest text-[11px] uppercase">Response</p>
			<div class="space-y-3">
				<div class="rounded-2xl border-base-content/8 bg-base-200/50 p-4 border">
					<div class="gap-2 mb-3 flex items-center">
						<Zap size={13} class="text-success" />
						<span class="text-xs font-bold">단순 알림 응답</span>
					</div>
					<pre class="font-mono text-base-content/70 text-[12px]">{'{ "sent": 1 }'}</pre>
				</div>

				<div class="rounded-2xl border-base-content/8 bg-base-200/50 p-4 border">
					<div class="gap-2 mb-3 flex items-center">
						<MessageSquare size={13} class="text-primary" />
						<span class="text-xs font-bold">ask 응답 (plain text)</span>
					</div>
					<pre
						class="font-mono text-base-content/70 text-[12px]">{'승인         # 모바일에서 선택한 값\n거절\ntimeout      # 시간 초과'}</pre>
				</div>
			</div>
		</section>

		<!-- 파라미터 -->
		<section>
			<p class="font-mono text-primary mb-3 tracking-widest text-[11px] uppercase">Parameters</p>
			<div class="rounded-2xl border-base-content/8 overflow-hidden border">
				{#each [{ param: 'msg', type: 'string', desc: '알림 메세지 본문', required: true }, { param: 'actions', type: 'string', desc: '리액션 선택지. 쉼표로 구분 (ask 전용)', required: false }, { param: 'timeout', type: 'number', desc: '응답 대기 시간 (초, 기본 300)', required: false }] as row, i}
					<div
						class="px-4 py-3 gap-4 flex items-start {i !== 0
							? 'border-base-content/5 border-t'
							: ''}"
					>
						<div class="w-24 shrink-0">
							<code class="font-mono text-primary text-[11px]">{row.param}</code>
							{#if row.required}
								<span class="ml-1 text-error font-bold text-[9px]">*</span>
							{/if}
						</div>
						<code class="font-mono text-base-content/30 w-14 shrink-0 text-[10px]">{row.type}</code>
						<span class="text-xs text-base-content/60">{row.desc}</span>
					</div>
				{/each}
			</div>
		</section>

		<!-- 엔드포인트 -->
		<section>
			<p class="font-mono text-primary mb-3 tracking-widest text-[11px] uppercase">Endpoints</p>
			<div class="rounded-2xl border-base-content/8 overflow-hidden border">
				{#each [{ method: 'POST', path: `/api/v1/push/${token}`, desc: '단순 알림 발송' }, { method: 'POST', path: `/api/v1/push/${token}/ask`, desc: '리액션 대기 (롱폴링)' }] as row, i}
					<div
						class="px-4 py-3 gap-3 flex items-center {i !== 0
							? 'border-base-content/5 border-t'
							: ''}"
					>
						<span
							class="font-mono font-bold text-success bg-success/10 px-2 py-0.5 rounded shrink-0 text-[10px]"
							>{row.method}</span
						>
						<code class="font-mono text-base-content/60 flex-1 truncate text-[11px]"
							>{row.path}</code
						>
						<span class="text-base-content/40 shrink-0 text-[11px]">{row.desc}</span>
					</div>
				{/each}
			</div>
		</section>
	</main>
</div>
