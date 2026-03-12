export interface Reaction {
  emoji: string
  count: number
  reacted_by_me: boolean
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
