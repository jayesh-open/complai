export default function ImportedPayrollPage() {
  return (
    <div data-testid="imported-payroll-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Imported Payroll Data</h2>
        <p className="text-body text-foreground-muted mt-1">
          Payroll and Form 16 data synced from HRMS for TDS 24Q and ITR filing
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <p className="text-body text-foreground-muted">
          No payroll data imported yet. Connect HRMS to begin syncing.
        </p>
      </div>
    </div>
  );
}
