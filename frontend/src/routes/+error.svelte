<script>
	import { goto } from '$app/navigation';
	import { page } from '$app/stores';
</script>

<div style="text-align: center; padding: 50px;">
	<h1>{$page.status} 에러가 발생했습니다!</h1>
	<p>{$page.error?.message || '페이지를 찾을 수 없습니다.'}</p>
	<button onclick={() => goto('/')}>홈으로 돌아가기</button>
</div>
