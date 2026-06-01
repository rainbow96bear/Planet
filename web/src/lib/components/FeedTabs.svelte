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
        background: var(--bg);
        border: 1px solid var(--border);
        border-radius: 10px;
        padding: 3px;
    }

    .tab-btn {
        padding: 0.4rem 1rem;
        border-radius: 7px;
        font-size: 0.82rem;
        font-weight: 600;
        cursor: pointer;
        border: none;
        background: transparent;
        color: var(--text-secondary);
        font-family: inherit;
        transition: all 0.2s;
    }

    .tab-btn:hover:not(.active) {
        color: var(--planet-primary);
        background: rgba(79, 156, 249, 0.05);
    }

    .tab-btn.active {
        background: var(--planet-primary);
        color: white;
        box-shadow: 0 1px 3px rgba(15, 23, 42, 0.08);
    }
</style>