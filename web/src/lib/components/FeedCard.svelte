<script lang="ts">
    import { addTaskReaction, removeTaskReaction } from '$lib/api/reaction'
    import type { Feed } from '$lib/types/feed'

    let { feed, onupdate }: {
        feed: Feed
        onupdate: (patch: Partial<Feed>) => void
    } = $props()

    let liked      = $state(feed.is_liked ?? false)
    let cheered    = $state(feed.is_cheered ?? false)
    let likeCount  = $state(feed.like_count ?? 0)
    let cheerCount = $state(feed.cheer_count ?? 0)

    async function toggleLike() {
        if (!feed.task_id) return
        const prevLiked = liked
        liked      = !liked
        likeCount += liked ? 1 : -1
        try {
            if (liked) {
                await addTaskReaction(feed.task_id, 'like')
            } else {
                await removeTaskReaction(feed.task_id, 'like')
            }
        } catch (e) {
            liked     = prevLiked
            likeCount += liked ? 1 : -1
        }
    }

    async function toggleCheer() {
        if (!feed.task_id) return
        const prevCheered = cheered
        const prevCount   = cheerCount
        cheered      = !cheered
        cheerCount  += cheered ? 1 : -1
        try {
            if (cheered) {
                await addTaskReaction(feed.task_id, 'cheer')
            } else {
                await removeTaskReaction(feed.task_id, 'cheer')
            }
            onupdate({ is_cheered: cheered, cheer_count: cheerCount })
        } catch (e) {
            cheered    = prevCheered
            cheerCount = prevCount
        }
    }

    function formatTime(dateStr: string) {
        const diff = Date.now() - new Date(dateStr).getTime()
        const min  = Math.floor(diff / 60000)
        const hour = Math.floor(diff / 3600000)
        const day  = Math.floor(diff / 86400000)
        if (min < 1)   return '방금 전'
        if (min < 60)  return `${min}분 전`
        if (hour < 24) return `${hour}시간 전`
        return `${day}일 전`
    }

    const iconMap: Record<string, string> = {
        'task.created':   '📝',
        'task.completed': '✅',
    }

    const labelMap: Record<string, string> = {
        'task.created':   '추가',
        'task.completed': '완료',
    }

    const textMap: Record<string, string> = {
        'task.created':   '할 일을 추가했습니다',
        'task.completed': '할 일을 완료했습니다',
    }
</script>

<article class="feed-card">
    <div class="card-main">
        <span class="feed-icon" aria-hidden="true">
            {iconMap[feed.type] ?? '🪐'}
        </span>

        <div class="feed-body">
            <div class="feed-meta">
                <a href={`/profile/${feed.actor_id}`} class="feed-actor">
                    @{feed.actor_nickname}
                </a>
                <span class="type-label">{labelMap[feed.type] ?? '활동'}</span>
            </div>

            <p class="feed-text">
                {textMap[feed.type] ?? '활동했습니다'}
            </p>

            {#if feed.task_title}
                <span class="task-badge">"{feed.task_title}"</span>
            {/if}
        </div>

        <time class="feed-time" datetime={feed.created_at}>
            {formatTime(feed.created_at)}
        </time>
    </div>

    <div class="reactions">
        <button
            class="reaction-btn"
            class:liked
            onclick={toggleLike}
            aria-label="좋아요"
            aria-pressed={liked}
        >
            <span class="reaction-icon">♥</span>
            <span class="reaction-count">{likeCount}</span>
            <span class="reaction-label">좋아요</span>
        </button>

        <button
            class="reaction-btn"
            class:cheered
            onclick={toggleCheer}
            aria-label="응원"
            aria-pressed={cheered}
        >
            <span class="reaction-icon">🔥</span>
            <span class="reaction-count">{cheerCount}</span>
            <span class="reaction-label">응원</span>
        </button>
    </div>
</article>

<style>
    .feed-card {
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 12px;
        padding: 1rem;
        transition: all 0.2s;
    }

    .feed-card:hover {
        border-color: var(--planet-primary);
        box-shadow: 0 4px 12px rgba(15, 23, 42, 0.05);
    }

    .card-main {
        display: flex;
        align-items: center;
        gap: 0.85rem;
    }

    .feed-icon {
        font-size: 1.25rem;
        width: 38px;
        height: 38px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: rgba(79, 156, 249, 0.08);
        border: 1px solid rgba(79, 156, 249, 0.15);
        border-radius: 50%;
        flex-shrink: 0;
    }

    .feed-body {
        flex: 1;
        min-width: 0;
    }

    .feed-meta {
        display: flex;
        align-items: baseline;
        gap: 0.4rem;
        flex-wrap: wrap;
        margin-bottom: 0.25rem;
    }

    .feed-actor {
        font-size: 0.88rem;
        font-weight: 700;
        color: var(--planet-primary);
        text-decoration: none;
        white-space: nowrap;
        transition: color 0.2s;
    }

    .feed-actor:hover {
        color: var(--planet-primary-hover);
    }

    .type-label {
        font-size: 0.75rem;
        font-weight: 600;
        color: var(--planet-secondary);
        background: rgba(175, 168, 230, 0.1);
        border: 1px solid rgba(175, 168, 230, 0.2);
        border-radius: 999px;
        padding: 2px 8px;
    }

    .feed-text {
        font-size: 0.88rem;
        color: var(--text-primary);
        margin: 0;
        word-break: break-word;
    }

    .task-badge {
        display: inline-block;
        margin-top: 0.5rem;
        font-size: 0.8rem;
        font-weight: 600;
        color: var(--text-primary);
        background: var(--bg);
        border: 1px solid var(--border);
        border-radius: 8px;
        padding: 4px 8px;
    }

    .feed-time {
        font-size: 0.75rem;
        color: var(--text-secondary);
        flex-shrink: 0;
        white-space: nowrap;
        align-self: flex-start;
    }

    .reactions {
        display: flex;
        gap: 8px;
        margin-top: 0.9rem;
        padding-top: 0.9rem;
        border-top: 1px solid var(--border);
    }

    .reaction-btn {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        padding: 6px 12px;
        border-radius: 999px;
        border: 1px solid var(--border);
        background: var(--surface);
        cursor: pointer;
        font-size: 0.8rem;
        font-family: inherit;
        color: var(--text-secondary);
        transition: all 0.15s;
    }

    .reaction-btn:hover {
        border-color: var(--planet-primary);
        color: var(--planet-primary);
        background: rgba(79, 156, 249, 0.05);
    }

    .reaction-icon {
        font-size: 0.85rem;
        line-height: 1;
    }

    .reaction-count {
        font-weight: 700;
    }

    /* 좋아요 */
    .reaction-btn.liked {
        background: rgba(239, 68, 68, 0.08);
        border-color: rgba(239, 68, 68, 0.25);
        color: #ef4444;
    }

    /* 응원 */
    .reaction-btn.cheered {
        background: rgba(245, 158, 11, 0.08);
        border-color: rgba(245, 158, 11, 0.25);
        color: #f59e0b;
    }

    .reaction-label {
        color: inherit;
    }
    
    .reaction-btn.liked {
        background: #3d1a1a;
        border-color: #e2524a60;
        color: #f09595;
    }

    .reaction-btn.cheered {
        background: #2e2010;
        border-color: #ef9f2760;
        color: #fac775;
    }

    @media (max-width: 520px) {
        .feed-card {
            padding: 0.85rem;
        }

        .card-main {
            gap: 0.65rem;
        }

        .feed-icon {
            width: 32px;
            height: 32px;
            font-size: 1rem;
        }

        .reaction-label {
            display: none;
        }
    }
</style>