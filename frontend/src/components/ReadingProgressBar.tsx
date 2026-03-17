import { RefObject } from 'react'
import { useReadingProgress } from '@hooks/useReadingProgress'
import styles from './ReadingProgressBar.module.css'

export function ReadingProgressBar({ targetRef }: { targetRef: RefObject<HTMLElement | null> }) {
  const progress = useReadingProgress(targetRef)

  return (
    <div className={styles.bar} style={{ width: `${progress}%` }} />
  )
}
