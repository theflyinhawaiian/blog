import { Link } from 'react-router-dom'
import { PostSummary } from '@typedef/post'
import styles from './PostCard.module.css'

interface Props {
  post: PostSummary
}

export function PostCard({ post }: Props) {
  const excerpt = post.excerpt?.Valid ? post.excerpt.String : ''
  const image = post.post_image?.Valid ? post.post_image.String : ''

  return (
    <Link to={`/posts/${post.slug}`} className={styles.card}>
      <article>
        {image && (
          <img src={image} alt={post.title} className={styles.cardImage} />
        )}
        <div className={styles.cardBody}>
          <h2 className={styles.cardTitle}>{post.title}</h2>
          {excerpt && <p className={styles.cardExcerpt}>{excerpt}</p>}
          <time className={styles.cardDate}>
            {new Date(post.created_at).toLocaleDateString()}
          </time>
        </div>
      </article>
    </Link>
  )
}
