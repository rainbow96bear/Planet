import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const POST: RequestHandler = async ({ params, fetch, cookies }) => {
    const { userid } = params

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/users/${userid}/orbit`,
        { method: 'POST' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}

export const DELETE: RequestHandler = async ({ params, fetch, cookies }) => {
    const { userid } = params

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/users/${userid}/orbit`,
        { method: 'DELETE' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}