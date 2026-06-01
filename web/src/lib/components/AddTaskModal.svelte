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
        background: rgba(31, 41, 55, 0.4);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 110;
    }

    .modal {
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 16px;
        width: 100%;
        max-width: 380px;
        padding: 1.5rem;
        box-shadow: 0 12px 32px rgba(31, 41, 55, 0.08);
    }

    .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;
        margin-bottom: 1.25rem;
    }

    .modal-date {
        display: flex;
        flex-direction: column;
        gap: 2px;
    }

    .label {
        font-size: 1rem;
        font-weight: 700;
        color: var(--text-primary);
    }

    .sub {
        font-size: 0.78rem;
        color: var(--text-secondary);
    }

    .modal-close {
        background: none;
        border: none;
        color: var(--text-secondary);
        font-size: 1rem;
        cursor: pointer;
        transition: color 0.2s;
        flex-shrink: 0;
    }

    .modal-close:hover {
        color: var(--text-primary);
    }

    /* 입력 */
    .input-group {
        position: relative;
        margin-bottom: 1rem;
    }

    .input-group input {
        width: 100%;
        background: var(--bg);
        border: 1px solid var(--border);
        border-radius: 8px;
        padding: 0.65rem 3rem 0.65rem 0.9rem;
        color: var(--text-primary);
        font-size: 0.9rem;
        outline: none;
        box-sizing: border-box;
        transition: border-color 0.2s;
    }

    .input-group input:focus {
        border-color: var(--planet-primary);
        box-shadow: 0 0 0 3px rgba(124, 199, 201, 0.15);
    }

    .input-group input::placeholder {
        color: var(--text-secondary);
    }

    .input-group input:disabled {
        opacity: 0.5;
    }

    .char-count {
        position: absolute;
        right: 0.75rem;
        top: 50%;
        transform: translateY(-50%);
        font-size: 0.7rem;
        color: var(--text-secondary);
        pointer-events: none;
    }

    /* 공개 토글 */
    .toggle-row {
        display: flex;
        align-items: center;
        justify-content: space-between;
        margin-bottom: 1.25rem;
        cursor: pointer;
    }

    .toggle-label {
        font-size: 0.85rem;
        color: var(--text-secondary);
    }

    .toggle {
        position: relative;
        width: 36px;
        height: 20px;
        border-radius: 999px;
        border: none;
        cursor: pointer;
        transition: background 0.2s;
        padding: 0;
    }

    .toggle.on {
        background: var(--planet-primary);
    }

    .toggle.off {
        background: #d1d5db;
    }

    .toggle-thumb {
        position: absolute;
        top: 3px;
        left: 3px;
        width: 14px;
        height: 14px;
        border-radius: 50%;
        background: white;
        transition: transform 0.2s;
    }

    .toggle.on .toggle-thumb {
        transform: translateX(16px);
    }

    /* 버튼 영역 */
    .modal-actions {
        display: flex;
        gap: 8px;
    }

    .btn-cancel {
        flex: 1;
        padding: 0.6rem;
        background: var(--surface);
        border: 1px solid var(--border);
        border-radius: 8px;
        color: var(--text-secondary);
        font-size: 0.875rem;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
    }

    .btn-cancel:hover:not(:disabled) {
        border-color: var(--planet-secondary);
        color: var(--planet-secondary);
    }

    .btn-cancel:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .btn-submit {
        flex: 2;
        padding: 0.6rem;
        background: var(--planet-primary);
        border: none;
        border-radius: 8px;
        color: white;
        font-size: 0.875rem;
        font-weight: 700;
        cursor: pointer;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: background 0.2s;
    }

    .btn-submit:hover:not(:disabled) {
        background: var(--planet-primary-dark);
    }

    .btn-submit:disabled {
        opacity: 0.4;
        cursor: not-allowed;
    }

    .spinner {
        display: inline-block;
        width: 14px;
        height: 14px;
        border: 2px solid rgba(255, 255, 255, 0.35);
        border-top-color: white;
        border-radius: 50%;
        animation: spin 0.6s linear infinite;
    }

    @keyframes spin {
        to {
            transform: rotate(360deg);
        }
    }

    .error-msg {
        font-size: 0.8rem;
        color: var(--danger);
        margin-bottom: 0.75rem;
        padding: 0.6rem 0.8rem;
        background: rgba(239, 68, 68, 0.08);
        border: 1px solid rgba(239, 68, 68, 0.15);
        border-radius: 8px;
    }
</style>