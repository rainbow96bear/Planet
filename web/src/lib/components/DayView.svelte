<script lang="ts">
	import type { Task, OrbitSchedule, OrbitLaneItem } from '$lib/types/task';
	import { DAY_MINUTES, startMinutes, endMinutesClamped } from '$lib/utils/schedule';
	import WeekStrip from './WeekStrip.svelte';

	const DAYS = ['일', '월', '화', '수', '목', '금', '토'];
	// 아주 짧은 Orbit 일정도 최소 이 정도 폭은 보이게 강제한다.
	const MIN_BAR_WIDTH_PERCENT = 3;

	let {
		year,
		month,
		selectedDay,
		weekAnchorDate,
		myTasksForDay,
		orbitSchedulesForDay,
		isOwner,
		myLayerOn,
		orbitLayerOn,
		dayViewError,
		activeOrbitScheduleId,
		onBack,
		onSelectDate,
		onPrevWeek,
		onNextWeek,
		onToggleTask,
		onDeleteTask,
		onEditClick,
		onAddClick,
		onOrbitRowClick
	}: {
		year: number;
		month: number;
		selectedDay: number;
		weekAnchorDate: Date;
		myTasksForDay: Task[];
		orbitSchedulesForDay: OrbitSchedule[];
		isOwner: boolean;
		myLayerOn: boolean;
		orbitLayerOn: boolean;
		dayViewError: string;
		activeOrbitScheduleId: string | null;
		onBack: () => void;
		onSelectDate: (d: Date) => void;
		onPrevWeek: () => void;
		onNextWeek: () => void;
		onToggleTask: (task: Task) => void;
		onDeleteTask: (taskId: string) => void;
		onEditClick: (task: Task) => void;
		onAddClick: (day: number, e: MouseEvent) => void;
		onOrbitRowClick: (item: OrbitLaneItem) => void;
	} = $props();

	function formatDayTitle(y: number, m: number, day: number) {
		const d = new Date(y, m - 1, day);
		return `${y}년 ${m}월 ${day}일 (${DAYS[d.getDay()]})`;
	}

	function formatTime(iso: string) {
		const d = new Date(iso);
		return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
	}

	function barStyle(schedule: OrbitSchedule): string {
		const start = startMinutes(schedule.start_at);
		const end = Math.max(endMinutesClamped(schedule.start_at, schedule.end_at), start + 1);
		const left = (start / DAY_MINUTES) * 100;
		const width = Math.max(((end - start) / DAY_MINUTES) * 100, MIN_BAR_WIDTH_PERCENT);
		return `left: ${left}%; width: ${width}%;`;
	}

	function handleRowClick(schedule: OrbitSchedule) {
		const startMin = startMinutes(schedule.start_at);
		const endMin = Math.max(endMinutesClamped(schedule.start_at, schedule.end_at), startMin + 1);
		onOrbitRowClick({ schedule, lane: 0, startMin, endMin });
	}
</script>

<div class="day-view">
	<button class="back-btn" onclick={onBack}>
		<svg
			class="back-icon"
			width="18"
			height="18"
			viewBox="0 0 24 24"
			fill="none"
			stroke="currentColor"
			stroke-width="2.5"
			stroke-linecap="round"
			stroke-linejoin="round"
		>
			<polyline points="15 18 9 12 15 6" />
		</svg>
		<span class="back-title">{formatDayTitle(year, month, selectedDay)}</span>
	</button>

	<WeekStrip
		{weekAnchorDate}
		{year}
		{month}
		{selectedDay}
		{onSelectDate}
		{onPrevWeek}
		{onNextWeek}
	/>

	{#if dayViewError}
		<p class="error-msg">{dayViewError}</p>
	{/if}

	{#if myLayerOn}
		<section class="day-section">
			<h3 class="day-section-title">
				<span class="my-dot"></span> My Schedule ({myTasksForDay.length})
			</h3>
			{#if myTasksForDay.length === 0}
				<p class="empty-msg">일정이 없습니다.</p>
			{:else}
				<ul class="day-item-list">
					{#each myTasksForDay as task (task.id)}
						<li class="day-item {task.is_completed ? 'completed' : ''}">
							<span class="day-item-time">
								{formatTime(task.start_at)} – {formatTime(task.end_at)}
							</span>
							{#if isOwner}
								<button
									class="btn-toggle"
									onclick={() => onToggleTask(task)}
									title={task.is_completed ? '완료 취소' : '완료로 표시'}
								>
									{task.is_completed ? '✓' : '○'}
								</button>
							{/if}
							<span class="day-item-title">{task.title}</span>
							{#if isOwner}
								<button class="btn-edit" onclick={() => onEditClick(task)} title="수정">✎</button>
								<button class="btn-delete" onclick={() => onDeleteTask(task.id)} title="삭제"
									>✕</button
								>
							{/if}
						</li>
					{/each}
				</ul>
			{/if}
			{#if isOwner}
				<button
					class="btn-open-add"
					onclick={(e) => {
						onAddClick(selectedDay, e);
					}}
				>
					<span>+</span> 할 일 추가
				</button>
			{/if}
		</section>
	{/if}

	{#if orbitLayerOn && isOwner}
		<section class="day-section">
			<h3 class="day-section-title">
				<span class="orbit-dot"></span> Orbit Schedule ({orbitSchedulesForDay.length})
			</h3>
			{#if orbitSchedulesForDay.length === 0}
				<p class="empty-msg">Orbit 일정이 없습니다.</p>
			{:else}
				<!-- 시간순으로 한 행에 하나씩 — 얇은 띠로 시간 흐름을 보여주고,
				     이름/제목은 옅은 텍스트로 가볍게만 인지되게 한다. -->
				<ul class="orbit-row-list">
					{#each orbitSchedulesForDay as schedule (schedule.id)}
						<li>
							<button
								class="orbit-row {activeOrbitScheduleId === schedule.id ? 'active' : ''}"
								onclick={() => handleRowClick(schedule)}
							>
								<span class="orbit-row-track">
									<span class="orbit-row-bar" style={barStyle(schedule)}></span>
								</span>
								<span class="orbit-row-label">
									{schedule.nickname} · {schedule.title}
								</span>
							</button>
						</li>
					{/each}
				</ul>
			{/if}
		</section>
	{/if}
</div>

<style>
	.day-view {
		display: flex;
		flex-direction: column;
		gap: var(--space-xl);
	}
	.back-btn {
		display: inline-flex;
		align-items: center;
		gap: 6px;
		align-self: flex-start;
		height: 34px;
		padding: 0 12px 0 8px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: 999px;
		color: var(--text-primary);
		font-size: 0.85rem;
		font-weight: 700;
		font-family: inherit;
		cursor: pointer;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast);
	}
	.back-btn:hover {
		background: var(--surface-alt);
		border-color: var(--planet-primary);
		color: var(--planet-primary);
	}
	.back-icon {
		flex-shrink: 0;
	}
	.back-title {
		white-space: nowrap;
	}
	/* ---------- Sections ---------- */
	.day-section-title {
		display: flex;
		align-items: center;
		gap: 8px;
		margin: 0 0 var(--space-md);
		color: var(--text-primary);
		font-size: 0.9rem;
		font-weight: 700;
	}
	.my-dot {
		display: inline-block;
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--planet-primary);
	}
	.orbit-dot {
		display: inline-block;
		width: 6px;
		height: 6px;
		border-radius: 50%;
		background: var(--text-secondary);
	}
	.day-item-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		list-style: none;
		margin: 0 0 var(--space-md);
		padding: 0;
	}
	.day-item {
		display: flex;
		align-items: center;
		gap: 0.65rem;
		padding: 0.75rem 0.9rem;
		background: var(--bg);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		transition:
			border-color var(--transition-normal),
			background var(--transition-normal);
	}
	.day-item:hover {
		border-color: var(--planet-primary);
	}
	.day-item.completed {
		color: var(--text-muted);
	}
	.day-item.completed .day-item-title {
		text-decoration: line-through;
	}
	.day-item-time {
		flex-shrink: 0;
		color: var(--text-secondary);
		font-size: 0.78rem;
		font-variant-numeric: tabular-nums;
	}
	.day-item-title {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
		font-size: 0.9rem;
		color: inherit;
	}
	/* ---------- Toggle / Delete (My Schedule) ---------- */
	.btn-toggle {
		width: 24px;
		height: 24px;
		border: none;
		background: transparent;
		color: var(--planet-primary);
		font-size: 0.95rem;
		cursor: pointer;
		border-radius: var(--radius-sm);
		flex-shrink: 0;
		transition:
			color var(--transition-normal),
			background var(--transition-normal);
	}
	.btn-toggle:hover {
		background: rgba(var(--planet-primary-rgb), 0.08);
		color: var(--planet-primary-hover);
	}
	.btn-edit {
		width: 24px;
		height: 24px;
		border: none;
		background: transparent;
		color: var(--text-muted);
		font-size: 0.85rem;
		border-radius: var(--radius-sm);
		cursor: pointer;
		flex-shrink: 0;
		transition:
			color var(--transition-normal),
			background var(--transition-normal);
	}
	.btn-edit:hover {
		background: rgba(var(--planet-primary-rgb), 0.08);
		color: var(--planet-primary);
	}
	.btn-delete {
		width: 24px;
		height: 24px;
		border: none;
		background: transparent;
		color: var(--text-muted);
		border-radius: var(--radius-sm);
		cursor: pointer;
		flex-shrink: 0;
		transition:
			color var(--transition-normal),
			background var(--transition-normal);
	}
	.btn-delete:hover {
		background: rgba(var(--danger-rgb), 0.08);
		color: var(--danger);
	}
	.btn-open-add {
		display: flex;
		align-items: center;
		justify-content: center;
		gap: 0.45rem;
		width: 100%;
		height: 42px;
		border-radius: var(--radius-md);
		border: 1px dashed var(--border);
		background: var(--surface);
		color: var(--planet-primary);
		font-size: 0.85rem;
		font-weight: 600;
		font-family: inherit;
		cursor: pointer;
		transition:
			border-color var(--transition-normal),
			background var(--transition-normal),
			color var(--transition-normal);
	}
	.btn-open-add:hover {
		border-color: var(--planet-primary);
		background: rgba(var(--planet-primary-rgb), 0.05);
	}
	.btn-open-add span {
		font-size: 1rem;
		line-height: 1;
	}
	/* ---------- Orbit — 시간순 행 리스트 ---------- */
	.orbit-row-list {
		display: flex;
		flex-direction: column;
		gap: 6px;
		list-style: none;
		margin: 0;
		padding: 0;
	}
	.orbit-row {
		display: flex;
		align-items: center;
		gap: 10px;
		width: 100%;
		padding: 6px 8px;
		background: transparent;
		border: none;
		border-radius: var(--radius-sm);
		cursor: pointer;
		text-align: left;
		transition: background var(--transition-fast);
	}
	.orbit-row:hover,
	.orbit-row.active {
		background: var(--surface-alt);
	}
	.orbit-row-track {
		position: relative;
		flex: 0 0 96px;
		height: 6px;
		border-radius: 3px;
		background: var(--surface-alt);
		overflow: hidden;
	}
	.orbit-row-bar {
		position: absolute;
		top: 0;
		height: 100%;
		border-radius: 3px;
		background: rgba(var(--planet-primary-rgb), 0.55);
	}
	.orbit-row-label {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
		color: var(--text-muted);
		font-size: 0.75rem;
	}
	/* ---------- Empty / Error ---------- */
	.empty-msg {
		margin: 0 0 var(--space-md);
		padding: var(--space-2xl) var(--space-md);
		text-align: center;
		font-size: 0.85rem;
		color: var(--text-secondary);
	}
	.error-msg {
		margin: 0 0 var(--space-md);
		padding: 0.65rem 0.85rem;
		font-size: 0.8rem;
		color: var(--danger);
		background: rgba(var(--danger-rgb), 0.08);
		border: 1px solid rgba(var(--danger-rgb), 0.18);
		border-radius: var(--radius-md);
	}
	/* ---------- Responsive ---------- */
	@media (max-width: 520px) {
		.day-item {
			padding: 0.7rem 0.8rem;
		}
		.day-item-title {
			font-size: 0.875rem;
		}
		.orbit-row-track {
			flex-basis: 72px;
		}
	}
</style>
