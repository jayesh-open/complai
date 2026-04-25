# Aura — Design System & UI Guidelines
## For Claude Code: Follow These Rules for EVERY UI Component

**This file defines the visual standard for Aura. Every page, component, and interaction must follow these rules. Do NOT deviate. Do NOT use generic/default shadcn styling. Apply these theme tokens and component patterns everywhere.**

---

## 1. THEME SYSTEM — 15 Variants across 5 Families

Aura offers 5 theme families with 3 variants each (15 total): **Dark** (Classic, Midnight Blue, Deep Purple), **Light** (Classic, Sky Blue, Pastel), **Ocean** (Deep Teal, Cyan Glow, Pacific), **Sunset** (Amber Gold, Rose, Dusk), **Forest** (Emerald, Mint, Deep Pine). Each variant uses a unique gradient for sidebar and content backgrounds. All variants verified for text readability — WCAG AA minimum contrast on all text, badges, charts, checkboxes, filters. The active theme is stored in Zustand (`appStore.theme`) and persisted in localStorage. All colors are applied via CSS custom properties on `<html>` element. ThemeProvider also sets `data-theme-family` attribute (dark/light/ocean/sunset/forest) for CSS overrides that apply to all variants in a family.

### Theme Color Interface

```tsx
// lib/themes.ts

export type ThemeMode =
  | "dark" | "dark-midnight" | "dark-purple"
  | "light" | "light-sky" | "light-pastel"
  | "ocean" | "ocean-cyan" | "ocean-pacific"
  | "sunset" | "sunset-rose" | "sunset-dusk"
  | "forest" | "forest-mint" | "forest-pine";

export interface ThemeColors {
  // Backgrounds
  bgPrimary: string;        // Main page background
  bgSecondary: string;      // Cards, panels, elevated surfaces
  bgTertiary: string;       // Inputs, nested elements, table rows on hover
  bgSidebar: string;        // Sidebar background
  bgOverlay: string;        // Modal/dialog overlay

  // Borders
  borderDefault: string;    // Default borders
  borderLight: string;      // Subtle dividers
  borderFocus: string;      // Input focus ring

  // Text
  textPrimary: string;      // Headings, primary content
  textSecondary: string;    // Body text, descriptions
  textMuted: string;        // Labels, placeholders, timestamps
  textDisabled: string;     // Disabled state text

  // Accent (Brand)
  accent: string;           // Primary brand color — amber
  accentHover: string;      // Accent button hover
  accentMuted: string;      // Accent at 10-15% opacity for backgrounds
  accentBorder: string;     // Accent at 20-25% opacity for borders
  accentText: string;       // Text color ON accent backgrounds

  // Semantic Colors
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

  // Extended Palette
  purple: string;
  purpleMuted: string;
  pink: string;
  pinkMuted: string;
  teal: string;
  tealMuted: string;
  orange: string;
  orangeMuted: string;

  // Shadows
  shadowSm: string;
  shadowMd: string;
  shadowLg: string;
  shadowAccent: string;     // Glow shadow for accent buttons
}
```

### Dark Theme — Premium Dark Fintech
Inspired by: Linear, Raycast, Vercel

```tsx
export const darkTheme: ThemeColors = {
  bgPrimary:    "#141519",
  bgSecondary:  "#1c1e25",
  bgTertiary:   "#25272f",
  bgSidebar:    "#111318",
  bgOverlay:    "rgba(0, 0, 0, 0.65)",

  borderDefault: "#323540",
  borderLight:   "#3d414c",
  borderFocus:   "#f59e0b",

  textPrimary:   "#e8e9ed",
  textSecondary: "#c8ccd4",
  textMuted:     "#9ca0ab",
  textDisabled:  "#3d4150",

  accent:        "#f59e0b",
  accentHover:   "#d97706",
  accentMuted:   "rgba(245, 158, 11, 0.10)",
  accentBorder:  "rgba(245, 158, 11, 0.22)",
  accentText:    "#0a0b0e",

  success:       "#10b981",
  successMuted:  "rgba(16, 185, 129, 0.10)",
  successBorder: "rgba(16, 185, 129, 0.25)",
  danger:        "#ef4444",
  dangerMuted:   "rgba(239, 68, 68, 0.10)",
  dangerBorder:  "rgba(239, 68, 68, 0.25)",
  warning:       "#f59e0b",
  warningMuted:  "rgba(245, 158, 11, 0.10)",
  warningBorder: "rgba(245, 158, 11, 0.25)",
  info:          "#3b82f6",
  infoMuted:     "rgba(59, 130, 246, 0.10)",
  infoBorder:    "rgba(59, 130, 246, 0.25)",

  purple:       "#8b5cf6",
  purpleMuted:  "rgba(139, 92, 246, 0.10)",
  pink:         "#ec4899",
  pinkMuted:    "rgba(236, 72, 153, 0.10)",
  teal:         "#14b8a6",
  tealMuted:    "rgba(20, 184, 166, 0.10)",
  orange:       "#f97316",
  orangeMuted:  "rgba(249, 115, 22, 0.10)",

  shadowSm:     "0 1px 3px rgba(0,0,0,0.4)",
  shadowMd:     "0 4px 12px rgba(0,0,0,0.5)",
  shadowLg:     "0 8px 32px rgba(0,0,0,0.6)",
  shadowAccent: "0 4px 16px rgba(245, 158, 11, 0.25)",
};
```

#### Dark Theme Component Rules
- Primary text: #e8e9ed. Secondary: #c8ccd4. Muted: #9ca0ab minimum.
- Hover text: #ffffff (pure white). Hover bg: rgba(255,255,255,0.10).
- NEVER dark text on dark bg on hover — hover must INCREASE contrast.
- Modals/dialogs/sheets: MUST have visible 1px border AND card background (noticeably lighter than page bg).
- Input fields: visible 1px border + slightly lighter bg than card bg.
- Dropdown hover: contrasting bg with light text, NEVER dark text on dark bg.
- Muted text minimum contrast: #9ca0ab — NEVER darker.

### Light Theme — Clean, Warm, Professional
Inspired by: Notion, Linear Light, Ramp. Background: very light warm, similar to Claude desktop.

```tsx
export const lightTheme: ThemeColors = {
  bgPrimary:    "#faf8f5",
  bgSecondary:  "#ffffff",
  bgTertiary:   "#f5f3f0",
  bgSidebar:    "#f0ebe3",
  bgOverlay:    "rgba(90, 70, 50, 0.25)",

  borderDefault: "#e0d8cc",
  borderLight:   "#eae3d8",
  borderFocus:   "#d97706",

  textPrimary:   "#1a1612",
  textSecondary: "#4a4239",
  textMuted:     "#6b6358",
  textDisabled:  "#c4bdb4",

  accent:        "#d97706",
  accentHover:   "#b45309",
  accentMuted:   "rgba(217, 119, 6, 0.08)",
  accentBorder:  "rgba(217, 119, 6, 0.20)",
  accentText:    "#ffffff",

  success:       "#059669",
  successMuted:  "rgba(5, 150, 105, 0.08)",
  successBorder: "rgba(5, 150, 105, 0.20)",
  danger:        "#dc2626",
  dangerMuted:   "rgba(220, 38, 38, 0.07)",
  dangerBorder:  "rgba(220, 38, 38, 0.18)",
  warning:       "#d97706",
  warningMuted:  "rgba(217, 119, 6, 0.08)",
  warningBorder: "rgba(217, 119, 6, 0.18)",
  info:          "#2563eb",
  infoMuted:     "rgba(37, 99, 235, 0.07)",
  infoBorder:    "rgba(37, 99, 235, 0.18)",

  purple:       "#7c3aed",
  purpleMuted:  "rgba(124, 58, 237, 0.07)",
  pink:         "#db2777",
  pinkMuted:    "rgba(219, 39, 119, 0.07)",
  teal:         "#0d9488",
  tealMuted:    "rgba(13, 148, 136, 0.07)",
  orange:       "#ea580c",
  orangeMuted:  "rgba(234, 88, 12, 0.07)",

  shadowSm:     "0 1px 3px rgba(140,120,100,0.08)",
  shadowMd:     "0 4px 12px rgba(140,120,100,0.10)",
  shadowLg:     "0 8px 32px rgba(140,120,100,0.12)",
  shadowAccent: "0 4px 16px rgba(217, 119, 6, 0.18)",
};
```

#### Light Theme Component Rules
- Primary text: #1a1612. Secondary: #4a4239. Muted: #6b6358.
- Cards: must have visible border OR shadow to distinguish from background.
- Input fields: #f5f3f0 background with visible border.
- Buttons and interactive elements have clear borders/shadows.

### Other Theme Families (summary)

**Ocean Theme:** Sidebar: linear-gradient(180deg, #111e30, #122840, #134270). Content: #152232. Cards: #1c2a40. Accent: #38bdf8 (sky blue). Text: #e8edf3 / #b0c4de / #7a9bbd. Border: #243a5a.

**Sunset Theme:** Sidebar: linear-gradient(180deg, #211510, #2d1810, #45261a). Content: #1e1510. Cards: #2a1f18. Accent: #f97316 (orange). Text: #f0e6e0 / #d4b8aa / #a08070. Border: #452e24.

**Forest Theme:** Sidebar: linear-gradient(180deg, #102015, #152918, #1a4028). Content: #131e16. Cards: #1a3224. Accent: #22c55e (green). Text: #e0f0e6 / #a8d4b8 / #6fa888. Border: #243a2c.

### Applying Theme — CSS Variables

```tsx
// lib/theme-provider.tsx
"use client";
import { useEffect } from "react";
import { useAppStore } from "@/store/app-store";
import { THEME_MAP, type ThemeColors } from "@/lib/themes";

export function ThemeProvider({ children }: { children: React.ReactNode }) {
  const { theme } = useAppStore();

  useEffect(() => {
    const colors: ThemeColors = THEME_MAP[theme] || THEME_MAP.dark;
    const root = document.documentElement;
    Object.entries(colors).forEach(([key, value]) => {
      const cssVar = `--${key.replace(/([A-Z])/g, "-$1").toLowerCase()}`;
      root.style.setProperty(cssVar, value);
    });
    root.setAttribute("data-theme", theme);
  }, [theme]);

  return <>{children}</>;
}
```

### Tailwind CSS Integration

```ts
// tailwind.config.ts
module.exports = {
  theme: {
    extend: {
      colors: {
        app: {
          bg:          "var(--bg-primary)",
          card:        "var(--bg-secondary)",
          input:       "var(--bg-tertiary)",
          sidebar:     "var(--bg-sidebar)",
          border:      "var(--border-default)",
          "border-lt": "var(--border-light)",
          accent:      "var(--accent)",
          "accent-h":  "var(--accent-hover)",
          "accent-m":  "var(--accent-muted)",
          "accent-b":  "var(--accent-border)",
          "accent-t":  "var(--accent-text)",
          success:     "var(--success)",
          "success-m": "var(--success-muted)",
          danger:      "var(--danger)",
          "danger-m":  "var(--danger-muted)",
          warning:     "var(--warning)",
          "warning-m": "var(--warning-muted)",
          info:        "var(--info)",
          "info-m":    "var(--info-muted)",
          purple:      "var(--purple)",
          "purple-m":  "var(--purple-muted)",
          teal:        "var(--teal)",
          "teal-m":    "var(--teal-muted)",
        },
        foreground: {
          DEFAULT:   "var(--text-primary)",
          secondary: "var(--text-secondary)",
          muted:     "var(--text-muted)",
          disabled:  "var(--text-disabled)",
        },
      },
      boxShadow: {
        "app-sm":     "var(--shadow-sm)",
        "app-md":     "var(--shadow-md)",
        "app-lg":     "var(--shadow-lg)",
        "app-accent": "var(--shadow-accent)",
      },
    },
  },
};

// USAGE IN COMPONENTS:
// <div className="bg-app-card border border-app-border rounded-2xl p-6">
// <h2 className="text-foreground font-bold">Title</h2>
// <p className="text-foreground-muted">Description</p>
// <button className="bg-app-accent text-app-accent-t rounded-lg px-5 py-2">Action</button>
```

---

## 2. TYPOGRAPHY

### Font Stack
```css
/* Primary: System font stack for maximum performance and native feel */
--font-primary: -apple-system, BlinkMacSystemFont, "Segoe UI", "Noto Sans",
                Helvetica, Arial, sans-serif, "Apple Color Emoji", "Segoe UI Emoji";

/* Monospace: For IDs, codes, reference numbers, technical values */
--font-mono: "SF Mono", "Cascadia Code", "Fira Code", "JetBrains Mono",
             ui-monospace, SFMono-Regular, Menlo, Consolas, monospace;
```

DO NOT import Google Fonts or any external font. System fonts load instantly and feel native.

### Type Scale

| Token | Size | Weight | Line Height | Usage |
|-------|------|--------|-------------|-------|
| `display` | 24px | 800 | 1.15 | Page main metric (key numbers on dashboard) |
| `heading-xl` | 18px | 700 | 1.2 | Page greeting |
| `heading-lg` | 16px | 700 | 1.25 | Section titles, card titles |
| `heading-md` | 14px | 700 | 1.3 | Page titles in header, dialog titles |
| `heading-sm` | 13px | 600 | 1.35 | Card headers, table group headers |
| `body` | 13px | 400 | 1.5 | Default body text, descriptions |
| `body-sm` | 12px | 400/500 | 1.45 | Table cells, form labels, secondary info |
| `caption` | 11px | 500 | 1.4 | Filter labels, timestamps, helper text |
| `overline` | 10px | 600 | 1.3 | Metric labels, badge text, uppercase labels |
| `tiny` | 9px | 600 | 1.2 | Keyboard shortcuts, version numbers |

### Typography Rules
- Page titles: `heading-md` (14px/700) — NOT larger. Keep the UI compact and data-dense.
- Metric values: `display` (24px/800) with `font-feature-settings: "tnum"` for tabular numbers
- IDs, codes, reference numbers: ALWAYS use `font-mono` class
- All uppercase labels: add `letter-spacing: 0.04em`
- Never use font-size larger than 24px anywhere except the main dashboard metric
- Body text: use `text-secondary` color, not `text-primary`. Reserve primary for headings and values.

### Tailwind Classes
```
text-display:    text-[24px] font-extrabold leading-[1.15]
text-heading-xl: text-[18px] font-bold leading-[1.2]
text-heading-lg: text-base font-bold leading-[1.25]
text-heading-md: text-sm font-bold leading-[1.3]
text-heading-sm: text-[13px] font-semibold leading-[1.35]
text-body:       text-[13px] font-normal leading-[1.5]
text-body-sm:    text-xs font-normal leading-[1.45]
text-caption:    text-[11px] font-medium leading-[1.4]
text-overline:   text-[10px] font-semibold leading-[1.3] tracking-wide uppercase
text-tiny:       text-[9px] font-semibold leading-[1.2]
```

---

## 3. SPACING SYSTEM

Use a consistent 4px base grid. All spacing must be multiples of 4.

| Token | Value | Usage |
|-------|-------|-------|
| `space-1` | 4px | Icon-to-text gap inside badges |
| `space-2` | 8px | Between inline elements, icon padding |
| `space-3` | 12px | Between list items, small card padding |
| `space-4` | 16px | Standard gap between cards, table cell padding |
| `space-5` | 20px | Card internal padding, section gaps |
| `space-6` | 24px | Large card padding, page section gaps |
| `space-7` | 28px | Page padding (main content area) |
| `space-8` | 32px | Between major page sections |
| `space-10` | 40px | Maximum internal spacing |

### Spacing Rules
- Page padding: `28px` (7 × 4) on all sides of the main content area
- Card padding: `22px` for standard cards, `24px` for large feature cards
- Card gap (grid/flex gap between cards): `16px`
- Table cell padding: `14px 18px` (vertical × horizontal)
- Form field gap: `16px` between fields, `24px` between field groups
- Sidebar item padding: `9px 14px`
- Button padding: `9px 20px` (standard), `6px 14px` (small), `12px 28px` (large)
- Badge padding: `2px 10px` (standard), `1px 6px` (small/xs)
- Border radius: `16px` (cards), `10px` (buttons, inputs, badges), `8px` (small elements), `6px` (tiny elements like chart bars)
- Modal width: `480px` (small), `640px` (medium), `800px` (large), `960px` (full form)

---

## 4. COMPONENT SPECIFICATIONS

### 4.1 Metric Card
```
┌─────────────────────────────────┐
│  [icon 38x38]        ↑ 12%     │   ← Trend top-right
│                                 │
│  1,284                          │   ← display size, font-extrabold
│  Active Users                   │   ← caption size, text-muted
│  +42 this week                  │   ← tiny size, text-disabled
└─────────────────────────────────┘
```
- Background: `bg-app-card`
- Border: `1px solid border-default`
- Border radius: `16px`
- Padding: `16px`
- Icon container: `30px × 30px`, rounded-[8px], background: `{color}15`, border: `1px solid {color}30`
- Amount: 20px bold, label: 10px uppercase muted with letter-spacing 0.5px
- Thin muted separator line with my-2 margin
- Subtitle/count: 12px muted below separator
- Hover: border brightens to accent
- Each metric card uses a different icon color (accent, info, success, danger)

### 4.2 Data Table
```
┌──────────────────────────────────────────────────────────────────────────┐
│  Name             │ Category     │ Status      │ Value      │ Actions   │  ← Header
├──────────────────────────────────────────────────────────────────────────┤
│  Item One         │ Type A       │ ● Active    │ 1,234      │ [Edit]    │  ← Data row
│  REF-001234       │              │             │            │           │  ← Sub-line
├──────────────────────────────────────────────────────────────────────────┤
│  Showing 1-20 of 342                              ← 1  2  3 ... 35  →  │  ← Footer
└──────────────────────────────────────────────────────────────────────────┘
```
- Container: `bg-app-card`, `rounded-[14px]`, `border border-app-border`, `overflow-hidden`
- Header row: `px-[18px] py-[10px]`, `border-b`, `text-overline` style, `text-foreground-muted`
- Data row: `px-[18px] py-[14px]`, `border-b border-app-border`
- Row hover: `bg-app-input` transition
- Row cursor: `pointer` (clickable rows navigate to detail)
- Primary cell text: `text-body-sm font-semibold text-foreground`
- Secondary cell text: `text-body-sm text-foreground-muted`
- Sub-line: `text-[11px] text-foreground-disabled font-mono`
- Pagination: `30px × 30px` buttons, `rounded-[6px]`, active: `bg-app-accent-m text-app-accent`
- IDs and reference numbers in cells: ALWAYS `font-mono`
- Numeric values: right-aligned, `font-semibold`, tabular numbers

#### Column Width Rules
| Column Type | Sizing | CSS | Examples |
|---|---|---|---|
| Names / descriptions | `minmax(Npx, fr)` — expands to fill | `truncate` | Name, Description |
| IDs / codes | `minmax(Npx, fr)` or fixed | `font-mono text-[12px] truncate` | ID, Reference # |
| Numbers | Fixed 80-100px | `text-right tabular-nums whitespace-nowrap` | Amounts, counts |
| Dates | Fixed 80-90px | `whitespace-nowrap text-[12px]` | DD/MM/YYYY |
| Status badges | Fixed 80-110px | `whitespace-nowrap text-center` | Active, Pending |
| Action buttons | Fixed by button count | `flex gap-1 justify-center` | Edit, Delete |

### 4.3 Status Badge
```
Dimensions: padding 2px 10px, border-radius 6px, font-size 10px, font-weight 600
Letter-spacing: 0.02em, text-transform: uppercase

Variants:
├── success:  bg=successMuted, text=success, border=successBorder     → Active, Verified, Complete
├── danger:   bg=dangerMuted, text=danger, border=dangerBorder        → Failed, Error, Blocked
├── warning:  bg=warningMuted, text=warning, border=warningBorder     → Pending, On Hold, Expiring
├── info:     bg=infoMuted, text=info, border=infoBorder              → Processing, Under Review
├── purple:   bg=purpleMuted, text=purple, border=purpleBorder        → Custom, Special
├── teal:     bg=tealMuted, text=teal, border=tealBorder              → Category-specific
├── default:  bg=rgba(107,112,128,0.1), text=textMuted                → Draft, Informational

XS variant: padding 1px 6px, font-size 9px (for inline badges in table cells)

ALWAYS use a 1px border with the *Border color. Never use a solid-fill badge.
```

### 4.4 Buttons

```
Primary (CTA):
  background: linear-gradient(135deg, accent, accentHover)
  color: accentText
  border: none
  border-radius: 10px
  padding: 8px 18px
  font-size: 12px, font-weight: 700
  box-shadow: shadowAccent
  hover: brightness(0.9)
  Only ONE primary button per visible screen area

  ⚠️ TEXT COLOR RULE (MANDATORY):
  ALWAYS add className="app-gradient-btn" to gradient/CTA buttons.
  NEVER use text-white, text-black, or inline color styles.
  The .app-gradient-btn class applies: color: var(--accent-text) !important
  This ensures correct contrast across ALL 15 theme variants automatically.

Secondary:
  background: transparent
  color: textSecondary
  border: 1px solid borderDefault
  border-radius: 10px
  padding: 8px 18px
  font-size: 12px, font-weight: 500
  hover: bg-app-input

Danger:
  background: dangerMuted
  color: danger
  border: none
  border-radius: 10px
  padding: 8px 18px
  font-size: 12px, font-weight: 600

Ghost:
  background: transparent
  color: textMuted
  border: none
  padding: 6px 12px
  hover: bg-app-input

Icon Button:
  width: 32px, height: 32px
  border-radius: 8px
  border: 1px solid borderDefault
  hover: bg-app-input

Small Button:
  padding: 5px 12px, font-size: 11px
```

### 4.4.1 Selected / Active Button States

All interactive elements that can be "selected" or "active" follow a consistent pattern across all 15 themes using CSS variables. The pattern is: accent-tinted background + accent text + accent border.

```
SELECTED STATE PATTERN (universal across all themes):

  ┌─ Unselected ────────────────────────────────────────────────┐
  │  background: transparent                                     │
  │  color: textMuted                                            │
  │  border: 1px solid borderDefault                             │
  │  font-weight: 450                                            │
  └──────────────────────────────────────────────────────────────┘

  ┌─ Selected / Active ─────────────────────────────────────────┐
  │  background: accentMuted (accent at 10% opacity)             │
  │  color: accent                                               │
  │  border: 1px solid accentBorder (accent at 22% opacity)      │
  │  font-weight: 600                                            │
  └──────────────────────────────────────────────────────────────┘

How this renders per theme family:

  Dark Classic:
    Unselected: transparent bg, #9ca0ab text, #2a2d38 border
    Selected:   rgba(245,158,11,0.10) bg, #f59e0b text, rgba(245,158,11,0.22) border

  Light Classic:
    Unselected: transparent bg, #6b6358 text, #e0d8cc border
    Selected:   rgba(217,119,6,0.08) bg, #d97706 text, rgba(217,119,6,0.20) border

  Ocean:
    Unselected: transparent bg, #7a9bbd text, #1e3a5f border
    Selected:   rgba(56,189,248,0.10) bg, #38bdf8 text, rgba(56,189,248,0.22) border

  Sunset:
    Unselected: transparent bg, #a08070 text, #3d2820 border
    Selected:   rgba(249,115,22,0.10) bg, #f97316 text, rgba(249,115,22,0.22) border

  Forest:
    Unselected: transparent bg, #6fa888 text, #1e3a28 border
    Selected:   rgba(34,197,94,0.10) bg, #22c55e text, rgba(34,197,94,0.22) border
```

**Where the selected state applies:**

```
Filter Pills / Chip Buttons:
  padding: 8px 14px, rounded-[8px]
  Unselected: bg-app-input, border border-app-border, text-foreground-muted
  Selected:   bg-accentMuted, border border-accentBorder, text-accent, font-weight 600
  Transition: background 150ms ease, color 150ms ease, border-color 150ms ease

Tab Buttons:
  Unselected: bg transparent, text-foreground-muted, font-weight 450, no bottom border
  Selected:   bg-app-card, text-foreground, font-weight 600,
              border-bottom: 2px solid accent (overlapping container border)

Toggle Buttons (on/off groups):
  Unselected: bg transparent, text-foreground-muted, border border-app-border
  Selected:   bg-accentMuted, text-accent, border border-accentBorder
  Group container: flex gap-0, first child rounded-l-[8px], last child rounded-r-[8px]
  Divider between: 1px solid borderDefault (hidden when adjacent button is selected)

Sidebar Nav Items:
  Unselected: bg transparent, text-foreground-secondary
  Selected:   bg-accentMuted, text-accent, 3px solid left border in accent color
  Hover (unselected): bg-app-input/40, text-foreground

Pagination Buttons:
  Unselected: bg transparent, text-foreground-muted
  Selected:   bg-accentMuted, text-accent, font-weight 600
  Size: 30px × 30px, rounded-[6px]

Dropdown Menu Selected Option:
  Unselected: bg transparent, text-foreground-secondary
  Selected:   bg-accentMuted, text-accent
  Hover: bg-app-input

Card Selection (clickable cards):
  Unselected: bg-app-card, border border-app-border
  Selected:   bg-app-card, border-2 border-accent, shadow-app-accent (subtle glow)
  Check icon: 16px accent-colored checkmark in top-right corner

Checkbox / Radio (custom styled):
  Unchecked: border-2 borderDefault, bg transparent
  Checked:   border-2 accent, bg accent, white checkmark icon
  Indeterminate: border-2 accent, bg accent, white dash icon

Theme Selector Cards (Settings → Appearance):
  Unselected: bg-app-card, border border-app-border, rounded-[12px]
  Selected:   bg-app-card, border-2 border-accent, check icon overlay
```

**Tailwind utility classes for selected states:**

```css
/* Reusable selected state classes */
.selected-pill {
  @apply bg-[var(--accent-muted)] text-[var(--accent)] 
         border border-[var(--accent-border)] font-semibold;
}

.unselected-pill {
  @apply bg-[var(--bg-tertiary)] text-[var(--text-muted)] 
         border border-[var(--border-default)] font-normal;
}

/* For sidebar / nav items */
.nav-active {
  @apply bg-[var(--accent-muted)] text-[var(--accent)] 
         border-l-[3px] border-l-[var(--accent)];
}
```

### 4.4.2 Action Guidance — Highlighting the Next Step

Guide users to the most important action on each screen using visual hierarchy and subtle animation. The user should always know "what do I do next?" without reading instructions.

```
GUIDANCE PATTERN HIERARCHY (strongest → subtlest):

1. PULSING PRIMARY CTA (strongest — use sparingly)
   ────────────────────────────────────────────────
   When: The user MUST take an action to proceed (e.g., "Give Consent",
         "Complete Setup", "Verify OTP")
   
   Style:
     background: linear-gradient(135deg, accent, accentHover)
     box-shadow: 0 0 0 0 rgba(accent, 0.4)
     animation: pulse-ring 2s ease-in-out infinite
   
   @keyframes pulse-ring {
     0%   { box-shadow: 0 0 0 0 rgba(var(--accent-rgb), 0.4); }
     70%  { box-shadow: 0 0 0 8px rgba(var(--accent-rgb), 0); }
     100% { box-shadow: 0 0 0 0 rgba(var(--accent-rgb), 0); }
   }
   
   Only ONE pulsing button per screen. Never pulse secondary actions.
   Pulse stops after the user hovers or clicks.

2. HIGHLIGHTED CARD / SECTION (moderate — draws eye)
   ────────────────────────────────────────────────
   When: A section needs attention but isn't blocking (e.g., 
         "Complete your profile", "3 items need review")
   
   Style:
     border: 1px solid accentBorder (instead of normal borderDefault)
     background: accentMuted (subtle tint, 8-10% opacity)
     Optional: small accent-colored dot or badge in corner
   
   Example:
   ┌─ accentBorder ──────────────────────────────────────┐
   │  ⚠️ 3 documents pending upload                     │
   │  Complete your verification to enable payments.     │
   │                                          [Upload →] │
   └─────────────────────────────────────────────────────┘

3. NUMBERED STEP INDICATOR (sequential guidance)
   ────────────────────────────────────────────────
   When: Multi-step process (onboarding, setup wizard, form steps)
   
   Style:
     Completed step: bg-success circle, white checkmark, success-colored label
     Current step:   bg-accent circle (pulsing), accent-colored label, font-weight 700
     Future step:    bg transparent, borderDefault circle, textMuted label
     Connector line: completed=success, current=accent dashed, future=borderDefault
   
   ①──────②──────③──────④
   ✓ Done   ● Current  ○ Next   ○ Last

4. BADGE NUDGE (subtle — "hey, look here")
   ────────────────────────────────────────────────
   When: A nav item or tab has pending items (e.g., "My Tasks (5)")
   
   Style:
     Count badge: min-w-[18px] h-[18px] rounded-full
     bg-accent text-accentText font-size 9px font-weight 700
     Position: inline after label, or absolute top-right of icon
   
   Animated entry: scale 0→1 with bounce (spring, 300ms) when count increases
   
   For destructive/urgent: bg-danger instead of bg-accent

5. EMPTY STATE WITH CTA (contextual guidance)
   ────────────────────────────────────────────────
   When: No data yet, user needs to create/add something
   
   Style:
     Center-aligned, muted icon (48px, 50% opacity)
     Heading: "No vendors yet" (heading-sm, textMuted)
     Description: "Start by onboarding your first vendor" (body-sm, textDisabled)
     CTA button: PRIMARY gradient (this is the main action on an empty screen)
     
   The CTA on an empty state IS the primary button — use the gradient style.
   It replaces the header-level "Create" button as the visual anchor.

6. TOOLTIP SPOTLIGHT (educational — first-time users)
   ────────────────────────────────────────────────
   When: Feature is new or user hasn't used it yet
   
   Style:
     Tooltip arrow pointing at the target element
     bg-accent text-accentText rounded-[8px] p-3
     "Try this: Click here to..." (body-sm, font-weight 500)
     [Got it] dismiss button (ghost, text-accentText)
     
   Shown once per user (dismissed state in localStorage)
   Max 1 tooltip visible at a time — never stack them

7. INLINE HELPER TEXT (gentlest — always visible)
   ────────────────────────────────────────────────
   When: Field or section needs context
   
   Style:
     text-foreground-muted, font-size 11px, margin-top 4px
     Below the element it describes
     Optional info icon (ℹ️ 12px) before the text
   
   Examples:
     "Enter GSTIN to auto-fetch company details"
     "Amount will be converted using today's exchange rate"
```

**Action guidance DO's and DON'Ts:**

```
DO:
  ✓ Use exactly ONE primary (pulsing) CTA per screen
  ✓ Use highlighted cards for sections that need attention
  ✓ Use badge nudges on nav items with pending counts
  ✓ Use empty states with clear CTAs when there's no data
  ✓ Stop the pulse animation on hover (user found it)
  ✓ Use step indicators for multi-step flows

DON'T:
  ✗ Never pulse more than one button
  ✗ Never use red/danger styling for guidance (red = error/destructive only)
  ✗ Never auto-open tooltips/popovers — only show on first visit
  ✗ Never animate multiple elements simultaneously (overwhelming)
  ✗ Never use flashing/blinking — pulse is subtle opacity, not strobe
  ✗ Never block content with guidance — it should enhance, not obscure
```

### 4.5 Input Fields

```
Default:
  background: bgTertiary
  border: 1px solid borderDefault
  border-radius: 8px
  padding: 8px 12px
  font-size: 12px
  color: textPrimary
  placeholder-color: textMuted
  focus: border-color → accent, ring: 0 0 0 2px accentMuted

Search Input:
  Same as default + 🔍 icon on left + keyboard shortcut badge (⌘K) on right
  min-width: 220px in header

Select/Dropdown:
  Same styling as input + chevron icon on right
  Dropdown menu: bg-app-card, border, rounded-[10px], shadow-app-md
  Option hover: bg-app-input
  Selected option: text-accent, bg-accentMuted

Filter Pill Buttons:
  padding: 8px 14px, rounded-[8px]
  bg-app-input, border border-app-border
  text-foreground-muted, font-size 11px
  Active: bg-accentMuted, border-accentBorder, text-accent
```

### 4.6 Cards

```
Standard Card:
  bg-app-card, border border-app-border, rounded-[16px], p-[22px]
  No box-shadow by default (shadows reserved for floating elements only)

Hoverable Card (clickable):
  Same + cursor-pointer
  hover: border-color transitions to borderLight (subtle lightening)

Feature Card:
  Same as standard
  hover: border-color transitions to the card's accent color at 40% opacity

Nested Card (inside a card):
  bg-app-input, border border-app-border, rounded-[10px], p-[10px 12px]
```

### 4.7 Tab Bar

```
Container: flex gap-1, border-bottom 1px solid borderDefault
Tab Button:
  padding: 8px 16px
  border-radius: 8px 8px 0 0
  font-size: 12px

  Inactive: bg transparent, text-foreground-muted, font-weight 450
  Active: bg-app-card, text-foreground, font-weight 600,
          border-bottom: 2px solid accent

  Count badge inside tab: padding 1px 7px, rounded-[5px], font-size 9px, font-weight 600
    Inactive: bg-app-input, text-foreground-disabled
    Active: bg-accentMuted, text-accent
```

### 4.8 Sidebar

```
Width: 240px expanded, 64px collapsed (transition: width 0.2s ease)
Background: bgSidebar
Border-right: 1px solid borderDefault

Logo Area (top):
  Padding: 18px
  Logo: 30x30 rounded-[8px] with gradient background
  Text: App name heading-md bold, tagline tiny uppercase

Section Labels:
  color: text-foreground-muted, font-size: 11px, padding: 6px 14px
  Hover: text-foreground-secondary
  Collapsible with chevron arrow, state persisted in localStorage

Nav Items:
  padding: 8px 14px, rounded-[10px], gap 12px (icon + label)
  Icon: 16px, width 20px centered
  Label: body-sm weight

  Default text: text-foreground-secondary
  Hover text: text-foreground
  Hover background: bg-app-input/40
  Active: text-accent, 3px solid left border, bg-accentMuted
  NEVER hardcode hex colors — always use Tailwind text-foreground-* classes

  Badge (count): min-width 18px, height 18px, rounded full,
    bg-accent, text-accentText, font-size 9px, font-weight 700

User Section (bottom):
  border-top 1px solid borderDefault, padding 14px 16px
  Avatar: 32x32 rounded-[8px] with gradient
  Name: body-sm font-semibold
  Role: tiny text-foreground-disabled

Collapsed mode:
  Items show only icons, centered
  Tooltips on hover (use shadcn Tooltip)
```

### 4.9 Modals / Dialogs

```
Overlay: bgOverlay with backdrop-blur-sm
Container: bg-app-card, rounded-[16px], border border-app-border, shadow-app-lg
  Header: px-6 py-4, border-bottom, heading-md, close button (X) top-right
  Body: px-6 py-5
  Footer: px-6 py-4, border-top, flex justify-end gap-3
  
Sizes: sm (480px), md (640px), lg (800px), xl (960px)
Animation: fade-in + scale from 0.95 to 1.0 (150ms ease-out)
```

### 4.10 Toast / Notifications

```
Position: bottom-right, stacked
Container: bg-app-card, border border-app-border, rounded-[12px], shadow-app-md
  padding: 14px 18px, max-width 380px
  Left color strip: 4px width, rounded-l, colored by type
  Title: body-sm font-semibold
  Description: caption text-foreground-muted
  Auto-dismiss: 4 seconds (success), 8 seconds (error), persistent (action required)
```

### 4.11 Empty State

```
Center-aligned in the container:
  Icon: 48px, text-foreground-disabled, opacity 0.5
  Title: heading-sm, text-foreground-muted, margin-top 12px
  Description: body-sm, text-foreground-disabled, max-width 320px, margin-top 4px
  Action button: secondary button, margin-top 16px
```

### 4.12 Loading Skeleton

```
Use animated shimmer skeletons matching the exact layout of content being loaded.

Skeleton element: rounded-[8px], bg-app-input
Shimmer: linear-gradient sweep animation (1.5s infinite)
  From: bg-app-input → lighter → bg-app-input

For tables: render 5 skeleton rows matching column widths
For metric cards: render card shape with circle (icon) + 2 rectangles (value + label)
For detail pages: render header skeleton + tab bar + content skeleton

ALWAYS show skeletons during data loading. NEVER show a blank page or spinner-only.
```

---

## 5. PAGE LAYOUT PATTERNS

### Standard List Page
```
┌──────────────────────────────────────────────────────────────┐
│ PAGE HEADER                                                   │
│ Title (heading-md)              [Filter] [Filter] [+ Create]  │
│ Subtitle (caption, text-muted)                                │
├──────────────────────────────────────────────────────────────┤
│ TAB BAR                                                       │
│ [All (342)] [Category A (298)] [Category B (28)] [Special]    │
├──────────────────────────────────────────────────────────────┤
│ FILTER BAR                                                    │
│ [🔍 Search...        ] [Status ▾] [Category ▾] [Date ▾]     │
├──────────────────────────────────────────────────────────────┤
│ DATA TABLE (see 4.2)                                          │
├──────────────────────────────────────────────────────────────┤
│ PAGINATION                                                    │
│ Showing 1-20 of 342                    ← 1 2 3 ... 35 →      │
└──────────────────────────────────────────────────────────────┘
```

### Standard Detail Page
```
┌──────────────────────────────────────────────────────────────┐
│ BREADCRUMB: Items › ITEM-2026-0847                            │
├──────────────────────────────────────────────────────────────┤
│ ENTITY HEADER CARD                                            │
│ ID + Badges           [Secondary] [Danger] [Primary Action]   │
│ Key info line                                                 │
│ ──────────────────────────────────                            │
│ Key Metrics Grid (4-6 columns)                                │
├──────────────────────────────────────────────────────────────┤
│ TAB BAR                                                       │
│ [Details] [Items] [History] [Validation] [Audit]              │
├──────────────────────────────────────────────────────────────┤
│ TAB CONTENT (varies by tab)                                   │
└──────────────────────────────────────────────────────────────┘
```

### Dashboard Page
```
┌──────────────────────────────────────────────────────────────┐
│ GREETING: Good morning, User                                  │
│ Subtitle: Here's your overview                                │
├──────────────────────────────────────────────────────────────┤
│ METRIC CARDS (4 in a row, flex wrap)                          │
│ [Metric 1] [Metric 2] [Metric 3] [Metric 4]                  │
├──────────────────────────────────────────────────────────────┤
│ CHARTS ROW (2-column grid)                                    │
│ [Chart 1]                       [Chart 2]                     │
├──────────────────────────────────────────────────────────────┤
│ ACTIVITY ROW (2-column grid)                                  │
│ [Activity Feed]                 [Summary Cards]               │
└──────────────────────────────────────────────────────────────┘
```

### Multi-Step Form
```
┌──────────────────────────────────────────────────────────────┐
│ STEP INDICATOR                                                │
│ ① Step 1 ──── ② Step 2 ──── ③ Step 3 ──── ④ Review          │
│ (active step highlighted with accent, completed with check)   │
├──────────────────────────────────────────────────────────────┤
│ FORM CONTENT (varies per step)                                │
├──────────────────────────────────────────────────────────────┤
│ FORM FOOTER                                                   │
│ [Save as Draft]                     [← Back]  [Next Step →]   │
└──────────────────────────────────────────────────────────────┘
```

---

## 6. ANIMATION & MOTION

| Element | Animation | Duration | Easing |
|---------|-----------|----------|--------|
| Page transition | fade-in + translateY(-4px → 0) | 200ms | ease-out |
| Card hover border | border-color transition | 150ms | ease |
| Table row hover bg | background transition | 150ms | ease |
| Modal enter | opacity 0→1 + scale 0.95→1.0 | 150ms | ease-out |
| Modal exit | opacity 1→0 + scale 1.0→0.95 | 100ms | ease-in |
| Toast enter | translateX(100% → 0) | 250ms | ease-out |
| Toast exit | translateX(0 → 100%) | 200ms | ease-in |
| Skeleton shimmer | gradient sweep left to right | 1500ms | linear, infinite |
| Sidebar collapse | width transition | 200ms | ease |
| Dropdown open | opacity + translateY(-4px → 0) | 120ms | ease-out |
| Button press | scale(0.97) | 75ms | ease |
| Chart bar | height 0→full | 400ms | ease-out, staggered 50ms per bar |
| Progress bar fill | width 0→value | 500ms | ease-out |
| CTA pulse ring | box-shadow expand + fade | 2000ms | ease-in-out, infinite |
| Badge nudge entry | scale 0→1 | 300ms | spring (bounce allowed here only) |
| Selected state transition | bg + color + border | 150ms | ease |
| Step indicator current | box-shadow opacity pulse | 2000ms | ease-in-out, infinite |

Keep animations subtle and fast. Professional software must feel snappy and precise. Never use bounce or spring physics — **exception:** badge count nudge uses a brief spring for playful entry.

---

## 7. RESPONSIVE BREAKPOINTS

| Breakpoint | Width | Layout Changes |
|-----------|-------|----------------|
| Desktop XL | ≥1440px | Full layout, all columns visible |
| Desktop | ≥1024px | Standard layout |
| Tablet | ≥768px | Sidebar collapses to icon-only, table scrolls horizontally, grid goes 1-col |
| Mobile | <768px | Sidebar → hamburger, single column, stack metric cards |

### Mobile-Specific Rules
- Sidebar: hidden by default, hamburger menu toggle
- Tables: horizontally scrollable with sticky first column
- Metric cards: 2 per row on tablet, 1 per row on mobile
- Forms: full-width single column
- Modals: full-screen on mobile (rounded-t-[16px], slides up from bottom)
- Action buttons: full-width sticky bottom bar on mobile

---

## 8. ACCESSIBILITY

- All interactive elements: `focus-visible` outline using `accent` color
- Color contrast: minimum 4.5:1 for text, 3:1 for large text (both themes pass)
- All form inputs have associated labels (visible or sr-only)
- Tables: proper `<thead>`, `<th scope="col">` structure
- Modals trap focus and close on Escape
- Toasts have role="alert" for screen readers
- Icons with meaning have aria-labels
- Status badges are not color-only — they include text labels
- All buttons: minimum 44px touch target on mobile

---

## 9. CHART STYLING (Recharts)

All charts MUST import from `@/lib/chart-theme.ts` — single source of truth.

```tsx
import {
  CHART_TOOLTIP_STYLE,  // contentStyle for <Tooltip>
  CHART_LABEL_STYLE,    // labelStyle for <Tooltip>
  CHART_AXIS_TICK,      // tick prop for <XAxis>/<YAxis> — { fontSize: 10, fill: "var(--text-muted)" }
  CHART_GRID_STROKE,    // stroke for <CartesianGrid> — "var(--border-default)"
  CHART_GRID_OPACITY,   // strokeOpacity — 0.4
  CHART_CURSOR,         // cursor for <Tooltip> — { fill: "var(--accent-muted)" }
  CHART_COLORS,         // palette — accent, info, success, purple, teal, warning, danger
} from "@/lib/chart-theme";
```

Tooltip spec: `bg-secondary` background, `border-default` 1px border, 10px radius, `text-primary` color, subtle shadow. Auto-adapts via CSS variables.

Bar chart: radius `[4,4,0,0]`, prefer gradient fills.
Line chart: strokeWidth 2, dot r=4.
Pie chart: innerRadius 50, outerRadius 80-90, paddingAngle 3.

---

## 10. FORMATTING HELPERS

Use these consistently across the ENTIRE app. Import from `@/lib/utils`.

```typescript
// Format date
formatDate(date: string | Date): string
  // → "15/02/2026" (DD/MM/YYYY)

// Format datetime
formatDateTime(date: string | Date): string
  // → "15/02/2026, 2:30 PM"

// Format relative time
formatRelative(date: string | Date): string
  // → "2 hours ago", "3 days ago", "Just now"

// Format numbers with locale-appropriate separators
formatNumber(value: number): string

// Format compact numbers for metric cards
formatCompact(value: number): string
  // 1234 → "1.2K", 1250000 → "1.25M" etc.
```

---

## SUMMARY: Top 12 UI Rules

1. **Use theme CSS variables everywhere** — never hardcode colors. `bg-app-card`, `text-foreground`, `border-app-border`.
2. **Metric values in display size** (24px/800), labels in overline size (10px/600 uppercase).
3. **IDs, codes, reference numbers always in `font-mono`**.
4. **Locale-appropriate formatting** for numbers, dates, currency.
5. **DD/MM/YYYY** dates everywhere — never ISO or US format.
6. **Status badges with border** — never solid fill. Use the muted background + border + colored text pattern.
7. **One primary CTA button per screen** — gradient amber with glow shadow. Everything else is secondary or ghost.
8. **Card border-radius: 16px**, button/input: 10px, badge: 6px — consistent everywhere.
9. **Loading skeletons** matching exact content layout — never blank screens or centered spinners.
10. **Subtle animations only** — fast transitions (150ms), no bounce, no spring. Professional software must feel snappy and precise.
11. **Selected states use accent tinting** — `bg-accentMuted` + `text-accent` + `border-accentBorder`. Same pattern for filter pills, tabs, sidebar items, pagination, toggle groups. Adapts automatically across all 15 themes.
12. **One pulsing CTA per screen maximum** — guide users to the next action with a subtle pulse ring animation. Use highlighted cards for sections needing attention, badge nudges for nav counts, and step indicators for multi-step flows. Never pulse more than one element.

---

## TECH STACK

| Layer | Technology | Notes |
|-------|-----------|-------|
| Frontend | Next.js 14, TypeScript, Tailwind CSS, shadcn/ui | App Router, 15-variant theme system |
| State | Zustand | Theme, auth, app state persisted in localStorage |
| Components | shadcn/ui | Customized with theme CSS variables |
| Icons | Lucide React | Consistent 14-16px sizes |
| Charts | Recharts | Themed via chart-theme.ts |
| Fonts | System font stack | No external fonts |
