import { Comment } from '../types/comment'
import { CommentItem } from './CommentItem'
import { CommentForm } from './CommentForm'
import styles from './CommentList.module.css'

interface Props {
  comments: Comment[]
  onSubmit: (content: string) => Promise<void>
  onReact: (commentId: number, emoji: string) => void
  isSubmitting: boolean
}

export function CommentList({ comments, onSubmit, onReact, isSubmitting }: Props) {
  return (
    <section className={styles.section}>
      <h2 className={styles.heading}>
        Comments ({comments.length})
      </h2>
      <CommentForm onSubmit={onSubmit} isSubmitting={isSubmitting} />
      {comments.length === 0 ? (
        <p className={styles.empty}>No comments yet. Be the first!</p>
      ) : (
        comments.map((c) => (
          <CommentItem key={c.id} comment={c} onReact={onReact} />
        ))
      )}
    </section>
  )
}
