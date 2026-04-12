<script lang="ts">
  import { goto } from '$app/navigation'
  import './page.css'

  let username = $state('')
  let nickname = $state('')
  let password = $state('')
  let error = $state('')
  let loading = $state(false)

</script>

<div class="login-container">
  <div class="login-card">
    <div class="login-logo">🪐 Planet</div>
    <div class="login-tagline">우주처럼 넓은 이야기를 나눠요</div>

    <h1 class="login-title">로그인</h1>

    {#if error}
      <p class="error-msg">{error}</p>
    {/if}

    <div class="field">
      <label for="email">아이디</label>
      <input id="email" type="email" bind:value={username} placeholder="영문과 숫자만 사용 가능, 최소 4자리" />
    </div>

    <div class="field">
      <label for="password">비밀번호</label>
      <input id="password" type="password" bind:value={password} placeholder="최소 8자리" />
    </div>

    <!-- <button class="btn-primary" on:click={handleSubmit} disabled={loading}>
      {loading ? '로그인 중...' : '로그인'}
    </button> -->

    <div class="login-footer">
      계정이 없으신가요? <a href="/signup">회원가입</a>
    </div>
  </div>
</div>