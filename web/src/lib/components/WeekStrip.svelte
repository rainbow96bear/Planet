<script lang="ts">
	import { getWeekDates } from '$lib/utils/date';

	const DAYS = ['일', '월', '화', '수', '목', '금', '토'];

	let {
		weekAnchorDate,
		year,
		month,
		selectedDay,
		onSelectDate,
		onPrevWeek,
		onNextWeek
	}: {
		weekAnchorDate: Date;
		year: number;
		month: number;
		selectedDay: number;
		onSelectDate: (d: Date) => void;
		onPrevWeek: () => void;
		onNextWeek: () => void;
	} = $props();

	function isSelected(d: Date) {
		return d.getFullYear() === year && d.getMonth() + 1 === month && d.getDate() === selectedDay;
	}

	const weekDates = $derived(getWeekDates(weekAnchorDate));
</script>

<div class="week-strip-nav">
	<button class="week-nav-btn" onclick={onPrevWeek} aria-label="이전 주">
		<svg
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="15 18 9 12 15 6" />
		</svg>
	</button>
	<div class="week-strip">
		{#each weekDates as d (d.toISOString())}
			<button
				class="week-strip-day {isSelected(d) ? 'selected' : ''}"
				onclick={() => onSelectDate(d)}
			>
				<span
					class="week-strip-dow {d.getDay() === 0 ? 'sunday' : d.getDay() === 6 ? 'saturday' : ''}"
				>
					{DAYS[d.getDay()]}
				</span>
				<span class="week-strip-date">{d.getDate()}</span>
			</button>
		{/each}
	</div>
	<button class="week-nav-btn" onclick={onNextWeek} aria-label="다음 주">
		<svg
			width="16"
			height="16"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="9 18 15 12 9 6" />
		</svg>
	</button>
</div>

<style>
	.week-strip-nav {
		display: flex;
		align-items: center;
		gap: var(--space-xs);
		padding: var(--space-sm) 0 var(--space-lg);
		border-bottom: 1px solid var(--border);
	}
	.week-nav-btn {
		display: flex;
		align-items: center;
		justify-content: center;
		flex-shrink: 0;
		width: 28px;
		height: 28px;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		color: var(--text-secondary);
		cursor: pointer;
		transition:
			background var(--transition-fast),
			color var(--transition-fast);
	}
	.week-nav-btn:hover {
		background: var(--surface-alt);
		color: var(--planet-primary);
	}
	.week-strip {
		flex: 1;
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: var(--space-xs);
	}
	.week-strip-day {
		display: flex;
		flex-direction: column;
		align-items: center;
		gap: 4px;
		padding: 6px 0;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		cursor: pointer;
		transition: background var(--transition-fast);
	}
	.week-strip-day:hover {
		background: var(--surface-alt);
	}
	.week-strip-dow {
		color: var(--text-muted);
		font-size: 0.68rem;
		font-weight: 600;
	}
	.week-strip-dow.sunday {
		color: var(--danger);
	}
	.week-strip-dow.saturday {
		color: var(--planet-primary);
	}
	.week-strip-date {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 28px;
		height: 28px;
		border-radius: 50%;
		color: var(--text-primary);
		font-size: 0.85rem;
		font-weight: 600;
	}
	.week-strip-day.selected .week-strip-date {
		background: var(--planet-primary);
		color: var(--text-on-dark);
	}
</style>
