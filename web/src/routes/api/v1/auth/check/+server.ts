import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import {
    GO_API_URL
} from '$env/static/private'


export const GET: RequestHandler = async ({ fetch, url }) => {

    const username = url.searchParams.get('username')

    const res = await fetch(`${GO_API_URL}/api/v1/auth/check?username=${username}`)

    const data = await res.json()
    return json(data, { status: res.status })
}