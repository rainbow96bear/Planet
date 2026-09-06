<script lang="ts">
	import logo from '$lib/assets/planet.png';
	import { goto } from '$app/navigation';
	import { resolve } from '$app/paths';
	import { onDestroy } from 'svelte';
	import { createOAuthUser, checkUsername } from '$lib/api/auth';
	import { validateNickname } from '$lib/utils/validation';
	import './page.css';
	import TermsAgreement from '$lib/components/TermsAgreement.svelte';

	let username = $state('');
	let nickname = $state('');
	let error = $state('');
	let loading = $state(false);
	let usernameMsg = $state('');
	let usernameOk = $state(false);
	let nicknameMsg = $state('');
	let agreeTerms = $state(false);
	let agreePrivacy = $state(false);

	let profileImageFile = $state<File | null>(null);
	let profileImagePreview = $state<string | null>(null);
	let profileImageError = $state('');

	const MAX_PROFILE_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB

	let debounceTimer: ReturnType<typeof setTimeout>;
	let usernameCheckId = 0; // 레이스 컨디션 방지용 토큰

	onDestroy(() => {
		clearTimeout(debounceTimer);
		if (profileImagePreview) URL.revokeObjectURL(profileImagePreview);
	});

	function onUsernameInput() {
		username = username.replace(/[^a-zA-Z0-9_]/g, '');
		usernameMsg = '';
		usernameOk = false;
		clearTimeout(debounceTimer);
		if (username.length < 4) {
			usernameMsg = '최소 4자리 이상 입력해주세요.';
			return;
		}
		const myCheckId = ++usernameCheckId;
		debounceTimer = setTimeout(async () => {
			try {
				const res = await checkUsername(username);
				if (myCheckId !== usernameCheckId) return; // 더 최신 입력이 있으면 이 결과는 버림
				usernameOk = res.available;
				usernameMsg = res.available ? '사용 가능한 아이디입니다.' : '이미 사용 중인 아이디입니다.';
			} catch {
				if (myCheckId !== usernameCheckId) return;
				usernameMsg = '확인 중 오류가 발생했습니다.';
			}
		}, 500);
	}

	function onNicknameInput() {
		nicknameMsg = validateNickname(nickname);
	}

	function onProfileImageChange(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		input.value = ''; // 같은 파일 재선택 가능하도록 초기화
		if (!file) return;

		profileImageError = '';

		if (!file.type.startsWith('image/')) {
			profileImageError = '이미지 파일만 업로드 가능합니다.';
			return;
		}
		if (file.size > MAX_PROFILE_IMAGE_SIZE) {
			profileImageError = '이미지 크기는 5MB 이하여야 합니다.';
			return;
		}

		if (profileImagePreview) URL.revokeObjectURL(profileImagePreview);
		profileImageFile = file;
		profileImagePreview = URL.createObjectURL(file);
	}

	function removeProfileImage() {
		if (profileImagePreview) URL.revokeObjectURL(profileImagePreview);
		profileImageFile = null;
		profileImagePreview = null;
		profileImageError = '';
	}

	async function handleSubmit(e: Event) {
		e.preventDefault();
		if (!usernameOk) {
			error = '아이디 중복 확인이 필요합니다.';
			return;
		}
		const nicknameError = validateNickname(nickname);
		if (nicknameError) {
			error = nicknameError;
			return;
		}
		if (!agreeTerms || !agreePrivacy) {
			error = '필수 약관에 모두 동의해주세요.';
			return;
		}
		error = '';
		loading = true;
		try {
			// 회원가입 요청 하나에 이미지까지 함께 전송 (multipart/form-data)
			await createOAuthUser({ username, nickname, agreeTerms, agreePrivacy }, profileImageFile);
			goto(resolve('/login'));
		} catch (e) {
			console.error('oauth signup failed', e);
			error = '회원가입에 실패했습니다. 다시 시도해주세요.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="login-container">
	<div class="login-card">
		<div class="login-logo"><img src={logo} alt="Planet" /> Planet</div>
		<div class="login-tagline">우주처럼 넓은 이야기를 나눠요</div>

		<h1 class="login-title">회원가입</h1>

		{#if error}
			<p class="error-msg">{error}</p>
		{/if}

		<form onsubmit={handleSubmit}>
			<div class="field profile-image-field">
				<label for="profile-image">프로필 이미지 <span class="field-optional">(선택)</span></label>
				<div class="profile-image-picker">
					{#if profileImagePreview}
						<img src={profileImagePreview} alt="프로필 미리보기" class="profile-image-preview" />
						<button
							type="button"
							class="profile-image-remove"
							onclick={removeProfileImage}
							aria-label="이미지 제거"
						>
							✕
						</button>
					{:else}
						<label
							for="profile-image"
							class="profile-image-placeholder"
							aria-label="프로필 이미지 선택"
						>
							<svg
								width="20"
								height="20"
								viewBox="0 0 24 24"
								fill="none"
								stroke="currentColor"
								stroke-width="2"
								stroke-linecap="round"
								stroke-linejoin="round"
							>
								<circle cx="12" cy="8" r="4" />
								<path d="M4 20c0-4 3.5-6 8-6s8 2 8 6" />
							</svg>
							<span>＋</span>
						</label>
					{/if}
					<input
						id="profile-image"
						type="file"
						accept="image/*"
						onchange={onProfileImageChange}
						hidden
					/>
				</div>
				{#if profileImageError}
					<span class="field-error">{profileImageError}</span>
				{/if}
			</div>

			<div class="field">
				<label for="username">아이디</label>
				<input
					id="username"
					type="text"
					autocomplete="username"
					bind:value={username}
					oninput={onUsernameInput}
					maxlength={20}
					placeholder="영문, 숫자, 언더바만 사용 가능, 4~20자"
				/>
				{#if usernameMsg}
					<span class={usernameOk ? 'field-ok' : 'field-error'}>{usernameMsg}</span>
				{/if}
			</div>

			<div class="field">
				<label for="nickname">닉네임</label>
				<input
					id="nickname"
					type="text"
					autocomplete="nickname"
					bind:value={nickname}
					oninput={onNicknameInput}
					maxlength={20}
					placeholder="한글, 영문, 숫자, 언더바 사용 가능, 2~20자"
				/>
				{#if nicknameMsg}
					<span class="field-error">{nicknameMsg}</span>
				{/if}
			</div>

			<TermsAgreement bind:agreeTerms bind:agreePrivacy />

			<button
				class="btn-primary"
				type="submit"
				disabled={loading || !usernameOk || !agreeTerms || !agreePrivacy}
			>
				{loading ? '가입 중...' : '가입하기'}
			</button>
		</form>

		<div class="login-footer">
			이미 계정이 있으신가요? <a href={resolve('/login')}>로그인</a>
		</div>
	</div>
</div>
