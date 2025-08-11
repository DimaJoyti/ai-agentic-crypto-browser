# 🔧 Troubleshooting Guide - AI-Agentic Crypto Browser

## 🚨 Common Deployment Issues

### **NPM Cache Issues in CI/CD**

#### **Error:** 
```
Error: Some specified paths were not resolved, unable to cache dependencies.
/opt/hostedtoolcache/node/18.20.8/x64/bin/npm config get cache
/home/runner/.npm
```

#### **Root Cause:**
This error occurs in GitHub Actions or other CI environments when the npm cache path cannot be resolved or accessed properly.

#### **Solutions:**

##### **1. Use Custom Deployment Script**
```bash
# Run the enhanced deployment script
./scripts/deploy-cloudflare.sh
```

##### **2. Manual CI Fix**
```yaml
# In your GitHub Actions workflow
- name: Setup Node.js with custom cache
  uses: actions/setup-node@v4
  with:
    node-version: '18'
    cache: 'npm'
    cache-dependency-path: 'web/package-lock.json'

- name: Configure npm cache
  run: |
    mkdir -p /tmp/.npm-cache
    npm config set cache /tmp/.npm-cache
    npm config set prefer-offline true
    npm config set audit false
    npm config set fund false

- name: Install dependencies
  working-directory: ./web
  run: |
    npm ci --prefer-offline --no-audit --progress=false
  env:
    NPM_CONFIG_CACHE: /tmp/.npm-cache
```

##### **3. Local Environment Fix**
```bash
# Clear npm cache
npm cache clean --force

# Set custom cache directory
export NPM_CONFIG_CACHE=/tmp/.npm-cache
mkdir -p /tmp/.npm-cache

# Install dependencies
cd web
npm ci --prefer-offline --no-audit
```

---

## 🔄 Build Issues

### **Next.js Build Failures**

#### **Error:** Hydration mismatches or SSR errors

#### **Solution:**
```bash
# Use the production build script
cd web
npm run build:cloudflare
```

#### **Environment Variables:**
```bash
export NODE_ENV=production
export NEXT_PUBLIC_CLOUDFLARE_DEPLOYMENT=true
export NEXT_PUBLIC_API_BASE_URL=https://ai-crypto-browser-api.gcp-inspiration.workers.dev
```

---

## 🌐 Cloudflare Pages Issues

### **Deployment Failures**

#### **Error:** Wrangler authentication issues

#### **Solution:**
```bash
# Login to Cloudflare
wrangler auth login

# Deploy with explicit project name
wrangler pages deploy web/out --project-name=ai-agentic-crypto-browser
```

### **Build Output Issues**

#### **Error:** No `out` directory found

#### **Solution:**
```bash
# Ensure Next.js is configured for static export
# Check next.config.js has:
output: 'export'

# Build and verify output
npm run build:cloudflare
ls -la web/out/
```

---

## 🔐 Authentication & API Issues

### **API Connection Failures**

#### **Error:** API endpoints not responding

#### **Check:**
1. **API Base URL**: Verify `NEXT_PUBLIC_API_BASE_URL` is correct
2. **CORS Settings**: Ensure API allows requests from your domain
3. **SSL/TLS**: Check certificate validity

#### **Solution:**
```javascript
// In your API calls, add error handling
try {
  const response = await fetch(`${process.env.NEXT_PUBLIC_API_BASE_URL}/endpoint`);
  if (!response.ok) {
    throw new Error(`HTTP ${response.status}: ${response.statusText}`);
  }
  return await response.json();
} catch (error) {
  console.error('API Error:', error);
  // Fallback logic
}
```

---

## 📱 Runtime Issues

### **Wallet Connection Problems**

#### **Error:** WalletConnect initialization failures

#### **Solution:**
```javascript
// Check environment variables
const projectId = process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID;
if (!projectId || projectId === 'your_walletconnect_project_id_here') {
  console.warn('WalletConnect project ID not configured');
  // Use fallback or disable wallet features
}
```

### **Hydration Mismatches**

#### **Error:** Text content does not match server-rendered HTML

#### **Solution:**
```javascript
// Use client-side only rendering for problematic components
import dynamic from 'next/dynamic';

const ClientOnlyComponent = dynamic(
  () => import('./ClientOnlyComponent'),
  { ssr: false }
);
```

---

## 🛠️ Development Environment

### **Local Development Issues**

#### **Error:** Module not found or import errors

#### **Solution:**
```bash
# Clear all caches and reinstall
rm -rf web/node_modules web/.next
cd web
npm install
npm run dev
```

### **Environment Variables**

#### **Missing or incorrect environment variables**

#### **Check:**
```bash
# Verify .env.local exists and has correct values
cat web/.env.local

# Required variables:
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
NEXT_PUBLIC_WS_URL=ws://localhost:8080/ws
NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID=your_project_id
```

---

## 🚀 Performance Issues

### **Large Bundle Sizes**

#### **Solution:**
```javascript
// Use dynamic imports for large components
const HeavyComponent = dynamic(() => import('./HeavyComponent'), {
  loading: () => <div>Loading...</div>
});

// Optimize images
import Image from 'next/image';
<Image src="/image.jpg" alt="Description" width={500} height={300} />
```

### **Slow Loading Times**

#### **Check:**
1. **Bundle analysis**: `npm run analyze`
2. **Image optimization**: Use WebP format
3. **Code splitting**: Implement lazy loading
4. **CDN caching**: Verify Cloudflare cache settings

---

## 📞 Getting Help

### **Debug Information to Collect:**

1. **Environment Details:**
   ```bash
   node --version
   npm --version
   npx next --version
   ```

2. **Build Logs:**
   ```bash
   npm run build:cloudflare > build.log 2>&1
   ```

3. **Network Information:**
   ```bash
   curl -I https://01becb04.ai-agentic-crypto-browser.pages.dev
   ```

### **Useful Commands:**

```bash
# Check deployment status
wrangler pages deployment list --project-name=ai-agentic-crypto-browser

# View build logs
npm run build:cloudflare --verbose

# Test API connectivity
curl -v https://ai-crypto-browser-api.gcp-inspiration.workers.dev/health

# Check environment variables
printenv | grep NEXT_PUBLIC
```

---

## ✅ Quick Fixes Checklist

- [ ] **Clear npm cache**: `npm cache clean --force`
- [ ] **Update dependencies**: `npm update`
- [ ] **Check environment variables**: Verify all required vars are set
- [ ] **Rebuild application**: `npm run build:cloudflare`
- [ ] **Test locally first**: `npm run dev`
- [ ] **Check API endpoints**: Verify backend is running
- [ ] **Review build output**: Ensure `out/` directory exists
- [ ] **Validate deployment**: Test live URL functionality

---

## 🔄 Emergency Rollback

If deployment fails and you need to rollback:

```bash
# Rollback to previous deployment (Cloudflare Dashboard)
# Or redeploy from a previous commit:
git checkout <previous-commit-hash>
./scripts/deploy-cloudflare.sh
```

---

**💡 Pro Tip:** Always test locally before deploying to production!
