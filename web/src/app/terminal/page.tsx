'use client'

import { useState } from 'react'
import { motion } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { 
  Terminal as TerminalIcon, 
  Settings, 
  Plus, 
  History,
  BookOpen,
  Zap,
  TrendingUp,
  Wallet,
  Brain
} from 'lucide-react'
import { Terminal } from '@/components/terminal/Terminal'

export default function TerminalPage() {
  const [activeTab, setActiveTab] = useState('terminal')
  const [sessions, setSessions] = useState<string[]>([])
  const [currentSession, setCurrentSession] = useState<string | null>(null)
  const [isFullscreen, setIsFullscreen] = useState(false)

  const handleNewSession = () => {
    // This will trigger the terminal to create a new session
    setCurrentSession(null)
  }

  const handleSessionChange = (sessionId: string) => {
    setCurrentSession(sessionId)
    if (!sessions.includes(sessionId)) {
      setSessions(prev => [...prev, sessionId])
    }
  }

  const commandCategories = [
    {
      name: 'System',
      icon: Settings,
      commands: ['status', 'help', 'clear', 'exit', 'config', 'logs', 'health', 'version'],
      description: 'System management and information commands'
    },
    {
      name: 'Trading',
      icon: TrendingUp,
      commands: ['buy', 'sell', 'portfolio', 'orders', 'history', 'balance', 'positions'],
      description: 'Trading operations and portfolio management'
    },
    {
      name: 'Market',
      icon: Zap,
      commands: ['price', 'chart', 'news', 'analysis', 'alerts'],
      description: 'Market data and analysis commands'
    },
    {
      name: 'AI',
      icon: Brain,
      commands: ['analyze', 'predict', 'sentiment', 'chat', 'learn'],
      description: 'AI-powered analysis and predictions'
    },
    {
      name: 'Web3',
      icon: Wallet,
      commands: ['wallet', 'connect', 'balance', 'transfer', 'defi'],
      description: 'Blockchain and Web3 operations'
    }
  ]

  if (isFullscreen) {
    return (
      <Terminal
        sessionId={currentSession || undefined}
        onSessionChange={handleSessionChange}
        fullscreen={true}
        onFullscreenChange={setIsFullscreen}
        className="h-screen"
      />
    )
  }

  return (
    <div className="container mx-auto p-6 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h1 className="text-3xl font-bold flex items-center gap-2">
            <TerminalIcon className="w-8 h-8" />
            Terminal
          </h1>
          <p className="text-muted-foreground">
            Command-line interface for AI-powered crypto trading
          </p>
        </div>
        
        <div className="flex items-center gap-2">
          <Button onClick={handleNewSession} className="flex items-center gap-2">
            <Plus className="w-4 h-4" />
            New Session
          </Button>
        </div>
      </div>

      {/* Session Info */}
      {currentSession && (
        <Card>
          <CardHeader>
            <CardTitle className="text-lg">Active Session</CardTitle>
          </CardHeader>
          <CardContent>
            <div className="flex items-center gap-4">
              <Badge variant="outline">
                Session: {currentSession.substring(0, 8)}...
              </Badge>
              <Badge variant="default">
                Connected
              </Badge>
            </div>
          </CardContent>
        </Card>
      )}

      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-3">
          <TabsTrigger value="terminal" className="flex items-center gap-2">
            <TerminalIcon className="w-4 h-4" />
            Terminal
          </TabsTrigger>
          <TabsTrigger value="commands" className="flex items-center gap-2">
            <BookOpen className="w-4 h-4" />
            Commands
          </TabsTrigger>
          <TabsTrigger value="history" className="flex items-center gap-2">
            <History className="w-4 h-4" />
            History
          </TabsTrigger>
        </TabsList>

        <TabsContent value="terminal" className="space-y-4">
          <Terminal
            sessionId={currentSession || undefined}
            onSessionChange={handleSessionChange}
            fullscreen={isFullscreen}
            onFullscreenChange={setIsFullscreen}
          />
        </TabsContent>

        <TabsContent value="commands" className="space-y-6">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {commandCategories.map((category) => {
              const IconComponent = category.icon
              return (
                <motion.div
                  key={category.name}
                  initial={{ opacity: 0, y: 20 }}
                  animate={{ opacity: 1, y: 0 }}
                  transition={{ duration: 0.3 }}
                >
                  <Card>
                    <CardHeader>
                      <CardTitle className="flex items-center gap-2 text-lg">
                        <IconComponent className="w-5 h-5" />
                        {category.name}
                      </CardTitle>
                      <CardDescription>
                        {category.description}
                      </CardDescription>
                    </CardHeader>
                    <CardContent>
                      <div className="space-y-2">
                        {category.commands.map((command) => (
                          <div
                            key={command}
                            className="flex items-center justify-between p-2 rounded bg-muted/50 hover:bg-muted cursor-pointer transition-colors"
                            onClick={() => {
                              // Switch to terminal tab and focus
                              setActiveTab('terminal')
                              // TODO: Insert command into terminal input
                            }}
                          >
                            <code className="text-sm font-mono">{command}</code>
                            <Button variant="ghost" size="sm">
                              Try
                            </Button>
                          </div>
                        ))}
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              )
            })}
          </div>

          {/* Quick Start Guide */}
          <Card>
            <CardHeader>
              <CardTitle>Quick Start Guide</CardTitle>
              <CardDescription>
                Get started with the terminal interface
              </CardDescription>
            </CardHeader>
            <CardContent className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div className="space-y-2">
                  <h4 className="font-semibold">Basic Commands</h4>
                  <div className="space-y-1 text-sm">
                    <div><code>help</code> - Show available commands</div>
                    <div><code>status</code> - Check system status</div>
                    <div><code>clear</code> - Clear terminal screen</div>
                    <div><code>exit</code> - Exit current session</div>
                  </div>
                </div>
                
                <div className="space-y-2">
                  <h4 className="font-semibold">Trading Examples</h4>
                  <div className="space-y-1 text-sm">
                    <div><code>price BTC</code> - Get Bitcoin price</div>
                    <div><code>portfolio</code> - View portfolio</div>
                    <div><code>buy BTC 0.1</code> - Buy 0.1 Bitcoin</div>
                    <div><code>analyze ETH</code> - AI analysis of Ethereum</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="history" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle>Session History</CardTitle>
              <CardDescription>
                View your recent terminal sessions and commands
              </CardDescription>
            </CardHeader>
            <CardContent>
              {sessions.length === 0 ? (
                <div className="text-center py-8 text-muted-foreground">
                  No sessions yet. Start by creating a new terminal session.
                </div>
              ) : (
                <div className="space-y-4">
                  {sessions.map((sessionId) => (
                    <div
                      key={sessionId}
                      className="flex items-center justify-between p-3 border rounded-lg hover:bg-muted/50 cursor-pointer"
                      onClick={() => setCurrentSession(sessionId)}
                    >
                      <div>
                        <div className="font-mono text-sm">
                          Session: {sessionId.substring(0, 16)}...
                        </div>
                        <div className="text-xs text-muted-foreground">
                          Created: {new Date().toLocaleString()}
                        </div>
                      </div>
                      <div className="flex items-center gap-2">
                        {sessionId === currentSession && (
                          <Badge variant="default" className="text-xs">
                            Active
                          </Badge>
                        )}
                        <Button variant="outline" size="sm">
                          Load
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>
              )}
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
