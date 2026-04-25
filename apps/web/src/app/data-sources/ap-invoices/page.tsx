export default function ImportedAPInvoicesPage() {
  return (
    <div data-testid="ap-invoices-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Imported AP Invoices</h2>
        <p className="text-body text-foreground-muted mt-1">
          AP invoices synced from Apex P2P for ITC reconciliation
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <p className="text-body text-foreground-muted">
          No AP invoices imported yet. Connect Apex P2P to begin syncing.
        </p>
      </div>
    </div>
  );
}
