import { Link } from 'react-router-dom'
import { PostSummary } from '../types/post'

interface Props {
  post: PostSummary
}

export function PostCard({ post }: Props) {
  const excerpt = post.excerpt?.Valid ? post.excerpt.String : ''
  const image = post.post_image?.Valid ? post.post_image.String : ''

  return (
    <article style={{ border: '1px solid #eee', borderRadius: '8px', overflow: 'hidden', marginBottom: '1.5rem' }}>
      {image && (
        <img src={image} alt={post.title} style={{ width: '100%', height: '200px', objectFit: 'cover' }} />
      )}
      <div style={{ padding: '1.25rem' }}>
        <Link to={`/posts/${post.slug}`} style={{ textDecoration: 'none', color: '#111' }}>
          <h2 style={{ margin: '0 0 0.5rem', fontSize: '1.4rem' }}>{post.title}</h2>
        </Link>
        {excerpt && <p style={{ color: '#555', margin: '0 0 0.75rem', lineHeight: 1.6 }}>{excerpt}</p>}
        <time style={{ fontSize: '0.85rem', color: '#999' }}>
          {new Date(post.created_at).toLocaleDateString()}
        </time>
      </div>
    </article>
  )
}
