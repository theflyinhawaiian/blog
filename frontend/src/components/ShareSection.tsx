import { useState } from 'react'
import styles from './ShareSection.module.css'

interface Props {
  url: string
  title: string
}

export function ShareSection({ url, title }: Props) {
  const encoded = encodeURIComponent(url)
  const encodedTitle = encodeURIComponent(title)
  const [copied, setCopied] = useState(false)

  function copyLink() {
    navigator.clipboard.writeText(url).then(() => {
      setCopied(true)
      setTimeout(() => setCopied(false), 2000)
    })
  }

  return (
    <div className={styles.section}>
      <span className={styles.label}>Share:</span>
      <a
        href={`https://www.facebook.com/sharer/sharer.php?u=${encoded}`}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.shareLink}
        style={{ background: '#1877f2' }}
      >
        Facebook
      </a>
      <a
        href={`https://twitter.com/intent/tweet?url=${encoded}&text=${encodedTitle}`}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.shareLink}
        style={{ background: '#1da1f2' }}
      >
        Twitter
      </a>
      <a
        href={`https://www.linkedin.com/sharing/share-offsite/?url=${encoded}`}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.shareLink}
        style={{ background: '#0077b5' }}
      >
        LinkedIn
      </a>
      <a
        href={`mailto:?subject=${encodedTitle}&body=${encoded}`}
        className={styles.shareLink}
        style={{ background: '#555' }}
      >
        Email
      </a>
      <button onClick={copyLink} className={`${styles.copyBtn} ${copied ? styles.copyBtnSuccess : ''}`}>
        {copied ? 'Copied!' : 'Copy link'}
      </button>
    </div>
  )
}
