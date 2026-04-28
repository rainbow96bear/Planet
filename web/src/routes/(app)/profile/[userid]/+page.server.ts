import type { PageServerLoad } from './$types'
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

    const [userRes, tasksRes] = await Promise.all([
        fetch(`${GO_API_URL}/api/v1/users/${userid}`, { headers }),
        fetch(`${GO_API_URL}/api/v1/users/${userid}/tasks?year=${year}&month=${month}`, { headers })
    ])

    const user = await userRes.json()
    const tasks = await tasksRes.json()

    return { user, tasks, year, month, me }
}