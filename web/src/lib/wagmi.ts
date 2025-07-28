import { createConfig, http } from 'wagmi'
import { mainnet, polygon, arbitrum, optimism, sepolia, base, avalanche, bsc, fantom, gnosis } from 'wagmi/chains'
import { injected, walletConnect, coinbaseWallet } from 'wagmi/connectors'

const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || ''

// Define supported chains as a tuple to satisfy Wagmi's type requirements
const chains = [mainnet, polygon, arbitrum, optimism, base, avalanche, bsc, fantom, gnosis, sepolia] as const

// Create transports for all supported chains
const transports = {
  [mainnet.id]: http(),
  [polygon.id]: http(),
  [arbitrum.id]: http(),
  [optimism.id]: http(),
  [base.id]: http(),
  [avalanche.id]: http(),
  [bsc.id]: http(),
  [fantom.id]: http(),
  [gnosis.id]: http(),
  [sepolia.id]: http(),
}

export const config = createConfig({
  chains,
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
  transports,
})
