import { useState, useEffect, useCallback } from 'react'
import { type Address } from 'viem'
import { 
  smartContractSecurityScanner,
  type SecurityScanResult,
  type ScanConfiguration,
  type RiskLevel,
  type Vulnerability,
  type SecurityScanEvent
} from '@/lib/security-scanner'
import { toast } from 'sonner'

export interface SecurityScannerState {
  scanResults: SecurityScanResult[]
  currentScan: SecurityScanResult | null
  scanHistory: Map<string, SecurityScanResult[]>
  isScanning: boolean
  error: string | null
  lastUpdate: number | null
}

export interface UseSecurityScannerOptions {
  enableNotifications?: boolean
  autoRefresh?: boolean
  refreshInterval?: number
}

export interface UseSecurityScannerReturn {
  // State
  state: SecurityScannerState
  
  // Scan Operations
  scanContract: (
    contractAddress: Address,
    chainId: number,
    config?: ScanConfiguration
  ) => Promise<SecurityScanResult>
  
  // Data Access
  getScanResult: (scanId: string) => SecurityScanResult | null
  getScanHistory: (contractAddress: Address) => SecurityScanResult[]
  
  // Analysis
  getRiskSummary: (result: SecurityScanResult) => RiskSummary
  getVulnerabilitySummary: (vulnerabilities: Vulnerability[]) => VulnerabilitySummary
  
  // Utilities
  clearError: () => void
  refresh: () => void
}

export interface RiskSummary {
  overallRisk: RiskLevel
  riskScore: number
  criticalIssues: number
  highIssues: number
  mediumIssues: number
  lowIssues: number
  auditStatus: 'audited' | 'not_audited' | 'partial'
  recommendations: number
}

export interface VulnerabilitySummary {
  total: number
  critical: number
  high: number
  medium: number
  low: number
  mostCommonType: string
  averageConfidence: number
  exploitabilityScore: number
}

export const useSecurityScanner = (
  options: UseSecurityScannerOptions = {}
): UseSecurityScannerReturn => {
  const {
    enableNotifications = true,
    autoRefresh = false,
    refreshInterval = 300000 // 5 minutes
  } = options

  const [state, setState] = useState<SecurityScannerState>({
    scanResults: [],
    currentScan: null,
    scanHistory: new Map(),
    isScanning: false,
    error: null,
    lastUpdate: null
  })

  // Handle security scan events
  const handleSecurityScanEvent = useCallback((event: SecurityScanEvent) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'scan_started':
          toast.info('Security Scan Started', {
            description: 'Analyzing smart contract for security vulnerabilities'
          })
          break
        case 'scan_completed':
          if (event.result) {
            const riskColor = getRiskColor(event.result.overallRisk)
            toast.success('Security Scan Completed', {
              description: `Risk Level: ${event.result.overallRisk.toUpperCase()} (${event.result.vulnerabilities.length} issues found)`,
              className: riskColor
            })
          }
          break
        case 'scan_failed':
          toast.error('Security Scan Failed', {
            description: event.error?.message || 'Failed to complete security scan'
          })
          break
        case 'vulnerability_found':
          if (event.vulnerability && event.vulnerability.severity === 'critical') {
            toast.error('Critical Vulnerability Found', {
              description: event.vulnerability.title,
              duration: 10000
            })
          }
          break
      }
    }

    // Update state after event
    if (event.result) {
      setState(prev => {
        const newResults = [...prev.scanResults]
        const existingIndex = newResults.findIndex(r => r.scanId === event.result!.scanId)
        
        if (existingIndex >= 0) {
          newResults[existingIndex] = event.result!
        } else {
          newResults.push(event.result!)
        }

        // Update history
        const newHistory = new Map(prev.scanHistory)
        const contractKey = event.result!.contractAddress.toLowerCase()
        const contractHistory = newHistory.get(contractKey) || []
        
        const historyIndex = contractHistory.findIndex(r => r.scanId === event.result!.scanId)
        if (historyIndex >= 0) {
          contractHistory[historyIndex] = event.result!
        } else {
          contractHistory.push(event.result!)
        }
        
        newHistory.set(contractKey, contractHistory)

        return {
          ...prev,
          scanResults: newResults,
          currentScan: event.result!,
          scanHistory: newHistory,
          isScanning: event.type === 'scan_started',
          error: event.type === 'scan_failed' ? event.error?.message || 'Scan failed' : null,
          lastUpdate: Date.now()
        }
      })
    } else if (event.type === 'scan_started') {
      setState(prev => ({ ...prev, isScanning: true, error: null }))
    } else if (event.type === 'scan_failed') {
      setState(prev => ({ 
        ...prev, 
        isScanning: false, 
        error: event.error?.message || 'Scan failed' 
      }))
    }
  }, [enableNotifications])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = smartContractSecurityScanner.addEventListener(handleSecurityScanEvent)

    return () => {
      unsubscribe()
    }
  }, [handleSecurityScanEvent])

  // Auto-refresh
  useEffect(() => {
    if (autoRefresh && refreshInterval > 0) {
      const interval = setInterval(() => {
        refresh()
      }, refreshInterval)

      return () => clearInterval(interval)
    }
  }, [autoRefresh, refreshInterval])

  // Scan contract
  const scanContract = useCallback(async (
    contractAddress: Address,
    chainId: number,
    config?: ScanConfiguration
  ): Promise<SecurityScanResult> => {
    setState(prev => ({ ...prev, isScanning: true, error: null }))

    try {
      const result = await smartContractSecurityScanner.scanContract(contractAddress, chainId, config)
      return result
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isScanning: false,
        error: errorMessage
      }))
      throw error
    }
  }, [])

  // Get scan result
  const getScanResult = useCallback((scanId: string): SecurityScanResult | null => {
    return smartContractSecurityScanner.getScanResult(scanId)
  }, [])

  // Get scan history
  const getScanHistory = useCallback((contractAddress: Address): SecurityScanResult[] => {
    return smartContractSecurityScanner.getScanHistory(contractAddress)
  }, [])

  // Get risk summary
  const getRiskSummary = useCallback((result: SecurityScanResult): RiskSummary => {
    const vulnerabilities = result.vulnerabilities
    
    return {
      overallRisk: result.overallRisk,
      riskScore: result.riskScore,
      criticalIssues: vulnerabilities.filter(v => v.severity === 'critical').length,
      highIssues: vulnerabilities.filter(v => v.severity === 'high').length,
      mediumIssues: vulnerabilities.filter(v => v.severity === 'medium').length,
      lowIssues: vulnerabilities.filter(v => v.severity === 'low' || v.severity === 'very_low').length,
      auditStatus: result.auditStatus.isAudited ? 'audited' : 'not_audited',
      recommendations: result.recommendations.length
    }
  }, [])

  // Get vulnerability summary
  const getVulnerabilitySummary = useCallback((vulnerabilities: Vulnerability[]): VulnerabilitySummary => {
    const total = vulnerabilities.length
    const critical = vulnerabilities.filter(v => v.severity === 'critical').length
    const high = vulnerabilities.filter(v => v.severity === 'high').length
    const medium = vulnerabilities.filter(v => v.severity === 'medium').length
    const low = vulnerabilities.filter(v => v.severity === 'low' || v.severity === 'very_low').length

    // Find most common vulnerability type
    const typeCounts = vulnerabilities.reduce((acc, v) => {
      acc[v.type] = (acc[v.type] || 0) + 1
      return acc
    }, {} as Record<string, number>)

    const mostCommonType = Object.entries(typeCounts)
      .sort(([,a], [,b]) => b - a)[0]?.[0] || 'none'

    const averageConfidence = total > 0 
      ? vulnerabilities.reduce((sum, v) => sum + v.confidence, 0) / total
      : 0

    const exploitabilityScore = total > 0
      ? vulnerabilities.reduce((sum, v) => sum + v.exploitability, 0) / total
      : 0

    return {
      total,
      critical,
      high,
      medium,
      low,
      mostCommonType,
      averageConfidence,
      exploitabilityScore
    }
  }, [])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  // Refresh
  const refresh = useCallback(() => {
    setState(prev => ({ ...prev, lastUpdate: Date.now() }))
  }, [])

  return {
    state,
    scanContract,
    getScanResult,
    getScanHistory,
    getRiskSummary,
    getVulnerabilitySummary,
    clearError,
    refresh
  }
}

// Helper function to get risk color
function getRiskColor(risk: RiskLevel): string {
  switch (risk) {
    case 'critical':
      return 'border-red-500 bg-red-50 text-red-900'
    case 'high':
      return 'border-orange-500 bg-orange-50 text-orange-900'
    case 'medium':
      return 'border-yellow-500 bg-yellow-50 text-yellow-900'
    case 'low':
      return 'border-blue-500 bg-blue-50 text-blue-900'
    case 'very_low':
      return 'border-green-500 bg-green-50 text-green-900'
    default:
      return 'border-gray-500 bg-gray-50 text-gray-900'
  }
}

// Simplified hook for quick security checks
export const useQuickSecurityCheck = () => {
  const { scanContract, state } = useSecurityScanner()

  const quickScan = useCallback(async (contractAddress: Address, chainId: number) => {
    const config: ScanConfiguration = {
      depth: 'basic',
      includeAuditCheck: true,
      includeGasAnalysis: false,
      includeComplianceCheck: false,
      customRules: [],
      excludePatterns: []
    }

    return scanContract(contractAddress, chainId, config)
  }, [scanContract])

  return {
    quickScan,
    isScanning: state.isScanning,
    currentScan: state.currentScan,
    error: state.error
  }
}

// Hook for vulnerability analysis
export const useVulnerabilityAnalysis = () => {
  const { state, getVulnerabilitySummary } = useSecurityScanner()

  const analyzeVulnerabilities = useCallback((scanResult: SecurityScanResult) => {
    const summary = getVulnerabilitySummary(scanResult.vulnerabilities)
    
    const analysis = {
      summary,
      riskDistribution: {
        critical: (summary.critical / summary.total) * 100,
        high: (summary.high / summary.total) * 100,
        medium: (summary.medium / summary.total) * 100,
        low: (summary.low / summary.total) * 100
      },
      topVulnerabilities: scanResult.vulnerabilities
        .sort((a, b) => {
          const severityOrder = { critical: 4, high: 3, medium: 2, low: 1, very_low: 0 }
          return severityOrder[b.severity] - severityOrder[a.severity]
        })
        .slice(0, 5),
      recommendations: scanResult.recommendations
        .filter(r => r.priority === 'high')
        .slice(0, 3)
    }

    return analysis
  }, [getVulnerabilitySummary])

  return {
    analyzeVulnerabilities,
    scanResults: state.scanResults
  }
}

// Hook for security monitoring
export const useSecurityMonitoring = () => {
  const { state, scanContract } = useSecurityScanner()

  const monitorContract = useCallback(async (
    contractAddress: Address,
    chainId: number,
    interval: number = 3600000 // 1 hour
  ) => {
    const monitor = setInterval(async () => {
      try {
        await scanContract(contractAddress, chainId, {
          depth: 'standard',
          includeAuditCheck: true,
          includeGasAnalysis: true,
          includeComplianceCheck: true,
          customRules: [],
          excludePatterns: []
        })
      } catch (error) {
        console.error('Security monitoring scan failed:', error)
      }
    }, interval)

    return () => clearInterval(monitor)
  }, [scanContract])

  return {
    monitorContract,
    scanHistory: state.scanHistory,
    isScanning: state.isScanning
  }
}
