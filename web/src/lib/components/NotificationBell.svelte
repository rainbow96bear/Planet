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
    : filter === 'follow' ? notifications.filter(n => n.type === 'follow')
    : notifications.filter(n => n.type !== 'follow')
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
                {notif.from_user.nickname.charAt(0)}
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

  .bell-btn {
    position: relative;
    width: 34px;
    height: 34px;
    border-radius: 50%;
    background: none;
    border: none;
    cursor: pointer;
    display: flex;
    align-items: center;
    justify-content: center;
    color: inherit;
    transition: background 0.15s;
  }
  .bell-btn:hover {
    background: rgba(0, 0, 0, 0.06);
  }

  .badge {
    position: absolute;
    top: -2px;
    right: -2px;
    background: #ef4444;
    color: #fff;
    font-size: 10px;
    font-weight: 700;
    min-width: 16px;
    height: 16px;
    border-radius: 999px;
    display: flex;
    align-items: center;
    justify-content: center;
    padding: 0 3px;
    border: 2px solid white;
    line-height: 1;
  }

  .dropdown {
    position: absolute;
    top: calc(100% + 10px);
    right: -8px;
    width: 340px;
    background: #fff;
    border: 1px solid #e5e7eb;
    border-radius: 12px;
    box-shadow: 0 4px 20px rgba(0, 0, 0, 0.12);
    z-index: 200;
    overflow: hidden;
  }

  .drop-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 14px 16px 10px;
    border-bottom: 1px solid #f3f4f6;
  }
  .drop-title {
    font-size: 15px;
    font-weight: 600;
  }
  .mark-all-btn {
    background: none;
    border: none;
    font-size: 12px;
    color: #3b82f6;
    cursor: pointer;
    padding: 4px 8px;
    border-radius: 6px;
  }
  .mark-all-btn:hover {
    background: #eff6ff;
  }

  .tabs {
    display: flex;
    border-bottom: 1px solid #f3f4f6;
    padding: 0 6px;
  }
  .tab {
    background: none;
    border: none;
    border-bottom: 2px solid transparent;
    font-size: 13px;
    color: #6b7280;
    padding: 8px 10px;
    cursor: pointer;
    transition: color 0.15s;
    margin-bottom: -1px;
  }
  .tab:hover {
    color: #111;
  }
  .tab.active {
    color: #111;
    font-weight: 600;
    border-bottom-color: #111;
  }

  .notif-list {
    list-style: none;
    margin: 0;
    padding: 0;
    max-height: 360px;
    overflow-y: auto;
  }

  .notif-empty {
    padding: 32px 16px;
    text-align: center;
    color: #9ca3af;
    font-size: 13px;
  }

  .notif-item {
    display: flex;
    align-items: flex-start;
    gap: 10px;
    padding: 12px 16px;
    border-bottom: 1px solid #f9fafb;
    cursor: pointer;
    transition: background 0.12s;
  }
  .notif-item:last-child {
    border-bottom: none;
  }
  .notif-item:hover {
    background: #f9fafb;
  }
  .notif-item.unread {
    background: #eff6ff;
  }
  .notif-item.unread:hover {
    background: #dbeafe;
  }

  .notif-avatar {
    position: relative;
    width: 40px;
    height: 40px;
    border-radius: 50%;
    background: #e5e7eb;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 14px;
    font-weight: 600;
    flex-shrink: 0;
  }
  .type-badge {
    position: absolute;
    bottom: -2px;
    right: -2px;
    font-size: 11px;
    background: white;
    border-radius: 50%;
    width: 18px;
    height: 18px;
    display: flex;
    align-items: center;
    justify-content: center;
    border: 1px solid #e5e7eb;
  }

  .notif-body {
    flex: 1;
    min-width: 0;
  }
  .notif-msg {
    margin: 0 0 3px;
    font-size: 13px;
    color: #111;
    line-height: 1.45;
  }
  .notif-time {
    font-size: 11px;
    color: #9ca3af;
  }

  .unread-dot {
    width: 8px;
    height: 8px;
    border-radius: 50%;
    background: #3b82f6;
    flex-shrink: 0;
    margin-top: 5px;
  }
</style>