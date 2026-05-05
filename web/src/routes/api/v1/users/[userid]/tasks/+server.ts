import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const GET: RequestHandler = async ({ params, url, fetch, cookies }) => {
    const { userid } = params
    const year = url.searchParams.get('year')
    const month = url.searchParams.get('month')

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/users/${userid}/tasks?year=${year}&month=${month}`,
        { method: 'GET' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}