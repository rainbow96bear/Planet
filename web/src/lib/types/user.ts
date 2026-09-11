export interface EnterOrbitResponse {
	id: string;
	is_orbiting: boolean;
}

export interface LeaveOrbitResponse {
	id: string;
	is_orbiting: boolean;
}

export interface UpdateProfileRequest {
	nickname: string;
}

export interface UpdateProfileResponse {
	id: string;
	nickname: string;
}

export interface UserProfile {
	id: string;
	username: string;
	nickname: string;
	profile_image: string | null;
}
