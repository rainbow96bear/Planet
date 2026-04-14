import { json, RequestHandler } from '@sveltejs/kit'
import { env } from 'process'

const GO_API = env.GO_API_URL

export const GET: RequestHandler = async ({ fetch, url }) => {

    const username = url.searchParams.get('username')

    const res = await fetch(`${GO_API}/api/v1/users/check?username=${username}`)

    const data = await res.json()
    return json(data, { status: res.status })
}