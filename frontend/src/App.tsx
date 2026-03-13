import { useState } from 'react'
import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { HelmetProvider } from 'react-helmet-async'
import { Header } from '@components/Header'
import { Footer } from '@components/Footer'
import { LoginModal } from '@components/LoginModal'
import { HomePage } from '@pages/HomePage'
import { PostPage } from '@pages/PostPage'
import { PrivacyPage } from '@pages/PrivacyPage'
import { LoginModalContext } from '@hooks/useLoginModal'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 60 * 1000,
    },
  },
})

export function App() {
  const [loginOpen, setLoginOpen] = useState(false)

  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <LoginModalContext.Provider value={{ openLoginModal: () => setLoginOpen(true) }}>
          <BrowserRouter>
            <Header />
            <div style={{ flex: 1 }}>
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/posts/:slug" element={<PostPage />} />
                <Route path="/privacy" element={<PrivacyPage />} />
              </Routes>
            </div>
            <Footer />
            <LoginModal open={loginOpen} onClose={() => setLoginOpen(false)} />
          </BrowserRouter>
        </LoginModalContext.Provider>
      </QueryClientProvider>
    </HelmetProvider>
  )
}
