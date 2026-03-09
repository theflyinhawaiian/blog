import { useEffect, useState } from 'react'

export function useReadingProgress(): number {
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    let rafId: number

    function update() {
      const el = document.documentElement
      const scrollTop = el.scrollTop || document.body.scrollTop
      const scrollHeight = el.scrollHeight - el.clientHeight
      if (scrollHeight <= 0) {
        setProgress(0)
        return
      }
      setProgress(Math.min(100, (scrollTop / scrollHeight) * 100))
    }

    function onScroll() {
      cancelAnimationFrame(rafId)
      rafId = requestAnimationFrame(update)
    }

    window.addEventListener('scroll', onScroll, { passive: true })
    update()

    return () => {
      window.removeEventListener('scroll', onScroll)
      cancelAnimationFrame(rafId)
    }
  }, [])

  return progress
}
