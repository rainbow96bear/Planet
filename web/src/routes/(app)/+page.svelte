<script lang="ts">
    import type { PageData } from './$types'
    import type { Feed } from '$lib/types/feed'
    import FeedCard from '$lib/components/FeedCard.svelte'
    import FeedTabs from '$lib/components/FeedTabs.svelte'
    import './page.css'

    let { data }: { data: PageData } = $props()

    const isLoggedIn = !!data.user
    let activeTab = $state<'feed' | 'explore'>(isLoggedIn ? 'feed' : 'explore')

    let feeds        = $state([...(data.feed       ?? [])])
    let exploreFeeds = $state([...(data.exploreFeed ?? [])])

    const currentFeeds = $derived(activeTab === 'feed' ? feeds : exploreFeeds)

    function updateFeed(feedId: string, patch: Partial<Feed>) {
        const list = activeTab === 'feed' ? feeds : exploreFeeds
        const idx = list.findIndex(f => f.feed_id === feedId)
        if (idx !== -1) list[idx] = { ...list[idx], ...patch }
    }
</script>

<div class="feed-container">
    <div class="feed-header">
        <h1 class="feed-title">피드</h1>
        <FeedTabs
            active={activeTab}
            showFollow={isLoggedIn}
            onchange={(tab) => activeTab = tab}
        />
    </div>

    {#if currentFeeds.length === 0}
        <div class="feed-empty">
            {#if activeTab === 'feed'}
                <p>아직 활동이 없습니다</p>
                <p class="feed-empty-sub">사람들을 팔로우하면 여기에 활동이 표시됩니다</p>
            {:else}
                <p>아직 활동이 없습니다</p>
            {/if}
        </div>
    {:else}
        <ul class="feed-list">
            {#each currentFeeds as feed (feed.feed_id)}
                <li>
                    <FeedCard {feed} onupdate={(patch) => updateFeed(feed.feed_id, patch)} />
                </li>
            {/each}
        </ul>
    {/if}
</div>