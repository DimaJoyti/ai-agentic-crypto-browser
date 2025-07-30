'use client'

import React from 'react'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'

export const AlertManagement: React.FC = () => {
  return (
    <Card>
      <CardHeader>
        <CardTitle>Alert Management</CardTitle>
        <CardDescription>Manage compliance and risk alerts</CardDescription>
      </CardHeader>
      <CardContent>
        <p className="text-muted-foreground">Alert management component coming soon...</p>
      </CardContent>
    </Card>
  )
}
