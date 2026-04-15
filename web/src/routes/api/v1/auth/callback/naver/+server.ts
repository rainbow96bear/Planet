import { json, redirect } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { 
  GO_API_URL,
  NAVER_CLIENT_ID,
  NAVER_CLIENT_SECRET
} from '$env/static/private'


export const GET: RequestHandler = async ({ url, fetch, cookies }) => {
  const code = url.searchParams.get('code')
  const state = url.searchParams.get('state')

  if (!code) {
    return json({ error: 'code가 없습니다' }, { status: 400 })
  }

  if (!state) {
    return json({ error: 'state가 일치하지 않습니다' }, { status: 400 })
  }

  // code → access_token 교환
  const params = new URLSearchParams()
  params.append('grant_type', 'authorization_code')
  params.append('client_id', NAVER_CLIENT_ID)
  params.append('client_secret', NAVER_CLIENT_SECRET)
  params.append('code', code)
  params.append('state', state)
  const tokenRes = await fetch('https://nid.naver.com/oauth2.0/token', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8',
    },
    body: params.toString(),
  })

  const tokenData = await tokenRes.json()
  const accessToken = tokenData.access_token

  // access_token → 유저 정보 조회
  const userRes = await fetch('https://openapi.naver.com/v1/nid/me', {
    headers: {
      'Authorization': `Bearer ${accessToken}`,
    },
  })

  const userData = await userRes.json()
  const providerId = String(userData.response.id)  // 카카오 고유 유저 ID

  const goRes = await fetch(`${GO_API_URL}/api/v1/auth/login/oauth`, {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
      provider: 'naver',
      provider_id: providerId,
    }),
  })
  const goData = await goRes.json()

  // 신규 사용자 → temp_token 쿠키에 저장 후 회원가입 페이지로
  if (goData.is_new_user) {
    cookies.set('temp_token', goData.temp_token, {
      httpOnly: true,
      path: '/',
      maxAge: 60 * 10, // 10분
    })
    redirect(302, '/signup/oauth')
  }

  // 기존 사용자 → access_token, refresh_token 쿠키에 저장 후 메인으로
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