import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'

export const GET: RequestHandler = async ({ params, url, fetch, cookies }) => {
    const { userid } = params
    const year = url.searchParams.get('year')
    const month = url.searchParams.get('month')
    const token = cookies.get('access_token')

    const headers: Record<string, string> = {}
    if (token) {
        headers['Authorization'] = `Bearer ${token}`
    }

    const res = await fetch(
        `${GO_API_URL}/api/v1/users/${userid}/tasks?year=${year}&month=${month}`,
        { headers }
    )

    const data = await res.json()
    return json(data, { status: res.status })
}