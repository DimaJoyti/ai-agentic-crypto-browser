# 🚀 Performance Optimization Guide

This guide covers all the performance optimizations implemented to make the UI run faster during development.

## 🎯 Quick Start

### Development Scripts

```bash
# Fastest development mode with all optimizations
npm run dev:fast

# Clean build cache and start fresh
npm run dev:clean

# Development with debugging enabled
npm run dev:debug

# Clean all caches and reinstall
npm run clean:all
```

### Environment Setup

1. Copy the development environment file:
```bash
cp .env.development .env.local
```

2. Enable performance mode in your browser:
   - Look for the orange ⚡ button in the bottom-right corner
   - Click to open the Performance Optimizer panel
   - Toggle optimizations as needed

## 🔧 Performance Optimizations

### 1. Next.js Configuration

**Development-specific optimizations in `next.config.js`:**

- ✅ **Disabled source maps** for faster builds
- ✅ **Disabled font optimization** in development
- ✅ **Disabled React Strict Mode** for faster refresh
- ✅ **Disabled chunk splitting** in development
- ✅ **Optimized webpack configuration**
- ✅ **Faster file watching**
- ✅ **Reduced bundle analysis**

### 2. Webpack Optimizations

**Development mode:**
- Disabled source maps (`devtool: false`)
- Disabled symlink resolution
- Optimized file watching
- Disabled chunk splitting
- Removed expensive optimizations

**Production mode:**
- Optimized chunk splitting
- Better error handling
- Improved loading timeouts

### 3. CSS & Styling

**Performance CSS classes:**
```css
.performance-mode * {
  animation-duration: 0s !important;
  transition-duration: 0s !important;
}

.reduce-motion * {
  animation: none !important;
  transition: none !important;
}

.performance-optimized {
  transform: translateZ(0);
  will-change: auto;
  contain: layout style paint;
}
```

### 4. Component Optimizations

**React optimizations:**
- Memoization with `React.memo`
- Debounced inputs
- Lazy loading components
- Virtual scrolling for large lists
- Efficient re-render prevention

## 🛠️ Development Tools

### Performance Optimizer Panel

Access via the ⚡ button (development only):

**Metrics displayed:**
- Render time (ms)
- Memory usage (MB)
- Component count
- Re-render count

**Optimization toggles:**
- Disable animations
- Reduce motion
- Enable lazy loading
- Virtual scrolling
- Component memoization
- Input debouncing

### Performance Hook

Use the `useDevOptimizations` hook in your components:

```tsx
import { useDevOptimizations } from '@/hooks/useDevOptimizations'

function MyComponent() {
  const {
    optimizations,
    metrics,
    createDebouncedCallback,
    memoize,
    useLazyLoading,
    PerformanceWrapper
  } = useDevOptimizations()

  // Debounced input handler
  const debouncedSearch = createDebouncedCallback((query: string) => {
    // Search logic
  }, 300)

  // Memoized expensive calculation
  const expensiveCalculation = memoize((data: any[]) => {
    return data.reduce((acc, item) => acc + item.value, 0)
  })

  // Lazy loading
  const { isVisible, setRef } = useLazyLoading()

  return (
    <PerformanceWrapper>
      <div ref={setRef}>
        {isVisible && <ExpensiveComponent />}
      </div>
    </PerformanceWrapper>
  )
}
```

## 📊 Performance Monitoring

### Real-time Metrics

The performance optimizer tracks:

1. **Render Time**: Component render duration
2. **Memory Usage**: JavaScript heap size
3. **FPS**: Frames per second
4. **Re-renders**: Component update frequency

### Performance Status

- 🟢 **Good**: FPS > 50, Memory < 50MB, Render < 8ms
- 🟡 **Fair**: FPS > 30, Memory < 100MB, Render < 16ms
- 🔴 **Poor**: Below fair thresholds

## 🎨 CSS Performance Classes

### Quick Performance Classes

```css
/* Fast loading skeleton */
.fast-skeleton { /* shimmer animation */ }

/* Efficient animations */
.efficient-fade { /* opacity transition */ }
.efficient-slide { /* transform transition */ }

/* GPU acceleration */
.gpu-accelerated { transform: translateZ(0); }

/* Memory efficient */
.simple-shadow { box-shadow: 0 2px 4px rgba(0,0,0,0.1); }
.simple-border { border: 1px solid rgba(0,0,0,0.1); }
```

### Development Helpers

```css
/* Visual debugging */
.dev-only { /* red dashed border */ }
.debug-performance * { /* red outline */ }
.dev-grid { /* grid overlay */ }
```

## 🚀 Best Practices

### 1. Component Optimization

```tsx
// ✅ Good: Memoized component
const OptimizedComponent = React.memo(({ data }) => {
  const memoizedValue = useMemo(() => 
    expensiveCalculation(data), [data]
  )
  
  return <div>{memoizedValue}</div>
})

// ❌ Avoid: Expensive operations in render
const SlowComponent = ({ data }) => {
  const value = expensiveCalculation(data) // Runs every render
  return <div>{value}</div>
}
```

### 2. Event Handling

```tsx
// ✅ Good: Debounced input
const debouncedHandler = createDebouncedCallback(handleSearch, 300)

// ❌ Avoid: Direct expensive operations
const handleInput = (e) => {
  expensiveSearch(e.target.value) // Runs on every keystroke
}
```

### 3. Lazy Loading

```tsx
// ✅ Good: Lazy loaded component
const LazyComponent = lazy(() => import('./HeavyComponent'))

// ✅ Good: Intersection observer
const { isVisible, setRef } = useLazyLoading()
```

## 🔍 Debugging Performance Issues

### 1. Check Performance Panel

1. Open Performance Optimizer (⚡ button)
2. Monitor real-time metrics
3. Identify bottlenecks (high render time, memory usage)
4. Toggle optimizations to test impact

### 2. Browser DevTools

1. **Performance tab**: Record and analyze performance
2. **Memory tab**: Check for memory leaks
3. **Network tab**: Optimize asset loading
4. **React DevTools**: Profile component renders

### 3. Common Issues

**Slow renders:**
- Enable memoization
- Reduce component complexity
- Use lazy loading

**High memory usage:**
- Check for memory leaks
- Optimize large data structures
- Use virtual scrolling

**Low FPS:**
- Disable animations
- Reduce DOM complexity
- Use CSS transforms instead of layout changes

## 📈 Performance Targets

### Development Targets

- **Render time**: < 16ms (60 FPS)
- **Memory usage**: < 100MB
- **Bundle size**: < 1MB (development)
- **Hot reload**: < 2s

### Production Targets

- **First Contentful Paint**: < 1.5s
- **Largest Contentful Paint**: < 2.5s
- **Cumulative Layout Shift**: < 0.1
- **First Input Delay**: < 100ms

## 🛡️ Troubleshooting

### Common Issues

1. **Slow development server**
   ```bash
   npm run dev:clean  # Clear cache and restart
   ```

2. **High memory usage**
   - Enable performance mode
   - Disable animations
   - Use lazy loading

3. **Slow hot reload**
   - Check file watching settings
   - Reduce bundle size
   - Clear Next.js cache

### Reset Performance Settings

```bash
# Clear all performance settings
localStorage.removeItem('dev-optimizations')
localStorage.removeItem('dev-performance-settings')

# Or use the reset button in Performance Optimizer panel
```

## 📚 Additional Resources

- [Next.js Performance](https://nextjs.org/docs/advanced-features/measuring-performance)
- [React Performance](https://react.dev/learn/render-and-commit)
- [Web Vitals](https://web.dev/vitals/)
- [Chrome DevTools Performance](https://developer.chrome.com/docs/devtools/performance/)
