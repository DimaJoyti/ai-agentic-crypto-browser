'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { usePathname, useRouter } from 'next/navigation'
import Link from 'next/link'
import { Button } from '@/components/ui/button'
import { Badge } from '@/components/ui/badge'
import {
  Menu,
  X,
  Home,
  TrendingUp,
  BarChart3,
  Wallet,
  Settings,
  User,
  Bell,
  Search,
  ChevronRight,
  Zap,
  Shield,
  Activity,
  Globe,
  Layers,
  Target,
  Briefcase,
  PieChart,
  Coins
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface NavigationItem {
  id: string
  label: string
  href: string
  icon: React.ComponentType<any>
  badge?: string
  children?: NavigationItem[]
  description?: string
  isNew?: boolean
  isPopular?: boolean
}

const navigationItems: NavigationItem[] = [
  {
    id: 'home',
    label: 'Home',
    href: '/',
    icon: Home,
    description: 'Dashboard overview'
  },
  {
    id: 'trading',
    label: 'Trading',
    href: '/trading',
    icon: TrendingUp,
    badge: 'Live',
    isPopular: true,
    description: 'Advanced trading platform',
    children: [
      { id: 'hft', label: 'HFT Trading', href: '/trading/hft', icon: Zap, isNew: true },
      { id: 'algorithms', label: 'Algorithms', href: '/trading/algorithms', icon: Target },
      { id: 'portfolio', label: 'Portfolio', href: '/trading/portfolio', icon: Briefcase },
      { id: 'risk', label: 'Risk Management', href: '/trading/risk', icon: Shield }
    ]
  },
  {
    id: 'analytics',
    label: 'Analytics',
    href: '/analytics',
    icon: BarChart3,
    description: 'Performance analytics',
    children: [
      { id: 'performance', label: 'Performance', href: '/analytics/performance', icon: Activity },
      { id: 'market', label: 'Market Analysis', href: '/analytics/market', icon: PieChart },
      { id: 'backtesting', label: 'Backtesting', href: '/analytics/backtesting', icon: BarChart3 }
    ]
  },
  {
    id: 'web3',
    label: 'Web3',
    href: '/web3',
    icon: Globe,
    badge: 'Beta',
    description: 'DeFi and Web3 tools',
    children: [
      { id: 'defi', label: 'DeFi Dashboard', href: '/web3/defi', icon: Layers },
      { id: 'nft', label: 'NFT Portfolio', href: '/web3/nft', icon: Coins },
      { id: 'wallet', label: 'Wallet Manager', href: '/web3/wallet', icon: Wallet }
    ]
  },
  {
    id: 'dashboard',
    label: 'Dashboard',
    href: '/dashboard',
    icon: Activity,
    description: 'Main dashboard'
  }
]

interface EnhancedNavigationProps {
  className?: string
}

export function EnhancedNavigation({ className }: EnhancedNavigationProps) {
  const [isOpen, setIsOpen] = useState(false)
  const [expandedItems, setExpandedItems] = useState<Set<string>>(new Set())
  const [searchQuery, setSearchQuery] = useState('')
  const [notifications, _setNotifications] = useState(3)
  const [_isCollapsed, setIsCollapsed] = useState(false)
  const pathname = usePathname()
  const router = useRouter()

  // Auto-collapse on mobile
  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth < 1024) {
        setIsCollapsed(true)
      } else {
        setIsCollapsed(false)
      }
    }

    handleResize()
    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener('resize', handleResize)
  }, [])

  // Auto-expand current section
  useEffect(() => {
    navigationItems.forEach(item => {
      if (item.children) {
        const hasActiveChild = item.children.some(child => pathname.startsWith(child.href))
        if (hasActiveChild) {
          setExpandedItems(prev => new Set([...Array.from(prev), item.id]))
        }
      }
    })
  }, [pathname])

  const toggleExpanded = (itemId: string) => {
    setExpandedItems(prev => {
      const newSet = new Set(prev)
      if (newSet.has(itemId)) {
        newSet.delete(itemId)
      } else {
        newSet.add(itemId)
      }
      return newSet
    })
  }

  const isActive = (href: string) => {
    if (href === '/') {
      return pathname === '/'
    }
    return pathname.startsWith(href)
  }

  const filteredItems = navigationItems.filter(item =>
    item.label.toLowerCase().includes(searchQuery.toLowerCase()) ||
    item.children?.some(child => 
      child.label.toLowerCase().includes(searchQuery.toLowerCase())
    )
  )

  return (
    <>
      {/* Mobile Menu Button */}
      <Button
        variant="ghost"
        size="icon"
        className="fixed top-4 left-4 z-50 lg:hidden bg-background/80 backdrop-blur-sm"
        onClick={() => setIsOpen(!isOpen)}
      >
        {isOpen ? <X className="h-5 w-5" /> : <Menu className="h-5 w-5" />}
      </Button>

      {/* Mobile Overlay */}
      <AnimatePresence>
        {isOpen && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-40 lg:hidden"
            onClick={() => setIsOpen(false)}
          />
        )}
      </AnimatePresence>

      {/* Navigation Sidebar */}
      <motion.nav
        initial={false}
        animate={{
          x: isOpen ? 0 : '-100%',
          opacity: isOpen ? 1 : 0
        }}
        className={cn(
          "fixed left-0 top-0 h-full w-80 bg-background/95 backdrop-blur-xl border-r border-border z-50",
          "lg:relative lg:translate-x-0 lg:opacity-100 lg:w-64",
          "flex flex-col shadow-xl lg:shadow-none",
          className
        )}
      >
        {/* Header */}
        <div className="p-6 border-b border-border">
          <div className="flex items-center justify-between">
            <div className="flex items-center space-x-3">
              <div className="w-8 h-8 bg-gradient-to-br from-blue-500 to-purple-600 rounded-lg flex items-center justify-center">
                <Zap className="h-4 w-4 text-white" />
              </div>
              <div>
                <h2 className="font-semibold text-lg">AI Browser</h2>
                <p className="text-xs text-muted-foreground">v2.0.0</p>
              </div>
            </div>
            <div className="relative">
              <Button variant="ghost" size="icon" className="relative">
                <Bell className="h-4 w-4" />
                {notifications > 0 && (
                  <Badge 
                    variant="destructive" 
                    className="absolute -top-1 -right-1 h-5 w-5 p-0 flex items-center justify-center text-xs"
                  >
                    {notifications}
                  </Badge>
                )}
              </Button>
            </div>
          </div>

          {/* Search */}
          <div className="mt-4 relative">
            <Search className="absolute left-3 top-1/2 transform -translate-y-1/2 h-4 w-4 text-muted-foreground" />
            <input
              type="text"
              placeholder="Search navigation..."
              value={searchQuery}
              onChange={(e) => setSearchQuery(e.target.value)}
              className="w-full pl-10 pr-4 py-2 bg-muted/50 border border-border rounded-lg text-sm focus:outline-none focus:ring-2 focus:ring-primary/20"
            />
          </div>
        </div>

        {/* Navigation Items */}
        <div className="flex-1 overflow-y-auto p-4 space-y-2">
          {filteredItems.map((item) => (
            <div key={item.id}>
              <div
                className={cn(
                  "flex items-center justify-between p-3 rounded-lg transition-all duration-200 cursor-pointer group",
                  isActive(item.href)
                    ? "bg-primary/10 text-primary border border-primary/20"
                    : "hover:bg-muted/50 text-foreground"
                )}
                onClick={() => {
                  if (item.children) {
                    toggleExpanded(item.id)
                  } else {
                    router.push(item.href)
                    setIsOpen(false)
                  }
                }}
              >
                <div className="flex items-center space-x-3 flex-1">
                  <item.icon className={cn(
                    "h-5 w-5 transition-colors",
                    isActive(item.href) ? "text-primary" : "text-muted-foreground group-hover:text-foreground"
                  )} />
                  <div className="flex-1">
                    <div className="flex items-center space-x-2">
                      <span className="font-medium">{item.label}</span>
                      {item.badge && (
                        <Badge variant={item.badge === 'Live' ? 'destructive' : 'secondary'} className="text-xs">
                          {item.badge}
                        </Badge>
                      )}
                      {item.isNew && (
                        <Badge variant="default" className="text-xs bg-green-500">
                          New
                        </Badge>
                      )}
                      {item.isPopular && (
                        <Badge variant="outline" className="text-xs">
                          Popular
                        </Badge>
                      )}
                    </div>
                    {item.description && (
                      <p className="text-xs text-muted-foreground mt-1">{item.description}</p>
                    )}
                  </div>
                </div>
                {item.children && (
                  <motion.div
                    animate={{ rotate: expandedItems.has(item.id) ? 90 : 0 }}
                    transition={{ duration: 0.2 }}
                  >
                    <ChevronRight className="h-4 w-4 text-muted-foreground" />
                  </motion.div>
                )}
              </div>

              {/* Submenu */}
              <AnimatePresence>
                {item.children && expandedItems.has(item.id) && (
                  <motion.div
                    initial={{ height: 0, opacity: 0 }}
                    animate={{ height: 'auto', opacity: 1 }}
                    exit={{ height: 0, opacity: 0 }}
                    transition={{ duration: 0.2 }}
                    className="overflow-hidden"
                  >
                    <div className="ml-8 mt-2 space-y-1">
                      {item.children.map((child) => (
                        <Link
                          key={child.id}
                          href={child.href}
                          onClick={() => setIsOpen(false)}
                          className={cn(
                            "flex items-center space-x-3 p-2 rounded-md transition-all duration-200 group",
                            isActive(child.href)
                              ? "bg-primary/10 text-primary"
                              : "hover:bg-muted/30 text-muted-foreground hover:text-foreground"
                          )}
                        >
                          <child.icon className="h-4 w-4" />
                          <span className="text-sm">{child.label}</span>
                          {child.isNew && (
                            <Badge variant="default" className="text-xs bg-green-500 ml-auto">
                              New
                            </Badge>
                          )}
                        </Link>
                      ))}
                    </div>
                  </motion.div>
                )}
              </AnimatePresence>
            </div>
          ))}
        </div>

        {/* Footer */}
        <div className="p-4 border-t border-border">
          <div className="flex items-center space-x-3 p-3 rounded-lg bg-muted/30">
            <div className="w-8 h-8 bg-gradient-to-br from-green-400 to-blue-500 rounded-full flex items-center justify-center">
              <User className="h-4 w-4 text-white" />
            </div>
            <div className="flex-1">
              <p className="font-medium text-sm">John Doe</p>
              <p className="text-xs text-muted-foreground">Premium User</p>
            </div>
            <Button variant="ghost" size="icon">
              <Settings className="h-4 w-4" />
            </Button>
          </div>
        </div>
      </motion.nav>
    </>
  )
}
