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
import { Checkbox } from '@/components/ui/checkbox'
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
  Zap,
  Award,
  Lock
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface KYCLevel {
  level: number
  name: string
  description: string
  requirements: string[]
  benefits: string[]
  maxAllocation: number
  processingTime: string
  isCompleted: boolean
  isAvailable: boolean
}

interface KYCDocument {
  id: string
  type: 'passport' | 'drivers_license' | 'national_id' | 'utility_bill' | 'bank_statement' | 'selfie' | 'proof_of_funds'
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
  investmentExperience: string
  riskTolerance: string
}

interface KYCStatus {
  level: number
  status: 'not_started' | 'in_progress' | 'under_review' | 'approved' | 'rejected'
  completionPercentage: number
  submittedAt?: number
  approvedAt?: number
  rejectedAt?: number
  rejectionReason?: string
  nextSteps: string[]
  eligibleProjects: string[]
}

export function LaunchpadKYC() {
  const [kycLevels, setKycLevels] = useState<KYCLevel[]>([])
  const [documents, setDocuments] = useState<KYCDocument[]>([])
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
    sourceOfFunds: '',
    investmentExperience: '',
    riskTolerance: ''
  })
  const [kycStatus, setKycStatus] = useState<KYCStatus>({
    level: 1,
    status: 'in_progress',
    completionPercentage: 65,
    nextSteps: ['Upload proof of address', 'Complete investment profile'],
    eligibleProjects: ['project1', 'project2']
  })
  const [currentStep, setCurrentStep] = useState(1)
  const [agreedToTerms, setAgreedToTerms] = useState(false)

  const { address, isConnected } = useAccount()

  useEffect(() => {
    if (!isConnected) return

    // Generate mock KYC data
    const mockKycLevels: KYCLevel[] = [
      {
        level: 1,
        name: 'Basic KYC',
        description: 'Basic verification for small allocations',
        requirements: ['Email verification', 'Phone verification', 'Basic personal information'],
        benefits: ['Access to public sales', 'Up to $1,000 allocation', 'Standard processing'],
        maxAllocation: 1000,
        processingTime: 'Instant',
        isCompleted: true,
        isAvailable: true
      },
      {
        level: 2,
        name: 'Enhanced KYC',
        description: 'Enhanced verification for medium allocations',
        requirements: ['Government ID', 'Selfie verification', 'Address verification', 'Investment profile'],
        benefits: ['Access to all public sales', 'Up to $10,000 allocation', 'Priority support'],
        maxAllocation: 10000,
        processingTime: '1-2 business days',
        isCompleted: false,
        isAvailable: true
      },
      {
        level: 3,
        name: 'Institutional KYC',
        description: 'Institutional verification for large allocations',
        requirements: ['Enhanced due diligence', 'Source of funds verification', 'Institutional documentation'],
        benefits: ['Access to private sales', 'Unlimited allocation', 'Dedicated support', 'Early access'],
        maxAllocation: 1000000,
        processingTime: '3-5 business days',
        isCompleted: false,
        isAvailable: false
      }
    ]

    const mockDocuments: KYCDocument[] = [
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
        id: 'proof_of_funds',
        type: 'proof_of_funds',
        name: 'Proof of Funds',
        description: 'Bank statement or proof of cryptocurrency holdings',
        status: 'not_uploaded',
        required: false
      }
    ]

    setKycLevels(mockKycLevels)
    setDocuments(mockDocuments)
  }, [isConnected])

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
      case 'utility_bill': case 'bank_statement': case 'proof_of_funds': return <Building2 className="w-4 h-4" />
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

  const handlePersonalInfoSubmit = () => {
    // Handle personal info submission
    setCurrentStep(2)
  }

  const handleDocumentSubmit = () => {
    // Handle document submission
    setCurrentStep(3)
  }

  const handleKYCSubmit = () => {
    // Handle final KYC submission
    setKycStatus(prev => ({ ...prev, status: 'under_review', submittedAt: Date.now() }))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Shield className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to access KYC verification for launchpad participation
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
          <h2 className="text-2xl font-bold">Launchpad KYC Verification</h2>
          <p className="text-muted-foreground">
            Complete KYC verification to participate in token launches
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Shield className="w-3 h-3 mr-1" />
            Level {kycStatus.level}
          </Badge>
          <Badge variant={getStatusBadgeVariant(kycStatus.status)}>
            {kycStatus.status.replace('_', ' ')}
          </Badge>
        </div>
      </div>

      {/* KYC Progress */}
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
              <span className="font-medium">{kycStatus.completionPercentage}%</span>
            </div>
            <Progress value={kycStatus.completionPercentage} className="h-3" />
            
            {kycStatus.nextSteps.length > 0 && (
              <div>
                <h4 className="font-medium mb-2">Next Steps:</h4>
                <ul className="space-y-1">
                  {kycStatus.nextSteps.map((step, index) => (
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
            level.level === kycStatus.level && "ring-2 ring-primary",
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
                  <h4 className="font-medium mb-2">Max Allocation</h4>
                  <div className="text-2xl font-bold text-primary">
                    ${level.maxAllocation.toLocaleString()}
                  </div>
                </div>

                <div>
                  <h4 className="font-medium mb-2">Benefits</h4>
                  <ul className="space-y-1">
                    {level.benefits.map((benefit, index) => (
                      <li key={index} className="text-sm text-muted-foreground flex items-center gap-2">
                        <CheckCircle className="w-3 h-3 text-green-500" />
                        {benefit}
                      </li>
                    ))}
                  </ul>
                </div>

                <div className="text-xs text-muted-foreground">
                  Processing: {level.processingTime}
                </div>

                {!level.isCompleted && level.isAvailable && (
                  <Button 
                    className="w-full" 
                    size="sm"
                    variant={level.level === kycStatus.level ? "default" : "outline"}
                  >
                    {level.level === kycStatus.level ? 'Continue' : 'Upgrade to Level ' + level.level}
                  </Button>
                )}
              </div>
            </CardContent>
          </Card>
        ))}
      </div>

      {/* KYC Form Steps */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <FileText className="w-5 h-5" />
            KYC Application
          </CardTitle>
        </CardHeader>
        <CardContent>
          {/* Step Indicator */}
          <div className="flex items-center justify-between mb-6">
            {[1, 2, 3].map((step) => (
              <div key={step} className="flex items-center">
                <div className={cn(
                  "w-8 h-8 rounded-full flex items-center justify-center text-sm font-medium",
                  step <= currentStep ? "bg-primary text-primary-foreground" : "bg-muted text-muted-foreground"
                )}>
                  {step}
                </div>
                {step < 3 && (
                  <div className={cn(
                    "w-16 h-1 mx-2",
                    step < currentStep ? "bg-primary" : "bg-muted"
                  )} />
                )}
              </div>
            ))}
          </div>

          {/* Step 1: Personal Information */}
          {currentStep === 1 && (
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Personal Information</h3>
              
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
                      <SelectItem value="crypto">Cryptocurrency</SelectItem>
                    </SelectContent>
                  </Select>
                </div>
              </div>

              <Button onClick={handlePersonalInfoSubmit} className="w-full">
                Continue to Document Upload
              </Button>
            </div>
          )}

          {/* Step 2: Document Upload */}
          {currentStep === 2 && (
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Document Verification</h3>
              
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

                    {document.uploadedAt && (
                      <div className="text-xs text-muted-foreground mt-2">
                        Uploaded: {formatTime(document.uploadedAt)}
                      </div>
                    )}
                  </div>
                ))}
              </div>

              <Button onClick={handleDocumentSubmit} className="w-full">
                Continue to Review
              </Button>
            </div>
          )}

          {/* Step 3: Review and Submit */}
          {currentStep === 3 && (
            <div className="space-y-4">
              <h3 className="text-lg font-medium">Review and Submit</h3>
              
              <Alert>
                <Shield className="h-4 w-4" />
                <AlertDescription>
                  Please review all information carefully before submitting. Once submitted, your application will be reviewed within 1-2 business days.
                </AlertDescription>
              </Alert>

              <div className="space-y-4">
                <div className="p-4 border rounded-lg">
                  <h4 className="font-medium mb-2">Personal Information</h4>
                  <div className="grid grid-cols-2 gap-2 text-sm">
                    <div>Name: {personalInfo.firstName} {personalInfo.lastName}</div>
                    <div>Nationality: {personalInfo.nationality}</div>
                    <div>Occupation: {personalInfo.occupation}</div>
                    <div>Source of Funds: {personalInfo.sourceOfFunds}</div>
                  </div>
                </div>

                <div className="p-4 border rounded-lg">
                  <h4 className="font-medium mb-2">Documents</h4>
                  <div className="space-y-1">
                    {documents.filter(d => d.status !== 'not_uploaded').map((doc) => (
                      <div key={doc.id} className="flex items-center justify-between text-sm">
                        <span>{doc.name}</span>
                        <Badge variant={getStatusBadgeVariant(doc.status)} className="text-xs">
                          {doc.status.replace('_', ' ')}
                        </Badge>
                      </div>
                    ))}
                  </div>
                </div>

                <div className="flex items-center space-x-2">
                  <Checkbox 
                    id="terms" 
                    checked={agreedToTerms}
                    onCheckedChange={(checked) => setAgreedToTerms(checked as boolean)}
                  />
                  <Label htmlFor="terms" className="text-sm">
                    I agree to the Terms of Service and Privacy Policy
                  </Label>
                </div>

                <Button 
                  onClick={handleKYCSubmit} 
                  className="w-full"
                  disabled={!agreedToTerms}
                >
                  <CheckCircle className="w-4 h-4 mr-2" />
                  Submit KYC Application
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>

      {/* Eligible Projects */}
      {kycStatus.eligibleProjects.length > 0 && (
        <Card>
          <CardHeader>
            <CardTitle className="flex items-center gap-2">
              <Award className="w-5 h-5" />
              Eligible Projects
            </CardTitle>
          </CardHeader>
          <CardContent>
            <p className="text-sm text-muted-foreground mb-4">
              Based on your current KYC level, you are eligible to participate in the following projects:
            </p>
            <div className="space-y-2">
              {kycStatus.eligibleProjects.map((projectId) => (
                <div key={projectId} className="flex items-center justify-between p-3 border rounded">
                  <div>
                    <div className="font-medium">Project {projectId}</div>
                    <div className="text-sm text-muted-foreground">Max allocation: $1,000</div>
                  </div>
                  <Button variant="outline" size="sm">
                    View Project
                  </Button>
                </div>
              ))}
            </div>
          </CardContent>
        </Card>
      )}
    </div>
  )
}
