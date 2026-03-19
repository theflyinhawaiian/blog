import { useState } from 'react'
import { useParams } from 'react-router-dom'
import { DateTime } from 'luxon'
import { Comment } from '@typedef/comment'
import { ReactionPicker } from './ReactionPicker'
import { useAuth } from '@hooks/useAuth'
import { useUpdateComment, useDeleteComment, useAddReaction } from '@hooks/useComments'
import { MarkdownContent } from './MarkdownContent'
import styles from './CommentItem.module.css'

function decodeHtml(html: string): string {
  const el = document.createElement('textarea')
  el.innerHTML = html
  return el.value
}

interface Props {
  comment: Comment
}

export function CommentItem({ comment }: Props) {
  const { slug = '' } = useParams<{ slug: string }>()
  const { user } = useAuth()
  const updateComment = useUpdateComment(slug)
  const deleteComment = useDeleteComment(slug)
  const addReaction = useAddReaction(slug)
  const [confirming, setConfirming] = useState(false)
  const [editing, setEditing] = useState(false)
  const [editContent, setEditContent] = useState(decodeHtml(comment.content))

  console.log(`user: ${JSON.stringify(user)}`);
  console.log(`comment's user: ${JSON.stringify(comment.user_id)}`);

  const isOwner = user?.id === comment.user_id

  function handleSave() {
    if (!editContent.trim() || editContent === comment.content) {
      setEditing(false)
      return
    }
    updateComment.mutate(
      { commentId: comment.id, content: editContent },
      { onSuccess: () => setEditing(false) },
    )
  }

  return (
    <div className={styles.item}>
      <div className={styles.itemHeader}>
        <strong className={styles.author}>{comment.display_name}</strong>
        <div className={styles.timestamps}>
          <time
            className={styles.date}
            title={DateTime.fromISO(comment.created_at).toLocaleString(DateTime.DATETIME_FULL)}
          >
            {DateTime.fromISO(comment.created_at).toRelative()}
          </time>
          {comment.updated_at && (
            <span
              className={styles.edited}
              title={DateTime.fromISO(comment.updated_at).toLocaleString(DateTime.DATETIME_FULL)}
            >
              edited {DateTime.fromISO(comment.updated_at).toRelative()}
            </span>
          )}
        </div>
      </div>

      {editing ? (
        <div className={styles.editBox}>
          <textarea
            className={styles.editTextarea}
            value={editContent}
            onChange={(e) => setEditContent(e.target.value)}
            rows={4}
          />
          <div className={styles.editActions}>
            <button
              className={styles.saveBtn}
              onClick={handleSave}
              disabled={updateComment.isPending}
            >
              {updateComment.isPending ? 'Saving...' : 'Save'}
            </button>
            <button
              className={styles.cancelBtn}
              onClick={() => { setEditing(false); setEditContent(decodeHtml(comment.content)) }}
            >
              Cancel
            </button>
          </div>
        </div>
      ) : (
        <MarkdownContent content={comment.content} className={`comment-content ${styles.content}`} />
      )}

      <div className={styles.footer}>
        <div className={styles.reactions}>
          {comment.reactions.filter((r) => r.count > 0).map((r) => (
            <button
              key={r.emoji}
              onClick={() => addReaction.mutate({ commentId: comment.id, emoji: r.emoji })}
              className={`${styles.reactionBtn} ${r.reacted_by_me ? styles.reactionBtnActive : ''}`}
            >
              {r.emoji} {r.count}
            </button>
          ))}
          {user && <ReactionPicker onSelect={(emoji) => addReaction.mutate({ commentId: comment.id, emoji })} />}
        </div>
        {isOwner && !editing && (
          <div className={styles.ownerActions}>
            <button className={styles.editBtn} onClick={() => setEditing(true)}>Edit</button>
            {confirming ? (
              <span className={styles.confirmDelete}>
                Delete?{' '}
                <button onClick={() => { deleteComment.mutate(comment.id); setConfirming(false) }}>Yes</button>
                {' '}
                <button onClick={() => setConfirming(false)}>No</button>
              </span>
            ) : (
              <button className={styles.deleteBtn} onClick={() => setConfirming(true)}>Delete</button>
            )}
          </div>
        )}
      </div>
    </div>
  )
}
