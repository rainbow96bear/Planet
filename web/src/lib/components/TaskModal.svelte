<script lang="ts">
	import type { Task } from '$lib/types/task';
	import { deleteTask, toggleTask } from '$lib/api/task';
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
		day: number;
		year: number;
		month: number;
		tasks: Task[];
		isOwner: boolean;
		onClose: () => void;
		onAddClick?: () => void;
		onDeleted?: (taskId: string) => void;
		onToggled?: (taskId: string) => void;
	} = $props();
	let error = $state('');
	async function handleDelete(taskId: string) {
		try {
			await deleteTask(taskId);
			tasks = tasks.filter((t) => t.id !== taskId);
			onDeleted?.(taskId);
		} catch {
			error = '삭제에 실패했습니다.';
		}
	}
	async function handleToggle(task: Task) {
		try {
			await toggleTask(task.id);
			tasks = tasks.map((t) => (t.id === task.id ? { ...t, is_completed: !t.is_completed } : t));
			onToggled?.(task.id);
		} catch (err) {
			console.error(err);
			error = '변경에 실패했습니다.';
		}
	}
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onClose();
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
							<button class="btn-delete" onclick={() => handleDelete(task.id)} title="삭제"
								>✕</button
							>
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
		display: flex;
		align-items: center;
		justify-content: center;
		padding: 1rem;
		background: rgba(34, 58, 94, 0.45);
		backdrop-filter: blur(2px);
		z-index: var(--z-modal);
	}
	.modal {
		width: 100%;
		max-width: 420px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		padding: var(--space-xl);
		box-shadow: var(--shadow-xl);
	}
	.modal-header {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-lg);
	}
	.modal-date {
		font-size: 1rem;
		font-weight: 700;
		color: var(--text-primary);
	}
	.modal-close {
		display: inline-flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		border: none;
		border-radius: var(--radius-md);
		background: transparent;
		color: var(--text-secondary);
		cursor: pointer;
		transition:
			background var(--transition-normal),
			color var(--transition-normal);
	}
	.modal-close:hover {
		background: var(--surface-hover);
		color: var(--text-primary);
	}
	/* ---------- Error ---------- */
	.error-msg {
		margin: 0 0 var(--space-md);
		padding: 0.65rem 0.85rem;
		font-size: 0.8rem;
		color: var(--danger);
		background: rgba(var(--danger-rgb), 0.08);
		border: 1px solid rgba(var(--danger-rgb), 0.18);
		border-radius: var(--radius-md);
	}
	/* ---------- List ---------- */
	.task-list {
		display: flex;
		flex-direction: column;
		gap: var(--space-sm);
		list-style: none;
		margin: 0 0 var(--space-lg);
		padding: 0;
	}
	.task-item {
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
	.task-item:hover {
		border-color: var(--planet-primary);
	}
	.task-item.completed {
		color: var(--text-muted);
	}
	.task-item.completed .task-title {
		text-decoration: line-through;
	}
	.task-title {
		flex: 1;
		min-width: 0;
		overflow: hidden;
		white-space: nowrap;
		text-overflow: ellipsis;
		font-size: 0.9rem;
		color: inherit;
	}
	/* ---------- Toggle ---------- */
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
	/* ---------- Delete ---------- */
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
	/* ---------- Empty ---------- */
	.empty-msg {
		margin: 0;
		padding: var(--space-2xl) var(--space-md);
		text-align: center;
		font-size: 0.85rem;
		color: var(--text-secondary);
	}
	/* ---------- Add ---------- */
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
	/* ---------- Responsive ---------- */
	@media (max-width: 520px) {
		.modal {
			padding: var(--space-lg);
		}
		.task-item {
			padding: 0.7rem 0.8rem;
		}
		.task-title {
			font-size: 0.875rem;
		}
	}
</style>
