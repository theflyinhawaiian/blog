import { useState } from 'react'
import { useAuth } from '@hooks/useAuth'
import { useLoginModal } from '@hooks/useLoginModal'
import styles from './CommentForm.module.css'

interface Props {
  onSubmit: (content: string) => Promise<void>
  isSubmitting: boolean
}

export function CommentForm({ onSubmit, isSubmitting }: Props) {
  const { user } = useAuth()
  const { openLoginModal } = useLoginModal()
  const [content, setContent] = useState('')

  if (!user) {
    return (
      <p className={styles.loginPrompt}>
        <button onClick={openLoginModal} className={styles.loginLink}>Log in</button> to leave a comment.
      </p>
    )
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim()) return
    await onSubmit(content)
    setContent('')
  }

  return (
    <form onSubmit={handleSubmit} className={styles.form}>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Write a comment..."
        rows={4}
        className={styles.textarea}
        required
      />
      <button
        type="submit"
        disabled={isSubmitting || !content.trim()}
        className={styles.submitBtn}
      >
        {isSubmitting ? 'Posting...' : 'Post comment'}
      </button>
    </form>
  )
}
