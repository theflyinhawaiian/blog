import { useParams } from 'react-router-dom'
import { Helmet } from 'react-helmet-async'
import ReactMarkdown from 'react-markdown'
import rehypeHighlight from 'rehype-highlight'
import 'highlight.js/styles/github.css'
import { usePost } from '../hooks/usePosts'
import { useComments, useCreateComment, useAddReaction } from '../hooks/useComments'
import { HeroImage } from '../components/HeroImage'
import { ReadingProgressBar } from '../components/ReadingProgressBar'
import { ShareSection } from '../components/ShareSection'
import { CommentList } from '../components/CommentList'

export function PostPage() {
  const { slug } = useParams<{ slug: string }>()
  const { data: post, isLoading, error } = usePost(slug ?? '')
  const { data: comments = [] } = useComments(slug ?? '')
  const createComment = useCreateComment(slug ?? '')
  const addReaction = useAddReaction(slug ?? '')

  if (isLoading) return <div style={{ padding: '2rem', textAlign: 'center' }}>Loading post...</div>
  if (error || !post) return <div style={{ padding: '2rem', color: 'red' }}>Post not found.</div>

  const postUrl = window.location.href
  const metaDesc = post.meta_description?.Valid ? post.meta_description.String : ''
  const canonicalUrl = post.canonical_url?.Valid ? post.canonical_url.String : postUrl
  const heroImage = post.post_image?.Valid ? post.post_image.String : ''

  return (
    <>
      <Helmet>
        <title>{post.title}</title>
        {metaDesc && <meta name="description" content={metaDesc} />}
        <link rel="canonical" href={canonicalUrl} />
        <meta property="og:title" content={post.title} />
        <meta property="og:type" content="article" />
        <meta property="og:url" content={canonicalUrl} />
        {heroImage && <meta property="og:image" content={heroImage} />}
        {metaDesc && <meta property="og:description" content={metaDesc} />}
        <meta name="twitter:card" content="summary_large_image" />
        <meta name="twitter:title" content={post.title} />
        {metaDesc && <meta name="twitter:description" content={metaDesc} />}
        {heroImage && <meta name="twitter:image" content={heroImage} />}
      </Helmet>

      <ReadingProgressBar />

      {heroImage && <HeroImage src={heroImage} alt={post.title} />}

      <main style={{ maxWidth: '800px', margin: '2rem auto', padding: '0 1rem' }}>
        <h1 style={{ fontSize: '2rem', marginBottom: '0.5rem' }}>{post.title}</h1>
        <time style={{ color: '#999', fontSize: '0.9rem', display: 'block', marginBottom: '2rem' }}>
          {new Date(post.created_at).toLocaleDateString()}
        </time>

        <div className="post-content" style={{ lineHeight: 1.8, fontSize: '1.05rem' }}>
          <ReactMarkdown rehypePlugins={[rehypeHighlight]}>{post.content}</ReactMarkdown>
        </div>

        <ShareSection url={postUrl} title={post.title} />

        <CommentList
          comments={comments}
          onSubmit={async (content) => { await createComment.mutateAsync(content) }}
          onReact={(commentId, emoji) => addReaction.mutate({ commentId, emoji })}
          isSubmitting={createComment.isPending}
        />
      </main>
    </>
  )
}
