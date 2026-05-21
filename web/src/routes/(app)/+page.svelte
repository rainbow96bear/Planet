<script lang="ts">
    import type { PageData } from './$types'
    import './page.css'

    let { data }: { data: PageData } = $props()
    const isLoggedIn = !!data.user
    let activeTab = $state<'feed' | 'explore'>(isLoggedIn ? 'feed' : 'explore')

    const activities = $derived(activeTab === 'feed' ? data.feed : data.exploreFeed)

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

    function getActivityText(type: string, taskTitle?: string) {
        switch (type) {
            case 'task.created':   return `"${taskTitle}" 할 일을 추가했습니다`
            case 'task.completed': return `"${taskTitle}" 할 일을 완료했습니다`
            case 'user.followed':  return '누군가를 팔로우했습니다'
            default:               return '활동했습니다'
        }
    }

    function getActivityIcon(type: string) {
        switch (type) {
            case 'task.created':   return '📝'
            case 'task.completed': return '✅'
            case 'user.followed':  return '➕'
            default:               return '🪐'
        }
    }
</script>

<div class="feed-container">
    <div class="feed-header">
        <h1 class="feed-title">피드</h1>

        <div class="feed-tabs">
            {#if isLoggedIn}
                <button
                    class="tab-btn {activeTab === 'feed' ? 'active' : ''}"
                    onclick={() => activeTab = 'feed'}
                >
                    팔로우
                </button>
            {/if}
            <button
                class="tab-btn {activeTab === 'explore' ? 'active' : ''}"
                onclick={() => activeTab = 'explore'}
            >
                탐색
            </button>
        </div>
    </div>

    {#if activities.length === 0}
        <div class="feed-empty">
            {#if activeTab === 'feed'}
                <p>아직 활동이 없습니다</p>
                <p class="feed-empty-sub">사람들을 팔로우하면 여기에 활동이 표시됩니다</p>
            {:else}
                <p>아직 활동이 없습니다</p>
            {/if}
        </div>
    {:else}
        <ul class="activity-list">
            {#each activities as activity (activity.activity_id)}
                <li class="activity-item">
                    <span class="activity-icon">{getActivityIcon(activity.type)}</span>
                    <div class="activity-body">
                        <a href={`/profile/${activity.actor_id}`} class="activity-actor">
                            @{activity.actor_nickname}
                        </a>
                        <span class="activity-text">
                            {getActivityText(activity.type, activity.task_title)}
                        </span>
                    </div>
                    <span class="activity-time">{formatTime(activity.created_at)}</span>
                </li>
            {/each}
        </ul>
    {/if}
</div>