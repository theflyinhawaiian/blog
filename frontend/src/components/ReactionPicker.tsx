import { useState, useRef, useEffect } from 'react'
import Picker from '@emoji-mart/react'
import data from '@emoji-mart/data'
import styles from './ReactionPicker.module.css'

interface Props {
  onSelect: (emoji: string) => void
}

export function ReactionPicker({ onSelect }: Props) {
  const [open, setOpen] = useState(false)
  const ref = useRef<HTMLDivElement>(null)

  useEffect(() => {
    function handleClickOutside(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false)
      }
    }
    if (open) document.addEventListener('mousedown', handleClickOutside)
    return () => document.removeEventListener('mousedown', handleClickOutside)
  }, [open])

  function handleEmojiSelect(emoji: { native: string }) {
    onSelect(emoji.native)
    setOpen(false)
  }

  return (
    <div ref={ref} className={styles.wrapper}>
      <button onClick={() => setOpen((v) => !v)} className={styles.trigger} title="Add reaction">
        ➕
      </button>
      {open && (
        <div className={styles.picker}>
          <Picker data={data} onEmojiSelect={handleEmojiSelect} theme="light" />
        </div>
      )}
    </div>
  )
}
