/** @type {import('next').NextConfig} */
const nextConfig = {
  experimental: {
    // appDir is now stable in Next.js 14, no longer needed
  },
  output: 'standalone',
  // Disable static optimization for authenticated pages
  trailingSlash: false,
  generateBuildId: async () => {
    return 'build-' + Date.now()
  },
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080',
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL || 'ws://localhost:8080',
    NEXT_PUBLIC_CHAIN_ID: process.env.NEXT_PUBLIC_CHAIN_ID || '1',
    NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID: process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || '',
  },
  images: {
    domains: ['localhost'],
  },
  webpack: (config) => {
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      net: false,
      tls: false,
    };
    return config;
  },
}

module.exports = nextConfig
