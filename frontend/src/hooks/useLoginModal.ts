import { createContext, useContext } from 'react'

interface LoginModalContextValue {
  openLoginModal: () => void
}

export const LoginModalContext = createContext<LoginModalContextValue>({
  openLoginModal: () => {},
})

export function useLoginModal() {
  return useContext(LoginModalContext)
}
