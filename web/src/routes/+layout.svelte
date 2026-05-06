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
        <a href={`/profile/${$page.data.user.userid}`} class="profile-icon">
          <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
            <circle cx="12" cy="8" r="4"/>
            <path d="M4 20c0-4 3.6-7 8-7s8 3 8 7"/>
          </svg>
        </a>
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