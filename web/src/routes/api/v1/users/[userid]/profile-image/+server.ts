import { json } from '@sveltejs/kit'
import type { RequestHandler } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'
import { fetchWithRefresh } from '$lib/server/fetchWithRefresh'

export const POST: RequestHandler = async ({ params, request, cookies, fetch }) => {
    const { userid } = params

    // 클라이언트가 보낸 multipart/form-data를 그대로 파싱
    const formData = await request.formData()
    const profileImage = formData.get('profile_image')

    if (!(profileImage instanceof File)) {
        return json({ error: 'profile_image 파일이 필요합니다.' }, { status: 400 })
    }

    // Go 백엔드로 전달할 새 FormData 구성 (스트림 재사용 불가하므로 새로 생성)
    const forwardForm = new FormData()
    forwardForm.append('profile_image', profileImage, profileImage.name)

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/users/${userid}/profile-image`,
        { method: 'POST', body: forwardForm },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}

export const DELETE: RequestHandler = async ({ params, cookies, fetch }) => {
    const { userid } = params

    const res = await fetchWithRefresh(
        `${GO_API_URL}/api/v1/users/${userid}/profile-image`,
        { method: 'DELETE' },
        cookies,
        fetch
    )

    const data = await res.json()
    return json(data, { status: res.status })
}