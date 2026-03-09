<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
	import { push } from '$lib/client/pushManager.svelte';
	import { auth, hasAgreedToTerms } from '$lib/client/auth/auth';
	import { onMount } from 'svelte';

	let { children } = $props();

	const { ready } = auth;
	let shouldShowContent = $state(false);

	onMount(async () => {
		await auth.init();
		await auth.whenReady();
	});

	// 경로나 인증 상태가 바뀔 때마다 체크
	$effect(() => {
		if (!$ready) {
			shouldShowContent = false;
			return;
		}

		const currentPath = $page.url.pathname;
		if (!auth.isAuthenticated()) {
			shouldShowContent = false;
			goto('/');
			return;
		}
		// [derived store] hasAgreedToTerms 에 의하여 상태변경을 감지하고 라우팅 수행.
		if (currentPath === '/app/welcome' && $hasAgreedToTerms) {
			shouldShowContent = false;
			goto('/app');
			return;
		}

		if (currentPath !== '/app/welcome' && !$hasAgreedToTerms) {
			shouldShowContent = false;
			goto('/app/welcome');
			return;
		}

		// 모든 체크 통과
		shouldShowContent = true;
	});
</script>

{#if !shouldShowContent}
	<div class="flex h-screen items-center justify-center">
		<span class="loading loading-spinner loading-lg"></span>
	</div>
{:else}
	<div class="flex h-screen flex-col">
		{#if push.permissionState !== null && push.permissionState !== 'granted'}
			<button
				type="button"
				onclick={() => push.handleSubscribe()}
				class="top-0 bg-warning text-warning-content px-6 py-3.5 sticky z-10 flex w-full items-center justify-between"
			>
				<div class="gap-3 flex items-center">
					<span class="h-2 w-2 relative flex">
						<span
							class="bg-warning-content animate-ping absolute inline-flex h-full w-full rounded-full opacity-75"
						></span>
						<span class="h-2 w-2 bg-warning-content relative inline-flex rounded-full"></span>
					</span>
					<p class="font-black tracking-tighter text-[11px] uppercase">
						{push.permissionState === 'denied'
							? '알림 권한이 차단됨 (브라우저 설정에서 해제 필요)'
							: '알림 권한이 필요합니다'}
					</p>
					{#if push.permissionState !== 'denied'}
						<span class="font-bold text-[10px] underline underline-offset-2">설정하기</span>
					{/if}
				</div>
			</button>
		{/if}

		<!-- <div class="flex-1 overflow-y-auto">
			<div class="max-w-3xl px-0 md:px-4 font-sans text-base-content mx-auto flex flex-col">
				{@render children()}
			</div>
		</div> -->
		<div class="flex-1">
			<div class="max-w-3xl px-0 md:px-4 mx-auto flex flex-col">
				{@render children()}
			</div>
		</div>
	</div>
{/if}
