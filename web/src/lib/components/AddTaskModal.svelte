<script lang="ts">
	import type { Task } from '$lib/types/task';
	import { createTask, updateTask } from '$lib/api/task';
	let {
		day,
		year,
		month,
		task,
		onClose,
		onSaved
	}: {
		day: number;
		year: number;
		month: number;
		// task가 있으면 수정 모드, 없으면 기존과 동일한 생성 모드.
		task?: Task;
		onClose: () => void;
		onSaved: (task: Task) => void;
	} = $props();

	const isEditMode = task !== undefined;

	function timeOf(iso: string) {
		const d = new Date(iso);
		return `${String(d.getHours()).padStart(2, '0')}:${String(d.getMinutes()).padStart(2, '0')}`;
	}

	let title = $state(task?.title ?? '');
	let isPublic = $state(task?.is_public ?? true);
	let allDay = $state(
		task ? timeOf(task.start_at) === '00:00' && timeOf(task.end_at) === '23:59' : true
	);
	let startTime = $state(task && !allDay ? timeOf(task.start_at) : '09:00');
	let endTime = $state(task && !allDay ? timeOf(task.end_at) : '10:00');
	let loading = $state(false);
	let error = $state('');
	let inputEl = $state<HTMLInputElement | null>(null);
	$effect(() => {
		inputEl?.focus();
	});

	const dateStr = $derived(
		`${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}`
	);

	async function handleSubmit() {
		if (!title.trim()) return;

		if (!allDay && startTime >= endTime) {
			error = '종료 시간은 시작 시간보다 늦어야 합니다.';
			return;
		}

		loading = true;
		error = '';
		try {
			const body = {
				title: title.trim(),
				start_at: allDay ? `${dateStr}T00:00:00Z` : `${dateStr}T${startTime}:00Z`,
				end_at: allDay ? `${dateStr}T23:59:59Z` : `${dateStr}T${endTime}:00Z`,
				is_public: isPublic
			};
			const saved = isEditMode ? await updateTask(task!.id, body) : await createTask(body);
			onSaved(saved);
			onClose();
		} catch {
			error = isEditMode ? '수정에 실패했습니다.' : '추가에 실패했습니다.';
		} finally {
			loading = false;
		}
	}
	function handleBackdropClick(e: MouseEvent) {
		if (e.target === e.currentTarget) onClose();
	}
	function handleKeydown(e: KeyboardEvent) {
		if (e.key === 'Enter') handleSubmit();
		if (e.key === 'Escape') onClose();
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
			<span class="modal-date">
				<span class="label">{isEditMode ? '할 일 수정' : '할 일 추가'}</span>
				<span class="sub">{year}년 {month}월 {day}일</span>
			</span>
			<button class="modal-close" onclick={onClose}>✕</button>
		</div>
		<div class="modal-body">
			{#if error}
				<p class="error-msg">{error}</p>
			{/if}
			<div class="input-group">
				<input
					bind:this={inputEl}
					bind:value={title}
					type="text"
					placeholder="할 일을 입력하세요"
					onkeydown={handleKeydown}
					maxlength={100}
					disabled={loading}
				/>
				<span class="char-count">{title.length}/100</span>
			</div>

			<label class="toggle-row">
				<span class="toggle-label">하루 종일</span>
				<button
					class="toggle {allDay ? 'on' : ''}"
					type="button"
					onclick={() => (allDay = !allDay)}
					aria-pressed={allDay}
					aria-label="하루 종일 여부 전환"
				>
					<span class="toggle-thumb"></span>
				</button>
			</label>

			{#if !allDay}
				<div class="time-row">
					<label class="time-field">
						<span class="time-field-label">시작</span>
						<input type="time" bind:value={startTime} disabled={loading} />
					</label>
					<span class="time-sep">–</span>
					<label class="time-field">
						<span class="time-field-label">종료</span>
						<input type="time" bind:value={endTime} disabled={loading} />
					</label>
				</div>
			{/if}

			<label class="toggle-row">
				<span class="toggle-label">공개</span>
				<button
					class="toggle {isPublic ? 'on' : 'off'}"
					type="button"
					onclick={() => (isPublic = !isPublic)}
					aria-pressed={isPublic}
					aria-label="공개 여부 전환"
				>
					<span class="toggle-thumb"></span>
				</button>
			</label>
			<div class="modal-actions">
				<button class="btn-cancel" onclick={onClose} disabled={loading}>취소</button>
				<button class="btn-submit" onclick={handleSubmit} disabled={loading || !title.trim()}>
					{#if loading}
						<span class="spinner"></span>
					{:else}
						{isEditMode ? '저장' : '추가'}
					{/if}
				</button>
			</div>
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
		background: rgba(34, 58, 94, 0.32);
		backdrop-filter: blur(2px);
		z-index: var(--z-modal);
	}
	.modal {
		width: 100%;
		max-width: 400px;
		padding: var(--space-2xl);
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-lg);
		box-sizing: border-box;
	}
	.modal-header {
		display: flex;
		justify-content: space-between;
		align-items: flex-start;
		margin-bottom: var(--space-xl);
	}
	.modal-date {
		display: flex;
		flex-direction: column;
		gap: var(--space-2xs);
	}
	.label {
		color: var(--text-primary);
		font-size: 1.1rem;
		font-weight: 700;
	}
	.sub {
		color: var(--text-secondary);
		font-size: 0.8rem;
	}
	.modal-close {
		display: flex;
		align-items: center;
		justify-content: center;
		width: 32px;
		height: 32px;
		padding: 0;
		background: transparent;
		border: none;
		border-radius: var(--radius-md);
		color: var(--text-secondary);
		cursor: pointer;
		transition:
			background var(--transition-fast),
			color var(--transition-fast);
	}
	.modal-close:hover {
		background: var(--surface-hover);
		color: var(--text-primary);
	}
	.input-group {
		position: relative;
		margin-bottom: var(--space-lg);
	}
	.input-group input {
		width: 100%;
		height: 42px;
		padding: 0 54px 0 14px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		color: var(--text-primary);
		font: inherit;
		box-sizing: border-box;
		transition:
			border-color var(--transition-fast),
			background var(--transition-fast),
			box-shadow var(--transition-fast);
	}
	.input-group input::placeholder {
		color: var(--text-muted);
	}
	.input-group input:focus {
		outline: none;
		border-color: var(--planet-primary);
		box-shadow: 0 0 0 3px rgba(var(--planet-primary-rgb), 0.12);
	}
	.input-group input:disabled {
		background: var(--surface-hover);
		color: var(--text-muted);
	}
	.char-count {
		position: absolute;
		top: 50%;
		right: 14px;
		transform: translateY(-50%);
		color: var(--text-muted);
		font-size: 0.75rem;
		pointer-events: none;
	}
	.toggle-row {
		display: flex;
		align-items: center;
		justify-content: space-between;
		margin-bottom: var(--space-md);
	}
	.toggle-label {
		color: var(--text-secondary);
		font-size: 0.875rem;
		font-weight: 500;
	}
	.toggle {
		position: relative;
		width: 40px;
		height: 22px;
		padding: 0;
		background: var(--border);
		border: none;
		border-radius: 999px;
		cursor: pointer;
		transition: background var(--transition-fast);
	}
	.toggle.on {
		background: var(--planet-primary);
	}
	.toggle-thumb {
		position: absolute;
		top: 2px;
		left: 2px;
		width: 18px;
		height: 18px;
		background: var(--text-on-dark);
		border-radius: 50%;
		transition: transform var(--transition-fast);
	}
	.toggle.on .toggle-thumb {
		transform: translateX(18px);
	}
	.time-row {
		display: flex;
		align-items: flex-end;
		gap: var(--space-sm);
		margin-bottom: var(--space-xl);
	}
	.time-field {
		flex: 1;
		display: flex;
		flex-direction: column;
		gap: 4px;
	}
	.time-field-label {
		color: var(--text-muted);
		font-size: 0.7rem;
		font-weight: 600;
	}
	.time-field input[type='time'] {
		width: 100%;
		height: 38px;
		padding: 0 10px;
		background: var(--surface);
		border: 1px solid var(--border);
		border-radius: var(--radius-md);
		color: var(--text-primary);
		font: inherit;
		box-sizing: border-box;
	}
	.time-field input[type='time']:focus {
		outline: none;
		border-color: var(--planet-primary);
		box-shadow: 0 0 0 3px rgba(var(--planet-primary-rgb), 0.12);
	}
	.time-sep {
		padding-bottom: 9px;
		color: var(--text-muted);
	}
	.error-msg {
		margin: 0 0 var(--space-md);
		padding: var(--space-sm) var(--space-md);
		background: rgba(var(--danger-rgb), 0.06);
		border: 1px solid rgba(var(--danger-rgb), 0.12);
		border-radius: var(--radius-md);
		color: var(--danger);
		font-size: 0.875rem;
	}
	.modal-actions {
		display: flex;
		gap: var(--space-sm);
	}
	.btn-cancel,
	.btn-submit {
		height: 40px;
		border-radius: var(--radius-md);
		font: inherit;
		font-size: 0.875rem;
		font-weight: 600;
		transition:
			background var(--transition-fast),
			border-color var(--transition-fast),
			color var(--transition-fast),
			opacity var(--transition-fast);
		cursor: pointer;
	}
	.btn-cancel {
		flex: 1;
		background: var(--surface);
		border: 1px solid var(--border);
		color: var(--text-secondary);
	}
	.btn-cancel:hover:not(:disabled) {
		background: var(--surface-hover);
		color: var(--text-primary);
	}
	.btn-submit {
		flex: 2;
		display: flex;
		align-items: center;
		justify-content: center;
		background: var(--planet-primary);
		border: 1px solid var(--planet-primary);
		color: var(--text-on-dark);
	}
	.btn-submit:hover:not(:disabled) {
		background: var(--planet-primary-hover);
		border-color: var(--planet-primary-hover);
	}
	.btn-cancel:disabled,
	.btn-submit:disabled {
		opacity: 0.5;
		cursor: not-allowed;
	}
	.spinner {
		width: 14px;
		height: 14px;
		border: 2px solid currentColor;
		border-top-color: transparent;
		border-radius: 50%;
		animation: spin 0.6s linear infinite;
	}
	@keyframes spin {
		to {
			transform: rotate(360deg);
		}
	}
	@media (max-width: 520px) {
		.modal {
			padding: var(--space-xl);
		}
		.modal-actions {
			flex-direction: column-reverse;
		}
		.btn-cancel,
		.btn-submit {
			width: 100%;
			flex: none;
		}
	}
</style>
