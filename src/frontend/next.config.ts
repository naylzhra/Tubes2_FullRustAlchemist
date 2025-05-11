import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  async rewrites() {
    return [
      {
        source: "/api/:path*",          // request from browser
        destination: "http://localhost:8080/api/:path*", // Gin server
      },
    ];
  },
  env: {
    NEXT_PUBLIC_BACKEND_URL: "http://localhost:8080",
  },
};

export default nextConfig;
