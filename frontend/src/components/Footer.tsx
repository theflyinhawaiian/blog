import { Link } from 'react-router-dom'
import styles from './Footer.module.css'

export function Footer() {
  return (
    <footer className={styles.footer}>
      <Link to="/privacy" className={styles.link}>Privacy Policy</Link>
    </footer>
  )
}
