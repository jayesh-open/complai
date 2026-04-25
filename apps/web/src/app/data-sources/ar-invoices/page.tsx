export default function ImportedARInvoicesPage() {
  return (
    <div data-testid="ar-invoices-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Imported AR Invoices</h2>
        <p className="text-body text-foreground-muted mt-1">
          AR invoices synced from Aura O2C for GSTR-1 filing
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <p className="text-body text-foreground-muted">
          No AR invoices imported yet. Connect Aura O2C to begin syncing.
        </p>
      </div>
    </div>
  );
}
