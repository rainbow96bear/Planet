export interface CreateUserBody{
    username:string
    nickname:string
    password:string
}

export interface CheckUsernameResponse {
    username: string
    available: boolean
}