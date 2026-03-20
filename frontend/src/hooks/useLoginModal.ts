import { createContext, useContext, useState } from 'react'

interface LoginModalContextValue {
  openLoginModal: (returnTo?: string) => void
  open: boolean
  returnTo: string | undefined
  close: () => void
}

export const LoginModalContext = createContext<LoginModalContextValue>({
  openLoginModal: () => {},
  open: false,
  returnTo: undefined,
  close: () => {},
})

export function useLoginModal() {
  return useContext(LoginModalContext)
}

export function useLoginModalState() {
  const [open, setOpen] = useState(false)
  const [returnTo, setReturnTo] = useState<string | undefined>()

  return {
    open,
    returnTo,
    openLoginModal: (r?: string) => { setReturnTo(r); setOpen(true) },
    close: () => { setOpen(false); setReturnTo(undefined) },
  }
}
