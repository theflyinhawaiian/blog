import { User } from '@typedef/user'
import { apiFetch } from './client'

const API_URL = import.meta.env.VITE_API_URL || 'http://localhost:8080'

export function fetchMe(): Promise<User> {
  return apiFetch<User>('/auth/me')
}

export function logout(): Promise<void> {
  return apiFetch<void>('/auth/logout', { method: 'POST' })
}

export function getLoginUrl(provider: string): string {
  return `${API_URL}/auth/${provider}`
}
