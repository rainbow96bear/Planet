import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const DELETE: RequestHandler = async ({ fetch, params, cookies }) => {
    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/tasks/${params.id}`,
        { method: 'DELETE' },
        cookies,
        fetch
    )

    return new Response(null, { status: res.status })
}