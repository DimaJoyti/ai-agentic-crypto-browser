import { type Address } from 'viem'

export interface PortfolioRiskAssessment {
  portfolioId: string
  userAddress: Address
  timestamp: string
  overallRisk: RiskLevel
  riskScore: number
  riskMetrics: RiskMetrics
  exposureAnalysis: ExposureAnalysis
  concentrationRisk: ConcentrationRisk
  liquidityRisk: LiquidityRisk
  marketRisk: MarketRisk
  counterpartyRisk: CounterpartyRisk
  technicalRisk: TechnicalRisk
  stressTests: StressTestResult[]
  riskMitigation: RiskMitigationPlan
  recommendations: RiskRecommendation[]
  alerts: RiskAlert[]
}

export enum RiskLevel {
  VERY_LOW = 'very_low',
  LOW = 'low',
  MEDIUM = 'medium',
  HIGH = 'high',
  CRITICAL = 'critical'
}

export interface RiskMetrics {
  valueAtRisk: VaRMetrics
  expectedShortfall: number
  sharpeRatio: number
  sortinoRatio: number
  maxDrawdown: number
  volatility: number
  beta: number
  correlation: CorrelationMatrix
  diversificationRatio: number
  riskAdjustedReturn: number
}

export interface VaRMetrics {
  var95: number
  var99: number
  timeHorizon: number
  confidence: number
  method: 'historical' | 'parametric' | 'monte_carlo'
}

export interface CorrelationMatrix {
  assets: string[]
  matrix: number[][]
  averageCorrelation: number
  maxCorrelation: number
  minCorrelation: number
}

export interface ExposureAnalysis {
  totalExposure: number
  exposureByAsset: AssetExposure[]
  exposureByProtocol: ProtocolExposure[]
  exposureByChain: ChainExposure[]
  exposureByCategory: CategoryExposure[]
  geographicExposure: GeographicExposure[]
  sectorExposure: SectorExposure[]
}

export interface AssetExposure {
  asset: string
  symbol: string
  contractAddress?: Address
  exposure: number
  percentage: number
  riskContribution: number
  volatility: number
  liquidity: number
}

export interface ProtocolExposure {
  protocol: string
  exposure: number
  percentage: number
  riskScore: number
  tvl: number
  auditStatus: string
  vulnerabilities: number
}

export interface ChainExposure {
  chainId: number
  chainName: string
  exposure: number
  percentage: number
  riskScore: number
  bridgeRisk: number
  validatorRisk: number
}

export interface CategoryExposure {
  category: string
  exposure: number
  percentage: number
  riskScore: number
  volatility: number
}

export interface GeographicExposure {
  region: string
  exposure: number
  percentage: number
  regulatoryRisk: number
  politicalRisk: number
}

export interface SectorExposure {
  sector: string
  exposure: number
  percentage: number
  riskScore: number
  correlation: number
}

export interface ConcentrationRisk {
  herfindahlIndex: number
  top5Concentration: number
  top10Concentration: number
  concentrationScore: number
  concentrationLevel: RiskLevel
  concentratedAssets: ConcentratedAsset[]
  diversificationGap: number
}

export interface ConcentratedAsset {
  asset: string
  percentage: number
  riskContribution: number
  recommendedReduction: number
}

export interface LiquidityRisk {
  liquidityScore: number
  liquidityLevel: RiskLevel
  illiquidAssets: IlliquidAsset[]
  liquidityBuffer: number
  liquidationTime: LiquidationTimeAnalysis
  marketImpact: MarketImpactAnalysis
}

export interface IlliquidAsset {
  asset: string
  liquidityScore: number
  dailyVolume: number
  bidAskSpread: number
  marketDepth: number
  liquidationTime: number
}

export interface LiquidationTimeAnalysis {
  immediate: number
  within1Hour: number
  within1Day: number
  within1Week: number
  moreThan1Week: number
}

export interface MarketImpactAnalysis {
  lowImpact: number
  mediumImpact: number
  highImpact: number
  averageSlippage: number
}

export interface MarketRisk {
  marketBeta: number
  marketCorrelation: number
  systematicRisk: number
  idiosyncraticRisk: number
  marketRegimes: MarketRegime[]
  cyclicalRisk: CyclicalRisk
  macroFactors: MacroFactor[]
}

export interface MarketRegime {
  regime: 'bull' | 'bear' | 'sideways' | 'volatile'
  probability: number
  expectedReturn: number
  expectedVolatility: number
  portfolioPerformance: number
}

export interface CyclicalRisk {
  cyclicalExposure: number
  defensiveExposure: number
  growthExposure: number
  valueExposure: number
  cyclicalityScore: number
}

export interface MacroFactor {
  factor: string
  exposure: number
  sensitivity: number
  riskContribution: number
}

export interface CounterpartyRisk {
  counterpartyScore: number
  counterpartyLevel: RiskLevel
  counterparties: CounterpartyExposure[]
  concentrationRisk: number
  creditRisk: CreditRisk
  operationalRisk: OperationalRisk
}

export interface CounterpartyExposure {
  counterparty: string
  exposure: number
  percentage: number
  creditRating: string
  riskScore: number
  defaultProbability: number
}

export interface CreditRisk {
  creditScore: number
  defaultProbability: number
  recoveryRate: number
  creditSpread: number
  creditVaR: number
}

export interface OperationalRisk {
  operationalScore: number
  keyPersonRisk: number
  systemRisk: number
  processRisk: number
  externalRisk: number
}

export interface TechnicalRisk {
  technicalScore: number
  technicalLevel: RiskLevel
  smartContractRisk: SmartContractRisk
  bridgeRisk: BridgeRisk
  oracleRisk: OracleRisk
  governanceRisk: GovernanceRisk
}

export interface SmartContractRisk {
  contractsAnalyzed: number
  vulnerabilities: number
  auditCoverage: number
  codeQuality: number
  upgradeRisk: number
}

export interface BridgeRisk {
  bridgeExposure: number
  bridgeCount: number
  bridgeRiskScore: number
  crossChainRisk: number
}

export interface OracleRisk {
  oracleExposure: number
  oracleCount: number
  oracleRiskScore: number
  manipulationRisk: number
}

export interface GovernanceRisk {
  governanceScore: number
  centralizationRisk: number
  votingPowerConcentration: number
  proposalRisk: number
}

export interface StressTestResult {
  scenario: string
  description: string
  probability: number
  portfolioImpact: number
  worstCaseValue: number
  recoveryTime: number
  affectedAssets: string[]
  mitigationStrategies: string[]
}

export interface RiskMitigationPlan {
  currentRiskLevel: RiskLevel
  targetRiskLevel: RiskLevel
  mitigationStrategies: MitigationStrategy[]
  hedgingRecommendations: HedgingRecommendation[]
  diversificationPlan: DiversificationPlan
  liquidityPlan: LiquidityPlan
  implementationTimeline: ImplementationStep[]
}

export interface MitigationStrategy {
  strategy: string
  description: string
  riskReduction: number
  cost: number
  timeframe: string
  priority: 'high' | 'medium' | 'low'
  implementation: string[]
}

export interface HedgingRecommendation {
  hedgeType: 'options' | 'futures' | 'swaps' | 'insurance'
  asset: string
  hedgeRatio: number
  cost: number
  effectiveness: number
  duration: string
}

export interface DiversificationPlan {
  currentDiversification: number
  targetDiversification: number
  rebalanceActions: RebalanceAction[]
  newAssetRecommendations: AssetRecommendation[]
}

export interface RebalanceAction {
  action: 'reduce' | 'increase' | 'maintain'
  asset: string
  currentWeight: number
  targetWeight: number
  amount: number
}

export interface AssetRecommendation {
  asset: string
  allocation: number
  rationale: string
  riskContribution: number
  expectedReturn: number
}

export interface LiquidityPlan {
  currentLiquidity: number
  targetLiquidity: number
  liquidityActions: LiquidityAction[]
  emergencyPlan: EmergencyLiquidityPlan
}

export interface LiquidityAction {
  action: 'increase' | 'decrease'
  asset: string
  amount: number
  method: string
  timeframe: string
}

export interface EmergencyLiquidityPlan {
  emergencyThreshold: number
  liquidationOrder: string[]
  expectedSlippage: number
  minimumLiquidity: number
}

export interface ImplementationStep {
  step: string
  description: string
  timeframe: string
  priority: number
  dependencies: string[]
  resources: string[]
}

export interface RiskRecommendation {
  id: string
  type: 'reduce_exposure' | 'diversify' | 'hedge' | 'rebalance' | 'monitor'
  priority: 'high' | 'medium' | 'low'
  title: string
  description: string
  rationale: string
  expectedImpact: number
  implementation: string
  timeframe: string
  cost: number
}

export interface RiskAlert {
  id: string
  severity: RiskLevel
  category: string
  title: string
  description: string
  threshold: number
  currentValue: number
  trend: 'increasing' | 'decreasing' | 'stable'
  action: string
  deadline?: string
}

export interface RiskConfiguration {
  riskThresholds: RiskThresholds
  stressTestScenarios: StressTestScenario[]
  alertSettings: AlertSettings
  mitigationPreferences: MitigationPreferences
}

export interface RiskThresholds {
  overall: { low: number; medium: number; high: number; critical: number }
  concentration: { low: number; medium: number; high: number; critical: number }
  liquidity: { low: number; medium: number; high: number; critical: number }
  volatility: { low: number; medium: number; high: number; critical: number }
}

export interface StressTestScenario {
  name: string
  description: string
  factors: ScenarioFactor[]
  probability: number
  enabled: boolean
}

export interface ScenarioFactor {
  factor: string
  change: number
  correlation: number
}

export interface AlertSettings {
  enableAlerts: boolean
  alertChannels: string[]
  alertFrequency: string
  thresholdOverrides: Record<string, number>
}

export interface MitigationPreferences {
  riskTolerance: 'conservative' | 'moderate' | 'aggressive'
  hedgingPreference: boolean
  diversificationTarget: number
  liquidityTarget: number
  maxConcentration: number
}

export class PortfolioRiskManager {
  private static instance: PortfolioRiskManager
  private riskAssessments = new Map<string, PortfolioRiskAssessment>()
  private configuration: RiskConfiguration
  private eventListeners = new Set<(event: RiskEvent) => void>()

  private constructor() {
    this.configuration = this.getDefaultConfiguration()
  }

  static getInstance(): PortfolioRiskManager {
    if (!PortfolioRiskManager.instance) {
      PortfolioRiskManager.instance = new PortfolioRiskManager()
    }
    return PortfolioRiskManager.instance
  }

  /**
   * Get default risk configuration
   */
  private getDefaultConfiguration(): RiskConfiguration {
    return {
      riskThresholds: {
        overall: { low: 20, medium: 40, high: 60, critical: 80 },
        concentration: { low: 15, medium: 25, high: 40, critical: 60 },
        liquidity: { low: 10, medium: 20, high: 35, critical: 50 },
        volatility: { low: 15, medium: 30, high: 50, critical: 75 }
      },
      stressTestScenarios: [
        {
          name: 'Market Crash',
          description: '50% market decline scenario',
          factors: [
            { factor: 'market_return', change: -0.5, correlation: 0.8 },
            { factor: 'volatility', change: 2.0, correlation: 0.6 }
          ],
          probability: 0.05,
          enabled: true
        },
        {
          name: 'DeFi Crisis',
          description: 'Major DeFi protocol failure',
          factors: [
            { factor: 'defi_return', change: -0.7, correlation: 0.9 },
            { factor: 'liquidity', change: -0.8, correlation: 0.7 }
          ],
          probability: 0.1,
          enabled: true
        }
      ],
      alertSettings: {
        enableAlerts: true,
        alertChannels: ['email', 'push'],
        alertFrequency: 'real_time',
        thresholdOverrides: {}
      },
      mitigationPreferences: {
        riskTolerance: 'moderate',
        hedgingPreference: true,
        diversificationTarget: 0.8,
        liquidityTarget: 0.2,
        maxConcentration: 0.3
      }
    }
  }

  /**
   * Assess portfolio risk
   */
  async assessPortfolioRisk(
    portfolioId: string,
    userAddress: Address,
    portfolioData: any
  ): Promise<PortfolioRiskAssessment> {
    try {
      // Emit assessment started event
      this.emitEvent({
        type: 'assessment_started',
        portfolioId,
        timestamp: Date.now()
      })

      // Calculate risk metrics
      const riskMetrics = await this.calculateRiskMetrics(portfolioData)
      
      // Analyze exposures
      const exposureAnalysis = await this.analyzeExposures(portfolioData)
      
      // Assess concentration risk
      const concentrationRisk = await this.assessConcentrationRisk(portfolioData)
      
      // Assess liquidity risk
      const liquidityRisk = await this.assessLiquidityRisk(portfolioData)
      
      // Assess market risk
      const marketRisk = await this.assessMarketRisk(portfolioData)
      
      // Assess counterparty risk
      const counterpartyRisk = await this.assessCounterpartyRisk(portfolioData)
      
      // Assess technical risk
      const technicalRisk = await this.assessTechnicalRisk(portfolioData)
      
      // Run stress tests
      const stressTests = await this.runStressTests(portfolioData)
      
      // Calculate overall risk score
      const riskScore = this.calculateOverallRiskScore(
        riskMetrics,
        concentrationRisk,
        liquidityRisk,
        marketRisk,
        counterpartyRisk,
        technicalRisk
      )
      
      const overallRisk = this.getRiskLevel(riskScore)
      
      // Generate mitigation plan
      const riskMitigation = await this.generateMitigationPlan(
        portfolioData,
        riskScore,
        concentrationRisk,
        liquidityRisk
      )
      
      // Generate recommendations
      const recommendations = this.generateRecommendations(
        riskScore,
        concentrationRisk,
        liquidityRisk,
        marketRisk
      )
      
      // Generate alerts
      const alerts = this.generateAlerts(
        riskScore,
        concentrationRisk,
        liquidityRisk,
        marketRisk
      )

      const assessment: PortfolioRiskAssessment = {
        portfolioId,
        userAddress,
        timestamp: new Date().toISOString(),
        overallRisk,
        riskScore,
        riskMetrics,
        exposureAnalysis,
        concentrationRisk,
        liquidityRisk,
        marketRisk,
        counterpartyRisk,
        technicalRisk,
        stressTests,
        riskMitigation,
        recommendations,
        alerts
      }

      // Store assessment
      this.riskAssessments.set(portfolioId, assessment)

      // Emit assessment completed event
      this.emitEvent({
        type: 'assessment_completed',
        portfolioId,
        assessment,
        timestamp: Date.now()
      })

      return assessment

    } catch (error) {
      // Emit assessment failed event
      this.emitEvent({
        type: 'assessment_failed',
        portfolioId,
        error: error as Error,
        timestamp: Date.now()
      })

      throw error
    }
  }

  /**
   * Calculate risk metrics
   */
  private async calculateRiskMetrics(portfolioData: any): Promise<RiskMetrics> {
    // Mock implementation - in real app, this would calculate actual metrics
    return {
      valueAtRisk: {
        var95: Math.random() * 0.1,
        var99: Math.random() * 0.15,
        timeHorizon: 1,
        confidence: 0.95,
        method: 'historical'
      },
      expectedShortfall: Math.random() * 0.12,
      sharpeRatio: Math.random() * 2 - 0.5,
      sortinoRatio: Math.random() * 2.5 - 0.5,
      maxDrawdown: Math.random() * 0.3,
      volatility: Math.random() * 0.4 + 0.1,
      beta: Math.random() * 2 + 0.5,
      correlation: {
        assets: ['BTC', 'ETH', 'USDC'],
        matrix: [
          [1.0, 0.8, 0.1],
          [0.8, 1.0, 0.2],
          [0.1, 0.2, 1.0]
        ],
        averageCorrelation: 0.37,
        maxCorrelation: 0.8,
        minCorrelation: 0.1
      },
      diversificationRatio: Math.random() * 0.5 + 0.5,
      riskAdjustedReturn: Math.random() * 0.2 - 0.05
    }
  }

  /**
   * Analyze exposures
   */
  private async analyzeExposures(portfolioData: any): Promise<ExposureAnalysis> {
    return {
      totalExposure: 100000,
      exposureByAsset: [
        {
          asset: 'Bitcoin',
          symbol: 'BTC',
          exposure: 40000,
          percentage: 40,
          riskContribution: 35,
          volatility: 0.6,
          liquidity: 0.9
        },
        {
          asset: 'Ethereum',
          symbol: 'ETH',
          exposure: 35000,
          percentage: 35,
          riskContribution: 30,
          volatility: 0.7,
          liquidity: 0.85
        }
      ],
      exposureByProtocol: [
        {
          protocol: 'Uniswap',
          exposure: 25000,
          percentage: 25,
          riskScore: 30,
          tvl: 5000000000,
          auditStatus: 'audited',
          vulnerabilities: 0
        }
      ],
      exposureByChain: [
        {
          chainId: 1,
          chainName: 'Ethereum',
          exposure: 75000,
          percentage: 75,
          riskScore: 25,
          bridgeRisk: 10,
          validatorRisk: 15
        }
      ],
      exposureByCategory: [
        {
          category: 'DeFi',
          exposure: 60000,
          percentage: 60,
          riskScore: 45,
          volatility: 0.8
        }
      ],
      geographicExposure: [
        {
          region: 'Global',
          exposure: 100000,
          percentage: 100,
          regulatoryRisk: 30,
          politicalRisk: 20
        }
      ],
      sectorExposure: [
        {
          sector: 'Cryptocurrency',
          exposure: 100000,
          percentage: 100,
          riskScore: 60,
          correlation: 0.8
        }
      ]
    }
  }

  /**
   * Assess concentration risk
   */
  private async assessConcentrationRisk(portfolioData: any): Promise<ConcentrationRisk> {
    const herfindahlIndex = 0.35 // Mock calculation
    
    return {
      herfindahlIndex,
      top5Concentration: 0.85,
      top10Concentration: 0.95,
      concentrationScore: herfindahlIndex * 100,
      concentrationLevel: this.getRiskLevel(herfindahlIndex * 100),
      concentratedAssets: [
        {
          asset: 'Bitcoin',
          percentage: 40,
          riskContribution: 35,
          recommendedReduction: 10
        }
      ],
      diversificationGap: 0.25
    }
  }

  /**
   * Assess liquidity risk
   */
  private async assessLiquidityRisk(portfolioData: any): Promise<LiquidityRisk> {
    return {
      liquidityScore: 75,
      liquidityLevel: RiskLevel.MEDIUM,
      illiquidAssets: [
        {
          asset: 'NFT Collection',
          liquidityScore: 30,
          dailyVolume: 1000,
          bidAskSpread: 0.15,
          marketDepth: 5000,
          liquidationTime: 7
        }
      ],
      liquidityBuffer: 0.15,
      liquidationTime: {
        immediate: 0.6,
        within1Hour: 0.8,
        within1Day: 0.9,
        within1Week: 0.95,
        moreThan1Week: 0.05
      },
      marketImpact: {
        lowImpact: 0.7,
        mediumImpact: 0.25,
        highImpact: 0.05,
        averageSlippage: 0.02
      }
    }
  }

  /**
   * Assess market risk
   */
  private async assessMarketRisk(portfolioData: any): Promise<MarketRisk> {
    return {
      marketBeta: 1.2,
      marketCorrelation: 0.8,
      systematicRisk: 0.6,
      idiosyncraticRisk: 0.4,
      marketRegimes: [
        {
          regime: 'bull',
          probability: 0.3,
          expectedReturn: 0.15,
          expectedVolatility: 0.25,
          portfolioPerformance: 0.18
        },
        {
          regime: 'bear',
          probability: 0.2,
          expectedReturn: -0.2,
          expectedVolatility: 0.4,
          portfolioPerformance: -0.25
        }
      ],
      cyclicalRisk: {
        cyclicalExposure: 0.7,
        defensiveExposure: 0.1,
        growthExposure: 0.8,
        valueExposure: 0.2,
        cyclicalityScore: 70
      },
      macroFactors: [
        {
          factor: 'Interest Rates',
          exposure: 0.6,
          sensitivity: -0.8,
          riskContribution: 0.25
        }
      ]
    }
  }

  /**
   * Assess counterparty risk
   */
  private async assessCounterpartyRisk(portfolioData: any): Promise<CounterpartyRisk> {
    return {
      counterpartyScore: 65,
      counterpartyLevel: RiskLevel.MEDIUM,
      counterparties: [
        {
          counterparty: 'Binance',
          exposure: 25000,
          percentage: 25,
          creditRating: 'A',
          riskScore: 30,
          defaultProbability: 0.02
        }
      ],
      concentrationRisk: 0.4,
      creditRisk: {
        creditScore: 70,
        defaultProbability: 0.05,
        recoveryRate: 0.6,
        creditSpread: 0.02,
        creditVaR: 0.08
      },
      operationalRisk: {
        operationalScore: 60,
        keyPersonRisk: 40,
        systemRisk: 50,
        processRisk: 45,
        externalRisk: 55
      }
    }
  }

  /**
   * Assess technical risk
   */
  private async assessTechnicalRisk(portfolioData: any): Promise<TechnicalRisk> {
    return {
      technicalScore: 55,
      technicalLevel: RiskLevel.MEDIUM,
      smartContractRisk: {
        contractsAnalyzed: 15,
        vulnerabilities: 2,
        auditCoverage: 0.8,
        codeQuality: 75,
        upgradeRisk: 30
      },
      bridgeRisk: {
        bridgeExposure: 0.2,
        bridgeCount: 3,
        bridgeRiskScore: 45,
        crossChainRisk: 40
      },
      oracleRisk: {
        oracleExposure: 0.6,
        oracleCount: 5,
        oracleRiskScore: 35,
        manipulationRisk: 25
      },
      governanceRisk: {
        governanceScore: 60,
        centralizationRisk: 50,
        votingPowerConcentration: 0.4,
        proposalRisk: 30
      }
    }
  }

  /**
   * Run stress tests
   */
  private async runStressTests(portfolioData: any): Promise<StressTestResult[]> {
    return this.configuration.stressTestScenarios
      .filter(scenario => scenario.enabled)
      .map(scenario => ({
        scenario: scenario.name,
        description: scenario.description,
        probability: scenario.probability,
        portfolioImpact: Math.random() * 0.6 - 0.3,
        worstCaseValue: 70000 + Math.random() * 20000,
        recoveryTime: Math.floor(Math.random() * 365) + 30,
        affectedAssets: ['BTC', 'ETH', 'DeFi Tokens'],
        mitigationStrategies: [
          'Increase cash reserves',
          'Diversify across uncorrelated assets',
          'Implement hedging strategies'
        ]
      }))
  }

  /**
   * Calculate overall risk score
   */
  private calculateOverallRiskScore(
    riskMetrics: RiskMetrics,
    concentrationRisk: ConcentrationRisk,
    liquidityRisk: LiquidityRisk,
    marketRisk: MarketRisk,
    counterpartyRisk: CounterpartyRisk,
    technicalRisk: TechnicalRisk
  ): number {
    // Weighted average of different risk components
    const weights = {
      concentration: 0.25,
      liquidity: 0.2,
      market: 0.2,
      counterparty: 0.15,
      technical: 0.2
    }

    return (
      concentrationRisk.concentrationScore * weights.concentration +
      (100 - liquidityRisk.liquidityScore) * weights.liquidity +
      marketRisk.systematicRisk * 100 * weights.market +
      (100 - counterpartyRisk.counterpartyScore) * weights.counterparty +
      technicalRisk.technicalScore * weights.technical
    )
  }

  /**
   * Get risk level from score
   */
  private getRiskLevel(score: number): RiskLevel {
    const thresholds = this.configuration.riskThresholds.overall
    
    if (score >= thresholds.critical) return RiskLevel.CRITICAL
    if (score >= thresholds.high) return RiskLevel.HIGH
    if (score >= thresholds.medium) return RiskLevel.MEDIUM
    if (score >= thresholds.low) return RiskLevel.LOW
    return RiskLevel.VERY_LOW
  }

  /**
   * Generate mitigation plan
   */
  private async generateMitigationPlan(
    portfolioData: any,
    riskScore: number,
    concentrationRisk: ConcentrationRisk,
    liquidityRisk: LiquidityRisk
  ): Promise<RiskMitigationPlan> {
    const currentRiskLevel = this.getRiskLevel(riskScore)
    const targetRiskLevel = this.getTargetRiskLevel(currentRiskLevel)

    return {
      currentRiskLevel,
      targetRiskLevel,
      mitigationStrategies: [
        {
          strategy: 'Diversification',
          description: 'Reduce concentration in top holdings',
          riskReduction: 15,
          cost: 500,
          timeframe: '2-4 weeks',
          priority: 'high',
          implementation: [
            'Sell 10% of Bitcoin position',
            'Invest in uncorrelated assets',
            'Rebalance monthly'
          ]
        }
      ],
      hedgingRecommendations: [
        {
          hedgeType: 'options',
          asset: 'BTC',
          hedgeRatio: 0.3,
          cost: 1200,
          effectiveness: 0.8,
          duration: '3 months'
        }
      ],
      diversificationPlan: {
        currentDiversification: 0.65,
        targetDiversification: 0.8,
        rebalanceActions: [
          {
            action: 'reduce',
            asset: 'BTC',
            currentWeight: 0.4,
            targetWeight: 0.3,
            amount: 10000
          }
        ],
        newAssetRecommendations: [
          {
            asset: 'Bonds',
            allocation: 0.1,
            rationale: 'Reduce volatility',
            riskContribution: 5,
            expectedReturn: 0.04
          }
        ]
      },
      liquidityPlan: {
        currentLiquidity: 0.15,
        targetLiquidity: 0.25,
        liquidityActions: [
          {
            action: 'increase',
            asset: 'USDC',
            amount: 10000,
            method: 'Convert from illiquid assets',
            timeframe: '1 week'
          }
        ],
        emergencyPlan: {
          emergencyThreshold: 0.1,
          liquidationOrder: ['USDC', 'ETH', 'BTC'],
          expectedSlippage: 0.02,
          minimumLiquidity: 0.05
        }
      },
      implementationTimeline: [
        {
          step: 'Immediate Risk Reduction',
          description: 'Reduce highest risk positions',
          timeframe: '1 week',
          priority: 1,
          dependencies: [],
          resources: ['Trading account', 'Market analysis']
        }
      ]
    }
  }

  /**
   * Get target risk level
   */
  private getTargetRiskLevel(currentLevel: RiskLevel): RiskLevel {
    const preferences = this.configuration.mitigationPreferences
    
    switch (preferences.riskTolerance) {
      case 'conservative':
        return RiskLevel.LOW
      case 'moderate':
        return currentLevel === RiskLevel.CRITICAL ? RiskLevel.HIGH : RiskLevel.MEDIUM
      case 'aggressive':
        return currentLevel === RiskLevel.CRITICAL ? RiskLevel.MEDIUM : currentLevel
      default:
        return RiskLevel.MEDIUM
    }
  }

  /**
   * Generate recommendations
   */
  private generateRecommendations(
    riskScore: number,
    concentrationRisk: ConcentrationRisk,
    liquidityRisk: LiquidityRisk,
    marketRisk: MarketRisk
  ): RiskRecommendation[] {
    const recommendations: RiskRecommendation[] = []

    // High concentration risk
    if (concentrationRisk.concentrationLevel === RiskLevel.HIGH || concentrationRisk.concentrationLevel === RiskLevel.CRITICAL) {
      recommendations.push({
        id: 'reduce_concentration',
        type: 'diversify',
        priority: 'high',
        title: 'Reduce Portfolio Concentration',
        description: 'Your portfolio is highly concentrated in a few assets',
        rationale: 'High concentration increases portfolio volatility and risk',
        expectedImpact: 20,
        implementation: 'Gradually reduce positions in top holdings and diversify',
        timeframe: '2-4 weeks',
        cost: 500
      })
    }

    // Low liquidity
    if (liquidityRisk.liquidityLevel === RiskLevel.HIGH || liquidityRisk.liquidityLevel === RiskLevel.CRITICAL) {
      recommendations.push({
        id: 'improve_liquidity',
        type: 'rebalance',
        priority: 'medium',
        title: 'Improve Portfolio Liquidity',
        description: 'Portfolio has low liquidity which may impact exit strategies',
        rationale: 'Better liquidity provides more flexibility during market stress',
        expectedImpact: 15,
        implementation: 'Convert some illiquid assets to more liquid alternatives',
        timeframe: '1-2 weeks',
        cost: 300
      })
    }

    // High market risk
    if (marketRisk.marketBeta > 1.5) {
      recommendations.push({
        id: 'hedge_market_risk',
        type: 'hedge',
        priority: 'medium',
        title: 'Hedge Market Risk',
        description: 'Portfolio has high market beta and systematic risk',
        rationale: 'Hedging can reduce downside risk during market downturns',
        expectedImpact: 25,
        implementation: 'Consider options or futures hedging strategies',
        timeframe: '1 week',
        cost: 1000
      })
    }

    return recommendations
  }

  /**
   * Generate alerts
   */
  private generateAlerts(
    riskScore: number,
    concentrationRisk: ConcentrationRisk,
    liquidityRisk: LiquidityRisk,
    marketRisk: MarketRisk
  ): RiskAlert[] {
    const alerts: RiskAlert[] = []

    // Overall risk alert
    if (riskScore >= this.configuration.riskThresholds.overall.high) {
      alerts.push({
        id: 'high_overall_risk',
        severity: this.getRiskLevel(riskScore),
        category: 'overall',
        title: 'High Portfolio Risk',
        description: 'Portfolio risk score exceeds acceptable threshold',
        threshold: this.configuration.riskThresholds.overall.high,
        currentValue: riskScore,
        trend: 'increasing',
        action: 'Review and implement risk mitigation strategies'
      })
    }

    // Concentration alert
    if (concentrationRisk.concentrationScore >= this.configuration.riskThresholds.concentration.high) {
      alerts.push({
        id: 'high_concentration',
        severity: concentrationRisk.concentrationLevel,
        category: 'concentration',
        title: 'High Concentration Risk',
        description: 'Portfolio is highly concentrated in few assets',
        threshold: this.configuration.riskThresholds.concentration.high,
        currentValue: concentrationRisk.concentrationScore,
        trend: 'stable',
        action: 'Diversify portfolio holdings'
      })
    }

    return alerts
  }

  /**
   * Get risk assessment
   */
  getRiskAssessment(portfolioId: string): PortfolioRiskAssessment | null {
    return this.riskAssessments.get(portfolioId) || null
  }

  /**
   * Update configuration
   */
  updateConfiguration(config: Partial<RiskConfiguration>): void {
    this.configuration = { ...this.configuration, ...config }
  }

  /**
   * Get configuration
   */
  getConfiguration(): RiskConfiguration {
    return { ...this.configuration }
  }

  /**
   * Emit event to listeners
   */
  private emitEvent(event: RiskEvent): void {
    for (const listener of Array.from(this.eventListeners)) {
      try {
        listener(event)
      } catch (error) {
        console.error('Error in risk event listener:', error)
      }
    }
  }

  /**
   * Add event listener
   */
  addEventListener(listener: (event: RiskEvent) => void): () => void {
    this.eventListeners.add(listener)
    
    return () => {
      this.eventListeners.delete(listener)
    }
  }

  /**
   * Clear all data
   */
  clear(): void {
    this.riskAssessments.clear()
  }

  /**
   * Cleanup resources
   */
  destroy(): void {
    this.clear()
    this.eventListeners.clear()
  }
}

export interface RiskEvent {
  type: 'assessment_started' | 'assessment_completed' | 'assessment_failed' | 'alert_triggered' | 'threshold_breached'
  portfolioId: string
  assessment?: PortfolioRiskAssessment
  alert?: RiskAlert
  error?: Error
  timestamp: number
}

// Export singleton instance
export const portfolioRiskManager = PortfolioRiskManager.getInstance()
