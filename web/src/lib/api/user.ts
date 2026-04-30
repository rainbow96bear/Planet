import { toggleFollowResponse } from "../types/user"

export const toggleFollow = async (userid: number): Promise<toggleFollowResponse> => {
    const res = await fetch(`/api/v1/users/${userid}/follow`, {
        method: 'POST',
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}