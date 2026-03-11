import styles from './HeroImage.module.css'

interface Props {
  src: string
  alt: string
}

export function HeroImage({ src, alt }: Props) {
  return (
    <div className={styles.wrapper}>
      <img src={src} alt={alt} className={styles.image} />
    </div>
  )
}
