export interface FollowResponse {
    id:number,
    is_following:boolean,
}

export interface UnfollowResponse {
    id:number,
    is_following:boolean,
}

export interface UpdateProfileRequest {
    nickname : string
}

export interface UpdateProfileResponse {
    id : number,
    nickname : string,
}