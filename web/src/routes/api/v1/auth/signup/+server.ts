import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import {
    GO_API_URL
} from '$env/static/private'

export const POST: RequestHandler = async ({ request, fetch }) => {
    const body = await request.json()

    const res = await fetch(`${GO_API_URL}/api/v1/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })

    const data = await res.json()
    return json(data, { status: res.status })
}