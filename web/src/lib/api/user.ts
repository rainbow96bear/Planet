import type { CheckUsernameResponse, CreateUserBody } from "$lib/types/user";

export const createUser = (body:CreateUserBody)=>
    fetch('/api/v1/users', {method:'POST', body:JSON.stringify(body)})

export const checkUsername = async (username: string): Promise<boolean> => {
    const res = await fetch(`/api/v1/users/check?username=${username}`)
    if (!res.ok) throw new Error('서버 오류')
    const json = await res.json()
    return json.data.available
}