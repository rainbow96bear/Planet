import type {
	EnterOrbitResponse,
	LeaveOrbitResponse,
	UpdateProfileRequest,
	UpdateProfileResponse
} from '$lib/types/user';

export const enterOrbit = async (userid: string): Promise<EnterOrbitResponse> => {
	const res = await fetch(`/api/v1/users/${userid}/orbit`, {
		method: 'POST'
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}

	return res.json();
};

export const leaveOrbit = async (userid: string): Promise<LeaveOrbitResponse> => {
	const res = await fetch(`/api/v1/users/${userid}/orbit`, {
		method: 'DELETE'
	});

	if (!res.ok) {
		const err = await res.json();
		throw new Error(err.error ?? '서버 오류');
	}

	return res.json();
};

export const updateProfile = async (
	userid: string,
	body: UpdateProfileRequest
): Promise<UpdateProfileResponse> => {
	const res = await fetch(`/api/v1/users/${userid}`, {
		method: 'PATCH',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify(body)
	});

	if (!res.ok) {
		const data = await res.json();
		throw new Error(data.error ?? '수정에 실패했습니다.');
	}

	return res.json();
};

export async function uploadProfileImage(
	userid: string,
	file: File
): Promise<{ profile_image: string }> {
	const formData = new FormData();
	formData.append('profile_image', file);

	const res = await fetch(`/api/v1/users/${userid}/profile-image`, {
		method: 'POST',
		body: formData
	});
	if (!res.ok) throw new Error('profile image upload failed');
	return res.json();
}

export async function deleteProfileImage(userid: string): Promise<void> {
	const res = await fetch(`/api/v1/users/${userid}/profile-image`, {
		method: 'DELETE'
	});
	if (!res.ok) throw new Error('profile image delete failed');
}
