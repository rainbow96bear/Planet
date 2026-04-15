export interface CreateUserBody{
    username:string
    nickname:string
    password:string
}

export interface CreateOAuthUserBody{
    username:string
    nickname:string
}


export interface CheckUsernameResponse {
    data : {
        username: string
        available: boolean
    }
}

export interface LoginBody {
    username: string
    password: string
}

export interface LoginResponse {
    username: string
    access_token: string
    refresh_token: string
}