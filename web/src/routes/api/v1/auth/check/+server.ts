import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const GET: RequestHandler = async ({ fetch, url, cookies }) => {
    const username = url.searchParams.get('username')

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/auth/check?username=${username}`,
        { method: 'GET' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}