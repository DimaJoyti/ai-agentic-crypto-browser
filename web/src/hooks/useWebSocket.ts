import { useState, useEffect, useRef, useCallback } from 'react'

interface WebSocketMessage {
  type: string
  data: any
  timestamp: string
}

interface UseWebSocketReturn {
  isConnected: boolean
  lastMessage: WebSocketMessage | null
  sendMessage: (message: any) => void
  connect: () => void
  disconnect: () => void
}

export function useWebSocket(url?: string): UseWebSocketReturn {
  const [isConnected, setIsConnected] = useState(false)
  const [lastMessage, setLastMessage] = useState<WebSocketMessage | null>(null)
  const ws = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout>()
  const reconnectAttempts = useRef(0)
  const maxReconnectAttempts = 5

  const wsUrl = url || `${process.env.NEXT_PUBLIC_WS_URL}/ws`

  const connect = useCallback(() => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      return
    }

    try {
      ws.current = new WebSocket(wsUrl)

      ws.current.onopen = () => {
        console.log('WebSocket connected')
        setIsConnected(true)
        reconnectAttempts.current = 0
      }

      ws.current.onmessage = (event) => {
        try {
          const message: WebSocketMessage = JSON.parse(event.data)
          setLastMessage(message)
        } catch (error) {
          console.error('Failed to parse WebSocket message:', error)
        }
      }

      ws.current.onclose = (event) => {
        console.log('WebSocket disconnected:', event.code, event.reason)
        setIsConnected(false)

        // Attempt to reconnect if not a manual close
        if (event.code !== 1000 && reconnectAttempts.current < maxReconnectAttempts) {
          const timeout = Math.pow(2, reconnectAttempts.current) * 1000 // Exponential backoff
          reconnectTimeoutRef.current = setTimeout(() => {
            reconnectAttempts.current++
            console.log(`Attempting to reconnect (${reconnectAttempts.current}/${maxReconnectAttempts})...`)
            connect()
          }, timeout)
        }
      }

      ws.current.onerror = (error) => {
        console.error('WebSocket error:', error)
      }
    } catch (error) {
      console.error('Failed to create WebSocket connection:', error)
    }
  }, [wsUrl])

  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
    }
    
    if (ws.current) {
      ws.current.close(1000, 'Manual disconnect')
      ws.current = null
    }
    
    setIsConnected(false)
    reconnectAttempts.current = 0
  }, [])

  const sendMessage = useCallback((message: any) => {
    if (ws.current?.readyState === WebSocket.OPEN) {
      try {
        const messageWithTimestamp = {
          ...message,
          timestamp: new Date().toISOString(),
        }
        ws.current.send(JSON.stringify(messageWithTimestamp))
      } catch (error) {
        console.error('Failed to send WebSocket message:', error)
      }
    } else {
      console.warn('WebSocket is not connected')
    }
  }, [])

  useEffect(() => {
    connect()

    return () => {
      disconnect()
    }
  }, [connect, disconnect])

  // Cleanup on unmount
  useEffect(() => {
    return () => {
      if (reconnectTimeoutRef.current) {
        clearTimeout(reconnectTimeoutRef.current)
      }
    }
  }, [])

  return {
    isConnected,
    lastMessage,
    sendMessage,
    connect,
    disconnect,
  }
}

// Hook for specific message types
export function useWebSocketSubscription(
  messageType: string,
  callback: (data: any) => void,
  deps: any[] = []
) {
  const { lastMessage } = useWebSocket()

  useEffect(() => {
    if (lastMessage && lastMessage.type === messageType) {
      callback(lastMessage.data)
    }
  }, [lastMessage, messageType, callback, ...deps])
}

// Hook for real-time notifications
export function useRealTimeNotifications() {
  const [notifications, setNotifications] = useState<any[]>([])

  useWebSocketSubscription('notification', (data) => {
    setNotifications(prev => [data, ...prev.slice(0, 9)]) // Keep last 10 notifications
  })

  const clearNotifications = useCallback(() => {
    setNotifications([])
  }, [])

  const removeNotification = useCallback((id: string) => {
    setNotifications(prev => prev.filter(n => n.id !== id))
  }, [])

  return {
    notifications,
    clearNotifications,
    removeNotification,
  }
}

// Hook for real-time task updates
export function useTaskUpdates() {
  const [tasks, setTasks] = useState<any[]>([])

  useWebSocketSubscription('task_update', (data) => {
    setTasks(prev => {
      const existingIndex = prev.findIndex(t => t.id === data.id)
      if (existingIndex >= 0) {
        // Update existing task
        const updated = [...prev]
        updated[existingIndex] = { ...updated[existingIndex], ...data }
        return updated
      } else {
        // Add new task
        return [data, ...prev]
      }
    })
  })

  useWebSocketSubscription('task_completed', (data) => {
    setTasks(prev => 
      prev.map(task => 
        task.id === data.id 
          ? { ...task, status: 'completed', completedAt: data.completedAt }
          : task
      )
    )
  })

  useWebSocketSubscription('task_failed', (data) => {
    setTasks(prev => 
      prev.map(task => 
        task.id === data.id 
          ? { ...task, status: 'failed', error: data.error }
          : task
      )
    )
  })

  return { tasks }
}

// Hook for real-time portfolio updates
export function usePortfolioUpdates() {
  const [portfolio, setPortfolio] = useState<any>(null)

  useWebSocketSubscription('portfolio_update', (data) => {
    setPortfolio(data)
  })

  useWebSocketSubscription('price_update', (data) => {
    setPortfolio((prev: any) => {
      if (!prev) return prev
      
      return {
        ...prev,
        tokens: prev.tokens?.map((token: any) => 
          token.symbol === data.symbol 
            ? { ...token, price: data.price, change24h: data.change24h }
            : token
        )
      }
    })
  })

  return { portfolio }
}
