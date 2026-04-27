import { GetTasksByMonthResponse } from "$lib/types/task";

export const getTasksByMonth = async (
    username: string,
    year: number,
    month: number
): Promise<GetTasksByMonthResponse> => {
    const res = await fetch(`/api/v1/users/${username}/tasks?year=${year}&month=${month}`)

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const createTask = async (body: CreateTaskBody): Promise<CreateTaskResponse> => {
    const res = await fetch('/api/v1/tasks', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const deleteTask = async (taskId: number): Promise<void> => {
    const res = await fetch(`/api/v1/tasks/${taskId}`, {
        method: 'DELETE',
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }
}