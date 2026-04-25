import type { NextConfig } from "next";

const nextConfig: NextConfig = {
  transpilePackages: ["@complai/ui-components", "@complai/shared-kernel"],
};

export default nextConfig;
