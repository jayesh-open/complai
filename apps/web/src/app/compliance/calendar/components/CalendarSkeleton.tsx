export function CalendarSkeleton() {
  return (
    <div className="border border-[var(--border-default)] rounded-xl overflow-hidden" data-testid="calendar-skeleton">
      <div className="grid grid-cols-7 border-b border-[var(--border-default)]">
        {Array.from({ length: 7 }).map((_, i) => (
          <div key={i} className="py-2 bg-[var(--bg-secondary)]">
            <div className="h-3 w-8 mx-auto rounded bg-[var(--bg-tertiary)] animate-pulse" />
          </div>
        ))}
      </div>
      <div className="grid grid-cols-7" style={{ gridAutoRows: "110px" }}>
        {Array.from({ length: 42 }).map((_, i) => (
          <div
            key={i}
            className="p-1.5 border-b border-r border-[var(--border-default)] bg-[var(--bg-primary)]"
          >
            <div className="h-3 w-4 rounded bg-[var(--bg-tertiary)] animate-pulse mb-2" />
            <div className="space-y-1">
              <div className="h-3 w-3/4 rounded bg-[var(--bg-tertiary)] animate-pulse" />
              {i % 3 === 0 && (
                <div className="h-3 w-1/2 rounded bg-[var(--bg-tertiary)] animate-pulse" />
              )}
            </div>
          </div>
        ))}
      </div>
    </div>
  );
}
