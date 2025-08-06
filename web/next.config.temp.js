/** @type {import('next').NextConfig} */
const nextConfig = {
  // Minimal config for debugging
  trailingSlash: false,
  generateBuildId: async () => {
    return 'debug-build'
  },
  webpack: (config, { dev }) => {
    // Basic fallbacks
    config.resolve.fallback = {
      ...config.resolve.fallback,
      fs: false,
      net: false,
      tls: false,
    };

    // Disable chunk optimization in dev mode
    if (dev) {
      config.optimization = {
        ...config.optimization,
        splitChunks: false,
      };
    }

    return config;
  },
}

module.exports = nextConfig
