import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';

export const PATCH: RequestHandler = async ({ fetch, cookies }) => {
	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/notifications/read-all`,
		{ method: 'PATCH' },
		cookies,
		fetch
	);

	const data = await res.json();
	return json(data, { status: res.status });
};
