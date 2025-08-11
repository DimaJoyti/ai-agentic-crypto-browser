'use client'

import React, { useState, useEffect } from 'react'
import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import { Label } from '@/components/ui/label'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Badge } from '@/components/ui/badge'
import { Progress } from '@/components/ui/progress'
import { Alert, AlertDescription } from '@/components/ui/alert'

import { 
  Shield, 
  User, 
  FileText,
  Camera,
  CheckCircle,
  AlertTriangle,
  Clock,
  Upload,
  Eye,
  Download,
  Globe,
  Building2,
  CreditCard,
  Phone,
  Mail,
  MapPin,
  Calendar,
  Zap
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface KYCLevel {
  level: number
  name: string
  description: string
  requirements: string[]
  limits: {
    daily: number
    monthly: number
    annual: number
  }
  features: string[]
  processingTime: string
  isCompleted: boolean
  isAvailable: boolean
}

interface VerificationDocument {
  id: string
  type: 'passport' | 'drivers_license' | 'national_id' | 'utility_bill' | 'bank_statement' | 'selfie'
  name: string
  description: string
  status: 'not_uploaded' | 'uploaded' | 'under_review' | 'approved' | 'rejected'
  uploadedAt?: number
  reviewedAt?: number
  rejectionReason?: string
  required: boolean
  fileSize?: string
  fileName?: string
}

interface PersonalInfo {
  firstName: string
  lastName: string
  dateOfBirth: string
  nationality: string
  phoneNumber: string
  email: string
  address: {
    street: string
    city: string
    state: string
    postalCode: string
    country: string
  }
  occupation: string
  sourceOfFunds: string
}

interface VerificationStatus {
  level: number
  status: 'not_started' | 'in_progress' | 'under_review' | 'approved' | 'rejected'
  completionPercentage: number
  submittedAt?: number
  approvedAt?: number
  rejectedAt?: number
  rejectionReason?: string
  nextSteps: string[]
}

export function KYCVerification() {
  const [currentLevel, setCurrentLevel] = useState(1)
  const [verificationStatus, setVerificationStatus] = useState<VerificationStatus>({
    level: 1,
    status: 'in_progress',
    completionPercentage: 65,
    nextSteps: ['Upload proof of address', 'Complete identity verification']
  })
  const [documents, setDocuments] = useState<VerificationDocument[]>([])
  const [personalInfo, setPersonalInfo] = useState<PersonalInfo>({
    firstName: '',
    lastName: '',
    dateOfBirth: '',
    nationality: '',
    phoneNumber: '',
    email: '',
    address: {
      street: '',
      city: '',
      state: '',
      postalCode: '',
      country: ''
    },
    occupation: '',
    sourceOfFunds: ''
  })

  const { address, isConnected } = useAccount()

  const kycLevels: KYCLevel[] = [
    {
      level: 1,
      name: 'Basic Verification',
      description: 'Email and phone verification for basic trading',
      requirements: ['Email verification', 'Phone verification'],
      limits: { daily: 1000, monthly: 5000, annual: 25000 },
      features: ['Basic trading', 'Crypto deposits/withdrawals'],
      processingTime: 'Instant',
      isCompleted: true,
      isAvailable: true
    },
    {
      level: 2,
      name: 'Identity Verification',
      description: 'Government ID verification for increased limits',
      requirements: ['Government ID', 'Selfie verification', 'Personal information'],
      limits: { daily: 10000, monthly: 50000, annual: 250000 },
      features: ['Fiat deposits/withdrawals', 'Higher trading limits', 'Credit card purchases'],
      processingTime: '1-2 business days',
      isCompleted: false,
      isAvailable: true
    },
    {
      level: 3,
      name: 'Enhanced Verification',
      description: 'Address and source of funds verification for maximum limits',
      requirements: ['Proof of address', 'Source of funds', 'Enhanced due diligence'],
      limits: { daily: 100000, monthly: 500000, annual: 2500000 },
      features: ['Maximum limits', 'OTC trading', 'Institutional features', 'Priority support'],
      processingTime: '3-5 business days',
      isCompleted: false,
      isAvailable: false
    }
  ]

  useEffect(() => {
    if (!isConnected) return

    // Generate mock document data
    const mockDocuments: VerificationDocument[] = [
      {
        id: 'email',
        type: 'passport',
        name: 'Email Verification',
        description: 'Verify your email address',
        status: 'approved',
        uploadedAt: Date.now() - 86400000,
        reviewedAt: Date.now() - 86400000 + 3600000,
        required: true
      },
      {
        id: 'phone',
        type: 'drivers_license',
        name: 'Phone Verification',
        description: 'Verify your phone number via SMS',
        status: 'approved',
        uploadedAt: Date.now() - 86400000,
        reviewedAt: Date.now() - 86400000 + 1800000,
        required: true
      },
      {
        id: 'government_id',
        type: 'passport',
        name: 'Government ID',
        description: 'Upload passport, driver\'s license, or national ID',
        status: 'under_review',
        uploadedAt: Date.now() - 3600000,
        required: true,
        fileSize: '2.4 MB',
        fileName: 'passport_front.jpg'
      },
      {
        id: 'selfie',
        type: 'selfie',
        name: 'Selfie Verification',
        description: 'Take a selfie holding your ID document',
        status: 'uploaded',
        uploadedAt: Date.now() - 1800000,
        required: true,
        fileSize: '1.8 MB',
        fileName: 'selfie_with_id.jpg'
      },
      {
        id: 'address_proof',
        type: 'utility_bill',
        name: 'Proof of Address',
        description: 'Upload utility bill, bank statement, or government letter',
        status: 'not_uploaded',
        required: true
      },
      {
        id: 'bank_statement',
        type: 'bank_statement',
        name: 'Bank Statement',
        description: 'Upload recent bank statement for source of funds verification',
        status: 'not_uploaded',
        required: false
      }
    ]

    setDocuments(mockDocuments)
  }, [isConnected])

  const formatCurrency = (amount: number) => {
    return new Intl.NumberFormat('en-US', {
      style: 'currency',
      currency: 'USD',
      minimumFractionDigits: 0,
      maximumFractionDigits: 0
    }).format(amount)
  }

  const formatTime = (timestamp: number) => {
    return new Date(timestamp).toLocaleString()
  }

  const getStatusColor = (status: string) => {
    switch (status) {
      case 'approved': return 'text-green-500'
      case 'under_review': case 'uploaded': return 'text-blue-500'
      case 'in_progress': return 'text-yellow-500'
      case 'rejected': case 'not_uploaded': return 'text-red-500'
      default: return 'text-muted-foreground'
    }
  }

  const getStatusBadgeVariant = (status: string) => {
    switch (status) {
      case 'approved': return 'default'
      case 'under_review': case 'uploaded': case 'in_progress': return 'secondary'
      case 'rejected': case 'not_uploaded': return 'destructive'
      default: return 'outline'
    }
  }

  const getStatusIcon = (status: string) => {
    switch (status) {
      case 'approved': return <CheckCircle className="w-4 h-4 text-green-500" />
      case 'under_review': case 'uploaded': return <Clock className="w-4 h-4 text-blue-500" />
      case 'in_progress': return <Clock className="w-4 h-4 text-yellow-500" />
      case 'rejected': case 'not_uploaded': return <AlertTriangle className="w-4 h-4 text-red-500" />
      default: return <Clock className="w-4 h-4 text-muted-foreground" />
    }
  }

  const getDocumentIcon = (type: string) => {
    switch (type) {
      case 'passport': case 'drivers_license': case 'national_id': return <FileText className="w-4 h-4" />
      case 'utility_bill': case 'bank_statement': return <Building2 className="w-4 h-4" />
      case 'selfie': return <Camera className="w-4 h-4" />
      default: return <FileText className="w-4 h-4" />
    }
  }

  const handleFileUpload = (documentId: string) => {
    // Simulate file upload
    setDocuments(prev => prev.map(doc => 
      doc.id === documentId 
        ? { ...doc, status: 'uploaded', uploadedAt: Date.now(), fileSize: '2.1 MB', fileName: 'document.jpg' }
        : doc
    ))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access KYC verification
          </p>
        </CardContent>
      </Card>
    )
  }

  return (
    <div className="space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">KYC Verification</h2>
          <p className="text-muted-foreground">
            Complete verification to unlock higher trading limits and features
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Level {verificationStatus.level}
          </Badge>
          <Badge variant={getStatusBadgeVariant(verificationStatus.status)}>
            {verificationStatus.status.replace('_', ' ')}
          </Badge>
        </div>
      </div>

      {/* Verification Progress */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="w-5 h-5" />
            Verification Progress
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Overall Progress</span>
              <span className="font-medium">{verificationStatus.completionPercentage}%</span>
            </div>
            <Progress value={verificationStatus.completionPercentage} className="h-3" />
            
            {verificationStatus.nextSteps.length > 0 && (
              <div>
                <h4 className="font-medium mb-2">Next Steps:</h4>
                <ul className="space-y-1">
                  {verificationStatus.nextSteps.map((step, index) => (
                    <li key={index} className="flex items-center gap-2 text-sm">
                      <div className="w-2 h-2 bg-primary rounded-full" />
                      {step}
                    </li>
                  ))}
                </ul>
              </div>
            )}
          </div>
        </CardContent>
      </Card>

      {/* KYC Levels */}
      <div className="grid grid-cols-1 md:grid-cols-3 gap-4">
        {kycLevels.map((level) => (
          <Card key={level.level} className={cn(
            "relative",
            level.level === currentLevel && "ring-2 ring-primary",
            !level.isAvailable && "opacity-60"
          )}>
            <CardHeader>
              <div className="flex items-center justify-between">
                <CardTitle className="text-lg">
                  Level {level.level}: {level.name}
                </CardTitle>
                {level.isCompleted && (
                  <CheckCircle className="w-5 h-5 text-green-500" />
                )}
              </div>
              <p className="text-sm text-muted-foreground">{level.description}</p>
            </CardHeader>
            <CardContent>
              <div className="space-y-4">
                <div>
                  <h4 className="font-medium mb-2">Trading Limits</h4>
                  <div className="space-y-1 text-sm">
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Daily</span>
                      <span>{formatCurrency(level.limits.daily)}</span>
                    </div>
                    <div className="flex justify-between">
                      <span className="text-muted-foreground">Monthly</span>
                      <span>{formatCurrency(level.limits.monthly)}</span>
                    </div>
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Features</h4>
                  <div className="flex flex-wrap gap-1">
                    {level.features.map((feature) => (
                      <Badge key={feature} variant="outline" className="text-xs">
                        {feature}
                      </Badge>
                    ))}
                  </div>
                </div>

                <div className="text-xs text-muted-foreground">
                  Processing: {level.processingTime}
                </div>

                {!level.isCompleted && level.isAvailable && (
                  <Button 
                    className="w-full" 
                    size="sm"
                    onClick={() => setCurrentLevel(level.level)}
                  >
                    {level.level === currentLevel ? 'Continue' : 'Start Level ' + level.level}
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* Document Upload Section */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            Document Verification
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            {documents.map((document) => (
              <div key={document.id} className="p-4 border rounded-lg">
                <div className="flex items-center justify-between mb-3">
                  <div className="flex items-center gap-2">
                    {getDocumentIcon(document.type)}
                    <span className="font-medium">{document.name}</span>
                  </div>
                  <div className="flex items-center gap-2">
                    {getStatusIcon(document.status)}
                    <Badge variant={getStatusBadgeVariant(document.status)} className="text-xs">
                      {document.status.replace('_', ' ')}
                    </Badge>
                  </div>
                </div>

                <p className="text-sm text-muted-foreground mb-3">
                  {document.description}
                </p>

                {document.status === 'not_uploaded' && (
                  <Button 
                    variant="outline" 
                    size="sm" 
                    className="w-full"
                    onClick={() => handleFileUpload(document.id)}
                  >
                    <Upload className="w-3 h-3 mr-2" />
                    Upload Document
                  </Button>
                )}

                {document.status === 'uploaded' && (
                  <div className="space-y-2">
                    <div className="flex items-center justify-between text-xs text-muted-foreground">
                      <span>{document.fileName}</span>
                      <span>{document.fileSize}</span>
                    </div>
                    <div className="flex gap-2">
                      <Button variant="outline" size="sm" className="flex-1">
                        <Eye className="w-3 h-3 mr-1" />
                        View
                      </Button>
                      <Button variant="outline" size="sm" className="flex-1">
                        <Upload className="w-3 h-3 mr-1" />
                        Replace
                      </Button>
                    </div>
                  </div>
                )}

                {document.status === 'under_review' && (
                  <div className="text-center py-2">
                    <Clock className="w-4 h-4 mx-auto mb-1 text-blue-500" />
                    <p className="text-xs text-muted-foreground">Under review</p>
                  </div>
                )}

                {document.status === 'approved' && (
                  <div className="text-center py-2">
                    <CheckCircle className="w-4 h-4 mx-auto mb-1 text-green-500" />
                    <p className="text-xs text-green-600">Approved</p>
                  </div>
                )}

                {document.status === 'rejected' && (
                  <div className="space-y-2">
                    <Alert className="border-red-200">
                      <AlertTriangle className="h-4 w-4" />
                      <AlertDescription className="text-xs">
                        {document.rejectionReason || 'Document rejected. Please upload a clearer image.'}
                      </AlertDescription>
                    </Alert>
                    <Button variant="outline" size="sm" className="w-full">
                      <Upload className="w-3 h-3 mr-2" />
                      Re-upload Document
                    </Button>
                  </div>
                )}

                {document.uploadedAt && (
                  <div className="text-xs text-muted-foreground mt-2">
                    Uploaded: {formatTime(document.uploadedAt)}
                  </div>
                )}
              </div>
            ))}
          </div>
        </CardContent>
      </Card>

      {/* Personal Information Form */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <User className="w-5 h-5" />
            Personal Information
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <Label htmlFor="firstName">First Name</Label>
              <Input
                id="firstName"
                value={personalInfo.firstName}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, firstName: e.target.value }))}
                placeholder="Enter your first name"
              />
            </div>

            <div>
              <Label htmlFor="lastName">Last Name</Label>
              <Input
                id="lastName"
                value={personalInfo.lastName}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, lastName: e.target.value }))}
                placeholder="Enter your last name"
              />
            </div>

            <div>
              <Label htmlFor="dateOfBirth">Date of Birth</Label>
              <Input
                id="dateOfBirth"
                type="date"
                value={personalInfo.dateOfBirth}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, dateOfBirth: e.target.value }))}
              />
            </div>

            <div>
              <Label htmlFor="nationality">Nationality</Label>
              <Select 
                value={personalInfo.nationality} 
                onValueChange={(value) => setPersonalInfo(prev => ({ ...prev, nationality: value }))}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select nationality" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="US">United States</SelectItem>
                  <SelectItem value="UK">United Kingdom</SelectItem>
                  <SelectItem value="CA">Canada</SelectItem>
                  <SelectItem value="AU">Australia</SelectItem>
                  <SelectItem value="DE">Germany</SelectItem>
                  <SelectItem value="FR">France</SelectItem>
                </SelectContent>
              </Select>
            </div>

            <div>
              <Label htmlFor="phoneNumber">Phone Number</Label>
              <Input
                id="phoneNumber"
                value={personalInfo.phoneNumber}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, phoneNumber: e.target.value }))}
                placeholder="+1 (555) 123-4567"
              />
            </div>

            <div>
              <Label htmlFor="email">Email Address</Label>
              <Input
                id="email"
                type="email"
                value={personalInfo.email}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, email: e.target.value }))}
                placeholder="your.email@example.com"
              />
            </div>

            <div className="md:col-span-2">
              <Label htmlFor="street">Street Address</Label>
              <Input
                id="street"
                value={personalInfo.address.street}
                onChange={(e) => setPersonalInfo(prev => ({ 
                  ...prev, 
                  address: { ...prev.address, street: e.target.value }
                }))}
                placeholder="123 Main Street"
              />
            </div>

            <div>
              <Label htmlFor="city">City</Label>
              <Input
                id="city"
                value={personalInfo.address.city}
                onChange={(e) => setPersonalInfo(prev => ({ 
                  ...prev, 
                  address: { ...prev.address, city: e.target.value }
                }))}
                placeholder="New York"
              />
            </div>

            <div>
              <Label htmlFor="state">State/Province</Label>
              <Input
                id="state"
                value={personalInfo.address.state}
                onChange={(e) => setPersonalInfo(prev => ({ 
                  ...prev, 
                  address: { ...prev.address, state: e.target.value }
                }))}
                placeholder="NY"
              />
            </div>

            <div>
              <Label htmlFor="occupation">Occupation</Label>
              <Input
                id="occupation"
                value={personalInfo.occupation}
                onChange={(e) => setPersonalInfo(prev => ({ ...prev, occupation: e.target.value }))}
                placeholder="Software Engineer"
              />
            </div>

            <div>
              <Label htmlFor="sourceOfFunds">Source of Funds</Label>
              <Select 
                value={personalInfo.sourceOfFunds} 
                onValueChange={(value) => setPersonalInfo(prev => ({ ...prev, sourceOfFunds: value }))}
              >
                <SelectTrigger>
                  <SelectValue placeholder="Select source" />
                </SelectTrigger>
                <SelectContent>
                  <SelectItem value="salary">Salary/Employment</SelectItem>
                  <SelectItem value="business">Business Income</SelectItem>
                  <SelectItem value="investment">Investment Returns</SelectItem>
                  <SelectItem value="inheritance">Inheritance</SelectItem>
                  <SelectItem value="savings">Personal Savings</SelectItem>
                  <SelectItem value="other">Other</SelectItem>
                </SelectContent>
              </Select>
            </div>
          </div>

          <div className="flex gap-4 mt-6">
            <Button className="flex-1">
              <CheckCircle className="w-4 h-4 mr-2" />
              Save Information
            </Button>
            <Button variant="outline" className="flex-1">
              <Download className="w-4 h-4 mr-2" />
              Download Data
            </Button>
          </div>
        </CardContent>
      </Card>

      {/* Verification Tips */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Zap className="w-5 h-5" />
            Verification Tips
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
            <div>
              <h4 className="font-medium mb-2">Document Requirements</h4>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>• High-quality, clear images</li>
                <li>• All corners visible</li>
                <li>• No glare or shadows</li>
                <li>• Documents must be valid and not expired</li>
                <li>• Information must match your profile</li>
              </ul>
            </div>
            <div>
              <h4 className="font-medium mb-2">Processing Times</h4>
              <ul className="space-y-1 text-sm text-muted-foreground">
                <li>• Level 1: Instant verification</li>
                <li>• Level 2: 1-2 business days</li>
                <li>• Level 3: 3-5 business days</li>
                <li>• Rejected documents: Re-upload immediately</li>
                <li>• Support available 24/7</li>
              </ul>
            </div>
          </div>
        </CardContent>
      </Card>
    </div>
  )
}
