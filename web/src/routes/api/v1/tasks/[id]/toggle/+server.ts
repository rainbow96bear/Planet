import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';

export const POST: RequestHandler = async ({ fetch, params, cookies }) => {
	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/tasks/${params.id}/toggle`,
		{ method: 'POST' },
		cookies,
		fetch
	);

	return new Response(null, { status: res.status });
};
