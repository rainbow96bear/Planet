import { json, redirect } from '@sveltejs/kit';
import type { RequestHandler } from '@sveltejs/kit';
import {
	GO_API_URL,
	KAKAO_REST_API_KEY,
	KAKAO_REDIRECT_URI,
	KAKAO_CLIENT_SECRET
} from '$env/static/private';
import { dev } from '$app/environment';

// 로컬 개발(http://localhost)에서는 secure 쿠키가 저장되지 않으므로 환경별로 분기한다.
const isSecure = !dev;

export const GET: RequestHandler = async ({ url, fetch, cookies }) => {
	const code = url.searchParams.get('code');
	if (!code) {
		return json({ error: 'code가 없습니다' }, { status: 400 });
	}

	// CSRF 방지: authorize 요청 시 심어둔 state와 콜백으로 돌아온 state가 일치하는지 확인.
	// (authorize 라우트에서 'oauth_state' 쿠키를 생성/저장하도록 별도 수정 필요)
	const returnedState = url.searchParams.get('state');
	const savedState = cookies.get('oauth_state');
	cookies.delete('oauth_state', { path: '/' });

	if (!returnedState || !savedState || returnedState !== savedState) {
		return json({ error: '잘못된 요청입니다 (state mismatch)' }, { status: 400 });
	}

	// 1. code → access_token 교환
	const params = new URLSearchParams();
	params.append('grant_type', 'authorization_code');
	params.append('client_id', KAKAO_REST_API_KEY);
	params.append('redirect_uri', KAKAO_REDIRECT_URI);
	params.append('code', code);
	params.append('client_secret', KAKAO_CLIENT_SECRET);

	const tokenRes = await fetch('https://kauth.kakao.com/oauth/token', {
		method: 'POST',
		headers: {
			'Content-Type': 'application/x-www-form-urlencoded;charset=utf-8'
		},
		body: params.toString()
	});

	if (!tokenRes.ok) {
		return json({ error: '카카오 토큰 교환에 실패했습니다' }, { status: 502 });
	}

	const tokenData = await tokenRes.json();
	const accessToken = tokenData.access_token;

	// access_token → 유저 정보 조회
	const userRes = await fetch('https://kapi.kakao.com/v2/user/me', {
		headers: {
			Authorization: `Bearer ${accessToken}`
		}
	});

	if (!userRes.ok) {
		return json({ error: '카카오 사용자 정보 조회에 실패했습니다' }, { status: 502 });
	}

	const userData = await userRes.json();
	const providerId = String(userData.id); // 카카오 고유 유저 ID

	const goRes = await fetch(`${GO_API_URL}/api/v1/auth/login/oauth`, {
		method: 'POST',
		headers: { 'Content-Type': 'application/json' },
		body: JSON.stringify({
			provider: 'kakao',
			provider_id: providerId
		})
	});

	if (!goRes.ok) {
		return json({ error: '로그인 처리에 실패했습니다' }, { status: 502 });
	}

	const goData = await goRes.json();

	// 신규 사용자 → temp_token 쿠키에 저장 후 회원가입 페이지로
	if (goData.is_new_user) {
		cookies.set('temp_token', goData.temp_token, {
			httpOnly: true,
			secure: isSecure,
			// sameSite: 'lax',
			path: '/',
			maxAge: 60 * 10 // 10분 — Go 백엔드 GenerateTempToken 만료 시간과 일치
		});
		redirect(302, '/signup/oauth');
	}

	// 기존 사용자 → access_token, refresh_token 쿠키에 저장 후 메인으로
	cookies.set('access_token', goData.access_token, {
		httpOnly: true,
		secure: isSecure,
		// sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 // 1시간 — Go 백엔드 GenerateAccessToken 만료 시간과 일치
	});
	cookies.set('refresh_token', goData.refresh_token, {
		httpOnly: true,
		secure: isSecure,
		// sameSite: 'lax',
		path: '/',
		maxAge: 60 * 60 * 24 * 7 // 7일 — Go 백엔드 GenerateRefreshToken 만료 시간과 일치 (기존 30일에서 수정)
	});

	redirect(302, '/');
};