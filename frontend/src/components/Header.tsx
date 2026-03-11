import { Link } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import styles from './Header.module.css'

export function Header() {
  const { user, isLoading, logout } = useAuth()

  return (
    <header className={styles.header}>
      <Link to="/" className={styles.logo}>
        Blog
      </Link>
      <nav>
        {isLoading ? null : user ? (
          <span className={styles.userNav}>
            <span className={styles.userName}>{user.display_name}</span>
            <button onClick={logout} className={styles.logoutBtn}>
              Logout
            </button>
          </span>
        ) : (
          <Link to="/login" className={styles.loginLink}>
            Login
          </Link>
        )}
      </nav>
    </header>
  )
}
