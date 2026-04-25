"use client";

import { useEffect, useCallback } from "react";
import { useRouter } from "next/navigation";
import { Command } from "cmdk";
import {
  LayoutDashboard, FileText, Building2, FileSpreadsheet,
  Receipt, BarChart3, Settings, Plus, History, Search,
  HelpCircle, FileCheck2, Truck,
} from "lucide-react";
import { cn } from "@/lib/utils";

interface CommandPaletteProps {
  open: boolean;
  onOpenChange: (open: boolean) => void;
}

interface CommandItem {
  id: string;
  label: string;
  icon: React.ElementType;
  href?: string;
  category: string;
  shortcut?: string;
}

const COMMANDS: CommandItem[] = [
  { id: "nav-dashboard", label: "Dashboard", icon: LayoutDashboard, href: "/dashboard", category: "Navigate", shortcut: "g d" },
  { id: "nav-invoices", label: "Invoices", icon: FileText, href: "/payables/invoices", category: "Navigate", shortcut: "g i" },
  { id: "nav-gst", label: "GST", icon: FileSpreadsheet, href: "/compliance/gst", category: "Navigate", shortcut: "g g" },
  { id: "nav-tds", label: "TDS", icon: Receipt, href: "/compliance/tds", category: "Navigate", shortcut: "g t" },
  { id: "nav-einvoice", label: "E-Invoicing", icon: FileCheck2, href: "/compliance/e-invoice", category: "Navigate" },
  { id: "nav-ewb", label: "E-Way Bill", icon: Truck, href: "/compliance/e-way-bill", category: "Navigate" },
  { id: "nav-vendors", label: "Vendors", icon: Building2, href: "/payables/vendors", category: "Navigate" },
  { id: "nav-reports", label: "Reports & Analytics", icon: BarChart3, href: "/insights/reports", category: "Navigate" },
  { id: "nav-settings", label: "Settings", icon: Settings, href: "/configure/settings", category: "Navigate" },
  { id: "nav-appearance", label: "Appearance", icon: Settings, href: "/configure/appearance", category: "Navigate" },
  { id: "create-invoice", label: "New Invoice", icon: Plus, href: "/payables/invoices/new", category: "Create", shortcut: "n i" },
  { id: "create-po", label: "New Purchase Order", icon: Plus, href: "/procurement/purchase-orders/new", category: "Create", shortcut: "n p" },
  { id: "create-vendor", label: "New Vendor", icon: Plus, href: "/payables/vendors/new", category: "Create" },
  { id: "filing-gstr1", label: "File GSTR-1", icon: FileSpreadsheet, href: "/compliance/gst/gstr-1", category: "Filings" },
  { id: "filing-gstr3b", label: "File GSTR-3B", icon: FileSpreadsheet, href: "/compliance/gst/gstr-3b", category: "Filings" },
  { id: "report-gst-summary", label: "GST Summary Report", icon: BarChart3, href: "/insights/reports/gst-summary", category: "Reports" },
  { id: "recent-placeholder", label: "Recent items will appear here", icon: History, category: "Recent" },
  { id: "help-docs", label: "Documentation", icon: HelpCircle, category: "Help" },
  { id: "help-shortcuts", label: "Keyboard Shortcuts", icon: HelpCircle, category: "Help" },
];

const CATEGORIES = ["Navigate", "Create", "Recent", "Filings", "Reports", "Help"];

export function CommandPalette({ open, onOpenChange }: CommandPaletteProps) {
  const router = useRouter();

  const handleSelect = useCallback(
    (item: CommandItem) => {
      if (item.href) {
        router.push(item.href);
      }
      onOpenChange(false);
    },
    [router, onOpenChange],
  );

  useEffect(() => {
    const down = (e: KeyboardEvent) => {
      if (e.key === "k" && (e.metaKey || e.ctrlKey)) {
        e.preventDefault();
        onOpenChange(!open);
      }
      if (e.key === "Escape") {
        onOpenChange(false);
      }
    };
    document.addEventListener("keydown", down);
    return () => document.removeEventListener("keydown", down);
  }, [open, onOpenChange]);

  if (!open) return null;

  return (
    <div className="fixed inset-0 z-50" data-testid="command-palette">
      <div
        className="absolute inset-0 bg-[var(--bg-overlay)] backdrop-blur-sm"
        onClick={() => onOpenChange(false)}
      />
      <div className="absolute top-[20%] left-1/2 -translate-x-1/2 w-full max-w-lg mx-4">
        <Command
          className={cn(
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "rounded-2xl shadow-[var(--shadow-lg)] overflow-hidden",
          )}
          label="Command palette"
        >
          <div className="flex items-center gap-2 px-4 border-b border-[var(--border-default)]">
            <Search className="w-4 h-4 text-foreground-muted flex-shrink-0" />
            <Command.Input
              placeholder="Type a command or search..."
              className="w-full py-3 text-sm bg-transparent text-foreground outline-none placeholder:text-foreground-muted"
              autoFocus
            />
          </div>
          <Command.List className="max-h-[320px] overflow-y-auto p-2">
            <Command.Empty className="py-8 text-center text-sm text-foreground-muted">
              No results found.
            </Command.Empty>
            {CATEGORIES.map((category) => {
              const items = COMMANDS.filter((c) => c.category === category);
              if (items.length === 0) return null;
              return (
                <Command.Group
                  key={category}
                  heading={category}
                  className="[&_[cmdk-group-heading]]:text-overline [&_[cmdk-group-heading]]:text-foreground-muted [&_[cmdk-group-heading]]:px-2 [&_[cmdk-group-heading]]:py-1.5"
                >
                  {items.map((item) => {
                    const Icon = item.icon;
                    return (
                      <Command.Item
                        key={item.id}
                        value={`${item.label} ${item.category}`}
                        onSelect={() => handleSelect(item)}
                        className={cn(
                          "flex items-center gap-3 px-3 py-2 rounded-lg cursor-pointer",
                          "text-foreground-secondary text-xs",
                          "data-[selected=true]:bg-[var(--accent-muted)] data-[selected=true]:text-[var(--accent)]",
                        )}
                      >
                        <Icon className="w-4 h-4 flex-shrink-0" />
                        <span className="flex-1">{item.label}</span>
                        {item.shortcut && (
                          <kbd className="text-tiny text-foreground-disabled bg-app-input px-1.5 py-0.5 rounded border border-app-border">
                            {item.shortcut}
                          </kbd>
                        )}
                      </Command.Item>
                    );
                  })}
                </Command.Group>
              );
            })}
          </Command.List>
        </Command>
      </div>
    </div>
  );
}
