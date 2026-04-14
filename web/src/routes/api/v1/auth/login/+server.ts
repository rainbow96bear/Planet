import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import {
    GO_API_URL
} from '$env/static/private'

export const POST: RequestHandler = async ({ request, fetch, cookies }) => {
    const body = await request.json()

    const res = await fetch(`${GO_API_URL}/api/v1/auth/login`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })

    const data = await res.json()

    if (res.ok) {
        cookies.set('access_token', data.data.access_token, {
            path: '/',
            httpOnly: true,
            secure: false,           // production에서 true
            maxAge: 60 * 60          // 1시간
        })
        cookies.set('refresh_token', data.data.refresh_token, {
            path: '/',
            httpOnly: true,
            secure: false,
            maxAge: 60 * 60 * 24 * 7 // 7일
        })
    }

    return json(data, { status: res.status })
}