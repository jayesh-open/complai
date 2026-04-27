"use client";

import { usePathname } from "next/navigation";
import Link from "next/link";
import {
  LayoutDashboard, ListTodo, Inbox, FileText,
  Receipt, FileCheck2, Truck,
  FileSpreadsheet, GitCompareArrows, BarChart3, Gauge, History,
  FolderOpen, Mail, Settings, Users, Workflow, ShieldAlert,
  Wallet, ChevronDown, ChevronRight, PanelLeftClose, PanelLeftOpen,
  Sparkles, RefreshCw, Link2, ArrowDownToLine,
} from "lucide-react";
import { useAppStore } from "@/store/app-store";
import { cn } from "@/lib/utils";

interface NavItem {
  label: string;
  href: string;
  icon: React.ElementType;
  badge?: number;
}

interface NavGroup {
  id: string;
  label?: string;
  items: NavItem[];
}

const NAV_GROUPS: NavGroup[] = [
  {
    id: "core",
    items: [
      { label: "Dashboard", href: "/dashboard", icon: LayoutDashboard },
      { label: "My Tasks", href: "/tasks", icon: ListTodo, badge: 5 },
      { label: "Inbox", href: "/inbox", icon: Inbox, badge: 12 },
    ],
  },
  {
    id: "compliance",
    label: "COMPLIANCE",
    items: [
      { label: "GST Returns", href: "/compliance/gst", icon: FileSpreadsheet },
      { label: "E-Invoicing", href: "/compliance/e-invoicing", icon: FileCheck2 },
      { label: "E-Way Bill", href: "/compliance/e-way-bill", icon: Truck },
      { label: "ITC Reconciliation", href: "/compliance/itc-recon", icon: GitCompareArrows },
      { label: "Vendor Compliance", href: "/compliance/vendor-compliance", icon: Gauge },
      { label: "TDS/TCS", href: "/compliance/tds", icon: Receipt },
      { label: "ITR", href: "/compliance/itr", icon: FileSpreadsheet },
      { label: "Secretarial", href: "/compliance/secretarial", icon: ShieldAlert },
    ],
  },
  {
    id: "insights",
    label: "INSIGHTS",
    items: [
      { label: "CFO Dashboard", href: "/insights/cfo", icon: Sparkles },
      { label: "Reports & Analytics", href: "/insights/reports", icon: BarChart3 },
      { label: "Audit Trail", href: "/insights/audit-trail", icon: History },
    ],
  },
  {
    id: "data-sources",
    label: "DATA SOURCES",
    items: [
      { label: "Connected Apps", href: "/data-sources/connected-apps", icon: Link2 },
      { label: "Sync Status", href: "/data-sources/sync-status", icon: RefreshCw },
      { label: "Imported AR Invoices", href: "/data-sources/ar-invoices", icon: ArrowDownToLine },
      { label: "Imported AP Invoices", href: "/data-sources/ap-invoices", icon: ArrowDownToLine },
      { label: "Imported Vendors", href: "/data-sources/vendors", icon: ArrowDownToLine },
      { label: "Imported Contracts", href: "/data-sources/contracts", icon: FileText },
      { label: "Imported Payroll Data", href: "/data-sources/payroll", icon: ArrowDownToLine },
    ],
  },
  {
    id: "documents",
    label: "DOCUMENTS",
    items: [
      { label: "Documents", href: "/documents", icon: FolderOpen },
      { label: "Email Inbox", href: "/documents/email-inbox", icon: Mail },
    ],
  },
  {
    id: "configure",
    label: "CONFIGURE",
    items: [
      { label: "Settings", href: "/configure/settings", icon: Settings },
      { label: "Users & Roles", href: "/configure/users", icon: Users },
      { label: "Approval Workflows", href: "/configure/workflows", icon: Workflow },
      { label: "GST Configuration", href: "/configure/gst", icon: FileSpreadsheet },
      { label: "TDS Configuration", href: "/configure/tds", icon: Receipt },
      { label: "Integrations", href: "/configure/integrations", icon: Link2 },
      { label: "Billing", href: "/configure/billing", icon: Wallet },
    ],
  },
];

export function Sidebar() {
  const pathname = usePathname();
  const collapsed = useAppStore((s) => s.sidebarCollapsed);
  const toggleSidebar = useAppStore((s) => s.toggleSidebar);
  const groupState = useAppStore((s) => s.sidebarGroupState);
  const toggleGroup = useAppStore((s) => s.toggleSidebarGroup);

  return (
    <aside
      className={cn(
        "h-screen flex flex-col border-r border-app-border bg-app-sidebar transition-all duration-200 flex-shrink-0",
        collapsed ? "w-16" : "w-60",
      )}
      data-testid="sidebar"
    >
      <div className="p-4 flex items-center gap-3">
        <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-app-accent to-app-accent-h flex items-center justify-center flex-shrink-0">
          <span className="text-app-accent-t font-bold text-sm">C</span>
        </div>
        {!collapsed && (
          <div className="min-w-0">
            <div className="text-heading-md text-foreground truncate">Complai</div>
            <div className="text-tiny uppercase tracking-wider text-foreground-muted">
              Compliance Platform
            </div>
          </div>
        )}
      </div>

      <nav className="flex-1 overflow-y-auto px-2 pb-2" data-testid="sidebar-nav">
        {NAV_GROUPS.map((group) => (
          <div key={group.id} className="mb-1" data-testid={`sidebar-group-${group.id}`}>
            {group.label && !collapsed && (
              <button
                onClick={() => toggleGroup(group.id)}
                className="flex items-center gap-1 w-full px-3 py-1.5 text-overline text-foreground-muted hover:text-foreground-secondary"
              >
                {groupState[group.id] ? (
                  <ChevronRight className="w-3 h-3" />
                ) : (
                  <ChevronDown className="w-3 h-3" />
                )}
                <span data-testid="sidebar-group-label">{group.label}</span>
              </button>
            )}
            {group.label && collapsed && (
              <div className="mx-3 my-2 border-t border-app-border" />
            )}
            {(!group.label || !groupState[group.id]) &&
              group.items.map((item) => {
                const active = pathname === item.href || pathname?.startsWith(item.href + "/");
                const Icon = item.icon;
                return (
                  <Link
                    key={item.href}
                    href={item.href}
                    title={collapsed ? item.label : undefined}
                    data-testid="sidebar-item"
                    className={cn(
                      "flex items-center gap-3 rounded-[10px] transition-colors duration-150 mb-0.5",
                      collapsed ? "justify-center px-2 py-2" : "px-3 py-2",
                      active
                        ? "nav-active bg-[var(--accent-muted)] text-[var(--accent)] font-semibold"
                        : "text-foreground-secondary hover:bg-[var(--bg-tertiary)] hover:text-foreground",
                    )}
                  >
                    <Icon className="w-4 h-4 flex-shrink-0" style={{ width: 16, minWidth: 16 }} />
                    {!collapsed && (
                      <>
                        <span className="text-body-sm truncate flex-1">{item.label}</span>
                        {item.badge !== undefined && item.badge > 0 && (
                          <span data-testid="sidebar-badge" className="min-w-[18px] h-[18px] rounded-full bg-app-accent text-app-accent-t text-[9px] font-bold flex items-center justify-center px-1">
                            {item.badge}
                          </span>
                        )}
                      </>
                    )}
                  </Link>
                );
              })}
          </div>
        ))}
      </nav>

      <div className="border-t border-app-border p-3">
        {collapsed ? (
          <button
            onClick={toggleSidebar}
            className="w-full flex justify-center py-1 text-foreground-muted hover:text-foreground"
          >
            <PanelLeftOpen className="w-4 h-4" />
          </button>
        ) : (
          <div className="flex items-center justify-between">
            <div className="flex items-center gap-2 min-w-0">
              <div className="w-8 h-8 rounded-lg bg-gradient-to-br from-app-purple to-app-info flex items-center justify-center flex-shrink-0">
                <span className="text-white text-xs font-bold">JH</span>
              </div>
              <div className="min-w-0">
                <div className="text-body-sm font-semibold text-foreground truncate">
                  Jayesh H
                </div>
                <div className="text-tiny text-foreground-disabled truncate">
                  Admin · Complai
                </div>
              </div>
            </div>
            <button
              onClick={toggleSidebar}
              className="text-foreground-muted hover:text-foreground p-1"
            >
              <PanelLeftClose className="w-4 h-4" />
            </button>
          </div>
        )}
      </div>
    </aside>
  );
}
