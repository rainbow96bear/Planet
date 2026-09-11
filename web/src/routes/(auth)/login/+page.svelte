<script lang="ts">
	import { resolve } from '$app/paths';
	import logo from '$lib/assets/planet.png';
	import './page.css';

	let loading = $state(false);
	let error = $state('');

	function handleKakaoLogin() {
		loading = true;
		window.location.href = '/api/v1/auth/login/oauth/kakao';
	}

	function handleNaverLogin() {
		loading = true;
		window.location.href = '/api/v1/auth/login/oauth/naver';
	}
</script>

<div class="login-container">
	<div class="login-card">
		<a href={resolve('/')} class="login-logo">
			<img src={logo} alt="Planet" height="28" />
			Planet
		</a>
		<div class="login-tagline">우주처럼 넓은 이야기를 나눠요</div>

		<h1 class="login-title">로그인</h1>

		{#if error}
			<p class="error-msg">{error}</p>
		{/if}

		<!-- 
		<form onsubmit={handleSubmit}>
		<div class="field">
			<label for="username">아이디</label>
			<input
			id="username"
			type="text"
			bind:value={username}
			placeholder="영문과 숫자만 사용 가능, 최소 4자리"
			/>
		</div>

		<div class="field">
			<label for="password">비밀번호</label>
			<input
			id="password"
			type="password"
			bind:value={password}
			placeholder="최소 8자리"
			/>
		</div>

		<button class="btn-primary" type="submit" disabled={loading}>
			{loading ? '로그인 중...' : '로그인'}
		</button>
		</form>

		<div class="divider"><span>또는</span></div>
		-->

		<button class="btn-kakao" onclick={handleKakaoLogin} disabled={loading}>
			카카오로 로그인
		</button>

		<button class="btn-naver" onclick={handleNaverLogin} disabled={loading}>
			네이버로 로그인
		</button>
	</div>
</div>
