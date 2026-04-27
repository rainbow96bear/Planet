export interface Task {
    id: number
    title: string
    description: string
    date: string
    is_completed: boolean
    is_public: boolean
}

export type GetTasksByMonthResponse = Task[]

export interface CreateTaskBody {
    title: string
    description?: string
    date: string
    is_public: boolean
}

export interface CreateTaskResponse {
    id: number
    title: string
    date: string
    is_completed: boolean
    is_public: boolean
    created_at: string
}