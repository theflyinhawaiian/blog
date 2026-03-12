import { Comment } from '@typedef/comment'
import { ReactionPicker } from './ReactionPicker'
import { useAuth } from '@hooks/useAuth'
import { MarkdownContent } from './MarkdownContent'
import styles from './CommentItem.module.css'

interface Props {
  comment: Comment
  onReact: (commentId: number, emoji: string) => void
}

export function CommentItem({ comment, onReact }: Props) {
  const { user } = useAuth()

  return (
    <div className={styles.item}>
      <div className={styles.itemHeader}>
        <strong className={styles.author}>{comment.display_name}</strong>
        <time className={styles.date}>
          {new Date(comment.created_at).toLocaleString()}
        </time>
      </div>
      <MarkdownContent content={comment.content} className={`comment-content ${styles.content}`} />
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
    </div>
  )
}
