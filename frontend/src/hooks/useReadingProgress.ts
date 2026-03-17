import { useEffect, useState, RefObject } from 'react'

export function useReadingProgress(targetRef: RefObject<HTMLElement | null>): number {
  const [progress, setProgress] = useState(0)

  useEffect(() => {
    let rafId: number

    function update() {
      const el = targetRef.current
      if (!el) return

      const elTop = el.getBoundingClientRect().top + window.scrollY
      const scrollable = el.offsetHeight - window.innerHeight
      if (scrollable <= 0) {
        setProgress(100)
        return
      }
      const scrolled = window.scrollY - elTop
      setProgress(Math.min(100, Math.max(0, (scrolled / scrollable) * 100)))
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
  }, [targetRef])

  return progress
}
