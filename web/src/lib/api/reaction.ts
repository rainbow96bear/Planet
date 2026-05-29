import { ReactionType, ToggleReactionBody } from "$lib/types/reaction"

export const addTaskReaction = async (
    taskId: string,
    type: ReactionType
) => {
    const res = await fetch(`/api/v1/tasks/${taskId}/reactions`, {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            type,
        } satisfies ToggleReactionBody),
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const removeTaskReaction = async (
    taskId: string,
    type: ReactionType
) => {
    const res = await fetch(`/api/v1/tasks/${taskId}/reactions`, {
        method: 'DELETE',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            type,
        } satisfies ToggleReactionBody),
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}