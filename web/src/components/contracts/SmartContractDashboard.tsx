'use client'

import { useState } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import { Input } from '@/components/ui/input'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Alert, AlertDescription } from '@/components/ui/alert'
import { 
  Code,
  Play,
  Upload,
  Search,
  Settings,
  AlertTriangle,
  CheckCircle,
  Clock,
  Zap,
  Shield,
  FileCode,
  Database,
  Activity,
  RefreshCw,
  ExternalLink,
  Copy,
  Eye,
  EyeOff
} from 'lucide-react'
import { useSmartContract, useContractRegistry } from '@/hooks/useSmartContract'
import { type SmartContract, type ContractCall, type ContractFunction } from '@/lib/smart-contract-integration'
import { cn } from '@/lib/utils'

export function SmartContractDashboard() {
  const [activeTab, setActiveTab] = useState('contracts')
  const [searchQuery, setSearchQuery] = useState('')
  const [selectedContract, setSelectedContract] = useState<SmartContract | null>(null)
  const [showABI, setShowABI] = useState(false)

  const {
    state,
    prepareCall,
    executeCall,
    clearError
  } = useSmartContract({
    autoLoad: true,
    enableNotifications: true,
    enableSimulation: true
  })

  const { searchContracts } = useContractRegistry()

  const filteredContracts = searchQuery 
    ? searchContracts(searchQuery)
    : state.contracts

  const getRiskLevelColor = (riskLevel: string) => {
    switch (riskLevel) {
      case 'low':
        return 'text-green-600 bg-green-100 dark:bg-green-900'
      case 'medium':
        return 'text-yellow-600 bg-yellow-100 dark:bg-yellow-900'
      case 'high':
        return 'text-red-600 bg-red-100 dark:bg-red-900'
      case 'critical':
        return 'text-red-700 bg-red-200 dark:bg-red-800'
      default:
        return 'text-gray-600 bg-gray-100 dark:bg-gray-900'
    }
  }

  const getStateMutabilityIcon = (stateMutability: string) => {
    switch (stateMutability) {
      case 'view':
      case 'pure':
        return <Eye className="w-4 h-4 text-blue-500" />
      case 'payable':
        return <Zap className="w-4 h-4 text-yellow-500" />
      case 'nonpayable':
        return <Activity className="w-4 h-4 text-green-500" />
      default:
        return <Code className="w-4 h-4 text-gray-500" />
    }
  }

  const formatAddress = (address: string) => {
    return `${address.slice(0, 6)}...${address.slice(-4)}`
  }

  const copyToClipboard = (text: string) => {
    navigator.clipboard.writeText(text)
  }

  const handleFunctionCall = async (func: ContractFunction, args: any[]) => {
    if (!selectedContract) return

    try {
      const call = await prepareCall(selectedContract.address, func.name, args)
      await executeCall(call)
    } catch (error) {
      console.error('Function call failed:', error)
    }
  }

  return (
    <div className="space-y-6">
      <div className="flex items-center justify-between">
        <div>
          <h2 className="text-2xl font-bold">Smart Contracts</h2>
          <p className="text-muted-foreground">
            Manage and interact with smart contracts across multiple chains
          </p>
        </div>
        <div className="flex items-center gap-2">
          <Button variant="outline" size="sm">
            <Upload className="w-4 h-4 mr-2" />
            Import ABI
          </Button>
          <Button variant="outline" size="sm">
            <Code className="w-4 h-4 mr-2" />
            Deploy Contract
          </Button>
          <Button variant="outline" size="sm">
            <Settings className="w-4 h-4 mr-2" />
            Settings
          </Button>
        </div>
      </div>

      {/* Error Alert */}
      {state.error && (
        <Alert variant="destructive">
          <AlertTriangle className="h-4 w-4" />
          <AlertDescription>
            {state.error}
            <Button variant="ghost" size="sm" onClick={clearError} className="ml-2">
              Dismiss
            </Button>
          </AlertDescription>
        </Alert>
      )}

      {/* Search */}
      <div className="flex items-center gap-4">
        <div className="relative flex-1">
          <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 text-muted-foreground w-4 h-4" />
          <Input
            placeholder="Search contracts by name, address, or tags..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="pl-10"
          />
        </div>
        <Button variant="outline" size="sm">
          <RefreshCw className="w-4 h-4 mr-2" />
          Refresh
        </Button>
      </div>

      {/* Statistics */}
      <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-4">
        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Total Contracts</p>
                <p className="text-2xl font-bold">{state.contracts.length}</p>
              </div>
              <FileCode className="w-8 h-8 text-blue-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Registered contracts
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Recent Calls</p>
                <p className="text-2xl font-bold">{state.calls.length}</p>
              </div>
              <Activity className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Function executions
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Deployments</p>
                <p className="text-2xl font-bold">{state.deployments.length}</p>
              </div>
              <Upload className="w-8 h-8 text-orange-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Contract deployments
            </div>
          </CardContent>
        </Card>

        <Card>
          <CardContent className="p-6">
            <div className="flex items-center justify-between">
              <div>
                <p className="text-sm font-medium text-muted-foreground">Verified</p>
                <p className="text-2xl font-bold text-green-600">
                  {state.contracts.filter(c => c.verified).length}
                </p>
              </div>
              <Shield className="w-8 h-8 text-green-500" />
            </div>
            <div className="mt-2 text-sm text-muted-foreground">
              Verified contracts
            </div>
          </CardContent>
        </Card>
      </div>

      {/* Main Content */}
      <Tabs value={activeTab} onValueChange={setActiveTab}>
        <TabsList className="grid w-full grid-cols-4">
          <TabsTrigger value="contracts">Contracts ({filteredContracts.length})</TabsTrigger>
          <TabsTrigger value="calls">Recent Calls ({state.calls.length})</TabsTrigger>
          <TabsTrigger value="deployments">Deployments ({state.deployments.length})</TabsTrigger>
          <TabsTrigger value="interact">Interact</TabsTrigger>
        </TabsList>

        <TabsContent value="contracts" className="space-y-4">
          {filteredContracts.length > 0 ? (
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              <AnimatePresence>
                {filteredContracts.map((contract, index) => (
                  <motion.div
                    key={contract.address}
                    initial={{ opacity: 0, y: 20 }}
                    animate={{ opacity: 1, y: 0 }}
                    exit={{ opacity: 0, y: -20 }}
                    transition={{ delay: index * 0.05 }}
                  >
                    <Card 
                      className={cn(
                        "transition-all duration-200 cursor-pointer hover:shadow-md",
                        selectedContract?.address === contract.address && "ring-2 ring-primary"
                      )}
                      onClick={() => setSelectedContract(contract)}
                    >
                      <CardHeader className="pb-3">
                        <div className="flex items-center justify-between">
                          <CardTitle className="text-lg">{contract.name}</CardTitle>
                          <div className="flex items-center gap-1">
                            {contract.verified && (
                              <Badge variant="default" className="text-xs">
                                <Shield className="w-3 h-3 mr-1" />
                                Verified
                              </Badge>
                            )}
                            {contract.proxy && (
                              <Badge variant="secondary" className="text-xs">
                                Proxy
                              </Badge>
                            )}
                          </div>
                        </div>
                        <CardDescription>
                          {contract.metadata.description || 'No description available'}
                        </CardDescription>
                      </CardHeader>
                      <CardContent>
                        <div className="space-y-3">
                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Address:</span>
                            <div className="flex items-center gap-1">
                              <span className="font-mono">{formatAddress(contract.address)}</span>
                              <Button
                                variant="ghost"
                                size="sm"
                                onClick={(e) => {
                                  e.stopPropagation()
                                  copyToClipboard(contract.address)
                                }}
                              >
                                <Copy className="w-3 h-3" />
                              </Button>
                            </div>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Chain:</span>
                            <span>{contract.chainId}</span>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Functions:</span>
                            <span>{contract.functions.length}</span>
                          </div>

                          <div className="flex items-center justify-between text-sm">
                            <span className="text-muted-foreground">Events:</span>
                            <span>{contract.events.length}</span>
                          </div>

                          {contract.metadata.security && (
                            <div className="flex items-center justify-between text-sm">
                              <span className="text-muted-foreground">Risk Level:</span>
                              <Badge className={cn("text-xs", getRiskLevelColor(contract.metadata.security.riskLevel))}>
                                {contract.metadata.security.riskLevel.toUpperCase()}
                              </Badge>
                            </div>
                          )}

                          <div className="flex flex-wrap gap-1 mt-2">
                            {contract.tags.slice(0, 3).map(tag => (
                              <Badge key={tag} variant="outline" className="text-xs">
                                {tag}
                              </Badge>
                            ))}
                            {contract.tags.length > 3 && (
                              <Badge variant="outline" className="text-xs">
                                +{contract.tags.length - 3}
                              </Badge>
                            )}
                          </div>
                        </div>
                      </CardContent>
                    </Card>
                  </motion.div>
                ))}
              </AnimatePresence>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <FileCode className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Contracts Found</h3>
                <p className="text-muted-foreground mb-4">
                  {searchQuery ? 'No contracts match your search criteria' : 'No contracts registered yet'}
                </p>
                <Button>
                  <Upload className="w-4 h-4 mr-2" />
                  Import Contract
                </Button>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="calls" className="space-y-4">
          {state.calls.length > 0 ? (
            <div className="space-y-3">
              {state.calls.map((call, index) => (
                <motion.div
                  key={call.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <Card>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex items-center gap-2">
                            {call.status === 'success' ? (
                              <CheckCircle className="w-4 h-4 text-green-500" />
                            ) : call.status === 'failed' || call.status === 'reverted' ? (
                              <AlertTriangle className="w-4 h-4 text-red-500" />
                            ) : (
                              <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
                            )}
                          </div>
                          
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <p className="font-medium">{call.functionName}</p>
                              <Badge variant="outline" className="text-xs">
                                {formatAddress(call.contract)}
                              </Badge>
                            </div>
                            
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <span>Status: {call.status}</span>
                              {call.gasUsed && <span>Gas: {call.gasUsed}</span>}
                              <span>{new Date(call.timestamp).toLocaleTimeString()}</span>
                            </div>

                            {call.error && (
                              <p className="text-sm text-red-600 mt-1">
                                Error: {call.error}
                              </p>
                            )}
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          {call.hash && (
                            <Button variant="ghost" size="sm">
                              <ExternalLink className="w-4 h-4" />
                            </Button>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Activity className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Recent Calls</h3>
                <p className="text-muted-foreground">
                  Contract function calls will appear here
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="deployments" className="space-y-4">
          {state.deployments.length > 0 ? (
            <div className="space-y-3">
              {state.deployments.map((deployment, index) => (
                <motion.div
                  key={deployment.id}
                  initial={{ opacity: 0, x: -20 }}
                  animate={{ opacity: 1, x: 0 }}
                  transition={{ delay: index * 0.05 }}
                >
                  <Card>
                    <CardContent className="p-4">
                      <div className="flex items-center justify-between">
                        <div className="flex items-center gap-3">
                          <div className="flex items-center gap-2">
                            {deployment.status === 'success' ? (
                              <CheckCircle className="w-4 h-4 text-green-500" />
                            ) : deployment.status === 'failed' ? (
                              <AlertTriangle className="w-4 h-4 text-red-500" />
                            ) : (
                              <Clock className="w-4 h-4 text-yellow-500 animate-pulse" />
                            )}
                          </div>
                          
                          <div>
                            <div className="flex items-center gap-2 mb-1">
                              <p className="font-medium">{deployment.name}</p>
                              <Badge variant="outline" className="text-xs">
                                Chain {deployment.chainId}
                              </Badge>
                            </div>
                            
                            <div className="flex items-center gap-4 text-sm text-muted-foreground">
                              <span>Status: {deployment.status}</span>
                              {deployment.address && (
                                <span>Address: {formatAddress(deployment.address)}</span>
                              )}
                              <span>{new Date(deployment.timestamp).toLocaleTimeString()}</span>
                            </div>

                            {deployment.error && (
                              <p className="text-sm text-red-600 mt-1">
                                Error: {deployment.error}
                              </p>
                            )}
                          </div>
                        </div>

                        <div className="flex items-center gap-2">
                          {deployment.hash && (
                            <Button variant="ghost" size="sm">
                              <ExternalLink className="w-4 h-4" />
                            </Button>
                          )}
                        </div>
                      </div>
                    </CardContent>
                  </Card>
                </motion.div>
              ))}
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Upload className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">No Deployments</h3>
                <p className="text-muted-foreground mb-4">
                  Contract deployments will appear here
                </p>
                <Button>
                  <Code className="w-4 h-4 mr-2" />
                  Deploy Contract
                </Button>
              </CardContent>
            </Card>
          )}
        </TabsContent>

        <TabsContent value="interact" className="space-y-4">
          {selectedContract ? (
            <div className="grid gap-6 lg:grid-cols-2">
              {/* Contract Info */}
              <Card>
                <CardHeader>
                  <CardTitle className="flex items-center justify-between">
                    {selectedContract.name}
                    <Button
                      variant="outline"
                      size="sm"
                      onClick={() => setShowABI(!showABI)}
                    >
                      {showABI ? <EyeOff className="w-4 h-4" /> : <Eye className="w-4 h-4" />}
                      {showABI ? 'Hide' : 'Show'} ABI
                    </Button>
                  </CardTitle>
                  <CardDescription>
                    {selectedContract.metadata.description}
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3">
                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Address:</span>
                      <div className="flex items-center gap-1">
                        <span className="font-mono">{formatAddress(selectedContract.address)}</span>
                        <Button
                          variant="ghost"
                          size="sm"
                          onClick={() => copyToClipboard(selectedContract.address)}
                        >
                          <Copy className="w-3 h-3" />
                        </Button>
                      </div>
                    </div>

                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Chain ID:</span>
                      <span>{selectedContract.chainId}</span>
                    </div>

                    <div className="flex items-center justify-between text-sm">
                      <span className="text-muted-foreground">Functions:</span>
                      <span>{selectedContract.functions.length}</span>
                    </div>

                    {showABI && (
                      <div className="mt-4">
                        <h4 className="font-medium mb-2">Contract ABI</h4>
                        <div className="bg-muted p-3 rounded-md max-h-40 overflow-y-auto">
                          <pre className="text-xs">
                            {JSON.stringify(selectedContract.abi, null, 2)}
                          </pre>
                        </div>
                      </div>
                    )}
                  </div>
                </CardContent>
              </Card>

              {/* Functions */}
              <Card>
                <CardHeader>
                  <CardTitle>Functions</CardTitle>
                  <CardDescription>
                    Available contract functions
                  </CardDescription>
                </CardHeader>
                <CardContent>
                  <div className="space-y-3 max-h-96 overflow-y-auto">
                    {selectedContract.functions.map((func, index) => (
                      <div
                        key={index}
                        className="flex items-center justify-between p-3 border rounded-lg"
                      >
                        <div className="flex items-center gap-3">
                          {getStateMutabilityIcon(func.stateMutability)}
                          <div>
                            <p className="font-medium">{func.name}</p>
                            <p className="text-xs text-muted-foreground">
                              {func.inputs.length} inputs • {func.outputs.length} outputs
                            </p>
                          </div>
                        </div>
                        <div className="flex items-center gap-2">
                          <Badge className={cn("text-xs", getRiskLevelColor(func.riskLevel))}>
                            {func.riskLevel}
                          </Badge>
                          <Button
                            variant="outline"
                            size="sm"
                            onClick={() => handleFunctionCall(func, [])}
                            disabled={state.isExecuting}
                          >
                            <Play className="w-3 h-3 mr-1" />
                            Call
                          </Button>
                        </div>
                      </div>
                    ))}
                  </div>
                </CardContent>
              </Card>
            </div>
          ) : (
            <Card>
              <CardContent className="p-12 text-center">
                <Code className="w-12 h-12 mx-auto text-muted-foreground mb-4" />
                <h3 className="text-lg font-medium mb-2">Select a Contract</h3>
                <p className="text-muted-foreground">
                  Choose a contract from the Contracts tab to interact with it
                </p>
              </CardContent>
            </Card>
          )}
        </TabsContent>
      </Tabs>
    </div>
  )
}
