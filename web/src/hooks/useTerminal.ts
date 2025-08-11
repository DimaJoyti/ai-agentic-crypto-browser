'use client'

import { useState, useEffect, useCallback, useRef } from 'react'

interface TerminalSession {
  id: string
  userId: string
  createdAt: string
  lastActive: string
  environment: Record<string, string>
}

interface CommandHistoryEntry {
  id: string
  command: string
  args: string[]
  output: string
  error?: string
  exitCode: number
  startTime: string
  endTime: string
  duration: number
}

interface WSMessage {
  type: string
  sessionId?: string
  data: any
  timestamp: string
}

interface UseTerminalOptions {
  sessionId?: string
  onSessionChange?: (sessionId: string) => void
  autoConnect?: boolean
  reconnectAttempts?: number
  reconnectDelay?: number
}

export function useTerminal({
  sessionId,
  onSessionChange,
  autoConnect = true,
  reconnectAttempts = 3,
  reconnectDelay = 1000
}: UseTerminalOptions = {}) {
  const [terminal, setTerminal] = useState<any>(null)
  const [isConnected, setIsConnected] = useState(false)
  const [currentSession, setCurrentSession] = useState<string | null>(sessionId || null)
  const [commandHistory, setCommandHistory] = useState<CommandHistoryEntry[]>([])
  const [error, setError] = useState<string | null>(null)
  const [isLoading, setIsLoading] = useState(false)
  
  const wsRef = useRef<WebSocket | null>(null)
  const reconnectTimeoutRef = useRef<NodeJS.Timeout | null>(null)
  const reconnectCountRef = useRef(0)

  // WebSocket connection
  const connect = useCallback(async () => {
    if (wsRef.current?.readyState === WebSocket.OPEN) {
      return
    }

    setIsLoading(true)
    setError(null)

    try {
      // Create WebSocket connection
      const wsUrl = `ws://localhost:8085/ws`
      const ws = new WebSocket(wsUrl)

      ws.onopen = () => {
        console.log('Terminal WebSocket connected')
        setIsConnected(true)
        setIsLoading(false)
        setError(null)
        reconnectCountRef.current = 0

        // Join or create session
        if (currentSession) {
          ws.send(JSON.stringify({
            type: 'session_join',
            data: { session_id: currentSession },
            timestamp: new Date().toISOString()
          }))
        } else {
          // Create new session
          createSession()
        }
      }

      ws.onmessage = (event) => {
        try {
          const message: WSMessage = JSON.parse(event.data)
          handleWebSocketMessage(message)
        } catch (err) {
          console.error('Failed to parse WebSocket message:', err)
        }
      }

      ws.onclose = (event) => {
        console.log('Terminal WebSocket disconnected:', event.code, event.reason)
        setIsConnected(false)
        wsRef.current = null

        // Attempt reconnection if not a clean close
        if (event.code !== 1000 && reconnectCountRef.current < reconnectAttempts) {
          reconnectCountRef.current++
          reconnectTimeoutRef.current = setTimeout(() => {
            console.log(`Reconnecting... (${reconnectCountRef.current}/${reconnectAttempts})`)
            connect()
          }, reconnectDelay * reconnectCountRef.current)
        }
      }

      ws.onerror = (error) => {
        console.error('Terminal WebSocket error:', error)
        setError('Connection error')
        setIsLoading(false)
      }

      wsRef.current = ws
    } catch (err) {
      console.error('Failed to connect to terminal:', err)
      setError('Failed to connect')
      setIsLoading(false)
    }
  }, [currentSession, reconnectAttempts, reconnectDelay])

  // Disconnect WebSocket
  const disconnect = useCallback(() => {
    if (reconnectTimeoutRef.current) {
      clearTimeout(reconnectTimeoutRef.current)
      reconnectTimeoutRef.current = null
    }

    if (wsRef.current) {
      wsRef.current.close(1000, 'User disconnected')
      wsRef.current = null
    }

    setIsConnected(false)
    reconnectCountRef.current = 0
  }, [])

  // Handle WebSocket messages
  const handleWebSocketMessage = useCallback((message: WSMessage) => {
    switch (message.type) {
      case 'welcome':
        console.log('Terminal service welcome:', message.data)
        break

      case 'command_output':
        // Add command output to history
        const outputData = message.data
        setCommandHistory(prev => [
          ...prev,
          {
            id: Date.now().toString(),
            command: outputData.command || '',
            args: outputData.args || [],
            output: outputData.output || '',
            error: outputData.error,
            exitCode: outputData.exit_code || 0,
            startTime: outputData.start_time || new Date().toISOString(),
            endTime: outputData.end_time || new Date().toISOString(),
            duration: outputData.duration || 0
          }
        ])
        break

      case 'session_created':
        const sessionData = message.data
        setCurrentSession(sessionData.session_id)
        onSessionChange?.(sessionData.session_id)
        break

      case 'session_joined':
        console.log('Joined session:', message.data)
        break

      case 'error':
        setError(message.data.message || 'Unknown error')
        break

      default:
        console.log('Unknown message type:', message.type, message.data)
    }
  }, [onSessionChange])

  // Create new session
  const createSession = useCallback(async () => {
    try {
      const response = await fetch('http://localhost:8085/api/v1/sessions', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({
          user_id: 'anonymous', // TODO: Get from auth context
          environment: {}
        })
      })

      if (!response.ok) {
        throw new Error('Failed to create session')
      }

      const data = await response.json()
      setCurrentSession(data.session_id)
      onSessionChange?.(data.session_id)

      // Join the session via WebSocket
      if (wsRef.current?.readyState === WebSocket.OPEN) {
        wsRef.current.send(JSON.stringify({
          type: 'session_join',
          data: { session_id: data.session_id },
          timestamp: new Date().toISOString()
        }))
      }
    } catch (err) {
      console.error('Failed to create session:', err)
      setError('Failed to create session')
    }
  }, [onSessionChange])

  // Send command
  const sendCommand = useCallback((command: string) => {
    if (!wsRef.current || wsRef.current.readyState !== WebSocket.OPEN) {
      setError('Not connected to terminal')
      return
    }

    if (!currentSession) {
      setError('No active session')
      return
    }

    // Add command to history immediately (optimistic update)
    const tempEntry: CommandHistoryEntry = {
      id: Date.now().toString(),
      command,
      args: command.split(' ').slice(1),
      output: '',
      exitCode: 0,
      startTime: new Date().toISOString(),
      endTime: new Date().toISOString(),
      duration: 0
    }

    setCommandHistory(prev => [...prev, tempEntry])

    // Send command via WebSocket
    wsRef.current.send(JSON.stringify({
      type: 'command',
      data: {
        command,
        session_id: currentSession
      },
      timestamp: new Date().toISOString()
    }))
  }, [currentSession])

  // Clear terminal
  const clearTerminal = useCallback(() => {
    setCommandHistory([])
  }, [])

  // Get session history
  const getSessionHistory = useCallback(async (sessionId: string) => {
    try {
      const response = await fetch(`http://localhost:8085/api/v1/sessions/${sessionId}/history`)
      if (!response.ok) {
        throw new Error('Failed to get session history')
      }
      const history = await response.json()
      setCommandHistory(history)
    } catch (err) {
      console.error('Failed to get session history:', err)
      setError('Failed to load history')
    }
  }, [])

  // Auto-connect on mount
  useEffect(() => {
    if (autoConnect) {
      connect()
    }

    return () => {
      disconnect()
    }
  }, [autoConnect, connect, disconnect])

  // Load session history when session changes
  useEffect(() => {
    if (currentSession && isConnected) {
      getSessionHistory(currentSession)
    }
  }, [currentSession, isConnected, getSessionHistory])

  return {
    terminal,
    isConnected,
    isLoading,
    currentSession,
    commandHistory,
    error,
    connect,
    disconnect,
    sendCommand,
    createSession,
    clearTerminal,
    getSessionHistory
  }
}
