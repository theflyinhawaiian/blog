interface Props {
  url: string
  title: string
}

export function ShareSection({ url, title }: Props) {
  const encoded = encodeURIComponent(url)
  const encodedTitle = encodeURIComponent(title)

  function copyLink() {
    navigator.clipboard.writeText(url)
  }

  return (
    <div style={{ margin: '2rem 0', display: 'flex', gap: '0.75rem', flexWrap: 'wrap', alignItems: 'center' }}>
      <span style={{ fontWeight: '600', marginRight: '0.5rem' }}>Share:</span>
      <a
        href={`https://www.facebook.com/sharer/sharer.php?u=${encoded}`}
        target="_blank"
        rel="noopener noreferrer"
        style={shareStyle('#1877f2')}
      >
        Facebook
      </a>
      <a
        href={`https://twitter.com/intent/tweet?url=${encoded}&text=${encodedTitle}`}
        target="_blank"
        rel="noopener noreferrer"
        style={shareStyle('#1da1f2')}
      >
        Twitter
      </a>
      <a
        href={`https://www.linkedin.com/sharing/share-offsite/?url=${encoded}`}
        target="_blank"
        rel="noopener noreferrer"
        style={shareStyle('#0077b5')}
      >
        LinkedIn
      </a>
      <a
        href={`mailto:?subject=${encodedTitle}&body=${encoded}`}
        style={shareStyle('#555')}
      >
        Email
      </a>
      <button onClick={copyLink} style={{ ...shareStyle('#333'), cursor: 'pointer', border: '1px solid #ccc', background: '#f5f5f5' }}>
        Copy link
      </button>
    </div>
  )
}

function shareStyle(color: string): React.CSSProperties {
  return {
    padding: '0.4rem 1rem',
    borderRadius: '4px',
    textDecoration: 'none',
    color: '#fff',
    background: color,
    fontSize: '0.875rem',
    display: 'inline-block',
  }
}
