import { useReadingProgress } from '@hooks/useReadingProgress'
import styles from './ReadingProgressBar.module.css'

export function ReadingProgressBar() {
  const progress = useReadingProgress()

  return (
    <div className={styles.bar} style={{ width: `${progress}%` }} />
  )
}
