import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const GET: RequestHandler = async ({ fetch, cookies }) => {
    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/notifications/unread-count`,
        { method: 'GET' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}