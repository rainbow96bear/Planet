import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';

export const PATCH: RequestHandler = async ({ fetch, params, request, cookies }) => {
	const body = await request.json();

	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/tasks/${params.id}`,
		{
			method: 'PATCH',
			headers: { 'Content-Type': 'application/json' },
			body: JSON.stringify(body)
		},
		cookies,
		fetch
	);

	const data = await res.json();
	return json(data, { status: res.status });
};

export const DELETE: RequestHandler = async ({ fetch, params, cookies }) => {
	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/tasks/${params.id}`,
		{ method: 'DELETE' },
		cookies,
		fetch
	);

	return new Response(null, { status: res.status });
};
