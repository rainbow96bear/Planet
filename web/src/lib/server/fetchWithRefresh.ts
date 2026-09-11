import { GO_API_URL } from '$env/static/private';
import type { Cookies } from '@sveltejs/kit';

export async function fetchWithRefresh(
	url: string,
	options: RequestInit,
	cookies: Cookies,
	fetchFn: typeof fetch
): Promise<Response> {
	const token = cookies.get('access_token');

	// access token 주입
	const authOptions = {
		...options,
		headers: {
			...(options.headers as Record<string, string>),
			Authorization: `Bearer ${token}`
		}
	};

	const res = await fetchFn(url, authOptions);

	// 401 아니면 그냥 반환
	if (res.status !== 401) return res;

	// refresh 시도
	const refreshToken = cookies.get('refresh_token');
	if (!refreshToken) return res;

	const refreshRes = await fetchFn(`${GO_API_URL}/api/v1/auth/refresh`, {
		method: 'POST',
		headers: { Authorization: `Bearer ${refreshToken}` }
	});

	if (!refreshRes.ok) {
		cookies.delete('access_token', { path: '/' });
		cookies.delete('refresh_token', { path: '/' });
		return res; // 원래 401 그대로 반환
	}

	const data = await refreshRes.json();

	// 쿠키 갱신
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

	return fetchFn(url, {
		...options,
		headers: {
			...(options.headers as Record<string, string>),
			Authorization: `Bearer ${data.access_token}`
		}
	});
}
