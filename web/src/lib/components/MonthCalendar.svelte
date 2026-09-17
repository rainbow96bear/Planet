<script lang="ts">
	import type { Task, OrbitSchedule } from '$lib/types/task';

	let {
		year,
		month,
		loading,
		isOwner,
		myLayerOn,
		orbitLayerOn,
		getTasksForDay,
		getOrbitSchedulesForDay,
		onPrevMonth,
		onNextMonth,
		onDayClick,
		onAddClick
	}: {
		year: number;
		month: number;
		loading: boolean;
		isOwner: boolean;
		myLayerOn: boolean;
		orbitLayerOn: boolean;
		getTasksForDay: (day: number) => Task[];
		getOrbitSchedulesForDay: (day: number) => OrbitSchedule[];
		onPrevMonth: () => void;
		onNextMonth: () => void;
		onDayClick: (day: number) => void;
		onAddClick: (day: number, e: MouseEvent) => void;
	} = $props();

	const MAX_VISIBLE_TASKS = 2;

	const DAYS = ['일', '월', '화', '수', '목', '금', '토'];

	function getCalendarDays(y: number, m: number) {
		const firstDay = new Date(y, m - 1, 1).getDay();
		const lastDate = new Date(y, m, 0).getDate();
		const days: (number | null)[] = [];
		for (let i = 0; i < firstDay; i++) days.push(null);
		for (let i = 1; i <= lastDate; i++) days.push(i);
		return days;
	}

	function getTodayDay(y: number, m: number) {
		const now = new Date();
		if (now.getFullYear() === y && now.getMonth() + 1 === m) {
			return now.getDate();
		}
		return null;
	}

	function handleCellKeydown(e: KeyboardEvent, day: number) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			onDayClick(day);
		}
	}
</script>

<div class="calendar-nav">
	<button class="nav-btn" onclick={onPrevMonth} disabled={loading}>◀</button>
	<span class="calendar-title">{year}년 {month}월</span>
	<button class="nav-btn" onclick={onNextMonth} disabled={loading}>▶</button>
</div>

{#if loading}
	<div class="calendar-loading">불러오는 중...</div>
{:else}
	<div class="calendar-grid">
		{#each DAYS as day (day)}
			<div class="calendar-day-header">{day}</div>
		{/each}

		{#each getCalendarDays(year, month) as day, i (i)}
			{@const isToday = day === getTodayDay(year, month)}
			{#if day === null}
				<div class="calendar-cell empty" aria-hidden="true"></div>
			{:else}
				{@const dayTasks = getTasksForDay(day)}
				<div
					class="calendar-cell {isToday ? 'today' : ''}"
					role="button"
					tabindex="0"
					onclick={() => onDayClick(day)}
					onkeydown={(e) => handleCellKeydown(e, day)}
				>
					<div class="cell-top">
						<span
							class="day-number {i % 7 === 0 ? 'sunday' : i % 7 === 6 ? 'saturday' : ''} {isToday
								? 'today-number'
								: ''}"
						>
							{day}
						</span>
						{#if isOwner}
							<button
								class="add-task-btn"
								onclick={(e) => onAddClick(day, e)}
								title="할 일 추가"
								aria-label="할 일 추가">+</button
							>
						{/if}
					</div>

					<!-- My Schedule: 내 일정이 우선이라 공간이 허락하는 만큼 제목을 그대로 보여준다 -->
					{#if myLayerOn && dayTasks.length > 0}
						<div class="my-list">
							{#each dayTasks.slice(0, MAX_VISIBLE_TASKS) as task (task.id)}
								<span class="my-chip {task.is_completed ? 'completed' : ''}">{task.title}</span>
							{/each}
							{#if dayTasks.length > MAX_VISIBLE_TASKS}
								<span class="my-more">+{dayTasks.length - MAX_VISIBLE_TASKS}</span>
							{/if}
						</div>
					{/if}

					{#if orbitLayerOn && getOrbitSchedulesForDay(day).length > 0}
						<span class="orbit-strip"></span>
					{/if}
				</div>
			{/if}
		{/each}
	</div>
{/if}

<style>
	.calendar-nav {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-2xl);
	}
	.calendar-title {
		font-size: 1.05rem;
		font-weight: 600;
		color: var(--text-primary);
	}
	.nav-btn {
		width: 36px;
		height: 36px;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		color: var(--text-secondary);
		cursor: pointer;
		transition:
			background var(--transition-fast),
			color var(--transition-fast),
			border-color var(--transition-fast);
	}
	.nav-btn:hover:not(:disabled) {
		background: var(--surface-alt);
		color: var(--planet-primary);
	}
	.nav-btn:disabled {
		opacity: 0.4;
		cursor: not-allowed;
	}
	.calendar-loading {
		padding: var(--space-4xl) 0;
		text-align: center;
		color: var(--text-secondary);
	}
	/* ---------- Grid ---------- */
	.calendar-grid {
		display: grid;
		grid-template-columns: repeat(7, 1fr);
		gap: var(--space-xs);
	}
	.calendar-day-header {
		padding: var(--space-sm) 0;
		text-align: center;
		color: var(--text-secondary);
		font-size: 0.75rem;
		font-weight: 600;
	}
	.calendar-cell {
		box-sizing: border-box;
		position: relative;
		display: flex;
		flex-direction: column;
		gap: 2px;
		width: 100%;
		min-height: 92px;
		padding: var(--space-xs);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		text-align: left;
		cursor: pointer;
		overflow: hidden;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast);
	}
	.calendar-cell:hover {
		background: var(--surface-alt);
	}
	.calendar-cell:hover .add-task-btn {
		opacity: 1;
	}
	.calendar-cell:focus-visible {
		outline: 2px solid var(--planet-primary);
		outline-offset: 2px;
	}
	.calendar-cell.empty {
		background: transparent;
		border-color: transparent;
		cursor: default;
	}
	.calendar-cell.today {
		background: rgba(var(--planet-primary-rgb), 0.05);
		border-color: var(--planet-primary);
	}
	/* ---------- Cell ---------- */
	.cell-top {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: 2px;
	}
	.day-number {
		color: var(--text-secondary);
		font-size: 0.8rem;
	}
	.day-number.sunday {
		color: var(--danger);
	}
	.day-number.saturday {
		color: var(--planet-primary);
	}
	.day-number.today-number {
		color: var(--planet-primary);
		font-weight: 700;
	}
	/* ---------- Add Button ---------- */
	.add-task-btn {
		width: 20px;
		height: 20px;
		display: flex;
		align-items: center;
		justify-content: center;
		opacity: 0;
		padding: 0;
		background: transparent;
		border: none;
		color: var(--text-muted);
		cursor: pointer;
		transition:
			opacity var(--transition-fast),
			color var(--transition-fast);
	}
	.add-task-btn:hover {
		color: var(--planet-primary);
		opacity: 1;
	}
	/* ---------- My Schedule — 제목 텍스트 ---------- */
	.my-list {
		display: flex;
		flex-direction: column;
		gap: 1px;
	}
	.my-chip {
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
		color: var(--text-primary);
		font-size: 0.7rem;
		line-height: 1.3;
	}
	.my-chip.completed {
		color: var(--text-muted);
		text-decoration: line-through;
	}
	.my-more {
		color: var(--text-muted);
		font-size: 0.65rem;
	}
	/* ---------- Orbit Schedule — 지금은 있음/없음만 표시 (Step 2에서 시간축으로 교체 예정) ---------- */
	.orbit-strip {
		display: block;
		height: 3px;
		margin-top: auto;
		border-radius: 2px;
		background: linear-gradient(
			90deg,
			rgba(var(--planet-primary-rgb), 0.35),
			rgba(var(--planet-primary-rgb), 0.15)
		);
	}
	/* ---------- Responsive ---------- */
	@media (max-width: 520px) {
		.calendar-grid {
			gap: 6px;
		}
		.calendar-cell {
			min-height: 74px;
		}
		.my-chip {
			font-size: 0.65rem;
		}
	}
</style>