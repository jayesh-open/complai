"use client";

import { useAppStore, type DensityMode } from "@/store/app-store";
import {
  THEME_MAP, THEME_DISPLAY_NAMES, THEME_FAMILIES,
  type ThemeMode,
} from "@/lib/themes";
import { cn } from "@/lib/utils";

export default function AppearancePage() {
  const theme = useAppStore((s) => s.theme);
  const setTheme = useAppStore((s) => s.setTheme);
  const density = useAppStore((s) => s.density);
  const setDensity = useAppStore((s) => s.setDensity);

  return (
    <div className="max-w-3xl" data-testid="appearance-page">
      <h2 className="text-heading-lg text-foreground mb-1">Appearance</h2>
      <p className="text-body text-foreground-muted mb-6">
        Customize the look and feel of your Complai workspace.
      </p>

      <section className="mb-8">
        <h3 className="text-heading-sm text-foreground mb-4">Theme</h3>
        {THEME_FAMILIES.map((family) => (
          <div key={family.family} className="mb-5">
            <h4 className="text-overline text-foreground-muted mb-2">{family.label}</h4>
            <div className="grid grid-cols-3 gap-3">
              {family.modes.map((mode) => {
                const colors = THEME_MAP[mode];
                const active = theme === mode;
                return (
                  <button
                    key={mode}
                    onClick={() => setTheme(mode)}
                    data-testid={`theme-card-${mode}`}
                    className={cn(
                      "relative rounded-xl p-3 border-2 transition-all text-left",
                      active
                        ? "border-[var(--accent)] shadow-app-accent"
                        : "border-app-border hover:border-app-border-lt",
                    )}
                    style={{ backgroundColor: colors.bgPrimary }}
                  >
                    {active && (
                      <div className="absolute top-2 right-2 w-5 h-5 rounded-full flex items-center justify-center"
                        style={{ backgroundColor: colors.accent }}
                      >
                        <span style={{ color: colors.accentText }} className="text-[10px]">✓</span>
                      </div>
                    )}
                    <div className="flex gap-1.5 mb-2">
                      <div className="w-4 h-4 rounded" style={{ backgroundColor: colors.bgSidebar, border: `1px solid ${colors.borderDefault}` }} />
                      <div className="w-4 h-4 rounded" style={{ backgroundColor: colors.accent }} />
                      <div className="w-4 h-4 rounded" style={{ backgroundColor: colors.success }} />
                    </div>
                    <div className="text-[11px] font-semibold" style={{ color: colors.textPrimary }}>
                      {THEME_DISPLAY_NAMES[mode]}
                    </div>
                    <div className="flex gap-1 mt-1">
                      {[colors.bgSecondary, colors.bgTertiary, colors.borderDefault].map((c, i) => (
                        <div key={i} className="w-3 h-2 rounded-sm" style={{ backgroundColor: c }} />
                      ))}
                    </div>
                  </button>
                );
              })}
            </div>
          </div>
        ))}
      </section>

      <section className="mb-8">
        <h3 className="text-heading-sm text-foreground mb-4">Density</h3>
        <div className="flex gap-3">
          {(["compact", "comfortable", "spacious"] as DensityMode[]).map((d) => (
            <button
              key={d}
              onClick={() => setDensity(d)}
              data-testid={`density-${d}`}
              aria-pressed={density === d}
              className={cn(
                "px-4 py-2 rounded-btn text-xs font-medium border transition-colors capitalize",
                density === d
                  ? "selected-pill"
                  : "unselected-pill hover:bg-app-input",
              )}
            >
              {d}
              <span className="block text-[9px] mt-0.5 opacity-70">
                {d === "compact" ? "40px rows" : d === "comfortable" ? "52px rows" : "64px rows"}
              </span>
            </button>
          ))}
        </div>
      </section>
    </div>
  );
}
