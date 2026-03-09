import { usePosts } from '../hooks/usePosts'
import { PostCard } from '../components/PostCard'

export function HomePage() {
  const { data: posts, isLoading, error } = usePosts()

  if (isLoading) return <div style={{ padding: '2rem', textAlign: 'center' }}>Loading posts...</div>
  if (error) return <div style={{ padding: '2rem', color: 'red' }}>Failed to load posts.</div>

  return (
    <main style={{ maxWidth: '800px', margin: '2rem auto', padding: '0 1rem' }}>
      <h1 style={{ marginBottom: '2rem' }}>Latest Posts</h1>
      {posts && posts.length > 0 ? (
        posts.map((post) => <PostCard key={post.id} post={post} />)
      ) : (
        <p style={{ color: '#888' }}>No posts yet.</p>
      )}
    </main>
  )
}
