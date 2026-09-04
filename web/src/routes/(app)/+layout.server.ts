import type { LayoutServerLoad } from './$types'
import { GO_API_URL } from '$env/static/private'

async function safeFetchJson<T>(fetch: typeof globalThis.fetch, url: string, headers: Record<string, string>): Promise<T[]> {
    try {
        const res = await fetch(url, { headers, signal: AbortSignal.timeout(3000) })
        if (!res.ok) return []
        const data = await res.json()
        return data ?? []
    } catch {
        return []
    }
}

export const load: LayoutServerLoad = async ({ parent, cookies, fetch }) => {
    const { user } = await parent()

    if (!user) return { user, feed: [], exploreFeed: [] }

    const token = cookies.get('access_token')
    const headers: Record<string, string> = token
        ? { Authorization: `Bearer ${token}` }
        : {}

    const [feed, exploreFeed] = await Promise.all([
        safeFetchJson(fetch, `${GO_API_URL}/api/v1/feed`, headers),
        safeFetchJson(fetch, `${GO_API_URL}/api/v1/feed/explore`, headers)
    ])

    return { user, feed, exploreFeed }
}