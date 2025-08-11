'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Switch } from '@/components/ui/switch'
import { Badge } from '@/components/ui/badge'
import { Slider } from '@/components/ui/slider'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Brain,
  Cpu,
  Zap,
  Settings,
  Play,
  Pause,
  Square,
  BarChart3,
  TrendingUp,
  TrendingDown,
  Target,
  Shield,
  Clock,
  Activity,
  Layers,
  Grid3X3,
  LineChart,
  PieChart,
  Calculator,
  Gauge,
  Bot,
  Sparkles,
  Rocket,
  Eye,
  Edit,
  Trash2,
  Copy,
  Download,
  Upload,
  RefreshCw,
  CheckCircle,
  XCircle,
  Plus,
  Minus,
  ArrowRight,
  ArrowDown,
  Code,
  Database,
  Wifi,
  WifiOff
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface StrategyNode {
  id: string
  type: 'indicator' | 'condition' | 'action' | 'logic'
  name: string
  parameters: Record<string, any>
  position: { x: number; y: number }
  connections: string[]
}

interface BacktestResult {
  period: string
  totalReturn: number
  annualizedReturn: number
  sharpeRatio: number
  maxDrawdown: number
  winRate: number
  totalTrades: number
  avgTrade: number
  profitFactor: number
  calmarRatio: number
  sortinoRatio: number
  volatility: number
  beta: number
  alpha: number
  var95: number
  expectedShortfall: number
}

interface StrategyTemplate {
  id: string
  name: string
  description: string
  category: 'trend_following' | 'mean_reversion' | 'momentum' | 'arbitrage' | 'market_making'
  difficulty: 'beginner' | 'intermediate' | 'advanced' | 'expert'
  nodes: StrategyNode[]
  expectedReturn: string
  riskLevel: 'low' | 'medium' | 'high'
  timeframe: string
  marketConditions: string[]
}

export function AlgorithmicStrategyBuilder() {
  const [strategyNodes, setStrategyNodes] = useState<StrategyNode[]>([])
  const [selectedNode, setSelectedNode] = useState<string | null>(null)
  const [backtestResults, setBacktestResults] = useState<BacktestResult | null>(null)
  const [isBacktesting, setIsBacktesting] = useState(false)
  const [strategyTemplates, setStrategyTemplates] = useState<StrategyTemplate[]>([])
  const [activeTab, setActiveTab] = useState('builder')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock strategy templates
    const mockTemplates: StrategyTemplate[] = [
      {
        id: 'template1',
        name: 'RSI Mean Reversion',
        description: 'Buy when RSI < 30, sell when RSI > 70 with stop loss and take profit',
        category: 'mean_reversion',
        difficulty: 'beginner',
        nodes: [],
        expectedReturn: '12-18%',
        riskLevel: 'medium',
        timeframe: '1H - 4H',
        marketConditions: ['Sideways', 'Low Volatility']
      },
      {
        id: 'template2',
        name: 'MACD Momentum',
        description: 'Trend following strategy using MACD crossovers with dynamic position sizing',
        category: 'momentum',
        difficulty: 'intermediate',
        nodes: [],
        expectedReturn: '15-25%',
        riskLevel: 'medium',
        timeframe: '4H - 1D',
        marketConditions: ['Trending', 'Medium Volatility']
      },
      {
        id: 'template3',
        name: 'Grid Trading Pro',
        description: 'Advanced grid trading with dynamic range adjustment and ML optimization',
        category: 'market_making',
        difficulty: 'advanced',
        nodes: [],
        expectedReturn: '8-15%',
        riskLevel: 'low',
        timeframe: '15M - 1H',
        marketConditions: ['Sideways', 'High Liquidity']
      },
      {
        id: 'template4',
        name: 'Statistical Arbitrage',
        description: 'Pairs trading strategy using cointegration and statistical models',
        category: 'arbitrage',
        difficulty: 'expert',
        nodes: [],
        expectedReturn: '20-35%',
        riskLevel: 'high',
        timeframe: '1M - 15M',
        marketConditions: ['Any', 'High Frequency']
      }
    ]

    setStrategyTemplates(mockTemplates)
  }, [isConnected])

  const addNode = (type: StrategyNode['type'], name: string) => {
    const newNode: StrategyNode = {
      id: `node_${Date.now()}`,
      type,
      name,
      parameters: {},
      position: { x: Math.random() * 400, y: Math.random() * 300 },
      connections: []
    }
    setStrategyNodes(prev => [...prev, newNode])
  }

  const removeNode = (nodeId: string) => {
    setStrategyNodes(prev => prev.filter(node => node.id !== nodeId))
    if (selectedNode === nodeId) {
      setSelectedNode(null)
    }
  }

  const runBacktest = async () => {
    setIsBacktesting(true)
    
    // Simulate backtest
    await new Promise(resolve => setTimeout(resolve, 3000))
    
    const mockResults: BacktestResult = {
      period: '1Y',
      totalReturn: 24.5,
      annualizedReturn: 24.5,
      sharpeRatio: 1.85,
      maxDrawdown: 12.3,
      winRate: 68.5,
      totalTrades: 247,
      avgTrade: 0.98,
      profitFactor: 2.15,
      calmarRatio: 1.99,
      sortinoRatio: 2.34,
      volatility: 13.2,
      beta: 0.85,
      alpha: 0.15,
      var95: 2.8,
      expectedShortfall: 4.2
    }
    
    setBacktestResults(mockResults)
    setIsBacktesting(false)
  }

  const nodeTypes = [
    { type: 'indicator', name: 'RSI', icon: LineChart },
    { type: 'indicator', name: 'MACD', icon: TrendingUp },
    { type: 'indicator', name: 'Bollinger Bands', icon: Layers },
    { type: 'indicator', name: 'Moving Average', icon: Activity },
    { type: 'condition', name: 'Price Above', icon: TrendingUp },
    { type: 'condition', name: 'Price Below', icon: TrendingDown },
    { type: 'condition', name: 'Crossover', icon: Target },
    { type: 'condition', name: 'Time Filter', icon: Clock },
    { type: 'action', name: 'Buy Market', icon: Plus },
    { type: 'action', name: 'Sell Market', icon: Minus },
    { type: 'action', name: 'Set Stop Loss', icon: Shield },
    { type: 'action', name: 'Set Take Profit', icon: Target },
    { type: 'logic', name: 'AND', icon: Grid3X3 },
    { type: 'logic', name: 'OR', icon: Layers },
    { type: 'logic', name: 'NOT', icon: XCircle }
  ]

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Brain className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access the algorithmic strategy builder
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Algorithmic Strategy Builder</h2>
          <p className="text-muted-foreground">
            Build, test, and deploy custom trading algorithms with visual programming
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Bot className="w-3 h-3 mr-1" />
            {strategyNodes.length} Nodes
          </Badge>
          <Button onClick={runBacktest} disabled={isBacktesting || strategyNodes.length === 0}>
            {isBacktesting ? (
              <>
                <RefreshCw className="w-4 h-4 mr-2 animate-spin" />
                Testing...
              </>
            ) : (
              <>
                <Play className="w-4 h-4 mr-2" />
                Backtest
              </>
            )}
          </Button>
        </div>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="builder">Strategy Builder</TabsTrigger>
          <TabsTrigger value="templates">Templates</TabsTrigger>
          <TabsTrigger value="backtest">Backtest Results</TabsTrigger>
          <TabsTrigger value="deployment">Deployment</TabsTrigger>
        </TabsList>

        <TabsContent value="builder" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-4 gap-4">
            {/* Node Palette */}
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Node Palette</CardTitle>
              </CardHeader>
              <CardContent className="space-y-3">
                {['indicator', 'condition', 'action', 'logic'].map((category) => (
                  <div key={category}>
                    <h4 className="font-medium text-xs uppercase tracking-wide text-muted-foreground mb-2">
                      {category}
                    </h4>
                    <div className="space-y-1">
                      {nodeTypes.filter(node => node.type === category).map((nodeType) => {
                        const Icon = nodeType.icon
                        return (
                          <Button
                            key={nodeType.name}
                            variant="outline"
                            size="sm"
                            className="w-full justify-start text-xs"
                            onClick={() => addNode(nodeType.type as StrategyNode['type'], nodeType.name)}
                          >
                            <Icon className="w-3 h-3 mr-2" />
                            {nodeType.name}
                          </Button>
                        )
                      })}
                    </div>
                  </div>
                ))}
              </CardContent>
            </Card>

            {/* Strategy Canvas */}
            <div className="lg:col-span-2">
              <Card className="h-[600px]">
                <CardHeader>
                  <CardTitle className="text-sm">Strategy Canvas</CardTitle>
                </CardHeader>
                <CardContent className="relative h-full">
                  <div className="absolute inset-4 border-2 border-dashed border-muted rounded-lg bg-muted/10">
                    {strategyNodes.length === 0 ? (
                      <div className="flex items-center justify-center h-full">
                        <div className="text-center">
                          <Bot className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
                          <h3 className="font-medium mb-2">Start Building Your Strategy</h3>
                          <p className="text-sm text-muted-foreground">
                            Drag nodes from the palette to create your algorithm
                          </p>
                        </div>
                      </div>
                    ) : (
                      <div className="relative h-full">
                        {strategyNodes.map((node) => (
                          <div
                            key={node.id}
                            className={cn(
                              "absolute p-3 border rounded-lg bg-background cursor-pointer transition-colors",
                              selectedNode === node.id ? "border-primary bg-primary/5" : "border-muted hover:border-primary/50"
                            )}
                            style={{
                              left: node.position.x,
                              top: node.position.y,
                              minWidth: '120px'
                            }}
                            onClick={() => setSelectedNode(node.id)}
                          >
                            <div className="flex items-center justify-between mb-1">
                              <span className="text-xs font-medium">{node.name}</span>
                              <Button
                                variant="ghost"
                                size="sm"
                                className="h-4 w-4 p-0"
                                onClick={(e) => {
                                  e.stopPropagation()
                                  removeNode(node.id)
                                }}
                              >
                                <XCircle className="w-3 h-3" />
                              </Button>
                            </div>
                            <Badge variant="outline" className="text-xs">
                              {node.type}
                            </Badge>
                          </div>
                        ))}
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>
            </div>

            {/* Node Properties */}
            <Card>
              <CardHeader>
                <CardTitle className="text-sm">Node Properties</CardTitle>
              </CardHeader>
              <CardContent>
                {selectedNode ? (
                  <div className="space-y-3">
                    {(() => {
                      const node = strategyNodes.find(n => n.id === selectedNode)
                      if (!node) return null
                      
                      return (
                        <div>
                          <h4 className="font-medium mb-2">{node.name}</h4>
                          <div className="space-y-2">
                            {node.type === 'indicator' && (
                              <>
                                <div>
                                  <Label className="text-xs">Period</Label>
                                  <Input type="number" placeholder="14" className="h-8" />
                                </div>
                                <div>
                                  <Label className="text-xs">Source</Label>
                                  <Select>
                                    <SelectTrigger className="h-8">
                                      <SelectValue placeholder="Close" />
                                    </SelectTrigger>
                                    <SelectContent>
                                      <SelectItem value="close">Close</SelectItem>
                                      <SelectItem value="open">Open</SelectItem>
                                      <SelectItem value="high">High</SelectItem>
                                      <SelectItem value="low">Low</SelectItem>
                                    </SelectContent>
                                  </Select>
                                </div>
                              </>
                            )}
                            {node.type === 'condition' && (
                              <>
                                <div>
                                  <Label className="text-xs">Threshold</Label>
                                  <Input type="number" placeholder="0" className="h-8" />
                                </div>
                                <div>
                                  <Label className="text-xs">Operator</Label>
                                  <Select>
                                    <SelectTrigger className="h-8">
                                      <SelectValue placeholder=">" />
                                    </SelectTrigger>
                                    <SelectContent>
                                      <SelectItem value="gt">Greater Than</SelectItem>
                                      <SelectItem value="lt">Less Than</SelectItem>
                                      <SelectItem value="eq">Equal To</SelectItem>
                                    </SelectContent>
                                  </Select>
                                </div>
                              </>
                            )}
                            {node.type === 'action' && (
                              <>
                                <div>
                                  <Label className="text-xs">Quantity (%)</Label>
                                  <Input type="number" placeholder="100" className="h-8" />
                                </div>
                                <div>
                                  <Label className="text-xs">Order Type</Label>
                                  <Select>
                                    <SelectTrigger className="h-8">
                                      <SelectValue placeholder="Market" />
                                    </SelectTrigger>
                                    <SelectContent>
                                      <SelectItem value="market">Market</SelectItem>
                                      <SelectItem value="limit">Limit</SelectItem>
                                      <SelectItem value="stop">Stop</SelectItem>
                                    </SelectContent>
                                  </Select>
                                </div>
                              </>
                            )}
                          </div>
                        </div>
                      )
                    })()}
                  </div>
                ) : (
                  <div className="text-center py-8">
                    <Settings className="w-8 h-8 mx-auto mb-2 text-muted-foreground opacity-50" />
                    <p className="text-xs text-muted-foreground">
                      Select a node to edit its properties
                    </p>
                  </div>
                )}
              </CardContent>
            </Card>
          </div>

          {/* Strategy Actions */}
          <div className="flex gap-3">
            <Button>
              <CheckCircle className="w-4 h-4 mr-2" />
              Save Strategy
            </Button>
            <Button variant="outline">
              <Download className="w-4 h-4 mr-2" />
              Export
            </Button>
            <Button variant="outline">
              <Upload className="w-4 h-4 mr-2" />
              Import
            </Button>
            <Button variant="outline">
              <Copy className="w-4 h-4 mr-2" />
              Clone
            </Button>
          </div>
        </TabsContent>

        <TabsContent value="templates" className="space-y-4">
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
            {strategyTemplates.map((template) => (
              <Card key={template.id} className="cursor-pointer hover:border-primary/50 transition-colors">
                <CardHeader>
                  <div className="flex items-center justify-between">
                    <h4 className="font-bold">{template.name}</h4>
                    <Badge variant={
                      template.difficulty === 'beginner' ? 'default' :
                      template.difficulty === 'intermediate' ? 'secondary' :
                      template.difficulty === 'advanced' ? 'outline' : 'destructive'
                    }>
                      {template.difficulty}
                    </Badge>
                  </div>
                  <p className="text-sm text-muted-foreground">{template.description}</p>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-sm">
                      <span>Category:</span>
                      <Badge variant="outline" className="text-xs">
                        {template.category.replace('_', ' ')}
                      </Badge>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Expected Return:</span>
                      <span className="font-medium text-green-500">{template.expectedReturn}</span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Risk Level:</span>
                      <span className={cn(
                        "font-medium",
                        template.riskLevel === 'low' ? 'text-green-500' :
                        template.riskLevel === 'medium' ? 'text-yellow-500' : 'text-red-500'
                      )}>
                        {template.riskLevel}
                      </span>
                    </div>
                    <div className="flex items-center justify-between text-sm">
                      <span>Timeframe:</span>
                      <span className="font-medium">{template.timeframe}</span>
                    </div>
                    
                    <div className="space-y-1">
                      <div className="text-sm font-medium">Market Conditions:</div>
                      <div className="flex flex-wrap gap-1">
                        {template.marketConditions.map((condition, index) => (
                          <Badge key={index} variant="outline" className="text-xs">
                            {condition}
                          </Badge>
                        ))}
                      </div>
                    </div>

                    <div className="flex gap-2 mt-4">
                      <Button size="sm" className="flex-1">
                        <Copy className="w-3 h-3 mr-1" />
                        Use Template
                      </Button>
                      <Button variant="outline" size="sm">
                        <Eye className="w-3 h-3" />
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="backtest" className="space-y-4">
          {backtestResults ? (
            <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-4 gap-4">
              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <TrendingUp className="w-4 h-4 text-green-500" />
                    <span className="text-sm text-muted-foreground">Total Return</span>
                  </div>
                  <div className="text-2xl font-bold text-green-500">
                    {backtestResults.totalReturn.toFixed(1)}%
                  </div>
                  <div className="text-xs text-muted-foreground">
                    Period: {backtestResults.period}
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <BarChart3 className="w-4 h-4 text-blue-500" />
                    <span className="text-sm text-muted-foreground">Sharpe Ratio</span>
                  </div>
                  <div className="text-2xl font-bold">{backtestResults.sharpeRatio.toFixed(2)}</div>
                  <div className="text-xs text-muted-foreground">
                    Risk-adjusted return
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <TrendingDown className="w-4 h-4 text-red-500" />
                    <span className="text-sm text-muted-foreground">Max Drawdown</span>
                  </div>
                  <div className="text-2xl font-bold text-red-500">
                    {backtestResults.maxDrawdown.toFixed(1)}%
                  </div>
                  <div className="text-xs text-muted-foreground">
                    Worst decline
                  </div>
                </CardContent>
              </Card>

              <Card>
                <CardContent className="p-4">
                  <div className="flex items-center gap-2 mb-2">
                    <Target className="w-4 h-4 text-purple-500" />
                    <span className="text-sm text-muted-foreground">Win Rate</span>
                  </div>
                  <div className="text-2xl font-bold">{backtestResults.winRate.toFixed(1)}%</div>
                  <div className="text-xs text-muted-foreground">
                    {backtestResults.totalTrades} trades
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <div className="text-center py-12">
              <BarChart3 className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
              <h3 className="text-lg font-medium mb-2">No Backtest Results</h3>
              <p className="text-muted-foreground mb-4">
                Run a backtest to see performance metrics and analysis
              </p>
              <Button onClick={runBacktest} disabled={strategyNodes.length === 0}>
                <Play className="w-4 h-4 mr-2" />
                Run Backtest
              </Button>
            </div>
          )}
        </TabsContent>

        <TabsContent value="deployment" className="space-y-4">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Rocket className="w-5 h-5" />
                Deploy Strategy
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-6">
                <Alert>
                  <Shield className="h-4 w-4" />
                  <AlertDescription>
                    Ensure your strategy has been thoroughly backtested before deploying with real funds.
                  </AlertDescription>
                </Alert>

                <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                  <div>
                    <Label htmlFor="initialCapital">Initial Capital (USD)</Label>
                    <Input
                      id="initialCapital"
                      type="number"
                      placeholder="1000"
                      step="100"
                    />
                  </div>
                  <div>
                    <Label htmlFor="maxRisk">Max Risk per Trade (%)</Label>
                    <Input
                      id="maxRisk"
                      type="number"
                      placeholder="2"
                      step="0.1"
                      max="10"
                    />
                  </div>
                  <div>
                    <Label htmlFor="tradingPair">Trading Pair</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Select pair" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="BTC/USDT">BTC/USDT</SelectItem>
                        <SelectItem value="ETH/USDT">ETH/USDT</SelectItem>
                        <SelectItem value="BNB/USDT">BNB/USDT</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                  <div>
                    <Label htmlFor="timeframe">Timeframe</Label>
                    <Select>
                      <SelectTrigger>
                        <SelectValue placeholder="Select timeframe" />
                      </SelectTrigger>
                      <SelectContent>
                        <SelectItem value="1m">1 Minute</SelectItem>
                        <SelectItem value="5m">5 Minutes</SelectItem>
                        <SelectItem value="15m">15 Minutes</SelectItem>
                        <SelectItem value="1h">1 Hour</SelectItem>
                        <SelectItem value="4h">4 Hours</SelectItem>
                        <SelectItem value="1d">1 Day</SelectItem>
                      </SelectContent>
                    </Select>
                  </div>
                </div>

                <div className="space-y-4">
                  <div className="flex items-center space-x-2">
                    <Switch id="paperTrading" />
                    <Label htmlFor="paperTrading">Start with paper trading</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Switch id="notifications" />
                    <Label htmlFor="notifications">Enable trade notifications</Label>
                  </div>
                  <div className="flex items-center space-x-2">
                    <Switch id="autoRebalance" />
                    <Label htmlFor="autoRebalance">Auto-rebalance portfolio</Label>
                  </div>
                </div>

                <div className="flex gap-3">
                  <Button className="flex-1">
                    <Rocket className="w-4 h-4 mr-2" />
                    Deploy Strategy
                  </Button>
                  <Button variant="outline">
                    <Eye className="w-4 h-4 mr-2" />
                    Preview
                  </Button>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
