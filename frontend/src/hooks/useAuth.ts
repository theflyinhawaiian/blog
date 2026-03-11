import { useQuery, useQueryClient } from '@tanstack/react-query'
import { fetchMe, logout as apiLogout } from '@api/auth'

export function useAuth() {
  const queryClient = useQueryClient()

  const { data: user, isLoading } = useQuery({
    queryKey: ['me'],
    queryFn: fetchMe,
    retry: false,
    staleTime: 5 * 60 * 1000,
  })

  async function logout() {
    await apiLogout()
    queryClient.setQueryData(['me'], null)
    queryClient.invalidateQueries({ queryKey: ['me'] })
  }

  return { user: user ?? null, isLoading, logout }
}
