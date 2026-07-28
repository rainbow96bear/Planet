<script lang="ts">
    import type { Task } from '$lib/types/task'
    import { createTask } from '$lib/api/task'
    let {
        day,
        year,
        month,
        onClose,
        onCreated
    }: {
        day: number
        year: number
        month: number
        onClose: () => void
        onCreated: (task: Task) => void
    } = $props()
    let title = $state('')
    let isPublic = $state(true)
    let loading = $state(false)
    let error = $state('')
    let inputEl = $state<HTMLInputElement | null>(null)
    $effect(() => {
        inputEl?.focus()
    })
    async function handleCreate() {
        if (!title.trim()) return
        loading = true
        error = ''
        try {
            const task = await createTask({
                title: title.trim(),
                date: `${year}-${String(month).padStart(2, '0')}-${String(day).padStart(2, '0')}T00:00:00Z`,
                is_public: isPublic
            })
            onCreated(task)
            onClose()
        } catch {
            error = '추가에 실패했습니다.'
        } finally {
            loading = false
        }
    }
    function handleBackdropClick(e: MouseEvent) {
        if (e.target === e.currentTarget) onClose()
    }
    function handleKeydown(e: KeyboardEvent) {
        if (e.key === 'Enter') handleCreate()
        if (e.key === 'Escape') onClose()
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
                <span class="label">할 일 추가</span>
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
                <span class="toggle-label">공개</span>
                <button
                    class="toggle {isPublic ? 'on' : 'off'}"
                    type="button"
                    onclick={() => isPublic = !isPublic}
                    aria-pressed={isPublic}
                >
                    <span class="toggle-thumb"></span>
                </button>
            </label>
            <div class="modal-actions">
                <button class="btn-cancel" onclick={onClose} disabled={loading}>취소</button>
                <button
                    class="btn-submit"
                    onclick={handleCreate}
                    disabled={loading || !title.trim()}
                >
                    {#if loading}
                        <span class="spinner"></span>
                    {:else}
                        추가
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
        background: rgba(34, 58, 94, .32);
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
    /* ==========================
    Header
    ========================== */
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
        font-size: .8rem;
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
    /* ==========================
    Input
    ========================== */
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
        box-shadow: 0 0 0 3px rgba(var(--planet-primary-rgb), .12);
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
        font-size: .75rem;
        pointer-events: none;
    }
    /* ==========================
    Toggle
    ========================== */
    .toggle-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: var(--space-xl);
    }
    .toggle-label {
        color: var(--text-secondary);
        font-size: .875rem;
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
    /* ==========================
    Error
    ========================== */
    .error-msg {
        margin: 0 0 var(--space-md);
        padding: var(--space-sm) var(--space-md);
        background: rgba(var(--danger-rgb), .06);
        border: 1px solid rgba(var(--danger-rgb), .12);
        border-radius: var(--radius-md);
        color: var(--danger);
        font-size: .875rem;
    }
    /* ==========================
    Actions
    ========================== */
    .modal-actions {
        display: flex;
        gap: var(--space-sm);
    }
    .btn-cancel,
    .btn-submit {
        height: 40px;
        border-radius: var(--radius-md);
        font: inherit;
        font-size: .875rem;
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
        opacity: .5;
        cursor: not-allowed;
    }
    /* ==========================
    Spinner
    ========================== */
    .spinner {
        width: 14px;
        height: 14px;
        border: 2px solid currentColor;
        border-top-color: transparent;
        border-radius: 50%;
        animation: spin .6s linear infinite;
    }
    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }
    /* ==========================
    Responsive
    ========================== */
    @media (max-width:520px){
        .modal{
            padding: var(--space-xl);
        }
        .modal-actions{
            flex-direction:column-reverse;
        }
        .btn-cancel,
        .btn-submit{
            width:100%;
            flex:none;
        }
    }
</style>