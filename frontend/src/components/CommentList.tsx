import { Comment } from '@typedef/comment'
import { CommentItem } from './CommentItem'
import { CommentForm } from './CommentForm'
import styles from './CommentList.module.css'

interface Props {
  comments: Comment[]
}

export function CommentList({ comments }: Props) {
  return (
    <section id="comments" className={styles.section}>
      <h2 className={styles.heading}>
        Comments ({comments.length})
      </h2>
      {comments.length === 0 ? (
        <p className={styles.empty}>No comments yet. Be the first!</p>
      ) : (
        comments.map((c) => (
          <CommentItem key={c.id} comment={c} />
        ))
      )}
      <CommentForm />
    </section>
  )
}
