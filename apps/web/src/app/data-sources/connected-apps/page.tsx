import { Link2 } from "lucide-react";

export default function ConnectedAppsPage() {
  return (
    <div data-testid="connected-apps-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Connected Apps</h2>
        <p className="text-body text-foreground-muted mt-1">
          Bank Open sibling applications connected to Complai
        </p>
      </div>
      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4">
        {[
          { name: "Apex P2P", status: "Connected", desc: "Vendor master, AP invoices, payments" },
          { name: "Aura O2C", status: "Pending", desc: "Customer master, AR invoices" },
          { name: "Bridge", status: "Pending", desc: "Contracts, secretarial obligations" },
          { name: "HRMS", status: "Pending", desc: "Payroll data, Form 16" },
        ].map((app) => (
          <div key={app.name} className="bg-app-card border border-app-border rounded-card p-5">
            <div className="flex items-center gap-3 mb-3">
              <div className="w-10 h-10 rounded-lg bg-app-input flex items-center justify-center">
                <Link2 className="w-5 h-5 text-foreground-muted" />
              </div>
              <div>
                <div className="text-body-sm font-semibold text-foreground">{app.name}</div>
                <span className={`text-[10px] font-semibold uppercase tracking-wide ${app.status === "Connected" ? "text-app-success" : "text-foreground-muted"}`}>
                  {app.status}
                </span>
              </div>
            </div>
            <p className="text-caption text-foreground-muted">{app.desc}</p>
          </div>
        ))}
      </div>
    </div>
  );
}
