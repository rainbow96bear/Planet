<script lang="ts">
    import type { Task } from '$lib/types/task'
    import { deleteTask, toggleTask } from '$lib/api/task'

    let {
        day,
        year,
        month,
        tasks = $bindable(),
        isOwner,
        onClose,
        onAddClick,
        onDeleted,
        onToggled
    }: {
        day: number
        year: number
        month: number
        tasks: Task[]
        isOwner: boolean
        onClose: () => void
        onAddClick?: () => void
        onDeleted?: (taskId: number) => void
        onToggled?: (taskId: number) => void
    } = $props()

    let error = $state('')

    async function handleDelete(taskId: number) {
        try {
            await deleteTask(taskId)
            tasks = tasks.filter(t => t.id !== taskId)
            onDeleted?.(taskId)
        } catch {
            error = '삭제에 실패했습니다.'
        }
    }

    async function handleToggle(task: Task) {
        try {
            await toggleTask(task.id)
            tasks = tasks.map(t =>
                t.id === task.id ? { ...t, is_completed: !t.is_completed } : t
            )
            onToggled?.(task.id)
        } catch {
            error = '변경에 실패했습니다.'
        }
    }

    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose()
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
                        {#if isOwner}
                            <button
                                class="btn-toggle"
                                onclick={() => handleToggle(task)}
                                title={task.is_completed ? '완료 취소' : '완료로 표시'}
                            >
                                {task.is_completed ? '✓' : '○'}
                            </button>
                        {/if}
                        <span class="task-title">{task.title}</span>
                        {#if isOwner}
                            <button
                                class="btn-delete"
                                onclick={() => handleDelete(task.id)}
                                title="삭제"
                            >✕</button>
                        {/if}
                    </li>
                {:else}
                    <p class="empty-msg">할 일이 없습니다.</p>
                {/each}
            </ul>

            {#if isOwner}
                <button class="btn-open-add" onclick={onAddClick}>
                    <span>+</span> 할 일 추가
                </button>
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
        align-items: center;
        gap: 0.5rem;
        background: #1a1a26;
        border: 1px solid #a89fd420;
        border-radius: 8px;
        padding: 0.6rem 0.8rem;
        font-size: 0.9rem;
        color: #ffffff;
    }

    .task-item.completed {
        color: #a89fd4;
    }

    .task-item.completed .task-title {
        text-decoration: line-through;
    }

    .task-title {
        flex: 1;
        min-width: 0;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
    }

    .btn-toggle {
        background: none;
        border: none;
        color: #a89fd4;
        font-size: 0.85rem;
        cursor: pointer;
        width: 20px;
        flex-shrink: 0;
        transition: color 0.2s;
    }

    .btn-toggle:hover {
        color: #b2ede6;
    }

    .btn-delete {
        background: none;
        border: none;
        color: #a89fd450;
        font-size: 0.75rem;
        cursor: pointer;
        flex-shrink: 0;
        transition: color 0.2s;
    }

    .btn-delete:hover {
        color: #e05c5c;
    }

    .empty-msg {
        font-size: 0.85rem;
        color: #a89fd4;
        text-align: center;
        padding: 1.5rem 0;
        margin: 0;
    }

    .btn-open-add {
        display: flex;
        align-items: center;
        justify-content: center;
        gap: 0.4rem;
        width: 100%;
        padding: 0.6rem;
        background: #b2ede610;
        border: 1px dashed #b2ede640;
        border-radius: 8px;
        color: #b2ede6;
        font-size: 0.85rem;
        font-weight: 600;
        cursor: pointer;
        transition: background 0.2s, border-color 0.2s;
    }

    .btn-open-add:hover {
        background: #b2ede620;
        border-color: #b2ede6;
    }

    .error-msg {
        font-size: 0.8rem;
        color: #e05c5c;
        margin-bottom: 0.75rem;
    }
</style>