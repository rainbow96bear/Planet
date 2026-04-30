import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'

export const POST: RequestHandler = async ({ params, url, fetch, cookies }) => {
    const { userid } = params

    const token = cookies.get('access_token')

    const headers: Record<string, string> = {}
    if (token) {
        headers['Authorization'] = `Bearer ${token}`
    }

    const res = await fetch(
        `${GO_API_URL}/api/v1/users/${userid}/follow`,
        { 
            method: 'POST',
            headers
        }
    )

    const data = await res.json()
    return json(data, { status: res.status })
}

export const DELETE: RequestHandler = async ({ params, url, fetch, cookies }) => {
    const { userid } = params

    const token = cookies.get('access_token')

    const headers: Record<string, string> = {}
    if (token) {
        headers['Authorization'] = `Bearer ${token}`
    }

    const res = await fetch(
        `${GO_API_URL}/api/v1/users/${userid}/follow`,
        { 
            method: 'DELETE',
            headers 
        }
    )

    const data = await res.json()
    return json(data, { status: res.status })
}