<script lang="ts">
	import type { PageData } from './$types';
	import type { Feed } from '$lib/types/feed';
	import FeedCard from '$lib/components/FeedCard.svelte';
	import FeedTabs from '$lib/components/FeedTabs.svelte';
	import './page.css';

	let { data }: { data: PageData } = $props();

	const isLoggedIn = $derived(!!data.user);
	let activeTab = $state<'feed' | 'explore'>(data.user ? 'feed' : 'explore');

	// data.feed/data.exploreFeed의 타입이 load 체인 어딘가에서 unknown으로
	// 흘러들어오는 경우를 대비해 Feed[]로 명시한다. 이러면 list, f 등
	// 파생되는 모든 변수의 타입이 자동으로 안전하게 이어진다.
	const feeds = $derived((data.feed ?? []) as Feed[]);
	const exploreFeeds = $derived((data.exploreFeed ?? []) as Feed[]);

	const currentFeeds = $derived(activeTab === 'feed' ? feeds : exploreFeeds);

	function updateFeed(feedId: string, patch: Partial<Feed>) {
		const list = activeTab === 'feed' ? feeds : exploreFeeds;
		const idx = list.findIndex((f) => f.feed_id === feedId);
		if (idx !== -1) list[idx] = { ...list[idx], ...patch };
	}
</script>

<div class="feed-container">
	<div class="feed-header">
		<h1 class="feed-title">피드</h1>
		<FeedTabs active={activeTab} showOrbit={isLoggedIn} onchange={(tab) => (activeTab = tab)} />
	</div>

	{#if currentFeeds.length === 0}
		<div class="feed-empty">
			{#if activeTab === 'feed'}
				<p>아직 활동이 없습니다</p>
				<p class="feed-empty-sub">사람들의 궤도에 들어가면 여기에 활동이 표시됩니다</p>
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
