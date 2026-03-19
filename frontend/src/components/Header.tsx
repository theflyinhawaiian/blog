import { Link } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import { useLoginModal } from '@hooks/useLoginModal'
import styles from './Header.module.css'

function getGreeting() {
  const hour = new Date().getHours()
  if (hour > 3 && hour < 12) return 'Good morning'
  if (hour < 18) return 'Good afternoon'
  return 'Good evening'
}

export function Header() {
  const { user, isLoading, logout } = useAuth()
  const { openLoginModal } = useLoginModal()

  return (
    <header className={styles.header}>
      <div />
      <Link to="/" className={styles.logo}>Pete's Blog</Link>
      <nav>
        {isLoading ? null : user ? (
          <span className={styles.userNav}>
            <span className={styles.userName}>{getGreeting()}, {user.display_name}!</span>
            <button onClick={logout} className={styles.logoutBtn}>
              Logout
            </button>
          </span>
        ) : (
          <button onClick={openLoginModal} className={styles.loginLink}>
            Login
          </button>
        )}
      </nav>
    </header>
  )
}
