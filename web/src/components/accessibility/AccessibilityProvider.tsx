'use client'

import { createContext, useContext, useEffect, useState, ReactNode } from 'react'
import { motion, AnimatePresence } from 'framer-motion'
import { Button } from '@/components/ui/button'
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card'
import { Switch } from '@/components/ui/switch'
import { Label } from '@/components/ui/label'
import { 
  Eye, 
  EyeOff, 
  Volume2, 
  VolumeX, 
  Type, 
  Contrast, 
  MousePointer,
  Keyboard,
  Accessibility,
  X,
  Settings,
  Play,
  Pause,
  SkipForward,
  SkipBack
} from 'lucide-react'
import { cn } from '@/lib/utils'

interface AccessibilitySettings {
  screenReader: boolean
  highContrast: boolean
  largeText: boolean
  reducedMotion: boolean
  keyboardNavigation: boolean
  voiceAnnouncements: boolean
  focusIndicators: boolean
  colorBlindMode: 'none' | 'protanopia' | 'deuteranopia' | 'tritanopia'
  fontSize: number
  lineHeight: number
  letterSpacing: number
}

interface AccessibilityContextType {
  settings: AccessibilitySettings
  updateSetting: <K extends keyof AccessibilitySettings>(key: K, value: AccessibilitySettings[K]) => void
  announceMessage: (message: string) => void
  isAccessibilityPanelOpen: boolean
  setIsAccessibilityPanelOpen: (open: boolean) => void
}

const defaultSettings: AccessibilitySettings = {
  screenReader: false,
  highContrast: false,
  largeText: false,
  reducedMotion: false,
  keyboardNavigation: true,
  voiceAnnouncements: false,
  focusIndicators: true,
  colorBlindMode: 'none',
  fontSize: 16,
  lineHeight: 1.5,
  letterSpacing: 0
}

const AccessibilityContext = createContext<AccessibilityContextType | undefined>(undefined)

export function useAccessibility() {
  const context = useContext(AccessibilityContext)
  if (!context) {
    throw new Error('useAccessibility must be used within AccessibilityProvider')
  }
  return context
}

interface AccessibilityProviderProps {
  children: ReactNode
}

export function AccessibilityProvider({ children }: AccessibilityProviderProps) {
  const [settings, setSettings] = useState<AccessibilitySettings>(defaultSettings)
  const [isAccessibilityPanelOpen, setIsAccessibilityPanelOpen] = useState(false)
  const [announcements, setAnnouncements] = useState<string[]>([])

  const updateSetting = <K extends keyof AccessibilitySettings>(
    key: K,
    value: AccessibilitySettings[K]
  ) => {
    setSettings(prev => ({ ...prev, [key]: value }))
    
    // Apply settings immediately
    applyAccessibilitySettings({ ...settings, [key]: value })
    
    // Announce change
    if (settings.voiceAnnouncements) {
      announceMessage(`${key} ${value ? 'enabled' : 'disabled'}`)
    }
  }

  const announceMessage = (message: string) => {
    setAnnouncements(prev => [...prev, message])
    
    // Speak the message if voice announcements are enabled
    if (settings.voiceAnnouncements && 'speechSynthesis' in window) {
      const utterance = new SpeechSynthesisUtterance(message)
      utterance.rate = 0.8
      utterance.volume = 0.7
      speechSynthesis.speak(utterance)
    }
    
    // Remove announcement after 3 seconds
    setTimeout(() => {
      setAnnouncements(prev => prev.filter(a => a !== message))
    }, 3000)
  }

  const applyAccessibilitySettings = (newSettings: AccessibilitySettings) => {
    const root = document.documentElement

    // High contrast
    if (newSettings.highContrast) {
      root.classList.add('high-contrast')
    } else {
      root.classList.remove('high-contrast')
    }

    // Large text
    if (newSettings.largeText) {
      root.classList.add('large-text')
    } else {
      root.classList.remove('large-text')
    }

    // Reduced motion
    if (newSettings.reducedMotion) {
      root.classList.add('reduce-motion')
    } else {
      root.classList.remove('reduce-motion')
    }

    // Focus indicators
    if (newSettings.focusIndicators) {
      root.classList.add('enhanced-focus')
    } else {
      root.classList.remove('enhanced-focus')
    }

    // Color blind mode
    root.classList.remove('protanopia', 'deuteranopia', 'tritanopia')
    if (newSettings.colorBlindMode !== 'none') {
      root.classList.add(newSettings.colorBlindMode)
    }

    // Font settings
    root.style.setProperty('--accessibility-font-size', `${newSettings.fontSize}px`)
    root.style.setProperty('--accessibility-line-height', newSettings.lineHeight.toString())
    root.style.setProperty('--accessibility-letter-spacing', `${newSettings.letterSpacing}px`)
  }

  // Keyboard shortcuts
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      // Alt + A to open accessibility panel
      if (e.altKey && e.key === 'a') {
        e.preventDefault()
        setIsAccessibilityPanelOpen(true)
        announceMessage('Accessibility panel opened')
      }

      // Alt + H for high contrast toggle
      if (e.altKey && e.key === 'h') {
        e.preventDefault()
        updateSetting('highContrast', !settings.highContrast)
      }

      // Alt + T for large text toggle
      if (e.altKey && e.key === 't') {
        e.preventDefault()
        updateSetting('largeText', !settings.largeText)
      }

      // Alt + M for reduced motion toggle
      if (e.altKey && e.key === 'm') {
        e.preventDefault()
        updateSetting('reducedMotion', !settings.reducedMotion)
      }

      // Alt + V for voice announcements toggle
      if (e.altKey && e.key === 'v') {
        e.preventDefault()
        updateSetting('voiceAnnouncements', !settings.voiceAnnouncements)
      }
    }

    window.addEventListener('keydown', handleKeyDown)
    return () => window.removeEventListener('keydown', handleKeyDown)
  }, [settings])

  // Apply settings on mount
  useEffect(() => {
    applyAccessibilitySettings(settings)
  }, [])

  return (
    <AccessibilityContext.Provider value={{
      settings,
      updateSetting,
      announceMessage,
      isAccessibilityPanelOpen,
      setIsAccessibilityPanelOpen
    }}>
      {children}
      
      {/* Live Region for Announcements */}
      <div
        aria-live="polite"
        aria-atomic="true"
        className="sr-only"
      >
        {announcements.map((announcement, index) => (
          <div key={index}>{announcement}</div>
        ))}
      </div>

      {/* Accessibility Panel */}
      <AccessibilityPanel />

      {/* Skip Links */}
      <SkipLinks />
    </AccessibilityContext.Provider>
  )
}

function AccessibilityPanel() {
  const { settings, updateSetting, isAccessibilityPanelOpen, setIsAccessibilityPanelOpen } = useAccessibility()

  return (
    <AnimatePresence>
      {isAccessibilityPanelOpen && (
        <>
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            exit={{ opacity: 0 }}
            className="fixed inset-0 bg-black/50 z-50"
            onClick={() => setIsAccessibilityPanelOpen(false)}
          />
          
          <motion.div
            initial={{ x: '-100%' }}
            animate={{ x: 0 }}
            exit={{ x: '-100%' }}
            transition={{ type: 'spring', damping: 20, stiffness: 300 }}
            className="fixed left-0 top-0 h-full w-96 bg-background border-r border-border z-50 overflow-y-auto"
          >
            <Card className="h-full rounded-none border-0">
              <CardHeader className="border-b">
                <div className="flex items-center justify-between">
                  <div className="flex items-center space-x-2">
                    <Accessibility className="h-5 w-5" />
                    <CardTitle>Accessibility Settings</CardTitle>
                  </div>
                  <Button
                    variant="ghost"
                    size="icon"
                    onClick={() => setIsAccessibilityPanelOpen(false)}
                    aria-label="Close accessibility panel"
                  >
                    <X className="h-4 w-4" />
                  </Button>
                </div>
                <CardDescription>
                  Customize your experience for better accessibility
                </CardDescription>
              </CardHeader>

              <CardContent className="space-y-6 p-6">
                {/* Visual Settings */}
                <div className="space-y-4">
                  <h3 className="font-medium flex items-center space-x-2">
                    <Eye className="h-4 w-4" />
                    <span>Visual</span>
                  </h3>
                  
                  <div className="space-y-3">
                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="high-contrast">High Contrast</Label>
                        <p className="text-sm text-muted-foreground">Increase color contrast</p>
                      </div>
                      <Switch
                        id="high-contrast"
                        checked={settings.highContrast}
                        onCheckedChange={(checked) => updateSetting('highContrast', checked)}
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="large-text">Large Text</Label>
                        <p className="text-sm text-muted-foreground">Increase text size</p>
                      </div>
                      <Switch
                        id="large-text"
                        checked={settings.largeText}
                        onCheckedChange={(checked) => updateSetting('largeText', checked)}
                      />
                    </div>

                    <div className="flex items-center justify-between">
                      <div>
                        <Label htmlFor="focus-indicators">Enhanced Focus</Label>
                        <p className="text-sm text-muted-foreground">Stronger focus indicators</p>
                      </div>
                      <Switch
                        id="focus-indicators"
                        checked={settings.focusIndicators}
                        onCheckedChange={(checked) => updateSetting('focusIndicators', checked)}
                      />
                    </div>

                    <div>
                      <Label htmlFor="color-blind-mode">Color Blind Support</Label>
                      <select
                        id="color-blind-mode"
                        value={settings.colorBlindMode}
                        onChange={(e) => updateSetting('colorBlindMode', e.target.value as any)}
                        className="w-full mt-1 p-2 border border-border rounded-md bg-background"
                      >
                        <option value="none">None</option>
                        <option value="protanopia">Protanopia (Red-blind)</option>
                        <option value="deuteranopia">Deuteranopia (Green-blind)</option>
                        <option value="tritanopia">Tritanopia (Blue-blind)</option>
                      </select>
                    </div>
                  </div>
                </div>

                {/* Motion Settings */}
                <div className="space-y-4">
                  <h3 className="font-medium flex items-center space-x-2">
                    <MousePointer className="h-4 w-4" />
                    <span>Motion</span>
                  </h3>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="reduced-motion">Reduced Motion</Label>
                      <p className="text-sm text-muted-foreground">Minimize animations</p>
                    </div>
                    <Switch
                      id="reduced-motion"
                      checked={settings.reducedMotion}
                      onCheckedChange={(checked) => updateSetting('reducedMotion', checked)}
                    />
                  </div>
                </div>

                {/* Audio Settings */}
                <div className="space-y-4">
                  <h3 className="font-medium flex items-center space-x-2">
                    <Volume2 className="h-4 w-4" />
                    <span>Audio</span>
                  </h3>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="voice-announcements">Voice Announcements</Label>
                      <p className="text-sm text-muted-foreground">Speak important updates</p>
                    </div>
                    <Switch
                      id="voice-announcements"
                      checked={settings.voiceAnnouncements}
                      onCheckedChange={(checked) => updateSetting('voiceAnnouncements', checked)}
                    />
                  </div>
                </div>

                {/* Navigation Settings */}
                <div className="space-y-4">
                  <h3 className="font-medium flex items-center space-x-2">
                    <Keyboard className="h-4 w-4" />
                    <span>Navigation</span>
                  </h3>
                  
                  <div className="flex items-center justify-between">
                    <div>
                      <Label htmlFor="keyboard-navigation">Keyboard Navigation</Label>
                      <p className="text-sm text-muted-foreground">Enhanced keyboard support</p>
                    </div>
                    <Switch
                      id="keyboard-navigation"
                      checked={settings.keyboardNavigation}
                      onCheckedChange={(checked) => updateSetting('keyboardNavigation', checked)}
                    />
                  </div>
                </div>

                {/* Keyboard Shortcuts Help */}
                <div className="space-y-2 p-4 bg-muted/50 rounded-lg">
                  <h4 className="font-medium text-sm">Keyboard Shortcuts</h4>
                  <div className="text-xs text-muted-foreground space-y-1">
                    <div>Alt + A: Open accessibility panel</div>
                    <div>Alt + H: Toggle high contrast</div>
                    <div>Alt + T: Toggle large text</div>
                    <div>Alt + M: Toggle reduced motion</div>
                    <div>Alt + V: Toggle voice announcements</div>
                  </div>
                </div>
              </CardContent>
            </Card>
          </motion.div>
        </>
      )}
    </AnimatePresence>
  )
}

function SkipLinks() {
  return (
    <div className="sr-only focus-within:not-sr-only">
      <a
        href="#main-content"
        className="fixed top-4 left-4 z-50 bg-primary text-primary-foreground px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-ring"
      >
        Skip to main content
      </a>
      <a
        href="#navigation"
        className="fixed top-4 left-32 z-50 bg-primary text-primary-foreground px-4 py-2 rounded-md focus:outline-none focus:ring-2 focus:ring-ring"
      >
        Skip to navigation
      </a>
    </div>
  )
}
