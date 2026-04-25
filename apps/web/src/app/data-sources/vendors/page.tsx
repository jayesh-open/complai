export default function ImportedVendorsPage() {
  return (
    <div data-testid="imported-vendors-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Imported Vendors</h2>
        <p className="text-body text-foreground-muted mt-1">
          Vendor master synced from Apex P2P for compliance scoring
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <p className="text-body text-foreground-muted">
          No vendors imported yet. Connect Apex P2P to begin syncing.
        </p>
      </div>
    </div>
  );
}
