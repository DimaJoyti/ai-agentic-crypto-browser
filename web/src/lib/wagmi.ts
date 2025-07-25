import { createConfig, http } from 'wagmi'
import { mainnet, polygon, arbitrum, optimism, sepolia } from 'wagmi/chains'
import { injected, walletConnect, coinbaseWallet } from 'wagmi/connectors'

const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || ''

export const config = createConfig({
  chains: [mainnet, polygon, arbitrum, optimism, sepolia],
  connectors: [
    injected(),
    walletConnect({ 
      projectId,
      metadata: {
        name: 'AI Agentic Browser',
        description: 'AI-powered web browser with Web3 integration',
        url: 'https://ai-agentic-browser.com',
        icons: ['https://ai-agentic-browser.com/icon.png']
      }
    }),
    coinbaseWallet({
      appName: 'AI Agentic Browser',
      appLogoUrl: 'https://ai-agentic-browser.com/icon.png'
    }),
  ],
  transports: {
    [mainnet.id]: http(),
    [polygon.id]: http(),
    [arbitrum.id]: http(),
    [optimism.id]: http(),
    [sepolia.id]: http(),
  },
})
