'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Rocket, 
  TrendingUp,
  Clock,
  Users,
  DollarSign,
  Target,
  Star,
  Calendar,
  CheckCircle,
  AlertTriangle,
  Zap,
  Shield,
  Globe,
  Award,
  Eye,
  Heart,
  Share2,
  ExternalLink,
  Twitter,
  MessageCircle,
  FileText,
  BarChart3,
  Lock,
  Unlock,
  Calculator
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface ProjectDetailsProps {
  projectId: string
}

interface VestingSchedule {
  percentage: number
  unlockDate: number
  description: string
  isUnlocked: boolean
  claimedAmount?: number
}

interface TeamMember {
  name: string
  role: string
  avatar: string
  linkedin?: string
  twitter?: string
  experience: string
}

interface Tokenomics {
  category: string
  percentage: number
  amount: number
  description: string
  vestingPeriod?: string
  color: string
}

interface UserParticipation {
  isParticipating: boolean
  allocatedAmount: number
  contributedAmount: number
  tokensReceived: number
  vestingSchedule: VestingSchedule[]
  claimableAmount: number
}

export function ProjectDetails({ projectId }: ProjectDetailsProps) {
  const [project, setProject] = useState<any>(null)
  const [userParticipation, setUserParticipation] = useState<UserParticipation | null>(null)
  const [participationAmount, setParticipationAmount] = useState('')
  const [activeTab, setActiveTab] = useState('overview')

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Mock project details
    const mockProject = {
      id: projectId,
      name: 'DeFiMax Protocol',
      symbol: 'DMAX',
      description: 'DeFiMax Protocol is a next-generation DeFi platform that combines automated yield optimization with cross-chain liquidity aggregation. Our innovative approach allows users to maximize their returns while minimizing risk through advanced algorithmic strategies.',
      longDescription: `DeFiMax Protocol represents the future of decentralized finance, offering a comprehensive suite of tools designed to optimize yield generation across multiple blockchain networks. 

Our platform features:
- Automated yield farming strategies
- Cross-chain liquidity aggregation
- Advanced risk management protocols
- Governance token with utility
- Institutional-grade security measures

The protocol has been audited by leading security firms and has undergone extensive testing to ensure the safety of user funds.`,
      logo: '/projects/defimax.png',
      website: 'https://defimax.io',
      twitter: '@defimax',
      telegram: 't.me/defimax',
      whitepaper: 'https://defimax.io/whitepaper.pdf',
      category: 'defi',
      status: 'live',
      launchType: 'ido',
      
      // Financial Details
      totalRaise: 2000000,
      targetRaise: 1500000,
      tokenPrice: 0.15,
      totalSupply: 100000000,
      tokensForSale: 10000000,
      minAllocation: 100,
      maxAllocation: 5000,
      
      // Timeline
      registrationStart: Date.now() - 86400000 * 7,
      registrationEnd: Date.now() - 86400000 * 2,
      saleStart: Date.now() - 86400000,
      saleEnd: Date.now() + 86400000 * 2,
      vestingStart: Date.now() + 86400000 * 7,
      
      // Progress
      currentRaised: 1250000,
      participantCount: 2847,
      
      // Team
      team: [
        {
          name: 'Alex Chen',
          role: 'CEO & Founder',
          avatar: '/team/alex.jpg',
          linkedin: 'linkedin.com/in/alexchen',
          twitter: '@alexchen',
          experience: '10+ years in DeFi and blockchain development'
        },
        {
          name: 'Sarah Johnson',
          role: 'CTO',
          avatar: '/team/sarah.jpg',
          linkedin: 'linkedin.com/in/sarahjohnson',
          experience: 'Former lead engineer at major DeFi protocols'
        },
        {
          name: 'Michael Rodriguez',
          role: 'Head of Product',
          avatar: '/team/michael.jpg',
          twitter: '@mrodriguez',
          experience: 'Product management at top fintech companies'
        }
      ],
      
      // Tokenomics
      tokenomics: [
        { category: 'Public Sale', percentage: 10, amount: 10000000, description: 'IDO and public distribution', color: '#3b82f6' },
        { category: 'Team', percentage: 15, amount: 15000000, description: 'Team allocation with 2-year vesting', vestingPeriod: '24 months', color: '#ef4444' },
        { category: 'Advisors', percentage: 5, amount: 5000000, description: 'Strategic advisors and partners', vestingPeriod: '12 months', color: '#f59e0b' },
        { category: 'Development', percentage: 25, amount: 25000000, description: 'Protocol development and maintenance', vestingPeriod: '36 months', color: '#10b981' },
        { category: 'Marketing', percentage: 10, amount: 10000000, description: 'Marketing and community growth', vestingPeriod: '18 months', color: '#8b5cf6' },
        { category: 'Liquidity', percentage: 20, amount: 20000000, description: 'DEX liquidity and market making', color: '#06b6d4' },
        { category: 'Treasury', percentage: 15, amount: 15000000, description: 'Protocol treasury and reserves', color: '#84cc16' }
      ],
      
      // Vesting
      vestingSchedule: [
        { percentage: 20, unlockDate: Date.now() + 86400000 * 7, description: 'TGE Unlock', isUnlocked: false },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 37, description: '1 Month', isUnlocked: false },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 67, description: '2 Months', isUnlocked: false },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 97, description: '3 Months', isUnlocked: false },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 127, description: '4 Months', isUnlocked: false }
      ],
      
      // Risk & Audit
      riskScore: 25,
      auditStatus: 'completed',
      auditReports: [
        { firm: 'CertiK', status: 'completed', score: 95, report: 'https://certik.io/defimax' },
        { firm: 'Quantstamp', status: 'completed', score: 92, report: 'https://quantstamp.com/defimax' }
      ],
      teamVerified: true,
      
      // Social
      likes: 1247,
      views: 15632,
      isLiked: false,
      isWatchlisted: true
    }

    const mockUserParticipation: UserParticipation = {
      isParticipating: true,
      allocatedAmount: 1000,
      contributedAmount: 1000,
      tokensReceived: 6666.67,
      claimableAmount: 1333.33,
      vestingSchedule: [
        { percentage: 20, unlockDate: Date.now() + 86400000 * 7, description: 'TGE Unlock', isUnlocked: false, claimedAmount: 0 },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 37, description: '1 Month', isUnlocked: false, claimedAmount: 0 },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 67, description: '2 Months', isUnlocked: false, claimedAmount: 0 },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 97, description: '3 Months', isUnlocked: false, claimedAmount: 0 },
        { percentage: 20, unlockDate: Date.now() + 86400000 * 127, description: '4 Months', isUnlocked: false, claimedAmount: 0 }
      ]
    }

    setProject(mockProject)
    setUserParticipation(mockUserParticipation)
  }, [projectId, isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString()
  }

  const calculateTokensReceived = (amount: number) => {
    if (!project) return 0
    return amount / project.tokenPrice
  }

  const handleParticipate = () => {
    const amount = parseFloat(participationAmount)
    if (amount >= project.minAllocation && amount <= project.maxAllocation) {
      // Handle participation logic
      console.log('Participating with amount:', amount)
    }
  }

  if (!isConnected || !project) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Rocket className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Loading Project Details</h3>
          <p className="text-muted-foreground">
            Please wait while we load the project information
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Project Header */}
      <Card>
        <CardContent className="p-6">
          <div className="flex items-start justify-between mb-4">
            <div className="flex items-center gap-4">
              <div className="w-16 h-16 bg-primary/10 rounded-full flex items-center justify-center">
                <Rocket className="w-8 h-8" />
              </div>
              <div>
                <h2 className="text-3xl font-bold">{project.name}</h2>
                <p className="text-xl text-muted-foreground">{project.symbol}</p>
                <div className="flex items-center gap-2 mt-2">
                  <Badge variant="default">Live</Badge>
                  <Badge variant="outline">IDO</Badge>
                  <Badge variant="outline" className="capitalize">{project.category}</Badge>
                </div>
              </div>
            </div>
            <div className="flex items-center gap-2">
              <Button variant="outline" size="sm">
                <Heart className="w-4 h-4 mr-2" />
                {project.likes}
              </Button>
              <Button variant="outline" size="sm">
                <Share2 className="w-4 h-4 mr-2" />
                Share
              </Button>
            </div>
          </div>

          <p className="text-muted-foreground mb-4">{project.description}</p>

          <div className="flex items-center gap-4">
            <Button variant="outline" size="sm">
              <Globe className="w-4 h-4 mr-2" />
              Website
            </Button>
            <Button variant="outline" size="sm">
              <Twitter className="w-4 h-4 mr-2" />
              Twitter
            </Button>
            <Button variant="outline" size="sm">
              <MessageCircle className="w-4 h-4 mr-2" />
              Telegram
            </Button>
            <Button variant="outline" size="sm">
              <FileText className="w-4 h-4 mr-2" />
              Whitepaper
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Sale Progress */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="w-5 h-5" />
            Sale Progress
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex justify-between text-sm">
              <span>Progress</span>
              <span>{((project.currentRaised / project.targetRaise) * 100).toFixed(1)}%</span>
            </div>
            <Progress value={(project.currentRaised / project.targetRaise) * 100} className="h-3" />
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm">
              <div>
                <div className="text-muted-foreground">Raised</div>
                <div className="font-bold text-lg">{formatCurrency(project.currentRaised)}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Target</div>
                <div className="font-bold text-lg">{formatCurrency(project.targetRaise)}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Participants</div>
                <div className="font-bold text-lg">{project.participantCount.toLocaleString()}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Time Left</div>
                <div className="font-bold text-lg">2d 14h</div>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Participation Card */}
      {project.status === 'live' && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Zap className="w-5 h-5" />
              Participate in Sale
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-4">
              <div className="grid grid-cols-2 gap-4 text-sm">
                <div>
                  <div className="text-muted-foreground">Token Price</div>
                  <div className="font-bold">${project.tokenPrice}</div>
                </div>
                <div>
                  <div className="text-muted-foreground">Min/Max Allocation</div>
                  <div className="font-bold">${project.minAllocation} - ${project.maxAllocation}</div>
                </div>
              </div>

              <div className="space-y-2">
                <Label htmlFor="amount">Investment Amount (USD)</Label>
                <Input
                  id="amount"
                  type="number"
                  placeholder="Enter amount"
                  value={participationAmount}
                  onChange={(e) => setParticipationAmount(e.target.value)}
                  min={project.minAllocation}
                  max={project.maxAllocation}
                />
                <div className="text-sm text-muted-foreground">
                  You will receive: {participationAmount ? calculateTokensReceived(parseFloat(participationAmount)).toFixed(2) : '0'} {project.symbol}
                </div>
              </div>

              <Button 
                className="w-full" 
                onClick={handleParticipate}
                disabled={!participationAmount || parseFloat(participationAmount) < project.minAllocation}
              >
                Participate Now
              </Button>

              <Alert>
                <Shield className="h-4 w-4" />
                <AlertDescription>
                  This project requires KYC verification. Make sure you have completed the verification process.
                </AlertDescription>
              </Alert>
            </div>
          </CardContent>
        </Card>
      )}

      {/* User Participation Status */}
      {userParticipation?.isParticipating && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <CheckCircle className="w-5 h-5 text-green-500" />
              Your Participation
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
              <div>
                <div className="text-muted-foreground">Contributed</div>
                <div className="font-bold">{formatCurrency(userParticipation.contributedAmount)}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Tokens Allocated</div>
                <div className="font-bold">{userParticipation.tokensReceived.toFixed(2)} {project.symbol}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Claimable Now</div>
                <div className="font-bold text-green-500">{userParticipation.claimableAmount.toFixed(2)} {project.symbol}</div>
              </div>
              <div>
                <div className="text-muted-foreground">Next Unlock</div>
                <div className="font-bold">{formatTime(userParticipation.vestingSchedule[0].unlockDate)}</div>
              </div>
            </div>

            {userParticipation.claimableAmount > 0 && (
              <Button className="w-full mb-4">
                <Unlock className="w-4 h-4 mr-2" />
                Claim {userParticipation.claimableAmount.toFixed(2)} {project.symbol}
              </Button>
            )}

            <div className="space-y-2">
              <div className="text-sm font-medium">Vesting Schedule</div>
              {userParticipation.vestingSchedule.map((vest, index) => (
                <div key={index} className="flex items-center justify-between p-2 border rounded">
                  <div className="flex items-center gap-2">
                    {vest.isUnlocked ? (
                      <Unlock className="w-4 h-4 text-green-500" />
                    ) : (
                      <Lock className="w-4 h-4 text-muted-foreground" />
                    )}
                    <span className="text-sm">{vest.description}</span>
                  </div>
                  <div className="text-right">
                    <div className="text-sm font-medium">
                      {((userParticipation.tokensReceived * vest.percentage) / 100).toFixed(2)} {project.symbol}
                    </div>
                    <div className="text-xs text-muted-foreground">{formatTime(vest.unlockDate)}</div>
                  </div>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}

      {/* Detailed Information Tabs */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="tokenomics">Tokenomics</TabsTrigger>
          <TabsTrigger value="team">Team</TabsTrigger>
          <TabsTrigger value="roadmap">Roadmap</TabsTrigger>
          <TabsTrigger value="audit">Security</TabsTrigger>
        </TabsList>

        <TabsContent value="overview">
          <Card>
            <CardHeader>
              <CardTitle>Project Overview</CardTitle>
            </CardHeader>
            <CardContent>
              <div className="prose max-w-none">
                <div className="whitespace-pre-line text-muted-foreground">
                  {project.longDescription}
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="tokenomics">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <BarChart3 className="w-5 h-5" />
                Token Distribution
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                {project.tokenomics.map((item: Tokenomics, index: number) => (
                  <div key={index} className="flex items-center justify-between p-3 border rounded">
                    <div className="flex items-center gap-3">
                      <div 
                        className="w-4 h-4 rounded-full" 
                        style={{ backgroundColor: item.color }}
                      />
                      <div>
                        <div className="font-medium">{item.category}</div>
                        <div className="text-sm text-muted-foreground">{item.description}</div>
                        {item.vestingPeriod && (
                          <div className="text-xs text-muted-foreground">Vesting: {item.vestingPeriod}</div>
                        )}
                      </div>
                    </div>
                    <div className="text-right">
                      <div className="font-bold">{item.percentage}%</div>
                      <div className="text-sm text-muted-foreground">
                        {(item.amount / 1000000).toFixed(1)}M tokens
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="team">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Users className="w-5 h-5" />
                Team Members
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
                {project.team.map((member: TeamMember, index: number) => (
                  <div key={index} className="p-4 border rounded-lg">
                    <div className="flex items-center gap-3 mb-3">
                      <div className="w-12 h-12 bg-primary/10 rounded-full flex items-center justify-center">
                        <Users className="w-6 h-6" />
                      </div>
                      <div>
                        <div className="font-bold">{member.name}</div>
                        <div className="text-sm text-muted-foreground">{member.role}</div>
                      </div>
                    </div>
                    <p className="text-sm text-muted-foreground mb-3">{member.experience}</p>
                    <div className="flex gap-2">
                      {member.linkedin && (
                        <Button variant="outline" size="sm">
                          <ExternalLink className="w-3 h-3" />
                        </Button>
                      )}
                      {member.twitter && (
                        <Button variant="outline" size="sm">
                          <Twitter className="w-3 h-3" />
                        </Button>
                      )}
                    </div>
                  </div>
                ))}
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="roadmap">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Calendar className="w-5 h-5" />
                Project Roadmap
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-start gap-4">
                  <div className="w-4 h-4 bg-green-500 rounded-full mt-1" />
                  <div>
                    <div className="font-medium">Q1 2024 - Protocol Launch</div>
                    <div className="text-sm text-muted-foreground">Core protocol deployment and initial features</div>
                  </div>
                </div>
                <div className="flex items-start gap-4">
                  <div className="w-4 h-4 bg-blue-500 rounded-full mt-1" />
                  <div>
                    <div className="font-medium">Q2 2024 - Cross-chain Integration</div>
                    <div className="text-sm text-muted-foreground">Multi-chain support and bridge functionality</div>
                  </div>
                </div>
                <div className="flex items-start gap-4">
                  <div className="w-4 h-4 bg-gray-300 rounded-full mt-1" />
                  <div>
                    <div className="font-medium">Q3 2024 - Advanced Features</div>
                    <div className="text-sm text-muted-foreground">Governance implementation and advanced strategies</div>
                  </div>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>

        <TabsContent value="audit">
          <Card>
            <CardHeader>
              <CardTitle className="flex items-center gap-2">
                <Shield className="w-5 h-5" />
                Security & Audits
              </CardTitle>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div className="flex items-center gap-2 mb-4">
                  <CheckCircle className="w-5 h-5 text-green-500" />
                  <span className="font-medium">Team Verified</span>
                </div>

                <div className="space-y-3">
                  {project.auditReports.map((audit: any, index: number) => (
                    <div key={index} className="flex items-center justify-between p-3 border rounded">
                      <div>
                        <div className="font-medium">{audit.firm} Audit</div>
                        <div className="text-sm text-muted-foreground">Security Score: {audit.score}/100</div>
                      </div>
                      <div className="flex items-center gap-2">
                        <Badge variant="default">Completed</Badge>
                        <Button variant="outline" size="sm">
                          <ExternalLink className="w-3 h-3 mr-1" />
                          View Report
                        </Button>
                      </div>
                    </div>
                  ))}
                </div>

                <div className="p-4 bg-green-50 dark:bg-green-950/20 border border-green-200 dark:border-green-800 rounded">
                  <div className="flex items-center gap-2 mb-2">
                    <Shield className="w-4 h-4 text-green-600" />
                    <span className="font-medium text-green-800 dark:text-green-200">Low Risk Score</span>
                  </div>
                  <p className="text-sm text-green-700 dark:text-green-300">
                    This project has a low risk score of {project.riskScore}/100 based on team verification, audit results, and tokenomics analysis.
                  </p>
                </div>
              </div>
            </CardContent>
          </Card>
        </TabsContent>
      </Tabs>
    </div>
  )
}
