export default function ImportedContractsPage() {
  return (
    <div data-testid="imported-contracts-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Imported Contracts</h2>
        <p className="text-body text-foreground-muted mt-1">
          Contracts synced from Bridge for TDS section determination
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <p className="text-body text-foreground-muted">
          No contracts imported yet. Connect Bridge to begin syncing.
        </p>
      </div>
    </div>
  );
}
