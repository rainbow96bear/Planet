import { redirect } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { 
  NAVER_CLIENT_ID, 
  NAVER_REDIRECT_URI, 
} from '$env/static/private'

export const GET: RequestHandler = async () => {
    const params = new URLSearchParams()
    params.append('response_type', 'code')
    params.append('client_id', NAVER_CLIENT_ID)
    params.append('redirect_uri', NAVER_REDIRECT_URI)

    throw redirect(302, `https://nid.naver.com/oauth2.0/authorize?${params}`)
}