import { usePosts } from '@hooks/usePosts'
import { PostCard } from '@components/PostCard'
import styles from './HomePage.module.css'

export function HomePage() {
  const { data: posts, isLoading, error } = usePosts()

  if (isLoading) return <div className={styles.loading}>Loading posts...</div>
  if (error) return <div className={styles.error}>Failed to load posts.</div>

  return (
    <main className={styles.main}>
      <h1 className={styles.heading}>Latest Posts</h1>
      {posts && posts.length > 0 ? (
        posts.map((post) => <PostCard key={post.id} post={post} />)
      ) : (
        <p className={styles.empty}>No posts yet.</p>
      )}
    </main>
  )
}
