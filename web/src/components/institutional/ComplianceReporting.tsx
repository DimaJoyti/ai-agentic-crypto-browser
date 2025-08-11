'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { 
  FileText, 
  Shield, 
  CheckCircle,
  AlertTriangle,
  Clock,
  Download,
  Upload,
  Eye,
  Settings,
  Calendar,
  BarChart3,
  Users,
  Globe,
  Award,
  Lock
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface ComplianceRequirement {
  id: string
  name: string
  description: string
  jurisdiction: string
  type: 'kyc' | 'aml' | 'reporting' | 'licensing' | 'audit'
  status: 'compliant' | 'pending' | 'overdue' | 'not_applicable'
  dueDate?: number
  lastUpdated: number
  completionPercent: number
  assignedTo: string
  priority: 'low' | 'medium' | 'high' | 'critical'
}

interface RegulatoryReport {
  id: string
  name: string
  type: 'daily' | 'weekly' | 'monthly' | 'quarterly' | 'annual' | 'ad_hoc'
  jurisdiction: string
  regulator: string
  status: 'draft' | 'pending_review' | 'submitted' | 'approved' | 'rejected'
  dueDate: number
  submittedDate?: number
  size: string
  format: 'xml' | 'csv' | 'pdf' | 'json'
  isAutomated: boolean
}

interface AuditTrail {
  id: string
  action: string
  user: string
  timestamp: number
  details: string
  ipAddress: string
  riskLevel: 'low' | 'medium' | 'high'
  category: 'trade' | 'transfer' | 'admin' | 'compliance' | 'security'
}

interface KYCRecord {
  id: string
  clientId: string
  clientName: string
  clientType: 'individual' | 'corporate' | 'institutional'
  status: 'pending' | 'approved' | 'rejected' | 'expired' | 'under_review'
  tier: 'basic' | 'enhanced' | 'institutional'
  lastReview: number
  nextReview: number
  riskRating: 'low' | 'medium' | 'high'
  jurisdiction: string
  documentsRequired: number
  documentsReceived: number
}

export function ComplianceReporting() {
  const [requirements, setRequirements] = useState<ComplianceRequirement[]>([])
  const [reports, setReports] = useState<RegulatoryReport[]>([])
  const [auditTrail, setAuditTrail] = useState<AuditTrail[]>([])
  const [kycRecords, setKycRecords] = useState<KYCRecord[]>([])
  const [activeTab, setActiveTab] = useState('overview')

  useEffect(() => {
    // Generate mock compliance data
    const mockRequirements: ComplianceRequirement[] = [
      {
        id: 'req1',
        name: 'MiFID II Transaction Reporting',
        description: 'Daily transaction reporting to ESMA under MiFID II regulations',
        jurisdiction: 'EU',
        type: 'reporting',
        status: 'compliant',
        dueDate: Date.now() + 86400000,
        lastUpdated: Date.now() - 3600000,
        completionPercent: 100,
        assignedTo: 'Compliance Team',
        priority: 'high'
      },
      {
        id: 'req2',
        name: 'CFTC Large Trader Reporting',
        description: 'Weekly large trader position reporting to CFTC',
        jurisdiction: 'US',
        type: 'reporting',
        status: 'pending',
        dueDate: Date.now() + 172800000,
        lastUpdated: Date.now() - 7200000,
        completionPercent: 75,
        assignedTo: 'Risk Team',
        priority: 'medium'
      },
      {
        id: 'req3',
        name: 'Enhanced KYC Review',
        description: 'Annual enhanced KYC review for institutional clients',
        jurisdiction: 'Global',
        type: 'kyc',
        status: 'overdue',
        dueDate: Date.now() - 86400000,
        lastUpdated: Date.now() - 172800000,
        completionPercent: 60,
        assignedTo: 'KYC Team',
        priority: 'critical'
      }
    ]

    const mockReports: RegulatoryReport[] = [
      {
        id: 'report1',
        name: 'Daily Transaction Report',
        type: 'daily',
        jurisdiction: 'EU',
        regulator: 'ESMA',
        status: 'submitted',
        dueDate: Date.now() + 86400000,
        submittedDate: Date.now() - 3600000,
        size: '2.4 MB',
        format: 'xml',
        isAutomated: true
      },
      {
        id: 'report2',
        name: 'Monthly Risk Report',
        type: 'monthly',
        jurisdiction: 'US',
        regulator: 'CFTC',
        status: 'pending_review',
        dueDate: Date.now() + 604800000,
        size: '15.7 MB',
        format: 'pdf',
        isAutomated: false
      },
      {
        id: 'report3',
        name: 'Quarterly Capital Report',
        type: 'quarterly',
        jurisdiction: 'UK',
        regulator: 'FCA',
        status: 'draft',
        dueDate: Date.now() + 1209600000,
        size: '8.3 MB',
        format: 'csv',
        isAutomated: false
      }
    ]

    const mockAuditTrail: AuditTrail[] = [
      {
        id: 'audit1',
        action: 'Large Trade Executed',
        user: 'trader@institution.com',
        timestamp: Date.now() - 1800000,
        details: 'Executed BTC trade for $5.2M',
        ipAddress: '192.168.1.100',
        riskLevel: 'medium',
        category: 'trade'
      },
      {
        id: 'audit2',
        action: 'KYC Document Uploaded',
        user: 'compliance@institution.com',
        timestamp: Date.now() - 3600000,
        details: 'Uploaded enhanced due diligence documents for Client ABC',
        ipAddress: '192.168.1.101',
        riskLevel: 'low',
        category: 'compliance'
      },
      {
        id: 'audit3',
        action: 'Risk Limit Modified',
        user: 'risk@institution.com',
        timestamp: Date.now() - 7200000,
        details: 'Increased VaR limit from $5M to $7M',
        ipAddress: '192.168.1.102',
        riskLevel: 'high',
        category: 'admin'
      }
    ]

    const mockKYCRecords: KYCRecord[] = [
      {
        id: 'kyc1',
        clientId: 'CLIENT_001',
        clientName: 'Institutional Fund ABC',
        clientType: 'institutional',
        status: 'approved',
        tier: 'institutional',
        lastReview: Date.now() - 86400000 * 180,
        nextReview: Date.now() + 86400000 * 185,
        riskRating: 'low',
        jurisdiction: 'US',
        documentsRequired: 12,
        documentsReceived: 12
      },
      {
        id: 'kyc2',
        clientId: 'CLIENT_002',
        clientName: 'Hedge Fund XYZ',
        clientType: 'institutional',
        status: 'under_review',
        tier: 'enhanced',
        lastReview: Date.now() - 86400000 * 30,
        nextReview: Date.now() + 86400000 * 335,
        riskRating: 'medium',
        jurisdiction: 'EU',
        documentsRequired: 8,
        documentsReceived: 6
      }
    ]

    setRequirements(mockRequirements)
    setReports(mockReports)
    setAuditTrail(mockAuditTrail)
    setKycRecords(mockKYCRecords)
  }, [])

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleDateString()
  }

  const formatDateTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'compliant': case 'approved': case 'submitted': return 'text-green-500'
      case 'pending': case 'pending_review': case 'under_review': case 'draft': return 'text-yellow-500'
      case 'overdue': case 'rejected': case 'expired': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'compliant': case 'approved': case 'submitted': return 'default'
      case 'pending': case 'pending_review': case 'under_review': case 'draft': return 'secondary'
      case 'overdue': case 'rejected': case 'expired': return 'destructive'
      default: return 'outline'
    }
  }

  const getPriorityColor = (priority: string) => {
    switch (priority) {
      case 'critical': return 'text-red-500'
      case 'high': return 'text-orange-500'
      case 'medium': return 'text-yellow-500'
      case 'low': return 'text-blue-500'
      default: return 'text-muted-foreground'
    }
  }

  const getRiskColor = (risk: string) => {
    switch (risk) {
      case 'high': return 'text-red-500'
      case 'medium': return 'text-yellow-500'
      case 'low': return 'text-green-500'
      default: return 'text-muted-foreground'
    }
  }

  const getComplianceScore = () => {
    const total = requirements.length
    const compliant = requirements.filter(r => r.status === 'compliant').length
    return Math.round((compliant / total) * 100)
  }

  const complianceScore = getComplianceScore()

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Compliance & Reporting</h2>
          <p className="text-muted-foreground">
            Regulatory compliance monitoring and automated reporting
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Compliance Score: {complianceScore}%
          </Badge>
          <Button variant="outline" size="sm">
            <Settings className="w-3 h-3 mr-1" />
            Configure
          </Button>
        </div>
      </div>

      {/* Overview Stats */}
      <div className="grid grid-cols-1 md:grid-cols-4 gap-4">
        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <CheckCircle className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Compliance Score</span>
            </div>
            <div className="text-2xl font-bold text-green-500">{complianceScore}%</div>
            <div className="text-xs text-muted-foreground">
              {requirements.filter(r => r.status === 'compliant').length}/{requirements.length} requirements
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <FileText className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Pending Reports</span>
            </div>
            <div className="text-2xl font-bold">
              {reports.filter(r => r.status === 'draft' || r.status === 'pending_review').length}
            </div>
            <div className="text-xs text-muted-foreground">
              Due this week
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <Users className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">KYC Reviews</span>
            </div>
            <div className="text-2xl font-bold">
              {kycRecords.filter(k => k.status === 'under_review').length}
            </div>
            <div className="text-xs text-muted-foreground">
              In progress
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-4">
            <div className="flex items-center gap-2 mb-2">
              <AlertTriangle className="w-4 h-4 text-muted-foreground" />
              <span className="text-sm text-muted-foreground">Overdue Items</span>
            </div>
            <div className="text-2xl font-bold text-red-500">
              {requirements.filter(r => r.status === 'overdue').length}
            </div>
            <div className="text-xs text-muted-foreground">
              Require attention
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Interface */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-5">
          <TabsTrigger value="overview">Overview</TabsTrigger>
          <TabsTrigger value="requirements">Requirements</TabsTrigger>
          <TabsTrigger value="reports">Reports</TabsTrigger>
          <TabsTrigger value="kyc">KYC/AML</TabsTrigger>
          <TabsTrigger value="audit">Audit Trail</TabsTrigger>
        </TabsList>

        <TabsContent value="overview" className="space-y-4">
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-4">
            {/* Compliance Status */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Shield className="w-5 h-5" />
                  Compliance Status
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-4">
                  <div className="flex justify-between items-center">
                    <span className="text-sm text-muted-foreground">Overall Compliance</span>
                    <span className="font-bold text-green-500">{complianceScore}%</span>
                  </div>
                  <Progress value={complianceScore} className="h-3" />
                  
                  <div className="space-y-2">
                    {['compliant', 'pending', 'overdue'].map((status) => {
                      const count = requirements.filter(r => r.status === status).length
                      const percentage = (count / requirements.length) * 100
                      return (
                        <div key={status} className="flex justify-between text-sm">
                          <span className="capitalize">{status}</span>
                          <span className={getStatusColor(status)}>{count} ({percentage.toFixed(0)}%)</span>
                        </div>
                      )
                    })}
                  </div>
                </div>
              </CardContent>
            </Card>

            {/* Recent Activity */}
            <Card>
              <CardHeader>
                <CardTitle className="flex items-center gap-2">
                  <Clock className="w-5 h-5" />
                  Recent Activity
                </CardTitle>
              </CardHeader>
              <CardContent>
                <div className="space-y-3">
                  {auditTrail.slice(0, 5).map((entry) => (
                    <div key={entry.id} className="flex items-start gap-3 p-3 border rounded">
                      <div className={cn(
                        "w-2 h-2 rounded-full mt-2",
                        entry.riskLevel === 'high' ? "bg-red-500" :
                        entry.riskLevel === 'medium' ? "bg-yellow-500" :
                        "bg-green-500"
                      )} />
                      <div className="flex-1">
                        <div className="font-medium">{entry.action}</div>
                        <div className="text-sm text-muted-foreground">{entry.details}</div>
                        <div className="text-xs text-muted-foreground mt-1">
                          {entry.user} • {formatDateTime(entry.timestamp)}
                        </div>
                      </div>
                    </div>
                  ))}
                </div>
              </CardContent>
            </Card>
          </div>
        </TabsContent>

        <TabsContent value="requirements" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Compliance Requirements</h3>
            <Button>
              <FileText className="w-4 h-4 mr-2" />
              Add Requirement
            </Button>
          </div>

          <div className="space-y-4">
            {requirements.map((req) => (
              <Card key={req.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{req.name}</h4>
                      <p className="text-sm text-muted-foreground">{req.description}</p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className={getPriorityColor(req.priority)}>
                        {req.priority}
                      </Badge>
                      <Badge variant={getStatusBadgeVariant(req.status)}>
                        {req.status}
                      </Badge>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Jurisdiction</div>
                      <div className="font-medium">{req.jurisdiction}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Type</div>
                      <div className="font-medium capitalize">{req.type}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Assigned To</div>
                      <div className="font-medium">{req.assignedTo}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Due Date</div>
                      <div className="font-medium">
                        {req.dueDate ? formatTime(req.dueDate) : 'N/A'}
                      </div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Completion</span>
                      <span className="font-medium">{req.completionPercent}%</span>
                    </div>
                    <Progress value={req.completionPercent} className="h-2" />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="reports" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Regulatory Reports</h3>
            <Button>
              <Upload className="w-4 h-4 mr-2" />
              Generate Report
            </Button>
          </div>

          <div className="space-y-4">
            {reports.map((report) => (
              <Card key={report.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{report.name}</h4>
                      <p className="text-sm text-muted-foreground">
                        {report.regulator} • {report.jurisdiction}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      {report.isAutomated && (
                        <Badge variant="outline">
                          <BarChart3 className="w-3 h-3 mr-1" />
                          Automated
                        </Badge>
                      )}
                      <Badge variant={getStatusBadgeVariant(report.status)}>
                        {report.status.replace('_', ' ')}
                      </Badge>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Type</div>
                      <div className="font-medium capitalize">{report.type}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Format</div>
                      <div className="font-medium uppercase">{report.format}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Size</div>
                      <div className="font-medium">{report.size}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Due Date</div>
                      <div className="font-medium">{formatTime(report.dueDate)}</div>
                    </div>
                  </div>

                  <div className="flex items-center justify-between">
                    <div className="text-sm text-muted-foreground">
                      {report.submittedDate && `Submitted: ${formatTime(report.submittedDate)}`}
                    </div>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm">
                        <Eye className="w-3 h-3 mr-1" />
                        View
                      </Button>
                      <Button variant="outline" size="sm">
                        <Download className="w-3 h-3 mr-1" />
                        Download
                      </Button>
                    </div>
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="kyc" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">KYC/AML Records</h3>
            <Button>
              <Users className="w-4 h-4 mr-2" />
              New KYC Review
            </Button>
          </div>

          <div className="space-y-4">
            {kycRecords.map((record) => (
              <Card key={record.id}>
                <CardContent className="p-6">
                  <div className="flex items-center justify-between mb-4">
                    <div>
                      <h4 className="font-bold">{record.clientName}</h4>
                      <p className="text-sm text-muted-foreground">
                        {record.clientId} • {record.jurisdiction}
                      </p>
                    </div>
                    <div className="flex items-center gap-2">
                      <Badge variant="outline" className={getRiskColor(record.riskRating)}>
                        {record.riskRating} risk
                      </Badge>
                      <Badge variant={getStatusBadgeVariant(record.status)}>
                        {record.status.replace('_', ' ')}
                      </Badge>
                    </div>
                  </div>

                  <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-sm mb-4">
                    <div>
                      <div className="text-muted-foreground">Client Type</div>
                      <div className="font-medium capitalize">{record.clientType}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Tier</div>
                      <div className="font-medium capitalize">{record.tier}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Last Review</div>
                      <div className="font-medium">{formatTime(record.lastReview)}</div>
                    </div>
                    <div>
                      <div className="text-muted-foreground">Next Review</div>
                      <div className="font-medium">{formatTime(record.nextReview)}</div>
                    </div>
                  </div>

                  <div className="space-y-2">
                    <div className="flex justify-between text-sm">
                      <span className="text-muted-foreground">Document Collection</span>
                      <span className="font-medium">
                        {record.documentsReceived}/{record.documentsRequired}
                      </span>
                    </div>
                    <Progress 
                      value={(record.documentsReceived / record.documentsRequired) * 100} 
                      className="h-2" 
                    />
                  </div>
                </CardContent>
              </Card>
            ))}
          </div>
        </TabsContent>

        <TabsContent value="audit" className="space-y-4">
          <div className="flex items-center justify-between">
            <h3 className="text-lg font-medium">Audit Trail</h3>
            <Button variant="outline">
              <Download className="w-4 h-4 mr-2" />
              Export Logs
            </Button>
          </div>

          <div className="space-y-4">
            {auditTrail.map((entry) => (
              <Card key={entry.id}>
                <CardContent className="p-6">
                  <div className="flex items-start justify-between">
                    <div className="flex items-start gap-3">
                      <div className={cn(
                        "w-3 h-3 rounded-full mt-1",
                        entry.riskLevel === 'high' ? "bg-red-500" :
                        entry.riskLevel === 'medium' ? "bg-yellow-500" :
                        "bg-green-500"
                      )} />
                      <div>
                        <h4 className="font-bold">{entry.action}</h4>
                        <p className="text-sm text-muted-foreground">{entry.details}</p>
                        <div className="grid grid-cols-2 md:grid-cols-4 gap-4 text-xs text-muted-foreground mt-2">
                          <span>User: {entry.user}</span>
                          <span>IP: {entry.ipAddress}</span>
                          <span>Category: {entry.category}</span>
                          <span>Time: {formatDateTime(entry.timestamp)}</span>
                        </div>
                      </div>
                    </div>
                    <Badge variant="outline" className={getRiskColor(entry.riskLevel)}>
                      {entry.riskLevel} risk
                    </Badge>
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
