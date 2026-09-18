<script lang="ts">
	import type { OrbitSchedule } from '$lib/types/task';
	import { dedupeParticipants } from '$lib/utils/orbit';

	let {
		schedules,
		rangeStartMin,
		rangeEndMin,
		densityBuckets,
		bucketMinutes,
		onClose
	}: {
		schedules: OrbitSchedule[];
		rangeStartMin: number;
		rangeEndMin: number;
		densityBuckets: number[];
		bucketMinutes: number;
		onClose: () => void;
	} = $props();

	const DOT_COLORS = ['#4EC7BC', '#3F7DD6', '#F2A65B', '#E96A5A', '#8B7FD1'];

	function formatClock(min: number): string {
		const h = Math.floor(min / 60) % 24;
		const m = min % 60;
		return `${String(h).padStart(2, '0')}:${String(m).padStart(2, '0')}`;
	}

	// 겹치는 스케줄 중 이름이 다른 사용자만 모아 아바타 스택을 만든다.
	// 같은 사람이 그 시간대에 여러 일정을 가진 경우 아바타는 한 번만 보여준다.
	function uniqueParticipants(list: OrbitSchedule[]) {
		const seen = new Map<string, OrbitSchedule>();
		for (const s of list) {
			if (!seen.has(s.user_id)) seen.set(s.user_id, s);
		}
		return Array.from(seen.values());
	}

	const participants = $derived(dedupeParticipants(schedules));
	const visibleParticipants = $derived(participants.slice(0, 5));
	const overflowCount = $derived(Math.max(participants.length - 5, 0));

	const maxDensity = $derived(Math.max(...densityBuckets, 1));

	function isActiveBucket(bucketIndex: number): boolean {
		const bucketStart = bucketIndex * bucketMinutes;
		const bucketEnd = bucketStart + bucketMinutes;
		return bucketStart < rangeEndMin && bucketEnd > rangeStartMin;
	}

	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onClose();
	}
</script>

<div
	class="modal-backdrop"
	role="presentation"
	onclick={handleBackdropClick}
	onkeydown={(e) => e.key === 'Escape' && onClose()}
>
	<div class="modal">
		<div class="modal-header">
			<span class="modal-time">
				{formatClock(rangeStartMin)} – {formatClock(rangeEndMin)}
				<span class="orbit-badge">Orbit</span>
			</span>
			<button class="modal-close" onclick={onClose} aria-label="닫기">✕</button>
		</div>

		<div class="modal-body">
			<div class="avatar-row">
				{#each visibleParticipants as p (p.user_id)}
					{#if p.profile_image}
						<img class="avatar" src={p.profile_image} alt={p.nickname} />
					{:else}
						<span class="avatar avatar-fallback">{p.nickname.charAt(0)}</span>
					{/if}
				{/each}
				{#if overflowCount > 0}
					<span class="avatar avatar-overflow">+{overflowCount}명</span>
				{/if}
			</div>

			<p class="summary-line">
				{schedules[0]?.title}
				{#if schedules.length > 1}
					<span class="summary-more">외 {schedules.length - 1}건</span>
				{/if}
			</p>

			<ul class="orbit-list">
				{#each schedules as schedule, i (schedule.id)}
					<li class="orbit-list-item">
						<span class="orbit-dot-color" style="background: {DOT_COLORS[i % DOT_COLORS.length]}"
						></span>
						<span class="orbit-list-title">{schedule.title}</span>
						<span class="orbit-list-nickname">{schedule.nickname}</span>
					</li>
				{/each}
			</ul>

			<div class="density-section">
				<h4 class="density-title">시간대별 Orbit 분포</h4>
				<div class="density-chart">
					{#each densityBuckets as count, i (i)}
						<span
							class="density-bar {isActiveBucket(i) ? 'active' : ''}"
							style="height: {Math.max((count / maxDensity) * 100, 4)}%"
						></span>
					{/each}
				</div>
				<div class="density-labels">
					{#each [0, 6, 12, 18, 23] as bucketGroup (bucketGroup)}
						<span>{String(Math.floor((bucketGroup * bucketMinutes) / 60)).padStart(2, '0')}:00</span
						>
					{/each}
				</div>
			</div>

			<button class="btn-close-modal" onclick={onClose}>닫기</button>
		</div>
	</div>
</div>

<style>
	.modal-backdrop {
		position: fixed;
		inset: 0;
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: rgba(34, 58, 94, 0.45);
		backdrop-filter: blur(2px);
		z-index: var(--z-modal);
	}
	.modal {
		width: 100%;
		max-width: 440px;
		max-height: 85vh;
		overflow-y: auto;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: var(--space-xl);
		box-shadow: var(--shadow-xl);
	}
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-lg);
	}
	.modal-time {
		display: inline-flex;
		align-items: center;
		gap: 8px;
		font-size: 1.1rem;
		font-weight: 700;
		color: var(--text-primary);
	}
	.orbit-badge {
		padding: 2px 8px;
		border-radius: 999px;
		background: rgba(var(--planet-primary-rgb), 0.12);
		color: var(--planet-primary);
		font-size: 0.68rem;
		font-weight: 700;
	}
	.modal-close {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: none;
		border-radius: var(--radius-md);
		background: transparent;
		color: var(--text-secondary);
		cursor: pointer;
	}
	.modal-close:hover {
		background: var(--surface-hover);
		color: var(--text-primary);
	}
	/* ---------- Avatars ---------- */
	.avatar-row {
		display: flex;
		align-items: center;
		margin-bottom: var(--space-md);
	}
	.avatar {
		width: 32px;
		height: 32px;
		border-radius: 50%;
		border: 2px solid var(--surface);
		object-fit: cover;
		margin-left: -8px;
	}
	.avatar:first-child {
		margin-left: 0;
	}
	.avatar-fallback {
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--planet-primary);
		color: var(--text-on-dark);
		font-size: 0.8rem;
		font-weight: 700;
	}
	.avatar-overflow {
		display: flex;
		align-items: center;
		justify-content: center;
		width: auto;
		height: 32px;
		padding: 0 10px;
		border-radius: 999px;
		background: var(--surface-alt);
		border: 1px solid var(--border);
		color: var(--text-secondary);
		font-size: 0.75rem;
		font-weight: 600;
		margin-left: 4px;
	}
	/* ---------- Summary ---------- */
	.summary-line {
		margin: 0 0 var(--space-md);
		color: var(--text-primary);
		font-size: 0.95rem;
		font-weight: 700;
	}
	.summary-more {
		color: var(--text-secondary);
		font-weight: 500;
	}
	/* ---------- List ---------- */
	.orbit-list {
		display: flex;
		flex-direction: column;
		gap: 6px;
		list-style: none;
		margin: 0 0 var(--space-lg);
		padding: 0;
	}
	.orbit-list-item {
		display: flex;
		align-items: center;
		gap: 8px;
		font-size: 0.85rem;
		color: var(--text-secondary);
	}
	.orbit-dot-color {
		width: 7px;
		height: 7px;
		border-radius: 50%;
		flex-shrink: 0;
	}
	.orbit-list-title {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
		color: var(--text-primary);
	}
	.orbit-list-nickname {
		flex-shrink: 0;
		font-size: 0.75rem;
		color: var(--text-muted);
	}
	/* ---------- Density chart ---------- */
	.density-section {
		margin-bottom: var(--space-lg);
	}
	.density-title {
		margin: 0 0 var(--space-sm);
		font-size: 0.8rem;
		font-weight: 700;
		color: var(--text-secondary);
	}
	.density-chart {
		display: flex;
		align-items: flex-end;
		gap: 2px;
		height: 56px;
		padding: 6px;
		background: var(--surface-alt);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
	}
	.density-bar {
		flex: 1;
		min-width: 2px;
		background: rgba(var(--planet-primary-rgb), 0.25);
		border-radius: 1px;
	}
	.density-bar.active {
		background: var(--planet-primary);
	}
	.density-labels {
		display: flex;
		justify-content: space-between;
		margin-top: 4px;
		color: var(--text-muted);
		font-size: 0.65rem;
	}
	/* ---------- Close button ---------- */
	.btn-close-modal {
		width: 100%;
		height: 42px;
		border-radius: var(--radius-md);
		border: 1px solid var(--border);
		background: var(--surface);
		color: var(--text-secondary);
		font: inherit;
		font-weight: 600;
		cursor: pointer;
	}
	.btn-close-modal:hover {
		background: var(--surface-hover);
		color: var(--text-primary);
	}
	@media (max-width: 520px) {
		.modal {
			padding: var(--space-lg);
		}
	}
</style>
