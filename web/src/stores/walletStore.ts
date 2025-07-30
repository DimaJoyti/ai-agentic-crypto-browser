import { create } from 'zustand'
import { persist, createJSONStorage } from 'zustand/middleware'
import { subscribeWithSelector } from 'zustand/middleware'
import { type Address } from 'viem'
import { type Connector } from 'wagmi'
import { 
  type WalletProvider, 
  type WalletConnectionState,
  WALLET_PROVIDERS 
} from '@/lib/wallet-providers'

// Types for wallet state management
export interface ConnectedWallet {
  id: string
  address: Address
  chainId: number
  connector: string
  provider: WalletProvider
  connectedAt: number
  lastUsed: number
  isActive: boolean
  balance?: string
  ensName?: string
  avatar?: string
}

export interface WalletSession {
  id: string
  walletId: string
  address: Address
  chainId: number
  startTime: number
  endTime?: number
  duration?: number
  transactionCount: number
  lastActivity: number
}

export interface WalletPreferences {
  autoConnect: boolean
  preferredWallet?: string
  preferredChain?: number
  showTestnets: boolean
  hideSmallBalances: boolean
  defaultSlippage: number
  gasPreference: 'slow' | 'standard' | 'fast' | 'custom'
  customGasPrice?: string
  notifications: {
    transactions: boolean
    priceAlerts: boolean
    news: boolean
  }
}

export interface WalletAnalytics {
  totalConnections: number
  totalTransactions: number
  totalValueTransacted: string
  favoriteChains: Record<number, number>
  walletUsageStats: Record<string, {
    connections: number
    lastUsed: number
    totalTime: number
  }>
  monthlyStats: Record<string, {
    connections: number
    transactions: number
    volume: string
  }>
}

interface WalletState {
  // Connection state
  isConnected: boolean
  isConnecting: boolean
  connectionError: string | null
  
  // Current wallet
  currentWallet: ConnectedWallet | null
  currentSession: WalletSession | null
  
  // Wallet management
  connectedWallets: ConnectedWallet[]
  walletHistory: ConnectedWallet[]
  sessions: WalletSession[]
  
  // User preferences
  preferences: WalletPreferences
  
  // Analytics
  analytics: WalletAnalytics
  
  // Available wallets
  availableWallets: WalletProvider[]
  
  // Actions
  setConnectionState: (state: Partial<WalletConnectionState>) => void
  setCurrentWallet: (wallet: ConnectedWallet | null) => void
  addConnectedWallet: (wallet: ConnectedWallet) => void
  removeConnectedWallet: (address: Address) => void
  updateWalletBalance: (address: Address, balance: string) => void
  updateWalletChain: (address: Address, chainId: number) => void
  
  // Session management
  startSession: (wallet: ConnectedWallet) => void
  endSession: () => void
  updateSessionActivity: () => void
  
  // Preferences
  updatePreferences: (preferences: Partial<WalletPreferences>) => void
  setAutoConnect: (enabled: boolean) => void
  setPreferredWallet: (walletId: string) => void
  
  // Analytics
  incrementConnections: (walletId: string) => void
  incrementTransactions: (chainId: number) => void
  addTransactionVolume: (amount: string, chainId: number) => void
  
  // Utility
  getWalletById: (id: string) => ConnectedWallet | undefined
  getWalletByAddress: (address: Address) => ConnectedWallet | undefined
  getMostUsedWallet: () => ConnectedWallet | undefined
  getRecentWallets: (limit?: number) => ConnectedWallet[]
  
  // Cleanup
  clearExpiredSessions: () => void
  clearWalletHistory: () => void
  resetAnalytics: () => void
  resetStore: () => void
}

// Default preferences
const defaultPreferences: WalletPreferences = {
  autoConnect: false,
  showTestnets: false,
  hideSmallBalances: true,
  defaultSlippage: 0.5,
  gasPreference: 'standard',
  notifications: {
    transactions: true,
    priceAlerts: false,
    news: false
  }
}

// Default analytics
const defaultAnalytics: WalletAnalytics = {
  totalConnections: 0,
  totalTransactions: 0,
  totalValueTransacted: '0',
  favoriteChains: {},
  walletUsageStats: {},
  monthlyStats: {}
}

// Create the store with persistence
export const useWalletStore = create<WalletState>()(
  subscribeWithSelector(
    persist(
      (set, get) => ({
        // Initial state
        isConnected: false,
        isConnecting: false,
        connectionError: null,
        currentWallet: null,
        currentSession: null,
        connectedWallets: [],
        walletHistory: [],
        sessions: [],
        preferences: defaultPreferences,
        analytics: defaultAnalytics,
        availableWallets: WALLET_PROVIDERS,

        // Connection state actions
        setConnectionState: (state) => set((prev) => ({
          isConnected: state.isConnected ?? prev.isConnected,
          isConnecting: state.isConnecting ?? prev.isConnecting,
          connectionError: state.error ?? prev.connectionError,
        })),

        // Wallet management
        setCurrentWallet: (wallet) => set({ currentWallet: wallet }),

        addConnectedWallet: (wallet) => set((state) => {
          const existing = state.connectedWallets.find(w => w.address === wallet.address)
          if (existing) {
            // Update existing wallet
            return {
              connectedWallets: state.connectedWallets.map(w => 
                w.address === wallet.address 
                  ? { ...w, ...wallet, lastUsed: Date.now() }
                  : w
              )
            }
          }
          
          // Add new wallet
          const newWallet = { ...wallet, lastUsed: Date.now() }
          return {
            connectedWallets: [...state.connectedWallets, newWallet],
            walletHistory: [newWallet, ...state.walletHistory.slice(0, 9)] // Keep last 10
          }
        }),

        removeConnectedWallet: (address) => set((state) => ({
          connectedWallets: state.connectedWallets.filter(w => w.address !== address),
          currentWallet: state.currentWallet?.address === address ? null : state.currentWallet
        })),

        updateWalletBalance: (address, balance) => set((state) => ({
          connectedWallets: state.connectedWallets.map(w =>
            w.address === address ? { ...w, balance } : w
          ),
          currentWallet: state.currentWallet?.address === address 
            ? { ...state.currentWallet, balance }
            : state.currentWallet
        })),

        updateWalletChain: (address, chainId) => set((state) => ({
          connectedWallets: state.connectedWallets.map(w =>
            w.address === address ? { ...w, chainId } : w
          ),
          currentWallet: state.currentWallet?.address === address 
            ? { ...state.currentWallet, chainId }
            : state.currentWallet
        })),

        // Session management
        startSession: (wallet) => {
          const session: WalletSession = {
            id: `session_${Date.now()}_${Math.random().toString(36).substr(2, 9)}`,
            walletId: wallet.id,
            address: wallet.address,
            chainId: wallet.chainId,
            startTime: Date.now(),
            transactionCount: 0,
            lastActivity: Date.now()
          }
          
          set((state) => ({
            currentSession: session,
            sessions: [session, ...state.sessions.slice(0, 49)] // Keep last 50 sessions
          }))
        },

        endSession: () => set((state) => {
          if (!state.currentSession) return state
          
          const endedSession = {
            ...state.currentSession,
            endTime: Date.now(),
            duration: Date.now() - state.currentSession.startTime
          }
          
          return {
            currentSession: null,
            sessions: state.sessions.map(s => 
              s.id === endedSession.id ? endedSession : s
            )
          }
        }),

        updateSessionActivity: () => set((state) => {
          if (!state.currentSession) return state
          
          return {
            currentSession: {
              ...state.currentSession,
              lastActivity: Date.now()
            }
          }
        }),

        // Preferences
        updatePreferences: (newPreferences) => set((state) => ({
          preferences: { ...state.preferences, ...newPreferences }
        })),

        setAutoConnect: (enabled) => set((state) => ({
          preferences: { ...state.preferences, autoConnect: enabled }
        })),

        setPreferredWallet: (walletId) => set((state) => ({
          preferences: { ...state.preferences, preferredWallet: walletId }
        })),

        // Analytics
        incrementConnections: (walletId) => set((state) => {
          const currentMonth = new Date().toISOString().slice(0, 7) // YYYY-MM
          
          return {
            analytics: {
              ...state.analytics,
              totalConnections: state.analytics.totalConnections + 1,
              walletUsageStats: {
                ...state.analytics.walletUsageStats,
                [walletId]: {
                  connections: (state.analytics.walletUsageStats[walletId]?.connections || 0) + 1,
                  lastUsed: Date.now(),
                  totalTime: state.analytics.walletUsageStats[walletId]?.totalTime || 0
                }
              },
              monthlyStats: {
                ...state.analytics.monthlyStats,
                [currentMonth]: {
                  connections: (state.analytics.monthlyStats[currentMonth]?.connections || 0) + 1,
                  transactions: state.analytics.monthlyStats[currentMonth]?.transactions || 0,
                  volume: state.analytics.monthlyStats[currentMonth]?.volume || '0'
                }
              }
            }
          }
        }),

        incrementTransactions: (chainId) => set((state) => {
          const currentMonth = new Date().toISOString().slice(0, 7)
          
          return {
            analytics: {
              ...state.analytics,
              totalTransactions: state.analytics.totalTransactions + 1,
              favoriteChains: {
                ...state.analytics.favoriteChains,
                [chainId]: (state.analytics.favoriteChains[chainId] || 0) + 1
              },
              monthlyStats: {
                ...state.analytics.monthlyStats,
                [currentMonth]: {
                  connections: state.analytics.monthlyStats[currentMonth]?.connections || 0,
                  transactions: (state.analytics.monthlyStats[currentMonth]?.transactions || 0) + 1,
                  volume: state.analytics.monthlyStats[currentMonth]?.volume || '0'
                }
              }
            }
          }
        }),

        addTransactionVolume: (amount, chainId) => set((state) => {
          const currentMonth = new Date().toISOString().slice(0, 7)
          const currentVolume = parseFloat(state.analytics.totalValueTransacted)
          const newVolume = currentVolume + parseFloat(amount)
          
          const monthlyVolume = parseFloat(state.analytics.monthlyStats[currentMonth]?.volume || '0')
          const newMonthlyVolume = monthlyVolume + parseFloat(amount)
          
          return {
            analytics: {
              ...state.analytics,
              totalValueTransacted: newVolume.toString(),
              monthlyStats: {
                ...state.analytics.monthlyStats,
                [currentMonth]: {
                  ...state.analytics.monthlyStats[currentMonth],
                  volume: newMonthlyVolume.toString()
                }
              }
            }
          }
        }),

        // Utility functions
        getWalletById: (id) => {
          return get().connectedWallets.find(w => w.id === id)
        },

        getWalletByAddress: (address) => {
          return get().connectedWallets.find(w => w.address === address)
        },

        getMostUsedWallet: () => {
          const { analytics, connectedWallets } = get()
          const mostUsedWalletId = Object.entries(analytics.walletUsageStats)
            .sort(([,a], [,b]) => b.connections - a.connections)[0]?.[0]
          
          return connectedWallets.find(w => w.id === mostUsedWalletId)
        },

        getRecentWallets: (limit = 5) => {
          return get().walletHistory.slice(0, limit)
        },

        // Cleanup functions
        clearExpiredSessions: () => set((state) => {
          const oneWeekAgo = Date.now() - (7 * 24 * 60 * 60 * 1000)
          return {
            sessions: state.sessions.filter(s => s.startTime > oneWeekAgo)
          }
        }),

        clearWalletHistory: () => set({ walletHistory: [] }),

        resetAnalytics: () => set({ analytics: defaultAnalytics }),

        resetStore: () => set({
          isConnected: false,
          isConnecting: false,
          connectionError: null,
          currentWallet: null,
          currentSession: null,
          connectedWallets: [],
          walletHistory: [],
          sessions: [],
          preferences: defaultPreferences,
          analytics: defaultAnalytics
        })
      }),
      {
        name: 'wallet-store',
        storage: createJSONStorage(() => localStorage),
        partialize: (state) => ({
          // Only persist these fields
          connectedWallets: state.connectedWallets,
          walletHistory: state.walletHistory,
          sessions: state.sessions,
          preferences: state.preferences,
          analytics: state.analytics
        }),
        version: 1,
        migrate: (persistedState: any, version: number) => {
          // Handle migrations between versions
          if (version === 0) {
            // Migration from version 0 to 1
            return {
              ...persistedState,
              analytics: defaultAnalytics
            }
          }
          return persistedState
        }
      }
    )
  )
)
