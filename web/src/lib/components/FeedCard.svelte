<script lang="ts">
	import { addTaskReaction, removeTaskReaction } from '$lib/api/reaction';
	import { resolve } from '$app/paths';
	import type { Feed } from '$lib/types/feed';
	let {
		feed,
		onupdate
	}: {
		feed: Feed;
		onupdate: (patch: Partial<Feed>) => void;
	} = $props();
	let liked = $state(feed.is_liked ?? false);
	let cheered = $state(feed.is_cheered ?? false);
	let likeCount = $state(feed.like_count ?? 0);
	let cheerCount = $state(feed.cheer_count ?? 0);

	// feed prop 자체가 다른 피드로 교체될 때(리스트 재정렬, 필터 전환 등)
	// 로컬 낙관적 업데이트 상태를 서버 기준 최신 값으로 다시 맞춘다.
	// feed.feed_id가 바뀔 때만 동기화해서, 이 컴포넌트 자신이 만든
	// onupdate로 인한 부모 갱신 루프에 휘말리지 않게 한다.
	$effect(() => {
		void feed.feed_id;
		liked = feed.is_liked ?? false;
		cheered = feed.is_cheered ?? false;
		likeCount = feed.like_count ?? 0;
		cheerCount = feed.cheer_count ?? 0;
	});

	async function toggleLike() {
		if (!feed.task_id) return;
		const prevLiked = liked;
		liked = !liked;
		likeCount += liked ? 1 : -1;
		try {
			if (liked) {
				await addTaskReaction(feed.task_id, 'like');
			} else {
				await removeTaskReaction(feed.task_id, 'like');
			}
			onupdate({ is_liked: liked, like_count: likeCount });
		} catch (e) {
			console.error(e);
			liked = prevLiked;
			likeCount += liked ? 1 : -1;
		}
	}
	async function toggleCheer() {
		if (!feed.task_id) return;
		const prevCheered = cheered;
		const prevCount = cheerCount;
		cheered = !cheered;
		cheerCount += cheered ? 1 : -1;
		try {
			if (cheered) {
				await addTaskReaction(feed.task_id, 'cheer');
			} else {
				await removeTaskReaction(feed.task_id, 'cheer');
			}
			onupdate({ is_cheered: cheered, cheer_count: cheerCount });
		} catch (e) {
			console.error(e);
			cheered = prevCheered;
			cheerCount = prevCount;
		}
	}
	function formatTime(dateStr: string) {
		const diff = Date.now() - new Date(dateStr).getTime();
		const min = Math.floor(diff / 60000);
		const hour = Math.floor(diff / 3600000);
		const day = Math.floor(diff / 86400000);
		if (min < 1) return '방금 전';
		if (min < 60) return `${min}분 전`;
		if (hour < 24) return `${hour}시간 전`;
		return `${day}일 전`;
	}
	const iconMap: Record<string, string> = {
		'task.created': '📝',
		'task.completed': '✅'
	};
	const labelMap: Record<string, string> = {
		'task.created': '추가',
		'task.completed': '완료'
	};
	const textMap: Record<string, string> = {
		'task.created': '할 일을 추가했습니다',
		'task.completed': '할 일을 완료했습니다'
	};
</script>

<article class="feed-card">
	<div class="card-main">
		<span class="feed-icon" aria-hidden="true">
			{iconMap[feed.type] ?? '🪐'}
		</span>
		<div class="feed-body">
			<div class="feed-meta">
				<a href={resolve(`/profile/${feed.actor_id}`)} class="feed-actor">
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
		border-radius: var(--radius-lg);
		padding: var(--space-lg);
		transition:
			border-color var(--transition-fast),
			background var(--transition-fast);
	}
	.feed-card:hover {
		border-color: var(--planet-primary);
	}
	.card-main {
		display: flex;
		align-items: flex-start;
		gap: var(--space-md);
	}
	.feed-icon {
		width: 40px;
		height: 40px;
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		font-size: 1.15rem;
		color: var(--planet-primary);
		background: rgba(var(--planet-primary-rgb), 0.08);
		border: 1px solid rgba(var(--planet-primary-rgb), 0.15);
		border-radius: 50%;
	}
	.feed-body {
		flex: 1;
		min-width: 0;
	}
	.feed-meta {
		display: flex;
		align-items: center;
		flex-wrap: wrap;
		gap: var(--space-xs);
		margin-bottom: 0.35rem;
	}
	.feed-actor {
		color: var(--planet-primary);
		font-size: 0.9rem;
		font-weight: 700;
		text-decoration: none;
		transition: color var(--transition-fast);
	}
	.feed-actor:hover {
		color: var(--planet-primary-hover);
	}
	.type-label {
		padding: 3px 10px;
		border-radius: 999px;
		background: rgba(var(--planet-secondary-rgb), 0.08);
		color: var(--planet-secondary);
		font-size: 0.72rem;
		font-weight: 600;
	}
	.feed-text {
		margin: 0;
		color: var(--text-primary);
		font-size: 0.9rem;
		line-height: 1.5;
		word-break: break-word;
	}
	.task-badge {
		display: inline-flex;
		margin-top: 0.65rem;
		padding: 0.35rem 0.7rem;
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		background: var(--surface-hover);
		color: var(--text-primary);
		font-size: 0.8rem;
		font-weight: 600;
	}
	.feed-time {
		margin-left: var(--space-sm);
		color: var(--text-muted);
		font-size: 0.75rem;
		white-space: nowrap;
		flex-shrink: 0;
	}
	.reactions {
		display: flex;
		gap: var(--space-sm);
		margin-top: var(--space-md);
		padding-top: var(--space-md);
		border-top: 1px solid var(--border);
	}
	.reaction-btn {
		display: inline-flex;
		align-items: center;
		gap: 0.4rem;
		padding: 0.5rem 0.9rem;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 999px;
		color: var(--text-secondary);
		font: inherit;
		font-size: 0.8rem;
		font-weight: 500;
		cursor: pointer;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast);
	}
	.reaction-btn:hover {
		background: var(--surface-hover);
		border-color: var(--planet-primary);
		color: var(--planet-primary);
	}
	.reaction-icon {
		font-size: 0.9rem;
		line-height: 1;
	}
	.reaction-count {
		font-weight: 700;
	}
	.reaction-label {
		color: inherit;
	}
	.reaction-btn.liked {
		background: rgba(var(--planet-secondary-rgb), 0.08);
		border-color: rgba(var(--planet-secondary-rgb), 0.18);
		color: var(--danger);
	}
	.reaction-btn.cheered {
		background: rgba(var(--planet-highlight-rgb), 0.12);
		border-color: rgba(var(--planet-highlight-rgb), 0.22);
		color: var(--planet-highlight-text);
	}
	@media (max-width: 520px) {
		.feed-card {
			padding: var(--space-md);
		}
		.card-main {
			gap: 0.75rem;
		}
		.feed-icon {
			width: 36px;
			height: 36px;
			font-size: 1rem;
		}
		.reaction-label {
			display: none;
		}
		.reactions {
			gap: var(--space-xs);
			flex-wrap: wrap;
		}
	}
</style>