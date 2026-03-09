<script lang="ts">
	import { agreeToTerms } from '$lib/api/user';
	import { auth } from '$lib/client/auth/auth';
	import { Check } from 'lucide-svelte';

	let agreed = $state(false);
	let isSubmitting = $state(false);

	async function handleSubmit() {
		if (!agreed || isSubmitting) return;

		isSubmitting = true;

		try {
			await agreeToTerms();
			await auth.init();
		} catch (e) {
			console.error(e);
		} finally {
			isSubmitting = false;
		}
	}
</script>

<div
	class="bg-base-100 font-sans text-base-content p-6 flex min-h-screen flex-col items-center justify-center"
>
	<main class="max-w-md w-full">
		<!-- 로고/타이틀 -->
		<div class="mb-8 text-center">
			<h1 class="text-4xl font-black tracking-tight gap-1 mb-5 flex items-center justify-center">
				Torchi<span class="text-primary">.</span>
			</h1>
			<h1 class="mb-1 text-xl font-black tracking-tight">환영합니다!</h1>
			<p class="text-sm opacity-60">서비스를 시작하기 전에 약관 동의가 필요합니다</p>
		</div>

		<!-- 약관 동의 카드 -->
		<div class="space-y-4">
			<!-- 수집 정보 안내 -->
			<div class="rounded-3xl border-base-content/5 bg-base-200/40 p-5 border">
				<h2 class="mb-3 text-xs font-bold tracking-[0.2em] uppercase opacity-40">수집하는 정보</h2>
				<div class="space-y-2 text-sm">
					<div class="gap-2 flex items-start">
						<div class="mt-0.5 h-1.5 w-1.5 bg-primary rounded-full"></div>
						<div>
							<span class="font-bold">항목:</span>
							<span class="opacity-70">이메일 주소</span>
						</div>
					</div>
					<div class="gap-2 flex items-start">
						<div class="mt-0.5 h-1.5 w-1.5 bg-primary rounded-full"></div>
						<div>
							<span class="font-bold">목적:</span>
							<span class="opacity-70">서비스 제공 및 사용자 식별</span>
						</div>
					</div>
					<div class="gap-2 flex items-start">
						<div class="mt-0.5 h-1.5 w-1.5 bg-primary rounded-full"></div>
						<div>
							<span class="font-bold">보유:</span>
							<span class="opacity-70">회원 탈퇴 시까지</span>
						</div>
					</div>
				</div>
			</div>

			<!-- 약관 동의 체크박스 -->
			<label
				class="group gap-4 rounded-3xl border-base-content/5 bg-base-200/40 p-5 hover:bg-base-200/70 flex cursor-pointer items-start border transition-all {agreed
					? 'border-primary/30 bg-primary/5'
					: ''}"
			>
				<div class="relative shrink-0">
					<input
						type="checkbox"
						bind:checked={agreed}
						class="peer h-6 w-6 rounded-lg border-base-content/20 bg-base-100 checked:border-primary checked:bg-primary focus:ring-primary/30 cursor-pointer appearance-none border-2 transition-all focus:ring-2 focus:outline-none"
					/>
					<div
						class="inset-0 text-primary-content pointer-events-none absolute flex items-center justify-center opacity-0 transition-opacity peer-checked:opacity-100"
					>
						<Check size={16} strokeWidth={3} />
					</div>
				</div>
				<div class="pt-0.5 flex-1">
					<p class="mb-1 text-sm font-bold">
						<span class="text-error">*</span> 필수 약관에 모두 동의합니다
					</p>
					<p class="mb-2 text-xs opacity-60">아래 약관을 확인하고 동의해주세요</p>
					<div class="gap-2 flex flex-wrap">
						<a
							href="/terms"
							target="_blank"
							class="text-xs font-bold text-primary decoration-primary/30 hover:decoration-primary underline underline-offset-2 transition-colors"
							onclick={(e) => e.stopPropagation()}
						>
							이용약관
						</a>
						<span class="text-xs opacity-30">·</span>
						<a
							href="/privacy"
							target="_blank"
							class="text-xs font-bold text-primary decoration-primary/30 hover:decoration-primary underline underline-offset-2 transition-colors"
							onclick={(e) => e.stopPropagation()}
						>
							개인정보처리방침
						</a>
					</div>
				</div>
			</label>

			<!-- 시작하기 버튼 -->
			<button
				onclick={handleSubmit}
				disabled={!agreed || isSubmitting}
				class="btn h-14 rounded-2xl font-bold w-full border-2 transition-all {agreed
					? 'btn-primary hover:scale-[1.02]'
					: 'btn-disabled opacity-40'}"
			>
				{#if isSubmitting}
					<span class="loading loading-spinner loading-sm"></span>
					처리중...
				{:else}
					동의하고 시작하기
				{/if}
			</button>

			<!-- 거부 안내 -->
			<p class="pt-2 text-xs text-center opacity-40">
				동의하지 않으면 서비스를 이용하실 수 없습니다
			</p>
		</div>
	</main>
</div>
