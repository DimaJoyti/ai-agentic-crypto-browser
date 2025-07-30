// Trading API client for HFT system

const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080'

interface ApiResponse<T> {
  data?: T
  error?: string
  message?: string
}

class TradingApiClient {
  private baseUrl: string

  constructor(baseUrl: string = API_BASE_URL) {
    this.baseUrl = baseUrl
  }

  private async request<T>(endpoint: string, options?: RequestInit): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`
    
    try {
      const response = await fetch(url, {
        headers: {
          'Content-Type': 'application/json',
          ...options?.headers,
        },
        ...options,
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      return data
    } catch (error) {
      console.error(`API request failed for ${endpoint}:`, error)
      throw error
    }
  }

  // HFT Engine endpoints
  async getHFTStatus() {
    return this.request('/api/hft/status')
  }

  async startHFTEngine() {
    return this.request('/api/hft/start', { method: 'POST' })
  }

  async stopHFTEngine() {
    return this.request('/api/hft/stop', { method: 'POST' })
  }

  async getHFTMetrics() {
    return this.request('/api/hft/metrics')
  }

  async getHFTConfig() {
    return this.request('/api/hft/config')
  }

  async updateHFTConfig(config: any) {
    return this.request('/api/hft/config', {
      method: 'PUT',
      body: JSON.stringify(config),
    })
  }

  // Trading endpoints
  async getOrders(limit?: number) {
    const params = limit ? `?limit=${limit}` : ''
    return this.request(`/api/trading/orders${params}`)
  }

  async createOrder(order: any) {
    return this.request('/api/trading/orders', {
      method: 'POST',
      body: JSON.stringify(order),
    })
  }

  async cancelOrder(orderId: string) {
    return this.request(`/api/trading/orders/${orderId}`, {
      method: 'DELETE',
    })
  }

  async getPositions() {
    return this.request('/api/trading/positions')
  }

  async getSignals() {
    return this.request('/api/trading/signals')
  }

  async getOrderbook(symbol: string, limit?: number) {
    const params = limit ? `?limit=${limit}` : ''
    return this.request(`/api/trading/orderbook/${symbol}${params}`)
  }

  // Portfolio endpoints
  async getPortfolioSummary() {
    return this.request('/api/portfolio/summary')
  }

  async getPortfolioPositions() {
    return this.request('/api/portfolio/positions')
  }

  async getPortfolioMetrics() {
    return this.request('/api/portfolio/metrics')
  }

  async getPortfolioRisk() {
    return this.request('/api/portfolio/risk')
  }

  // Strategy endpoints
  async getStrategies() {
    return this.request('/api/strategies')
  }

  async createStrategy(strategy: any) {
    return this.request('/api/strategies', {
      method: 'POST',
      body: JSON.stringify(strategy),
    })
  }

  async getStrategy(id: string) {
    return this.request(`/api/strategies/${id}`)
  }

  async updateStrategy(id: string, strategy: any) {
    return this.request(`/api/strategies/${id}`, {
      method: 'PUT',
      body: JSON.stringify(strategy),
    })
  }

  async deleteStrategy(id: string) {
    return this.request(`/api/strategies/${id}`, {
      method: 'DELETE',
    })
  }

  async startStrategy(id: string) {
    return this.request(`/api/strategies/${id}/start`, {
      method: 'POST',
    })
  }

  async stopStrategy(id: string) {
    return this.request(`/api/strategies/${id}/stop`, {
      method: 'POST',
    })
  }

  async getStrategyPerformance(id: string) {
    return this.request(`/api/strategies/${id}/performance`)
  }

  // Risk Management endpoints
  async getRiskLimits() {
    return this.request('/api/risk/limits')
  }

  async createRiskLimit(limit: any) {
    return this.request('/api/risk/limits', {
      method: 'POST',
      body: JSON.stringify(limit),
    })
  }

  async getRiskLimit(id: string) {
    return this.request(`/api/risk/limits/${id}`)
  }

  async updateRiskLimit(id: string, limit: any) {
    return this.request(`/api/risk/limits/${id}`, {
      method: 'PUT',
      body: JSON.stringify(limit),
    })
  }

  async deleteRiskLimit(id: string) {
    return this.request(`/api/risk/limits/${id}`, {
      method: 'DELETE',
    })
  }

  async getRiskViolations() {
    return this.request('/api/risk/violations')
  }

  async getRiskMetrics() {
    return this.request('/api/risk/metrics')
  }

  async emergencyStop(reason: string, options: any = {}) {
    return this.request('/api/risk/emergency-stop', {
      method: 'POST',
      body: JSON.stringify({
        reason,
        ...options,
      }),
    })
  }

  // MCP Integration endpoints
  async getMCPInsights() {
    return this.request('/api/mcp/insights')
  }

  async getMCPInsight(symbol: string) {
    return this.request(`/api/mcp/insights/${symbol}`)
  }

  async getMCPSentiment(symbol: string) {
    return this.request(`/api/mcp/sentiment/${symbol}`)
  }

  async getMCPNews(symbol: string) {
    return this.request(`/api/mcp/news/${symbol}`)
  }

  // TradingView endpoints
  async getTradingViewCharts() {
    return this.request('/api/tradingview/charts')
  }

  async createTradingViewChart(chart: any) {
    return this.request('/api/tradingview/charts', {
      method: 'POST',
      body: JSON.stringify(chart),
    })
  }

  async getTradingViewChart(id: string) {
    return this.request(`/api/tradingview/charts/${id}`)
  }

  async deleteTradingViewChart(id: string) {
    return this.request(`/api/tradingview/charts/${id}`, {
      method: 'DELETE',
    })
  }

  async getTradingViewSignals() {
    return this.request('/api/tradingview/signals')
  }

  async getTradingViewIndicators(symbol: string) {
    return this.request(`/api/tradingview/indicators/${symbol}`)
  }

  // System endpoints
  async getSystemStatus() {
    return this.request('/api/system/status')
  }

  async getSystemConfig() {
    return this.request('/api/system/config')
  }
}

// Export singleton instance
export const tradingApi = new TradingApiClient()

// Export types for use in components
export type {
  ApiResponse,
}

// Export the class for custom instances
export { TradingApiClient }
