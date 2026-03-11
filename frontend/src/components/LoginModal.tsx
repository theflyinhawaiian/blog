import { useEffect } from 'react'
import { getLoginUrl } from '@api/auth'
import styles from './LoginModal.module.css'

const PROVIDERS = [
  { id: 'google', label: 'Continue with Google', color: '#4285f4' },
  { id: 'apple', label: 'Continue with Apple', color: '#000' },
  { id: 'facebook', label: 'Continue with Facebook', color: '#1877f2' },
  { id: 'linkedin', label: 'Continue with LinkedIn', color: '#0077b5' },
  { id: 'github', label: 'Continue with GitHub', color: '#333' },
]

interface Props {
  open: boolean
  onClose: () => void
}

export function LoginModal({ open, onClose }: Props) {
  useEffect(() => {
    function onKeyDown(e: KeyboardEvent) {
      if (e.key === 'Escape') onClose()
    }
    if (open) document.addEventListener('keydown', onKeyDown)
    return () => document.removeEventListener('keydown', onKeyDown)
  }, [open, onClose])

  if (!open) return null

  return (
    <div className={styles.overlay} onClick={onClose}>
      <div className={styles.modal} onClick={(e) => e.stopPropagation()}>
        <button className={styles.closeBtn} onClick={onClose} aria-label="Close">✕</button>
        <h2 className={styles.heading}>Sign In</h2>
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
      </div>
    </div>
  )
}
