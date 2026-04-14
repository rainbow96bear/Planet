import { json, redirect } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { env } from 'process'

const KAKAO_REST_API_KEY = env.KAKAO_REST_API_KEY ?? ""
const KAKAO_REDIRECT_URI = env.KAKAO_REDIRECT_URI ?? ""
const KAKAO_CLIENT_SECRET = env.KAKAO_CLIENT_SECRET ?? ""
const GO_API_URL = env.GO_API_URL ?? ""

export const GET: RequestHandler = async ({ url, fetch, cookies }) => {
  const code = url.searchParams.get('code')

  if (!code) {
    return json({ error: 'code가 없습니다' }, { status: 400 })
  }

  // 1. code → access_token 교환
  const params = new URLSearchParams()
  params.append('grant_type', 'authorization_code')
  params.append('client_id', KAKAO_REST_API_KEY)
  params.append('redirect_uri', KAKAO_REDIRECT_URI)
  params.append('code', code)
  params.append('client_secret', KAKAO_CLIENT_SECRET)

  const tokenRes = await fetch('https://kauth.kakao.com/oauth/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8',
    },
    body: params.toString(),
  })

  const tokenData = await tokenRes.json()
  const accessToken = tokenData.access_token

  // 2. access_token → 유저 정보 조회
  const userRes = await fetch('https://kapi.kakao.com/v2/user/me', {
    headers: {
      'Authorization': `Bearer ${accessToken}`,
    },
  })

  const userData = await userRes.json()

  const providerId = String(userData.id)  // 카카오 고유 유저 ID
  const email = userData.kakao_account?.email

  const goRes = await fetch(`${GO_API_URL}/api/v1/auth/oauth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      provider: 'kakao',
      provider_id: providerId,
      email,
    }),
  })
  const goData = await goRes.json()

  // 7. 신규 사용자 → temp_token 쿠키에 저장 후 회원가입 페이지로
  if (goData.is_new_user) {
    cookies.set('temp_token', goData.temp_token, {
      httpOnly: true,
      path: '/',
      maxAge: 60 * 10, // 10분
    })
    redirect(302, '/signup/oauth')
  }

  // 8. 기존 사용자 → access_token, refresh_token 쿠키에 저장 후 메인으로
  cookies.set('access_token', goData.access_token, {
    httpOnly: true,
    path: '/',
    maxAge: 60 * 60, // 1시간
  })
  cookies.set('refresh_token', goData.refresh_token, {
    httpOnly: true,
    path: '/',
    maxAge: 60 * 60 * 24 * 30, // 30일
  })

  redirect(302, '/')
}