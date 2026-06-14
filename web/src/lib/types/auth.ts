export interface CreateUserBody{
    username:string
    nickname:string
    password:string
    agreeTerms: boolean
    agreePrivacy: boolean
}

export interface CreateOAuthUserBody {
    username: string
    nickname: string
    agreeTerms: boolean
    agreePrivacy: boolean
}

export interface CheckUsernameResponse {
    username: string
    available: boolean
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