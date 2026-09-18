<script lang="ts">
	import { getTasksByMonth, getOrbitSchedulesByMonth, deleteTask, toggleTask } from '$lib/api/task';
	import { enterOrbit, leaveOrbit } from '$lib/api/user';
	import type { PageData } from './$types';
	import type { Task, OrbitSchedule, OrbitLaneItem } from '$lib/types/task';
	import ProfileHeader from '$lib/components/ProfileHeader.svelte';
	import MonthCalendar from '$lib/components/MonthCalendar.svelte';
	import DayView from '$lib/components/DayView.svelte';
	import AddTaskModal from '$lib/components/AddTaskModal.svelte';
	import OrbitDetailModal from '$lib/components/OrbitDetailModal.svelte';
	import { shiftDate } from '$lib/utils/date';
	import { DAY_MINUTES, startMinutes, endMinutesClamped } from '$lib/utils/schedule';
	import './page.css';
	import { page } from '$app/stores';

	const BUCKET_MINUTES = 30;

	let { data }: { data: PageData } = $props();
	const userid = $derived($page.params.userid || '');
	const isOwner = $derived(userid === data.me?.userid);
	let tasks = $state<Task[]>(Array.isArray(data.tasks) ? data.tasks : []);
	let orbitSchedules = $state<OrbitSchedule[]>(
		Array.isArray(data.orbitSchedules) ? data.orbitSchedules : []
	);
	let year = $state(data.year);
	let month = $state(data.month);
	let loading = $state(false);
	let isOrbiting = $state(data.profileUser.is_orbiting ?? false);
	let orbitLoading = $state(false);
	let dayViewError = $state('');

	// Layer ON/OFF — 조회와 분리된 순수 표시 상태 (문서 9번 원칙)
	let myLayerOn = $state(true);
	let orbitLayerOn = $state(true);

	// 월간 캘린더 ↔ 특정 날짜 상세(Day View) 전환.
	let viewMode = $state<'month' | 'day'>('month');
	let selectedDay = $state<number | null>(null);
	// Day View의 주간 스트립이 보여주는 "기준 주". prev/next week 이동은 이것만 바꾸고,
	// 실제 selectedDay는 스트립에서 날짜를 직접 클릭해야 바뀐다 — 둘러보기와 확정을 분리.
	let weekAnchorDate = $state<Date>(new Date());
	let addDay = $state<number | null>(null);
	// 수정 모드로 열린 일정. 있으면 AddTaskModal이 수정 모드로 렌더링된다.
	let editingTask = $state<Task | null>(null);
	// 클릭한 Orbit 행 — 이 값이 있으면 상세 모달이 열린다.
	let activeOrbitItem = $state<OrbitLaneItem | null>(null);

	const orbitCount = $derived(data.profileUser.orbit ?? 0);
	const gravityCount = $derived(data.profileUser.gravity ?? 0);
	const modalYear = $derived(editingTask ? new Date(editingTask.start_at).getFullYear() : year);
	const modalMonth = $derived(editingTask ? new Date(editingTask.start_at).getMonth() + 1 : month);
	const modalDay = $derived(editingTask ? new Date(editingTask.start_at).getDate() : (addDay ?? 1));

	$effect(() => {
		tasks = Array.isArray(data.tasks) ? data.tasks : [];
		orbitSchedules = Array.isArray(data.orbitSchedules) ? data.orbitSchedules : [];
	});

	$effect(() => {
		void userid;
		year = data.year;
		month = data.month;
		isOrbiting = data.profileUser.is_orbiting ?? false;
		myLayerOn = true;
		orbitLayerOn = true;
		viewMode = 'month';
		selectedDay = null;
	});

	async function fetchMonth(targetYear: number, targetMonth: number) {
		loading = true;
		try {
			const [newTasks, newOrbitSchedules] = await Promise.all([
				getTasksByMonth(userid, targetYear, targetMonth),
				isOwner ? getOrbitSchedulesByMonth(userid, targetYear, targetMonth) : Promise.resolve([])
			]);
			tasks = newTasks;
			orbitSchedules = newOrbitSchedules;
		} finally {
			loading = false;
		}
	}

	async function prevMonth() {
		if (month === 1) {
			year -= 1;
			month = 12;
		} else {
			month -= 1;
		}
		await fetchMonth(year, month);
	}

	async function nextMonth() {
		if (month === 12) {
			year += 1;
			month = 1;
		} else {
			month += 1;
		}
		await fetchMonth(year, month);
	}

	function openDayView(day: number) {
		selectedDay = day;
		weekAnchorDate = new Date(year, month - 1, day);
		viewMode = 'day';
		activeOrbitItem = null;
		dayViewError = '';
	}

	function backToMonth() {
		viewMode = 'month';
		selectedDay = null;
	}

	async function selectDateInStrip(d: Date) {
		const targetYear = d.getFullYear();
		const targetMonth = d.getMonth() + 1;
		const targetDay = d.getDate();

		if (targetYear !== year || targetMonth !== month) {
			year = targetYear;
			month = targetMonth;
			await fetchMonth(year, month);
		}

		selectedDay = targetDay;
		weekAnchorDate = d;
		activeOrbitItem = null;
		dayViewError = '';
	}

	function goToPrevWeek() {
		weekAnchorDate = shiftDate(weekAnchorDate, -7);
	}

	function goToNextWeek() {
		weekAnchorDate = shiftDate(weekAnchorDate, 7);
	}

	// 특정 시간 구간과 겹치는 스케줄만 모은다 (Orbit 상세 모달용).
	function getOverlappingSchedules(
		schedules: OrbitSchedule[],
		rangeStart: number,
		rangeEnd: number
	): OrbitSchedule[] {
		return schedules.filter((s) => {
			const start = startMinutes(s.start_at);
			const end = Math.max(endMinutesClamped(s.start_at, s.end_at), start + 1);
			return start < rangeEnd && end > rangeStart;
		});
	}

	// 하루를 30분 단위 버킷으로 나눠, 각 구간에 몇 개의 Orbit 일정이 걸쳐있는지 센다.
	function computeDensityBuckets(schedules: OrbitSchedule[]): number[] {
		const bucketCount = DAY_MINUTES / BUCKET_MINUTES;
		const density = new Array(bucketCount).fill(0);
		for (const s of schedules) {
			const start = startMinutes(s.start_at);
			const end = Math.max(endMinutesClamped(s.start_at, s.end_at), start + 1);
			const startBucket = Math.floor(start / BUCKET_MINUTES);
			const endBucket = Math.ceil(end / BUCKET_MINUTES);
			for (let b = startBucket; b < endBucket && b < bucketCount; b++) {
				density[b]++;
			}
		}
		return density;
	}

	function openOrbitDetail(item: OrbitLaneItem) {
		activeOrbitItem = item;
	}

	function closeOrbitDetail() {
		activeOrbitItem = null;
	}

	async function handleDeleteTask(taskId: string) {
		try {
			await deleteTask(taskId);
			tasks = tasks.filter((t) => t.id !== taskId);
		} catch {
			dayViewError = '삭제에 실패했습니다.';
		}
	}

	async function handleToggleTask(task: Task) {
		try {
			await toggleTask(task.id);
			tasks = tasks.map((t) => (t.id === task.id ? { ...t, is_completed: !t.is_completed } : t));
		} catch {
			dayViewError = '변경에 실패했습니다.';
		}
	}

	async function handleEnterOrbit() {
		orbitLoading = true;
		try {
			await enterOrbit(userid);
			isOrbiting = true;
		} catch (e) {
			console.error(e);
		} finally {
			orbitLoading = false;
		}
	}

	async function handleLeaveOrbit() {
		orbitLoading = true;
		try {
			await leaveOrbit(userid);
			isOrbiting = false;
		} catch (e) {
			console.error(e);
		} finally {
			orbitLoading = false;
		}
	}

	function handleTaskSaved(task: Task) {
		tasks = [...tasks.filter((t) => t.id !== task.id), task];
	}

	function handleEditClick(task: Task) {
		editingTask = task;
	}

	// 문서 6번 정책: 일정은 시작일(start_at) 기준 셀/리스트에만 표시한다.
	function getTasksForDay(day: number) {
		return tasks
			.filter((t) => new Date(t.start_at).getDate() === day)
			.sort((a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime());
	}

	function getOrbitSchedulesForDay(day: number) {
		return orbitSchedules
			.filter((s) => new Date(s.start_at).getDate() === day)
			.sort((a, b) => new Date(a.start_at).getTime() - new Date(b.start_at).getTime());
	}

	function openAddModal(day: number, e: MouseEvent) {
		e.stopPropagation();
		addDay = day;
	}

	const orbitOverlapGroup = $derived(
		activeOrbitItem && selectedDay !== null
			? getOverlappingSchedules(
					getOrbitSchedulesForDay(selectedDay),
					activeOrbitItem.startMin,
					activeOrbitItem.endMin
				)
			: []
	);
	const orbitDensityBuckets = $derived(
		selectedDay !== null ? computeDensityBuckets(getOrbitSchedulesForDay(selectedDay)) : []
	);
</script>

<div class="profile-container">
	<ProfileHeader
		profileUser={data.profileUser}
		{orbitCount}
		{gravityCount}
		{isOwner}
		{isOrbiting}
		{orbitLoading}
		onEnterOrbit={handleEnterOrbit}
		onLeaveOrbit={handleLeaveOrbit}
	/>

	<div class="calendar-card">
		<div class="layer-pills">
			<button
				class="layer-pill {myLayerOn ? 'active my' : ''}"
				onclick={() => (myLayerOn = !myLayerOn)}
				aria-pressed={myLayerOn}
			>
				{#if myLayerOn}✓{/if} My Schedule
			</button>
			{#if isOwner}
				<button
					class="layer-pill {orbitLayerOn ? 'active orbit' : ''}"
					onclick={() => (orbitLayerOn = !orbitLayerOn)}
					aria-pressed={orbitLayerOn}
				>
					{#if orbitLayerOn}✓{/if} Orbit Schedule
				</button>
			{/if}
		</div>

		{#if viewMode === 'month'}
			<MonthCalendar
				{year}
				{month}
				{loading}
				{isOwner}
				{myLayerOn}
				{orbitLayerOn}
				{getTasksForDay}
				{getOrbitSchedulesForDay}
				onPrevMonth={prevMonth}
				onNextMonth={nextMonth}
				onDayClick={openDayView}
				onAddClick={openAddModal}
			/>
		{:else if selectedDay !== null}
			<DayView
				{year}
				{month}
				{selectedDay}
				{weekAnchorDate}
				myTasksForDay={getTasksForDay(selectedDay)}
				orbitSchedulesForDay={getOrbitSchedulesForDay(selectedDay)}
				{isOwner}
				{myLayerOn}
				{orbitLayerOn}
				{dayViewError}
				activeOrbitScheduleId={activeOrbitItem?.schedule.id ?? null}
				onBack={backToMonth}
				onSelectDate={selectDateInStrip}
				onPrevWeek={goToPrevWeek}
				onNextWeek={goToNextWeek}
				onToggleTask={handleToggleTask}
				onDeleteTask={handleDeleteTask}
				onEditClick={handleEditClick}
				onAddClick={openAddModal}
				onOrbitRowClick={openOrbitDetail}
			/>
		{/if}
	</div>
</div>

{#if addDay !== null || editingTask !== null}
	<AddTaskModal
		day={modalDay}
		year={modalYear}
		month={modalMonth}
		task={editingTask ?? undefined}
		onClose={() => {
			addDay = null;
			editingTask = null;
		}}
		onSaved={handleTaskSaved}
	/>
{/if}

{#if activeOrbitItem}
	<OrbitDetailModal
		schedules={orbitOverlapGroup}
		rangeStartMin={activeOrbitItem.startMin}
		rangeEndMin={activeOrbitItem.endMin}
		densityBuckets={orbitDensityBuckets}
		bucketMinutes={BUCKET_MINUTES}
		onClose={closeOrbitDetail}
	/>
{/if}
