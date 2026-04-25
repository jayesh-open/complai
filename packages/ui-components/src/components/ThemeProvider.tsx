"use client";
import { useEffect, type ReactNode } from 'react';
import { useThemeStore } from '../lib/theme-store';
import { THEME_MAP, getThemeFamily, type ThemeColors } from '../lib/themes';

export function ThemeProvider({ children }: { children: ReactNode }) {
  const { theme } = useThemeStore();

  useEffect(() => {
    const colors: ThemeColors = THEME_MAP[theme] || THEME_MAP.light;
    const root = document.documentElement;
    Object.entries(colors).forEach(([key, value]) => {
      const cssVar = `--${key.replace(/([A-Z])/g, '-$1').toLowerCase()}`;
      root.style.setProperty(cssVar, value);
    });
    root.setAttribute('data-theme', theme);
    root.setAttribute('data-theme-family', getThemeFamily(theme));
  }, [theme]);

  return <>{children}</>;
}
