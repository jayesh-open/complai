import type { Preview } from "@storybook/react";
import "../src/app/globals.css";
import { THEME_MAP } from "@complai/ui-components";

const lightColors = THEME_MAP.light;

const preview: Preview = {
  decorators: [
    (Story) => {
      const root = document.documentElement;
      Object.entries(lightColors).forEach(([key, value]) => {
        const cssVar = `--${key.replace(/([A-Z])/g, "-$1").toLowerCase()}`;
        root.style.setProperty(cssVar, value);
      });
      root.setAttribute("data-theme", "light");
      root.setAttribute("data-theme-family", "light");
      return Story();
    },
  ],
  parameters: {
    controls: { matchers: { color: /(background|color)$/i, date: /Date$/i } },
  },
};

export default preview;
