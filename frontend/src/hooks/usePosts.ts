import { useQuery } from '@tanstack/react-query'
import { fetchPosts, fetchPost } from '@api/posts'

export function usePosts() {
  return useQuery({
    queryKey: ['posts'],
    queryFn: fetchPosts,
  })
}

export function usePost(slug: string) {
  return useQuery({
    queryKey: ['posts', slug],
    queryFn: () => fetchPost(slug),
    enabled: !!slug,
  })
}
