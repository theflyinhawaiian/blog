import { BrowserRouter, Routes, Route } from 'react-router-dom'
import { QueryClient, QueryClientProvider } from '@tanstack/react-query'
import { HelmetProvider } from 'react-helmet-async'
import { Header } from '@components/Header'
import { Footer } from '@components/Footer'
import { LoginModalProvider } from '@components/LoginModalProvider'
import { ScrollToHash } from '@components/ScrollToHash'
import { BackToTop } from '@components/BackToTop'
import { HomePage } from '@pages/HomePage'
import { PostPage } from '@pages/PostPage'
import { TagPage } from '@pages/TagPage'
import { PrivacyPage } from '@pages/PrivacyPage'

const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      staleTime: 60 * 1000,
    },
  },
})

export function App() {
  return (
    <HelmetProvider>
      <QueryClientProvider client={queryClient}>
        <LoginModalProvider>
          <BrowserRouter>
            <ScrollToHash />
            <Header />
            <div style={{ flex: 1 }}>
              <Routes>
                <Route path="/" element={<HomePage />} />
                <Route path="/posts/:slug" element={<PostPage />} />
                <Route path="/tags/:tag" element={<TagPage />} />
                <Route path="/privacy" element={<PrivacyPage />} />
              </Routes>
            </div>
            <Footer />
            <BackToTop />
          </BrowserRouter>
        </LoginModalProvider>
      </QueryClientProvider>
    </HelmetProvider>
  )
}
