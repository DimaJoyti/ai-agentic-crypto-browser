import { type Address } from 'viem'
import { defiProtocolService } from './defi-protocols'

export enum PositionType {
  LENDING = 'lending',
  BORROWING = 'borrowing',
  LIQUIDITY = 'liquidity',
  STAKING = 'staking',
  YIELD_FARMING = 'yield_farming',
  DERIVATIVES = 'derivatives'
}

export enum PositionStatus {
  ACTIVE = 'active',
  INACTIVE = 'inactive',
  AT_RISK = 'at_risk',
  LIQUIDATABLE = 'liquidatable',
  EXPIRED = 'expired'
}

export interface TokenPosition {
  address: Address
  symbol: string
  name: string
  decimals: number
  amount: string
  value: string
  price: string
  change24h: string
}

export interface BasePosition {
  id: string
  type: PositionType
  protocol: string
  chainId: number
  status: PositionStatus
  createdAt: number
  updatedAt: number
  userAddress: Address
  totalValue: string
  pnl: string
  pnlPercentage: string
  apy?: string
  rewards?: TokenPosition[]
}

export interface LendingPosition extends BasePosition {
  type: PositionType.LENDING
  asset: TokenPosition
  supplied: string
  earned: string
  supplyApy: string
  utilizationRate: string
  collateralFactor: string
}

export interface BorrowingPosition extends BasePosition {
  type: PositionType.BORROWING
  asset: TokenPosition
  borrowed: string
  debt: string
  borrowApy: string
  healthFactor: string
  liquidationThreshold: string
  collateral: TokenPosition[]
}

export interface LiquidityPosition extends BasePosition {
  type: PositionType.LIQUIDITY
  pool: {
    address: Address
    token0: TokenPosition
    token1: TokenPosition
    fee: number
    priceRange: {
      min: string
      max: string
      current: string
    }
  }
  liquidity: string
  fees24h: string
  feesTotal: string
  impermanentLoss: string
  inRange: boolean
}

export interface StakingPosition extends BasePosition {
  type: PositionType.STAKING
  asset: TokenPosition
  staked: string
  rewards: TokenPosition[]
  lockPeriod?: number
  unlockDate?: number
  slashingRisk: string
  validatorAddress?: Address
}

export interface YieldFarmingPosition extends BasePosition {
  type: PositionType.YIELD_FARMING
  pool: {
    address: Address
    tokens: TokenPosition[]
    lpToken: TokenPosition
  }
  staked: string
  rewards: TokenPosition[]
  multiplier: string
  lockPeriod?: number
  harvestable: string
}

export type DeFiPosition = 
  | LendingPosition 
  | BorrowingPosition 
  | LiquidityPosition 
  | StakingPosition 
  | YieldFarmingPosition

export interface PositionAlert {
  id: string
  positionId: string
  type: 'health_factor' | 'liquidation' | 'price_range' | 'unlock' | 'harvest'
  severity: 'low' | 'medium' | 'high' | 'critical'
  title: string
  message: string
  threshold: string
  currentValue: string
  createdAt: number
  acknowledged: boolean
}

export interface PositionSummary {
  totalPositions: number
  totalValue: string
  totalPnl: string
  totalPnlPercentage: string
  totalRewards: string
  activeAlerts: number
  positionsByType: Record<PositionType, number>
  positionsByProtocol: Record<string, number>
  positionsByChain: Record<number, number>
}

export class PositionManager {
  private static instance: PositionManager
  private positions: Map<string, DeFiPosition> = new Map()
  private alerts: Map<string, PositionAlert> = new Map()
  private updateIntervals: Map<string, NodeJS.Timeout> = new Map()

  private constructor() {
    this.initializeMockPositions()
    this.startPositionTracking()
  }

  static getInstance(): PositionManager {
    if (!PositionManager.instance) {
      PositionManager.instance = new PositionManager()
    }
    return PositionManager.instance
  }

  private initializeMockPositions() {
    // Mock lending position
    const lendingPosition: LendingPosition = {
      id: 'lending-1',
      type: PositionType.LENDING,
      protocol: 'aave-v3',
      chainId: 1,
      status: PositionStatus.ACTIVE,
      createdAt: Date.now() - 86400000 * 30, // 30 days ago
      updatedAt: Date.now(),
      userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      totalValue: '10000',
      pnl: '320.50',
      pnlPercentage: '3.21',
      apy: '3.2%',
      asset: {
        address: '0xA0b86a33E6441b8435b662f0E2d0c2837b0b3c0' as Address,
        symbol: 'USDC',
        name: 'USD Coin',
        decimals: 6,
        amount: '10000',
        value: '10000',
        price: '1.00',
        change24h: '0.01'
      },
      supplied: '10000',
      earned: '320.50',
      supplyApy: '3.2%',
      utilizationRate: '75%',
      collateralFactor: '80%'
    }

    // Mock borrowing position
    const borrowingPosition: BorrowingPosition = {
      id: 'borrowing-1',
      type: PositionType.BORROWING,
      protocol: 'aave-v3',
      chainId: 1,
      status: PositionStatus.ACTIVE,
      createdAt: Date.now() - 86400000 * 15, // 15 days ago
      updatedAt: Date.now(),
      userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      totalValue: '5000',
      pnl: '-75.25',
      pnlPercentage: '-1.51',
      asset: {
        address: '0x6B175474E89094C44Da98b954EedeAC495271d0F' as Address,
        symbol: 'DAI',
        name: 'Dai Stablecoin',
        decimals: 18,
        amount: '5000',
        value: '5000',
        price: '1.00',
        change24h: '-0.02'
      },
      borrowed: '5000',
      debt: '5075.25',
      borrowApy: '4.8%',
      healthFactor: '2.1',
      liquidationThreshold: '85%',
      collateral: [
        {
          address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' as Address,
          symbol: 'WETH',
          name: 'Wrapped Ether',
          decimals: 18,
          amount: '4.2',
          value: '10080',
          price: '2400',
          change24h: '2.5'
        }
      ]
    }

    // Mock liquidity position
    const liquidityPosition: LiquidityPosition = {
      id: 'liquidity-1',
      type: PositionType.LIQUIDITY,
      protocol: 'uniswap-v3',
      chainId: 1,
      status: PositionStatus.ACTIVE,
      createdAt: Date.now() - 86400000 * 7, // 7 days ago
      updatedAt: Date.now(),
      userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      totalValue: '2500',
      pnl: '125.75',
      pnlPercentage: '5.03',
      apy: '8.5%',
      pool: {
        address: '0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640' as Address,
        token0: {
          address: '0xA0b86a33E6441b8435b662f0E2d0c2837b0b3c0' as Address,
          symbol: 'USDC',
          name: 'USD Coin',
          decimals: 6,
          amount: '1250',
          value: '1250',
          price: '1.00',
          change24h: '0.01'
        },
        token1: {
          address: '0xC02aaA39b223FE8D0A0e5C4F27eAD9083C756Cc2' as Address,
          symbol: 'WETH',
          name: 'Wrapped Ether',
          decimals: 18,
          amount: '0.52',
          value: '1248',
          price: '2400',
          change24h: '2.5'
        },
        fee: 500, // 0.05%
        priceRange: {
          min: '2350',
          max: '2450',
          current: '2400'
        }
      },
      liquidity: '2500',
      fees24h: '12.50',
      feesTotal: '87.50',
      impermanentLoss: '-2.3%',
      inRange: true
    }

    // Mock staking position
    const stakingPosition: StakingPosition = {
      id: 'staking-1',
      type: PositionType.STAKING,
      protocol: 'lido',
      chainId: 1,
      status: PositionStatus.ACTIVE,
      createdAt: Date.now() - 86400000 * 60, // 60 days ago
      updatedAt: Date.now(),
      userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
      totalValue: '12000',
      pnl: '456.80',
      pnlPercentage: '3.81',
      apy: '3.8%',
      asset: {
        address: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
        symbol: 'stETH',
        name: 'Liquid staked Ether 2.0',
        decimals: 18,
        amount: '5.02',
        value: '12048',
        price: '2400',
        change24h: '2.5'
      },
      staked: '5.0',
      rewards: [
        {
          address: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
          symbol: 'stETH',
          name: 'Liquid staked Ether 2.0',
          decimals: 18,
          amount: '0.02',
          value: '48',
          price: '2400',
          change24h: '2.5'
        }
      ],
      slashingRisk: 'Low'
    }

    this.positions.set(lendingPosition.id, lendingPosition)
    this.positions.set(borrowingPosition.id, borrowingPosition)
    this.positions.set(liquidityPosition.id, liquidityPosition)
    this.positions.set(stakingPosition.id, stakingPosition)

    // Generate alerts
    this.generateAlerts()
  }

  private generateAlerts() {
    // Health factor alert for borrowing position
    const healthAlert: PositionAlert = {
      id: 'alert-1',
      positionId: 'borrowing-1',
      type: 'health_factor',
      severity: 'medium',
      title: 'Health Factor Warning',
      message: 'Your health factor is approaching the liquidation threshold',
      threshold: '1.5',
      currentValue: '2.1',
      createdAt: Date.now() - 3600000, // 1 hour ago
      acknowledged: false
    }

    // Price range alert for liquidity position
    const rangeAlert: PositionAlert = {
      id: 'alert-2',
      positionId: 'liquidity-1',
      type: 'price_range',
      severity: 'low',
      title: 'Price Range Update',
      message: 'Your liquidity position is still in range but approaching the edge',
      threshold: '2450',
      currentValue: '2400',
      createdAt: Date.now() - 1800000, // 30 minutes ago
      acknowledged: false
    }

    this.alerts.set(healthAlert.id, healthAlert)
    this.alerts.set(rangeAlert.id, rangeAlert)
  }

  private startPositionTracking() {
    // Update positions every 30 seconds
    const interval = setInterval(() => {
      this.updatePositions()
    }, 30000)
    
    this.updateIntervals.set('main', interval)
  }

  private updatePositions() {
    // Update position values and check for alerts
    this.positions.forEach(position => {
      position.updatedAt = Date.now()
      
      // Update PnL based on price changes
      this.updatePositionPnL(position)
      
      // Check for new alerts
      this.checkPositionAlerts(position)
    })
  }

  private updatePositionPnL(position: DeFiPosition) {
    // Mock PnL updates - in real app, would calculate based on current prices
    const randomChange = (Math.random() - 0.5) * 0.1 // ±5% random change
    const currentPnl = parseFloat(position.pnl)
    const newPnl = currentPnl * (1 + randomChange)
    
    position.pnl = newPnl.toFixed(2)
    position.pnlPercentage = ((newPnl / parseFloat(position.totalValue)) * 100).toFixed(2)
  }

  private checkPositionAlerts(position: DeFiPosition) {
    if (position.type === PositionType.BORROWING) {
      const borrowPosition = position as BorrowingPosition
      const healthFactor = parseFloat(borrowPosition.healthFactor)
      
      if (healthFactor < 1.5 && healthFactor > 1.2) {
        position.status = PositionStatus.AT_RISK
      } else if (healthFactor <= 1.2) {
        position.status = PositionStatus.LIQUIDATABLE
      } else {
        position.status = PositionStatus.ACTIVE
      }
    }
  }

  // Public methods
  getAllPositions(userAddress?: Address): DeFiPosition[] {
    const positions = Array.from(this.positions.values())
    return userAddress 
      ? positions.filter(p => p.userAddress.toLowerCase() === userAddress.toLowerCase())
      : positions
  }

  getPositionsByType(type: PositionType, userAddress?: Address): DeFiPosition[] {
    return this.getAllPositions(userAddress).filter(p => p.type === type)
  }

  getPositionsByProtocol(protocol: string, userAddress?: Address): DeFiPosition[] {
    return this.getAllPositions(userAddress).filter(p => p.protocol === protocol)
  }

  getPositionsByChain(chainId: number, userAddress?: Address): DeFiPosition[] {
    return this.getAllPositions(userAddress).filter(p => p.chainId === chainId)
  }

  getPosition(id: string): DeFiPosition | undefined {
    return this.positions.get(id)
  }

  getPositionSummary(userAddress?: Address): PositionSummary {
    const positions = this.getAllPositions(userAddress)
    
    const totalValue = positions.reduce((sum, p) => sum + parseFloat(p.totalValue), 0)
    const totalPnl = positions.reduce((sum, p) => sum + parseFloat(p.pnl), 0)
    const totalRewards = positions.reduce((sum, p) => {
      if (p.rewards) {
        return sum + p.rewards.reduce((rewardSum, r) => rewardSum + parseFloat(r.value), 0)
      }
      return sum
    }, 0)

    const positionsByType = positions.reduce((acc, p) => {
      acc[p.type] = (acc[p.type] || 0) + 1
      return acc
    }, {} as Record<PositionType, number>)

    const positionsByProtocol = positions.reduce((acc, p) => {
      acc[p.protocol] = (acc[p.protocol] || 0) + 1
      return acc
    }, {} as Record<string, number>)

    const positionsByChain = positions.reduce((acc, p) => {
      acc[p.chainId] = (acc[p.chainId] || 0) + 1
      return acc
    }, {} as Record<number, number>)

    return {
      totalPositions: positions.length,
      totalValue: totalValue.toFixed(2),
      totalPnl: totalPnl.toFixed(2),
      totalPnlPercentage: ((totalPnl / totalValue) * 100).toFixed(2),
      totalRewards: totalRewards.toFixed(2),
      activeAlerts: this.getActiveAlerts().length,
      positionsByType,
      positionsByProtocol,
      positionsByChain
    }
  }

  getActiveAlerts(positionId?: string): PositionAlert[] {
    const alerts = Array.from(this.alerts.values()).filter(a => !a.acknowledged)
    return positionId ? alerts.filter(a => a.positionId === positionId) : alerts
  }

  acknowledgeAlert(alertId: string): void {
    const alert = this.alerts.get(alertId)
    if (alert) {
      alert.acknowledged = true
    }
  }

  destroy(): void {
    this.updateIntervals.forEach(interval => clearInterval(interval))
    this.updateIntervals.clear()
  }
}

// Export singleton instance
export const positionManager = PositionManager.getInstance()
