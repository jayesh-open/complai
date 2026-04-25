import type { StorybookConfig } from "@storybook/react-vite";
import path from "path";

const config: StorybookConfig = {
  stories: ["../src/**/*.stories.@(ts|tsx)"],
  addons: [
    "@storybook/addon-essentials",
    "@storybook/addon-a11y",
  ],
  framework: {
    name: "@storybook/react-vite",
    options: {},
  },
  viteFinal: async (config) => {
    const repoRoot = path.resolve(__dirname, "../../..");

    config.resolve = config.resolve || {};
    config.resolve.alias = {
      ...config.resolve.alias,
      "@": path.resolve(__dirname, "../src"),
      "@complai/ui-components": path.resolve(repoRoot, "packages/ui-components/src/index.ts"),
    };

    config.server = config.server || {};
    config.server.fs = config.server.fs || {};
    config.server.fs.allow = [
      ...(config.server.fs.allow || []),
      repoRoot,
    ];

    config.esbuild = {
      ...config.esbuild,
      jsx: "automatic",
    };

    return config;
  },
};

export default config;
