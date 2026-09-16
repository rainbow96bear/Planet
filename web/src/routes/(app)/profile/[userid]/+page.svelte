<script lang="ts">
	import { getTasksByMonth, getOrbitSchedulesByMonth, deleteTask, toggleTask } from '$lib/api/task';
	import { enterOrbit, leaveOrbit } from '$lib/api/user';
	import { resolve } from '$app/paths';
	import type { PageData } from './$types';
	import type { Task, OrbitSchedule, OrbitLaneItem } from '$lib/types/task';
	import AddTaskModal from '$lib/components/AddTaskModal.svelte';
	import OrbitDetailModal from '$lib/components/OrbitDetailModal.svelte';
	import './page.css';
	import { page } from '$app/stores';

	const DAY_MINUTES = 24 * 60;
	const BUCKET_MINUTES = 30;

	let { data }: { data: PageData } = $props();
	const userid = $derived($page.params.userid || '');
	const isOwner = $derived(userid === data.me?.userid);
	const orbitLaneItems = $derived(
		selectedDay !== null ? assignOrbitLanes(getOrbitSchedulesForDay(selectedDay)) : []
	);
	const orbitLaneCount = $derived(
		orbitLaneItems.length > 0 ? Math.max(...orbitLaneItems.map((i) => i.lane)) + 1 : 0
	);
	const orbitOverlapGroup = $derived(
		activeOrbitLane && selectedDay !== null
			? getOverlappingSchedules(
					getOrbitSchedulesForDay(selectedDay),
					activeOrbitLane.startMin,
					activeOrbitLane.endMin
				)
			: []
	);
	const orbitDensityBuckets = $derived(
		selectedDay !== null ? computeDensityBuckets(getOrbitSchedulesForDay(selectedDay)) : []
	);

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

	// 월간 캘린더 ↔ 특정 날짜 상세(Day View) 전환. 기존 TaskModal(오버레이)을 대체.
	let viewMode = $state<'month' | 'day'>('month');
	let selectedDay = $state<number | null>(null);
	let addDay = $state<number | null>(null);
	// Orbit 일정 하나를 눌러 상세(작성자/제목)를 펼친 상태 — 문서 7번 "기본 숨김, 선택 시 상세" 반영.
	// 정식 상세 모달(참여자 아바타/분포 히스토그램)은 Phase 7에서 만들고, 지금은 인라인 펼침으로 대체.
	let activeOrbitLane = $state<OrbitLaneItem | null>(null);

	const orbitCount = $derived(data.profileUser.orbit ?? 0);
	const gravityCount = $derived(data.profileUser.gravity ?? 0);

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
		viewMode = 'day';
		expandedOrbitId = null;
		dayViewError = '';
	}

	function backToMonth() {
		viewMode = 'month';
		selectedDay = null;
	}

	function handleDayCellKeydown(e: KeyboardEvent, day: number) {
		if (e.key === 'Enter' || e.key === ' ') {
			e.preventDefault();
			openDayView(day);
		}
	}

	function openOrbitDetail(item: OrbitLaneItem) {
		activeOrbitLane = item;
	}

	function closeOrbitDetail() {
		activeOrbitLane = null;
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

	function handleTaskCreated(task: Task) {
		tasks = [...tasks, task];
	}

	const DAYS = ['일', '월', '화', '수', '목', '금', '토'];

	function getCalendarDays(year: number, month: number) {
		const firstDay = new Date(year, month - 1, 1).getDay();
		const lastDate = new Date(year, month, 0).getDate();
		const days: (number | null)[] = [];
		for (let i = 0; i < firstDay; i++) days.push(null);
		for (let i = 1; i <= lastDate; i++) days.push(i);
		return days;
	}

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

	function getTodayDay() {
		const now = new Date();
		if (now.getFullYear() === year && now.getMonth() + 1 === month) {
			return now.getDate();
		}
		return null;
	}

	function openAddModal(day: number, e: MouseEvent) {
		e.stopPropagation();
		addDay = day;
	}

	function formatTime(iso: string) {
		const d = new Date(iso);
		return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
	}

	// 선택된 날짜가 속한 주(일~토)의 Date 배열을 구한다.
	function getWeekDates(y: number, m: number, day: number): Date[] {
		const base = new Date(y, m - 1, day);
		const dow = base.getDay(); // 0 = 일요일
		const sunday = new Date(base);
		sunday.setDate(base.getDate() - dow);
		return Array.from({ length: 7 }, (_, i) => {
			const d = new Date(sunday);
			d.setDate(sunday.getDate() + i);
			return d;
		});
	}

	function isSameSelectedDate(d: Date) {
		return (
			selectedDay !== null &&
			d.getFullYear() === year &&
			d.getMonth() + 1 === month &&
			d.getDate() === selectedDay
		);
	}

	// 주간 스트립에서 날짜 클릭 — 같은 달이면 selectedDay만 바꾸고,
	// 다른 달(주가 월 경계에 걸치는 경우)이면 해당 월 데이터를 새로 불러온다.
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
		expandedOrbitId = null;
		dayViewError = '';
	}

	function formatDayTitle(y: number, m: number, day: number) {
		const d = new Date(y, m - 1, day);
		return `${y}년 ${m}월 ${day}일 (${DAYS[d.getDay()]})`;
	}

	function startMinutes(iso: string): number {
		const d = new Date(iso);
		return d.getHours() * 60 + d.getMinutes();
	}

	// 자정을 넘는 일정은 그날 타임라인에서 24:00로 잘라 표시한다 (실제 end_at 데이터는 그대로 유지 — 문서 5번 원칙).
	function endMinutesClamped(startIso: string, endIso: string): number {
		const start = new Date(startIso);
		const end = new Date(endIso);
		if (
			end.getFullYear() !== start.getFullYear() ||
			end.getMonth() !== start.getMonth() ||
			end.getDate() !== start.getDate()
		) {
			return DAY_MINUTES;
		}
		return end.getHours() * 60 + end.getMinutes();
	}

	// 겹치는 Orbit 일정을 레인(가로줄)에 배치한다.
	// schedules는 이미 start_at 기준 정렬되어 들어온다(getOrbitSchedulesForDay).
	// 그리디 방식: 시작 시각이 어떤 레인의 마지막 종료 시각보다 크거나 같으면 그 레인 재사용, 없으면 새 레인 생성.
	function assignOrbitLanes(schedules: OrbitSchedule[]): OrbitLaneItem[] {
		const laneEndTimes: number[] = [];
		const result: OrbitLaneItem[] = [];

		for (const schedule of schedules) {
			const startMin = startMinutes(schedule.start_at);
			const endMin = Math.max(endMinutesClamped(schedule.start_at, schedule.end_at), startMin + 1);

			let laneIndex = laneEndTimes.findIndex((endTime) => startMin >= endTime);
			if (laneIndex === -1) {
				laneIndex = laneEndTimes.length;
				laneEndTimes.push(endMin);
			} else {
				laneEndTimes[laneIndex] = endMin;
			}

			result.push({ schedule, lane: laneIndex, startMin, endMin });
		}

		return result;
	}

	// 특정 시간 구간(startMin~endMin)과 겹치는 스케줄만 모은다.
	// 레인 인덱스가 아니라 실제 시간 겹침을 기준으로 다시 계산한다 (레인 배치와 별개 목적).
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
</script>

<div class="profile-container">
	<div class="profile-header">
		<div class="profile-image">
			{#if data.profileUser.profile_image}
				<img
					src={data.profileUser.profile_image}
					alt="{data.profileUser.nickname}님의 프로필 이미지"
					class="profile-image-photo"
				/>
			{:else}
				<span class="profile-image-placeholder">🪐</span>
			{/if}
		</div>
		<div class="profile-info">
			<h1 class="profile-nickname">{data.profileUser.nickname}</h1>
			<span class="profile-username">@{data.profileUser.username}</span>

			<div class="orbit-stats">
				<span class="orbit-stat">
					<strong>{orbitCount}</strong>
					<span class="orbit-stat-label">Orbit</span>
				</span>
				<span class="orbit-stat">
					<strong>{gravityCount}</strong>
					<span class="orbit-stat-label">Gravity</span>
				</span>
			</div>
		</div>

		<div class="profile-actions">
			{#if isOwner}
				<a href={resolve('/settings/profile')} class="action-btn secondary">
					<svg
						width="14"
						height="14"
						viewBox="0 0 24 24"
						fill="none"
						stroke="currentColor"
						stroke-width="2"
						stroke-linecap="round"
						stroke-linejoin="round"
					>
						<path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7" />
						<path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z" />
					</svg>
					프로필 수정
				</a>
			{:else if orbitLoading}
				<button class="action-btn primary" disabled aria-label="처리 중">
					<span class="spinner"></span>
				</button>
			{:else if isOrbiting}
				<button class="action-btn orbiting" onclick={handleLeaveOrbit}>
					<span class="btn-text-default">In Orbit</span>
					<span class="btn-text-hover">Leave Orbit</span>
				</button>
			{:else}
				<button class="action-btn primary" onclick={handleEnterOrbit}> Enter Orbit </button>
			{/if}
		</div>
	</div>

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
			<div class="calendar-nav">
				<button class="nav-btn" onclick={prevMonth} disabled={loading}>◀</button>
				<span class="calendar-title">{year}년 {month}월</span>
				<button class="nav-btn" onclick={nextMonth} disabled={loading}>▶</button>
			</div>

			{#if loading}
				<div class="calendar-loading">불러오는 중...</div>
			{:else}
				<div class="calendar-grid">
					{#each DAYS as day (day)}
						<div class="calendar-day-header">{day}</div>
					{/each}

					{#each getCalendarDays(year, month) as day, i (i)}
						{@const isToday = day === getTodayDay()}
						{#if day === null}
							<div class="calendar-cell empty" aria-hidden="true"></div>
						{:else}
							<div
								class="calendar-cell {isToday ? 'today' : ''}"
								role="button"
								tabindex="0"
								onclick={() => openDayView(day)}
								onkeydown={(e) => handleDayCellKeydown(e, day)}
							>
								<div class="cell-top">
									<span
										class="day-number {i % 7 === 0
											? 'sunday'
											: i % 7 === 6
												? 'saturday'
												: ''} {isToday ? 'today-number' : ''}"
									>
										{day}
									</span>
									{#if isOwner}
										<button
											class="add-task-btn"
											onclick={(e) => openAddModal(day, e)}
											title="할 일 추가"
											aria-label="할 일 추가">+</button
										>
									{/if}
								</div>
								<!-- 점/띠 정식 렌더링은 Phase 6에서 -->
								{#if myLayerOn && getTasksForDay(day).length > 0}
									<span class="my-dot"></span>
								{/if}
								{#if orbitLayerOn && getOrbitSchedulesForDay(day).length > 0}
									<span class="orbit-strip"></span>
								{/if}
							</div>
						{/if}
					{/each}
				</div>
			{/if}
		{:else if selectedDay !== null}
			<div class="day-view">
				<div class="day-view-header">
					<button class="back-btn" onclick={backToMonth}>‹ {formatDayTitle(year, month, selectedDay)}</button>
				</div>

				<div class="week-strip">
					{#each getWeekDates(year, month, selectedDay) as d (d.toISOString())}
						<button
							class="week-strip-day {isSameSelectedDate(d) ? 'selected' : ''}"
							onclick={() => selectDateInStrip(d)}
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

				{#if dayViewError}
					<p class="error-msg">{dayViewError}</p>
				{/if}

				{#if myLayerOn}
					<section class="day-section">
						<h3 class="day-section-title">
							<span class="my-dot"></span> My Schedule ({getTasksForDay(selectedDay).length})
						</h3>
						{#if getTasksForDay(selectedDay).length === 0}
							<p class="empty-msg">일정이 없습니다.</p>
						{:else}
							<ul class="day-item-list">
								{#each getTasksForDay(selectedDay) as task (task.id)}
									<li class="day-item {task.is_completed ? 'completed' : ''}">
										<span class="day-item-time">
											{formatTime(task.start_at)} – {formatTime(task.end_at)}
										</span>
										{#if isOwner}
											<button
												class="btn-toggle"
												onclick={() => handleToggleTask(task)}
												title={task.is_completed ? '완료 취소' : '완료로 표시'}
											>
												{task.is_completed ? '✓' : '○'}
											</button>
										{/if}
										<span class="day-item-title">{task.title}</span>
										{#if isOwner}
											<button
												class="btn-delete"
												onclick={() => handleDeleteTask(task.id)}
												title="삭제">✕</button
											>
										{/if}
									</li>
								{/each}
							</ul>
						{/if}
						{#if isOwner}
							<button class="btn-open-add" onclick={(e) => openAddModal(selectedDay, e)}>
								<span>+</span> 할 일 추가
							</button>
						{/if}
					</section>
				{/if}

				{#if orbitLayerOn && isOwner}
					<section class="day-section">
						<h3 class="day-section-title">
							<span class="orbit-dot"></span> Orbit Schedule ({getOrbitSchedulesForDay(selectedDay)
								.length})
						</h3>
						{#if orbitLaneItems.length === 0}
							<p class="empty-msg">Orbit 일정이 없습니다.</p>
						{:else}
							<div class="orbit-timeline">
								<div class="orbit-timeline-ruler">
									{#each [0, 6, 12, 18, 24] as hour}
										<span class="ruler-mark" style="left: {(hour / 24) * 100}%">
											{String(hour).padStart(2, '0')}:00
										</span>
									{/each}
								</div>
								<div class="orbit-timeline-track" style="height: {orbitLaneCount * 34}px">
									{#each orbitLaneItems as item (item.schedule.id)}
										<button
											class="orbit-lane-bar {activeOrbitLane?.schedule.id === item.schedule.id ? 'expanded' : ''}"
											style="left: {(item.startMin / DAY_MINUTES) * 100}%;
												width: {((item.endMin - item.startMin) / DAY_MINUTES) * 100}%;
												top: {item.lane * 34}px;"
											onclick={() => openOrbitDetail(item)}
											title="{formatTime(item.schedule.start_at)} – {formatTime(item.schedule.end_at)}"
										></button>
									{/each}
								</div>
							</div>
						{/if}
					</section>
				{/if}
			</div>
		{/if}
	</div>
</div>

{#if addDay !== null}
	<AddTaskModal
		day={addDay}
		{year}
		{month}
		onClose={() => (addDay = null)}
		onCreated={handleTaskCreated}
	/>
{/if}

{#if activeOrbitLane}
	<OrbitDetailModal
		schedules={orbitOverlapGroup}
		rangeStartMin={activeOrbitLane.startMin}
		rangeEndMin={activeOrbitLane.endMin}
		densityBuckets={orbitDensityBuckets}
		bucketMinutes={BUCKET_MINUTES}
		onClose={closeOrbitDetail}
	/>
{/if}