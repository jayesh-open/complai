"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import { type ThemeMode, DEFAULT_THEME } from "@/lib/themes";

export type DensityMode = "compact" | "comfortable" | "spacious";

interface AppState {
  theme: ThemeMode;
  density: DensityMode;
  sidebarCollapsed: boolean;
  sidebarGroupState: Record<string, boolean>;
  setTheme: (theme: ThemeMode) => void;
  setDensity: (density: DensityMode) => void;
  toggleSidebar: () => void;
  setSidebarCollapsed: (collapsed: boolean) => void;
  toggleSidebarGroup: (group: string) => void;
}

export const useAppStore = create<AppState>()(
  persist(
    (set) => ({
      theme: DEFAULT_THEME,
      density: "compact",
      sidebarCollapsed: false,
      sidebarGroupState: {},
      setTheme: (theme) => set({ theme }),
      setDensity: (density) => set({ density }),
      toggleSidebar: () => set((s) => ({ sidebarCollapsed: !s.sidebarCollapsed })),
      setSidebarCollapsed: (collapsed) => set({ sidebarCollapsed: collapsed }),
      toggleSidebarGroup: (group) =>
        set((s) => ({
          sidebarGroupState: {
            ...s.sidebarGroupState,
            [group]: !s.sidebarGroupState[group],
          },
        })),
    }),
    { name: "complai-app-store" },
  ),
);
