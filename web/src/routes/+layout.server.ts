// +layout.server.ts
import { 
  GO_API_URL,
} from '$env/static/private'

export const load: LayoutServerLoad = async ({ cookies, fetch }) => {
    const accessToken = cookies.get('access_token')
    if (!accessToken) return { user: null }

    try {
        const payload = JSON.parse(atob(accessToken.split('.')[1]))
        
        // 만료 안됨
        if (payload.exp * 1000 > Date.now()) {
            return { user: { id: payload.sub, username: payload.username, nickname: payload.nickname } }
        }

        // 만료됨 → refresh 시도
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
        
        const newPayload = JSON.parse(atob(data.access_token.split('.')[1]))
        return { user: { id: newPayload.sub, username: newPayload.username, nickname: newPayload.nickname } }

    } catch {
        return { user: null }
    }
}