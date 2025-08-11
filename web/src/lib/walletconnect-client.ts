'use client'

import { walletConnect } from 'wagmi/connectors'

// Client-side only WalletConnect connector factory
export const createWalletConnectConnector = (projectId: string) => {
  if (typeof window === 'undefined') {
    throw new Error('WalletConnect can only be initialized on the client side')
  }

  // Additional checks for browser APIs
  if (typeof indexedDB === 'undefined' || typeof localStorage === 'undefined') {
    throw new Error('WalletConnect requires browser storage APIs')
  }

  try {
    return walletConnect({
      projectId,
      metadata: {
        name: 'AI Agentic Browser',
        description: 'AI-powered web browser with Web3 integration and DeFi capabilities',
        url: window.location.origin,
        icons: [`${window.location.origin}/icons/icon-192x192.svg`]
      },
      showQrModal: true,
      qrModalOptions: {
        themeMode: 'light',
        themeVariables: {
          '--wcm-z-index': '1000'
        }
      }
    })
  } catch (error) {
    console.error('Failed to create WalletConnect connector:', error)
    throw error
  }
}

// Check if WalletConnect is available (client-side only)
export const isWalletConnectAvailable = () => {
  try {
    return typeof window !== 'undefined' &&
           typeof indexedDB !== 'undefined' &&
           typeof localStorage !== 'undefined' &&
           Boolean(process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID) &&
           process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID !== 'your_walletconnect_project_id_here'
  } catch (error) {
    console.warn('WalletConnect availability check failed:', error)
    return false
  }
}
