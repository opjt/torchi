<script lang="ts">
	import { ChevronLeft } from 'lucide-svelte';
</script>

<div class="bg-base-100 font-sans text-base-content min-h-screen">
	<header
		class="top-0 bg-base-100/80 px-6 py-6 backdrop-blur-md sticky z-20 flex items-center justify-between"
	>
		<div class="flex items-center">
			<button
				onclick={() => window.history.back()}
				class="mr-2 -ml-2 p-2 opacity-50 transition-opacity hover:opacity-100"
				title="뒤로가기"
			>
				<ChevronLeft />
			</button>
			<h1 class="text-xl font-black tracking-tight">이용약관</h1>
		</div>
	</header>

	<main class="max-w-3xl px-6 pt-4 pb-20 mx-auto">
		<div class="mb-8 rounded-3xl border-base-content/5 bg-base-200/40 p-6 border">
			<h2 class="mb-2 text-2xl font-black">Torchi 서비스 이용약관</h2>
			<p class="text-sm opacity-60">시행일: 2026년 1월 10일</p>
		</div>

		<div class="space-y-8">
			<!-- 제1조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제1조 (목적)</h3>
				<p class="leading-relaxed opacity-70">
					본 약관은 <span class="font-bold opacity-100">Torchi</span>(이하 "운영자")이 제공하는
					Torchi 서비스(이하 "서비스")의 이용 조건 및 절차에 관한 사항을 규정함을 목적으로 합니다.
				</p>
			</section>

			<!-- 제2조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제2조 (서비스의 제공 및 변경)</h3>
				<ol class="space-y-3 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span
							>본 서비스는 개발자를 위한 웹 푸시 알림 테스트 및 관리 기능을
							<span class="font-bold opacity-100">무료</span>로 제공합니다.</span
						>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span
							>운영자는 서비스의 안정적인 운영을 위해 노력하지만,
							<span class="font-bold opacity-100">100% 가동률을 보장하지는 않습니다.</span></span
						>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">3.</span>
						<span
							>운영자의 사정에 따라 서비스의 기능이 변경되거나, 사전 고지 후 서비스가 종료될 수
							있습니다.</span
						>
					</li>
				</ol>
			</section>

			<!-- 제3조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제3조 (면책 조항)</h3>
				<ol class="space-y-3 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span
							>운영자는 천재지변, 디도스(DDoS) 공격, IDC 장애, 기간통신사업자의 회선 장애 등
							불가항력으로 인해 서비스를 제공할 수 없는 경우에 대한
							<span class="font-bold opacity-100">책임이 면제</span>됩니다.</span
						>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span
							>운영자는
							<span class="font-bold opacity-100">무료로 제공되는 서비스</span> 이용과 관련하여
							이용자에게 발생한 어떠한 손해(데이터 유실, 알림 미수신에 따른 영업 손실 등)에 대해서도
							<span class="font-bold opacity-100">책임을 지지 않습니다.</span></span
						>
					</li>
				</ol>
			</section>

			<!-- 제4조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제4조 (이용자의 의무)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					이용자는 다음 행위를 해서는 안 되며, 적발 시 사전 통보 없이 이용 제한 또는 계정이 삭제될
					수 있습니다.
				</p>
				<ol class="space-y-3 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span>비정상적인 방법으로 과도한 트래픽을 유발하여 서버에 부하를 주는 행위</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span>타인의 정보를 도용하여 가입하는 행위</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">3.</span>
						<span>서비스를 이용하여 불법적인 정보(도박, 음란물 등)를 유통하는 행위</span>
					</li>
				</ol>
			</section>

			<!-- 제5조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제5조 (준거법 및 재판관할)</h3>
				<p class="leading-relaxed opacity-70">
					본 약관과 관련하여 발생한 분쟁에 대해서는 대한민국 법을 준거법으로 하며, 관할 법원은
					운영자의 거주지 관할 법원으로 합니다.
				</p>
			</section>
		</div>

		<!-- 하단 안내 -->
		<!-- <div class="mt-12 rounded-3xl border-base-content/5 bg-base-200/40 p-6 border text-center">
			<p class="text-sm opacity-60">문의사항이 있으시면 설정 페이지에서 피드백을 남겨주세요.</p>
		</div> -->
	</main>
</div>
