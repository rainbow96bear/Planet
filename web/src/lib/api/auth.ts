import type { CreateUserBody, CreateOAuthUserBody, CheckUsernameResponse, LoginBody, LoginResponse } from "$lib/types/auth";

export const createUser = async (
    body: CreateUserBody,
    profileImage?: File | null
) => {
    const formData = new FormData()
    Object.entries(body).forEach(([key, value]) => {
        formData.append(key, String(value))
    })
    if (profileImage) {
        formData.append('profile_image', profileImage)
    }

    const res = await fetch('/api/v1/auth/signup', {
        method: 'POST',
        body: formData
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const createOAuthUser = async (
    body: CreateOAuthUserBody,
    profileImage?: File | null
) => {
    const formData = new FormData()
    Object.entries(body).forEach(([key, value]) => {
        formData.append(key, String(value))
    })
    if (profileImage) {
        formData.append('profile_image', profileImage)
    }

    const res = await fetch('/api/v1/auth/signup/oauth', {
        method: 'POST',
        body: formData
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const checkUsername = async (username: string): Promise<CheckUsernameResponse> => {
    const res = await fetch(`/api/v1/auth/check?username=${username}`)
    if (!res.ok) throw new Error('서버 오류')
    return res.json()
}

export const login = async (body: LoginBody): Promise<LoginResponse> => {
    const res = await fetch('/api/v1/auth/login', {
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

export const kakaoOAuthLogin = async() => {
    const res = await fetch('/api/v1/auth/login/oauth/kakao')
}

export const naverOAuthLogin = async () => {
    const res = await fetch('/api/v1/auth/login/oauth/naver')
}

export const logout = async () => {
    const res = await fetch('/api/v1/auth/logout', {
        method: 'POST',
    })
}