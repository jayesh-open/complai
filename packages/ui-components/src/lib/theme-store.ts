import { create } from 'zustand';
import { persist } from 'zustand/middleware';
import { type ThemeMode, DEFAULT_THEME } from './themes';

export type DensityMode = 'compact' | 'comfortable' | 'spacious';

interface ThemeState {
  theme: ThemeMode;
  density: DensityMode;
  sidebarCollapsed: boolean;
  setTheme: (theme: ThemeMode) => void;
  setDensity: (density: DensityMode) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
}

export const useThemeStore = create<ThemeState>()(
  persist(
    (set) => ({
      theme: DEFAULT_THEME,
      density: 'compact',
      sidebarCollapsed: false,
      setTheme: (theme) => set({ theme }),
      setDensity: (density) => set({ density }),
      toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
    }),
    { name: 'complai-theme' }
  )
);
