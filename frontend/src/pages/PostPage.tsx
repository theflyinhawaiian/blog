import { useRef } from 'react'
import { Link, useParams } from 'react-router-dom'
import { Helmet } from 'react-helmet-async'
import { usePost } from '@hooks/usePosts'
import { useComments } from '@hooks/useComments'
import { MarkdownContent } from '@components/MarkdownContent'
import { HeroImage } from '@components/HeroImage'
import { ReadingProgressBar } from '@components/ReadingProgressBar'
import { ShareSection } from '@components/ShareSection'
import { CommentList } from '@components/CommentList'
import styles from './PostPage.module.css'

export function PostPage() {
  const contentRef = useRef<HTMLElement>(null)
  const { slug } = useParams<{ slug: string }>()
  const { data: post, isLoading, error } = usePost(slug ?? '')
  const { data: comments = [] } = useComments(slug ?? '')

  if (isLoading) return <div className={styles.loading}>Loading post...</div>
  if (error || !post) return <div className={styles.error}>Post not found.</div>

  const postUrl = window.location.href
  const metaDesc = post.meta_description?.Valid ? post.meta_description.String : ''
  const canonicalUrl = post.canonical_url?.Valid ? post.canonical_url.String : postUrl
  const heroImage = post.post_image?.Valid ? post.post_image.String : ''

  return (
    <>
      <Helmet>
        <title>Pete's blog - {post.title}</title>
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

      <ReadingProgressBar targetRef={contentRef} />

      {heroImage && <HeroImage src={heroImage} alt={post.title} />}

      <main ref={contentRef} className={styles.main}>
        <h1 className={styles.title}>{post.title}</h1>
        <time className={styles.date}>
          {new Date(post.created_at).toLocaleDateString()}
        </time>

        {post.tags && post.tags.length > 0 && (
          <div className={styles.tags}>
            {post.tags.map(tag => (
              <Link key={tag} to={`/tags/${tag}`} className={styles.tag}>
                #{tag}
              </Link>
            ))}
          </div>
        )}

        <MarkdownContent content={post.content} className={`post-content ${styles.content}`} />

        <ShareSection url={postUrl} title={post.title} />

        <CommentList comments={comments} />
      </main>
    </>
  )
}
