<script lang="ts">
  import favicon from '$lib/assets/favicon.png'
  import { page } from '$app/stores'
  import { logout } from '$lib/api/auth'
  import '../app.css'

  let { children } = $props()

  async function handleLogout() {
    await logout()
    window.location.reload()
  }
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<div class="app-layout">
  <header class="header">
    <a href="/" class="logo">🪐 Planet</a>
    <nav class="nav">
      {#if $page.data.user}
        <span class="nickname">{$page.data.user.nickname}</span>
        <button onclick={handleLogout} class="btn-logout">로그아웃</button>
      {:else}
        <a href="/login" class="btn-login">로그인</a>
      {/if}
    </nav>
  </header>

  <main class="main">
    {@render children()}
  </main>
</div>