import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@complai/ui-components", "@complai/shared-kernel"],
  async rewrites() {
    return [
      {
        source: "/api/v1/:path*",
        destination: `${process.env.BFF_URL ?? "http://localhost:4000"}/api/v1/:path*`,
      },
    ];
  },
};

export default nextConfig;
