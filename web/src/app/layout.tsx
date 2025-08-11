import type { Metadata } from 'next'
import { Inter } from 'next/font/google'
import './globals.css'
import { Providers } from '@/components/providers'
import { AccessibilityProvider } from '@/components/accessibility/AccessibilityProvider'
import { Toaster } from 'sonner'
import { ChunkErrorBoundary } from '@/components/ChunkErrorBoundary'

const inter = Inter({ subsets: ['latin'] })

export const metadata: Metadata = {
  metadataBase: new URL(process.env.NEXT_PUBLIC_APP_URL || 'http://localhost:3000'),
  title: 'AI Agentic Crypto Browser',
  description: 'An intelligent web browser powered by AI agents with Web3 integration and advanced trading capabilities',
  keywords: ['AI', 'browser', 'automation', 'Web3', 'cryptocurrency', 'DeFi', 'trading', 'blockchain'],
  authors: [{ name: 'AI Agentic Browser Team' }],
  manifest: '/manifest.json',
  appleWebApp: {
    capable: true,
    statusBarStyle: 'default',
    title: 'AI Browser',
  },
  formatDetection: {
    telephone: false,
  },
  openGraph: {
    type: 'website',
    siteName: 'AI Agentic Crypto Browser',
    title: 'AI Agentic Crypto Browser',
    description: 'Advanced AI-powered trading and Web3 browser',
    images: [
      {
        url: '/icons/icon-512x512.svg',
        width: 512,
        height: 512,
        alt: 'AI Agentic Crypto Browser',
      },
    ],
  },
  twitter: {
    card: 'summary_large_image',
    title: 'AI Agentic Crypto Browser',
    description: 'Advanced AI-powered trading and Web3 browser',
    images: ['/icons/icon-512x512.svg'],
  },
}

export const viewport = {
  width: 'device-width',
  initialScale: 1,
  maximumScale: 1,
  userScalable: false,
  themeColor: [
    { media: '(prefers-color-scheme: light)', color: '#ffffff' },
    { media: '(prefers-color-scheme: dark)', color: '#000000' },
  ],
}

export default function RootLayout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html lang="en" suppressHydrationWarning>
      <head>
        <link rel="manifest" href="/manifest.json" />
        <link rel="apple-touch-icon" href="/icons/icon-192x192.svg" />
        <meta name="apple-mobile-web-app-capable" content="yes" />
        <meta name="apple-mobile-web-app-status-bar-style" content="default" />
        <meta name="apple-mobile-web-app-title" content="AI Browser" />
        <meta name="mobile-web-app-capable" content="yes" />
        <meta name="msapplication-TileColor" content="#3b82f6" />
        <meta name="msapplication-tap-highlight" content="no" />
      </head>
<<<<<<< HEAD
      <body className={inter.className} suppressHydrationWarning>
=======
      <body className={inter.className}>
>>>>>>> d850c235d1b366ccb3b4e75eebc09fc566798249
        <ChunkErrorBoundary>
          <Providers>
            <AccessibilityProvider>
              <div className="min-h-screen bg-background" id="main-content">
                {children}
              </div>
              <Toaster richColors position="top-right" />
            </AccessibilityProvider>
          </Providers>
        </ChunkErrorBoundary>

        {/* Service Worker Registration */}
        <script
          dangerouslySetInnerHTML={{
            __html: `
              if ('serviceWorker' in navigator) {
                window.addEventListener('load', function() {
                  navigator.serviceWorker.register('/sw.js')
                    .then(function(registration) {
                      console.log('SW registered: ', registration);
                    })
                    .catch(function(registrationError) {
                      console.log('SW registration failed: ', registrationError);
                    });
                });
              }
            `,
          }}
        />
      </body>
    </html>
  )
}
