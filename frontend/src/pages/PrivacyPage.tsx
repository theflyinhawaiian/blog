import { useState } from 'react'
import { useNavigate } from 'react-router-dom'
import { Helmet } from 'react-helmet-async'
import { useAuth } from '@hooks/useAuth'
import { ConfirmModal } from '@components/ConfirmModal'
import styles from './PrivacyPage.module.css'

export function PrivacyPage() {
  const { user, deleteAccount } = useAuth()
  const navigate = useNavigate()
  const [confirmOpen, setConfirmOpen] = useState(false)

  async function handleDeleteAccount() {
    await deleteAccount()
    setConfirmOpen(false)
    navigate('/')
  }

  return (
    <>
      <Helmet>
        <title>Privacy Policy</title>
      </Helmet>

      <main className={styles.main}>
        <h1>Privacy Policy</h1>

        <p>I take data privacy seriously! This blog collects minimal data necessary to provide its features. When you sign in via a third-party provider (Google, GitHub, Facebook, or LinkedIn), only your display name and a provider-specific identifier to identify your account are stored and this info is encrypted at rest. This site never even sees your passwords.</p>

        <p>I will never sell or share your data with third parties. Data is used solely to operate this blog.</p>

        <p>If you wish to have your account and associated data removed, you may do so below. Your comments will remain but will be anonymized.</p>

        {user && (
          <div className={styles.deleteSection}>
            <h2>Delete Your Data</h2>
            <p>This will permanently delete your account and log you out. Your comments will remain but will be attributed to "Deleted user".</p>
            <button className={styles.deleteBtn} onClick={() => setConfirmOpen(true)}>
              Delete my account
            </button>
          </div>
        )}
      </main>

      <ConfirmModal
        open={confirmOpen}
        title="Delete your account?"
        message="This will permanently delete your account and all associated data. Your comments will remain but will be anonymized. This cannot be undone."
        confirmLabel="Yes, delete my account"
        onConfirm={handleDeleteAccount}
        onClose={() => setConfirmOpen(false)}
      />
    </>
  )
}
