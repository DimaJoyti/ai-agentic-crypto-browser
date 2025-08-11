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
  Rocket, 
  Upload,
  FileText,
  Users,
  DollarSign,
  Calendar,
  Shield,
  Globe,
  Twitter,
  MessageCircle,
  Github,
  CheckCircle,
  AlertTriangle,
  Clock,
  Target,
  BarChart3,
  Award,
  Building2,
  Mail,
  Phone,
  MapPin,
  Plus,
  Trash2,
  Eye
} from 'lucide-react'
import { cn } from '@/lib/utils'
import { useAccount } from 'wagmi'

interface ProjectFormData {
  // Basic Information
  projectName: string
  tokenSymbol: string
  description: string
  longDescription: string
  category: string
  website: string
  whitepaper: string
  
  // Social Links
  twitter: string
  telegram: string
  discord: string
  github: string
  medium: string
  
  // Token Economics
  totalSupply: number
  tokenPrice: number
  tokensForSale: number
  softCap: number
  hardCap: number
  minAllocation: number
  maxAllocation: number
  
  // Launch Details
  launchType: string
  saleStartDate: string
  saleEndDate: string
  vestingPeriod: number
  tgeUnlock: number
  
  // Team Information
  teamMembers: Array<{
    name: string
    role: string
    linkedin: string
    experience: string
  }>
  
  // Tokenomics
  tokenDistribution: Array<{
    category: string
    percentage: number
    description: string
    vestingMonths: number
  }>
  
  // Legal & Compliance
  jurisdiction: string
  kycRequired: boolean
  accreditedOnly: boolean
  restrictedCountries: string[]
  
  // Documents
  documents: Array<{
    type: string
    name: string
    uploaded: boolean
    required: boolean
  }>
}

interface SubmissionStatus {
  step: number
  completionPercentage: number
  status: 'draft' | 'submitted' | 'under_review' | 'approved' | 'rejected'
  submittedAt?: number
  reviewNotes?: string[]
  nextSteps?: string[]
}

export function ProjectSubmission() {
  const [formData, setFormData] = useState<ProjectFormData>({
    projectName: '',
    tokenSymbol: '',
    description: '',
    longDescription: '',
    category: '',
    website: '',
    whitepaper: '',
    twitter: '',
    telegram: '',
    discord: '',
    github: '',
    medium: '',
    totalSupply: 0,
    tokenPrice: 0,
    tokensForSale: 0,
    softCap: 0,
    hardCap: 0,
    minAllocation: 0,
    maxAllocation: 0,
    launchType: '',
    saleStartDate: '',
    saleEndDate: '',
    vestingPeriod: 0,
    tgeUnlock: 0,
    teamMembers: [{ name: '', role: '', linkedin: '', experience: '' }],
    tokenDistribution: [
      { category: 'Public Sale', percentage: 0, description: '', vestingMonths: 0 },
      { category: 'Team', percentage: 0, description: '', vestingMonths: 0 },
      { category: 'Advisors', percentage: 0, description: '', vestingMonths: 0 },
      { category: 'Development', percentage: 0, description: '', vestingMonths: 0 },
      { category: 'Marketing', percentage: 0, description: '', vestingMonths: 0 },
      { category: 'Liquidity', percentage: 0, description: '', vestingMonths: 0 }
    ],
    jurisdiction: '',
    kycRequired: true,
    accreditedOnly: false,
    restrictedCountries: [],
    documents: [
      { type: 'whitepaper', name: 'Whitepaper', uploaded: false, required: true },
      { type: 'pitch_deck', name: 'Pitch Deck', uploaded: false, required: true },
      { type: 'tokenomics', name: 'Tokenomics Document', uploaded: false, required: true },
      { type: 'legal_opinion', name: 'Legal Opinion', uploaded: false, required: true },
      { type: 'audit_report', name: 'Smart Contract Audit', uploaded: false, required: false },
      { type: 'team_kyc', name: 'Team KYC Documents', uploaded: false, required: true }
    ]
  })

  const [submissionStatus, setSubmissionStatus] = useState<SubmissionStatus>({
    step: 1,
    completionPercentage: 15,
    status: 'draft'
  })

  const [currentStep, setCurrentStep] = useState(1)
  const [agreedToTerms, setAgreedToTerms] = useState(false)

  const { address, isConnected } = useAccount()

  const steps = [
    { id: 1, name: 'Basic Info', icon: FileText },
    { id: 2, name: 'Token Economics', icon: DollarSign },
    { id: 3, name: 'Team & Tokenomics', icon: Users },
    { id: 4, name: 'Legal & Compliance', icon: Shield },
    { id: 5, name: 'Documents', icon: Upload },
    { id: 6, name: 'Review & Submit', icon: CheckCircle }
  ]

  const categories = [
    'DeFi', 'Gaming', 'NFT', 'Metaverse', 'Infrastructure', 'DAO', 'Social', 'AI/ML'
  ]

  const launchTypes = [
    { value: 'ido', label: 'IDO (Initial DEX Offering)' },
    { value: 'ino', label: 'INO (Initial NFT Offering)' },
    { value: 'fair_launch', label: 'Fair Launch' },
    { value: 'dutch_auction', label: 'Dutch Auction' },
    { value: 'lottery', label: 'Lottery System' }
  ]

  const updateFormData = (field: string, value: any) => {
    setFormData(prev => ({ ...prev, [field]: value }))
  }

  const addTeamMember = () => {
    setFormData(prev => ({
      ...prev,
      teamMembers: [...prev.teamMembers, { name: '', role: '', linkedin: '', experience: '' }]
    }))
  }

  const removeTeamMember = (index: number) => {
    setFormData(prev => ({
      ...prev,
      teamMembers: prev.teamMembers.filter((_, i) => i !== index)
    }))
  }

  const updateTeamMember = (index: number, field: string, value: string) => {
    setFormData(prev => ({
      ...prev,
      teamMembers: prev.teamMembers.map((member, i) => 
        i === index ? { ...member, [field]: value } : member
      )
    }))
  }

  const updateTokenDistribution = (index: number, field: string, value: any) => {
    setFormData(prev => ({
      ...prev,
      tokenDistribution: prev.tokenDistribution.map((item, i) => 
        i === index ? { ...item, [field]: value } : item
      )
    }))
  }

  const getTotalDistribution = () => {
    return formData.tokenDistribution.reduce((sum, item) => sum + item.percentage, 0)
  }

  const handleDocumentUpload = (index: number) => {
    setFormData(prev => ({
      ...prev,
      documents: prev.documents.map((doc, i) => 
        i === index ? { ...doc, uploaded: true } : doc
      )
    }))
  }

  const calculateCompletionPercentage = () => {
    let completed = 0
    let total = 0

    // Basic info (20%)
    total += 20
    if (formData.projectName && formData.tokenSymbol && formData.description && formData.category) {
      completed += 20
    }

    // Token economics (20%)
    total += 20
    if (formData.totalSupply && formData.tokenPrice && formData.hardCap) {
      completed += 20
    }

    // Team & tokenomics (20%)
    total += 20
    if (formData.teamMembers.some(m => m.name && m.role) && getTotalDistribution() === 100) {
      completed += 20
    }

    // Legal & compliance (20%)
    total += 20
    if (formData.jurisdiction && formData.launchType) {
      completed += 20
    }

    // Documents (20%)
    total += 20
    const requiredDocs = formData.documents.filter(d => d.required)
    const uploadedRequiredDocs = requiredDocs.filter(d => d.uploaded)
    if (uploadedRequiredDocs.length === requiredDocs.length) {
      completed += 20
    }

    return Math.round((completed / total) * 100)
  }

  useEffect(() => {
    const percentage = calculateCompletionPercentage()
    setSubmissionStatus(prev => ({ ...prev, completionPercentage: percentage }))
  }, [formData])

  const handleSubmit = () => {
    setSubmissionStatus(prev => ({
      ...prev,
      status: 'submitted',
      submittedAt: Date.now(),
      nextSteps: [
        'Initial review by our team (1-2 business days)',
        'Due diligence process (3-5 business days)',
        'Final approval and launch scheduling'
      ]
    }))
  }

  if (!isConnected) {
    return (
      <Card>
        <CardContent className="p-8 text-center">
          <Rocket className="w-12 h-12 mx-auto mb-4 text-muted-foreground opacity-50" />
          <h3 className="text-lg font-medium mb-2">Connect Wallet Required</h3>
          <p className="text-muted-foreground">
            Connect your wallet to submit your project for launchpad consideration
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
          <h2 className="text-2xl font-bold">Submit Your Project</h2>
          <p className="text-muted-foreground">
            Apply to launch your token on our launchpad platform
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Badge variant="outline">
            <Clock className="w-3 h-3 mr-1" />
            {submissionStatus.status}
          </Badge>
          <Badge variant="outline">
            {submissionStatus.completionPercentage}% Complete
          </Badge>
        </div>
      </div>

      {/* Progress Overview */}
      <Card>
        <CardHeader>
          <CardTitle className="flex items-center gap-2">
            <Target className="w-5 h-5" />
            Application Progress
          </CardTitle>
        </CardHeader>
        <CardContent>
          <div className="space-y-4">
            <div className="flex items-center justify-between">
              <span className="text-sm text-muted-foreground">Overall Completion</span>
              <span className="font-medium">{submissionStatus.completionPercentage}%</span>
            </div>
            <Progress value={submissionStatus.completionPercentage} className="h-3" />
            
            <div className="grid grid-cols-2 md:grid-cols-6 gap-2">
              {steps.map((step) => {
                const Icon = step.icon
                return (
                  <div
                    key={step.id}
                    className={cn(
                      "p-3 border rounded-lg text-center cursor-pointer transition-colors",
                      currentStep === step.id ? "border-primary bg-primary/5" : "border-muted"
                    )}
                    onClick={() => setCurrentStep(step.id)}
                  >
                    <Icon className={cn(
                      "w-5 h-5 mx-auto mb-1",
                      currentStep === step.id ? "text-primary" : "text-muted-foreground"
                    )} />
                    <div className="text-xs font-medium">{step.name}</div>
                  </div>
                )
              })}
            </div>
          </div>
        </CardContent>
      </Card>

      {/* Form Steps */}
      <Card>
        <CardHeader>
          <CardTitle>
            Step {currentStep}: {steps.find(s => s.id === currentStep)?.name}
          </CardTitle>
        </CardHeader>
        <CardContent>
          {/* Step 1: Basic Information */}
          {currentStep === 1 && (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="projectName">Project Name *</Label>
                  <Input
                    id="projectName"
                    value={formData.projectName}
                    onChange={(e) => updateFormData('projectName', e.target.value)}
                    placeholder="DeFi Protocol"
                  />
                </div>

                <div>
                  <Label htmlFor="tokenSymbol">Token Symbol *</Label>
                  <Input
                    id="tokenSymbol"
                    value={formData.tokenSymbol}
                    onChange={(e) => updateFormData('tokenSymbol', e.target.value)}
                    placeholder="DEFI"
                    maxLength={10}
                  />
                </div>

                <div>
                  <Label htmlFor="category">Category *</Label>
                  <Select value={formData.category} onValueChange={(value) => updateFormData('category', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select category" />
                    </SelectTrigger>
                    <SelectContent>
                      {categories.map((cat) => (
                        <SelectItem key={cat} value={cat.toLowerCase()}>{cat}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label htmlFor="website">Website *</Label>
                  <Input
                    id="website"
                    value={formData.website}
                    onChange={(e) => updateFormData('website', e.target.value)}
                    placeholder="https://yourproject.com"
                  />
                </div>
              </div>

              <div>
                <Label htmlFor="description">Short Description *</Label>
                <textarea
                  id="description"
                  value={formData.description}
                  onChange={(e) => updateFormData('description', e.target.value)}
                  placeholder="Brief description of your project (max 200 characters)"
                  maxLength={200}
                  rows={3}
                  className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                />
                <div className="text-xs text-muted-foreground mt-1">
                  {formData.description.length}/200 characters
                </div>
              </div>

              <div>
                <Label htmlFor="longDescription">Detailed Description *</Label>
                <textarea
                  id="longDescription"
                  value={formData.longDescription}
                  onChange={(e) => updateFormData('longDescription', e.target.value)}
                  placeholder="Detailed description of your project, technology, and vision"
                  rows={6}
                  className="w-full px-3 py-2 border rounded-lg focus:outline-none focus:ring-2 focus:ring-primary"
                />
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="twitter">Twitter</Label>
                  <Input
                    id="twitter"
                    value={formData.twitter}
                    onChange={(e) => updateFormData('twitter', e.target.value)}
                    placeholder="@yourproject"
                  />
                </div>

                <div>
                  <Label htmlFor="telegram">Telegram</Label>
                  <Input
                    id="telegram"
                    value={formData.telegram}
                    onChange={(e) => updateFormData('telegram', e.target.value)}
                    placeholder="t.me/yourproject"
                  />
                </div>
              </div>
            </div>
          )}

          {/* Step 2: Token Economics */}
          {currentStep === 2 && (
            <div className="space-y-4">
              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="totalSupply">Total Supply *</Label>
                  <Input
                    id="totalSupply"
                    type="number"
                    value={formData.totalSupply}
                    onChange={(e) => updateFormData('totalSupply', parseFloat(e.target.value))}
                    placeholder="100000000"
                  />
                </div>

                <div>
                  <Label htmlFor="tokenPrice">Token Price (USD) *</Label>
                  <Input
                    id="tokenPrice"
                    type="number"
                    step="0.001"
                    value={formData.tokenPrice}
                    onChange={(e) => updateFormData('tokenPrice', parseFloat(e.target.value))}
                    placeholder="0.10"
                  />
                </div>

                <div>
                  <Label htmlFor="tokensForSale">Tokens for Sale *</Label>
                  <Input
                    id="tokensForSale"
                    type="number"
                    value={formData.tokensForSale}
                    onChange={(e) => updateFormData('tokensForSale', parseFloat(e.target.value))}
                    placeholder="10000000"
                  />
                </div>

                <div>
                  <Label htmlFor="hardCap">Hard Cap (USD) *</Label>
                  <Input
                    id="hardCap"
                    type="number"
                    value={formData.hardCap}
                    onChange={(e) => updateFormData('hardCap', parseFloat(e.target.value))}
                    placeholder="1000000"
                  />
                </div>

                <div>
                  <Label htmlFor="minAllocation">Min Allocation (USD) *</Label>
                  <Input
                    id="minAllocation"
                    type="number"
                    value={formData.minAllocation}
                    onChange={(e) => updateFormData('minAllocation', parseFloat(e.target.value))}
                    placeholder="100"
                  />
                </div>

                <div>
                  <Label htmlFor="maxAllocation">Max Allocation (USD) *</Label>
                  <Input
                    id="maxAllocation"
                    type="number"
                    value={formData.maxAllocation}
                    onChange={(e) => updateFormData('maxAllocation', parseFloat(e.target.value))}
                    placeholder="5000"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="launchType">Launch Type *</Label>
                  <Select value={formData.launchType} onValueChange={(value) => updateFormData('launchType', value)}>
                    <SelectTrigger>
                      <SelectValue placeholder="Select launch type" />
                    </SelectTrigger>
                    <SelectContent>
                      {launchTypes.map((type) => (
                        <SelectItem key={type.value} value={type.value}>{type.label}</SelectItem>
                      ))}
                    </SelectContent>
                  </Select>
                </div>

                <div>
                  <Label htmlFor="tgeUnlock">TGE Unlock (%)</Label>
                  <Input
                    id="tgeUnlock"
                    type="number"
                    min="0"
                    max="100"
                    value={formData.tgeUnlock}
                    onChange={(e) => updateFormData('tgeUnlock', parseFloat(e.target.value))}
                    placeholder="20"
                  />
                </div>
              </div>

              <div className="grid grid-cols-1 md:grid-cols-2 gap-4">
                <div>
                  <Label htmlFor="saleStartDate">Sale Start Date</Label>
                  <Input
                    id="saleStartDate"
                    type="datetime-local"
                    value={formData.saleStartDate}
                    onChange={(e) => updateFormData('saleStartDate', e.target.value)}
                  />
                </div>

                <div>
                  <Label htmlFor="saleEndDate">Sale End Date</Label>
                  <Input
                    id="saleEndDate"
                    type="datetime-local"
                    value={formData.saleEndDate}
                    onChange={(e) => updateFormData('saleEndDate', e.target.value)}
                  />
                </div>
              </div>
            </div>
          )}

          {/* Navigation Buttons */}
          <div className="flex justify-between mt-6">
            <Button 
              variant="outline" 
              onClick={() => setCurrentStep(Math.max(1, currentStep - 1))}
              disabled={currentStep === 1}
            >
              Previous
            </Button>
            
            {currentStep < 6 ? (
              <Button onClick={() => setCurrentStep(Math.min(6, currentStep + 1))}>
                Next
              </Button>
            ) : (
              <Button 
                onClick={handleSubmit}
                disabled={!agreedToTerms || submissionStatus.completionPercentage < 80}
              >
                Submit Application
              </Button>
            )}
          </div>
        </CardContent>
      </Card>

      {/* Submission Status */}
      {submissionStatus.status === 'submitted' && (
        <Alert>
          <CheckCircle className="h-4 w-4" />
          <AlertDescription>
            Your project application has been submitted successfully! Our team will review it within 1-2 business days.
          </AlertDescription>
        </Alert>
      )}
    </div>
  )
}
