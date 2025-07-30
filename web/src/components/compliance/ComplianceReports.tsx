'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export const ComplianceReports: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Compliance Reports</CardTitle>
        <CardDescription>Generate and manage compliance reports</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">Compliance reports component coming soon...</p>
      </CardContent>
    </Card>
  )
}
