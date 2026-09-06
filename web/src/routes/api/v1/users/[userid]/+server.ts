import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';
import { json, type RequestHandler } from '@sveltejs/kit';

export const PATCH: RequestHandler = async ({ params, request, fetch, cookies }) => {
	const { userid } = params;
	const body = await request.json();

	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/users/${userid}`,
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
