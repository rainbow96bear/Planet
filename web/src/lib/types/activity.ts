export interface Activity {
    activity_id: string
    actor_id: string
    actor_nickname: string
    type: 'task.created' | 'task.completed'
    task_title?: string
    created_at: string
    reactions?: {
        like: number
        cheer: number
    }
}