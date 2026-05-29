import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const POST: RequestHandler = async ({ fetch, params, request, cookies }) => {
    const body = await request.json()

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/tasks/${params.id}/reactions`,
        { 
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
        },
        cookies,
        fetch
    )

    return new Response(null, { status: res.status })
}

export const DELETE: RequestHandler = async ({ fetch, params, request, cookies }) => {
    const body = await request.json()

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/tasks/${params.id}/reactions`,
        { 
            method: 'DELETE',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(body)
         },
        cookies,
        fetch
    )

    return new Response(null, { status: res.status })
}