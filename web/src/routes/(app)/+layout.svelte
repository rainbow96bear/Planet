<script lang="ts">
  import favicon from '$lib/assets/favicon.png'
  import logo from '$lib/assets/planet.png'
  import { page } from '$app/stores'
  import { goto } from '$app/navigation'
  import { logout } from '$lib/api/auth'
  import './layout.css'
	import NotificationBell from '$lib/components/NotificationBell.svelte';

  let { children } = $props()

  let searchQuery = $state('')

  async function handleLogout() {
    await logout()
    window.location.reload()
  }

  function handleSearch(e: Event) {
    e.preventDefault()
    if (!searchQuery.trim()) return
    goto(`/search?q=${encodeURIComponent(searchQuery.trim())}`)
  }
</script>

<svelte:head>
  <link rel="icon" href={favicon} />
</svelte:head>

<div class="app-layout">
  <header class="header">
    <a href="/" class="logo"><img src={logo} alt="Planet" height="32" />Planet</a>

    <form class="search-form" onsubmit={handleSearch}>
      <div class="search-wrapper">
        <svg class="search-icon" width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2.5" stroke-linecap="round" stroke-linejoin="round">
          <circle cx="11" cy="11" r="8"/>
          <line x1="21" y1="21" x2="16.65" y2="16.65"/>
        </svg>
        <input
          class="search-input"
          type="text"
          placeholder="사용자 검색"
          bind:value={searchQuery}
        />
      </div>
    </form>

    <nav class="nav">
      {#if $page.data.user}
      <NotificationBell />
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