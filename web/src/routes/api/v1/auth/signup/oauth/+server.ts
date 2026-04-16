import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import {
    GO_API_URL
} from '$env/static/private'

export const POST: RequestHandler = async ({ request, fetch, cookies }) => {
    const body = await request.json()
    const tempToken = cookies.get('temp_token')  // 여기서 꺼내서

    if (!tempToken) {
        return json({ error: '인증 정보가 없습니다' }, { status: 401 })
    }

    const res = await fetch(`${GO_API_URL}/api/v1/auth/signup/oauth`, {
        method: 'POST',
        headers: { 
            'Content-Type': 'application/json',
            'Authorization': `Bearer ${tempToken}`  // 여기서 Go로 전달
        },
        body: JSON.stringify(body)
    })
    const data = await res.json()
    
    cookies.delete('temp_token', { path: '/' })
    
    return json(data, { status: res.status })
}