'use client'

import { useState, useEffect } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { useTheme } from 'next-themes'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Badge } from '@/components/ui/badge'
import { Tabs, TabsContent, TabsList, TabsTrigger } from '@/components/ui/tabs'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { 
  Palette, 
  Monitor, 
  Sun, 
  Moon, 
  Smartphone,
  Eye,
  Zap,
  Settings,
  Paintbrush,
  Layout,
  Type,
  Contrast,
  Volume2,
  VolumeX,
  Vibrate,
  Bell,
  X,
  Check,
  RotateCcw
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface ThemeCustomizerProps {
  isOpen: boolean
  onClose: () => void
}

interface CustomTheme {
  id: string
  name: string
  description: string
  colors: {
    primary: string
    secondary: string
    accent: string
    background: string
    foreground: string
  }
  preview: string
}

const predefinedThemes: CustomTheme[] = [
  {
    id: 'default',
    name: 'Default',
    description: 'Clean and professional',
    colors: {
      primary: 'hsl(222.2 84% 4.9%)',
      secondary: 'hsl(210 40% 96%)',
      accent: 'hsl(210 40% 96%)',
      background: 'hsl(0 0% 100%)',
      foreground: 'hsl(222.2 84% 4.9%)'
    },
    preview: 'bg-gradient-to-br from-slate-50 to-slate-100'
  },
  {
    id: 'ocean',
    name: 'Ocean Blue',
    description: 'Calm and focused',
    colors: {
      primary: 'hsl(210 100% 50%)',
      secondary: 'hsl(210 100% 95%)',
      accent: 'hsl(195 100% 50%)',
      background: 'hsl(210 100% 98%)',
      foreground: 'hsl(210 100% 15%)'
    },
    preview: 'bg-gradient-to-br from-blue-50 to-cyan-100'
  },
  {
    id: 'forest',
    name: 'Forest Green',
    description: 'Natural and calming',
    colors: {
      primary: 'hsl(120 60% 30%)',
      secondary: 'hsl(120 60% 95%)',
      accent: 'hsl(90 60% 50%)',
      background: 'hsl(120 60% 98%)',
      foreground: 'hsl(120 60% 15%)'
    },
    preview: 'bg-gradient-to-br from-green-50 to-emerald-100'
  },
  {
    id: 'sunset',
    name: 'Sunset Orange',
    description: 'Warm and energetic',
    colors: {
      primary: 'hsl(25 95% 53%)',
      secondary: 'hsl(25 95% 95%)',
      accent: 'hsl(45 95% 60%)',
      background: 'hsl(25 95% 98%)',
      foreground: 'hsl(25 95% 15%)'
    },
    preview: 'bg-gradient-to-br from-orange-50 to-yellow-100'
  },
  {
    id: 'purple',
    name: 'Royal Purple',
    description: 'Elegant and sophisticated',
    colors: {
      primary: 'hsl(270 95% 40%)',
      secondary: 'hsl(270 95% 95%)',
      accent: 'hsl(300 95% 50%)',
      background: 'hsl(270 95% 98%)',
      foreground: 'hsl(270 95% 15%)'
    },
    preview: 'bg-gradient-to-br from-purple-50 to-pink-100'
  },
  {
    id: 'dark',
    name: 'Dark Mode',
    description: 'Easy on the eyes',
    colors: {
      primary: 'hsl(210 40% 98%)',
      secondary: 'hsl(217.2 32.6% 17.5%)',
      accent: 'hsl(217.2 32.6% 17.5%)',
      background: 'hsl(222.2 84% 4.9%)',
      foreground: 'hsl(210 40% 98%)'
    },
    preview: 'bg-gradient-to-br from-slate-800 to-slate-900'
  }
]

export function ThemeCustomizer({ isOpen, onClose }: ThemeCustomizerProps) {
  const { theme, setTheme } = useTheme()
  const [selectedTheme, setSelectedTheme] = useState('default')
  const [customSettings, setCustomSettings] = useState({
    fontSize: 'medium',
    animations: true,
    sounds: true,
    vibrations: true,
    notifications: true,
    highContrast: false,
    reducedMotion: false,
    compactMode: false,
    autoTheme: true
  })

  const applyTheme = (themeData: CustomTheme) => {
    const root = document.documentElement
    Object.entries(themeData.colors).forEach(([key, value]) => {
      root.style.setProperty(`--${key}`, value)
    })
    setSelectedTheme(themeData.id)
  }

  const resetToDefault = () => {
    const defaultTheme = predefinedThemes.find(t => t.id === 'default')
    if (defaultTheme) {
      applyTheme(defaultTheme)
    }
    setCustomSettings({
      fontSize: 'medium',
      animations: true,
      sounds: true,
      vibrations: true,
      notifications: true,
      highContrast: false,
      reducedMotion: false,
      compactMode: false,
      autoTheme: true
    })
  }

  return (
    <AnimatePresence>
      {isOpen && (
        <>
          {/* Overlay */}
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-50"
            onClick={onClose}
          />

          {/* Customizer Panel */}
          <motion.div
            initial={{ x: '100%' }}
            animate={{ x: 0 }}
            exit={{ x: '100%' }}
            transition={{ type: 'spring', damping: 20, stiffness: 300 }}
            className="fixed right-0 top-0 h-full w-96 bg-background border-l border-border z-50 overflow-y-auto"
          >
            <div className="p-6">
              {/* Header */}
              <div className="flex items-center justify-between mb-6">
                <div className="flex items-center space-x-3">
                  <div className="w-8 h-8 bg-gradient-to-br from-purple-500 to-pink-500 rounded-lg flex items-center justify-center">
                    <Palette className="h-4 w-4 text-white" />
                  </div>
                  <div>
                    <h2 className="font-semibold text-lg">Theme Customizer</h2>
                    <p className="text-sm text-muted-foreground">Personalize your experience</p>
                  </div>
                </div>
                <Button variant="ghost" size="icon" onClick={onClose}>
                  <X className="h-4 w-4" />
                </Button>
              </div>

              <Tabs defaultValue="themes" className="w-full">
                <TabsList className="grid w-full grid-cols-3">
                  <TabsTrigger value="themes">Themes</TabsTrigger>
                  <TabsTrigger value="layout">Layout</TabsTrigger>
                  <TabsTrigger value="accessibility">A11y</TabsTrigger>
                </TabsList>

                {/* Themes Tab */}
                <TabsContent value="themes" className="space-y-4">
                  <div>
                    <h3 className="font-medium mb-3">Color Themes</h3>
                    <div className="grid grid-cols-2 gap-3">
                      {predefinedThemes.map((themeData) => (
                        <motion.div
                          key={themeData.id}
                          whileHover={{ scale: 1.02 }}
                          whileTap={{ scale: 0.98 }}
                          className={cn(
                            "relative p-3 rounded-lg border-2 cursor-pointer transition-all",
                            selectedTheme === themeData.id
                              ? "border-primary bg-primary/5"
                              : "border-border hover:border-primary/50"
                          )}
                          onClick={() => applyTheme(themeData)}
                        >
                          <div className={cn("w-full h-16 rounded-md mb-2", themeData.preview)} />
                          <h4 className="font-medium text-sm">{themeData.name}</h4>
                          <p className="text-xs text-muted-foreground">{themeData.description}</p>
                          {selectedTheme === themeData.id && (
                            <div className="absolute top-2 right-2 w-5 h-5 bg-primary rounded-full flex items-center justify-center">
                              <Check className="h-3 w-3 text-primary-foreground" />
                            </div>
                          )}
                        </motion.div>
                      ))}
                    </div>
                  </div>

                  <div>
                    <h3 className="font-medium mb-3">System Theme</h3>
                    <div className="flex space-x-2">
                      <Button
                        variant={theme === 'light' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setTheme('light')}
                        className="flex-1"
                      >
                        <Sun className="h-4 w-4 mr-2" />
                        Light
                      </Button>
                      <Button
                        variant={theme === 'dark' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setTheme('dark')}
                        className="flex-1"
                      >
                        <Moon className="h-4 w-4 mr-2" />
                        Dark
                      </Button>
                      <Button
                        variant={theme === 'system' ? 'default' : 'outline'}
                        size="sm"
                        onClick={() => setTheme('system')}
                        className="flex-1"
                      >
                        <Monitor className="h-4 w-4 mr-2" />
                        Auto
                      </Button>
                    </div>
                  </div>
                </TabsContent>

                {/* Layout Tab */}
                <TabsContent value="layout" className="space-y-4">
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="compact-mode">Compact Mode</Label>
                        <p className="text-sm text-muted-foreground">Reduce spacing and padding</p>
                      </div>
                      <Switch
                        id="compact-mode"
                        checked={customSettings.compactMode}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, compactMode: checked }))
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="animations">Animations</Label>
                        <p className="text-sm text-muted-foreground">Enable smooth transitions</p>
                      </div>
                      <Switch
                        id="animations"
                        checked={customSettings.animations}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, animations: checked }))
                        }
                      />
                    </div>

                    <div>
                      <Label>Font Size</Label>
                      <div className="flex space-x-2 mt-2">
                        {['small', 'medium', 'large'].map((size) => (
                          <Button
                            key={size}
                            variant={customSettings.fontSize === size ? 'default' : 'outline'}
                            size="sm"
                            onClick={() => setCustomSettings(prev => ({ ...prev, fontSize: size }))}
                            className="flex-1 capitalize"
                          >
                            <Type className="h-4 w-4 mr-2" />
                            {size}
                          </Button>
                        ))}
                      </div>
                    </div>
                  </div>
                </TabsContent>

                {/* Accessibility Tab */}
                <TabsContent value="accessibility" className="space-y-4">
                  <div className="space-y-4">
                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="high-contrast">High Contrast</Label>
                        <p className="text-sm text-muted-foreground">Increase color contrast</p>
                      </div>
                      <Switch
                        id="high-contrast"
                        checked={customSettings.highContrast}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, highContrast: checked }))
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="reduced-motion">Reduced Motion</Label>
                        <p className="text-sm text-muted-foreground">Minimize animations</p>
                      </div>
                      <Switch
                        id="reduced-motion"
                        checked={customSettings.reducedMotion}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, reducedMotion: checked }))
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="sounds">Sound Effects</Label>
                        <p className="text-sm text-muted-foreground">Audio feedback</p>
                      </div>
                      <Switch
                        id="sounds"
                        checked={customSettings.sounds}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, sounds: checked }))
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="vibrations">Haptic Feedback</Label>
                        <p className="text-sm text-muted-foreground">Vibration on mobile</p>
                      </div>
                      <Switch
                        id="vibrations"
                        checked={customSettings.vibrations}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, vibrations: checked }))
                        }
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="notifications">Notifications</Label>
                        <p className="text-sm text-muted-foreground">System notifications</p>
                      </div>
                      <Switch
                        id="notifications"
                        checked={customSettings.notifications}
                        onCheckedChange={(checked) =>
                          setCustomSettings(prev => ({ ...prev, notifications: checked }))
                        }
                      />
                    </div>
                  </div>
                </TabsContent>
              </Tabs>

              {/* Actions */}
              <div className="flex space-x-2 mt-6">
                <Button onClick={resetToDefault} variant="outline" className="flex-1">
                  <RotateCcw className="h-4 w-4 mr-2" />
                  Reset
                </Button>
                <Button onClick={onClose} className="flex-1">
                  <Check className="h-4 w-4 mr-2" />
                  Apply
                </Button>
              </div>
            </div>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}
