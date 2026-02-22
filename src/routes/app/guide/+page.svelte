<script lang="ts">
	import { page } from '$app/stores';
	import { Copy, Check, Clock, Terminal } from 'lucide-svelte';

	let token = $derived($page.url.searchParams.get('token') ?? '{endpoint}');
	let copiedKey: string | null = $state(null);
	let activeTab = $state<'push' | 'ask'>('push');
	let activeAskTab = $state<'basic' | 'script' | 'timeout'>('basic');

	async function copy(text: string, key: string) {
		await navigator.clipboard.writeText(text);
		copiedKey = key;
		setTimeout(() => (copiedKey = null), 1500);
	}

	type AskExample = {
		key: 'basic' | 'script' | 'timeout';
		label: string;
		description: string;
		command: string;
	};

	let askExamples: AskExample[] = $derived([
		{
			key: 'basic',
			label: '기본',
			description: '알림을 보내고 응답 대기.',
			command: `curl -s "https://torchi.app/api/v1/push/${token}/ask" \\\n  -d 'msg=배포할까요?' \\\n  -d 'actions=승인,거절'`,
		},
		{
			key: 'script',
			label: '스크립트',
			description: '응답값으로 후속 작업 분기.',
			command: `reaction=$(curl -s "https://torchi.app/api/v1/push/${token}/ask" \\\n  -d 'msg=프로덕션 배포할까요?' \\\n  -d 'actions=승인,거절')\n\nif [ "$reaction" = "승인" ]; then\n  ./deploy.sh\nfi`,
		},
		{
			key: 'timeout',
			label: '타임아웃',
			description: '응답 대기 시간을 직접 지정.',
			command: `curl -s "https://torchi.app/api/v1/push/${token}/ask" \\\n  -d 'msg=DB 마이그레이션 실행할까요?' \\\n  -d 'actions=승인,거절' \\\n  -d 'timeout=60'`,
		},
	]);

	let activeAskExample: AskExample = $derived(
		askExamples.find((e) => e.key === activeAskTab) ?? askExamples[0],
	);
</script>

<div class="bg-base-100 font-sans text-base-content flex min-h-screen flex-col">
	<header
		class="top-0 border-base-content/8 bg-base-100/80 px-5 backdrop-blur-md sticky z-20 flex items-center justify-between border-b"
		style="height: 73px;"
	>
		<h1 class="text-xl font-black tracking-tight gap-0.5 flex">
			Torchi<span class="text-primary">.</span>
			<span class="text-base-content/30 font-mono text-sm ml-2 font-normal mb-0.5 self-end"
				>guide</span
			>
		</h1>
		<a
			href="/app/setting"
			class="text-xs font-bold opacity-40 transition-opacity hover:opacity-100"
		>
			← 설정으로
		</a>
	</header>

	<main class="px-5 pb-20 pt-8 max-w-2xl space-y-10 mx-auto w-full flex-1">
		<!-- 소개 -->
		<section>
			<p class="font-mono text-primary mb-2 tracking-widest text-[11px] uppercase">Quick Start</p>
			<h2 class="text-2xl font-black mb-3">보내고 끝내거나, 답을 기다리거나.</h2>
			<p class="text-sm text-base-content/50 leading-relaxed">
				단순 알림부터 승인/거절 리액션까지.<br />
				curl 한 줄로 시작할 수 있어요.
			</p>
		</section>

		<!-- 토큰 -->
		<!-- <section>
			<p class="font-mono text-base-content/30 tracking-widest mb-2 text-[10px] uppercase">
				Your Token
			</p>
			<div
				class="font-mono text-xs bg-base-200/60 border-base-content/8 rounded-xl px-4 py-3 text-base-content/60 gap-2 flex items-center justify-between border"
			>
				<span class="truncate">{token}</span>
				<button
					onclick={() => copy(token, 'token')}
					class="shrink-0 opacity-50 transition-opacity hover:opacity-100"
				>
					{#if copiedKey === 'token'}
						<Check size={12} class="text-success" />
					{:else}
						<Copy size={12} />
					{/if}
				</button>
			</div>
		</section> -->

		<!-- 메인 탭 -->
		<section>
			<div class="gap-1 mb-6 border-base-content/8 flex border-b">
				<button
					onclick={() => (activeTab = 'push')}
					class="px-4 py-2.5 text-xs font-bold -mb-px border-b-2 whitespace-nowrap transition-all
						{activeTab === 'push'
						? 'border-primary text-primary'
						: 'text-base-content/40 hover:text-base-content/70 border-transparent'}"
				>
					단순 알림
				</button>
				<button
					onclick={() => (activeTab = 'ask')}
					class="px-4 py-2.5 text-xs font-bold -mb-px border-b-2 whitespace-nowrap transition-all
						{activeTab === 'ask'
						? 'border-primary text-primary'
						: 'text-base-content/40 hover:text-base-content/70 border-transparent'}"
				>
					리액션 알림
				</button>
			</div>

			{#if activeTab === 'push'}
				<div class="space-y-4">
					<div>
						<div class="gap-3 mb-4 bg-base-300/50 rounded p-1.5 flex items-center">
							<span
								class="font-mono font-bold bg-primary/15 text-primary px-2 py-1 rounded text-[10px]"
								>POST</span
							>
							<code class="font-mono text-xs text-base-content/50 truncate"
								>/api/v1/push/{token}</code
							>
						</div>
						<p class="text-sm text-base-content/50 leading-relaxed">
							메세지를 보내고 끝. 응답을 기다리지 않아요.<br />
							<code class="font-mono text-xs">-d</code> 뒤에 오는 body 전체가 메세지로 전송됩니다.
						</p>
					</div>

					<div class="rounded-2xl border-base-content/8 overflow-hidden border">
						<div class="px-4 py-3 border-base-content/5 flex items-center justify-between border-b">
							<div class="gap-2 flex items-center">
								<Terminal size={13} class="text-success opacity-70" />
								<span class="text-xs text-base-content/50">메세지를 보내고 끝.</span>
							</div>
							<button
								onclick={() =>
									copy(
										`curl "https://torchi.app/api/v1/push/${token}" -d '안녕하세요!'`,
										'push-basic',
									)}
								class="gap-1.5 font-bold flex items-center text-[11px] opacity-40 transition-opacity hover:opacity-100"
							>
								{#if copiedKey === 'push-basic'}
									<Check size={12} class="text-success" /><span class="text-success">복사됨</span>
								{:else}
									<Copy size={12} />복사
								{/if}
							</button>
						</div>
						<pre
							class="px-4 py-3 bg-base-200 font-mono leading-relaxed text-base-content/80 overflow-x-auto text-[12px] whitespace-pre">curl "https://torchi.app/api/v1/push/{token}" \
  -d '안녕하세요!'</pre>
					</div>

					<!-- <div class="rounded-2xl border-base-content/8 bg-base-200/40 overflow-hidden border">
						<div class="px-4 py-2.5 border-base-content/5 border-b">
							<span class="font-mono text-base-content/40 text-[10px]">응답</span>
						</div>
						<pre
							class="px-4 py-4 font-mono text-base-content/75 text-[12px]">{'{ "sent": 1 }'}</pre>
					</div> -->
				</div>
			{:else}
				<div class="space-y-4">
					<div>
						<div class="gap-3 mb-4 bg-base-300/50 rounded p-1.5 flex items-center">
							<span
								class="font-mono font-bold bg-primary/15 text-primary px-2 py-1 rounded text-[10px]"
								>POST</span
							>
							<code class="font-mono text-xs text-base-content/50 truncate"
								>/api/v1/push/{token}/ask</code
							>
						</div>
						<p class="text-sm text-base-content/50 mb-2 leading-relaxed">
							알림을 보내고 모바일에서 응답할 때까지 커넥션을 유지해요.<br />
							응답이 오면 선택한 값이 <strong class="text-base-content/70">plain text</strong>로
							반환되고 종료됩니다.
						</p>
						<!-- <div class="gap-2 flex items-center">
							<Clock size={12} class="text-base-content/30" />
							<span class="text-base-content/40 text-[11px]"
								>기본 타임아웃 300초 · 초과 시 <code class="font-mono">timeout</code> 반환</span
							>
						</div> -->
					</div>

					<!-- 카드 하나로 통합 -->
					<div class="rounded-2xl border-base-content/8 overflow-hidden border">
						<!-- 탭 바 + 복사 -->
						<div
							class="px-4 pl-1.5 border-base-content/8 flex items-center justify-between border-b"
						>
							<div class="flex">
								{#each askExamples as ex}
									<button
										onclick={() => (activeAskTab = ex.key)}
										class="px-4 py-3 font-bold -mb-px border-b-2 text-[11px] transition-all
                        {activeAskTab === ex.key
											? 'border-primary text-primary'
											: 'text-base-content/35 hover:text-base-content/60 border-transparent'}"
									>
										{ex.label}
									</button>
								{/each}
							</div>
							<button
								onclick={() =>
									copy(activeAskExample.command.replace(/\\\n\s*/g, ' '), activeAskTab)}
								class="gap-1.5 font-bold flex items-center text-[11px] opacity-40 transition-opacity hover:opacity-100"
							>
								{#if copiedKey === activeAskTab}
									<Check size={12} class="text-success" /><span class="text-success">복사됨</span>
								{:else}
									<Copy size={12} />복사
								{/if}
							</button>
						</div>

						<!-- 코드 -->
						<pre
							class="bg-base-200 px-4 py-3 font-mono leading-relaxed text-base-content/80 overflow-x-auto text-[12px] whitespace-pre">{activeAskExample.command}</pre>
					</div>

					<!-- 응답 -->
					<!-- <div class="rounded-2xl border-base-content/8 bg-base-200/40 overflow-hidden border">
						<div class="px-4 py-2.5 border-base-content/5 border-b">
							<span class="font-mono text-base-content/40 text-[10px]">응답 (plain text)</span>
						</div>
						<pre
							class="px-4 py-4 font-mono leading-relaxed text-base-content/75 overflow-x-auto text-[12px] whitespace-pre">승인        # 모바일에서 선택한 값
거절
timeout     # 타임아웃 초과</pre>
					</div> -->

					<!-- 파라미터 -->
					<div>
						<p class="font-mono text-base-content/70 tracking-widest mb-3 text-[10px] uppercase">
							Parameters
						</p>
						<div class="rounded-2xl border-base-content/8 overflow-hidden border">
							{#each [{ param: 'msg', required: true, desc: '알림 메세지 본문' }, { param: 'actions', required: true, desc: '리액션 선택지. 쉼표로 구분 (예: 승인,거절)' }, { param: 'timeout', required: false, desc: '응답 대기 시간 (초, 기본 300)' }] as row, i}
								<div
									class="px-4 py-3 gap-4 flex items-center {i !== 0
										? 'border-base-content/5 border-t'
										: ''}"
								>
									<div class="w-20 gap-0.5 flex shrink-0 items-center">
										<code class="font-mono text-primary text-[11px]">{row.param}</code>
										{#if row.required}
											<span class="text-error font-bold text-[9px]">*</span>
										{/if}
									</div>
									<span class="text-xs text-base-content/50">{row.desc}</span>
								</div>
							{/each}
						</div>
					</div>
				</div>
			{/if}
		</section>
	</main>
</div>
