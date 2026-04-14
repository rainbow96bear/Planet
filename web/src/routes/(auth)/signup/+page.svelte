<script lang="ts">
  import { goto } from '$app/navigation'
  import { onDestroy } from 'svelte'
  import { createUser, checkUsername } from '$lib/api/auth'
  import './page.css'

  let username = $state('')
  let nickname = $state('')
  let password = $state('')
  let error = $state('')
  let loading = $state(false)
  let usernameMsg = $state('')
  let usernameOk = $state(false)

  let debounceTimer: ReturnType<typeof setTimeout>

  onDestroy(() => clearTimeout(debounceTimer))

  function onUsernameInput() {
    usernameMsg = ''
    usernameOk = false
    clearTimeout(debounceTimer)
    if (username.length < 4) {
      usernameMsg = '최소 4자리 이상 입력해주세요.'
      return
    }
    debounceTimer = setTimeout(async () => {
      try {
        const res = await checkUsername(username)
        usernameOk = res.data.available
        usernameMsg = res.data.available ? '사용 가능한 아이디입니다.' : '이미 사용 중인 아이디입니다.'
      } catch {
        usernameMsg = '확인 중 오류가 발생했습니다.'
      }
    }, 500)
  }

  async function handleSubmit(e: Event) {
    e.preventDefault()
    if (!usernameOk) {
      error = '아이디 중복 확인이 필요합니다.'
      return
    }
    if (nickname.length < 2) {
      error = '닉네임은 최소 2자리입니다.'
      return
    }
    if (password.length < 8) {
      error = '비밀번호는 최소 8자리입니다.'
      return
    }
    error = ''
    loading = true
    try {
      await createUser({ username, nickname, password })
      goto('/login')
    } catch {
      error = '회원가입에 실패했습니다. 다시 시도해주세요.'
    } finally {
      loading = false
    }
  }
</script>

<div class="login-container">
  <div class="login-card">
    <div class="login-logo">🪐 Planet</div>
    <div class="login-tagline">우주처럼 넓은 이야기를 나눠요</div>

    <h1 class="login-title">회원가입</h1>

    {#if error}
      <p class="error-msg">{error}</p>
    {/if}

    <form onsubmit={handleSubmit}>
      <div class="field">
        <label for="username">아이디</label>
        <input
          id="username"
          type="text"
          bind:value={username}
          oninput={onUsernameInput}
          placeholder="영문과 숫자만 사용 가능, 최소 4자리"
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
          bind:value={nickname}
          placeholder="표시될 이름을 입력해주세요"
        />
      </div>

      <button class="btn-primary" type="submit" disabled={loading || !usernameOk}>
        {loading ? '가입 중...' : '가입하기'}
      </button>
    </form>

    <div class="login-footer">
      이미 계정이 있으신가요? <a href="/login">로그인</a>
    </div>
  </div>
</div>