<script lang="ts">
    import type { PageData } from './$types'
    import FeedCard from '$lib/components/FeedCard.svelte'
    import FeedTabs from '$lib/components/FeedTabs.svelte'
    import './page.css'

    let { data }: { data: PageData } = $props()

    const isLoggedIn = !!data.user
    let activeTab = $state<'feed' | 'explore'>(isLoggedIn ? 'feed' : 'explore')

    const feeds = $derived(activeTab === 'feed' ? data.feed : data.exploreFeed)
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

    {#if feeds.length === 0}
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
            {#each feeds as feed (feed.feeds_id)}
                <li>
                    <FeedCard {feed} />
                </li>
            {/each}
        </ul>
    {/if}
</div>