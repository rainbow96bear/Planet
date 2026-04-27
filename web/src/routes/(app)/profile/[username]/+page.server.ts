import type { PageServerLoad } from './$types'
import { GO_API_URL } from '$env/static/private'

export const load: PageServerLoad = async ({ params, fetch, cookies }) => {
    const { username } = params
    const token = cookies.get('access_token')

    const now = new Date()
    const year = now.getFullYear()
    const month = now.getMonth() + 1

    const headers: Record<string, string> = {}
    if (token) {
        headers['Authorization'] = `Bearer ${token}`
    }

    const [userRes, tasksRes] = await Promise.all([
        fetch(`${GO_API_URL}/api/v1/users/${username}`, { headers }),
        fetch(`${GO_API_URL}/api/v1/users/${username}/tasks?year=${year}&month=${month}`, { headers })
    ])

    const user = await userRes.json()
    const tasks = await tasksRes.json()

    return { user, tasks, year, month }
}