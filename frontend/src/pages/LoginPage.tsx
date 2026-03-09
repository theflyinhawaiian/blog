import { getLoginUrl } from '../api/auth'
import { useAuth } from '../hooks/useAuth'
import { Navigate } from 'react-router-dom'

const PROVIDERS = [
  { id: 'google', label: 'Continue with Google', color: '#4285f4' },
  { id: 'apple', label: 'Continue with Apple', color: '#000' },
  { id: 'facebook', label: 'Continue with Facebook', color: '#1877f2' },
  { id: 'linkedin', label: 'Continue with LinkedIn', color: '#0077b5' },
  { id: 'github', label: 'Continue with GitHub', color: '#333' },
]

export function LoginPage() {
  const { user, isLoading } = useAuth()

  if (!isLoading && user) {
    return <Navigate to="/" replace />
  }

  return (
    <main style={{ maxWidth: '400px', margin: '4rem auto', padding: '2rem', textAlign: 'center' }}>
      <h1 style={{ marginBottom: '2rem' }}>Sign In</h1>
      <div style={{ display: 'flex', flexDirection: 'column', gap: '0.75rem' }}>
        {PROVIDERS.map((p) => (
          <a
            key={p.id}
            href={getLoginUrl(p.id)}
            style={{
              display: 'block',
              padding: '0.75rem 1.5rem',
              background: p.color,
              color: '#fff',
              textDecoration: 'none',
              borderRadius: '6px',
              fontWeight: '500',
              fontSize: '0.95rem',
            }}
          >
            {p.label}
          </a>
        ))}
      </div>
    </main>
  )
}
