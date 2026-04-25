import { RefreshCw } from "lucide-react";

export default function SyncStatusPage() {
  return (
    <div data-testid="sync-status-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Sync Status</h2>
        <p className="text-body text-foreground-muted mt-1">
          Real-time sync status for all connected data sources
        </p>
      </div>
      <div className="bg-app-card border border-app-border rounded-card p-5">
        <div className="flex items-center gap-3 text-foreground-muted">
          <RefreshCw className="w-5 h-5" />
          <span className="text-body">Sync dashboard will be available when sibling gateways are connected.</span>
        </div>
      </div>
    </div>
  );
}
