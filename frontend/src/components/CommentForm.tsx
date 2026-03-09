import { useState } from 'react'
import { Link } from 'react-router-dom'
import { useAuth } from '../hooks/useAuth'

interface Props {
  onSubmit: (content: string) => Promise<void>
  isSubmitting: boolean
}

export function CommentForm({ onSubmit, isSubmitting }: Props) {
  const { user } = useAuth()
  const [content, setContent] = useState('')

  if (!user) {
    return (
      <p style={{ color: '#666', margin: '1rem 0' }}>
        <Link to="/login" style={{ color: '#0070f3' }}>Log in</Link> to leave a comment.
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
    <form onSubmit={handleSubmit} style={{ margin: '1.5rem 0' }}>
      <textarea
        value={content}
        onChange={(e) => setContent(e.target.value)}
        placeholder="Write a comment..."
        rows={4}
        style={{ width: '100%', padding: '0.75rem', borderRadius: '4px', border: '1px solid #ccc', fontSize: '0.95rem', resize: 'vertical', boxSizing: 'border-box' }}
        required
      />
      <button
        type="submit"
        disabled={isSubmitting || !content.trim()}
        style={{ marginTop: '0.5rem', padding: '0.5rem 1.25rem', background: '#0070f3', color: '#fff', border: 'none', borderRadius: '4px', cursor: 'pointer', fontSize: '0.95rem', opacity: isSubmitting ? 0.6 : 1 }}
      >
        {isSubmitting ? 'Posting...' : 'Post comment'}
      </button>
    </form>
  )
}
