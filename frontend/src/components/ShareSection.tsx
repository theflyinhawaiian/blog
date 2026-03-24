import { useState } from 'react'
import { FaFacebook, FaXTwitter, FaLinkedin, FaEnvelope, FaLink, FaCheck } from 'react-icons/fa6'
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
        aria-label="Share on Facebook"
      >
        <FaFacebook />
      </a>
      <a
        href={`https://twitter.com/intent/tweet?url=${encoded}&text=${encodedTitle}`}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.shareLink}
        style={{ background: '#000' }}
        aria-label="Share on X"
      >
        <FaXTwitter />
      </a>
      <a
        href={`https://www.linkedin.com/sharing/share-offsite/?url=${encoded}`}
        target="_blank"
        rel="noopener noreferrer"
        className={styles.shareLink}
        style={{ background: '#0077b5' }}
        aria-label="Share on LinkedIn"
      >
        <FaLinkedin />
      </a>
      <a
        href={`mailto:?subject=${encodedTitle}&body=${encoded}`}
        className={styles.shareLink}
        style={{ background: '#555' }}
        aria-label="Share via email"
      >
        <FaEnvelope />
      </a>
      <button onClick={copyLink} className={`${styles.copyBtn} ${copied ? styles.copyBtnSuccess : ''}`} aria-label={copied ? 'Copied' : 'Copy link'}>
        {copied ? <FaCheck /> : <FaLink />}
      </button>
    </div>
  )
}
