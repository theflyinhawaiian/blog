import { Comment } from '../types/comment'
import { CommentItem } from './CommentItem'
import { CommentForm } from './CommentForm'

interface Props {
  comments: Comment[]
  onSubmit: (content: string) => Promise<void>
  onReact: (commentId: number, emoji: string) => void
  isSubmitting: boolean
}

export function CommentList({ comments, onSubmit, onReact, isSubmitting }: Props) {
  return (
    <section style={{ marginTop: '3rem' }}>
      <h2 style={{ borderBottom: '2px solid #eee', paddingBottom: '0.5rem' }}>
        Comments ({comments.length})
      </h2>
      <CommentForm onSubmit={onSubmit} isSubmitting={isSubmitting} />
      {comments.length === 0 ? (
        <p style={{ color: '#888' }}>No comments yet. Be the first!</p>
      ) : (
        comments.map((c) => (
          <CommentItem key={c.id} comment={c} onReact={onReact} />
        ))
      )}
    </section>
  )
}
