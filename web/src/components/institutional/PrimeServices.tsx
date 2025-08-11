'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { 
  Building2, 
  TrendingUp, 
  Shield,
  DollarSign,
  Clock,
  Users,
  BarChart3,
  Zap,
  Globe,
  Award,
  Target,
  Briefcase,
  CreditCard,
  Banknote,
  PieChart
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface PrimeBrokerageService {
  id: string
  name: string
  description: string
  type: 'lending' | 'borrowing' | 'margin' | 'derivatives' | 'custody'
  status: 'active' | 'pending' | 'suspended'
  tier: 'standard' | 'premium' | 'institutional'
  minRequirement: number
  currentExposure: number
  maxExposure: number
  interestRate: number
  fees: {
    management: number
    performance: number
    transaction: number
  }
}

interface CreditFacility {
  id: string
  type: 'revolving' | 'term' | 'bridge'
  currency: string
  approvedLimit: number
  utilizedAmount: number
  availableAmount: number
  interestRate: number
  maturityDate: number
  collateralRatio: number
  status: 'active' | 'pending' | 'expired'
}

interface InstitutionalReport {
  id: string
  type: 'daily' | 'weekly' | 'monthly' | 'quarterly'
  title: string
  description: string
  generatedAt: number
  size: string
  format: 'pdf' | 'excel' | 'csv'
  isCustom: boolean
}

export function PrimeServices() {
  const [services, setServices] = useState<PrimeBrokerageService[]>([])
  const [creditFacilities, setCreditFacilities] = useState<CreditFacility[]>([])
  const [reports, setReports] = useState<InstitutionalReport[]>([])
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    // Generate mock prime services data
    const mockServices: PrimeBrokerageService[] = [
      {
        id: 'service1',
        name: 'Securities Lending',
        description: 'Earn yield by lending digital assets to qualified borrowers',
        type: 'lending',
        status: 'active',
        tier: 'institutional',
        minRequirement: 1000000,
        currentExposure: 15000000,
        maxExposure: 50000000,
        interestRate: 8.5,
        fees: { management: 0.25, performance: 15, transaction: 0.05 }
      },
      {
        id: 'service2',
        name: 'Margin Trading',
        description: 'Access leveraged trading with institutional-grade risk management',
        type: 'margin',
        status: 'active',
        tier: 'premium',
        minRequirement: 500000,
        currentExposure: 8500000,
        maxExposure: 25000000,
        interestRate: 12.0,
        fees: { management: 0.15, performance: 10, transaction: 0.03 }
      },
      {
        id: 'service3',
        name: 'Derivatives Trading',
        description: 'Access to futures, options, and structured products',
        type: 'derivatives',
        status: 'pending',
        tier: 'institutional',
        minRequirement: 2000000,
        currentExposure: 0,
        maxExposure: 100000000,
        interestRate: 0,
        fees: { management: 0.35, performance: 20, transaction: 0.08 }
      }
    ]

    const mockCreditFacilities: CreditFacility[] = [
      {
        id: 'credit1',
        type: 'revolving',
        currency: 'USD',
        approvedLimit: 50000000,
        utilizedAmount: 12500000,
        availableAmount: 37500000,
        interestRate: 6.5,
        maturityDate: Date.now() + 365 * 24 * 60 * 60 * 1000,
        collateralRatio: 150,
        status: 'active'
      },
      {
        id: 'credit2',
        type: 'term',
        currency: 'EUR',
        approvedLimit: 25000000,
        utilizedAmount: 8000000,
        availableAmount: 17000000,
        interestRate: 7.2,
        maturityDate: Date.now() + 180 * 24 * 60 * 60 * 1000,
        collateralRatio: 125,
        status: 'active'
      }
    ]

    const mockReports: InstitutionalReport[] = [
      {
        id: 'report1',
        type: 'daily',
        title: 'Daily Risk Report',
        description: 'Comprehensive risk metrics and exposure analysis',
        generatedAt: Date.now() - 3600000,
        size: '2.4 MB',
        format: 'pdf',
        isCustom: false
      },
      {
        id: 'report2',
        type: 'monthly',
        title: 'Monthly Performance Report',
        description: 'Detailed performance analysis and attribution',
        generatedAt: Date.now() - 86400000,
        size: '8.7 MB',
        format: 'excel',
        isCustom: true
      }
    ]

    setServices(mockServices)
    setCreditFacilities(mockCreditFacilities)
    setReports(mockReports)
  }, [])

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

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'active': return 'text-green-500'
      case 'pending': return 'text-yellow-500'
      case 'suspended': case 'expired': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'active': return 'default'
      case 'pending': return 'secondary'
      case 'suspended': case 'expired': return 'destructive'
      default: return 'outline'
    }
  }

  const getTierColor = (tier: string) => {
    switch (tier) {
      case 'institutional': return 'text-purple-600'
      case 'premium': return 'text-blue-600'
      case 'standard': return 'text-gray-600'
      default: return 'text-muted-foreground'
    }
  }

  const getTotalExposure = () => {
    return services.reduce((sum, service) => sum + service.currentExposure, 0)
  }

  const getTotalCreditUtilization = () => {
    return creditFacilities.reduce((sum, facility) => sum + facility.utilizedAmount, 0)
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Prime Services</h2>
          <p className="text-muted-foreground">
            Institutional-grade financial services and credit facilities
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Award className="w-3 h-3 mr-1" />
            Prime Client
          </Badge>
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Tier 1 Institution
          </Badge>
        </div>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <DollarSign className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Total Exposure</span>
            </div>
            <div className="text-2xl font-bold">{formatCurrency(getTotalExposure())}</div>
            <div className="text-xs text-green-500">
              Across {services.filter(s => s.status === 'active').length} services
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <CreditCard className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Credit Utilized</span>
            </div>
            <div className="text-2xl font-bold">{formatCurrency(getTotalCreditUtilization())}</div>
            <div className="text-xs text-muted-foreground">
              {creditFacilities.length} facilities
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <TrendingUp className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Avg Yield</span>
            </div>
            <div className="text-2xl font-bold text-green-500">
              {(services.reduce((sum, s) => sum + s.interestRate, 0) / services.length).toFixed(1)}%
            </div>
            <div className="text-xs text-muted-foreground">
              Weighted average
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <BarChart3 className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Reports Generated</span>
            </div>
            <div className="text-2xl font-bold">{reports.length}</div>
            <div className="text-xs text-muted-foreground">
              This month
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="services">Prime Services</TabsTrigger>
          <TabsTrigger value="credit">Credit Facilities</TabsTrigger>
          <TabsTrigger value="reports">Reports</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Service Utilization */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <PieChart className="w-5 h-5" />
                  Service Utilization
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {services.filter(s => s.status === 'active').map((service) => (
                    <div key={service.id} className="space-y-2">
                      <div className="flex items-center justify-between">
                        <span className="font-medium">{service.name}</span>
                        <span className="text-sm text-muted-foreground">
                          {((service.currentExposure / service.maxExposure) * 100).toFixed(1)}%
                        </span>
                      </div>
                      <Progress 
                        value={(service.currentExposure / service.maxExposure) * 100} 
                        className="h-2" 
                      />
                      <div className="flex justify-between text-xs text-muted-foreground">
                        <span>{formatCurrency(service.currentExposure)}</span>
                        <span>Max: {formatCurrency(service.maxExposure)}</span>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>

            {/* Credit Overview */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Banknote className="w-5 h-5" />
                  Credit Overview
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  {creditFacilities.map((facility) => (
                    <div key={facility.id} className="p-3 border rounded">
                      <div className="flex items-center justify-between mb-2">
                        <div>
                          <div className="font-medium capitalize">{facility.type} Credit</div>
                          <div className="text-sm text-muted-foreground">{facility.currency}</div>
                        </div>
                        <Badge variant={getStatusBadgeVariant(facility.status)}>
                          {facility.status}
                        </Badge>
                      </div>
                      <div className="space-y-2">
                        <div className="flex justify-between text-sm">
                          <span className="text-muted-foreground">Utilized</span>
                          <span className="font-medium">
                            {formatCurrency(facility.utilizedAmount)}
                          </span>
                        </div>
                        <Progress 
                          value={(facility.utilizedAmount / facility.approvedLimit) * 100} 
                          className="h-2" 
                        />
                        <div className="flex justify-between text-xs text-muted-foreground">
                          <span>Available: {formatCurrency(facility.availableAmount)}</span>
                          <span>Rate: {facility.interestRate}%</span>
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="services" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Prime Brokerage Services</h3>
            <Button>
              <Briefcase className="w-4 h-4 mr-2" />
              Request New Service
            </Button>
          </div>

          <div className="space-y-4">
            {services.map((service) => (
              <Card key={service.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <Building2 className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{service.name}</h4>
                        <p className="text-sm text-muted-foreground">{service.description}</p>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className={getTierColor(service.tier)}>
                        {service.tier}
                      </Badge>
                      <Badge variant={getStatusBadgeVariant(service.status)}>
                        {service.status}
                      </Badge>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Current Exposure</div>
                      <div className="font-medium">{formatCurrency(service.currentExposure)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Max Exposure</div>
                      <div className="font-medium">{formatCurrency(service.maxExposure)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Interest Rate</div>
                      <div className="font-medium text-green-500">{service.interestRate}%</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Management Fee</div>
                      <div className="font-medium">{service.fees.management}%</div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="text-sm text-muted-foreground">
                      Min Requirement: {formatCurrency(service.minRequirement)}
                    </div>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        View Details
                      </Button>
                      {service.status === 'active' && (
                        <Button size="sm">
                          Manage
                        </Button>
                      )}
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="credit" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Credit Facilities</h3>
            <Button>
              <CreditCard className="w-4 h-4 mr-2" />
              Apply for Credit
            </Button>
          </div>

          <div className="space-y-4">
            {creditFacilities.map((facility) => (
              <Card key={facility.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold capitalize">{facility.type} Credit Facility</h4>
                      <p className="text-sm text-muted-foreground">
                        {facility.currency} • Expires {formatTime(facility.maturityDate)}
                      </p>
                    </div>
                    <Badge variant={getStatusBadgeVariant(facility.status)}>
                      {facility.status}
                    </Badge>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Approved Limit</div>
                      <div className="font-medium">{formatCurrency(facility.approvedLimit)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Utilized</div>
                      <div className="font-medium">{formatCurrency(facility.utilizedAmount)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Available</div>
                      <div className="font-medium text-green-500">
                        {formatCurrency(facility.availableAmount)}
                      </div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Interest Rate</div>
                      <div className="font-medium">{facility.interestRate}%</div>
                    </div>
                  </div>

                  <div className="space-y-2 mb-4">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Utilization</span>
                      <span className="font-medium">
                        {((facility.utilizedAmount / facility.approvedLimit) * 100).toFixed(1)}%
                      </span>
                    </div>
                    <Progress 
                      value={(facility.utilizedAmount / facility.approvedLimit) * 100} 
                      className="h-2" 
                    />
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="text-sm text-muted-foreground">
                      Collateral Ratio: {facility.collateralRatio}%
                    </div>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        Draw Down
                      </Button>
                      <Button variant="outline" size="sm">
                        Repay
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="reports" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Institutional Reports</h3>
            <Button>
              <BarChart3 className="w-4 h-4 mr-2" />
              Generate Custom Report
            </Button>
          </div>

          <div className="space-y-4">
            {reports.map((report) => (
              <Card key={report.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between">
                    <div className="flex items-center gap-3">
                      <div className="w-10 h-10 bg-primary/10 rounded-full flex items-center justify-center">
                        <BarChart3 className="w-5 h-5" />
                      </div>
                      <div>
                        <h4 className="font-bold">{report.title}</h4>
                        <p className="text-sm text-muted-foreground">{report.description}</p>
                        <div className="flex items-center gap-4 text-xs text-muted-foreground mt-1">
                          <span>Generated: {formatTime(report.generatedAt)}</span>
                          <span>Size: {report.size}</span>
                          <span>Format: {report.format.toUpperCase()}</span>
                        </div>
                      </div>
                    </div>
                    <div className="flex items-center gap-2">
                      {report.isCustom && (
                        <Badge variant="outline">Custom</Badge>
                      )}
                      <Badge variant="outline" className="capitalize">
                        {report.type}
                      </Badge>
                      <Button variant="outline" size="sm">
                        Download
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>
      </Tabs>
    </div>
  )
}
