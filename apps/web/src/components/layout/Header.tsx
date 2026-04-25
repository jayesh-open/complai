"use client";

import { usePathname } from "next/navigation";
import { Search, Bell, Menu } from "lucide-react";
import { cn } from "@/lib/utils";

interface HeaderProps {
  onCommandPaletteOpen: () => void;
  onMobileMenuToggle?: () => void;
}

const PAGE_TITLES: Record<string, string> = {
  "/dashboard": "Dashboard",
  "/tasks": "My Tasks",
  "/inbox": "Inbox",
  "/compliance/gst": "GST Returns",
  "/compliance/e-invoice": "E-Invoicing",
  "/compliance/e-way-bill": "E-Way Bill",
  "/compliance/itc-recon": "ITC Reconciliation",
  "/compliance/vendor-compliance": "Vendor Compliance",
  "/compliance/tds": "TDS/TCS",
  "/compliance/itr": "ITR",
  "/compliance/secretarial": "Secretarial",
  "/insights/cfo": "CFO Dashboard",
  "/insights/reports": "Reports & Analytics",
  "/insights/audit-trail": "Audit Trail",
  "/data-sources/connected-apps": "Connected Apps",
  "/data-sources/sync-status": "Sync Status",
  "/data-sources/ar-invoices": "Imported AR Invoices",
  "/data-sources/ap-invoices": "Imported AP Invoices",
  "/data-sources/vendors": "Imported Vendors",
  "/data-sources/contracts": "Imported Contracts",
  "/data-sources/payroll": "Imported Payroll Data",
  "/documents": "Documents",
  "/configure/settings": "Settings",
  "/configure/appearance": "Appearance",
  "/configure/users": "Users & Roles",
  "/configure/integrations": "Integrations",
};

export function Header({ onCommandPaletteOpen, onMobileMenuToggle }: HeaderProps) {
  const pathname = usePathname();
  const title = PAGE_TITLES[pathname ?? ""] ?? "Complai";

  return (
    <header
      className="h-[52px] border-b border-app-border bg-app-card flex items-center justify-between px-5 flex-shrink-0"
      data-testid="header"
    >
      <div className="flex items-center gap-3">
        {onMobileMenuToggle && (
          <button
            onClick={onMobileMenuToggle}
            className="md:hidden p-1 text-foreground-muted"
            data-testid="mobile-menu-btn"
          >
            <Menu className="w-5 h-5" />
          </button>
        )}
        <h1 className="text-heading-md text-foreground" data-testid="page-title">{title}</h1>
      </div>

      <div className="flex items-center gap-3">
        <button
          onClick={onCommandPaletteOpen}
          className={cn(
            "flex items-center gap-2 px-3 py-1.5 rounded-lg border",
            "bg-app-input border-app-border text-foreground-muted text-xs",
            "hover:border-[var(--border-focus)] transition-colors min-w-[220px]",
          )}
          data-testid="search-trigger"
        >
          <Search className="w-3.5 h-3.5" />
          <span className="flex-1 text-left">Search...</span>
          <kbd className="text-tiny bg-app-bg px-1.5 py-0.5 rounded border border-app-border">
            ⌘K
          </kbd>
        </button>

        <button className="relative p-2 rounded-lg border border-app-border text-foreground-muted hover:bg-app-input transition-colors">
          <Bell className="w-4 h-4" />
          <span className="absolute -top-1 -right-1 min-w-[16px] h-4 rounded-full bg-app-danger text-white text-[9px] font-bold flex items-center justify-center px-1">
            12
          </span>
        </button>

        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-app-purple to-app-info flex items-center justify-center cursor-pointer">
          <span className="text-white text-xs font-bold">JH</span>
        </div>
      </div>
    </header>
  );
}
