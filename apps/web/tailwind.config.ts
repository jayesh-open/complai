import type { Config } from "tailwindcss";

const config: Config = {
  content: [
    "./src/**/*.{ts,tsx}",
    "../../packages/ui-components/src/**/*.{ts,tsx}",
  ],
  theme: {
    extend: {
      fontFamily: {
        sans: [
          "-apple-system",
          "BlinkMacSystemFont",
          "Segoe UI",
          "Noto Sans",
          "Helvetica",
          "Arial",
          "sans-serif",
          "Apple Color Emoji",
          "Segoe UI Emoji",
        ],
        mono: [
          "SF Mono",
          "Cascadia Code",
          "Fira Code",
          "JetBrains Mono",
          "ui-monospace",
          "SFMono-Regular",
          "Menlo",
          "Consolas",
          "monospace",
        ],
      },
      colors: {
        app: {
          bg: "var(--bg-primary)",
          card: "var(--bg-secondary)",
          input: "var(--bg-tertiary)",
          sidebar: "var(--bg-sidebar)",
          border: "var(--border-default)",
          "border-lt": "var(--border-light)",
          accent: "var(--accent)",
          "accent-h": "var(--accent-hover)",
          "accent-m": "var(--accent-muted)",
          "accent-b": "var(--accent-border)",
          "accent-t": "var(--accent-text)",
          success: "var(--success)",
          "success-m": "var(--success-muted)",
          "success-b": "var(--success-border)",
          danger: "var(--danger)",
          "danger-m": "var(--danger-muted)",
          "danger-b": "var(--danger-border)",
          warning: "var(--warning)",
          "warning-m": "var(--warning-muted)",
          "warning-b": "var(--warning-border)",
          info: "var(--info)",
          "info-m": "var(--info-muted)",
          "info-b": "var(--info-border)",
          purple: "var(--purple)",
          "purple-m": "var(--purple-muted)",
          teal: "var(--teal)",
          "teal-m": "var(--teal-muted)",
          orange: "var(--orange)",
          "orange-m": "var(--orange-muted)",
        },
        foreground: {
          DEFAULT: "var(--text-primary)",
          secondary: "var(--text-secondary)",
          muted: "var(--text-muted)",
          disabled: "var(--text-disabled)",
        },
      },
      boxShadow: {
        "app-sm": "var(--shadow-sm)",
        "app-md": "var(--shadow-md)",
        "app-lg": "var(--shadow-lg)",
        "app-accent": "var(--shadow-accent)",
      },
      borderRadius: {
        card: "16px",
        btn: "10px",
        badge: "6px",
      },
    },
  },
  plugins: [],
};

export default config;
