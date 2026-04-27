<script lang="ts">
    import type { Task } from '$lib/types/task'
    import { createTask, deleteTask } from '$lib/api/task'

    let {
        day,
        year,
        month,
        username,
        tasks = $bindable(),
        isOwner,
        onClose
    }: {
        day: number
        year: number
        month: number
        username: string
        tasks: Task[]
        isOwner: boolean
        onClose: () => void
    } = $props()

    let newTitle = $state('')
    let loading = $state(false)
    let error = $state('')

    async function handleCreate() {
        if (!newTitle.trim()) return
        loading = true
        try {
            const task = await createTask({
                title: newTitle,
                date: `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`,
                is_public: true
            })
            tasks = [...tasks, task]
            newTitle = ''
        } catch {
            error = '추가에 실패했습니다.'
        } finally {
            loading = false
        }
    }

    async function handleDelete(taskId: number) {
        try {
            await deleteTask(taskId)
            tasks = tasks.filter(t => t.id !== taskId)
        } catch {
            error = '삭제에 실패했습니다.'
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose()
    }
</script>

<div class="modal-backdrop" onclick={handleBackdropClick}>
    <div class="modal">
        <div class="modal-header">
            <span class="modal-date">{year}년 {month}월 {day}일</span>
            <button class="modal-close" onclick={onClose}>✕</button>
        </div>

        <div class="modal-body">
            {#if error}
                <p class="error-msg">{error}</p>
            {/if}

            <ul class="task-list">
                {#each tasks as task}
                    <li class="task-item {task.is_completed ? 'completed' : ''}">
                        <span>{task.title}</span>
                        {#if isOwner}
                            <button class="btn-delete" onclick={() => handleDelete(task.id)}>✕</button>
                        {/if}
                    </li>
                {:else}
                    <p class="empty-msg">할 일이 없습니다.</p>
                {/each}
            </ul>

            {#if isOwner}
                <div class="task-input">
                    <input
                        type="text"
                        bind:value={newTitle}
                        placeholder="할 일을 입력하세요"
                        onkeydown={(e) => e.key === 'Enter' && handleCreate()}
                    />
                    <button class="btn-add" onclick={handleCreate} disabled={loading}>
                        {loading ? '...' : '+'}
                    </button>
                </div>
            {/if}
        </div>
    </div>
</div>

<style>
    .modal-backdrop {
        position: fixed;
        inset: 0;
        background: #00000080;
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 100;
    }

    .modal {
        background: #12121a;
        border: 1px solid #a89fd430;
        border-radius: 16px;
        width: 100%;
        max-width: 420px;
        padding: 1.5rem;
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        margin-bottom: 1.25rem;
    }

    .modal-date {
        font-size: 1rem;
        font-weight: 600;
        color: #ffffff;
    }

    .modal-close {
        background: none;
        border: none;
        color: #a89fd4;
        font-size: 1rem;
        cursor: pointer;
        transition: color 0.2s;
    }

    .modal-close:hover {
        color: #ffffff;
    }

    .task-list {
        list-style: none;
        padding: 0;
        margin: 0 0 1rem 0;
        display: flex;
        flex-direction: column;
        gap: 8px;
    }

    .task-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        background: #1a1a26;
        border: 1px solid #a89fd420;
        border-radius: 8px;
        padding: 0.6rem 0.8rem;
        font-size: 0.9rem;
        color: #ffffff;
    }

    .task-item.completed {
        color: #a89fd4;
        text-decoration: line-through;
    }

    .btn-delete {
        background: none;
        border: none;
        color: #a89fd4;
        font-size: 0.8rem;
        cursor: pointer;
        transition: color 0.2s;
    }

    .btn-delete:hover {
        color: #e05c5c;
    }

    .empty-msg {
        font-size: 0.85rem;
        color: #a89fd4;
        text-align: center;
        padding: 1rem 0;
    }

    .task-input {
        display: flex;
        gap: 8px;
    }

    .task-input input {
        flex: 1;
        background: #1a1a26;
        border: 1px solid #a89fd440;
        border-radius: 8px;
        padding: 0.6rem 0.9rem;
        color: #fff;
        font-size: 0.9rem;
        outline: none;
        transition: border-color 0.2s;
    }

    .task-input input:focus {
        border-color: #b2ede6;
    }

    .task-input input::placeholder {
        color: #3a3a50;
    }

    .btn-add {
        background: #b2ede6;
        border: none;
        border-radius: 8px;
        color: #0a0a0f;
        font-size: 1.2rem;
        font-weight: 700;
        width: 40px;
        cursor: pointer;
        transition: opacity 0.2s;
    }

    .btn-add:hover:not(:disabled) {
        opacity: 0.85;
    }

    .btn-add:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .error-msg {
        font-size: 0.8rem;
        color: #e05c5c;
        margin-bottom: 0.75rem;
    }
</style>