import { json, RequestHandler } from '@sveltejs/kit'
import { env } from 'process'

const GO_API = env.GO_API_URL

export const POST: RequestHandler = async ({ request, fetch }) => {
    const body = await request.json()

    const res = await fetch(`${GO_API}/api/v1/users`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(body)
    })

    const data = await res.json()
    return json(data, { status: res.status })
}