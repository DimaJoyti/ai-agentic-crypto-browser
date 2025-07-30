import { type Address, type Hash } from 'viem'

export interface LendingProtocol {
  id: string
  name: string
  version: string
  chainId: number
  contractAddresses: ProtocolContracts
  supportedFeatures: LendingFeature[]
  riskParameters: RiskParameters
  metadata: ProtocolMetadata
}

export interface ProtocolContracts {
  lendingPool: Address
  dataProvider: Address
  priceOracle: Address
  incentivesController?: Address
  wethGateway?: Address
  flashLoanReceiver?: Address
}

export interface ProtocolMetadata {
  description: string
  website: string
  documentation: string
  logo: string
  tvl: string
  totalBorrowed: string
  totalSupplied: string
  activeUsers: number
  security: SecurityInfo
}

export interface SecurityInfo {
  audited: boolean
  auditors: string[]
  lastAudit?: string
  bugBounty?: string
  insuranceFund?: string
  riskScore: number // 0-100
}

export interface RiskParameters {
  maxLTV: number // Maximum Loan-to-Value ratio
  liquidationThreshold: number
  liquidationPenalty: number
  reserveFactor: number
  optimalUtilization: number
  baseVariableBorrowRate: number
  variableRateSlope1: number
  variableRateSlope2: number
}

export enum LendingFeature {
  LENDING = 'lending',
  BORROWING = 'borrowing',
  FLASH_LOANS = 'flash_loans',
  COLLATERAL_SWAP = 'collateral_swap',
  DEBT_SWAP = 'debt_swap',
  LIQUIDATIONS = 'liquidations',
  GOVERNANCE = 'governance',
  INCENTIVES = 'incentives',
  CREDIT_DELEGATION = 'credit_delegation'
}

export interface LendingAsset {
  address: Address
  symbol: string
  name: string
  decimals: number
  chainId: number
  aTokenAddress: Address
  stableDebtTokenAddress?: Address
  variableDebtTokenAddress?: Address
  interestRateStrategyAddress: Address
  isActive: boolean
  isFrozen: boolean
  isPaused: boolean
  canBeCollateral: boolean
  canBeBorrowed: boolean
  canBeFlashLoaned: boolean
  ltv: number
  liquidationThreshold: number
  liquidationBonus: number
  reserveFactor: number
  usageAsCollateralEnabled: boolean
  borrowingEnabled: boolean
  stableBorrowRateEnabled: boolean
  flashLoanEnabled: boolean
  rates: AssetRates
  caps: AssetCaps
  priceUSD: number
  marketSize: string
  totalBorrowed: string
  availableLiquidity: string
  utilizationRate: number
}

export interface AssetRates {
  supplyAPY: number
  variableBorrowAPY: number
  stableBorrowAPY: number
  liquidityRate: number
  variableBorrowRate: number
  stableBorrowRate: number
  averageStableBorrowRate: number
  incentiveAPY?: number
  totalAPY: number
}

export interface AssetCaps {
  supplyCap: string
  borrowCap: string
  debtCeiling: string
}

export interface UserPosition {
  id: string
  protocolId: string
  userAddress: Address
  totalCollateralETH: string
  totalDebtETH: string
  availableBorrowsETH: string
  currentLiquidationThreshold: number
  ltv: number
  healthFactor: number
  supplies: SupplyPosition[]
  borrows: BorrowPosition[]
  netAPY: number
  totalSupplyUSD: number
  totalBorrowUSD: number
  netWorthUSD: number
  lastUpdate: number
  isAtRisk: boolean
  riskLevel: 'low' | 'medium' | 'high' | 'critical'
}

export interface SupplyPosition {
  asset: LendingAsset
  amount: string
  amountUSD: number
  aTokenBalance: string
  currentATokenBalance: string
  supplyAPY: number
  incentiveAPY: number
  totalAPY: number
  isCollateral: boolean
  canWithdraw: string
  canWithdrawUSD: number
}

export interface BorrowPosition {
  asset: LendingAsset
  amount: string
  amountUSD: number
  borrowAPY: number
  incentiveAPY: number
  netAPY: number
  rateMode: 'variable' | 'stable'
  canRepay: string
  canRepayUSD: number
  stableRate?: number
}

export interface LendingTransaction {
  id: string
  hash?: Hash
  protocolId: string
  userAddress: Address
  type: TransactionType
  asset: LendingAsset
  amount: string
  amountUSD: number
  gasLimit: string
  gasPrice: string
  gasUsed?: string
  status: 'pending' | 'confirmed' | 'failed' | 'reverted'
  timestamp: number
  blockNumber?: number
  error?: string
  revertReason?: string
  metadata?: TransactionMetadata
}

export interface TransactionMetadata {
  interestRateMode?: 'variable' | 'stable'
  useAsCollateral?: boolean
  onBehalfOf?: Address
  referralCode?: number
  flashLoanParams?: FlashLoanParams
}

export interface FlashLoanParams {
  assets: Address[]
  amounts: string[]
  modes: number[]
  params: string
  onBehalfOf: Address
  referralCode: number
}

export enum TransactionType {
  SUPPLY = 'supply',
  WITHDRAW = 'withdraw',
  BORROW = 'borrow',
  REPAY = 'repay',
  SWAP_BORROW_RATE = 'swap_borrow_rate',
  SET_USE_AS_COLLATERAL = 'set_use_as_collateral',
  LIQUIDATION_CALL = 'liquidation_call',
  FLASH_LOAN = 'flash_loan'
}

export interface YieldOpportunity {
  id: string
  protocolId: string
  asset: LendingAsset
  strategy: YieldStrategy
  currentAPY: number
  projectedAPY: number
  riskScore: number
  tvl: string
  minimumDeposit: string
  lockupPeriod?: number
  autoCompounding: boolean
  incentives: IncentiveInfo[]
  fees: FeeStructure
}

export interface YieldStrategy {
  type: 'supply' | 'leverage' | 'yield_farming' | 'liquidity_mining'
  description: string
  complexity: 'simple' | 'intermediate' | 'advanced'
  riskLevel: 'low' | 'medium' | 'high'
  steps: StrategyStep[]
}

export interface StrategyStep {
  action: string
  description: string
  estimatedGas: string
  riskFactors: string[]
}

export interface IncentiveInfo {
  token: Address
  symbol: string
  apy: number
  distributionEnd?: number
  claimable: string
}

export interface FeeStructure {
  managementFee: number
  performanceFee: number
  withdrawalFee: number
  flashLoanFee: number
}

export interface LiquidationOpportunity {
  id: string
  protocolId: string
  userAddress: Address
  healthFactor: number
  totalCollateralETH: string
  totalDebtETH: string
  liquidationThreshold: number
  maxLiquidatableDebt: string
  collateralAsset: LendingAsset
  debtAsset: LendingAsset
  liquidationBonus: number
  profitability: number
  gasEstimate: string
  deadline: number
}

export interface LendingConfig {
  defaultSlippage: number
  maxSlippage: number
  healthFactorTarget: number
  minHealthFactor: number
  autoRebalance: boolean
  enableLiquidationProtection: boolean
  maxLeverage: number
  preferredRateMode: 'variable' | 'stable'
  enableFlashLoans: boolean
  gasMultiplier: number
}

export class LendingIntegration {
  private static instance: LendingIntegration
  private protocols = new Map<string, LendingProtocol>()
  private assets = new Map<string, LendingAsset>()
  private positions = new Map<string, UserPosition>()
  private transactions = new Map<string, LendingTransaction>()
  private opportunities = new Map<string, YieldOpportunity>()
  private liquidations = new Map<string, LiquidationOpportunity>()
  private config: LendingConfig
  private eventListeners = new Set<(event: LendingEvent) => void>()

  private constructor() {
    this.config = {
      defaultSlippage: 0.5,
      maxSlippage: 5,
      healthFactorTarget: 2.0,
      minHealthFactor: 1.1,
      autoRebalance: false,
      enableLiquidationProtection: true,
      maxLeverage: 3,
      preferredRateMode: 'variable',
      enableFlashLoans: true,
      gasMultiplier: 1.2
    }

    this.initializeProtocols()
  }

  static getInstance(): LendingIntegration {
    if (!LendingIntegration.instance) {
      LendingIntegration.instance = new LendingIntegration()
    }
    return LendingIntegration.instance
  }

  /**
   * Initialize supported lending protocols
   */
  private initializeProtocols(): void {
    // Aave V3
    this.protocols.set('aave-v3', {
      id: 'aave-v3',
      name: 'Aave V3',
      version: '3.0.0',
      chainId: 1,
      contractAddresses: {
        lendingPool: '0x87870Bca3F3fD6335C3F4ce8392D69350B4fA4E2',
        dataProvider: '0x7B4EB56E7CD4b454BA8ff71E4518426369a138a3',
        priceOracle: '0x54586bE62E3c3580375aE3723C145253060Ca0C2',
        incentivesController: '0x8164Cc65827dcFe994AB23944CBC90e0aa80bFcb',
        wethGateway: '0xD322A49006FC828F9B5B37Ab215F99B4E5caB19C'
      },
      supportedFeatures: [
        LendingFeature.LENDING,
        LendingFeature.BORROWING,
        LendingFeature.FLASH_LOANS,
        LendingFeature.COLLATERAL_SWAP,
        LendingFeature.DEBT_SWAP,
        LendingFeature.LIQUIDATIONS,
        LendingFeature.INCENTIVES
      ],
      riskParameters: {
        maxLTV: 80,
        liquidationThreshold: 85,
        liquidationPenalty: 5,
        reserveFactor: 10,
        optimalUtilization: 80,
        baseVariableBorrowRate: 0,
        variableRateSlope1: 4,
        variableRateSlope2: 60
      },
      metadata: {
        description: 'Leading decentralized lending protocol with advanced features',
        website: 'https://aave.com',
        documentation: 'https://docs.aave.com',
        logo: '/logos/aave.svg',
        tvl: '10500000000',
        totalBorrowed: '7200000000',
        totalSupplied: '10500000000',
        activeUsers: 125000,
        security: {
          audited: true,
          auditors: ['Trail of Bits', 'OpenZeppelin', 'Consensys Diligence'],
          lastAudit: '2022-12-01',
          bugBounty: 'https://github.com/aave/bug-bounty',
          insuranceFund: '550000000',
          riskScore: 15
        }
      }
    })

    // Compound V3
    this.protocols.set('compound-v3', {
      id: 'compound-v3',
      name: 'Compound V3',
      version: '3.0.0',
      chainId: 1,
      contractAddresses: {
        lendingPool: '0xc3d688B66703497DAA19211EEdff47f25384cdc3',
        dataProvider: '0x0000000000000000000000000000000000000000',
        priceOracle: '0xdbd020CAeF83eFd542f4De03e3cF0C28A4428bd5'
      },
      supportedFeatures: [
        LendingFeature.LENDING,
        LendingFeature.BORROWING,
        LendingFeature.LIQUIDATIONS,
        LendingFeature.GOVERNANCE
      ],
      riskParameters: {
        maxLTV: 75,
        liquidationThreshold: 80,
        liquidationPenalty: 8,
        reserveFactor: 15,
        optimalUtilization: 85,
        baseVariableBorrowRate: 2,
        variableRateSlope1: 5,
        variableRateSlope2: 50
      },
      metadata: {
        description: 'Autonomous interest rate protocol with governance',
        website: 'https://compound.finance',
        documentation: 'https://docs.compound.finance',
        logo: '/logos/compound.svg',
        tvl: '3200000000',
        totalBorrowed: '2100000000',
        totalSupplied: '3200000000',
        activeUsers: 85000,
        security: {
          audited: true,
          auditors: ['Trail of Bits', 'OpenZeppelin'],
          lastAudit: '2022-08-01',
          bugBounty: 'https://compound.finance/security',
          riskScore: 20
        }
      }
    })

    // Morpho
    this.protocols.set('morpho', {
      id: 'morpho',
      name: 'Morpho',
      version: '1.0.0',
      chainId: 1,
      contractAddresses: {
        lendingPool: '0x777777c9898D384F785Ee44Acfe945efDFf5f3E0',
        dataProvider: '0x0000000000000000000000000000000000000000',
        priceOracle: '0x0000000000000000000000000000000000000000'
      },
      supportedFeatures: [
        LendingFeature.LENDING,
        LendingFeature.BORROWING,
        LendingFeature.GOVERNANCE
      ],
      riskParameters: {
        maxLTV: 85,
        liquidationThreshold: 90,
        liquidationPenalty: 3,
        reserveFactor: 5,
        optimalUtilization: 90,
        baseVariableBorrowRate: 0,
        variableRateSlope1: 3,
        variableRateSlope2: 40
      },
      metadata: {
        description: 'Lending optimizer built on top of Aave and Compound',
        website: 'https://morpho.xyz',
        documentation: 'https://docs.morpho.xyz',
        logo: '/logos/morpho.svg',
        tvl: '850000000',
        totalBorrowed: '520000000',
        totalSupplied: '850000000',
        activeUsers: 12000,
        security: {
          audited: true,
          auditors: ['Spearbit', 'Cantina'],
          lastAudit: '2023-06-01',
          bugBounty: 'https://immunefi.com/bounty/morpho',
          riskScore: 25
        }
      }
    })
  }

  /**
   * Get user position across all protocols
   */
  async getUserPosition(userAddress: Address, protocolId?: string): Promise<UserPosition[]> {
    const targetProtocols = protocolId ? [protocolId] : Array.from(this.protocols.keys())
    const positions: UserPosition[] = []

    for (const protocol of targetProtocols) {
      try {
        const position = await this.getUserPositionForProtocol(userAddress, protocol)
        if (position) {
          positions.push(position)
        }
      } catch (error) {
        console.error(`Failed to get position for ${protocol}:`, error)
      }
    }

    return positions
  }

  /**
   * Get user position for specific protocol
   */
  private async getUserPositionForProtocol(userAddress: Address, protocolId: string): Promise<UserPosition | null> {
    const protocol = this.protocols.get(protocolId)
    if (!protocol) {
      throw new Error(`Protocol not found: ${protocolId}`)
    }

    // Mock implementation - in real app, this would call actual protocol contracts
    const mockPosition = this.generateMockPosition(userAddress, protocolId)
    
    const positionKey = `${protocolId}_${userAddress.toLowerCase()}`
    this.positions.set(positionKey, mockPosition)

    return mockPosition
  }

  /**
   * Generate mock position data
   */
  private generateMockPosition(userAddress: Address, protocolId: string): UserPosition {
    const mockSupplies: SupplyPosition[] = [
      {
        asset: this.getMockAsset('ETH'),
        amount: '5.5',
        amountUSD: 9900,
        aTokenBalance: '5.5123',
        currentATokenBalance: '5.5123',
        supplyAPY: 3.2,
        incentiveAPY: 0.8,
        totalAPY: 4.0,
        isCollateral: true,
        canWithdraw: '3.2',
        canWithdrawUSD: 5760
      }
    ]

    const mockBorrows: BorrowPosition[] = [
      {
        asset: this.getMockAsset('USDC'),
        amount: '5000',
        amountUSD: 5000,
        borrowAPY: 4.5,
        incentiveAPY: 0.5,
        netAPY: 4.0,
        rateMode: 'variable',
        canRepay: '5000',
        canRepayUSD: 5000
      }
    ]

    const totalSupplyUSD = mockSupplies.reduce((sum, supply) => sum + supply.amountUSD, 0)
    const totalBorrowUSD = mockBorrows.reduce((sum, borrow) => sum + borrow.amountUSD, 0)
    const ltv = totalBorrowUSD / totalSupplyUSD * 100
    const healthFactor = (totalSupplyUSD * 0.85) / totalBorrowUSD // Simplified calculation

    return {
      id: `${protocolId}_${userAddress}`,
      protocolId,
      userAddress,
      totalCollateralETH: '5.5',
      totalDebtETH: '2.78',
      availableBorrowsETH: '1.95',
      currentLiquidationThreshold: 85,
      ltv,
      healthFactor,
      supplies: mockSupplies,
      borrows: mockBorrows,
      netAPY: 0.5,
      totalSupplyUSD,
      totalBorrowUSD,
      netWorthUSD: totalSupplyUSD - totalBorrowUSD,
      lastUpdate: Date.now(),
      isAtRisk: healthFactor < 1.5,
      riskLevel: healthFactor < 1.2 ? 'critical' : healthFactor < 1.5 ? 'high' : healthFactor < 2.0 ? 'medium' : 'low'
    }
  }

  /**
   * Get mock asset data
   */
  private getMockAsset(symbol: string): LendingAsset {
    const mockAssets: Record<string, Partial<LendingAsset>> = {
      'ETH': {
        symbol: 'ETH',
        name: 'Ethereum',
        decimals: 18,
        priceUSD: 1800,
        ltv: 80,
        liquidationThreshold: 85,
        rates: {
          supplyAPY: 3.2,
          variableBorrowAPY: 4.5,
          stableBorrowAPY: 5.2,
          liquidityRate: 3.2,
          variableBorrowRate: 4.5,
          stableBorrowRate: 5.2,
          averageStableBorrowRate: 5.2,
          incentiveAPY: 0.8,
          totalAPY: 4.0
        }
      },
      'USDC': {
        symbol: 'USDC',
        name: 'USD Coin',
        decimals: 6,
        priceUSD: 1,
        ltv: 85,
        liquidationThreshold: 90,
        rates: {
          supplyAPY: 2.8,
          variableBorrowAPY: 4.5,
          stableBorrowAPY: 5.0,
          liquidityRate: 2.8,
          variableBorrowRate: 4.5,
          stableBorrowRate: 5.0,
          averageStableBorrowRate: 5.0,
          incentiveAPY: 0.5,
          totalAPY: 3.3
        }
      }
    }

    const baseAsset = mockAssets[symbol] || mockAssets['ETH']
    
    return {
      address: '0x0000000000000000000000000000000000000000',
      aTokenAddress: '0x0000000000000000000000000000000000000001',
      stableDebtTokenAddress: '0x0000000000000000000000000000000000000002',
      variableDebtTokenAddress: '0x0000000000000000000000000000000000000003',
      interestRateStrategyAddress: '0x0000000000000000000000000000000000000004',
      chainId: 1,
      isActive: true,
      isFrozen: false,
      isPaused: false,
      canBeCollateral: true,
      canBeBorrowed: true,
      canBeFlashLoaned: true,
      liquidationBonus: 5,
      reserveFactor: 10,
      usageAsCollateralEnabled: true,
      borrowingEnabled: true,
      stableBorrowRateEnabled: true,
      flashLoanEnabled: true,
      caps: {
        supplyCap: '1000000',
        borrowCap: '800000',
        debtCeiling: '50000000'
      },
      marketSize: '500000000',
      totalBorrowed: '300000000',
      availableLiquidity: '200000000',
      utilizationRate: 60,
      ...baseAsset
    } as LendingAsset
  }

  /**
   * Execute lending transaction
   */
  async executeLendingTransaction(
    protocolId: string,
    type: TransactionType,
    asset: LendingAsset,
    amount: string,
    userAddress: Address,
    metadata?: TransactionMetadata
  ): Promise<LendingTransaction> {
    const protocol = this.protocols.get(protocolId)
    if (!protocol) {
      throw new Error(`Protocol not found: ${protocolId}`)
    }

    const transaction: LendingTransaction = {
      id: `tx_${Date.now()}_${Math.random().toString(36).substring(2, 11)}`,
      protocolId,
      userAddress,
      type,
      asset,
      amount,
      amountUSD: parseFloat(amount) * asset.priceUSD,
      gasLimit: this.estimateGas(type, protocolId),
      gasPrice: '20000000000', // 20 gwei
      status: 'pending',
      timestamp: Date.now(),
      metadata
    }

    this.transactions.set(transaction.id, transaction)

    try {
      // Execute the transaction (mock implementation)
      const result = await this.performLendingTransaction(transaction)
      
      transaction.status = 'confirmed'
      transaction.hash = result.hash
      transaction.blockNumber = result.blockNumber
      transaction.gasUsed = result.gasUsed

      // Emit success event
      this.emitEvent({
        type: 'transaction_success',
        transaction,
        timestamp: Date.now()
      })

    } catch (error) {
      transaction.status = 'failed'
      transaction.error = (error as Error).message

      // Emit failure event
      this.emitEvent({
        type: 'transaction_failed',
        transaction,
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }

    return transaction
  }

  /**
   * Estimate gas for transaction type
   */
  private estimateGas(type: TransactionType, protocolId: string): string {
    const baseGas: Record<TransactionType, number> = {
      [TransactionType.SUPPLY]: 150000,
      [TransactionType.WITHDRAW]: 180000,
      [TransactionType.BORROW]: 200000,
      [TransactionType.REPAY]: 160000,
      [TransactionType.SWAP_BORROW_RATE]: 120000,
      [TransactionType.SET_USE_AS_COLLATERAL]: 100000,
      [TransactionType.LIQUIDATION_CALL]: 300000,
      [TransactionType.FLASH_LOAN]: 400000
    }

    const protocolMultiplier = protocolId === 'aave-v3' ? 1.0 : 1.2
    return Math.floor(baseGas[type] * protocolMultiplier * this.config.gasMultiplier).toString()
  }

  /**
   * Perform lending transaction (mock implementation)
   */
  private async performLendingTransaction(transaction: LendingTransaction): Promise<{
    hash: Hash
    blockNumber: number
    gasUsed: string
  }> {
    // Simulate network delay
    await new Promise(resolve => setTimeout(resolve, 2000 + Math.random() * 3000))

    // Simulate 95% success rate
    if (Math.random() < 0.95) {
      return {
        hash: `0x${Math.random().toString(16).substring(2, 66)}` as Hash,
        blockNumber: Math.floor(Math.random() * 1000000) + 18000000,
        gasUsed: (parseInt(transaction.gasLimit) * (0.8 + Math.random() * 0.2)).toString()
      }
    } else {
      throw new Error('Transaction failed: Insufficient collateral')
    }
  }

  /**
   * Get yield opportunities
   */
  async getYieldOpportunities(minAPY?: number, maxRisk?: number): Promise<YieldOpportunity[]> {
    const opportunities = Array.from(this.opportunities.values())
    
    return opportunities.filter(opp => {
      if (minAPY && opp.currentAPY < minAPY) return false
      if (maxRisk && opp.riskScore > maxRisk) return false
      return true
    }).sort((a, b) => b.currentAPY - a.currentAPY)
  }

  /**
   * Get liquidation opportunities
   */
  async getLiquidationOpportunities(minProfit?: number): Promise<LiquidationOpportunity[]> {
    const opportunities = Array.from(this.liquidations.values())
    
    return opportunities.filter(opp => {
      if (minProfit && opp.profitability < minProfit) return false
      return opp.healthFactor < 1.0
    }).sort((a, b) => b.profitability - a.profitability)
  }

  /**
   * Get supported assets
   */
  getSupportedAssets(protocolId?: string): LendingAsset[] {
    const assets = Array.from(this.assets.values())
    return protocolId ? assets.filter(asset => asset.chainId === 1) : assets // Simplified filtering
  }

  /**
   * Get protocols
   */
  getProtocols(chainId?: number): LendingProtocol[] {
    const protocols = Array.from(this.protocols.values())
    return chainId ? protocols.filter(protocol => protocol.chainId === chainId) : protocols
  }

  /**
   * Get transaction
   */
  getTransaction(id: string): LendingTransaction | null {
    return this.transactions.get(id) || null
  }

  /**
   * Update configuration
   */
  updateConfig(config: Partial<LendingConfig>): void {
    this.config = { ...this.config, ...config }
  }

  /**
   * Get configuration
   */
  getConfig(): LendingConfig {
    return { ...this.config }
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: LendingEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in lending event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: LendingEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.positions.clear()
    this.transactions.clear()
    this.opportunities.clear()
    this.liquidations.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface LendingEvent {
  type: 'transaction_success' | 'transaction_failed' | 'position_updated' | 'liquidation_risk' | 'yield_opportunity'
  transaction?: LendingTransaction
  position?: UserPosition
  opportunity?: YieldOpportunity
  liquidation?: LiquidationOpportunity
  error?: Error
  timestamp: number
}

// Export singleton instance
export const lendingIntegration = LendingIntegration.getInstance()
