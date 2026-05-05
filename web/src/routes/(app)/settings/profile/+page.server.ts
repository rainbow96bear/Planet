import { redirect } from '@sveltejs/kit'
import type { PageServerLoad } from '../../profile/[userid]/$types'

export const load: PageServerLoad = async ({ parent }) => {
    const { user } = await parent()
    if (!user) redirect(302, '/login')
}