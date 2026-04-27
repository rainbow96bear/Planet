<script lang="ts">
    import { getTasksByMonth } from '$lib/api/task'
    import type { PageData } from './$types'

    let { data }: { data: PageData } = $props()
    
    const username = $page.params.username
    let tasks = $state(data.tasks)
    let year = $state(data.year)
    let month = $state(data.month)

    async function prevMonth() {
        if (month === 1) {
            year -= 1
            month = 12
        } else {
            month -= 1
        }
        tasks = await getTasksByMonth(username, year, month)
    }

    async function nextMonth() {
        if (month === 12) {
            year += 1
            month = 1
        } else {
            month += 1
        }
        tasks = await getTasksByMonth(username, year, month)
    }
</script>

<div>
    <h1>{data.user.nickname}</h1>

    <div class="calendar-header">
        <button onclick={prevMonth}>◀</button>
        <span>{year}년 {month}월</span>
        <button onclick={nextMonth}>▶</button>
    </div>

    <div class="calendar">
        <!-- 달력 렌더링 -->
        {#each tasks as task}
            <div>{task.title}</div>
        {/each}
    </div>
</div>