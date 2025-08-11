'use client'

import React, { useEffect, useRef, useState, useCallback } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Terminal as TerminalIcon,
  Maximize2,
  Minimize2,
  X,
  Copy,
  ClipboardPaste,
  RotateCcw
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useTerminal } from '@/hooks/useTerminal'

interface TerminalProps {
  className?: string
  sessionId?: string
  onSessionChange?: (sessionId: string) => void
  theme?: 'dark' | 'light' | 'matrix' | 'retro'
  fullscreen?: boolean
  onFullscreenChange?: (fullscreen: boolean) => void
}

export function Terminal({ 
  className, 
  sessionId, 
  onSessionChange,
  theme = 'dark',
  fullscreen = false,
  onFullscreenChange
}: TerminalProps) {
  const terminalRef = useRef<HTMLDivElement>(null)
  const [isMaximized, setIsMaximized] = useState(fullscreen)
  const [currentTheme] = useState(theme)

  const {
    terminal,
    isConnected,
    currentSession,
    commandHistory,
    sendCommand,
    clearTerminal,
    error
  } = useTerminal({
    sessionId,
    onSessionChange,
    autoConnect: true
  })

  // Initialize terminal when component mounts
  useEffect(() => {
    if (terminalRef.current && !terminal) {
      // Initialize terminal (we'll implement this with a simple div for now)
      // In a real implementation, this would use xterm.js
      initializeTerminal()
    }
  }, [terminalRef.current])

  // Handle fullscreen changes
  useEffect(() => {
    setIsMaximized(fullscreen)
  }, [fullscreen])

  const initializeTerminal = useCallback(() => {
    // This is a simplified implementation
    // In production, you would use xterm.js here
    if (terminalRef.current) {
      terminalRef.current.focus()
    }
  }, [])

  const handleMaximize = () => {
    const newMaximized = !isMaximized
    setIsMaximized(newMaximized)
    onFullscreenChange?.(newMaximized)
  }

  const handleClear = () => {
    clearTerminal()
  }

  const handleCopy = () => {
    // Copy selected text or last output
    const selection = window.getSelection()?.toString()
    if (selection) {
      navigator.clipboard.writeText(selection)
    }
  }

  const handlePaste = async () => {
    try {
      const text = await navigator.clipboard.readText()
      if (text) {
        sendCommand(text)
      }
    } catch (err) {
      console.error('Failed to paste:', err)
    }
  }

  const getThemeClasses = () => {
    switch (currentTheme) {
      case 'light':
        return 'bg-white text-black'
      case 'matrix':
        return 'bg-black text-green-400 font-mono'
      case 'retro':
        return 'bg-blue-900 text-amber-300 font-mono'
      default:
        return 'bg-gray-900 text-green-400 font-mono'
    }
  }

  return (
    <Card className={cn(
      'terminal-container',
      isMaximized && 'fixed inset-0 z-50 rounded-none',
      className
    )}>
      <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
        <CardTitle className="flex items-center gap-2">
          <TerminalIcon className="w-4 h-4" />
          Terminal
          {currentSession && (
            <Badge variant="outline" className="text-xs">
              {currentSession.substring(0, 8)}
            </Badge>
          )}
        </CardTitle>
        
        <div className="flex items-center gap-1">
          {/* Connection Status */}
          <Badge variant={isConnected ? 'default' : 'destructive'} className="text-xs">
            {isConnected ? 'Connected' : 'Disconnected'}
          </Badge>
          
          {/* Terminal Controls */}
          <Button variant="ghost" size="sm" onClick={handleCopy} title="Copy">
            <Copy className="w-3 h-3" />
          </Button>
          
          <Button variant="ghost" size="sm" onClick={handlePaste} title="Paste">
            <ClipboardPaste className="w-3 h-3" />
          </Button>
          
          <Button variant="ghost" size="sm" onClick={handleClear} title="Clear">
            <RotateCcw className="w-3 h-3" />
          </Button>
          
          <Button variant="ghost" size="sm" onClick={handleMaximize} title="Maximize">
            {isMaximized ? <Minimize2 className="w-3 h-3" /> : <Maximize2 className="w-3 h-3" />}
          </Button>
          
          {isMaximized && (
            <Button variant="ghost" size="sm" onClick={() => setIsMaximized(false)} title="Close">
              <X className="w-3 h-3" />
            </Button>
          )}
        </div>
      </CardHeader>
      
      <CardContent className="p-0">
        {error && (
          <div className="p-4 bg-red-50 border-b border-red-200 text-red-700 text-sm">
            Error: {error}
          </div>
        )}
        
        <div 
          ref={terminalRef}
          className={cn(
            'terminal-content p-4 min-h-[400px] max-h-[600px] overflow-auto focus:outline-none',
            isMaximized && 'min-h-[calc(100vh-120px)] max-h-[calc(100vh-120px)]',
            getThemeClasses()
          )}
          tabIndex={0}
        >
          {/* Terminal Output Area */}
          <TerminalOutput
            history={commandHistory}
            isConnected={isConnected}
          />

          {/* Command Input */}
          <TerminalInput
            onCommand={sendCommand}
            disabled={!isConnected}
          />
        </div>
      </CardContent>
    </Card>
  )
}

// Terminal Output Component
interface TerminalOutputProps {
  history: any[]
  isConnected: boolean
}

function TerminalOutput({ history, isConnected }: TerminalOutputProps) {
  return (
    <div className="terminal-output space-y-1">
      {/* Welcome Message */}
      {history.length === 0 && (
        <div className="text-sm opacity-70">
          <div>AI-Agentic Crypto Browser Terminal v1.0</div>
          <div>Type 'help' for available commands.</div>
          <div className="mt-2">
            {isConnected ? '✅ Connected to terminal service' : '❌ Disconnected'}
          </div>
        </div>
      )}
      
      {/* Command History */}
      {history.map((entry, index) => (
        <div key={index} className="terminal-entry">
          {/* Command */}
          <div className="flex items-center gap-2">
            <span className="text-blue-400">$</span>
            <span>{entry.command}</span>
          </div>
          
          {/* Output */}
          {entry.output && (
            <div className="ml-4 whitespace-pre-wrap text-sm">
              {entry.output}
            </div>
          )}
          
          {/* Error */}
          {entry.error && (
            <div className="ml-4 text-red-400 text-sm">
              Error: {entry.error}
            </div>
          )}
        </div>
      ))}
    </div>
  )
}

// Terminal Input Component
interface TerminalInputProps {
  onCommand: (command: string) => void
  disabled: boolean
}

function TerminalInput({ onCommand, disabled }: TerminalInputProps) {
  const [input, setInput] = useState('')
  const [historyIndex, setHistoryIndex] = useState(-1)
  const [commandHistory, setCommandHistory] = useState<string[]>([])
  const inputRef = useRef<HTMLInputElement>(null)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    if (input.trim() && !disabled) {
      onCommand(input.trim())
      setCommandHistory(prev => [...prev, input.trim()])
      setInput('')
      setHistoryIndex(-1)
    }
  }

  const handleKeyDown = (e: React.KeyboardEvent) => {
    if (e.key === 'ArrowUp') {
      e.preventDefault()
      if (historyIndex < commandHistory.length - 1) {
        const newIndex = historyIndex + 1
        setHistoryIndex(newIndex)
        setInput(commandHistory[commandHistory.length - 1 - newIndex])
      }
    } else if (e.key === 'ArrowDown') {
      e.preventDefault()
      if (historyIndex > 0) {
        const newIndex = historyIndex - 1
        setHistoryIndex(newIndex)
        setInput(commandHistory[commandHistory.length - 1 - newIndex])
      } else if (historyIndex === 0) {
        setHistoryIndex(-1)
        setInput('')
      }
    }
  }

  // Auto-focus input
  useEffect(() => {
    if (inputRef.current && !disabled) {
      inputRef.current.focus()
    }
  }, [disabled])

  return (
    <form onSubmit={handleSubmit} className="terminal-input flex items-center gap-2 mt-2">
      <span className="text-blue-400">$</span>
      <input
        ref={inputRef}
        type="text"
        value={input}
        onChange={(e) => setInput(e.target.value)}
        onKeyDown={handleKeyDown}
        disabled={disabled}
        className="flex-1 bg-transparent border-none outline-none text-inherit font-inherit"
        placeholder={disabled ? "Connecting..." : "Enter command..."}
        autoComplete="off"
        spellCheck={false}
      />
    </form>
  )
}
