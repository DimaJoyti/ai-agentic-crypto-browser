/** @type {import('next').NextConfig} */
const nextConfig = {
  // Configure for Cloudflare Pages deployment
  output: process.env.NODE_ENV === 'production' ? 'export' : undefined,
  trailingSlash: true,
  skipTrailingSlashRedirect: true,

  // Reduce experimental features that might cause chunk issues
  experimental: {
    optimizeCss: false, // Disable to prevent chunk conflicts
    // Remove optimizePackageImports temporarily
  },
  // Use stable build ID to prevent chunk loading issues
  generateBuildId: async () => {
    return 'ai-browser-stable'
  },

  // Add chunk loading timeout and retry configuration
  onDemandEntries: {
    // Period (in ms) where the server will keep pages in the buffer
    maxInactiveAge: 25 * 1000,
    // Number of pages that should be kept simultaneously without being disposed
    pagesBufferLength: 2,
  },

  // PWA and Performance Headers
  async headers() {
    return [
      {
        source: '/sw.js',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=0, must-revalidate',
          },
          {
            key: 'Service-Worker-Allowed',
            value: '/',
          },
        ],
      },
      {
        source: '/manifest.json',
        headers: [
          {
            key: 'Cache-Control',
            value: 'public, max-age=31536000, immutable',
          },
        ],
      },
      {
        source: '/(.*)',
        headers: [
          {
            key: 'X-Frame-Options',
            value: 'DENY',
          },
          {
            key: 'X-Content-Type-Options',
            value: 'nosniff',
          },
          {
            key: 'Referrer-Policy',
            value: 'origin-when-cross-origin',
          },
        ],
      },
    ]
  },

  // Compression
  compress: true,
  env: {
    NEXT_PUBLIC_API_URL: process.env.NEXT_PUBLIC_API_URL || (process.env.NODE_ENV === 'production' ? 'https://ai-crypto-browser-api.gcp-inspiration.workers.dev' : 'http://localhost:8080'),
    NEXT_PUBLIC_WS_URL: process.env.NEXT_PUBLIC_WS_URL || (process.env.NODE_ENV === 'production' ? 'wss://ai-crypto-browser-api.gcp-inspiration.workers.dev' : 'ws://localhost:8080'),
    NEXT_PUBLIC_CHAIN_ID: process.env.NEXT_PUBLIC_CHAIN_ID || '1',
    NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID: process.env.NEXT_PUBLIC_WALLETCONNECT_PROJECT_ID || '',
    NEXT_PUBLIC_CLOUDFLARE_DEPLOYMENT: process.env.NEXT_PUBLIC_CLOUDFLARE_DEPLOYMENT || 'false',
  },
  images: {
    unoptimized: process.env.NODE_ENV === 'production', // Disable optimization for static export
    domains: ['localhost', 'your-domain.com'],
    formats: ['image/webp', 'image/avif'],
    deviceSizes: [640, 750, 828, 1080, 1200, 1920, 2048, 3840],
    imageSizes: [16, 32, 48, 64, 96, 128, 256, 384],
  },
  webpack: (config, { dev, isServer }) => {
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      net: false,
      tls: false,
    };

    // Improve chunk loading reliability
    if (!isServer) {
      // Configure chunk loading with better error handling
      config.output = {
        ...config.output,
        // Add chunk loading timeout
        chunkLoadTimeout: 30000, // 30 seconds
        // Improve chunk loading global variable
        chunkLoadingGlobal: 'webpackChunkAiBrowser',
      };

      // Optimize chunk splitting for better loading
      config.optimization = {
        ...config.optimization,
        splitChunks: {
          chunks: 'all',
          minSize: 20000,
          maxSize: 244000,
          cacheGroups: {
            default: {
              minChunks: 1,
              priority: -20,
              reuseExistingChunk: true,
            },
            vendor: {
              test: /[\\/]node_modules[\\/]/,
              name: 'vendors',
              priority: -10,
              chunks: 'all',
              enforce: true,
            },
            // Separate chunk for large libraries
            react: {
              test: /[\\/]node_modules[\\/](react|react-dom)[\\/]/,
              name: 'react',
              priority: 20,
              chunks: 'all',
            },
          },
        },
      };
    }

    return config;
  },
}

module.exports = nextConfig
