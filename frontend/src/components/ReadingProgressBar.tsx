import { useReadingProgress } from '../hooks/useReadingProgress'

export function ReadingProgressBar() {
  const progress = useReadingProgress()

  return (
    <div
      style={{
        position: 'fixed',
        bottom: 0,
        left: 0,
        width: `${progress}%`,
        height: '4px',
        background: '#0070f3',
        transition: 'width 0.1s linear',
        zIndex: 1000,
      }}
    />
  )
}
