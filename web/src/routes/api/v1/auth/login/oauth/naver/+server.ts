import { redirect } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { 
  NAVER_CLIENT_ID, 
  NAVER_REDIRECT_URI, 
} from '$env/static/private'

export const GET: RequestHandler = async ({ cookies }) => {
    const state = crypto.randomUUID()
    
    // 콜백에서 검증할 수 있도록 쿠키에 저장
    cookies.set('oauth_state', state, {
        httpOnly: true,
        path: '/',
        maxAge: 60 * 10, // 10분
    })

    const params = new URLSearchParams()
    params.append('response_type', 'code')
    params.append('client_id', NAVER_CLIENT_ID)
    params.append('redirect_uri', NAVER_REDIRECT_URI)
    params.append('state', state)

    throw redirect(302, `https://nid.naver.com/oauth2.0/authorize?${params}`)
}