/** @type {import('next').NextConfig} */
const nextConfig = {
  reactStrictMode: true,

  // trailingSlash: true,
  async rewrites() {
    return [
      {
        source: '/v1/:path*',
        destination: `${process.env.API_HOST || 'http://localhost:8080'}/v1/:path*`,
      }
    ]
  }
}

module.exports = nextConfig
