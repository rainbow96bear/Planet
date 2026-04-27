export interface Task {
    id: number
    title: string
    description: string
    date: string
    is_completed: boolean
    is_public: boolean
}

export type GetTasksByMonthResponse = Task[]