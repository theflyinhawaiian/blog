import { Helmet } from 'react-helmet-async'
import { useParams } from 'react-router-dom'
import { usePosts } from '@hooks/usePosts'
import { PostCard } from '@components/PostCard'
import styles from './HomePage.module.css'

export function TagPage() {
  const { tag } = useParams<{ tag: string }>()
  const { data: posts, isLoading, error } = usePosts()
  const filtered = posts?.filter(p => p.tags?.includes(tag ?? ''))

  if (isLoading) return <div className={styles.loading}>Loading posts...</div>
  if (error) return <div className={styles.error}>Failed to load posts.</div>

  return (
    <>
      <Helmet>
        <title>#{tag} — A random collection of thoughts</title>
      </Helmet>
      <main className={styles.main}>
        <h1 className={styles.heading}>#{tag}</h1>
        {filtered && filtered.length > 0 ? (
          filtered.map(post => <PostCard key={post.id} post={post} />)
        ) : (
          <p className={styles.empty}>No posts tagged #{tag}.</p>
        )}
      </main>
    </>
  )
}
