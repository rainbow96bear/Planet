import { json, redirect, RequestHandler } from '@sveltejs/kit'
import { env } from 'process'

const KAKAO_REST_API_KEY = env.KAKAO_REST_API_KEY ?? ""
const KAKAO_REDIRECT_URI = env.KAKAO_REDIRECT_URI ?? ""
const KAKAO_CLIENT_SECRET = env.KAKAO_CLIENT_SECRET ?? ""

export const GET: RequestHandler = async ({ request, fetch, cookies }) => {

    const res = await fetch(`https://kauth.kakao.com/oauth/authorize?response_type=code&client_id=${KAKAO_REST_API_KEY}&redirect_uri=${KAKAO_REDIRECT_URI}`)

    const params = new URLSearchParams()

    params.append('response_type', 'code')
    params.append('client_id', KAKAO_REST_API_KEY)
    params.append('redirect_uri', KAKAO_REDIRECT_URI)
    params.append('client_secret',KAKAO_CLIENT_SECRET) // 선택


    redirect(302, `https://kauth.kakao.com/oauth/authorize?${params}`)
}