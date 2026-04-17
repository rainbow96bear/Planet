// +layout.server.ts
import { decodeJwtPayload } from '$lib/utils/jwt'
import { 
  GO_API_URL,
} from '$env/static/private'

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
    const accessToken = cookies.get('access_token')
    if (!accessToken) return { user: null }

    try {
        const payload = decodeJwtPayload(accessToken)
        
        if (payload.exp * 1000 > Date.now()) {
            return { user: { id: payload.sub, username: payload.username, nickname: payload.nickname } }
        }

        const refreshToken = cookies.get('refresh_token')
        if (!refreshToken) return { user: null }

        const res = await fetch(`${GO_API_URL}/api/v1/auth/refresh`, {
            method: 'POST',
            headers: { 'Authorization': `Bearer ${refreshToken}` }
        })

        if (!res.ok) {
            cookies.delete('access_token', { path: '/' })
            cookies.delete('refresh_token', { path: '/' })
            return { user: null }
        }

        const data = await res.json()
        cookies.set('access_token', data.access_token, {
             httpOnly: true, path: '/', maxAge: 60 * 60 
        })
        cookies.set('refresh_token', data.refresh_token, { 
            httpOnly: true, path: '/', maxAge: 60 * 60 * 24 * 30 
        })
        
        const newPayload = decodeJwtPayload(data.access_token)
        return { user: { id: newPayload.sub, username: newPayload.username, nickname: newPayload.nickname } }

    } catch {
        return { user: null }
    }
}