import { redirect } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { 
  KAKAO_REST_API_KEY, 
  KAKAO_REDIRECT_URI, 
  KAKAO_CLIENT_SECRET 
} from '$env/static/private'

export const GET: RequestHandler = async () => {

    const params = new URLSearchParams()
    params.append('response_type', 'code')
    params.append('client_id', KAKAO_REST_API_KEY)
    params.append('redirect_uri', KAKAO_REDIRECT_URI)
    params.append('client_secret', KAKAO_CLIENT_SECRET)

    throw redirect(302, `https://kauth.kakao.com/oauth/authorize?${params}`)
}