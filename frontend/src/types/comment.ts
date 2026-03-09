export interface Reaction {
  id: number
  comment_id: number
  emoji: string
  count: number
}

export interface Comment {
  id: number
  post_id: number
  user_id: number
  display_name: string
  content: string
  created_at: string
  reactions: Reaction[]
}
