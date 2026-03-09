interface Props {
  src: string
  alt: string
}

export function HeroImage({ src, alt }: Props) {
  return (
    <div style={{ width: '100%', maxHeight: '400px', overflow: 'hidden', marginBottom: '2rem' }}>
      <img
        src={src}
        alt={alt}
        style={{ width: '100%', height: '400px', objectFit: 'cover' }}
      />
    </div>
  )
}
