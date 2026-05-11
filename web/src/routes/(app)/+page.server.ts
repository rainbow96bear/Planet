import type { PageServerLoad } from './$types'

export const load: PageServerLoad = async ({ parent }) => {
    const { user, feed, exploreFeed } = await parent()
    return { user, feed, exploreFeed }
}