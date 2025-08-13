/** @type {import('next').NextConfig} */
const nextConfig = {
  // Configure for Cloudflare Pages deployment
  output: process.env.NODE_ENV === 'production' ? 'export' : undefined,
  trailingSlash: true,
  skipTrailingSlashRedirect: true,

  // Development optimizations
  ...(process.env.NODE_ENV === 'development' && {
    // Faster development builds
    swcMinify: false,
    // Disable source maps in development for faster builds
    productionBrowserSourceMaps: false,
    // Optimize for development
    optimizeFonts: false,
    // Faster refresh
    reactStrictMode: false,
  }),

  // Reduce experimental features that might cause chunk issues
  experimental: {
    optimizeCss: false, // Disable to prevent chunk conflicts
    // Development optimizations
    ...(process.env.NODE_ENV === 'development' && {
      // Faster compilation
      turbo: {
        rules: {
          '*.svg': {
            loaders: ['@svgr/webpack'],
            as: '*.js',
          },
        },
      },
      // Disable expensive optimizations in dev
      optimizePackageImports: [],
      // Faster builds
      webVitalsAttribution: [],
    }),
  },

  // Use stable build ID to prevent chunk loading issues
  generateBuildId: async () => {
    return process.env.NODE_ENV === 'development' ? 'dev-build' : 'ai-browser-stable'
  },

  // Optimized development settings
  onDemandEntries: {
    // Period (in ms) where the server will keep pages in the buffer
    maxInactiveAge: process.env.NODE_ENV === 'development' ? 60 * 1000 : 25 * 1000,
    // Number of pages that should be kept simultaneously without being disposed
    pagesBufferLength: process.env.NODE_ENV === 'development' ? 5 : 2,
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

    // Development-specific optimizations for faster builds
    if (dev) {
      // Disable source maps in development for faster builds
      config.devtool = false;

      // Optimize module resolution for faster builds
      config.resolve.symlinks = false;
      config.resolve.cacheWithContext = false;

      // Faster file watching
      config.watchOptions = {
        poll: false,
        ignored: /node_modules/,
      };

      // Disable chunk splitting in development for faster builds
      config.optimization = {
        ...config.optimization,
        splitChunks: false,
        removeAvailableModules: false,
        removeEmptyChunks: false,
        mergeDuplicateChunks: false,
      };
    }

    // Production optimizations
    if (!isServer && !dev) {
      // Configure chunk loading with better error handling
      config.output = {
        ...config.output,
        // Add chunk loading timeout
        chunkLoadTimeout: 30000, // 30 seconds
        // Improve chunk loading global variable
        chunkLoadingGlobal: 'webpackChunkAiBrowser',
      };

      // Optimize chunk splitting for better loading (production only)
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
