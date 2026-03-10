<script lang="ts">
	import { loginWithGithub } from '$lib/client/auth/github-auth';
	import { page } from '$app/stores';
	import { Bell, LoaderCircle, Play, Share, SquarePlus, X } from 'lucide-svelte';
	// 아이콘 추가
	import { onMount } from 'svelte';
	import { fly } from 'svelte/transition';
	// 애니메이션 추가
	import { guestLogin } from '$lib/client/auth/lifecycle';
	import { push } from '$lib/client/pushManager.svelte';

	// PWA 설치 관련 변수
	let deferredPrompt: any;
	let showInstallPopup = false;
	let isIOS = false;

	let isWebContext = true;

	const demoMessage = 'Hello world!';

	async function handleSubscribe() {
		await push.handleDemoPush(demoMessage);
	}

	// PWA 설치 함수
	async function installPwa() {
		if (!deferredPrompt) return;

		deferredPrompt.prompt();
		const { outcome } = await deferredPrompt.userChoice;

		if (outcome === 'accepted') {
			showInstallPopup = false;
			deferredPrompt = null;
		}
	}

	function closePopup() {
		hideInstallPopupForOneDay();
		showInstallPopup = false;
	}
	const INSTALL_POPUP_HIDE_KEY = 'hideInstallPopupUntil';
	const ONE_DAY = 1000 * 60 * 60 * 24;

	function shouldShowInstallPopup() {
		const hideUntil = localStorage.getItem(INSTALL_POPUP_HIDE_KEY);
		if (!hideUntil) return true;

		return Date.now() > Number(hideUntil);
	}

	function hideInstallPopupForOneDay() {
		const tomorrow = Date.now() + ONE_DAY;
		localStorage.setItem(INSTALL_POPUP_HIDE_KEY, String(tomorrow));
	}
	let origin = 'httptest'; // 기본값 (서버 에러 방지용)

	onMount(async () => {
		origin = window.location.origin;
		// 1. 현재 앱이 설치된 상태(Standalone)인지 확인
		const isStandalone = window.matchMedia('(display-mode: standalone)').matches;

		// 2. iOS 여부 확인 (User Agent 체크)
		const ua = window.navigator.userAgent.toLowerCase();
		isIOS =
			/iphone|ipad|ipod/.test(ua) ||
			!!(
				window.navigator.maxTouchPoints &&
				window.navigator.maxTouchPoints > 1 &&
				/macintosh|macintel/.test(ua)
			);

		// 3. iOS이고 아직 설치 안 했다면 -> 팝업 강제 표시 (iOS는 이벤트가 없으므로)
		if (isIOS && !isStandalone && shouldShowInstallPopup()) {
			// 약간의 딜레이를 주어 자연스럽게 등장
			setTimeout(() => {
				showInstallPopup = true;
			}, 800);
		}
		// PWA 설치 이벤트 리스너
		window.addEventListener('beforeinstallprompt', (e) => {
			console.log('성공: beforeinstallprompt 이벤트가 발생했습니다!');
			// 기본 브라우저 동작 방지 (자동으로 뜨는 미니 바 숨김)
			e.preventDefault();
			deferredPrompt = e;

			// 웹 환경 조건이 맞고, 설치 이벤트가 발생했을 때만 팝업 표시
			if (isWebContext && shouldShowInstallPopup()) {
				showInstallPopup = true;
			}
		});

		// await auth.whenReady();
		// if (auth.isAuthenticated()) {
		// 	goto('/app');
		// }
	});
</script>

<div
	class="max-w-md bg-base-100 md:p-4 font-sans text-base-content p-3 relative mx-auto flex min-h-screen flex-col"
>
	<div class="flex flex-1 flex-col items-center justify-center text-center">
		<!-- <div class="mb-8 relative">
			<div
				class="h-20 w-20 rounded-3xl bg-neutral text-neutral-content shadow-2xl relative z-10 flex items-center justify-center overflow-hidden"
			>
				<Flame size={38} class=" fill-primary/10" />
			</div>
			<div class="inset-0 bg-primary/40 blur-2xl absolute z-0 scale-110 rounded-full"></div>
			<div
				class="-top-1 -right-1 h-6 w-6 bg-primary border-base-100 absolute z-20 rounded-full border-4"
			></div>
		</div> -->

		<div class="space-y-2">
			<h1 class="text-4xl font-black tracking-tight gap-1 flex items-center justify-center">
				Torchi<span class="text-primary">.</span>
			</h1>
			<p class="leading-relaxed opacity-70">
				웹훅을 보내고, 폰에서 바로 응답하세요.<br />
				가장 단순한 <span class="font-bold text-primary">Interactive Push</span> 플랫폼
			</p>
		</div>

		<div class="mt-10 space-y-4 w-full text-left">
			<div class="px-1 flex items-center justify-between">
				<span
					class="font-bold tracking-widest gap-1 flex items-center text-[11px] uppercase opacity-50"
				>
					<Bell size={12} /> QUICK Test
				</span>
			</div>
			<div class="group relative">
				<div class="top-3 left-3 absolute z-10">
					<div class="gap-1 flex">
						<span class="h-2.5 w-2.5 bg-neutral-content/40 rounded-full"></span>
						<span class="h-2.5 w-2.5 bg-neutral-content/40 rounded-full"></span>
						<span class="h-2.5 w-2.5 bg-neutral-content/40 rounded-full"></span>
					</div>
				</div>

				<div
					class="rounded-2xl bg-neutral p-4 pt-8 pb-6 font-mono text-neutral-content/90 shadow-2xl border-white/5 relative overflow-hidden border text-[13px]"
				>
					<span class="">curl</span>
					<span class=""> "{$page.url.origin}/api/v1/demo" \ </span>
					<div>&nbsp;-d <span class="text-success/90">'{demoMessage}'</span></div>

					<button
						on:click={handleSubscribe}
						disabled={push.isToggling}
						class="btn right-2 bottom-2 gap-2 rounded-xl shadow-lg btn-sm btn-primary disabled:bg-primary disabled:text-primary-content absolute border-none normal-case transition-all hover:scale-105 active:scale-95 disabled:opacity-80"
					>
						{#if push.isToggling}
							<LoaderCircle size={14} class="animate-spin" />
							Sending...
						{:else}
							<Play size={14} fill="currentColor" />
							Run
						{/if}
					</button>
				</div>
			</div>

			<p class="text-center text-[11px] opacity-50">버튼을 누르면 푸시 알림이 발송됩니다.</p>
		</div>

		<div class="mt-12 space-y-3 w-full">
			<button
				on:click={loginWithGithub}
				class="btn h-14 gap-3 rounded-2xl btn-neutral hover:bg-neutral/90 shadow-xl flex w-full items-center justify-center border-none transition-all"
			>
				<svg
					xmlns="http://www.w3.org/2000/svg"
					width="20"
					height="20"
					fill="currentColor"
					viewBox="0 0 16 16"
				>
					<path
						d="M8 0C3.58 0 0 3.58 0 8c0 3.54 2.29 6.53 5.47 7.59.4.07.55-.17.55-.38 0-.19-.01-.82-.01-1.49-2.01.37-2.53-.49-2.69-.94-.09-.23-.48-.94-.82-1.13-.28-.15-.68-.52-.01-.53.63-.01 1.08.58 1.23.82.72 1.21 1.87.87 2.33.66.07-.52.28-.87.51-1.07-1.78-.2-3.64-.89-3.64-3.95 0-.87.31-1.59.82-2.15-.08-.2-.36-1.02.08-2.12 0 0 .67-.21 2.2.82.64-.18 1.32-.27 2-.27.68 0 1.36.09 2 .27 1.53-1.04 2.2-.82 2.2-.82.44 1.1.16 1.92.08 2.12.51.56.82 1.27.82 2.15 0 3.07-1.87 3.75-3.65 3.95.29.25.54.73.54 1.48 0 1.07-.01 1.93-.01 2.2 0 .21.15.46.55.38A8.012 8.012 0 0 0 16 8c0-4.42-3.58-8-8-8z"
					/>
				</svg>
				<span class="font-bold text-[15px]">GitHub로 시작하기</span>
			</button>

			<button
				on:click={guestLogin}
				class="btn h-14 rounded-2xl btn-outline border-base-content/20 hover:bg-base-200 w-full opacity-70"
			>
				<span class="font-semibold text-[15px]">게스트로 시작하기</span>
			</button>
		</div>
	</div>

	<footer class="py-6 mt-auto text-center text-[12px] opacity-30">© 2025 Torchi Platform</footer>
	{#if showInstallPopup}
		<button
			title="popup-outside"
			class="inset-0 bg-black/30 fixed z-40 transition-opacity"
			transition:fly={{ duration: 200 }}
			on:click={closePopup}
		></button>

		<div
			transition:fly={{ y: '100%', duration: 400, opacity: 1 }}
			class="bottom-0 left-0 right-0 max-w-md fixed z-50 mx-auto w-full"
		>
			<div
				class="border-base-content/10 text-neutral-content bg-base-100 flex w-full flex-col rounded-t-[32px] border-t"
			>
				<button title="close" class="pt-4 pb-2 flex w-full justify-center" on:click={closePopup}
				></button>

				<button
					class="btn btn-circle btn-ghost btn-sm right-4 top-4 text-base-content hover:bg-white/10 absolute"
					on:click={closePopup}
				>
					<X size={24} />
				</button>

				<div class="px-8 gap-2 flex flex-1 flex-col items-center justify-center text-center">
					<div
						class="h-12 w-12 rounded-xl text-primary bg-white drop-shadow-lg flex items-center justify-center"
					>
						<img src="/logo/icon-512x512.png" alt="logo" class="rounded-xl" />
					</div>

					<p class="mt-3 leading-relaxed text-base-content/50 mb-2 text-[15px]">
						{#if isIOS}
							iOS 정책상 <strong>Safari</strong>에서 홈 화면에 추가해야<br />
							알림을 받을 수 있습니다.
						{:else}
							홈 화면에 추가하여<br />서비스를 더 빠르게 경험하세요
						{/if}
					</p>
					{#if isIOS}
						<div class="mb-10 gap-3 text-sm text-base-content/80 flex flex-col">
							<div
								class="gap-4 rounded-xl bg-base-content/10 p-3 px-4 border-base-content/5 flex items-center border"
							>
								<span
									class="h-8 w-8 bg-primary/20 font-bold text-primary flex items-center justify-center rounded-full"
								>
									1
								</span>
								<span>
									하단 툴바의 <Share size={16} class="mx-1 inline opacity-70" />
									<strong class="text-base-content">공유</strong> 버튼 터치
								</span>
							</div>

							<div
								class="gap-4 rounded-xl bg-base-content/10 p-3 px-4 border-base-content/5 flex items-center border"
							>
								<span
									class="h-8 w-8 bg-primary/20 font-bold text-primary flex items-center justify-center rounded-full"
								>
									2
								</span>
								<span>
									메뉴에서 <SquarePlus size={16} class="mx-1 inline opacity-70" />
									<strong class="text-base-content">홈 화면에 추가</strong> 선택
								</span>
							</div>
						</div>
					{:else}
						<div class=" pt-0 pb-10 w-full">
							<button
								on:click={installPwa}
								class="btn btn-neutral rounded-2xl text-md w-full transition-transform active:scale-95"
							>
								홈 화면에 추가
							</button>
						</div>
					{/if}
				</div>
			</div>
		</div>
	{/if}
</div>
