export interface Notification {
  id: number
  type: 'follow' | 'comment' | 'reaction'
  message: string
  is_read: boolean
  created_at: string
  from_user: { userid: string; nickname: string; username: string }
}