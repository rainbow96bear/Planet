<script lang="ts">
    import { getTasksByMonth } from '$lib/api/task'
    import { page } from '$app/stores'
    import type { PageData } from './$types'
    import TaskModal from '$lib/components/TaskModal.svelte'
    import './page.css'

    let { data }: { data: PageData } = $props()

    const username = $page.params.username
    const isOwner = data.user.username === data.me?.username

    let tasks = $state(data.tasks)
    let year = $state(data.year)
    let month = $state(data.month)
    let loading = $state(false)
    let selectedDay = $state<number | null>(null)

    async function prevMonth() {
        if (month === 1) { year -= 1; month = 12 }
        else { month -= 1 }
        loading = true
        tasks = await getTasksByMonth(username, year, month)
        loading = false
    }

    async function nextMonth() {
        if (month === 12) { year += 1; month = 1 }
        else { month += 1 }
        loading = true
        tasks = await getTasksByMonth(username, year, month)
        loading = false
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

    function openModal(day: number) {
        selectedDay = day
    }

    function closeModal() {
        selectedDay = null
    }
</script>

<div class="profile-container">
    <div class="profile-header">
        <div class="profile-avatar">🪐</div>
        <div class="profile-info">
            <h1 class="profile-nickname">{data.user.nickname}</h1>
            <span class="profile-username">@{username}</span>
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
                    <button
                        class="calendar-cell {day === null ? 'empty' : ''}"
                        onclick={() => day && openModal(day)}
                    >
                        {#if day !== null}
                            <span class="day-number {i % 7 === 0 ? 'sunday' : i % 7 === 6 ? 'saturday' : ''}">
                                {day}
                            </span>
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

{#if selectedDay}
    <TaskModal
        day={selectedDay}
        {year}
        {month}
        {username}
        tasks={getTasksForDay(selectedDay)}
        {isOwner}
        onClose={closeModal}
    />
{/if}