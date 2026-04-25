export type ThemeMode =
  | "dark" | "dark-midnight" | "dark-purple"
  | "light" | "light-sky" | "light-pastel"
  | "ocean" | "ocean-cyan" | "ocean-pacific"
  | "sunset" | "sunset-rose" | "sunset-dusk"
  | "forest" | "forest-mint" | "forest-pine";

export interface ThemeColors {
  bgPrimary: string;
  bgSecondary: string;
  bgTertiary: string;
  bgSidebar: string;
  bgOverlay: string;
  borderDefault: string;
  borderLight: string;
  borderFocus: string;
  textPrimary: string;
  textSecondary: string;
  textMuted: string;
  textDisabled: string;
  accent: string;
  accentHover: string;
  accentMuted: string;
  accentBorder: string;
  accentText: string;
  success: string;
  successMuted: string;
  successBorder: string;
  danger: string;
  dangerMuted: string;
  dangerBorder: string;
  warning: string;
  warningMuted: string;
  warningBorder: string;
  info: string;
  infoMuted: string;
  infoBorder: string;
  purple: string;
  purpleMuted: string;
  pink: string;
  pinkMuted: string;
  teal: string;
  tealMuted: string;
  orange: string;
  orangeMuted: string;
  shadowSm: string;
  shadowMd: string;
  shadowLg: string;
  shadowAccent: string;
}

const darkTheme: ThemeColors = {
  bgPrimary: "#141519", bgSecondary: "#1c1e25", bgTertiary: "#25272f",
  bgSidebar: "#111318", bgOverlay: "rgba(0,0,0,0.65)",
  borderDefault: "#323540", borderLight: "#3d414c", borderFocus: "#f59e0b",
  textPrimary: "#e8e9ed", textSecondary: "#c8ccd4", textMuted: "#9ca0ab", textDisabled: "#3d4150",
  accent: "#f59e0b", accentHover: "#d97706", accentMuted: "rgba(245,158,11,0.10)",
  accentBorder: "rgba(245,158,11,0.22)", accentText: "#0a0b0e",
  success: "#10b981", successMuted: "rgba(16,185,129,0.10)", successBorder: "rgba(16,185,129,0.25)",
  danger: "#ef4444", dangerMuted: "rgba(239,68,68,0.10)", dangerBorder: "rgba(239,68,68,0.25)",
  warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  info: "#3b82f6", infoMuted: "rgba(59,130,246,0.10)", infoBorder: "rgba(59,130,246,0.25)",
  purple: "#8b5cf6", purpleMuted: "rgba(139,92,246,0.10)",
  pink: "#ec4899", pinkMuted: "rgba(236,72,153,0.10)",
  teal: "#14b8a6", tealMuted: "rgba(20,184,166,0.10)",
  orange: "#f97316", orangeMuted: "rgba(249,115,22,0.10)",
  shadowSm: "0 1px 3px rgba(0,0,0,0.4)", shadowMd: "0 4px 12px rgba(0,0,0,0.5)",
  shadowLg: "0 8px 32px rgba(0,0,0,0.6)", shadowAccent: "0 4px 16px rgba(245,158,11,0.25)",
};

const lightTheme: ThemeColors = {
  bgPrimary: "#faf8f5", bgSecondary: "#ffffff", bgTertiary: "#f5f3f0",
  bgSidebar: "#f0ebe3", bgOverlay: "rgba(90,70,50,0.25)",
  borderDefault: "#e0d8cc", borderLight: "#eae3d8", borderFocus: "#d97706",
  textPrimary: "#1a1612", textSecondary: "#4a4239", textMuted: "#6b6358", textDisabled: "#c4bdb4",
  accent: "#d97706", accentHover: "#b45309", accentMuted: "rgba(217,119,6,0.08)",
  accentBorder: "rgba(217,119,6,0.20)", accentText: "#ffffff",
  success: "#059669", successMuted: "rgba(5,150,105,0.08)", successBorder: "rgba(5,150,105,0.20)",
  danger: "#dc2626", dangerMuted: "rgba(220,38,38,0.07)", dangerBorder: "rgba(220,38,38,0.18)",
  warning: "#d97706", warningMuted: "rgba(217,119,6,0.08)", warningBorder: "rgba(217,119,6,0.18)",
  info: "#2563eb", infoMuted: "rgba(37,99,235,0.07)", infoBorder: "rgba(37,99,235,0.18)",
  purple: "#7c3aed", purpleMuted: "rgba(124,58,237,0.07)",
  pink: "#db2777", pinkMuted: "rgba(219,39,119,0.07)",
  teal: "#0d9488", tealMuted: "rgba(13,148,136,0.07)",
  orange: "#ea580c", orangeMuted: "rgba(234,88,12,0.07)",
  shadowSm: "0 1px 3px rgba(140,120,100,0.08)", shadowMd: "0 4px 12px rgba(140,120,100,0.10)",
  shadowLg: "0 8px 32px rgba(140,120,100,0.12)", shadowAccent: "0 4px 16px rgba(217,119,6,0.18)",
};

function darkVariant(base: ThemeColors, overrides: Partial<ThemeColors>): ThemeColors {
  return { ...base, ...overrides };
}

function lightVariant(base: ThemeColors, overrides: Partial<ThemeColors>): ThemeColors {
  return { ...base, ...overrides };
}

export const THEME_MAP: Record<ThemeMode, ThemeColors> = {
  dark: darkTheme,
  "dark-midnight": darkVariant(darkTheme, {
    bgPrimary: "#0f1420", bgSecondary: "#171d2c", bgTertiary: "#1e2538",
    bgSidebar: "#0c1018", borderDefault: "#253050", borderLight: "#304060",
    accent: "#3b82f6", accentHover: "#2563eb", accentMuted: "rgba(59,130,246,0.10)",
    accentBorder: "rgba(59,130,246,0.22)", accentText: "#ffffff",
    borderFocus: "#3b82f6", shadowAccent: "0 4px 16px rgba(59,130,246,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "dark-purple": darkVariant(darkTheme, {
    bgPrimary: "#15131e", bgSecondary: "#1e1b2a", bgTertiary: "#272435",
    bgSidebar: "#110f18", borderDefault: "#352e50", borderLight: "#403858",
    accent: "#8b5cf6", accentHover: "#7c3aed", accentMuted: "rgba(139,92,246,0.10)",
    accentBorder: "rgba(139,92,246,0.22)", accentText: "#ffffff",
    borderFocus: "#8b5cf6", shadowAccent: "0 4px 16px rgba(139,92,246,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  light: lightTheme,
  "light-sky": lightVariant(lightTheme, {
    bgPrimary: "#f5f9ff", bgSecondary: "#ffffff", bgTertiary: "#edf3fc",
    bgSidebar: "#e3edf8", borderDefault: "#c8d8ec", borderLight: "#dce6f2",
    accent: "#0284c7", accentHover: "#0369a1", accentMuted: "rgba(2,132,199,0.08)",
    accentBorder: "rgba(2,132,199,0.20)", accentText: "#ffffff",
    borderFocus: "#0284c7", shadowAccent: "0 4px 16px rgba(2,132,199,0.18)",
    warning: "#d97706", warningMuted: "rgba(217,119,6,0.08)", warningBorder: "rgba(217,119,6,0.18)",
  }),
  "light-pastel": lightVariant(lightTheme, {
    bgPrimary: "#fdf5f5", bgSecondary: "#ffffff", bgTertiary: "#f9eeee",
    bgSidebar: "#f5e8e8", borderDefault: "#e8d0d0", borderLight: "#f0dede",
    accent: "#e11d48", accentHover: "#be123c", accentMuted: "rgba(225,29,72,0.07)",
    accentBorder: "rgba(225,29,72,0.18)", accentText: "#ffffff",
    borderFocus: "#e11d48", shadowAccent: "0 4px 16px rgba(225,29,72,0.18)",
    warning: "#d97706", warningMuted: "rgba(217,119,6,0.08)", warningBorder: "rgba(217,119,6,0.18)",
  }),
  ocean: darkVariant(darkTheme, {
    bgPrimary: "#152232", bgSecondary: "#1c2a40", bgTertiary: "#243650",
    bgSidebar: "#111e30", bgOverlay: "rgba(0,15,30,0.65)",
    borderDefault: "#243a5a", borderLight: "#2e4868", borderFocus: "#38bdf8",
    textPrimary: "#e8edf3", textSecondary: "#b0c4de", textMuted: "#7a9bbd", textDisabled: "#354860",
    accent: "#38bdf8", accentHover: "#0ea5e9", accentMuted: "rgba(56,189,248,0.10)",
    accentBorder: "rgba(56,189,248,0.22)", accentText: "#0a1520",
    shadowAccent: "0 4px 16px rgba(56,189,248,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "ocean-cyan": darkVariant(darkTheme, {
    bgPrimary: "#0f2028", bgSecondary: "#162a35", bgTertiary: "#1e3542",
    bgSidebar: "#0b1a22", bgOverlay: "rgba(0,15,25,0.65)",
    borderDefault: "#1e3a4a", borderLight: "#284858", borderFocus: "#06b6d4",
    textPrimary: "#e0f0f5", textSecondary: "#a0c8d8", textMuted: "#6a9ab0", textDisabled: "#2a4050",
    accent: "#06b6d4", accentHover: "#0891b2", accentMuted: "rgba(6,182,212,0.10)",
    accentBorder: "rgba(6,182,212,0.22)", accentText: "#0a1518",
    shadowAccent: "0 4px 16px rgba(6,182,212,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "ocean-pacific": darkVariant(darkTheme, {
    bgPrimary: "#101828", bgSecondary: "#182238", bgTertiary: "#202e48",
    bgSidebar: "#0c1420", bgOverlay: "rgba(0,10,25,0.65)",
    borderDefault: "#203050", borderLight: "#2a3e60", borderFocus: "#6366f1",
    textPrimary: "#e4e8f0", textSecondary: "#a8b8d0", textMuted: "#7088a8", textDisabled: "#303e55",
    accent: "#6366f1", accentHover: "#4f46e5", accentMuted: "rgba(99,102,241,0.10)",
    accentBorder: "rgba(99,102,241,0.22)", accentText: "#ffffff",
    shadowAccent: "0 4px 16px rgba(99,102,241,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  sunset: darkVariant(darkTheme, {
    bgPrimary: "#1e1510", bgSecondary: "#2a1f18", bgTertiary: "#352820",
    bgSidebar: "#211510", bgOverlay: "rgba(20,10,0,0.65)",
    borderDefault: "#452e24", borderLight: "#553a2e", borderFocus: "#f97316",
    textPrimary: "#f0e6e0", textSecondary: "#d4b8aa", textMuted: "#a08070", textDisabled: "#4a3528",
    accent: "#f97316", accentHover: "#ea580c", accentMuted: "rgba(249,115,22,0.10)",
    accentBorder: "rgba(249,115,22,0.22)", accentText: "#0e0a08",
    shadowAccent: "0 4px 16px rgba(249,115,22,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "sunset-rose": darkVariant(darkTheme, {
    bgPrimary: "#1e1215", bgSecondary: "#2a1c20", bgTertiary: "#352428",
    bgSidebar: "#1a1012", bgOverlay: "rgba(20,5,10,0.65)",
    borderDefault: "#45282e", borderLight: "#55343a", borderFocus: "#f43f5e",
    textPrimary: "#f0e0e4", textSecondary: "#d4aab4", textMuted: "#a07080", textDisabled: "#4a2830",
    accent: "#f43f5e", accentHover: "#e11d48", accentMuted: "rgba(244,63,94,0.10)",
    accentBorder: "rgba(244,63,94,0.22)", accentText: "#ffffff",
    shadowAccent: "0 4px 16px rgba(244,63,94,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "sunset-dusk": darkVariant(darkTheme, {
    bgPrimary: "#1a1318", bgSecondary: "#251e24", bgTertiary: "#30262e",
    bgSidebar: "#161015", bgOverlay: "rgba(15,5,15,0.65)",
    borderDefault: "#40303a", borderLight: "#4e3c48", borderFocus: "#d946ef",
    textPrimary: "#f0e4f0", textSecondary: "#d0b4d0", textMuted: "#a07898", textDisabled: "#402838",
    accent: "#d946ef", accentHover: "#c026d3", accentMuted: "rgba(217,70,239,0.10)",
    accentBorder: "rgba(217,70,239,0.22)", accentText: "#ffffff",
    shadowAccent: "0 4px 16px rgba(217,70,239,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  forest: darkVariant(darkTheme, {
    bgPrimary: "#131e16", bgSecondary: "#1a3224", bgTertiary: "#22402e",
    bgSidebar: "#102015", bgOverlay: "rgba(0,15,5,0.65)",
    borderDefault: "#243a2c", borderLight: "#2e4838", borderFocus: "#22c55e",
    textPrimary: "#e0f0e6", textSecondary: "#a8d4b8", textMuted: "#6fa888", textDisabled: "#2a4530",
    accent: "#22c55e", accentHover: "#16a34a", accentMuted: "rgba(34,197,94,0.10)",
    accentBorder: "rgba(34,197,94,0.22)", accentText: "#0a150e",
    shadowAccent: "0 4px 16px rgba(34,197,94,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "forest-mint": darkVariant(darkTheme, {
    bgPrimary: "#121e1a", bgSecondary: "#1a2e24", bgTertiary: "#223c30",
    bgSidebar: "#0f1a15", bgOverlay: "rgba(0,15,8,0.65)",
    borderDefault: "#1e3a28", borderLight: "#284834", borderFocus: "#34d399",
    textPrimary: "#e0f5ec", textSecondary: "#a0d8c0", textMuted: "#68b090", textDisabled: "#284538",
    accent: "#34d399", accentHover: "#10b981", accentMuted: "rgba(52,211,153,0.10)",
    accentBorder: "rgba(52,211,153,0.22)", accentText: "#0a1510",
    shadowAccent: "0 4px 16px rgba(52,211,153,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
  "forest-pine": darkVariant(darkTheme, {
    bgPrimary: "#101a14", bgSecondary: "#18281e", bgTertiary: "#203428",
    bgSidebar: "#0c1510", bgOverlay: "rgba(0,10,5,0.65)",
    borderDefault: "#1a3020", borderLight: "#243c2a", borderFocus: "#15803d",
    textPrimary: "#d8ede0", textSecondary: "#98c8a8", textMuted: "#608878", textDisabled: "#253828",
    accent: "#15803d", accentHover: "#166534", accentMuted: "rgba(21,128,61,0.10)",
    accentBorder: "rgba(21,128,61,0.22)", accentText: "#ffffff",
    shadowAccent: "0 4px 16px rgba(21,128,61,0.25)",
    warning: "#f59e0b", warningMuted: "rgba(245,158,11,0.10)", warningBorder: "rgba(245,158,11,0.25)",
  }),
};

export const DEFAULT_THEME: ThemeMode = "light";

export function getThemeFamily(mode: ThemeMode): string {
  if (mode.startsWith("dark")) return "dark";
  if (mode.startsWith("light")) return "light";
  if (mode.startsWith("ocean")) return "ocean";
  if (mode.startsWith("sunset")) return "sunset";
  if (mode.startsWith("forest")) return "forest";
  return "light";
}

export const THEME_DISPLAY_NAMES: Record<ThemeMode, string> = {
  dark: "Dark Classic",
  "dark-midnight": "Midnight Blue",
  "dark-purple": "Deep Purple",
  light: "Light Classic",
  "light-sky": "Sky Blue",
  "light-pastel": "Pastel Rose",
  ocean: "Ocean Deep Teal",
  "ocean-cyan": "Ocean Cyan",
  "ocean-pacific": "Ocean Pacific",
  sunset: "Sunset Amber",
  "sunset-rose": "Sunset Rose",
  "sunset-dusk": "Sunset Dusk",
  forest: "Forest Emerald",
  "forest-mint": "Forest Mint",
  "forest-pine": "Forest Pine",
};

export const THEME_FAMILIES = [
  { family: "dark", label: "Dark", modes: ["dark", "dark-midnight", "dark-purple"] as ThemeMode[] },
  { family: "light", label: "Light", modes: ["light", "light-sky", "light-pastel"] as ThemeMode[] },
  { family: "ocean", label: "Ocean", modes: ["ocean", "ocean-cyan", "ocean-pacific"] as ThemeMode[] },
  { family: "sunset", label: "Sunset", modes: ["sunset", "sunset-rose", "sunset-dusk"] as ThemeMode[] },
  { family: "forest", label: "Forest", modes: ["forest", "forest-mint", "forest-pine"] as ThemeMode[] },
];
