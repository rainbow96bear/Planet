import { verifyJwt } from '$lib/utils/jwt';
import { GO_API_URL } from '$env/static/private';
import type { LayoutServerLoad } from './$types';
import type { UserProfile } from '$lib/types/user';

async function fetchProfile(fetch: typeof globalThis.fetch, accessToken: string) {
	try {
		const res = await fetch(`${GO_API_URL}/api/v1/users/me`, {
			headers: { Authorization: `Bearer ${accessToken}` }
		});
		if (!res.ok) return null;
		return (await res.json()) as UserProfile;
	} catch {
		return null; // 프로필 조회 실패가 세션을 죽이면 안 됨
	}
}

export const load: LayoutServerLoad = async ({ cookies, fetch, depends }) => {
	depends('app:user');

	const accessToken = cookies.get('access_token');
	const refreshToken = cookies.get('refresh_token');

	if (!accessToken && !refreshToken) return { user: null };

	try {
		if (accessToken) {
			try {
				const payload = await verifyJwt(accessToken);
				const profile = await fetchProfile(fetch, accessToken);
				return {
					user: {
						userid: payload.userid,
						username: payload.username,
						nickname: profile?.nickname ?? null,
						profileImage: profile?.profile_image ?? null
					}
				};
			} catch {
				// 서명 위조 or 만료 → refresh 흐름으로 진행
			}
		}

		if (!refreshToken) return { user: null };

		const res = await fetch(`${GO_API_URL}/api/v1/auth/refresh`, {
			method: 'POST',
			headers: { Authorization: `Bearer ${refreshToken}` }
		});

		if (!res.ok) {
			cookies.delete('access_token', { path: '/' });
			cookies.delete('refresh_token', { path: '/' });
			return { user: null };
		}

		const data = await res.json();
		cookies.set('access_token', data.access_token, {
			httpOnly: true,
			path: '/',
			maxAge: 60 * 60
		});
		cookies.set('refresh_token', data.refresh_token, {
			httpOnly: true,
			path: '/',
			maxAge: 60 * 60 * 24 * 30
		});

		const newPayload = await verifyJwt(data.access_token);
		const profile = await fetchProfile(fetch, data.access_token);
		return {
			user: {
				userid: newPayload.userid,
				username: newPayload.username,
				nickname: profile?.nickname ?? null,
				profile_image: profile?.profile_image ?? null
			}
		};
	} catch {
		return { user: null };
	}
};
