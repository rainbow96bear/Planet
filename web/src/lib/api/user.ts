import type { CreateUserBody, CreateUserResponse, CheckUsernameResponse, LoginBody, LoginResponse } from "$lib/types/user";

export const createUser = async (body: CreateUserBody): Promise<CreateUserResponse> => {
    const res = await fetch('/api/v1/users', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },  // 누락됨
        body: JSON.stringify(body)
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}

export const checkUsername = async (username: string): Promise<CheckUsernameResponse> => {
    const res = await fetch(`/api/v1/users/check?username=${username}`)
    if (!res.ok) throw new Error('서버 오류')
    return res.json()
}

export const login = async (body: LoginBody): Promise<LoginResponse> => {
    const res = await fetch('/api/v1/auth/login', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify(body)
    })

    if (!res.ok) {
        const err = await res.json()
        throw new Error(err.error ?? '서버 오류')
    }

    return res.json()
}