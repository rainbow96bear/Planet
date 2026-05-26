<script lang="ts">
    import type { Activity } from '$lib/types'

    let { activity }: { activity: Activity } = $props()

    let liked   = $state(false)
    let cheered = $state(false)
    let likeCount  = $state(activity.reactions?.like ?? 0)
    let cheerCount = $state(activity.reactions?.cheer ?? 0)

    function toggleLike() {
        liked = !liked
        likeCount += liked ? 1 : -1
    }

    function toggleCheer() {
        cheered = !cheered
        cheerCount += cheered ? 1 : -1
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
        'user.followed':  '➕',
    }

    const labelMap: Record<string, string> = {
        'task.created':   '추가',
        'task.completed': '완료',
        'user.followed':  '팔로우',
    }

    const textMap: Record<string, string> = {
        'task.created':   '할 일을 추가했습니다',
        'task.completed': '할 일을 완료했습니다',
        'user.followed':  '누군가를 팔로우했습니다',
    }
</script>

<article class="activity-card">
    <div class="card-main">
        <span class="activity-icon" aria-hidden="true">
            {iconMap[activity.type] ?? '🪐'}
        </span>

        <div class="activity-body">
            <div class="activity-meta">
                <a href={`/profile/${activity.actor_id}`} class="activity-actor">
                    @{activity.actor_nickname}
                </a>
                <span class="type-label">{labelMap[activity.type] ?? '활동'}</span>
            </div>

            <p class="activity-text">
                {textMap[activity.type] ?? '활동했습니다'}
            </p>

            {#if activity.task_title}
                <span class="task-badge">"{activity.task_title}"</span>
            {/if}
        </div>

        <time class="activity-time" datetime={activity.created_at}>
            {formatTime(activity.created_at)}
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
    .activity-card {
        background: #12121a;
        border: 1px solid #a89fd420;
        border-radius: 12px;
        padding: 0.9rem 1rem;
        transition: border-color 0.2s;
    }

    .activity-card:hover {
        border-color: #a89fd440;
    }

    /* ── 메인 행 ── */
    .card-main {
        display: flex;
        align-items: center;
        gap: 0.85rem;
    }

    /* ── 아이콘 ── */
    .activity-icon {
        font-size: 1.25rem;
        width: 38px;
        height: 38px;
        display: flex;
        align-items: center;
        justify-content: center;
        background: #1a1a26;
        border: 1px solid #a89fd425;
        border-radius: 50%;
        flex-shrink: 0;
    }

    /* ── 본문 ── */
    .activity-body {
        flex: 1;
        min-width: 0;
    }

    .activity-meta {
        display: flex;
        align-items: baseline;
        gap: 0.4rem;
        flex-wrap: wrap;
        margin-bottom: 0.2rem;
    }

    .activity-actor {
        font-size: 0.88rem;
        font-weight: 700;
        color: #b2ede6;
        text-decoration: none;
        white-space: nowrap;
        transition: color 0.2s;
    }

    .activity-actor:hover {
        color: #d4f5f0;
    }

    .type-label {
        font-size: 0.75rem;
        font-weight: 600;
        color: #a89fd4;
        background: #a89fd415;
        border: 1px solid #a89fd425;
        border-radius: 20px;
        padding: 1px 7px;
    }

    .activity-text {
        font-size: 0.85rem;
        color: #a89fd4;
        margin: 0;
        word-break: break-word;
    }

    .task-badge {
        display: inline-block;
        margin-top: 0.35rem;
        font-size: 0.8rem;
        font-weight: 600;
        color: #d0c8f0;
        background: #1a1a26;
        border: 1px solid #a89fd430;
        border-radius: 6px;
        padding: 2px 8px;
        word-break: keep-all;
    }

    /* ── 시간 ── */
    .activity-time {
        font-size: 0.75rem;
        color: #a89fd460;
        flex-shrink: 0;
        white-space: nowrap;
        align-self: flex-start;
    }

    /* ── 리액션 ── */
    .reactions {
        display: flex;
        gap: 6px;
        margin-top: 0.75rem;
        padding-top: 0.75rem;
        border-top: 1px solid #a89fd418;
    }

    .reaction-btn {
        display: inline-flex;
        align-items: center;
        gap: 5px;
        padding: 4px 12px;
        border-radius: 20px;
        border: 1px solid #a89fd425;
        background: transparent;
        cursor: pointer;
        font-size: 0.8rem;
        font-family: inherit;
        color: #a89fd4;
        transition: all 0.15s;
    }

    .reaction-btn:hover {
        background: #a89fd415;
        border-color: #a89fd445;
        color: #d0c8f0;
    }

    .reaction-icon {
        font-size: 0.85rem;
        line-height: 1;
    }

    .reaction-count {
        font-weight: 700;
        color: inherit;
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

    /* ── 반응형 ── */
    @media (max-width: 520px) {
        .activity-card {
            padding: 0.75rem;
        }

        .card-main {
            gap: 0.65rem;
        }

        .activity-icon {
            width: 32px;
            height: 32px;
            font-size: 1rem;
        }

        .reaction-label {
            display: none;
        }
    }
</style>