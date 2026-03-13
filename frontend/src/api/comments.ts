import { Comment, Reaction } from '../types/comment'
import { apiFetch } from './client'

export async function fetchComments(slug: string): Promise<Comment[]> {
  const result = await apiFetch<Comment[]>(`/api/posts/${slug}/comments`)
  return result ?? []
}

export function createComment(slug: string, content: string): Promise<Comment> {
  return apiFetch<Comment>(`/api/posts/${slug}/comments`, {
    method: 'POST',
    body: JSON.stringify({ content }),
  })
}

export function updateComment(commentId: number, content: string): Promise<Comment> {
  return apiFetch<Comment>(`/api/comments/${commentId}`, {
    method: 'PATCH',
    body: JSON.stringify({ content }),
  })
}

export function deleteComment(commentId: number): Promise<void> {
  return apiFetch<void>(`/api/comments/${commentId}`, { method: 'DELETE' })
}

export function addReaction(commentId: number, emoji: string): Promise<Reaction> {
  return apiFetch<Reaction>(`/api/comments/${commentId}/reactions`, {
    method: 'POST',
    body: JSON.stringify({ emoji }),
  })
}
