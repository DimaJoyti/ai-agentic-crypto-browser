'use client'

import * as React from "react"
import { Slot } from "@radix-ui/react-slot"
import { cva, type VariantProps } from "class-variance-authority"
import { motion, HTMLMotionProps } from "framer-motion"
import { cn } from "@/lib/utils"
import { Loader2 } from "lucide-react"

const enhancedButtonVariants = cva(
  "inline-flex items-center justify-center whitespace-nowrap rounded-md text-sm font-medium ring-offset-background transition-all duration-300 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring focus-visible:ring-offset-2 disabled:pointer-events-none disabled:opacity-50 relative overflow-hidden group",
  {
    variants: {
      variant: {
        default: "bg-gradient-to-r from-blue-600 to-purple-600 text-white hover:from-blue-700 hover:to-purple-700 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95 before:absolute before:inset-0 before:bg-gradient-to-r before:from-white/20 before:to-transparent before:opacity-0 hover:before:opacity-100 before:transition-opacity before:duration-300",
        destructive: "bg-gradient-to-r from-red-600 to-red-700 text-white hover:from-red-700 hover:to-red-800 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95",
        outline: "border-2 border-input bg-background hover:bg-accent hover:text-accent-foreground hover:border-primary/50 hover:shadow-md hover:scale-105 active:scale-95 backdrop-blur-sm",
        secondary: "bg-gradient-to-r from-slate-200 to-slate-300 dark:from-slate-700 dark:to-slate-800 text-secondary-foreground hover:from-slate-300 hover:to-slate-400 dark:hover:from-slate-600 dark:hover:to-slate-700 shadow-md hover:shadow-lg hover:scale-105 active:scale-95",
        ghost: "hover:bg-accent hover:text-accent-foreground hover:shadow-md hover:scale-105 active:scale-95 backdrop-blur-sm",
        link: "text-primary underline-offset-4 hover:underline hover:text-primary/80 transition-colors duration-300",
        gradient: "bg-gradient-to-r from-emerald-500 to-teal-600 text-white hover:from-emerald-600 hover:to-teal-700 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95",
        glass: "bg-white/10 dark:bg-black/10 backdrop-blur-md border border-white/20 dark:border-white/10 text-foreground hover:bg-white/20 dark:hover:bg-black/20 hover:border-white/30 dark:hover:border-white/20 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95",
        premium: "bg-gradient-to-r from-yellow-400 via-yellow-500 to-yellow-600 text-black hover:from-yellow-500 hover:via-yellow-600 hover:to-yellow-700 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95 font-semibold",
        neon: "bg-gradient-to-r from-cyan-400 to-blue-500 text-white hover:from-cyan-500 hover:to-blue-600 shadow-lg hover:shadow-cyan-500/50 hover:shadow-xl hover:scale-105 active:scale-95 animate-pulse-glow",
        success: "bg-gradient-to-r from-green-500 to-emerald-600 text-white hover:from-green-600 hover:to-emerald-700 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95",
        warning: "bg-gradient-to-r from-orange-500 to-red-500 text-white hover:from-orange-600 hover:to-red-600 shadow-lg hover:shadow-xl hover:scale-105 active:scale-95",
      },
      size: {
        default: "h-10 px-4 py-2",
        sm: "h-9 rounded-md px-3 text-xs",
        lg: "h-12 rounded-lg px-8 text-base font-semibold",
        xl: "h-14 rounded-xl px-10 text-lg font-bold",
        icon: "h-10 w-10",
        "icon-sm": "h-8 w-8",
        "icon-lg": "h-12 w-12",
      },
      glow: {
        none: "",
        blue: "hover:shadow-blue-500/50",
        purple: "hover:shadow-purple-500/50",
        emerald: "hover:shadow-emerald-500/50",
        yellow: "hover:shadow-yellow-500/50",
        cyan: "hover:shadow-cyan-500/50",
        red: "hover:shadow-red-500/50",
      },
    },
    defaultVariants: {
      variant: "default",
      size: "default",
      glow: "none",
    },
  }
)

export interface EnhancedButtonProps
  extends Omit<HTMLMotionProps<"button">, "size">,
    VariantProps<typeof enhancedButtonVariants> {
  asChild?: boolean
  loading?: boolean
  loadingText?: string
  ripple?: boolean
  pulse?: boolean
}

const EnhancedButton = React.forwardRef<HTMLButtonElement, EnhancedButtonProps>(
  ({ 
    className, 
    variant, 
    size, 
    glow,
    asChild = false, 
    loading = false,
    loadingText,
    ripple = true,
    pulse = false,
    children,
    disabled,
    onClick,
    ...props 
  }, ref) => {
    const [ripples, setRipples] = React.useState<Array<{ id: number; x: number; y: number }>>([])
    const rippleId = React.useRef(0)

    const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
      if (ripple && !disabled && !loading) {
        const button = event.currentTarget
        const rect = button.getBoundingClientRect()
        const x = event.clientX - rect.left
        const y = event.clientY - rect.top
        
        const newRipple = { id: rippleId.current++, x, y }
        setRipples(prev => [...prev, newRipple])
        
        setTimeout(() => {
          setRipples(prev => prev.filter(r => r.id !== newRipple.id))
        }, 600)
      }
      
      if (onClick && !disabled && !loading) {
        onClick(event)
      }
    }

    if (asChild) {
      return (
        <Slot
          className={cn(
            enhancedButtonVariants({ variant, size, glow, className }),
            loading && "cursor-not-allowed",
            pulse && "animate-pulse-glow"
          )}
          ref={ref}
          onClick={handleClick}
        >
          <>{children}</>
        </Slot>
      )
    }

    return (
      <motion.button
        className={cn(
          enhancedButtonVariants({ variant, size, glow, className }),
          loading && "cursor-not-allowed",
          pulse && "animate-pulse-glow"
        )}
        ref={ref}
        disabled={disabled || loading}
        onClick={handleClick}
        whileHover={{ scale: disabled || loading ? 1 : 1.05 }}
        whileTap={{ scale: disabled || loading ? 1 : 0.95 }}
        transition={{ type: "spring", stiffness: 400, damping: 17 }}
        {...props}
      >
        {/* Ripple effects */}
        {ripples.map((ripple) => (
          <motion.span
            key={ripple.id}
            className="absolute bg-white/30 rounded-full pointer-events-none"
            style={{
              left: ripple.x - 10,
              top: ripple.y - 10,
              width: 20,
              height: 20,
            }}
            initial={{ scale: 0, opacity: 1 }}
            animate={{ scale: 4, opacity: 0 }}
            transition={{ duration: 0.6, ease: "easeOut" }}
          />
        ))}

        {/* Loading state */}
        {loading && (
          <motion.div
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            className="flex items-center gap-2"
          >
            <Loader2 className="h-4 w-4 animate-spin" />
            {loadingText && <span>{loadingText}</span>}
          </motion.div>
        )}

        {/* Normal content */}
        {!loading && <>{children}</>}

        {/* Shimmer effect for premium variant */}
        {variant === "premium" && (
          <div className="absolute inset-0 bg-gradient-to-r from-transparent via-white/20 to-transparent animate-shimmer" />
        )}
      </motion.button>
    )
  }
)

EnhancedButton.displayName = "EnhancedButton"

export { EnhancedButton, enhancedButtonVariants }
