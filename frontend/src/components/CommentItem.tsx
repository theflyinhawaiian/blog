import { useState } from 'react'
import { Comment } from '@typedef/comment'
import { ReactionPicker } from './ReactionPicker'
import { useAuth } from '@hooks/useAuth'
import { MarkdownContent } from './MarkdownContent'
import styles from './CommentItem.module.css'

interface Props {
  comment: Comment
  onReact: (commentId: number, emoji: string) => void
  onDelete: (commentId: number) => void
}

export function CommentItem({ comment, onReact, onDelete }: Props) {
  const { user } = useAuth()
  const [confirming, setConfirming] = useState(false)

  return (
    <div className={styles.item}>
      <div className={styles.itemHeader}>
        <strong className={styles.author}>{comment.display_name}</strong>
        <time className={styles.date}>
          {new Date(comment.created_at).toLocaleString()}
        </time>
      </div>
      <MarkdownContent content={comment.content} className={`comment-content ${styles.content}`} />
      <div className={styles.footer}>
        <div className={styles.reactions}>
          {comment.reactions.filter((r) => r.count > 0).map((r) => (
            <button
              key={r.emoji}
              onClick={() => onReact(comment.id, r.emoji)}
              className={`${styles.reactionBtn} ${r.reacted_by_me ? styles.reactionBtnActive : ''}`}
            >
              {r.emoji} {r.count}
            </button>
          ))}
          {user && <ReactionPicker onSelect={(emoji) => onReact(comment.id, emoji)} />}
        </div>
        {user?.id === comment.user_id && (
          confirming ? (
            <span className={styles.confirmDelete}>
              Delete?{' '}
              <button onClick={() => { onDelete(comment.id); setConfirming(false) }}>Yes</button>
              {' '}
              <button onClick={() => setConfirming(false)}>No</button>
            </span>
          ) : (
            <button className={styles.deleteBtn} onClick={() => setConfirming(true)}>
              Delete
            </button>
          )
        )}
      </div>
    </div>
  )
}
