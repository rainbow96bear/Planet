<script lang="ts">
  import { page } from '$app/stores'
  import './layout.css'

  let { children } = $props()
</script>

<div class="app-layout">
  <header class="header">
    <a href="/" class="logo">🪐 Planet</a>
    <nav class="nav">
      {#if $page.data.user}
        <span class="nickname">{$page.data.user.nickname}</span>
        <a href="/logout" class="btn-logout">로그아웃</a>
      {:else}
        <a href="/login" class="btn-login">로그인</a>
      {/if}
    </nav>
  </header>

  <main class="main">
    {@render children()}
  </main>
</div>