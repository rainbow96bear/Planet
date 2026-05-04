import type { CreateTaskBody, CreateTaskResponse, GetTasksByMonthResponse, ToggleTaskResponse } from "$lib/types/task"

export const getTasksByMonth = async (
    userid: number,
    year: number,
    month: number
): Promise<GetTasksByMonthResponse> => {
    const res = await fetch(`/api/v1/users/${userid}/tasks?year=${year}&month=${month}`)

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const createTask = async (body: CreateTaskBody): Promise<CreateTaskResponse> => {
    console.log(body)
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

export const toggleTask = async (taskId: number): Promise<ToggleTaskResponse> => {
    const res = await fetch(`/api/v1/tasks/${taskId}/toggle`, {
        method: 'POST',
    })
 
    if (!res.ok) throw new Error('Failed to toggle task')
 
    return res.json()
}