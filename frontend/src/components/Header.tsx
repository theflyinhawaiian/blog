import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

export function Header() {
  const { user, isLoading, logout } = useAuth()

  return (
    <header style={{ borderBottom: '1px solid #eee', padding: '1rem 2rem', display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
      <Link to="/" style={{ textDecoration: 'none', fontWeight: 'bold', fontSize: '1.25rem', color: '#111' }}>
        Blog
      </Link>
      <nav>
        {isLoading ? null : user ? (
          <span style={{ display: 'flex', alignItems: 'center', gap: '1rem' }}>
            <span style={{ color: '#555' }}>{user.display_name}</span>
            <button onClick={logout} style={{ cursor: 'pointer', border: '1px solid #ccc', background: 'none', padding: '0.25rem 0.75rem', borderRadius: '4px' }}>
              Logout
            </button>
          </span>
        ) : (
          <Link to="/login" style={{ textDecoration: 'none', border: '1px solid #333', padding: '0.25rem 0.75rem', borderRadius: '4px', color: '#333' }}>
            Login
          </Link>
        )}
      </nav>
    </header>
  )
}
