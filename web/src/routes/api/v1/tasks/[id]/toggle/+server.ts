import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import {
    GO_API_URL
} from '$env/static/private'

export const POST: RequestHandler = async ({ fetch, params, cookies }) => {
    const token = cookies.get('access_token')

    const res = await fetch(`${GO_API_URL}/api/v1/tasks/${params.id}/toggle`, {
        method: 'POST',
        headers: {
            'Authorization': `Bearer ${token}`
        }
    })

    return new Response(null, { status: res.status })
}