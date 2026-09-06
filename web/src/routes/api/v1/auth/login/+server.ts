import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';

export const POST: RequestHandler = async ({ request, fetch, cookies }) => {
	const body = await request.json();

	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/auth/login`,
		{
			method: 'POST',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		},
		cookies,
		fetch
	);

	const data = await res.json();
	if (res.ok) {
		cookies.set('access_token', data.access_token, {
			path: '/',
			httpOnly: true,
			secure: false,
			maxAge: 60 * 60
		});
		cookies.set('refresh_token', data.refresh_token, {
			path: '/',
			httpOnly: true,
			secure: false,
			maxAge: 60 * 60 * 24 * 7
		});
	}

	return json(data, { status: res.status });
};
