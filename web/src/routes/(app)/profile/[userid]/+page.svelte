<script lang="ts">
    import { getTasksByMonth } from '$lib/api/task'
    import { follow, unfollow } from '$lib/api/user'
    import type { PageData } from './$types'
    import type { Task } from '$lib/types/task'
    import TaskModal from '$lib/components/TaskModal.svelte'
    import AddTaskModal from '$lib/components/AddTaskModal.svelte'
    import './page.css'
	import { page } from '$app/stores';

    let { data }: { data: PageData } = $props()
    const userid = $derived($page.params.userid || "")
    const isOwner = $derived(userid === data.me?.userid)
    let tasks = $state<Task[]>(Array.isArray(data.tasks) ? data.tasks : [])
    let year = $state(data.year)
    let month = $state(data.month)
    let loading = $state(false)
    let isFollowing = $state(data.profileUser.is_following ?? false)
    let followLoading = $state(false)

    let selectedDay = $state<number | null>(null)
    let addDay = $state<number | null>(null)

    $effect(() => {
        tasks = Array.isArray(data.tasks) ? data.tasks : []
    })

    async function prevMonth() {
        if (month === 1) { year -= 1; month = 12 }
        else { month -= 1 }
        loading = true
        tasks = await getTasksByMonth(userid, year, month)
        loading = false
    }

    async function nextMonth() {
        if (month === 12) { year += 1; month = 1 }
        else { month += 1 }
        loading = true
        tasks = await getTasksByMonth(userid, year, month)
        loading = false
    }

    async function handleFollow() {
        followLoading = true
        try {
            await follow(userid)
            isFollowing = true
        } catch (e) {
            console.error(e)
        } finally {
            followLoading = false
        }
    }

    async function handleUnfollow() {
        followLoading = true
        try {
            await unfollow(userid)
            isFollowing = false
        } catch (e) {
            console.error(e)
        } finally {
            followLoading = false
        }
    }

    // 삭제 콜백: 페이지 tasks에서 제거
    function handleTaskDeleted(taskId: string) {
        tasks = tasks.filter(t => t.id !== taskId)
    }

    // 토글 콜백: 페이지 tasks에서 완료 상태 반전
    function handleTaskToggled(taskId: string) {
        tasks = tasks.map(t =>
            t.id === taskId ? { ...t, is_completed: !t.is_completed } : t
        )
    }

    const DAYS = ['일', '월', '화', '수', '목', '금', '토']

    function getCalendarDays(year: number, month: number) {
        const firstDay = new Date(year, month - 1, 1).getDay()
        const lastDate = new Date(year, month, 0).getDate()
        const days: (number | null)[] = []
        for (let i = 0; i < firstDay; i++) days.push(null)
        for (let i = 1; i <= lastDate; i++) days.push(i)
        return days
    }

    function getTasksForDay(day: number) {
        return tasks.filter(t => new Date(t.date).getDate() === day)
    }

    function getTodayDay() {
        const now = new Date()
        if (now.getFullYear() === year && now.getMonth() + 1 === month) {
            return now.getDate()
        }
        return null
    }

    function openTaskModal(day: number) {
        selectedDay = day
    }

    function openAddModal(day: number, e: MouseEvent) {
        e.stopPropagation()
        addDay = day
    }

    function handleAddClick() {
        if (selectedDay !== null) {
            addDay = selectedDay
        }
    }

    function handleTaskCreated(task: Task) {
        tasks = [...tasks, task]
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
        </div>

        <div class="profile-actions">
            {#if isOwner}
                <a href="/settings/profile" class="action-btn secondary">
                    <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                        <path d="M11 4H4a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-7"/>
                        <path d="M18.5 2.5a2.121 2.121 0 0 1 3 3L12 15l-4 1 1-4 9.5-9.5z"/>
                    </svg>
                    프로필 수정
                </a>
            {:else}
                {#if followLoading}
                    <button class="action-btn primary" disabled>
                        <span class="spinner"></span>
                    </button>
                {:else if isFollowing}
                    <button class="action-btn following" onclick={handleUnfollow}>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <polyline points="20 6 9 17 4 12"/>
                        </svg>
                        팔로잉
                    </button>
                {:else}
                    <button class="action-btn primary" onclick={handleFollow}>
                        <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round">
                            <path d="M16 21v-2a4 4 0 0 0-4-4H6a4 4 0 0 0-4 4v2"/>
                            <circle cx="9" cy="7" r="4"/>
                            <line x1="19" y1="8" x2="19" y2="14"/>
                            <line x1="22" y1="11" x2="16" y2="11"/>
                        </svg>
                        팔로우
                    </button>
                {/if}
            {/if}
        </div>
    </div>

    <div class="calendar-card">
        <div class="calendar-nav">
            <button class="nav-btn" onclick={prevMonth} disabled={loading}>◀</button>
            <span class="calendar-title">{year}년 {month}월</span>
            <button class="nav-btn" onclick={nextMonth} disabled={loading}>▶</button>
        </div>

        {#if loading}
            <div class="calendar-loading">불러오는 중...</div>
        {:else}
            <div class="calendar-grid">
                {#each DAYS as day}
                    <div class="calendar-day-header">{day}</div>
                {/each}

                {#each getCalendarDays(year, month) as day, i}
                    {@const isToday = day === getTodayDay()}
                    <button
                        class="calendar-cell {day === null ? 'empty' : ''} {isToday ? 'today' : ''}"
                        onclick={() => day && openTaskModal(day)}
                        disabled={day === null}
                    >
                        {#if day !== null}
                            <div class="cell-top">
                                <span class="day-number {i % 7 === 0 ? 'sunday' : i % 7 === 6 ? 'saturday' : ''} {isToday ? 'today-number' : ''}">
                                    {day}
                                </span>
                                {#if isOwner}
                                    <button
                                        class="add-task-btn"
                                        onclick={(e) => openAddModal(day, e)}
                                        title="할 일 추가"
                                        aria-label="할 일 추가"
                                    >+</button>
                                {/if}
                            </div>
                            <div class="task-list">
                                {#each getTasksForDay(day) as task}
                                    <div class="task-chip {task.is_completed ? 'completed' : ''}">
                                        {task.title}
                                    </div>
                                {/each}
                            </div>
                        {/if}
                    </button>
                {/each}
            </div>
        {/if}
    </div>
</div>

{#if selectedDay !== null && addDay === null}
    <TaskModal
        day={selectedDay}
        {year}
        {month}
        tasks={getTasksForDay(selectedDay)}
        {isOwner}
        onClose={() => selectedDay = null}
        onAddClick={handleAddClick}
        onDeleted={handleTaskDeleted}
        onToggled={handleTaskToggled}
    />
{/if}

{#if addDay !== null}
    <AddTaskModal
        day={addDay}
        {year}
        {month}
        onClose={() => addDay = null}
        onCreated={handleTaskCreated}
    />
{/if}