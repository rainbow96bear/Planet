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
        display: inline-flex;
        align-items: center;
        gap: var(--space-2xs);

        padding: var(--space-2xs);

        background: var(--surface-hover);

        border: 1px solid var(--border);
        border-radius: var(--radius-md);
    }

    .tab-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;

        min-width: 72px;
        height: 36px;

        padding: 0 1rem;

        border: none;
        border-radius: var(--radius-sm);

        background: transparent;

        color: var(--text-secondary);

        font: inherit;
        font-size: 0.85rem;
        font-weight: 600;

        cursor: pointer;

        transition:
            background var(--transition-fast),
            color var(--transition-fast);
    }

    .tab-btn:hover:not(.active) {
        background: rgba(var(--planet-primary-rgb), 0.08);
        color: var(--planet-primary);
    }

    .tab-btn.active {
        background: var(--planet-primary);
        color: var(--text-on-dark);
    }

    .tab-btn:focus-visible {
        outline: none;
        box-shadow: 0 0 0 3px rgba(var(--planet-primary-rgb), 0.18);
    }

    @media (max-width: 520px) {
        .feed-tabs {
            width: 100%;
        }

        .tab-btn {
            flex: 1;
            min-width: 0;
        }
    }
</style>