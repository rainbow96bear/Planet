<script lang="ts">
	import { page } from '$app/stores';
	import { goto } from '$app/navigation';
	import { invalidate } from '$app/navigation';
	import './page.css';
	import { updateProfile, uploadProfileImage, deleteProfileImage } from '$lib/api/user.js';

	let { data } = $props();

	const user = $derived(data.user);

	let nickname = $state(user?.nickname ?? '');
	let username = $state(user?.username ?? '');
	// let bio = $state(user?.bio ?? '')

	let loading = $state(false);
	let error = $state('');
	let success = $state(false);

	let profileImage = $state<string | null>(user?.profile_image ?? null);
	let imageLoading = $state(false);
	let imageError = $state('');

	const MAX_PROFILE_IMAGE_SIZE = 5 * 1024 * 1024; // 5MB

	async function onProfileImageChange(e: Event) {
		const input = e.target as HTMLInputElement;
		const file = input.files?.[0];
		input.value = ''; // 같은 파일 재선택 가능하도록 초기화
		if (!file) return;

		imageError = '';

		if (!file.type.startsWith('image/')) {
			imageError = '이미지 파일만 업로드 가능합니다.';
			return;
		}
		if (file.size > MAX_PROFILE_IMAGE_SIZE) {
			imageError = '이미지 크기는 5MB 이하여야 합니다.';
			return;
		}

		imageLoading = true;
		try {
			const res = await uploadProfileImage(user!.userid, file);
			profileImage = res.profile_image;
			await invalidate('app:user');
		} catch (e) {
			imageError = e instanceof Error ? e.message : '이미지 업로드에 실패했습니다.';
		} finally {
			imageLoading = false;
		}
	}

	async function handleRemoveImage() {
		imageLoading = true;
		imageError = '';
		try {
			await deleteProfileImage(user!.userid);
			profileImage = null;
			await invalidate('app:user');
		} catch (e) {
			imageError = e instanceof Error ? e.message : '이미지 삭제에 실패했습니다.';
		} finally {
			imageLoading = false;
		}
	}

	async function handleSubmit() {
		if (!nickname.trim()) {
			error = '닉네임을 입력해주세요.';
			return;
		}
		loading = true;
		error = '';
		success = false;

		try {
			await updateProfile(user!.userid, { nickname: nickname.trim() });
			success = true;
			await invalidate('app:user');
		} catch (e) {
			error = e instanceof Error ? e.message : '수정에 실패했습니다.';
		} finally {
			loading = false;
		}
	}
</script>

<div class="settings-container">
	<div class="settings-header">
		<a href={`/profile/${user?.userid}`} class="back-btn">
			<svg
				width="16"
				height="16"
				viewBox="0 0 24 24"
				fill="none"
				stroke="currentColor"
				stroke-width="2"
				stroke-linecap="round"
				stroke-linejoin="round"
			>
				<polyline points="15 18 9 12 15 6" />
			</svg>
			돌아가기
		</a>
		<h1 class="settings-title">프로필 수정</h1>
	</div>

	<div class="settings-card">
		<div class="profile-image-section">
			<div class="profile-image-picker">
				{#if profileImage}
					<img src={profileImage} alt="프로필 이미지" class="profile-image-image" />
				{:else}
					<div class="profile-image-placeholder-circle">🪐</div>
				{/if}
				<label for="profile-image-upload" class="profile-image-edit-btn" aria-label="이미지 변경">
					<svg
						width="14"
						height="14"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M12 20h9" />
						<path d="M16.5 3.5a2.121 2.121 0 0 1 3 3L7 19l-4 1 1-4Z" />
					</svg>
				</label>
				<input
					id="profile-image-upload"
					type="file"
					accept="image/*"
					onchange={onProfileImageChange}
					disabled={imageLoading}
					hidden
				/>
			</div>
			<div class="profile-image-info">
				<span class="profile-image-name">{nickname || user?.nickname}</span>
				<span class="profile-image-username">@{username}</span>
				{#if profileImage}
					<button
						type="button"
						class="profile-image-remove-link"
						onclick={handleRemoveImage}
						disabled={imageLoading}
					>
						이미지 삭제
					</button>
				{/if}
			</div>
		</div>
		{#if imageError}
			<p class="form-error">{imageError}</p>
		{/if}

		<div class="divider"></div>

		<form
			class="settings-form"
			onsubmit={(e) => {
				e.preventDefault();
				handleSubmit();
			}}
		>
			<div class="form-group">
				<label class="form-label" for="username">아이디</label>
				<input id="username" class="form-input" type="text" value={username} disabled />
				<span class="form-hint">아이디는 변경할 수 없습니다.</span>
			</div>

			<div class="form-group">
				<label class="form-label" for="nickname">닉네임</label>
				<input
					id="nickname"
					class="form-input"
					type="text"
					bind:value={nickname}
					placeholder="닉네임을 입력하세요"
					maxlength={20}
					disabled={loading}
				/>
				<span class="form-hint">{nickname.length} / 20</span>
			</div>

			{#if error}
				<p class="form-error">{error}</p>
			{/if}

			{#if success}
				<p class="form-success">프로필이 수정되었습니다.</p>
			{/if}

			<div class="form-actions">
				<a href={`/profile/${user?.userid}`} class="btn-cancel">취소</a>
				<button type="submit" class="btn-submit" disabled={loading}>
					{#if loading}
						<span class="spinner"></span>
					{:else}
						저장
					{/if}
				</button>
			</div>
		</form>
	</div>
</div>
