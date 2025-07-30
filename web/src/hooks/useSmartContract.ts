import { useState, useEffect, useCallback } from 'react'
import { useAccount, useChainId } from 'wagmi'
import { type Address, type Abi } from 'viem'
import { 
  smartContractIntegration,
  type SmartContract,
  type ContractCall,
  type ContractDeployment,
  type ContractInteractionConfig
} from '@/lib/smart-contract-integration'
import { toast } from 'sonner'

export interface SmartContractState {
  contracts: SmartContract[]
  calls: ContractCall[]
  deployments: ContractDeployment[]
  isLoading: boolean
  isExecuting: boolean
  isDeploying: boolean
  config: ContractInteractionConfig
  error: string | null
  lastUpdate: number | null
}

export interface UseSmartContractOptions {
  autoLoad?: boolean
  enableNotifications?: boolean
  enableSimulation?: boolean
  enableGasOptimization?: boolean
}

export interface UseSmartContractReturn {
  // State
  state: SmartContractState
  
  // Contract Management
  registerContract: (contract: Omit<SmartContract, 'functions' | 'events' | 'errors' | 'createdAt' | 'updatedAt'>) => SmartContract
  getContract: (address: Address, chainId?: number) => SmartContract | null
  searchContracts: (query: string) => SmartContract[]
  
  // Contract Calls
  prepareCall: (address: Address, functionName: string, args: any[], options?: Partial<ContractCall>) => Promise<ContractCall>
  executeCall: (call: ContractCall) => Promise<ContractCall>
  getCall: (id: string) => ContractCall | null
  
  // Contract Deployment
  deployContract: (deployment: Omit<ContractDeployment, 'id' | 'timestamp' | 'status'>) => Promise<ContractDeployment>
  getDeployment: (id: string) => ContractDeployment | null
  
  // Configuration
  updateConfig: (config: Partial<ContractInteractionConfig>) => void
  
  // Utilities
  refresh: () => void
  clearError: () => void
}

export const useSmartContract = (
  options: UseSmartContractOptions = {}
): UseSmartContractReturn => {
  const {
    autoLoad = true,
    enableNotifications = true,
    enableSimulation = true,
    enableGasOptimization = true
  } = options

  const { address } = useAccount()
  const chainId = useChainId()

  const [state, setState] = useState<SmartContractState>({
    contracts: [],
    calls: [],
    deployments: [],
    isLoading: false,
    isExecuting: false,
    isDeploying: false,
    config: smartContractIntegration.getConfig(),
    error: null,
    lastUpdate: null
  })

  // Update state from smart contract integration
  const updateState = useCallback(() => {
    try {
      const contracts = smartContractIntegration.getAllContracts()
      const config = smartContractIntegration.getConfig()

      setState(prev => ({
        ...prev,
        contracts,
        config,
        error: null,
        lastUpdate: Date.now()
      }))
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        error: errorMessage
      }))
    }
  }, [])

  // Handle contract events
  const handleContractEvent = useCallback((event: any) => {
    if (enableNotifications) {
      switch (event.type) {
        case 'call_success':
          toast.success('Contract Call Successful', {
            description: `Function ${event.call.functionName} executed successfully`
          })
          break
        case 'call_failed':
          toast.error('Contract Call Failed', {
            description: `Function ${event.call.functionName} failed: ${event.error.message}`
          })
          break
        case 'deployment_success':
          toast.success('Contract Deployed', {
            description: `Contract ${event.deployment.name} deployed at ${event.deployment.address?.slice(0, 10)}...`
          })
          break
        case 'deployment_failed':
          toast.error('Deployment Failed', {
            description: `Contract ${event.deployment.name} deployment failed: ${event.error.message}`
          })
          break
      }
    }

    // Update state after event
    updateState()
  }, [enableNotifications, updateState])

  // Initialize and setup event listeners
  useEffect(() => {
    // Add event listener
    const unsubscribe = smartContractIntegration.addEventListener(handleContractEvent)

    // Update configuration
    smartContractIntegration.updateConfig({
      enableSimulation,
      enableGasOptimization
    })

    // Initial state update
    if (autoLoad) {
      updateState()
    }

    return () => {
      unsubscribe()
    }
  }, [autoLoad, enableSimulation, enableGasOptimization, handleContractEvent, updateState])

  // Register contract
  const registerContract = useCallback((
    contract: Omit<SmartContract, 'functions' | 'events' | 'errors' | 'createdAt' | 'updatedAt'>
  ): SmartContract => {
    try {
      const registeredContract = smartContractIntegration.registerContract(contract)
      updateState()

      if (enableNotifications) {
        toast.success('Contract Registered', {
          description: `Contract ${contract.name} registered successfully`
        })
      }

      return registeredContract
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to register contract', { description: errorMessage })
      }
      throw error
    }
  }, [enableNotifications, updateState])

  // Get contract
  const getContract = useCallback((address: Address, contractChainId?: number): SmartContract | null => {
    const targetChainId = contractChainId || chainId || 1
    return smartContractIntegration.getContract(address, targetChainId)
  }, [chainId])

  // Search contracts
  const searchContracts = useCallback((query: string): SmartContract[] => {
    return smartContractIntegration.searchContracts(query)
  }, [])

  // Prepare contract call
  const prepareCall = useCallback(async (
    address: Address,
    functionName: string,
    args: any[],
    options: Partial<ContractCall> = {}
  ): Promise<ContractCall> => {
    setState(prev => ({ ...prev, isLoading: true, error: null }))

    try {
      const targetChainId = chainId || 1
      const call = await smartContractIntegration.prepareCall(
        address,
        targetChainId,
        functionName,
        args,
        options
      )

      setState(prev => ({ ...prev, isLoading: false }))

      if (enableNotifications) {
        toast.info('Call Prepared', {
          description: `Function ${functionName} ready for execution`
        })
      }

      return call
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isLoading: false,
        error: errorMessage
      }))

      if (enableNotifications) {
        toast.error('Failed to prepare call', { description: errorMessage })
      }
      throw error
    }
  }, [chainId, enableNotifications])

  // Execute contract call
  const executeCall = useCallback(async (call: ContractCall): Promise<ContractCall> => {
    setState(prev => ({ ...prev, isExecuting: true, error: null }))

    try {
      const result = await smartContractIntegration.executeCall(call)
      
      setState(prev => ({
        ...prev,
        isExecuting: false,
        calls: [...prev.calls.filter(c => c.id !== call.id), result]
      }))

      return result
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isExecuting: false,
        error: errorMessage
      }))
      throw error
    }
  }, [])

  // Get contract call
  const getCall = useCallback((id: string): ContractCall | null => {
    return smartContractIntegration.getCall(id)
  }, [])

  // Deploy contract
  const deployContract = useCallback(async (
    deployment: Omit<ContractDeployment, 'id' | 'timestamp' | 'status'>
  ): Promise<ContractDeployment> => {
    setState(prev => ({ ...prev, isDeploying: true, error: null }))

    try {
      const result = await smartContractIntegration.deployContract({
        ...deployment,
        deployer: address || '0x0000000000000000000000000000000000000000',
        chainId: deployment.chainId || chainId || 1
      })

      setState(prev => ({
        ...prev,
        isDeploying: false,
        deployments: [...prev.deployments, result]
      }))

      return result
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({
        ...prev,
        isDeploying: false,
        error: errorMessage
      }))
      throw error
    }
  }, [address, chainId])

  // Get deployment
  const getDeployment = useCallback((id: string): ContractDeployment | null => {
    return smartContractIntegration.getDeployment(id)
  }, [])

  // Update configuration
  const updateConfig = useCallback((config: Partial<ContractInteractionConfig>) => {
    try {
      smartContractIntegration.updateConfig(config)
      setState(prev => ({ ...prev, config: smartContractIntegration.getConfig() }))

      if (enableNotifications) {
        toast.success('Configuration Updated', {
          description: 'Smart contract settings have been updated'
        })
      }
    } catch (error) {
      const errorMessage = (error as Error).message
      setState(prev => ({ ...prev, error: errorMessage }))

      if (enableNotifications) {
        toast.error('Failed to update configuration', { description: errorMessage })
      }
    }
  }, [enableNotifications])

  // Refresh state
  const refresh = useCallback(() => {
    updateState()
  }, [updateState])

  // Clear error
  const clearError = useCallback(() => {
    setState(prev => ({ ...prev, error: null }))
  }, [])

  return {
    state,
    registerContract,
    getContract,
    searchContracts,
    prepareCall,
    executeCall,
    getCall,
    deployContract,
    getDeployment,
    updateConfig,
    refresh,
    clearError
  }
}

// Simplified hook for contract interactions
export const useContractCall = (address: Address, abi: Abi, functionName: string) => {
  const { prepareCall, executeCall, state } = useSmartContract()
  const chainId = useChainId()

  const call = useCallback(async (args: any[], options?: Partial<ContractCall>) => {
    const prepared = await prepareCall(address, functionName, args, options)
    return executeCall(prepared)
  }, [address, functionName, prepareCall, executeCall])

  return {
    call,
    isExecuting: state.isExecuting,
    error: state.error
  }
}

// Hook for contract deployment
export const useContractDeployment = () => {
  const { deployContract, state } = useSmartContract()

  const deploy = useCallback(async (
    name: string,
    bytecode: string,
    abi: Abi,
    constructorArgs: any[] = [],
    options: Partial<ContractDeployment> = {}
  ) => {
    return deployContract({
      name,
      bytecode,
      abi,
      constructorArgs,
      deployer: '0x0000000000000000000000000000000000000000', // Will be set by hook
      chainId: 1, // Will be set by hook
      ...options
    })
  }, [deployContract])

  return {
    deploy,
    isDeploying: state.isDeploying,
    deployments: state.deployments,
    error: state.error
  }
}

// Hook for contract registry
export const useContractRegistry = () => {
  const { state, registerContract, searchContracts, getContract } = useSmartContract()

  return {
    contracts: state.contracts,
    registerContract,
    searchContracts,
    getContract,
    isLoading: state.isLoading
  }
}
