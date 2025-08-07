'use client'

import React, { Component, ErrorInfo, ReactNode } from 'react'
import { chunkErrorHandler } from '@/lib/chunk-error-handler'

interface Props {
  children: ReactNode
}

interface State {
  hasError: boolean
  error?: Error
  errorInfo?: ErrorInfo
}

export class ChunkErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props)
    this.state = { hasError: false }
  }

  static getDerivedStateFromError(error: Error): State {
    // Check if this is a chunk loading error or WalletConnect SSR error
    const isChunkError = error.message?.includes('Loading chunk') ||
                        error.message?.includes('timeout') ||
                        error.name === 'ChunkLoadError'

    const isWalletConnectSSRError = error.message?.includes('indexedDB is not defined') ||
                                   error.message?.includes('localStorage is not defined') ||
                                   error.message?.includes('window is not defined')

    return {
      hasError: true,
      error: (isChunkError || isWalletConnectSSRError) ? error : undefined
    }
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error('ChunkErrorBoundary caught an error:', error, errorInfo)

    // Check if this is a chunk loading error
    const isChunkError = error.message?.includes('Loading chunk') ||
                        error.message?.includes('timeout') ||
                        error.name === 'ChunkLoadError'

    // Check if this is a WalletConnect SSR error
    const isWalletConnectSSRError = error.message?.includes('indexedDB is not defined') ||
                                   error.message?.includes('localStorage is not defined') ||
                                   error.message?.includes('window is not defined')

    if (isChunkError) {
      console.log('Chunk loading error detected in error boundary')
      // Let the chunk error handler deal with it
      this.handleChunkError()
    } else if (isWalletConnectSSRError) {
      console.log('WalletConnect SSR error detected, reloading page')
      // For SSR errors, reload the page to ensure client-side rendering
      this.handleChunkError()
    } else {
      // For non-chunk errors, update state to show error UI
      this.setState({
        hasError: true,
        error,
        errorInfo
      })
    }
  }

  private handleChunkError = () => {
    // Reset the error boundary state
    this.setState({ hasError: false })
    
    // Trigger a page reload after a short delay
    setTimeout(() => {
      window.location.reload()
    }, 1000)
  }

  private handleRetry = () => {
    // Reset error state and try again
    this.setState({ hasError: false })
    chunkErrorHandler.reset()
  }

  private handleReload = () => {
    window.location.reload()
  }

  render() {
    if (this.state.hasError) {
      // Check if it's a chunk error or WalletConnect SSR error
      const isChunkOrSSRError = this.state.error?.message?.includes('Loading chunk') ||
                               this.state.error?.message?.includes('timeout') ||
                               this.state.error?.message?.includes('indexedDB is not defined') ||
                               this.state.error?.message?.includes('localStorage is not defined') ||
                               this.state.error?.message?.includes('window is not defined')

      if (isChunkOrSSRError) {
        return (
          <div className="min-h-screen flex items-center justify-center bg-gray-50">
            <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6 text-center">
              <div className="mb-4">
                <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-yellow-100">
                  <svg className="h-6 w-6 text-yellow-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 9v2m0 4h.01m-6.938 4h13.856c1.54 0 2.502-1.667 1.732-2.5L13.732 4c-.77-.833-1.964-.833-2.732 0L3.732 16.5c-.77.833.192 2.5 1.732 2.5z" />
                  </svg>
                </div>
              </div>
              <h3 className="text-lg font-medium text-gray-900 mb-2">
                Loading Issue
              </h3>
              <p className="text-sm text-gray-500 mb-6">
                There was a problem loading some application resources. This usually resolves with a page refresh.
              </p>
              <div className="flex space-x-3">
                <button
                  onClick={this.handleRetry}
                  className="flex-1 bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
                >
                  Try Again
                </button>
                <button
                  onClick={this.handleReload}
                  className="flex-1 bg-gray-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
                >
                  Refresh Page
                </button>
              </div>
            </div>
          </div>
        )
      }

      // For other errors, show a generic error boundary
      return (
        <div className="min-h-screen flex items-center justify-center bg-gray-50">
          <div className="max-w-md w-full bg-white shadow-lg rounded-lg p-6 text-center">
            <div className="mb-4">
              <div className="mx-auto flex items-center justify-center h-12 w-12 rounded-full bg-red-100">
                <svg className="h-6 w-6 text-red-600" fill="none" viewBox="0 0 24 24" stroke="currentColor">
                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8v4m0 4h.01M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                </svg>
              </div>
            </div>
            <h3 className="text-lg font-medium text-gray-900 mb-2">
              Something went wrong
            </h3>
            <p className="text-sm text-gray-500 mb-6">
              An unexpected error occurred. Please try refreshing the page.
            </p>
            <div className="space-y-3">
              <button
                onClick={this.handleRetry}
                className="w-full bg-blue-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-blue-700 focus:outline-none focus:ring-2 focus:ring-blue-500"
              >
                Try Again
              </button>
              <button
                onClick={this.handleReload}
                className="w-full bg-gray-600 text-white px-4 py-2 rounded-md text-sm font-medium hover:bg-gray-700 focus:outline-none focus:ring-2 focus:ring-gray-500"
              >
                Refresh Page
              </button>
              {process.env.NODE_ENV === 'development' && this.state.error && (
                <details className="mt-4 text-left">
                  <summary className="text-sm text-gray-600 cursor-pointer">
                    Error Details (Development)
                  </summary>
                  <pre className="mt-2 text-xs text-gray-800 bg-gray-100 p-2 rounded overflow-auto max-h-32">
                    {this.state.error.toString()}
                    {this.state.errorInfo?.componentStack}
                  </pre>
                </details>
              )}
            </div>
          </div>
        </div>
      )
    }

    return this.props.children
  }
}

// Client-side initialization component
export function ChunkErrorInitializer() {
  React.useEffect(() => {
    // Initialize the chunk error handler
    chunkErrorHandler.reset()
  }, [])

  return null
}
