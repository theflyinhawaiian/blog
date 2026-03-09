import { Comment, Reaction } from '../types/comment'
import { apiFetch } from './client'

export function fetchComments(slug: string): Promise<Comment[]> {
  return apiFetch<Comment[]>(`/api/posts/${slug}/comments`)
}

export function createComment(slug: string, content: string): Promise<Comment> {
  return apiFetch<Comment>(`/api/posts/${slug}/comments`, {
    method: 'POST',
    body: JSON.stringify({ content }),
  })
}

export function addReaction(commentId: number, emoji: string): Promise<Reaction> {
  return apiFetch<Reaction>(`/api/comments/${commentId}/reactions`, {
    method: 'POST',
    body: JSON.stringify({ emoji }),
  })
}
