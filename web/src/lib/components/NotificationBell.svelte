<script lang="ts">
	import { page } from '$app/stores';
	import { getUnreadCount, getNotifications, markAllRead } from '$lib/api/notification';
	import type { Notification } from '$lib/types/notification';
	let unreadCount = $state(0);
	let notifications = $state<Notification[]>([]);
	let open = $state(false);
	let loading = $state(false);
	let filter = $state<'all' | 'orbit' | 'reaction'>('all');
	let pollingTimer: ReturnType<typeof setInterval> | null = null;
	let markingAll = $state(false);

	const filtered = $derived(
		filter === 'all'
			? notifications
			: filter === 'orbit'
				? notifications.filter((n) => n.type === 'orbit_entered')
				: notifications.filter((n) => n.type !== 'orbit_entered')
	);
	$effect(() => {
		if ($page.data.user) {
			fetchCount();
			pollingTimer = setInterval(fetchCount, 30_000);
		} else {
			unreadCount = 0;
			if (pollingTimer) clearInterval(pollingTimer);
		}
		return () => {
			if (pollingTimer) clearInterval(pollingTimer);
		};
	});
	async function fetchCount() {
		unreadCount = await getUnreadCount();
	}
	async function toggle() {
		if (open) {
			open = false;
			return;
		}
		open = true;
		loading = true;
		try {
			notifications = await getNotifications();
			unreadCount = notifications.filter((n) => !n.is_read).length;
		} catch (e) {
			console.error(e);
		} finally {
			loading = false;
		}
	}

	async function handleMarkAll() {
		if (markingAll) return;
		markingAll = true;
		const prevNotifications = notifications;
		const prevUnreadCount = unreadCount;
		notifications = notifications.map((n) => ({ ...n, is_read: true }));
		unreadCount = 0;
		try {
			await markAllRead();
		} catch (e) {
			console.error(e);
			notifications = prevNotifications;
			unreadCount = prevUnreadCount;
		} finally {
			markingAll = false;
		}
	}
	function timeAgo(dateStr: string) {
		const diff = Date.now() - new Date(dateStr).getTime();
		const min = Math.floor(diff / 60000);
		if (min < 1) return '방금';
		if (min < 60) return `${min}분 전`;
		const h = Math.floor(min / 60);
		if (h < 24) return `${h}시간 전`;
		return `${Math.floor(h / 24)}일 전`;
	}
	const TYPE_ICON: Record<string, string> = {
		orbit_entered: '👤',
		comment: '💬',
		reaction: '❤️'
	};
	const TAB_LABEL: Record<string, string> = {
		all: '전체',
		orbit: 'Orbit',
		reaction: '댓글·반응'
	};
</script>

<svelte:window
	onclick={(e) => {
		if (open && !(e.target as Element).closest('.notif-wrap')) open = false;
	}}
/>
<div class="notif-wrap">
	<button class="bell-btn" onclick={toggle} aria-label="알림 {unreadCount}개">
		<svg
			width="20"
			height="20"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<path d="M18 8A6 6 0 0 0 6 8c0 7-3 9-3 9h18s-3-2-3-9" />
			<path d="M13.73 21a2 2 0 0 1-3.46 0" />
		</svg>
		{#if unreadCount > 0}
			<span class="badge">{unreadCount > 99 ? '99+' : unreadCount}</span>
		{/if}
	</button>
	{#if open}
		<div class="dropdown" role="dialog" aria-label="알림 목록">
			<div class="drop-header">
				<span class="drop-title">알림</span>
				<button class="mark-all-btn" onclick={handleMarkAll} disabled={markingAll}>모두 읽음</button>
			</div>
			<div class="tabs">
				{#each ['all', 'orbit', 'reaction'] as const as t (t)}
					<button class="tab" class:active={filter === t} onclick={() => (filter = t)}>
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
					{#each filtered as notif (notif.id)}
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
	/* ==========================
    Bell
  ========================== */
	.bell-btn {
		position: relative;
		display: flex;
		align-items: center;
		justify-content: center;
		width: 36px;
		height: 36px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 50%;
		color: var(--text-secondary);
		cursor: pointer;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast);
	}
	.bell-btn:hover {
		background: var(--surface-hover);
		border-color: var(--planet-primary);
		color: var(--planet-primary);
	}
	/* ==========================
    Badge
  ========================== */
	.badge {
		position: absolute;
		top: -4px;
		right: -4px;
		min-width: 18px;
		height: 18px;
		padding: 0 5px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--danger);
		color: var(--text-on-dark);
		border: 2px solid var(--surface);
		border-radius: 999px;
		font-size: 0.68rem;
		font-weight: 700;
		line-height: 1;
	}
	/* ==========================
    Dropdown
  ========================== */
	.dropdown {
		position: absolute;
		top: calc(100% + 10px);
		right: 0;
		width: 340px;
		max-width: calc(100vw - 2rem);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		overflow: hidden;
		z-index: var(--z-dropdown);
	}
	/* ==========================
    Header
  ========================== */
	.drop-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		padding: var(--space-md) var(--space-lg);
		border-bottom: 1px solid var(--border);
	}
	.drop-title {
		color: var(--text-primary);
		font-size: 0.95rem;
		font-weight: 700;
	}
	.mark-all-btn {
		padding: 0.4rem 0.7rem;
		background: transparent;
		border: 1px solid transparent;
		border-radius: var(--radius-md);
		color: var(--planet-primary);
		font-size: 0.75rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			background var(--transition-fast),
			color var(--transition-fast);
	}
	.mark-all-btn:hover {
		background: rgba(var(--planet-primary-rgb), 0.08);
	}
	/* ==========================
    Tabs
  ========================== */
	.tabs {
		display: flex;
		gap: var(--space-2xs);
		padding: var(--space-xs);
		background: var(--surface-hover);
		border-bottom: 1px solid var(--border);
	}
	.tab {
		flex: 1;
		height: 34px;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		color: var(--text-secondary);
		font: inherit;
		font-size: 0.8rem;
		font-weight: 600;
		cursor: pointer;
		transition:
			background var(--transition-fast),
			color var(--transition-fast);
	}
	.tab:hover:not(.active) {
		background: rgba(var(--planet-primary-rgb), 0.08);
		color: var(--planet-primary);
	}
	.tab.active {
		background: var(--planet-primary);
		color: var(--text-on-dark);
	}
	/* ==========================
    List
  ========================== */
	.notif-list {
		margin: 0;
		padding: 0;
		list-style: none;
		max-height: 360px;
		overflow-y: auto;
	}
	.notif-list::-webkit-scrollbar {
		width: 5px;
	}
	.notif-list::-webkit-scrollbar-thumb {
		background: var(--border);
		border-radius: 999px;
	}
	.notif-empty {
		padding: var(--space-2xl);
		text-align: center;
		color: var(--text-muted);
		font-size: 0.85rem;
	}
	/* ==========================
    Item
  ========================== */
	.notif-item {
		display: flex;
		align-items: flex-start;
		gap: 0.75rem;
		padding: 0.9rem var(--space-md);
		border-bottom: 1px solid var(--border);
		cursor: pointer;
		transition: background var(--transition-fast);
	}
	.notif-item:last-child {
		border-bottom: none;
	}
	.notif-item:hover {
		background: var(--surface-hover);
	}
	.notif-item.unread {
		background: rgba(var(--planet-primary-rgb), 0.05);
	}
	.notif-item.unread:hover {
		background: rgba(var(--planet-primary-rgb), 0.08);
	}
	/* ==========================
    Avatar
  ========================== */
	.notif-avatar {
		position: relative;
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		background: var(--surface-hover);
		border: 1px solid var(--border);
		border-radius: 50%;
		color: var(--planet-primary);
		font-size: 0.9rem;
		font-weight: 700;
	}
	.type-badge {
		position: absolute;
		right: -2px;
		bottom: -2px;
		width: 18px;
		height: 18px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 50%;
		font-size: 0.7rem;
	}
	/* ==========================
    Body
  ========================== */
	.notif-body {
		flex: 1;
		min-width: 0;
	}
	.notif-msg {
		margin: 0 0 var(--space-2xs);
		color: var(--text-primary);
		font-size: 0.85rem;
		line-height: 1.5;
		word-break: break-word;
	}
	.notif-time {
		color: var(--text-muted);
		font-size: 0.75rem;
	}
	/* ==========================
    Unread Dot
  ========================== */
	.unread-dot {
		width: 8px;
		height: 8px;
		margin-top: 0.45rem;
		border-radius: 50%;
		background: var(--planet-primary);
		flex-shrink: 0;
	}
	/* ==========================
    Responsive
  ========================== */
	@media (max-width: 520px) {
		.dropdown {
			width: min(340px, calc(100vw - 1rem));
			right: -8px;
		}
		.drop-header {
			padding: 0.9rem 1rem;
		}
		.notif-item {
			padding: 0.85rem 1rem;
		}
	}
</style>
