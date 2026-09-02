import type { PageServerLoad } from './$types'
import { error } from '@sveltejs/kit'
import { GO_API_URL } from '$env/static/private'

export const load: PageServerLoad = async ({ params, fetch, cookies, parent }) => {
    const { userid } = params
    const { user: me } = await parent()
    const token = cookies.get('access_token')

    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1

    const headers: Record<string, string> = token
        ? { Authorization: `Bearer ${token}` }
        : {}

    const [profileUserRes, tasksRes] = await Promise.all([
        fetch(`${GO_API_URL}/api/v1/users/${userid}`, { headers }),
        fetch(`${GO_API_URL}/api/v1/users/${userid}/tasks?year=${year}&month=${month}`, { headers })
    ])

    // 프로필 조회는 이 페이지의 핵심 데이터라, 실패하면 여기서 명확히 에러 처리
    if (!profileUserRes.ok) {
        const body = await profileUserRes.text()
        console.error(`GetProfile failed (${profileUserRes.status}):`, body)
        throw error(
            profileUserRes.status === 404 ? 404 : 500,
            profileUserRes.status === 404 ? '존재하지 않는 사용자입니다.' : '프로필 정보를 불러오지 못했습니다.'
        )
    }

    // tasks는 부가 데이터라, 실패해도 페이지 전체를 죽이지 않고 빈 배열로 대체
    let tasks = []
    if (tasksRes.ok) {
        tasks = await tasksRes.json()
    } else {
        const body = await tasksRes.text()
        console.error(`GetTasksByMonth failed (${tasksRes.status}):`, body)
    }

    const profileUser = await profileUserRes.json()
    return { profileUser, tasks, year, month, me }
}