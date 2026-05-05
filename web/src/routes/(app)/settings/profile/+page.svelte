<script lang="ts">
    import { page } from '$app/stores'
    import { goto } from '$app/navigation'
    import './page.css'
	import { updateProfile } from '$lib/api/user.js';

    let { data } = $props()

    const user = data.user

    let nickname = $state(user?.nickname ?? '')
    let username = $state(user?.username ?? '')
    let bio = $state(user?.bio ?? '')

    let loading = $state(false)
    let error = $state('')
    let success = $state(false)

    async function handleSubmit() {
        if (!nickname.trim()) {
            error = '닉네임을 입력해주세요.'
            return
        }
        loading = true
        error = ''
        success = false

        try {
            await updateProfile(user?.userid!, { nickname: nickname.trim() })
            success = true
        } catch (e) {
            error = e instanceof Error ? e.message : '수정에 실패했습니다.'
        } finally {
            loading = false
        }
    }
</script>

<div class="settings-container">
    <div class="settings-header">
        <a href={`/profile/${user?.userid}`} class="back-btn">
            <svg width="16" height="16" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
                <polyline points="15 18 9 12 15 6"/>
            </svg>
            돌아가기
        </a>
        <h1 class="settings-title">프로필 수정</h1>
    </div>

    <div class="settings-card">
        <div class="avatar-section">
            <div class="avatar">🪐</div>
            <div class="avatar-info">
                <span class="avatar-name">{nickname || user?.nickname}</span>
                <span class="avatar-username">@{username}</span>
            </div>
        </div>

        <div class="divider"></div>

        <form class="settings-form" onsubmit={(e) => { e.preventDefault(); handleSubmit() }}>
            <div class="form-group">
                <label class="form-label" for="username">아이디</label>
                <input
                    id="username"
                    class="form-input"
                    type="text"
                    value={username}
                    disabled
                />
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

            <!-- <div class="form-group">
                <label class="form-label" for="bio">자기소개</label>
                <textarea
                    id="bio"
                    class="form-input form-textarea"
                    bind:value={bio}
                    placeholder="자기소개를 입력하세요"
                    maxlength={150}
                    rows={3}
                    disabled={loading}
                ></textarea>
                <span class="form-hint">{bio.length} / 150</span>
            </div> -->

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