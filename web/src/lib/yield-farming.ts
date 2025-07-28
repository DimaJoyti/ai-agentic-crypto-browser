import { createPublicClient, http, type Address, type Hash } from 'viem'
import { SUPPORTED_CHAINS } from './chains'

export enum YieldStrategy {
  SINGLE_STAKING = 'single_staking',
  LIQUIDITY_MINING = 'liquidity_mining',
  YIELD_FARMING = 'yield_farming',
  LIQUID_STAKING = 'liquid_staking',
  LENDING_YIELD = 'lending_yield',
  AUTOCOMPOUNDING = 'autocompounding'
}

export enum RiskLevel {
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  VERY_HIGH = 'very_high'
}

export interface YieldToken {
  address: Address
  symbol: string
  name: string
  decimals: number
  price: string
  logoURI?: string
}

export interface RewardToken extends YieldToken {
  amount: string
}

export interface YieldFarm {
  id: string
  protocol: string
  name: string
  strategy: YieldStrategy
  chainId: number
  contractAddress: Address
  stakingToken: YieldToken
  rewardTokens: YieldToken[]
  apy: string
  tvl: string
  multiplier?: string
  lockPeriod?: number // in seconds
  minimumStake: string
  maximumStake?: string
  riskLevel: RiskLevel
  isActive: boolean
  endDate?: number
  features: string[]
  description: string
}

export interface StakingPool {
  id: string
  protocol: string
  name: string
  chainId: number
  contractAddress: Address
  stakingToken: YieldToken
  rewardToken: YieldToken
  apy: string
  tvl: string
  lockPeriod?: number
  slashingRisk: boolean
  minimumStake: string
  validatorCount?: number
  riskLevel: RiskLevel
  isActive: boolean
  features: string[]
}

export interface UserYieldPosition {
  id: string
  farmId: string
  userAddress: Address
  stakedAmount: string
  pendingRewards: RewardToken[]
  claimableRewards: RewardToken[]
  multiplier: string
  startTime: number
  lastClaimTime: number
  lockEndTime?: number
  autoCompound: boolean
  totalEarned: string
  currentValue: string
  pnl: string
  pnlPercentage: string
}

export interface YieldOptimization {
  currentFarm: YieldFarm
  suggestedFarms: YieldFarm[]
  potentialGains: string
  migrationCost: string
  recommendation: 'stay' | 'migrate' | 'diversify'
  reason: string
}

export class YieldFarmingService {
  private static instance: YieldFarmingService
  private clients: Map<number, any> = new Map()
  private farms: Map<string, YieldFarm> = new Map()
  private stakingPools: Map<string, StakingPool> = new Map()
  private userPositions: Map<string, UserYieldPosition[]> = new Map()

  private constructor() {
    this.initializeClients()
    this.initializeFarms()
    this.initializeStakingPools()
    this.initializeMockPositions()
  }

  static getInstance(): YieldFarmingService {
    if (!YieldFarmingService.instance) {
      YieldFarmingService.instance = new YieldFarmingService()
    }
    return YieldFarmingService.instance
  }

  private initializeClients() {
    Object.values(SUPPORTED_CHAINS).forEach(chain => {
      if (!chain.isTestnet || chain.id === 11155111) {
        try {
          const client = createPublicClient({
            chain: {
              id: chain.id,
              name: chain.name,
              network: chain.shortName.toLowerCase(),
              nativeCurrency: chain.nativeCurrency,
              rpcUrls: chain.rpcUrls
            } as any,
            transport: http()
          })
          this.clients.set(chain.id, client)
        } catch (error) {
          console.warn(`Failed to initialize yield farming client for chain ${chain.id}:`, error)
        }
      }
    })
  }

  private initializeFarms() {
    // Uniswap V3 Liquidity Mining
    this.farms.set('uniswap-usdc-eth', {
      id: 'uniswap-usdc-eth',
      protocol: 'Uniswap V3',
      name: 'USDC/ETH Liquidity Pool',
      strategy: YieldStrategy.LIQUIDITY_MINING,
      chainId: 1,
      contractAddress: '0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640' as Address,
      stakingToken: {
        address: '0x88e6A0c2dDD26FEEb64F039a2c41296FcB3f5640' as Address,
        symbol: 'UNI-V3-USDC-ETH',
        name: 'Uniswap V3 USDC/ETH LP',
        decimals: 18,
        price: '1.0'
      },
      rewardTokens: [
        {
          address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984' as Address,
          symbol: 'UNI',
          name: 'Uniswap',
          decimals: 18,
          price: '8.50'
        }
      ],
      apy: '12.5%',
      tvl: '$125M',
      multiplier: '1.0x',
      minimumStake: '0.01',
      riskLevel: RiskLevel.MEDIUM,
      isActive: true,
      features: ['Concentrated Liquidity', 'Fee Rewards', 'UNI Rewards'],
      description: 'Provide liquidity to the USDC/ETH pool and earn trading fees plus UNI rewards'
    })

    // Compound Yield Farming
    this.farms.set('compound-usdc', {
      id: 'compound-usdc',
      protocol: 'Compound V3',
      name: 'USDC Supply',
      strategy: YieldStrategy.LENDING_YIELD,
      chainId: 1,
      contractAddress: '0xc3d688B66703497DAA19211EEdff47f25384cdc3' as Address,
      stakingToken: {
        address: '0xA0b86a33E6441b8435b662f0E2d0c2837b0b3c0' as Address,
        symbol: 'USDC',
        name: 'USD Coin',
        decimals: 6,
        price: '1.00'
      },
      rewardTokens: [
        {
          address: '0xc00e94Cb662C3520282E6f5717214004A7f26888' as Address,
          symbol: 'COMP',
          name: 'Compound',
          decimals: 18,
          price: '45.20'
        }
      ],
      apy: '4.8%',
      tvl: '$85M',
      minimumStake: '1',
      riskLevel: RiskLevel.LOW,
      isActive: true,
      features: ['Lending Rewards', 'COMP Rewards', 'Auto-compound'],
      description: 'Supply USDC to Compound and earn lending interest plus COMP rewards'
    })

    // Curve Yield Farming
    this.farms.set('curve-3pool', {
      id: 'curve-3pool',
      protocol: 'Curve Finance',
      name: '3Pool (DAI/USDC/USDT)',
      strategy: YieldStrategy.YIELD_FARMING,
      chainId: 1,
      contractAddress: '0xbEbc44782C7dB0a1A60Cb6fe97d0b483032FF1C7' as Address,
      stakingToken: {
        address: '0x6c3F90f043a72FA612cbac8115EE7e52BDe6E490' as Address,
        symbol: '3Crv',
        name: 'Curve.fi DAI/USDC/USDT',
        decimals: 18,
        price: '1.02'
      },
      rewardTokens: [
        {
          address: '0xD533a949740bb3306d119CC777fa900bA034cd52' as Address,
          symbol: 'CRV',
          name: 'Curve DAO Token',
          decimals: 18,
          price: '0.85'
        }
      ],
      apy: '8.2%',
      tvl: '$180M',
      multiplier: '1.5x',
      minimumStake: '10',
      riskLevel: RiskLevel.LOW,
      isActive: true,
      features: ['Stable Yield', 'CRV Rewards', 'Low Impermanent Loss'],
      description: 'Stake 3Pool LP tokens to earn CRV rewards with minimal impermanent loss'
    })

    // Convex Autocompounding
    this.farms.set('convex-3pool', {
      id: 'convex-3pool',
      protocol: 'Convex Finance',
      name: 'Curve 3Pool Autocompound',
      strategy: YieldStrategy.AUTOCOMPOUNDING,
      chainId: 1,
      contractAddress: '0x689440f2Ff927E1f24c72F1087E1FAF471eCe1c8' as Address,
      stakingToken: {
        address: '0x6c3F90f043a72FA612cbac8115EE7e52BDe6E490' as Address,
        symbol: '3Crv',
        name: 'Curve.fi DAI/USDC/USDT',
        decimals: 18,
        price: '1.02'
      },
      rewardTokens: [
        {
          address: '0xD533a949740bb3306d119CC777fa900bA034cd52' as Address,
          symbol: 'CRV',
          name: 'Curve DAO Token',
          decimals: 18,
          price: '0.85'
        },
        {
          address: '0x4e3FBD56CD56c3e72c1403e103b45Db9da5B9D2B' as Address,
          symbol: 'CVX',
          name: 'Convex Token',
          decimals: 18,
          price: '2.15'
        }
      ],
      apy: '15.8%',
      tvl: '$95M',
      multiplier: '2.5x',
      minimumStake: '10',
      riskLevel: RiskLevel.MEDIUM,
      isActive: true,
      features: ['Auto-compound', 'Boosted Rewards', 'CRV + CVX'],
      description: 'Auto-compound Curve 3Pool rewards with boosted CRV and CVX rewards'
    })
  }

  private initializeStakingPools() {
    // Ethereum 2.0 Staking via Lido
    this.stakingPools.set('lido-eth', {
      id: 'lido-eth',
      protocol: 'Lido',
      name: 'Ethereum 2.0 Staking',
      chainId: 1,
      contractAddress: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
      stakingToken: {
        address: '0x0000000000000000000000000000000000000000' as Address,
        symbol: 'ETH',
        name: 'Ethereum',
        decimals: 18,
        price: '2400'
      },
      rewardToken: {
        address: '0xae7ab96520DE3A18E5e111B5EaAb095312D7fE84' as Address,
        symbol: 'stETH',
        name: 'Liquid staked Ether 2.0',
        decimals: 18,
        price: '2400'
      },
      apy: '3.8%',
      tvl: '$14.2B',
      slashingRisk: true,
      minimumStake: '0.01',
      validatorCount: 8500,
      riskLevel: RiskLevel.LOW,
      isActive: true,
      features: ['Liquid Staking', 'Daily Rewards', 'No Lock Period']
    })

    // Rocket Pool ETH Staking
    this.stakingPools.set('rocketpool-eth', {
      id: 'rocketpool-eth',
      protocol: 'Rocket Pool',
      name: 'Decentralized ETH Staking',
      chainId: 1,
      contractAddress: '0xae78736Cd615f374D3085123A210448E74Fc6393' as Address,
      stakingToken: {
        address: '0x0000000000000000000000000000000000000000' as Address,
        symbol: 'ETH',
        name: 'Ethereum',
        decimals: 18,
        price: '2400'
      },
      rewardToken: {
        address: '0xae78736Cd615f374D3085123A210448E74Fc6393' as Address,
        symbol: 'rETH',
        name: 'Rocket Pool ETH',
        decimals: 18,
        price: '2420'
      },
      apy: '4.1%',
      tvl: '$2.8B',
      slashingRisk: true,
      minimumStake: '0.01',
      validatorCount: 2100,
      riskLevel: RiskLevel.MEDIUM,
      isActive: true,
      features: ['Decentralized', 'Node Operator Rewards', 'Premium to ETH']
    })
  }

  private initializeMockPositions() {
    const mockPositions: UserYieldPosition[] = [
      {
        id: 'pos-1',
        farmId: 'uniswap-usdc-eth',
        userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        stakedAmount: '2500',
        pendingRewards: [
          {
            address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984' as Address,
            symbol: 'UNI',
            name: 'Uniswap',
            decimals: 18,
            price: '8.50',
            amount: '12.5'
          }
        ],
        claimableRewards: [
          {
            address: '0x1f9840a85d5aF5bf1D1762F925BDADdC4201F984' as Address,
            symbol: 'UNI',
            name: 'Uniswap',
            decimals: 18,
            price: '8.50',
            amount: '8.2'
          }
        ],
        multiplier: '1.0',
        startTime: Date.now() - 86400000 * 30, // 30 days ago
        lastClaimTime: Date.now() - 86400000 * 7, // 7 days ago
        autoCompound: false,
        totalEarned: '156.80',
        currentValue: '2656.80',
        pnl: '156.80',
        pnlPercentage: '6.27'
      },
      {
        id: 'pos-2',
        farmId: 'convex-3pool',
        userAddress: '0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1' as Address,
        stakedAmount: '5000',
        pendingRewards: [
          {
            address: '0xD533a949740bb3306d119CC777fa900bA034cd52' as Address,
            symbol: 'CRV',
            name: 'Curve DAO Token',
            decimals: 18,
            price: '0.85',
            amount: '45.2'
          },
          {
            address: '0x4e3FBD56CD56c3e72c1403e103b45Db9da5B9D2B' as Address,
            symbol: 'CVX',
            name: 'Convex Token',
            decimals: 18,
            price: '2.15',
            amount: '18.5'
          }
        ],
        claimableRewards: [],
        multiplier: '2.5',
        startTime: Date.now() - 86400000 * 45, // 45 days ago
        lastClaimTime: Date.now() - 86400000 * 1, // 1 day ago (auto-compound)
        autoCompound: true,
        totalEarned: '425.60',
        currentValue: '5425.60',
        pnl: '425.60',
        pnlPercentage: '8.51'
      }
    ]

    this.userPositions.set('0x742d35Cc6634C0532925a3b8D4C9db96C4b4d8b1', mockPositions)
  }

  // Public methods
  getAllFarms(chainId?: number): YieldFarm[] {
    const farms = Array.from(this.farms.values())
    return chainId ? farms.filter(farm => farm.chainId === chainId) : farms
  }

  getFarmsByStrategy(strategy: YieldStrategy, chainId?: number): YieldFarm[] {
    return this.getAllFarms(chainId).filter(farm => farm.strategy === strategy)
  }

  getFarmsByRisk(riskLevel: RiskLevel, chainId?: number): YieldFarm[] {
    return this.getAllFarms(chainId).filter(farm => farm.riskLevel === riskLevel)
  }

  getTopFarmsByAPY(limit: number = 10, chainId?: number): YieldFarm[] {
    return this.getAllFarms(chainId)
      .sort((a, b) => parseFloat(b.apy) - parseFloat(a.apy))
      .slice(0, limit)
  }

  getAllStakingPools(chainId?: number): StakingPool[] {
    const pools = Array.from(this.stakingPools.values())
    return chainId ? pools.filter(pool => pool.chainId === chainId) : pools
  }

  getUserPositions(userAddress: Address): UserYieldPosition[] {
    return this.userPositions.get(userAddress.toLowerCase()) || []
  }

  async stakeFarm(
    farmId: string,
    amount: string,
    userAddress: Address
  ): Promise<Hash> {
    const farm = this.farms.get(farmId)
    if (!farm) {
      throw new Error('Farm not found')
    }

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    
    // Add position to user positions
    const userPositions = this.getUserPositions(userAddress)
    const newPosition: UserYieldPosition = {
      id: `pos-${Date.now()}`,
      farmId,
      userAddress,
      stakedAmount: amount,
      pendingRewards: [],
      claimableRewards: [],
      multiplier: farm.multiplier || '1.0',
      startTime: Date.now(),
      lastClaimTime: Date.now(),
      autoCompound: false,
      totalEarned: '0',
      currentValue: amount,
      pnl: '0',
      pnlPercentage: '0'
    }

    userPositions.push(newPosition)
    this.userPositions.set(userAddress.toLowerCase(), userPositions)

    return mockTxHash as Hash
  }

  async unstakeFarm(
    positionId: string,
    amount: string,
    userAddress: Address
  ): Promise<Hash> {
    const userPositions = this.getUserPositions(userAddress)
    const positionIndex = userPositions.findIndex(pos => pos.id === positionId)
    
    if (positionIndex === -1) {
      throw new Error('Position not found')
    }

    const position = userPositions[positionIndex]
    const stakedAmount = parseFloat(position.stakedAmount)
    const unstakeAmount = parseFloat(amount)

    if (unstakeAmount > stakedAmount) {
      throw new Error('Insufficient staked amount')
    }

    // Update position
    if (unstakeAmount === stakedAmount) {
      // Remove position completely
      userPositions.splice(positionIndex, 1)
    } else {
      // Reduce staked amount
      position.stakedAmount = (stakedAmount - unstakeAmount).toString()
      position.currentValue = (parseFloat(position.currentValue) - unstakeAmount).toString()
    }

    this.userPositions.set(userAddress.toLowerCase(), userPositions)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  async claimRewards(
    positionId: string,
    userAddress: Address
  ): Promise<Hash> {
    const userPositions = this.getUserPositions(userAddress)
    const position = userPositions.find(pos => pos.id === positionId)
    
    if (!position) {
      throw new Error('Position not found')
    }

    // Move claimable rewards to total earned
    const claimableValue = position.claimableRewards.reduce((sum, reward) => {
      return sum + (parseFloat(reward.amount) * parseFloat(reward.price))
    }, 0)

    position.totalEarned = (parseFloat(position.totalEarned) + claimableValue).toString()
    position.claimableRewards = []
    position.lastClaimTime = Date.now()

    this.userPositions.set(userAddress.toLowerCase(), userPositions)

    // Mock transaction hash
    const mockTxHash = '0x' + Array(64).fill(0).map(() => Math.floor(Math.random() * 16).toString(16)).join('')
    return mockTxHash as Hash
  }

  getYieldOptimization(positionId: string, userAddress: Address): YieldOptimization | null {
    const userPositions = this.getUserPositions(userAddress)
    const position = userPositions.find(pos => pos.id === positionId)
    
    if (!position) return null

    const currentFarm = this.farms.get(position.farmId)
    if (!currentFarm) return null

    // Find better farms
    const allFarms = this.getAllFarms(currentFarm.chainId)
    const betterFarms = allFarms
      .filter(farm => 
        farm.id !== currentFarm.id && 
        parseFloat(farm.apy) > parseFloat(currentFarm.apy) &&
        farm.isActive
      )
      .sort((a, b) => parseFloat(b.apy) - parseFloat(a.apy))
      .slice(0, 3)

    if (betterFarms.length === 0) {
      return {
        currentFarm,
        suggestedFarms: [],
        potentialGains: '0',
        migrationCost: '0',
        recommendation: 'stay',
        reason: 'Current farm has the best available APY'
      }
    }

    const bestFarm = betterFarms[0]
    const currentAPY = parseFloat(currentFarm.apy)
    const bestAPY = parseFloat(bestFarm.apy)
    const potentialGains = ((bestAPY - currentAPY) / currentAPY * 100).toFixed(2)

    return {
      currentFarm,
      suggestedFarms: betterFarms,
      potentialGains: `${potentialGains}%`,
      migrationCost: '0.005 ETH',
      recommendation: parseFloat(potentialGains) > 2 ? 'migrate' : 'stay',
      reason: parseFloat(potentialGains) > 2 
        ? `Migration could increase APY by ${potentialGains}%`
        : 'Migration gains are too small to justify gas costs'
    }
  }
}

// Export singleton instance
export const yieldFarmingService = YieldFarmingService.getInstance()
