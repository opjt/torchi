<script lang="ts">
	import { goto } from '$app/navigation';
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
			<h1 class="text-xl font-black tracking-tight">개인정보처리방침</h1>
		</div>
	</header>

	<main class="max-w-3xl px-6 pt-4 pb-20 mx-auto">
		<div class="mb-8 rounded-3xl border-base-content/5 bg-base-200/40 p-6 border">
			<h2 class="mb-2 text-2xl font-black">Torchi 개인정보처리방침</h2>
			<p class="text-sm opacity-60">시행일: 2026년 1월 10일</p>
		</div>

		<div class="mb-8 rounded-3xl border-base-content/5 bg-base-200/40 p-5 border">
			<p class="leading-relaxed opacity-70">
				Torchi(이하 "서비스")는 이용자의 개인정보를 중요시하며, 「개인정보 보호법」을 준수하고
				있습니다. 본 개인정보처리방침은 이용자가 제공한 개인정보가 어떠한 용도와 방식으로 이용되고
				있으며, 개인정보보호를 위해 어떠한 조치가 취해지고 있는지 알려드립니다.
			</p>
		</div>

		<div class="space-y-8">
			<!-- 제1조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제1조 (개인정보의 수집 항목 및 수집 방법)</h3>

				<div class="space-y-4">
					<div>
						<p class="mb-2 text-sm font-bold opacity-100">1. 수집하는 개인정보 항목</p>
						<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>
									<span class="font-bold opacity-100">소셜 로그인(GitHub 등) 시:</span>
									이메일 주소, 소셜 계정 고유 식별자(ID), 프로필 사진(선택), 닉네임
								</span>
							</li>
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>
									<span class="font-bold opacity-100">서비스 이용 시:</span>
									웹 푸시 엔드포인트(Endpoint URL) 및 인증 키, 서비스 이용 기록, 접속 로그, 쿠키, 접속
									IP 정보
								</span>
							</li>
						</ul>
					</div>

					<div>
						<p class="mb-2 text-sm font-bold opacity-100">2. 개인정보 수집 방법</p>
						<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>홈페이지 내 소셜 로그인(OAuth) 절차</span>
							</li>
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>서비스 이용 과정에서 생성정보 수집 툴을 통한 자동 수집</span>
							</li>
						</ul>
					</div>
				</div>
			</section>

			<!-- 제2조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제2조 (개인정보의 수집 및 이용 목적)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 수집한 개인정보를 다음의 목적을 위해 활용합니다.
				</p>
				<ol class="space-y-2 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span>
							<span class="font-bold opacity-100">서비스 제공:</span>
							푸시 알림 발송, 엔드포인트 관리, 회원제 서비스 제공 및 본인 인증
						</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span>
							<span class="font-bold opacity-100">회원 관리:</span>
							개인 식별, 불량 회원의 부정 이용 방지, 민원 처리 및 공지사항 전달
						</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">3.</span>
						<span>
							<span class="font-bold opacity-100">보안 및 서비스 안정성 확보:</span>
							비인가 접근 방지, 서비스 남용 및 부정 이용 방지
						</span>
					</li>
				</ol>
			</section>

			<!-- 제3조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제3조 (개인정보의 보유 및 이용 기간)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 원칙적으로 회원 탈퇴 시 개인정보를 지체 없이 파기합니다. 단, 관련 법령에 따라
					다음과 같이 일정 기간 보관합니다.
				</p>

				<div class="space-y-4">
					<div>
						<p class="mb-2 text-sm font-bold opacity-100">1. 회원 정보</p>
						<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>
									<span class="font-bold opacity-100">보존 항목:</span>
									이메일, GitHub ID, 푸시 엔드포인트 정보
								</span>
							</li>
							<li class="gap-2 flex">
								<span class="shrink-0">•</span>
								<span>
									<span class="font-bold opacity-100">보존 기간:</span>
									회원 탈퇴 시 즉시 삭제
								</span>
							</li>
						</ul>
					</div>

					<div>
						<p class="mb-2 text-sm font-bold opacity-100">2. 법령에 따른 보존</p>
						<div class="overflow-x-auto">
							<table class="text-sm w-full border-collapse">
								<thead>
									<tr class="border-base-content/10 border-b">
										<th class="px-3 py-2 font-bold text-left opacity-100">보존 항목</th>
										<th class="px-3 py-2 font-bold text-left opacity-100">보존 기간</th>
										<th class="px-3 py-2 font-bold text-left opacity-100">법적 근거</th>
									</tr>
								</thead>
								<tbody class="opacity-70">
									<tr class="border-base-content/5 border-b">
										<td class="px-3 py-2">접속 로그</td>
										<td class="px-3 py-2">3개월</td>
										<td class="px-3 py-2">통신비밀보호법 제15조의2</td>
									</tr>
									<tr>
										<td class="px-3 py-2">부정 이용 기록</td>
										<td class="px-3 py-2">1년</td>
										<td class="px-3 py-2">서비스 부정 이용 방지 (IP, 접속 시각, 이용 내역)</td>
									</tr>
								</tbody>
							</table>
						</div>
					</div>
				</div>
			</section>

			<!-- 제4조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제4조 (개인정보의 파기 절차 및 방법)</h3>
				<ol class="space-y-2 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span>
							<span class="font-bold opacity-100">파기절차:</span>
							이용자가 입력한 정보는 목적 달성(탈퇴 등) 후 즉시 DB에서 영구 삭제됩니다.
						</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span>
							<span class="font-bold opacity-100">파기방법:</span>
							전자적 파일 형태로 저장된 개인정보는 기록을 재생할 수 없는 기술적 방법을 사용하여 삭제합니다.
						</span>
					</li>
				</ol>
			</section>

			<!-- 제5조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제5조 (개인정보의 제3자 제공)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 원칙적으로 이용자의 개인정보를 제3자에게 제공하지 않습니다. 다만, 다음의 경우에는
					예외로 합니다.
				</p>
				<ol class="space-y-2 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span>이용자가 사전에 동의한 경우</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span>
							법령의 규정에 의거하거나, 수사 목적으로 법령에 정해진 절차와 방법에 따라 수사기관의
							요구가 있는 경우
						</span>
					</li>
				</ol>
			</section>

			<!-- 제6조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제6조 (개인정보 처리 위탁)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 원활한 서비스 제공을 위해 다음과 같이 개인정보 처리를 위탁하고 있습니다.
				</p>
				<div class="overflow-x-auto">
					<table class="text-sm w-full border-collapse">
						<thead>
							<tr class="border-base-content/10 border-b">
								<th class="px-3 py-2 font-bold text-left opacity-100">수탁업체</th>
								<th class="px-3 py-2 font-bold text-left opacity-100">위탁업무 내용</th>
								<th class="px-3 py-2 font-bold text-left opacity-100">개인정보 항목</th>
							</tr>
						</thead>
						<tbody class="opacity-70">
							<tr class="border-base-content/5 border-b">
								<td class="px-3 py-2">GitHub</td>
								<td class="px-3 py-2">소셜 로그인 인증</td>
								<td class="px-3 py-2">이메일, GitHub ID, 프로필 사진</td>
							</tr>
							<tr>
								<td class="px-3 py-2">AWS</td>
								<td class="px-3 py-2">서버 호스팅 및 데이터 보관</td>
								<td class="px-3 py-2">모든 수집 정보</td>
							</tr>
						</tbody>
					</table>
				</div>
			</section>

			<!-- 제7조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제7조 (개인정보의 국외 이전)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 제6조에 따른 개인정보 처리 위탁으로 인해 개인정보가 국외로 이전됩니다.
				</p>
				<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이전 받는 자:</span>
							GitHub Inc., Amazon Web Services Inc.
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이전 항목:</span>
							이메일, GitHub ID, 프로필 사진, 서비스 이용 데이터
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이전 국가:</span>
							미국
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이전 목적:</span>
							소셜 로그인 인증, 서비스 운영 및 데이터 보관
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이전 일시 및 방법:</span>
							서비스 이용 시 네트워크를 통한 전송
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">보유 및 이용 기간:</span>
							회원 탈퇴 시 또는 위탁 계약 종료 시까지
						</span>
					</li>
				</ul>
			</section>

			<!-- 제8조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제8조 (정보주체의 권리·의무 및 행사방법)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					이용자는 언제든지 다음의 권리를 행사할 수 있습니다.
				</p>
				<ol class="space-y-2 leading-relaxed opacity-70">
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">1.</span>
						<span>개인정보 열람 요구</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">2.</span>
						<span>개인정보 오류 등이 있을 경우 정정 요구</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">3.</span>
						<span>개인정보 삭제 요구</span>
					</li>
					<li class="gap-3 flex">
						<span class="font-bold shrink-0 opacity-100">4.</span>
						<span>개인정보 처리정지 요구</span>
					</li>
				</ol>
				<p class="mt-3 leading-relaxed opacity-70">
					위 권리는 서비스 내 '회원탈퇴' 기능 또는 개인정보보호책임자에게 이메일을 통해 행사할 수
					있으며, 서비스는 지체 없이 조치하겠습니다.
				</p>
			</section>

			<!-- 제9조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제9조 (개인정보의 안전성 확보조치)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 개인정보의 안전성 확보를 위해 다음과 같은 조치를 취하고 있습니다.
				</p>
				<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">관리적 조치:</span>
							개인정보에 대한 접근 권한을 관리자 1인으로 제한
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">기술적 조치:</span>
							모든 통신 구간 암호화(HTTPS), 중요 데이터(토큰) 암호화 저장
						</span>
					</li>
				</ul>
			</section>

			<!-- 제10조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제10조 (개인정보보호책임자)</h3>
				<p class="mb-3 leading-relaxed opacity-70">
					서비스는 개인정보 처리에 관한 업무를 총괄하고 불만 처리를 위해 아래와 같이 책임자를
					지정하고 있습니다.
				</p>
				<ul class="space-y-2 pl-4 leading-relaxed opacity-70">
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">책임자:</span>
							<span class="font-mono text-sm">박정태</span>
						</span>
					</li>
					<li class="gap-2 flex">
						<span class="shrink-0">•</span>
						<span>
							<span class="font-bold opacity-100">이메일:</span>
							<span class="font-mono text-sm">jtpark1957@gmail.com</span>
						</span>
					</li>
				</ul>
				<p class="mt-3 leading-relaxed opacity-70">
					정보주체는 서비스 이용 중 발생한 모든 개인정보 보호 관련 문의, 불만처리, 피해구제 등에
					관한 사항을 개인정보보호책임자에게 문의하실 수 있습니다.
				</p>
			</section>

			<!-- 제11조 -->
			<section>
				<h3 class="mb-4 text-base font-bold">제11조 (개인정보처리방침의 변경)</h3>
				<p class="leading-relaxed opacity-70">
					이 개인정보처리방침은 2026년 1월 10일부터 적용되며, 법령 및 방침에 따른 변경사항이 있는
					경우 변경사항 시행 7일 전부터 서비스 내 공지사항을 통해 고지하겠습니다.
				</p>
			</section>
		</div>
	</main>
</div>
