import { useState } from 'react'
import { useLocation, useParams } from 'react-router-dom'
import { useAuth } from '@hooks/useAuth'
import { useLoginModal } from '@hooks/useLoginModal'
import { useCreateComment } from '@hooks/useComments'
import styles from './CommentForm.module.css'

export function CommentForm() {
  const { slug = '' } = useParams<{ slug: string }>()
  const { user } = useAuth()
  const { openLoginModal } = useLoginModal()
  const location = useLocation()
  const createComment = useCreateComment(slug)
  const [content, setContent] = useState('')

  if (!user) {
    return (
      <p className={styles.loginPrompt}>
        <button onClick={() => openLoginModal(location.pathname + '#comments')} className={styles.loginLink}>Log in</button> to leave a comment.
      </p>
    )
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault()
    if (!content.trim()) return
    await createComment.mutateAsync(content)
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
      <div className={styles.actions}>
        <button
          type="submit"
          disabled={createComment.isPending || !content.trim()}
          className={styles.submitBtn}
        >
          {createComment.isPending ? 'Posting...' : 'Post comment'}
        </button>
        <button
          type="button"
          disabled={!content}
          onClick={() => setContent('')}
          className={styles.clearBtn}
        >
          Clear
        </button>
      </div>
    </form>
  )
}
