import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';
import type { PageServerLoad } from './$types';

export const load: PageServerLoad = async ({ fetch, cookies, url }) => {
	const q = url.searchParams.get('q') ?? '';

	if (!q.trim()) return { users: [], q };

	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/search/users?q=${encodeURIComponent(q)}`,
		{ method: 'GET' },
		cookies,
		fetch
	);
	const data = await res.json();
	return { users: data ?? [], q };
};
