import type { LayoutServerLoad } from './$types'
import { GO_API_URL } from '$env/static/private'

export const load: LayoutServerLoad = async ({ parent, cookies, fetch }) => {
    const { user } = await parent()

    if (!user) return { user, feed: [], exploreFeed: [] }

    const token = cookies.get('access_token')
    const headers: Record<string, string> = token
        ? { Authorization: `Bearer ${token}` }
        : {}

    const [feedRes, exploreRes] = await Promise.all([
        fetch(`${GO_API_URL}/api/v1/feed`, { headers }),
        fetch(`${GO_API_URL}/api/v1/feed/explore`, { headers })
    ])

    const [feed, exploreFeed] = await Promise.all([
        feedRes.ok ? feedRes.json() : [],
        exploreRes.ok ? exploreRes.json() : []
    ])
    return { user, feed: feed ?? [], exploreFeed: exploreFeed ?? [] }
}