export default function PublicITRLayout({ children }: { children: React.ReactNode }) {
  return (
    <div className="fixed inset-0 z-[100] bg-[var(--bg-primary)] overflow-y-auto">
      <div className="min-h-screen flex flex-col">
        <header className="border-b border-[var(--border-default)] bg-[var(--bg-secondary)] px-6 py-3 flex items-center gap-3">
          <div className="w-7 h-7 rounded-lg bg-[var(--accent)] flex items-center justify-center">
            <span className="text-[var(--accent-text)] text-xs font-bold">C</span>
          </div>
          <span className="text-sm font-semibold text-[var(--text-primary)]">Complai</span>
          <span className="text-[10px] px-2 py-0.5 rounded-full bg-[var(--info-muted)] text-[var(--info)] font-medium">
            Employee Portal
          </span>
        </header>
        <main className="flex-1 p-7 max-w-5xl mx-auto w-full">
          {children}
        </main>
        <footer className="border-t border-[var(--border-default)] px-6 py-3 text-center">
          <p className="text-[10px] text-[var(--text-muted)]">
            Powered by Complai · Secure link · Do not share
          </p>
        </footer>
      </div>
    </div>
  );
}
