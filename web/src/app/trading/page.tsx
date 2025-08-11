'use client'

import { TradingInterface } from '@/components/trading/TradingInterface'
import { TradingDashboard } from '@/components/trading/TradingDashboard'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { BarChart3, TrendingUp, Zap } from 'lucide-react'

export default function TradingPage() {
  return (
    <div className="min-h-screen bg-background">
      <div className="container mx-auto p-4">
        <div className="mb-6">
          <div className="flex items-center justify-between">
            <div>
              <h1 className="text-3xl font-bold">Trading</h1>
              <p className="text-muted-foreground">Advanced cryptocurrency trading platform</p>
            </div>
            <div className="flex items-center gap-2">
              <Badge variant="default" className="bg-green-100 text-green-700">
                <TrendingUp className="w-3 h-3 mr-1" />
                Market Open
              </Badge>
              <Badge variant="outline">
                <Zap className="w-3 h-3 mr-1" />
                Real-time Data
              </Badge>
            </div>
          </div>
        </div>

        <Tabs defaultValue="advanced" className="space-y-4">
          <TabsList>
            <TabsTrigger value="advanced" className="flex items-center gap-2">
              <BarChart3 className="w-4 h-4" />
              Advanced Trading
            </TabsTrigger>
            <TabsTrigger value="dashboard">
              Trading Dashboard
            </TabsTrigger>
          </TabsList>

          <TabsContent value="advanced">
            <TradingInterface />
          </TabsContent>

          <TabsContent value="dashboard">
            <TradingDashboard />
          </TabsContent>
        </Tabs>
      </div>
    </div>
  )
}
