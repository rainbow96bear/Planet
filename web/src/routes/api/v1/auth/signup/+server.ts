import { json } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import { GO_API_URL } from '$env/static/private';
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh';

export const POST: RequestHandler = async ({ request, fetch, cookies }) => {
	// multipart/form-data로 들어오므로 JSON 파싱 대신 formData로 그대로 전달
	const formData = await request.formData();

	const res = await fetchWithRefresh(
		`${GO_API_URL}/api/v1/auth/signup`,
		{
			method: 'POST',
			body: formData
			// Content-Type 헤더는 지정하지 않음 — fetch가 FormData를 보고
			// boundary 포함된 multipart Content-Type을 자동으로 설정함
		},
		cookies,
		fetch
	);

	const data = await res.json();
	return json(data, { status: res.status });
};
