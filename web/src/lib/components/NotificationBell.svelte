<script lang="ts">
  import { page } from '$app/stores'
  import { getUnreadCount, getNotifications, markAllRead } from '$lib/api/notification'
  import type { Notification } from '$lib/types/notification'

  let unreadCount = $state(0)
  let notifications = $state<Notification[]>([])
  let open = $state(false)
  let loading = $state(false)
  let filter = $state<'all' | 'follow' | 'reaction'>('all')
  let pollingTimer: ReturnType<typeof setInterval> | null = null

  const filtered = $derived(
    filter === 'all' ? notifications
    : filter === 'follow' ? notifications.filter(n => n.type === 'followed')
    : notifications.filter(n => n.type !== 'followed')
  )

  $effect(() => {
    if ($page.data.user) {
      fetchCount()
      pollingTimer = setInterval(fetchCount, 30_000)
    } else {
      unreadCount = 0
      if (pollingTimer) clearInterval(pollingTimer)
    }
    return () => { if (pollingTimer) clearInterval(pollingTimer) }
  })

  async function fetchCount() {
    unreadCount = await getUnreadCount()
  }

  async function toggle() {
    if (open) { open = false; return }
    open = true
    loading = true
    notifications = await getNotifications()
    console.log("notifications : ", notifications)
    unreadCount = 0
    loading = false
    await markAllRead()
  }

  function handleMarkAll() {
    notifications = notifications.map(n => ({ ...n, is_read: true }))
  }

  function timeAgo(dateStr: string) {
    const diff = Date.now() - new Date(dateStr).getTime()
    const min = Math.floor(diff / 60000)
    if (min < 1) return '방금'
    if (min < 60) return `${min}분 전`
    const h = Math.floor(min / 60)
    if (h < 24) return `${h}시간 전`
    return `${Math.floor(h / 24)}일 전`
  }

  const TYPE_ICON: Record<string, string> = {
    follow: '👤', comment: '💬', reaction: '❤️'
  }

  const TAB_LABEL: Record<string, string> = {
    all: '전체', follow: '팔로우', reaction: '댓글·반응'
  }
</script>

<svelte:window onclick={(e) => {
  if (open && !(e.target as Element).closest('.notif-wrap')) open = false
}} />

<div class="notif-wrap">
  <button class="bell-btn" onclick={toggle} aria-label="알림 {unreadCount}개">
    <svg width="20" height="20" viewBox="0 0 24 24" fill="none"
      stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
      <path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9"/>
      <path d="M13.73 21a2 2 0 0 1-3.46 0"/>
    </svg>
    {#if unreadCount > 0}
      <span class="badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
    {/if}
  </button>

  {#if open}
    <div class="dropdown" role="dialog" aria-label="알림 목록">
      <div class="drop-header">
        <span class="drop-title">알림</span>
        <button class="mark-all-btn" onclick={handleMarkAll}>모두 읽음</button>
      </div>

      <div class="tabs">
        {#each (['all', 'follow', 'reaction'] as const) as t}
          <button
            class="tab"
            class:active={filter === t}
            onclick={() => filter = t}
          >
            {TAB_LABEL[t]}
          </button>
        {/each}
      </div>

      <ul class="notif-list">
        {#if loading}
          <li class="notif-empty">불러오는 중...</li>
        {:else if filtered.length === 0}
          <li class="notif-empty">알림이 없습니다</li>
        {:else}
          {#each filtered as notif}
            <li class="notif-item" class:unread={!notif.is_read}>
              <div class="notif-avatar">
                {notif.actor_nickname.charAt(0)}
                <span class="type-badge">{TYPE_ICON[notif.type]}</span>
              </div>
              <div class="notif-body">
                <p class="notif-msg">{notif.message}</p>
                <span class="notif-time">{timeAgo(notif.created_at)}</span>
              </div>
              {#if !notif.is_read}<span class="unread-dot"></span>{/if}
            </li>
          {/each}
        {/if}
      </ul>
    </div>
  {/if}
</div>

<style>
  .notif-wrap {
    position: relative;
  }

  /* ── 벨 버튼 ── */
  .bell-btn {
    position: relative;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    background: var(--surface);
    border: 1px solid var(--border);
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    transition: all 0.2s;
  }

  .bell-btn:hover {
    border-color: var(--planet-primary);
    color: var(--planet-primary);
  }

  /* ── 뱃지 ── */
  .badge {
    position: absolute;
    top: -3px;
    right: -3px;
    background: var(--danger);
    color: white;
    font-size: 10px;
    font-weight: 700;
    min-width: 16px;
    height: 16px;
    border-radius: 999px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 3px;
    border: 2px solid var(--surface);
    line-height: 1;
  }

  /* ── 드롭다운 ── */
  .dropdown {
    position: absolute;
    top: calc(100% + 10px);
    right: -8px;
    width: 340px;
    background: var(--surface);
    border: 1px solid var(--border);
    border-radius: 16px;
    box-shadow: 0 12px 32px rgba(15, 23, 42, 0.08);
    z-index: 200;
    overflow: hidden;
  }

  /* ── 헤더 ── */
  .drop-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px 10px;
    border-bottom: 1px solid var(--border);
  }

  .drop-title {
    font-size: 15px;
    font-weight: 600;
    color: var(--text-primary);
  }

  .mark-all-btn {
    background: none;
    border: none;
    font-size: 12px;
    color: var(--planet-primary);
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 6px;
    transition: background 0.15s;
  }

  .mark-all-btn:hover {
    background: rgba(79, 156, 249, 0.08);
  }

  /* ── 탭 ── */
  .tabs {
    display: flex;
    border-bottom: 1px solid var(--border);
    padding: 0 6px;
  }

  .tab {
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 13px;
    color: var(--text-secondary);
    padding: 8px 10px;
    cursor: pointer;
    transition: color 0.15s;
    margin-bottom: -1px;
  }

  .tab:hover {
    color: var(--planet-primary);
  }

  .tab.active {
    color: var(--planet-primary);
    font-weight: 600;
    border-bottom-color: var(--planet-primary);
  }

  /* ── 알림 목록 ── */
  .notif-list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 360px;
    overflow-y: auto;
  }

  .notif-list::-webkit-scrollbar {
    width: 4px;
  }

  .notif-list::-webkit-scrollbar-thumb {
    background: var(--border);
    border-radius: 999px;
  }

  .notif-empty {
    padding: 32px 16px;
    text-align: center;
    color: var(--text-secondary);
    font-size: 13px;
  }

  /* ── 알림 아이템 ── */
  .notif-item {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    padding: 12px 16px;
    border-bottom: 1px solid var(--border);
    cursor: pointer;
    transition: background 0.12s;
  }

  .notif-item:last-child {
    border-bottom: none;
  }

  .notif-item:hover {
    background: var(--surface-hover);
  }

  .notif-item.unread {
    background: rgba(79, 156, 249, 0.05);
  }

  .notif-item.unread:hover {
    background: rgba(79, 156, 249, 0.08);
  }

  /* ── 아바타 ── */
  .notif-avatar {
    position: relative;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: var(--bg);
    border: 1px solid var(--border);
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    color: var(--planet-primary);
    flex-shrink: 0;
  }

  .type-badge {
    position: absolute;
    bottom: -2px;
    right: -2px;
    font-size: 11px;
    background: var(--surface);
    border-radius: 50%;
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid var(--border);
  }

  /* ── 본문 ── */
  .notif-body {
    flex: 1;
    min-width: 0;
  }

  .notif-msg {
    margin: 0 0 3px;
    font-size: 13px;
    color: var(--text-primary);
    line-height: 1.45;
  }

  .notif-time {
    font-size: 11px;
    color: var(--text-secondary);
  }

  /* ── 읽지 않음 점 ── */
  .unread-dot {
    width: 7px;
    height: 7px;
    border-radius: 50%;
    background: var(--planet-primary);
    flex-shrink: 0;
    margin-top: 6px;
  }
</style>