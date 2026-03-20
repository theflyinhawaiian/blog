import { LoginModalContext, useLoginModalState } from '@hooks/useLoginModal'
import { LoginModal } from './LoginModal'

interface Props {
  children: React.ReactNode
}

export function LoginModalProvider({ children }: Props) {
  const state = useLoginModalState()

  return (
    <LoginModalContext.Provider value={state}>
      {children}
      <LoginModal open={state.open} returnTo={state.returnTo} onClose={state.close} />
    </LoginModalContext.Provider>
  )
}
