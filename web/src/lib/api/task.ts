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