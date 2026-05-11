export interface FollowResponse {
    id:string,
    is_following:boolean,
}

export interface UnfollowResponse {
    id:string,
    is_following:boolean,
}

export interface UpdateProfileRequest {
    nickname : string
}

export interface UpdateProfileResponse {
    id : string,
    nickname : string,
}