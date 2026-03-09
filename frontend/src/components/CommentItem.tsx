import { Comment } from '../types/comment'
import { ReactionPicker } from './ReactionPicker'
import { useAuth } from '../hooks/useAuth'

interface Props {
  comment: Comment
  onReact: (commentId: number, emoji: string) => void
}

export function CommentItem({ comment, onReact }: Props) {
  const { user } = useAuth()

  return (
    <div style={{ borderBottom: '1px solid #eee', padding: '1rem 0' }}>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: '0.4rem' }}>
        <strong style={{ fontSize: '0.9rem' }}>{comment.display_name}</strong>
        <time style={{ fontSize: '0.8rem', color: '#999' }}>
          {new Date(comment.created_at).toLocaleString()}
        </time>
      </div>
      <div
        style={{ lineHeight: 1.6, color: '#333' }}
        dangerouslySetInnerHTML={{ __html: comment.content }}
      />
      <div style={{ marginTop: '0.5rem', display: 'flex', gap: '0.4rem', flexWrap: 'wrap', alignItems: 'center' }}>
        {comment.reactions.map((r) => (
          <button
            key={r.emoji}
            onClick={() => onReact(comment.id, r.emoji)}
            style={{ background: '#f5f5f5', border: '1px solid #ddd', borderRadius: '12px', padding: '0.15rem 0.5rem', cursor: 'pointer', fontSize: '0.85rem' }}
          >
            {r.emoji} {r.count}
          </button>
        ))}
        {user && <ReactionPicker onSelect={(emoji) => onReact(comment.id, emoji)} />}
      </div>
    </div>
  )
}
