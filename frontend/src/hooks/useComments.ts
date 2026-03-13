import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'
import { fetchComments, createComment, deleteComment, addReaction } from '@api/comments'

export function useComments(slug: string) {
  return useQuery({
    queryKey: ['comments', slug],
    queryFn: () => fetchComments(slug),
    enabled: !!slug,
  })
}

export function useCreateComment(slug: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (content: string) => createComment(slug, content),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', slug] })
    },
  })
}

export function useDeleteComment(slug: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: (commentId: number) => deleteComment(commentId),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', slug] })
    },
  })
}

export function useAddReaction(slug: string) {
  const queryClient = useQueryClient()
  return useMutation({
    mutationFn: ({ commentId, emoji }: { commentId: number; emoji: string }) =>
      addReaction(commentId, emoji),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ['comments', slug] })
    },
  })
}
