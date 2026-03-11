import { getLoginUrl } from '../api/auth'
import { useAuth } from '../hooks/useAuth'
import { Navigate } from 'react-router-dom'
import styles from './LoginPage.module.css'

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
    <main className={styles.main}>
      <h1 className={styles.heading}>Sign In</h1>
      <div className={styles.providers}>
        {PROVIDERS.map((p) => (
          <a
            key={p.id}
            href={getLoginUrl(p.id)}
            className={styles.providerBtn}
            style={{ background: p.color }}
          >
            {p.label}
          </a>
        ))}
      </div>
    </main>
  )
}
