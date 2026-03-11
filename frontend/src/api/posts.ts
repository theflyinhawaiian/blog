import { Post, PostSummary } from '@typedef/post'
import { apiFetch } from './client'

export function fetchPosts(): Promise<PostSummary[]> {
  return apiFetch<PostSummary[]>('/api/posts')
}

export function fetchPost(slug: string): Promise<Post> {
  return apiFetch<Post>(`/api/posts/${slug}`)
}
