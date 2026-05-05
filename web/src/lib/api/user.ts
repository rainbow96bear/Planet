import type { FollowResponse, UnfollowResponse, UpdateProfileRequest } from "$lib/types/user"

export const follow = async (userid: number): Promise<FollowResponse> => {
    const res = await fetch(`/api/v1/users/${userid}/follow`, {
        method: 'POST',
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}


export const unfollow = async (userid: number): Promise<UnfollowResponse> => {
    const res = await fetch(`/api/v1/users/${userid}/follow`, {
        method: 'DELETE',
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const updateProfile = async (userid: number, body: UpdateProfileRequest): Promise<UpdateProfileResponse> => {
    const res = await fetch(`/api/v1/users/${userid}`, {
        method: 'PATCH',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })
    
    if (!res.ok) {
        const data = await res.json()
        throw new Error(data.error ?? '수정에 실패했습니다.')
    }

    return res.json()
}