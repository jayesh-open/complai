"use client";

import {
  FileText, Receipt, FileCheck2, AlertTriangle,
  TrendingUp, TrendingDown,
} from "lucide-react";
import { cn, formatINR, formatCompact } from "@/lib/utils";

interface MetricCardProps {
  icon: React.ReactNode;
  iconColor: string;
  value: string;
  label: string;
  subtitle?: string;
  trend?: { value: string; favorable: boolean };
}

function MetricCard({ icon, iconColor, value, label, subtitle, trend }: MetricCardProps) {
  return (
    <div
      data-testid="kpi-card"
      className={cn(
        "bg-app-card border border-app-border rounded-card p-4",
        "hover:border-app-accent transition-colors duration-150",
      )}
    >
      <div className="flex items-start justify-between mb-3">
        <div
          className="w-[30px] h-[30px] rounded-lg flex items-center justify-center border"
          style={{
            background: `color-mix(in srgb, ${iconColor} 15%, transparent)`,
            borderColor: `color-mix(in srgb, ${iconColor} 30%, transparent)`,
          }}
        >
          {icon}
        </div>
        {trend && (
          <span
            className={cn(
              "text-[11px] font-semibold flex items-center gap-0.5",
              trend.favorable ? "text-app-success" : "text-app-danger",
            )}
          >
            {trend.favorable ? <TrendingUp className="w-3 h-3" /> : <TrendingDown className="w-3 h-3" />}
            {trend.value}
          </span>
        )}
      </div>
      <div
        className="text-[20px] font-bold leading-tight text-foreground"
        style={{ fontFeatureSettings: '"tnum"' }}
      >
        {value}
      </div>
      <div className="text-overline text-foreground-muted mt-1">{label}</div>
      {subtitle && (
        <>
          <div className="border-t border-app-border-lt my-2" />
          <div className="text-xs text-foreground-muted">{subtitle}</div>
        </>
      )}
    </div>
  );
}

export default function DashboardPage() {
  return (
    <div data-testid="dashboard-page">
      <div className="mb-6">
        <h2 className="text-heading-xl text-foreground">Good morning, Jayesh</h2>
        <p className="text-body text-foreground-muted mt-1">
          Here&apos;s your compliance overview for April 2026
        </p>
      </div>

      <div className="grid grid-cols-1 sm:grid-cols-2 xl:grid-cols-4 gap-4 mb-6">
        <MetricCard
          icon={<FileText className="w-4 h-4" style={{ color: "var(--accent)" }} />}
          iconColor="var(--accent)"
          value="228"
          label="Total Invoices"
          subtitle={formatCompact(87600000)}
          trend={{ value: "12%", favorable: true }}
        />
        <MetricCard
          icon={<Receipt className="w-4 h-4" style={{ color: "var(--info)" }} />}
          iconColor="var(--info)"
          value="40"
          label="Pending Approval"
          subtitle={formatCompact(15200000)}
          trend={{ value: "8%", favorable: false }}
        />
        <MetricCard
          icon={<FileCheck2 className="w-4 h-4" style={{ color: "var(--success)" }} />}
          iconColor="var(--success)"
          value={formatINR(4520000)}
          label="ITC Claimed YTD"
          subtitle="Eligible: ₹52.8L"
          trend={{ value: "15%", favorable: true }}
        />
        <MetricCard
          icon={<AlertTriangle className="w-4 h-4" style={{ color: "var(--danger)" }} />}
          iconColor="var(--danger)"
          value="7"
          label="Overdue Filings"
          subtitle="GSTR-1: 3 · TDS: 4"
        />
      </div>

      <div className="grid grid-cols-1 lg:grid-cols-2 gap-4 mb-6">
        <div className="bg-app-card border border-app-border rounded-card p-5">
          <h3 className="text-heading-sm text-foreground mb-4">Compliance Health</h3>
          <div className="space-y-3">
            {[
              { module: "GST", status: "On Track", color: "success" as const, filedOnTime: "11/12" },
              { module: "TDS", status: "1 Overdue", color: "danger" as const, filedOnTime: "2/4" },
              { module: "E-Invoice", status: "Active", color: "success" as const, filedOnTime: "100%" },
              { module: "E-Way Bill", status: "Active", color: "success" as const, filedOnTime: "95%" },
            ].map((item) => (
              <div
                key={item.module}
                className="flex items-center justify-between py-2 px-3 rounded-lg bg-app-input"
              >
                <span className="text-body-sm font-semibold text-foreground">{item.module}</span>
                <div className="flex items-center gap-3">
                  <span className="text-caption text-foreground-muted">{item.filedOnTime}</span>
                  <span
                    className={cn(
                      "inline-flex items-center px-2.5 py-0.5 text-[10px] font-semibold uppercase tracking-wide rounded-badge border",
                      item.color === "success" && "bg-app-success-m text-app-success border-app-success-b",
                      item.color === "danger" && "bg-app-danger-m text-app-danger border-app-danger-b",
                    )}
                  >
                    {item.status}
                  </span>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-app-card border border-app-border rounded-card p-5">
          <h3 className="text-heading-sm text-foreground mb-4">Action Items</h3>
          <div className="space-y-2">
            {[
              { text: "File GSTR-3B for March 2026", due: "Due: 20/04/2026", urgent: true },
              { text: "Review 15 pending invoices", due: "From: Accounts team", urgent: false },
              { text: "Vendor compliance score review", due: "3 vendors below threshold", urgent: false },
              { text: "TDS Return Q4 — 24Q", due: "Due: 31/05/2026", urgent: false },
              { text: "GSTR-9 Annual Return FY 2025-26", due: "Due: 31/12/2026", urgent: false },
            ].map((item, i) => (
              <div
                key={i}
                className="flex items-start gap-3 py-2.5 px-3 rounded-lg hover:bg-app-input transition-colors cursor-pointer"
              >
                <div
                  className={cn(
                    "w-1.5 h-1.5 rounded-full mt-2 flex-shrink-0",
                    item.urgent ? "bg-app-danger" : "bg-app-accent",
                  )}
                />
                <div className="min-w-0">
                  <div className="text-body-sm font-medium text-foreground">{item.text}</div>
                  <div className="text-caption text-foreground-muted">{item.due}</div>
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>
    </div>
  );
}
