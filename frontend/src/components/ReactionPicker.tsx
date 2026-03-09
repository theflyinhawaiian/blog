import { useState, useRef, useEffect } from 'react'
import Picker from '@emoji-mart/react'
import data from '@emoji-mart/data'

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
    <div ref={ref} style={{ position: 'relative', display: 'inline-block' }}>
      <button
        onClick={() => setOpen((v) => !v)}
        style={{ background: 'none', border: '1px solid #ddd', borderRadius: '4px', cursor: 'pointer', padding: '0.2rem 0.5rem', fontSize: '1rem' }}
        title="Add reaction"
      >
        😊
      </button>
      {open && (
        <div style={{ position: 'absolute', bottom: '2.5rem', left: 0, zIndex: 100 }}>
          <Picker data={data} onEmojiSelect={handleEmojiSelect} theme="light" />
        </div>
      )}
    </div>
  )
}
