export interface Task {
    id: string
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
    id: string
    title: string
    date: string
    is_completed: boolean
    is_public: boolean
    created_at: string
}

export interface ToggleTaskResponse {
    id: string,
    is_completed: boolean,
}