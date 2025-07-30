'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Button } from '@/components/ui/button'
import { Progress } from '@/components/ui/progress'
import { 
  Shield, 
  AlertTriangle, 
  CheckCircle, 
  Clock,
  FileText,
  TrendingUp,
  Users,
  BarChart3
} from 'lucide-react'

interface ComplianceOverviewProps {
  frameworks: any[]
  violations: any[]
  alerts: any[]
  metrics: any
}

export const ComplianceOverview: React.FC<ComplianceOverviewProps> = ({
  frameworks,
  violations,
  alerts: _alerts,
  metrics
}) => {
  const getStatusColor = (status: string) => {
    switch (status.toLowerCase()) {
      case 'compliant':
        return 'bg-green-500'
      case 'partial':
        return 'bg-yellow-500'
      case 'non_compliant':
        return 'bg-red-500'
      default:
        return 'bg-gray-500'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status.toLowerCase()) {
      case 'compliant':
        return 'default'
      case 'partial':
        return 'secondary'
      case 'non_compliant':
        return 'destructive'
      default:
        return 'outline'
    }
  }

  return (
    <div className="space-y-6">
      {/* Compliance Frameworks */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Shield className="h-5 w-5" />
            Compliance Frameworks
          </CardTitle>
          <CardDescription>
            Status of regulatory compliance frameworks
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            {frameworks.map((framework) => (
              <div key={framework.id} className="flex items-center justify-between p-4 border rounded-lg">
                <div className="flex items-center gap-4">
                  <div className={`w-3 h-3 rounded-full ${getStatusColor(framework.status)}`} />
                  <div>
                    <h4 className="font-medium">{framework.name}</h4>
                    <p className="text-sm text-muted-foreground">
                      {framework.jurisdiction} • Last updated: {new Date(framework.last_update).toLocaleDateString()}
                    </p>
                  </div>
                </div>
                <div className="flex items-center gap-3">
                  <div className="text-right">
                    <p className="text-sm font-medium">{framework.compliance_rate || 85}%</p>
                    <Progress value={framework.compliance_rate || 85} className="w-20 h-2" />
                  </div>
                  <Badge variant={getStatusBadgeVariant(framework.status)}>
                    {framework.status}
                  </Badge>
                  <Button variant="outline" size="sm">
                    View Details
                  </Button>
                </div>
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Compliance Summary */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <BarChart3 className="h-5 w-5" />
              Compliance Score
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center">
              <div className="text-4xl font-bold text-green-600 mb-2">
                {metrics?.complianceScore || 85}
              </div>
              <Progress value={metrics?.complianceScore || 85} className="mb-4" />
              <p className="text-sm text-muted-foreground">
                Overall compliance rating
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <TrendingUp className="h-5 w-5" />
              Risk Assessment
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="text-center">
              <div className="text-4xl font-bold text-yellow-600 mb-2">
                {metrics?.riskScore || 25}
              </div>
              <div className="flex items-center justify-center gap-2 mb-4">
                <div className={`w-3 h-3 rounded-full ${
                  (metrics?.riskScore || 25) < 30 ? 'bg-green-500' :
                  (metrics?.riskScore || 25) < 60 ? 'bg-yellow-500' : 'bg-red-500'
                }`} />
                <span className="text-sm font-medium">
                  {(metrics?.riskScore || 25) < 30 ? 'Low Risk' :
                   (metrics?.riskScore || 25) < 60 ? 'Medium Risk' : 'High Risk'}
                </span>
              </div>
              <p className="text-sm text-muted-foreground">
                Current risk level
              </p>
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Users className="h-5 w-5" />
              Audit Activity
            </CardTitle>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex justify-between">
                <span className="text-sm">Events Today</span>
                <span className="font-medium">247</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm">High Risk Events</span>
                <span className="font-medium text-red-600">3</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm">Failed Logins</span>
                <span className="font-medium">12</span>
              </div>
              <div className="flex justify-between">
                <span className="text-sm">Data Exports</span>
                <span className="font-medium">5</span>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Recent Activity */}
      <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <AlertTriangle className="h-5 w-5" />
              Critical Issues
            </CardTitle>
            <CardDescription>
              Issues requiring immediate attention
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              {violations.filter(v => v.severity === 'CRITICAL' || v.severity === 'HIGH').slice(0, 3).map((violation) => (
                <div key={violation.id} className="flex items-center gap-3 p-3 border rounded-lg">
                  <AlertTriangle className="h-4 w-4 text-red-500" />
                  <div className="flex-1">
                    <p className="font-medium text-sm">{violation.description}</p>
                    <p className="text-xs text-muted-foreground">
                      {violation.framework} • {new Date(violation.timestamp).toLocaleDateString()}
                    </p>
                  </div>
                  <Badge variant="destructive">{violation.severity}</Badge>
                </div>
              ))}
              {violations.filter(v => v.severity === 'CRITICAL' || v.severity === 'HIGH').length === 0 && (
                <div className="text-center py-4 text-muted-foreground">
                  <CheckCircle className="h-8 w-8 mx-auto mb-2 text-green-500" />
                  <p>No critical issues found</p>
                </div>
              )}
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Clock className="h-5 w-5" />
              Upcoming Deadlines
            </CardTitle>
            <CardDescription>
              Reports and reviews due soon
            </CardDescription>
          </CardHeader>
          <CardContent>
            <div className="space-y-3">
              <div className="flex items-center gap-3 p-3 border rounded-lg">
                <FileText className="h-4 w-4 text-blue-500" />
                <div className="flex-1">
                  <p className="font-medium text-sm">Monthly BSA Report</p>
                  <p className="text-xs text-muted-foreground">
                    Due in 5 days
                  </p>
                </div>
                <Badge variant="secondary">Pending</Badge>
              </div>
              <div className="flex items-center gap-3 p-3 border rounded-lg">
                <FileText className="h-4 w-4 text-blue-500" />
                <div className="flex-1">
                  <p className="font-medium text-sm">Quarterly Risk Review</p>
                  <p className="text-xs text-muted-foreground">
                    Due in 12 days
                  </p>
                </div>
                <Badge variant="outline">Scheduled</Badge>
              </div>
              <div className="flex items-center gap-3 p-3 border rounded-lg">
                <FileText className="h-4 w-4 text-blue-500" />
                <div className="flex-1">
                  <p className="font-medium text-sm">EU AMLD5 Compliance Check</p>
                  <p className="text-xs text-muted-foreground">
                    Due in 18 days
                  </p>
                </div>
                <Badge variant="outline">Scheduled</Badge>
              </div>
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Compliance Trends */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <BarChart3 className="h-5 w-5" />
            Compliance Trends
          </CardTitle>
          <CardDescription>
            Monthly compliance and risk metrics
          </CardDescription>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-4 gap-6">
            <div className="text-center">
              <div className="text-2xl font-bold text-green-600 mb-1">
                {metrics?.violationsThisMonth || 5}
              </div>
              <p className="text-sm text-muted-foreground">Violations This Month</p>
              <div className="flex items-center justify-center gap-1 mt-1">
                <TrendingUp className="h-3 w-3 text-red-500" />
                <span className="text-xs text-red-500">+2 from last month</span>
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-blue-600 mb-1">
                {metrics?.resolvedViolations || 4}
              </div>
              <p className="text-sm text-muted-foreground">Resolved Violations</p>
              <div className="flex items-center justify-center gap-1 mt-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                <span className="text-xs text-green-500">+1 from last month</span>
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-purple-600 mb-1">
                {metrics?.reportsGenerated || 12}
              </div>
              <p className="text-sm text-muted-foreground">Reports Generated</p>
              <div className="flex items-center justify-center gap-1 mt-1">
                <TrendingUp className="h-3 w-3 text-green-500" />
                <span className="text-xs text-green-500">Same as last month</span>
              </div>
            </div>
            <div className="text-center">
              <div className="text-2xl font-bold text-orange-600 mb-1">
                2.5h
              </div>
              <p className="text-sm text-muted-foreground">Avg Resolution Time</p>
              <div className="flex items-center justify-center gap-1 mt-1">
                <TrendingUp className="h-3 w-3 text-red-500" />
                <span className="text-xs text-red-500">+0.5h from last month</span>
              </div>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
