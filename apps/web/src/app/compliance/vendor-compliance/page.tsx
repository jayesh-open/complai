"use client";

import { useState, useCallback, useMemo } from "react";
import {
  ArrowLeft, RefreshCw, Loader2, ChevronRight,
  CheckCircle2, AlertTriangle, XCircle, Shield,
} from "lucide-react";
import Link from "next/link";
import { cn } from "@/lib/utils";
import { VendorComplianceScoreCard } from "@complai/ui-components";

type ViewMode = "list" | "detail";
type CategoryFilter = "all" | "A" | "B" | "C" | "D";

interface ScoreDimension {
  label: string;
  score: number;
  maxScore: number;
  status: "pass" | "warn" | "fail";
  note?: string;
}

interface VendorScore {
  id: string;
  name: string;
  legalName: string;
  gstin: string;
  pan: string;
  state: string;
  category: string;
  totalScore: number;
  riskLevel: string;
  lastReviewed: string;
  registrationStatus: string;
  msmeRegistered: boolean;
  vendorCategory: string;
  dimensions: ScoreDimension[];
}

function generateMockVendors(): VendorScore[] {
  const states: Record<string, string> = {
    "29": "Karnataka", "27": "Maharashtra", "33": "Tamil Nadu",
    "07": "Delhi", "24": "Gujarat", "09": "Uttar Pradesh",
    "19": "West Bengal", "08": "Rajasthan", "37": "Andhra Pradesh",
    "36": "Telangana",
  };
  const stateCodes = Object.keys(states);
  const categories = ["Manufacturer", "Service Provider", "Trader", "Import Supplier", "Logistics"];

  const vendors: VendorScore[] = [];

  const exemplary = [
    { name: "Tata Steel Ltd", pan: "AABCT1234A" },
    { name: "Infosys Ltd", pan: "AABCI5678B" },
    { name: "Reliance Industries Ltd", pan: "AABCR9012C" },
    { name: "Wipro Ltd", pan: "AABCW3456D" },
    { name: "Hindustan Unilever Ltd", pan: "AABCH7890E" },
    { name: "Larsen & Toubro Ltd", pan: "AABCL1234F" },
    { name: "Bajaj Auto Ltd", pan: "AABCB5678G" },
    { name: "Asian Paints Ltd", pan: "AABCA9012H" },
    { name: "Maruti Suzuki India Ltd", pan: "AABCM3456I" },
    { name: "Tech Mahindra Ltd", pan: "AABCT7890J" },
  ];
  exemplary.forEach((v, i) => {
    const sc = stateCodes[i % stateCodes.length];
    vendors.push({
      id: `v-${i + 1}`,
      name: v.name,
      legalName: v.name,
      gstin: `${sc}${v.pan}1Z5`,
      pan: v.pan,
      state: states[sc],
      category: "A",
      totalScore: 92 + (i % 8),
      riskLevel: "Low",
      lastReviewed: "25/04/2026",
      registrationStatus: "Active",
      msmeRegistered: i % 3 === 0,
      vendorCategory: categories[i % categories.length],
      dimensions: [
        { label: "Filing regularity", score: 28 + (i % 3), maxScore: 30, status: "pass" },
        { label: "IRN compliance", score: 19 + (i % 2), maxScore: 20, status: "pass" },
        { label: "Mismatch rate", score: 18 + (i % 3), maxScore: 20, status: "pass" },
        { label: "Payment behaviour", score: 13 + (i % 3), maxScore: 15, status: "pass" },
        { label: "Document hygiene", score: 13 + (i % 3), maxScore: 15, status: "pass" },
      ],
    });
  });

  const good = [
    "QuickParts Industries Pvt Ltd", "GlobalTech Solutions LLP", "Metro Logistics Pvt Ltd",
    "Diamond Trading Co.", "Sunrise Pharma Pvt Ltd", "National Steel Corp",
    "Pacific Imports Pvt Ltd", "Green Energy Solutions", "Bharat Electronics Pvt Ltd",
    "United Services Ltd", "Mahindra Components", "Star Packaging Pvt Ltd",
    "Delta Systems LLP", "Premier Chemicals Ltd", "Atlas Engineering Pvt Ltd",
  ];
  good.forEach((name, i) => {
    const idx = i + 10;
    const sc = stateCodes[idx % stateCodes.length];
    const pan = `AABCG${String(idx).padStart(4, "0")}K`;
    const score = 62 + (i % 18);
    vendors.push({
      id: `v-${idx + 1}`,
      name,
      legalName: name,
      gstin: `${sc}${pan}1Z5`,
      pan,
      state: states[sc],
      category: "B",
      totalScore: score,
      riskLevel: score >= 80 ? "Low" : "Medium",
      lastReviewed: "25/04/2026",
      registrationStatus: "Active",
      msmeRegistered: i % 4 === 0,
      vendorCategory: categories[idx % categories.length],
      dimensions: [
        { label: "Filing regularity", score: 20 + (i % 8), maxScore: 30, status: i % 3 === 0 ? "warn" : "pass", note: i % 3 === 0 ? "1-2 late filings" : undefined },
        { label: "IRN compliance", score: 14 + (i % 6), maxScore: 20, status: "pass" },
        { label: "Mismatch rate", score: 14 + (i % 6), maxScore: 20, status: i % 4 === 0 ? "warn" : "pass", note: i % 4 === 0 ? "2-5% mismatch" : undefined },
        { label: "Payment behaviour", score: 10 + (i % 5), maxScore: 15, status: "pass" },
        { label: "Document hygiene", score: 10 + (i % 5), maxScore: 15, status: i % 5 === 0 ? "warn" : "pass" },
      ],
    });
  });

  const average = [
    "Sagar Textiles OPC", "Patel Construction Co.", "Sharma Electronics",
    "Raj Kumar Traders", "Kavitha Imports LLP", "Balaji Steels",
    "Gowda Logistics", "Nair Pharma Pvt Ltd", "Joshi Engineering",
    "Chopra Plastics", "Verma Trading Co.", "Gupta Chemicals",
    "Reddy Enterprises", "Iyer Services", "Menon Exports Pvt Ltd",
  ];
  average.forEach((name, i) => {
    const idx = i + 25;
    const sc = stateCodes[idx % stateCodes.length];
    const pan = `AABCC${String(idx).padStart(4, "0")}L`;
    vendors.push({
      id: `v-${idx + 1}`,
      name,
      legalName: name,
      gstin: `${sc}${pan}1Z5`,
      pan,
      state: states[sc],
      category: "C",
      totalScore: 42 + (i % 18),
      riskLevel: "High",
      lastReviewed: "25/04/2026",
      registrationStatus: i % 5 === 0 ? "Suspended" : "Active",
      msmeRegistered: false,
      vendorCategory: categories[idx % categories.length],
      dimensions: [
        { label: "Filing regularity", score: 12 + (i % 8), maxScore: 30, status: "warn", note: "3-5 late filings" },
        { label: "IRN compliance", score: 8 + (i % 8), maxScore: 20, status: "warn", note: "Partial IRN adoption" },
        { label: "Mismatch rate", score: 8 + (i % 8), maxScore: 20, status: "warn", note: "5-15% mismatch rate" },
        { label: "Payment behaviour", score: 6 + (i % 6), maxScore: 15, status: "warn", note: "Occasional delays" },
        { label: "Document hygiene", score: 6 + (i % 6), maxScore: 15, status: "fail", note: "Missing GST cert" },
      ],
    });
  });

  const poor = [
    "ABC Traders", "XYZ Enterprises", "Midnight Logistics",
    "Shadow Imports", "Ghost Corp OPC", "Fake Pharma LLP",
    "Dubious Trading Co.", "Phantom Services", "Hollow Steel Corp",
    "Shell Chemicals Pvt Ltd",
  ];
  poor.forEach((name, i) => {
    const idx = i + 40;
    const sc = stateCodes[idx % stateCodes.length];
    const pan = `AABCD${String(idx).padStart(4, "0")}M`;
    vendors.push({
      id: `v-${idx + 1}`,
      name,
      legalName: name,
      gstin: `${sc}${pan}1Z5`,
      pan,
      state: states[sc],
      category: "D",
      totalScore: 15 + (i % 24),
      riskLevel: "Critical",
      lastReviewed: "25/04/2026",
      registrationStatus: i % 3 === 0 ? "Cancelled" : "Active",
      msmeRegistered: false,
      vendorCategory: categories[idx % categories.length],
      dimensions: [
        { label: "Filing regularity", score: 3 + (i % 6), maxScore: 30, status: "fail", note: "Frequent missed filings" },
        { label: "IRN compliance", score: i % 4, maxScore: 20, status: "fail", note: "No IRN compliance" },
        { label: "Mismatch rate", score: 3 + (i % 5), maxScore: 20, status: "fail", note: "20-35% mismatch" },
        { label: "Payment behaviour", score: 2 + (i % 5), maxScore: 15, status: "fail", note: "Chronic overdue" },
        { label: "Document hygiene", score: 2 + (i % 5), maxScore: 15, status: "fail", note: "Missing GSTIN cert, PAN unverified" },
      ],
    });
  });

  return vendors;
}

const MOCK_VENDORS = generateMockVendors();

function CategoryBadge({ category }: { category: string }) {
  const config: Record<string, { bg: string; text: string; label: string }> = {
    A: { bg: "bg-[var(--success-muted)]", text: "text-[var(--success)]", label: "A — Exemplary" },
    B: { bg: "bg-[var(--accent-muted)]", text: "text-[var(--accent)]", label: "B — Good" },
    C: { bg: "bg-[var(--warning-muted)]", text: "text-[var(--warning)]", label: "C — Moderate" },
    D: { bg: "bg-[var(--danger-muted)]", text: "text-[var(--danger)]", label: "D — Poor" },
  };
  const c = config[category] ?? config.D;
  return (
    <span data-testid="category-badge" className={cn("px-2 py-0.5 rounded-full text-[10px] font-bold", c.bg, c.text)}>
      {c.label}
    </span>
  );
}

function RiskBadge({ level }: { level: string }) {
  const config: Record<string, string> = {
    Low: "text-[var(--success)]",
    Medium: "text-[var(--warning)]",
    High: "text-[var(--danger)]",
    Critical: "text-[var(--danger)]",
  };
  return (
    <span className={cn("text-[10px] font-semibold", config[level] ?? config.High)}>
      {level}
    </span>
  );
}

export default function VendorCompliancePage() {
  const [view, setView] = useState<ViewMode>("list");
  const [selectedVendor, setSelectedVendor] = useState<VendorScore | null>(null);
  const [filter, setFilter] = useState<CategoryFilter>("all");
  const [syncing, setSyncing] = useState(false);
  const [searchQuery, setSearchQuery] = useState("");

  const filteredVendors = useMemo(() => {
    let result = MOCK_VENDORS;
    if (filter !== "all") {
      result = result.filter((v) => v.category === filter);
    }
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      result = result.filter(
        (v) => v.name.toLowerCase().includes(q) || v.gstin.toLowerCase().includes(q) || v.pan.toLowerCase().includes(q),
      );
    }
    return result;
  }, [filter, searchQuery]);

  const summary = useMemo(() => {
    const s = { total: MOCK_VENDORS.length, catA: 0, catB: 0, catC: 0, catD: 0, avgScore: 0 };
    let totalScore = 0;
    for (const v of MOCK_VENDORS) {
      totalScore += v.totalScore;
      if (v.category === "A") s.catA++;
      else if (v.category === "B") s.catB++;
      else if (v.category === "C") s.catC++;
      else s.catD++;
    }
    s.avgScore = Math.round(totalScore / MOCK_VENDORS.length);
    return s;
  }, []);

  const handleSync = useCallback(async () => {
    setSyncing(true);
    await new Promise((r) => setTimeout(r, 1500));
    setSyncing(false);
  }, []);

  const handleVendorClick = useCallback((vendor: VendorScore) => {
    setSelectedVendor(vendor);
    setView("detail");
  }, []);

  if (view === "detail" && selectedVendor) {
    return (
      <div className="space-y-6" data-testid="vendor-detail-view">
        <div className="flex items-center gap-3">
          <button
            onClick={() => { setView("list"); setSelectedVendor(null); }}
            className="text-foreground-muted hover:text-foreground transition-colors"
            data-testid="back-to-list"
          >
            <ArrowLeft className="w-5 h-5" />
          </button>
          <div>
            <h1 className="text-heading-xl text-foreground">{selectedVendor.name}</h1>
            <p className="text-body-sm text-foreground-muted">
              {selectedVendor.gstin} &middot; {selectedVendor.state}
            </p>
          </div>
          <CategoryBadge category={selectedVendor.category} />
        </div>

        <div className="grid grid-cols-4 gap-4">
          {[
            { label: "Compliance Score", value: `${selectedVendor.totalScore}/100`, icon: Shield },
            { label: "Category", value: selectedVendor.category, icon: CheckCircle2 },
            { label: "Risk Level", value: selectedVendor.riskLevel, icon: AlertTriangle },
            { label: "Registration", value: selectedVendor.registrationStatus, icon: CheckCircle2 },
          ].map((kpi) => (
            <div key={kpi.label} className="bg-app-card border border-app-border rounded-card p-4">
              <div className="flex items-center gap-2 text-xs text-foreground-muted mb-1">
                <kpi.icon className="w-3.5 h-3.5" />
                {kpi.label}
              </div>
              <div className="text-heading-lg text-foreground">{kpi.value}</div>
            </div>
          ))}
        </div>

        <VendorComplianceScoreCard
          vendorName={selectedVendor.name}
          gstin={selectedVendor.gstin}
          state={selectedVendor.state}
          totalScore={selectedVendor.totalScore}
          riskLevel={selectedVendor.riskLevel}
          category={selectedVendor.category}
          lastReviewed={selectedVendor.lastReviewed}
          dimensions={selectedVendor.dimensions}
        />

        <div className="bg-app-card border border-app-border rounded-card p-5">
          <h3 className="text-heading-md text-foreground mb-3">Vendor Details</h3>
          <div className="grid grid-cols-2 gap-3">
            {[
              { label: "Legal Name", value: selectedVendor.legalName },
              { label: "PAN", value: selectedVendor.pan },
              { label: "GSTIN", value: selectedVendor.gstin },
              { label: "State", value: selectedVendor.state },
              { label: "Category", value: selectedVendor.vendorCategory },
              { label: "MSME Registered", value: selectedVendor.msmeRegistered ? "Yes" : "No" },
              { label: "Registration Status", value: selectedVendor.registrationStatus },
              { label: "Last Reviewed", value: selectedVendor.lastReviewed },
            ].map((d) => (
              <div key={d.label} className="flex gap-2 text-xs">
                <span className="text-foreground-muted min-w-[130px]">{d.label}:</span>
                <span className="text-foreground font-medium font-mono">{d.value}</span>
              </div>
            ))}
          </div>
        </div>
      </div>
    );
  }

  return (
    <div className="space-y-6" data-testid="vendor-compliance-list">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/compliance/gst" className="text-foreground-muted hover:text-foreground transition-colors">
            <ArrowLeft className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="text-heading-xl text-foreground">Vendor Compliance</h1>
            <p className="text-body-sm text-foreground-muted">
              100-point compliance scoring across 5 dimensions. Data synced from Apex P2P.
            </p>
          </div>
        </div>
        <button
          data-testid="sync-button"
          onClick={handleSync}
          disabled={syncing}
          className={cn(
            "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
            "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
            syncing && "opacity-50 cursor-not-allowed",
          )}
        >
          {syncing ? <Loader2 className="w-4 h-4 animate-spin" /> : <RefreshCw className="w-4 h-4" />}
          {syncing ? "Syncing..." : "Sync from Apex"}
        </button>
      </div>

      <div className="grid grid-cols-5 gap-4" data-testid="kpi-summary">
        <div className="bg-app-card border border-app-border rounded-card p-4">
          <div className="text-xs text-foreground-muted mb-1">Total Vendors</div>
          <div className="text-heading-lg text-foreground">{summary.total}</div>
        </div>
        <div className="bg-app-card border border-app-border rounded-card p-4">
          <div className="text-xs text-foreground-muted mb-1">Avg Score</div>
          <div className="text-heading-lg text-foreground">{summary.avgScore}/100</div>
        </div>
        <div className="bg-app-card border border-app-border rounded-card p-4">
          <div className="flex items-center gap-1 text-xs text-[var(--success)] mb-1">
            <CheckCircle2 className="w-3 h-3" /> Cat A + B
          </div>
          <div className="text-heading-lg text-foreground">{summary.catA + summary.catB}</div>
        </div>
        <div className="bg-app-card border border-app-border rounded-card p-4">
          <div className="flex items-center gap-1 text-xs text-[var(--warning)] mb-1">
            <AlertTriangle className="w-3 h-3" /> Cat C
          </div>
          <div className="text-heading-lg text-foreground">{summary.catC}</div>
        </div>
        <div className="bg-app-card border border-app-border rounded-card p-4">
          <div className="flex items-center gap-1 text-xs text-[var(--danger)] mb-1">
            <XCircle className="w-3 h-3" /> Cat D
          </div>
          <div className="text-heading-lg text-foreground">{summary.catD}</div>
        </div>
      </div>

      <div className="bg-app-card border border-app-border rounded-card">
        <div className="px-4 py-3 border-b border-app-border flex items-center gap-3">
          <div className="flex gap-1">
            {(["all", "A", "B", "C", "D"] as CategoryFilter[]).map((cat) => (
              <button
                key={cat}
                data-testid={`filter-${cat}`}
                onClick={() => setFilter(cat)}
                className={cn(
                  "px-3 py-1 rounded-lg text-[11px] font-semibold transition-colors",
                  filter === cat
                    ? "bg-[var(--accent-muted)] text-[var(--accent)]"
                    : "text-foreground-muted hover:bg-[var(--bg-tertiary)]",
                )}
              >
                {cat === "all" ? "All" : `Cat ${cat}`}
              </button>
            ))}
          </div>
          <input
            type="text"
            placeholder="Search vendors..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            data-testid="vendor-search"
            className="ml-auto px-3 py-1.5 rounded-lg text-xs border border-app-border bg-[var(--bg-primary)] text-foreground placeholder:text-foreground-disabled focus:outline-none focus:border-[var(--accent)] w-64"
          />
          <span className="text-[11px] text-foreground-muted">{filteredVendors.length} vendors</span>
        </div>

        <div data-testid="vendor-table" className="overflow-hidden">
          <table className="w-full text-xs">
            <thead>
              <tr className="bg-[var(--bg-tertiary)] border-b border-[var(--border-default)]">
                <th className="text-left px-4 py-2.5 text-[var(--text-muted)] font-semibold">Vendor</th>
                <th className="text-left px-4 py-2.5 text-[var(--text-muted)] font-semibold">GSTIN</th>
                <th className="text-left px-4 py-2.5 text-[var(--text-muted)] font-semibold">State</th>
                <th className="text-center px-4 py-2.5 text-[var(--text-muted)] font-semibold">Score</th>
                <th className="text-center px-4 py-2.5 text-[var(--text-muted)] font-semibold">Category</th>
                <th className="text-center px-4 py-2.5 text-[var(--text-muted)] font-semibold">Risk</th>
                <th className="text-center px-4 py-2.5 text-[var(--text-muted)] font-semibold">Status</th>
                <th className="w-8"></th>
              </tr>
            </thead>
            <tbody>
              {filteredVendors.map((vendor) => (
                <tr
                  key={vendor.id}
                  data-testid="vendor-row"
                  onClick={() => handleVendorClick(vendor)}
                  className="border-b border-[var(--border-default)] last:border-0 hover:bg-[var(--bg-tertiary)] transition-colors cursor-pointer"
                >
                  <td className="px-4 py-2.5">
                    <div className="font-medium text-foreground">{vendor.name}</div>
                    <div className="text-[10px] text-foreground-muted">{vendor.vendorCategory}</div>
                  </td>
                  <td className="px-4 py-2.5 font-mono text-foreground">{vendor.gstin}</td>
                  <td className="px-4 py-2.5 text-foreground">{vendor.state}</td>
                  <td className="px-4 py-2.5 text-center">
                    <span
                      className="font-bold"
                      style={{
                        color: vendor.totalScore >= 80 ? "var(--success)" : vendor.totalScore >= 60 ? "var(--warning)" : "var(--danger)",
                      }}
                    >
                      {vendor.totalScore}
                    </span>
                  </td>
                  <td className="px-4 py-2.5 text-center">
                    <CategoryBadge category={vendor.category} />
                  </td>
                  <td className="px-4 py-2.5 text-center">
                    <RiskBadge level={vendor.riskLevel} />
                  </td>
                  <td className="px-4 py-2.5 text-center">
                    <span className={cn(
                      "text-[10px] font-semibold",
                      vendor.registrationStatus === "Active" ? "text-[var(--success)]" : "text-[var(--danger)]",
                    )}>
                      {vendor.registrationStatus}
                    </span>
                  </td>
                  <td className="px-4 py-2.5">
                    <ChevronRight className="w-3.5 h-3.5 text-foreground-muted" />
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
