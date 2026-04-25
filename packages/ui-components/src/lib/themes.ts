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

// --- Shared semantic colors for dark-family themes ---
const DARK_SEMANTICS = {
  success: "#10b981",
  successMuted: "rgba(16,185,129,0.10)",
  successBorder: "rgba(16,185,129,0.25)",
  danger: "#ef4444",
  dangerMuted: "rgba(239,68,68,0.10)",
  dangerBorder: "rgba(239,68,68,0.25)",
  warning: "#f59e0b",
  warningMuted: "rgba(245,158,11,0.10)",
  warningBorder: "rgba(245,158,11,0.25)",
  info: "#3b82f6",
  infoMuted: "rgba(59,130,246,0.10)",
  infoBorder: "rgba(59,130,246,0.25)",
  purple: "#8b5cf6",
  purpleMuted: "rgba(139,92,246,0.10)",
  pink: "#ec4899",
  pinkMuted: "rgba(236,72,153,0.10)",
  teal: "#14b8a6",
  tealMuted: "rgba(20,184,166,0.10)",
  orange: "#f97316",
  orangeMuted: "rgba(249,115,22,0.10)",
} as const;

// --- Shared semantic colors for light-family themes ---
const LIGHT_SEMANTICS = {
  success: "#059669",
  successMuted: "rgba(5,150,105,0.08)",
  successBorder: "rgba(5,150,105,0.20)",
  danger: "#dc2626",
  dangerMuted: "rgba(220,38,38,0.07)",
  dangerBorder: "rgba(220,38,38,0.18)",
  warning: "#d97706",
  warningMuted: "rgba(217,119,6,0.08)",
  warningBorder: "rgba(217,119,6,0.18)",
  info: "#2563eb",
  infoMuted: "rgba(37,99,235,0.07)",
  infoBorder: "rgba(37,99,235,0.18)",
  purple: "#7c3aed",
  purpleMuted: "rgba(124,58,237,0.07)",
  pink: "#db2777",
  pinkMuted: "rgba(219,39,119,0.07)",
  teal: "#0d9488",
  tealMuted: "rgba(13,148,136,0.07)",
  orange: "#ea580c",
  orangeMuted: "rgba(234,88,12,0.07)",
} as const;

// --- Dark shadows ---
const DARK_SHADOWS = {
  shadowSm: "0 1px 3px rgba(0,0,0,0.4)",
  shadowMd: "0 4px 12px rgba(0,0,0,0.5)",
  shadowLg: "0 8px 32px rgba(0,0,0,0.6)",
} as const;

// --- Light shadows ---
const LIGHT_SHADOWS = {
  shadowSm: "0 1px 3px rgba(140,120,100,0.08)",
  shadowMd: "0 4px 12px rgba(140,120,100,0.10)",
  shadowLg: "0 8px 32px rgba(140,120,100,0.12)",
} as const;

// ===== THEME DEFINITIONS =====

const dark: ThemeColors = {
  bgPrimary: "#141519",
  bgSecondary: "#1c1e25",
  bgTertiary: "#25272f",
  bgSidebar: "#111318",
  bgOverlay: "rgba(0,0,0,0.65)",
  borderDefault: "#323540",
  borderLight: "#3d414c",
  borderFocus: "#f59e0b",
  textPrimary: "#e8e9ed",
  textSecondary: "#c8ccd4",
  textMuted: "#9ca0ab",
  textDisabled: "#3d4150",
  accent: "#f59e0b",
  accentHover: "#d97706",
  accentMuted: "rgba(245,158,11,0.10)",
  accentBorder: "rgba(245,158,11,0.22)",
  accentText: "#0a0b0e",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(245,158,11,0.25)",
};

const darkMidnight: ThemeColors = {
  bgPrimary: "#0f1420",
  bgSecondary: "#171d2c",
  bgTertiary: "#1e2538",
  bgSidebar: "#0c1018",
  bgOverlay: "rgba(0,0,0,0.70)",
  borderDefault: "#283040",
  borderLight: "#34405a",
  borderFocus: "#3b82f6",
  textPrimary: "#e0e4ef",
  textSecondary: "#b0b8cc",
  textMuted: "#7a86a0",
  textDisabled: "#384060",
  accent: "#3b82f6",
  accentHover: "#2563eb",
  accentMuted: "rgba(59,130,246,0.10)",
  accentBorder: "rgba(59,130,246,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(59,130,246,0.25)",
};

const darkPurple: ThemeColors = {
  bgPrimary: "#15131e",
  bgSecondary: "#1e1b2a",
  bgTertiary: "#272336",
  bgSidebar: "#110f18",
  bgOverlay: "rgba(0,0,0,0.68)",
  borderDefault: "#302a45",
  borderLight: "#3d3558",
  borderFocus: "#8b5cf6",
  textPrimary: "#e6e2f0",
  textSecondary: "#c0b8d4",
  textMuted: "#8a80a5",
  textDisabled: "#3a3350",
  accent: "#8b5cf6",
  accentHover: "#7c3aed",
  accentMuted: "rgba(139,92,246,0.10)",
  accentBorder: "rgba(139,92,246,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(139,92,246,0.25)",
};

const light: ThemeColors = {
  bgPrimary: "#faf8f5",
  bgSecondary: "#ffffff",
  bgTertiary: "#f5f3f0",
  bgSidebar: "#f0ebe3",
  bgOverlay: "rgba(90,70,50,0.25)",
  borderDefault: "#e0d8cc",
  borderLight: "#eae3d8",
  borderFocus: "#d97706",
  textPrimary: "#1a1612",
  textSecondary: "#4a4239",
  textMuted: "#6b6358",
  textDisabled: "#c4bdb4",
  accent: "#d97706",
  accentHover: "#b45309",
  accentMuted: "rgba(217,119,6,0.08)",
  accentBorder: "rgba(217,119,6,0.20)",
  accentText: "#ffffff",
  ...LIGHT_SEMANTICS,
  ...LIGHT_SHADOWS,
  shadowAccent: "0 4px 16px rgba(217,119,6,0.18)",
};

const lightSky: ThemeColors = {
  bgPrimary: "#f5f9ff",
  bgSecondary: "#ffffff",
  bgTertiary: "#edf3fc",
  bgSidebar: "#e8f0fa",
  bgOverlay: "rgba(40,60,100,0.20)",
  borderDefault: "#d0ddef",
  borderLight: "#e0eaf5",
  borderFocus: "#0284c7",
  textPrimary: "#0f172a",
  textSecondary: "#334155",
  textMuted: "#64748b",
  textDisabled: "#c0cde0",
  accent: "#0284c7",
  accentHover: "#0369a1",
  accentMuted: "rgba(2,132,199,0.08)",
  accentBorder: "rgba(2,132,199,0.20)",
  accentText: "#ffffff",
  ...LIGHT_SEMANTICS,
  shadowSm: "0 1px 3px rgba(100,120,160,0.08)",
  shadowMd: "0 4px 12px rgba(100,120,160,0.10)",
  shadowLg: "0 8px 32px rgba(100,120,160,0.12)",
  shadowAccent: "0 4px 16px rgba(2,132,199,0.18)",
};

const lightPastel: ThemeColors = {
  bgPrimary: "#fdf5f5",
  bgSecondary: "#ffffff",
  bgTertiary: "#faf0f0",
  bgSidebar: "#f5e8e8",
  bgOverlay: "rgba(100,50,60,0.20)",
  borderDefault: "#e8d0d4",
  borderLight: "#f0e0e2",
  borderFocus: "#e11d48",
  textPrimary: "#1a1215",
  textSecondary: "#4a3540",
  textMuted: "#6b5560",
  textDisabled: "#c4b0b5",
  accent: "#e11d48",
  accentHover: "#be123c",
  accentMuted: "rgba(225,29,72,0.08)",
  accentBorder: "rgba(225,29,72,0.20)",
  accentText: "#ffffff",
  ...LIGHT_SEMANTICS,
  shadowSm: "0 1px 3px rgba(140,100,110,0.08)",
  shadowMd: "0 4px 12px rgba(140,100,110,0.10)",
  shadowLg: "0 8px 32px rgba(140,100,110,0.12)",
  shadowAccent: "0 4px 16px rgba(225,29,72,0.18)",
};

const ocean: ThemeColors = {
  bgPrimary: "#152232",
  bgSecondary: "#1c2a40",
  bgTertiary: "#243450",
  bgSidebar: "#111e30",
  bgOverlay: "rgba(0,10,30,0.70)",
  borderDefault: "#2a3a55",
  borderLight: "#354868",
  borderFocus: "#38bdf8",
  textPrimary: "#e0eaf5",
  textSecondary: "#a8bcd0",
  textMuted: "#708aa8",
  textDisabled: "#304060",
  accent: "#38bdf8",
  accentHover: "#0ea5e9",
  accentMuted: "rgba(56,189,248,0.10)",
  accentBorder: "rgba(56,189,248,0.22)",
  accentText: "#0a1520",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(56,189,248,0.25)",
};

const oceanCyan: ThemeColors = {
  bgPrimary: "#0f2028",
  bgSecondary: "#162a35",
  bgTertiary: "#1e3545",
  bgSidebar: "#0b1a22",
  bgOverlay: "rgba(0,10,20,0.72)",
  borderDefault: "#253a48",
  borderLight: "#30485a",
  borderFocus: "#06b6d4",
  textPrimary: "#daf0f5",
  textSecondary: "#a0c8d8",
  textMuted: "#6898b0",
  textDisabled: "#2a4050",
  accent: "#06b6d4",
  accentHover: "#0891b2",
  accentMuted: "rgba(6,182,212,0.10)",
  accentBorder: "rgba(6,182,212,0.22)",
  accentText: "#0a1520",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(6,182,212,0.25)",
};

const oceanPacific: ThemeColors = {
  bgPrimary: "#101828",
  bgSecondary: "#182238",
  bgTertiary: "#202e4a",
  bgSidebar: "#0c1420",
  bgOverlay: "rgba(0,5,25,0.72)",
  borderDefault: "#283858",
  borderLight: "#354570",
  borderFocus: "#6366f1",
  textPrimary: "#e0e5f5",
  textSecondary: "#a8b2d0",
  textMuted: "#7080a8",
  textDisabled: "#303e60",
  accent: "#6366f1",
  accentHover: "#4f46e5",
  accentMuted: "rgba(99,102,241,0.10)",
  accentBorder: "rgba(99,102,241,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(99,102,241,0.25)",
};

const sunset: ThemeColors = {
  bgPrimary: "#1e1510",
  bgSecondary: "#2a1f18",
  bgTertiary: "#352820",
  bgSidebar: "#211510",
  bgOverlay: "rgba(10,5,0,0.68)",
  borderDefault: "#40302a",
  borderLight: "#503e35",
  borderFocus: "#f97316",
  textPrimary: "#f0e4d8",
  textSecondary: "#d0b8a0",
  textMuted: "#a08060",
  textDisabled: "#4a3828",
  accent: "#f97316",
  accentHover: "#ea580c",
  accentMuted: "rgba(249,115,22,0.10)",
  accentBorder: "rgba(249,115,22,0.22)",
  accentText: "#0e0a05",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(249,115,22,0.25)",
};

const sunsetRose: ThemeColors = {
  bgPrimary: "#1e1215",
  bgSecondary: "#2a1c20",
  bgTertiary: "#35242a",
  bgSidebar: "#1a1012",
  bgOverlay: "rgba(10,0,5,0.68)",
  borderDefault: "#402830",
  borderLight: "#503540",
  borderFocus: "#f43f5e",
  textPrimary: "#f0e0e4",
  textSecondary: "#d0a8b0",
  textMuted: "#a07080",
  textDisabled: "#4a2838",
  accent: "#f43f5e",
  accentHover: "#e11d48",
  accentMuted: "rgba(244,63,94,0.10)",
  accentBorder: "rgba(244,63,94,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(244,63,94,0.25)",
};

const sunsetDusk: ThemeColors = {
  bgPrimary: "#1a1318",
  bgSecondary: "#251e24",
  bgTertiary: "#302830",
  bgSidebar: "#161015",
  bgOverlay: "rgba(8,0,10,0.68)",
  borderDefault: "#3a2838",
  borderLight: "#4a3548",
  borderFocus: "#d946ef",
  textPrimary: "#f0e0f0",
  textSecondary: "#c8a8c8",
  textMuted: "#987098",
  textDisabled: "#402840",
  accent: "#d946ef",
  accentHover: "#c026d3",
  accentMuted: "rgba(217,70,239,0.10)",
  accentBorder: "rgba(217,70,239,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(217,70,239,0.25)",
};

const forest: ThemeColors = {
  bgPrimary: "#131e16",
  bgSecondary: "#1a3224",
  bgTertiary: "#22402e",
  bgSidebar: "#102015",
  bgOverlay: "rgba(0,10,5,0.70)",
  borderDefault: "#284038",
  borderLight: "#355048",
  borderFocus: "#22c55e",
  textPrimary: "#e0f0e4",
  textSecondary: "#a8d0b0",
  textMuted: "#70a880",
  textDisabled: "#284838",
  accent: "#22c55e",
  accentHover: "#16a34a",
  accentMuted: "rgba(34,197,94,0.10)",
  accentBorder: "rgba(34,197,94,0.22)",
  accentText: "#0a1510",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(34,197,94,0.25)",
};

const forestMint: ThemeColors = {
  bgPrimary: "#121e1a",
  bgSecondary: "#1a2e24",
  bgTertiary: "#223e30",
  bgSidebar: "#0f1a15",
  bgOverlay: "rgba(0,10,8,0.70)",
  borderDefault: "#254035",
  borderLight: "#305045",
  borderFocus: "#34d399",
  textPrimary: "#ddf0ea",
  textSecondary: "#a0d0c0",
  textMuted: "#68a890",
  textDisabled: "#254535",
  accent: "#34d399",
  accentHover: "#10b981",
  accentMuted: "rgba(52,211,153,0.10)",
  accentBorder: "rgba(52,211,153,0.22)",
  accentText: "#0a1510",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(52,211,153,0.25)",
};

const forestPine: ThemeColors = {
  bgPrimary: "#101a14",
  bgSecondary: "#18281e",
  bgTertiary: "#20382a",
  bgSidebar: "#0c1510",
  bgOverlay: "rgba(0,8,4,0.72)",
  borderDefault: "#223828",
  borderLight: "#2d4835",
  borderFocus: "#15803d",
  textPrimary: "#d8f0e0",
  textSecondary: "#a0c8a8",
  textMuted: "#689878",
  textDisabled: "#224030",
  accent: "#15803d",
  accentHover: "#166534",
  accentMuted: "rgba(21,128,61,0.10)",
  accentBorder: "rgba(21,128,61,0.22)",
  accentText: "#ffffff",
  ...DARK_SEMANTICS,
  ...DARK_SHADOWS,
  shadowAccent: "0 4px 16px rgba(21,128,61,0.25)",
};

export const THEME_MAP: Record<ThemeMode, ThemeColors> = {
  "dark": dark,
  "dark-midnight": darkMidnight,
  "dark-purple": darkPurple,
  "light": light,
  "light-sky": lightSky,
  "light-pastel": lightPastel,
  "ocean": ocean,
  "ocean-cyan": oceanCyan,
  "ocean-pacific": oceanPacific,
  "sunset": sunset,
  "sunset-rose": sunsetRose,
  "sunset-dusk": sunsetDusk,
  "forest": forest,
  "forest-mint": forestMint,
  "forest-pine": forestPine,
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
