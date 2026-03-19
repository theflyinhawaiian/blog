export interface PostSummary {
  id: number
  title: string
  slug: string
  excerpt?: { String: string; Valid: boolean }
  post_image?: { String: string; Valid: boolean }
  tags?: string[]
  created_at: string
}

export interface Post extends PostSummary {
  content: string
  meta_description?: { String: string; Valid: boolean }
  canonical_url?: { String: string; Valid: boolean }
  updated_at: string
}
