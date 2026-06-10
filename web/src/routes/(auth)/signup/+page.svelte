<script lang="ts">
  import logo from '$lib/assets/planet.png'
  import { goto } from '$app/navigation'
  import { onDestroy } from 'svelte'
  import { createUser, checkUsername } from '$lib/api/auth'
  import { validateNickname } from '$lib/utils/validation'
  import TermsAgreement from '$lib/components/auth/TermsAgreement.svelte'
  import './page.css'

  let username = $state('')
  let nickname = $state('')
  let password = $state('')
  let passwordConfirm = $state('')
  let error = $state('')
  let loading = $state(false)
  let usernameMsg = $state('')
  let usernameOk = $state(false)
  let nicknameMsg = $state('')
  let passwordMsg = $state('')
  let passwordConfirmMsg = $state('')
  let agreeTerms = $state(false)
  let agreePrivacy = $state(false)

  let debounceTimer: ReturnType<typeof setTimeout>

  onDestroy(() => clearTimeout(debounceTimer))

  function onUsernameInput() {
    username = username.replace(/[^a-zA-Z0-9_]/g, '')
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
        usernameOk = res.available
        usernameMsg = res.available ? '사용 가능한 아이디입니다.' : '이미 사용 중인 아이디입니다.'
      } catch {
        usernameMsg = '확인 중 오류가 발생했습니다.'
      }
    }, 500)
  }

  function onNicknameInput() {
    nicknameMsg = validateNickname(nickname)
  }

  function onPasswordInput() {
    if (password.length > 0 && password.length < 8) {
      passwordMsg = '최소 8자리 이상 입력해주세요.'
    } else {
      passwordMsg = ''
    }
    if (passwordConfirm.length > 0) {
      onPasswordConfirmInput()
    }
  }

  function onPasswordConfirmInput() {
    if (passwordConfirm.length === 0) {
      passwordConfirmMsg = ''
      return
    }
    passwordConfirmMsg =
      password === passwordConfirm ? '비밀번호가 일치합니다.' : '비밀번호가 일치하지 않습니다.'
  }

  const passwordConfirmOk = $derived(password.length >= 8 && password === passwordConfirm)

  async function handleSubmit(e: Event) {
    e.preventDefault()
    if (!usernameOk) {
      error = '아이디 중복 확인이 필요합니다.'
      return
    }
    const nicknameError = validateNickname(nickname)
    if (nicknameError) {
      error = nicknameError
      return
    }
    if (password.length < 8) {
      error = '비밀번호는 최소 8자리입니다.'
      return
    }
    if (!passwordConfirmOk) {
      error = '비밀번호가 일치하지 않습니다.'
      return
    }
    if (!agreeTerms || !agreePrivacy) {
      error = '필수 약관에 모두 동의해주세요.'
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
    <div class="login-logo"><img src={logo} alt="Planet" height="32" /> Planet</div>
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
          placeholder="영문, 숫자, 언더바만 사용 가능, 최소 4자리"
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
          oninput={onNicknameInput}
          placeholder="한글, 영문, 숫자, 언더바 사용 가능, 최소 2자리"
        />
        {#if nicknameMsg}
          <span class="field-error">{nicknameMsg}</span>
        {/if}
      </div>

      <div class="field">
        <label for="password">비밀번호</label>
        <input
          id="password"
          type="password"
          bind:value={password}
          oninput={onPasswordInput}
          placeholder="최소 8자리"
        />
        {#if passwordMsg}
          <span class="field-error">{passwordMsg}</span>
        {/if}
      </div>

      <div class="field">
        <label for="password-confirm">비밀번호 확인</label>
        <input
          id="password-confirm"
          type="password"
          bind:value={passwordConfirm}
          oninput={onPasswordConfirmInput}
          placeholder="비밀번호를 다시 입력해주세요"
        />
        {#if passwordConfirmMsg}
          <span class={passwordConfirmOk ? 'field-ok' : 'field-error'}>{passwordConfirmMsg}</span>
        {/if}
      </div>

      <TermsAgreement bind:agreeTerms bind:agreePrivacy />

      <button
        class="btn-primary"
        type="submit"
        disabled={loading || !usernameOk || !passwordConfirmOk || !agreeTerms || !agreePrivacy}
      >
        {loading ? '가입 중...' : '가입하기'}
      </button>
    </form>

    <div class="login-footer">
      이미 계정이 있으신가요? <a href="/login">로그인</a>
    </div>
  </div>
</div>