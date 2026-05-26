<script lang="ts">
    type Tab = 'feed' | 'explore'

    let {
        active,
        showFollow = true,
        onchange
    }: {
        active: Tab
        showFollow?: boolean
        onchange: (tab: Tab) => void
    } = $props()
</script>

<div class="feed-tabs">
    {#if showFollow}
        <button
            class="tab-btn"
            class:active={active === 'feed'}
            onclick={() => onchange('feed')}
        >
            팔로우
        </button>
    {/if}
    <button
        class="tab-btn"
        class:active={active === 'explore'}
        onclick={() => onchange('explore')}
    >
        탐색
    </button>
</div>

<style>
    .feed-tabs {
        display: flex;
        gap: 4px;
        background: #1a1a26;
        border: 1px solid #a89fd420;
        border-radius: 10px;
        padding: 3px;
    }

    .tab-btn {
        padding: 0.35rem 0.9rem;
        border-radius: 7px;
        font-size: 0.82rem;
        font-weight: 600;
        cursor: pointer;
        border: none;
        background: none;
        color: #a89fd4;
        font-family: inherit;
        transition: background 0.2s, color 0.2s;
    }

    .tab-btn:hover:not(.active) {
        color: #d0c8f0;
    }

    .tab-btn.active {
        background: #a89fd425;
        color: #ffffff;
    }
</style>