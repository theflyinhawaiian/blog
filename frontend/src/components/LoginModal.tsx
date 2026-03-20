import { useEffect } from 'react'
import { FaGoogle, FaLinkedin, FaGithub, FaTwitter } from 'react-icons/fa'
import { getLoginUrl } from '@api/auth'
import styles from './LoginModal.module.css'

const PROVIDERS = [
  { id: 'google',   label: 'Continue with Google',   color: '#4285f4', Icon: FaGoogle },
  //{ id: 'facebook', label: 'Continue with Facebook', color: '#1877f2', Icon: FaFacebook },
  { id: 'linkedin', label: 'Continue with LinkedIn', color: '#0077b5', Icon: FaLinkedin },
  { id: 'github',   label: 'Continue with GitHub',   color: '#333',    Icon: FaGithub },
  { id: 'twitter',  label: 'Continue with Twitter',  color: '#1da1f2', Icon: FaTwitter },
]

interface Props {
  open: boolean
  onClose: () => void
  returnTo?: string
}

export function LoginModal({ open, onClose, returnTo }: Props) {
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
              href={getLoginUrl(p.id, returnTo)}
              className={styles.providerBtn}
              style={{ background: p.color }}
            >
              <p.Icon aria-hidden /> {p.label}
            </a>
          ))}
        </div>
      </div>
    </div>
  )
}
