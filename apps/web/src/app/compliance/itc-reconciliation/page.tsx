"use client";

import { useState, useCallback, useMemo } from "react";
import {
  ArrowLeft, RefreshCw, Download, CheckCircle2, XCircle,
  AlertTriangle, Loader2, Sparkles, Search, Filter,
  ChevronDown, ChevronRight, Link2, SkipForward, Check,
  ArrowUpDown, Mail, FileSpreadsheet, Eye,
} from "lucide-react";
import Link from "next/link";
import { cn, formatINR } from "@/lib/utils";

type MatchBucket = "matched" | "mismatch" | "partial" | "missing_2b" | "missing_pr" | "duplicate";
type IMSAction = "accept" | "reject" | "pending";

interface PREntry {
  id: string;
  supplierGstin: string;
  supplierName: string;
  invoiceNumber: string;
  invoiceDate: string;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  totalValue: number;
  hsn: string;
}

interface GSTR2BEntry {
  id: string;
  supplierGstin: string;
  supplierName: string;
  invoiceNumber: string;
  invoiceDate: string;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  totalValue: number;
  hsn: string;
  imsAction: IMSAction;
}

interface ReconMatch {
  id: string;
  bucket: MatchBucket;
  confidence: number;
  prEntry: PREntry | null;
  gstr2bEntry: GSTR2BEntry | null;
  mismatchReasons: string[];
  aiSuggestion: boolean;
  accepted: boolean | null;
}

interface BucketCount {
  bucket: MatchBucket;
  label: string;
  count: number;
  color: string;
  icon: React.ReactNode;
}

const MOCK_SUPPLIERS = [
  { gstin: "29AABCA1234A1Z5", name: "Steel Corp India Pvt Ltd" },
  { gstin: "27BBBFB5678B2Y4", name: "Precision Auto Components" },
  { gstin: "33CCCGC9012C3X3", name: "Chennai Electronics Ltd" },
  { gstin: "07DDDID3456D4W2", name: "Delhi Packaging Solutions" },
  { gstin: "24EEEJE7890E5V1", name: "Gujarat Chemicals Co" },
  { gstin: "36FFFKF1234F6U0", name: "Telangana Textiles" },
  { gstin: "06GGGLG5678G7T9", name: "Haryana Machine Tools" },
  { gstin: "19HHHNH9012H8S8", name: "West Bengal Paper Mills" },
  { gstin: "32IIIMI3456I9R7", name: "Kerala Spices Export" },
  { gstin: "21JJJOJ7890J0Q6", name: "Odisha Mining Corp" },
];

function generateMockData(): ReconMatch[] {
  const matches: ReconMatch[] = [];

  for (let i = 0; i < 25; i++) {
    const supplier = MOCK_SUPPLIERS[i % MOCK_SUPPLIERS.length];
    const invNum = `INV-2026-${String(i + 1).padStart(4, "0")}`;
    const taxable = Math.round((50000 + Math.random() * 450000) * 100) / 100;
    const rate = [5, 12, 18, 28][Math.floor(Math.random() * 4)] / 100;
    const isInter = Math.random() > 0.5;
    const tax = Math.round(taxable * rate * 100) / 100;

    matches.push({
      id: `match-${i}`,
      bucket: "matched",
      confidence: 100,
      prEntry: {
        id: `pr-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: invNum, invoiceDate: "15/04/2026",
        taxableValue: taxable, cgst: isInter ? 0 : tax / 2, sgst: isInter ? 0 : tax / 2,
        igst: isInter ? tax : 0, totalValue: taxable + tax, hsn: "7208",
      },
      gstr2bEntry: {
        id: `2b-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: invNum, invoiceDate: "15/04/2026",
        taxableValue: taxable, cgst: isInter ? 0 : tax / 2, sgst: isInter ? 0 : tax / 2,
        igst: isInter ? tax : 0, totalValue: taxable + tax, hsn: "7208", imsAction: "accept",
      },
      mismatchReasons: [],
      aiSuggestion: false,
      accepted: null,
    });
  }

  for (let i = 0; i < 8; i++) {
    const supplier = MOCK_SUPPLIERS[(i + 3) % MOCK_SUPPLIERS.length];
    const invNum = `INV-2026-${String(i + 100).padStart(4, "0")}`;
    const prTaxable = Math.round((80000 + Math.random() * 200000) * 100) / 100;
    const diff = Math.round((500 + Math.random() * 3000) * 100) / 100;
    const rate = 0.18;
    const isInter = i % 2 === 0;
    const prTax = Math.round(prTaxable * rate * 100) / 100;
    const g2bTax = Math.round((prTaxable + diff) * rate * 100) / 100;
    const reasons: string[] = [];
    if (i < 4) reasons.push("Amount mismatch");
    if (i >= 4 && i < 6) reasons.push("Date mismatch");
    if (i >= 6) reasons.push("Amount mismatch", "Date mismatch");

    matches.push({
      id: `mismatch-${i}`,
      bucket: "mismatch",
      confidence: 60 + Math.floor(Math.random() * 30),
      prEntry: {
        id: `pr-m-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: invNum, invoiceDate: "15/04/2026",
        taxableValue: prTaxable, cgst: isInter ? 0 : prTax / 2, sgst: isInter ? 0 : prTax / 2,
        igst: isInter ? prTax : 0, totalValue: prTaxable + prTax, hsn: "8483",
      },
      gstr2bEntry: {
        id: `2b-m-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: invNum, invoiceDate: i >= 4 ? "17/04/2026" : "15/04/2026",
        taxableValue: prTaxable + diff, cgst: isInter ? 0 : g2bTax / 2, sgst: isInter ? 0 : g2bTax / 2,
        igst: isInter ? g2bTax : 0, totalValue: prTaxable + diff + g2bTax, hsn: "8483", imsAction: "pending",
      },
      mismatchReasons: reasons,
      aiSuggestion: i < 3,
      accepted: null,
    });
  }

  for (let i = 0; i < 3; i++) {
    const supplier = MOCK_SUPPLIERS[(i + 7) % MOCK_SUPPLIERS.length];
    matches.push({
      id: `partial-${i}`,
      bucket: "partial",
      confidence: 45 + Math.floor(Math.random() * 20),
      prEntry: {
        id: `pr-p-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `INV-2026-${String(i + 200).padStart(4, "0")}`, invoiceDate: "20/04/2026",
        taxableValue: 350000, cgst: 31500, sgst: 31500, igst: 0, totalValue: 413000, hsn: "3926",
      },
      gstr2bEntry: {
        id: `2b-p-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `INV-2026-${String(i + 200).padStart(4, "0")}A`, invoiceDate: "20/04/2026",
        taxableValue: 200000, cgst: 18000, sgst: 18000, igst: 0, totalValue: 236000, hsn: "3926",
        imsAction: "pending",
      },
      mismatchReasons: ["Partial amount — possible split invoice"],
      aiSuggestion: true,
      accepted: null,
    });
  }

  for (let i = 0; i < 5; i++) {
    const supplier = MOCK_SUPPLIERS[i % MOCK_SUPPLIERS.length];
    matches.push({
      id: `miss2b-${i}`,
      bucket: "missing_2b",
      confidence: 0,
      prEntry: {
        id: `pr-x-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `INV-2026-${String(i + 300).padStart(4, "0")}`, invoiceDate: "25/04/2026",
        taxableValue: 120000 + i * 30000, cgst: 0, sgst: 0,
        igst: Math.round((120000 + i * 30000) * 0.18 * 100) / 100,
        totalValue: Math.round((120000 + i * 30000) * 1.18 * 100) / 100, hsn: "8471",
      },
      gstr2bEntry: null,
      mismatchReasons: ["Not found in GSTR-2B"],
      aiSuggestion: false,
      accepted: null,
    });
  }

  for (let i = 0; i < 4; i++) {
    const supplier = MOCK_SUPPLIERS[(i + 5) % MOCK_SUPPLIERS.length];
    matches.push({
      id: `misspr-${i}`,
      bucket: "missing_pr",
      confidence: 0,
      prEntry: null,
      gstr2bEntry: {
        id: `2b-y-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `SUP-2026-${String(i + 400).padStart(4, "0")}`, invoiceDate: "22/04/2026",
        taxableValue: 95000 + i * 20000, cgst: 8550 + i * 1800, sgst: 8550 + i * 1800,
        igst: 0, totalValue: 112100 + i * 23600, hsn: "4811", imsAction: "pending",
      },
      mismatchReasons: ["Not found in Purchase Register"],
      aiSuggestion: false,
      accepted: null,
    });
  }

  for (let i = 0; i < 2; i++) {
    const supplier = MOCK_SUPPLIERS[i];
    matches.push({
      id: `dup-${i}`,
      bucket: "duplicate",
      confidence: 95,
      prEntry: {
        id: `pr-d-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `INV-2026-${String(i + 1).padStart(4, "0")}`, invoiceDate: "15/04/2026",
        taxableValue: 250000, cgst: 22500, sgst: 22500, igst: 0, totalValue: 295000, hsn: "7208",
      },
      gstr2bEntry: {
        id: `2b-d-${i}`, supplierGstin: supplier.gstin, supplierName: supplier.name,
        invoiceNumber: `INV-2026-${String(i + 1).padStart(4, "0")}`, invoiceDate: "15/04/2026",
        taxableValue: 250000, cgst: 22500, sgst: 22500, igst: 0, totalValue: 295000, hsn: "7208",
        imsAction: "accept",
      },
      mismatchReasons: ["Duplicate invoice detected"],
      aiSuggestion: false,
      accepted: null,
    });
  }

  return matches;
}

export default function ITCReconciliationWorkspace() {
  const [matches, setMatches] = useState<ReconMatch[]>(() => generateMockData());
  const [activeBucket, setActiveBucket] = useState<MatchBucket | "all">("all");
  const [searchQuery, setSearchQuery] = useState("");
  const [expandedRow, setExpandedRow] = useState<string | null>(null);
  const [loading, setLoading] = useState(false);
  const [reconRunning, setReconRunning] = useState(false);
  const [showAiOnly, setShowAiOnly] = useState(false);
  const [reasonFilter, setReasonFilter] = useState<string>("all");

  const bucketCounts: BucketCount[] = useMemo(() => [
    { bucket: "matched", label: "Matched", count: matches.filter((m) => m.bucket === "matched").length, color: "var(--success)", icon: <CheckCircle2 className="w-3.5 h-3.5" /> },
    { bucket: "mismatch", label: "Mismatch", count: matches.filter((m) => m.bucket === "mismatch").length, color: "var(--warning)", icon: <AlertTriangle className="w-3.5 h-3.5" /> },
    { bucket: "partial", label: "Partial", count: matches.filter((m) => m.bucket === "partial").length, color: "var(--accent)", icon: <Link2 className="w-3.5 h-3.5" /> },
    { bucket: "missing_2b", label: "Missing in 2B", count: matches.filter((m) => m.bucket === "missing_2b").length, color: "var(--danger)", icon: <XCircle className="w-3.5 h-3.5" /> },
    { bucket: "missing_pr", label: "Missing in PR", count: matches.filter((m) => m.bucket === "missing_pr").length, color: "var(--danger)", icon: <XCircle className="w-3.5 h-3.5" /> },
    { bucket: "duplicate", label: "Duplicate", count: matches.filter((m) => m.bucket === "duplicate").length, color: "var(--text-muted)", icon: <FileSpreadsheet className="w-3.5 h-3.5" /> },
  ], [matches]);

  const totalCount = matches.length;
  const matchedPct = totalCount > 0 ? Math.round((bucketCounts[0].count / totalCount) * 100) : 0;

  const filteredMatches = useMemo(() => {
    let result = matches;
    if (activeBucket !== "all") {
      result = result.filter((m) => m.bucket === activeBucket);
    }
    if (searchQuery) {
      const q = searchQuery.toLowerCase();
      result = result.filter((m) => {
        const prMatch = m.prEntry && (
          m.prEntry.invoiceNumber.toLowerCase().includes(q) ||
          m.prEntry.supplierName.toLowerCase().includes(q) ||
          m.prEntry.supplierGstin.toLowerCase().includes(q)
        );
        const g2bMatch = m.gstr2bEntry && (
          m.gstr2bEntry.invoiceNumber.toLowerCase().includes(q) ||
          m.gstr2bEntry.supplierName.toLowerCase().includes(q) ||
          m.gstr2bEntry.supplierGstin.toLowerCase().includes(q)
        );
        return prMatch || g2bMatch;
      });
    }
    if (showAiOnly) {
      result = result.filter((m) => m.aiSuggestion);
    }
    if (reasonFilter !== "all") {
      result = result.filter((m) => m.mismatchReasons.some((r) => r.toLowerCase().includes(reasonFilter.toLowerCase())));
    }
    return result;
  }, [matches, activeBucket, searchQuery, showAiOnly, reasonFilter]);

  const handleRunRecon = useCallback(async () => {
    setReconRunning(true);
    setLoading(true);
    await new Promise((r) => setTimeout(r, 2000));
    setMatches(generateMockData());
    setLoading(false);
    setReconRunning(false);
  }, []);

  const handleAcceptMatch = useCallback((matchId: string) => {
    setMatches((prev) => prev.map((m) => m.id === matchId ? { ...m, accepted: true } : m));
  }, []);

  const handleRejectMatch = useCallback((matchId: string) => {
    setMatches((prev) => prev.map((m) => m.id === matchId ? { ...m, accepted: false } : m));
  }, []);

  const handleBulkAcceptMatched = useCallback(() => {
    setMatches((prev) => prev.map((m) => m.bucket === "matched" ? { ...m, accepted: true } : m));
  }, []);

  const handleIMSAction = useCallback((matchId: string, action: IMSAction) => {
    setMatches((prev) => prev.map((m) => {
      if (m.id !== matchId || !m.gstr2bEntry) return m;
      return { ...m, gstr2bEntry: { ...m.gstr2bEntry, imsAction: action } };
    }));
  }, []);

  const getEntryDisplay = (entry: PREntry | GSTR2BEntry | null, side: "pr" | "2b") => {
    if (!entry) {
      return (
        <div className="flex-1 p-3 bg-[var(--bg-tertiary)] rounded-lg border border-dashed border-[var(--border-default)] flex items-center justify-center min-h-[120px]">
          <span className="text-xs text-[var(--text-muted)] italic">
            {side === "pr" ? "No Purchase Register entry" : "No GSTR-2B entry"}
          </span>
        </div>
      );
    }
    return (
      <div className="flex-1 p-3 bg-[var(--bg-tertiary)] rounded-lg border border-[var(--border-default)] space-y-1.5">
        <div className="flex items-center justify-between">
          <span className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">
            {side === "pr" ? "Purchase Register" : "GSTR-2B"}
          </span>
          {side === "2b" && "imsAction" in entry && (
            <span className={cn(
              "text-[9px] font-bold uppercase px-1.5 py-0.5 rounded",
              (entry as GSTR2BEntry).imsAction === "accept" && "bg-[var(--success-muted)] text-[var(--success)]",
              (entry as GSTR2BEntry).imsAction === "reject" && "bg-[var(--danger)]/10 text-[var(--danger)]",
              (entry as GSTR2BEntry).imsAction === "pending" && "bg-[var(--warning-muted)] text-[var(--warning)]",
            )}>
              IMS: {(entry as GSTR2BEntry).imsAction}
            </span>
          )}
        </div>
        <div className="space-y-1 text-xs">
          {[
            { label: "Invoice", value: entry.invoiceNumber },
            { label: "Supplier", value: entry.supplierName },
            { label: "GSTIN", value: entry.supplierGstin, mono: true },
            { label: "Date", value: entry.invoiceDate },
            { label: "Taxable", value: formatINR(entry.taxableValue), mono: true },
            { label: "CGST", value: formatINR(entry.cgst), mono: true },
            { label: "SGST", value: formatINR(entry.sgst), mono: true },
            { label: "IGST", value: formatINR(entry.igst), mono: true },
            { label: "Total", value: formatINR(entry.totalValue), mono: true, bold: true },
          ].map((d) => (
            <div key={d.label} className="flex justify-between">
              <span className="text-[var(--text-muted)]">{d.label}</span>
              <span className={cn("text-foreground", d.mono && "font-mono", d.bold && "font-bold")}>{d.value}</span>
            </div>
          ))}
        </div>
      </div>
    );
  };

  const getBucketBadge = (bucket: MatchBucket) => {
    const config: Record<MatchBucket, { label: string; className: string }> = {
      matched: { label: "Matched", className: "bg-[var(--success-muted)] text-[var(--success)]" },
      mismatch: { label: "Mismatch", className: "bg-[var(--warning-muted)] text-[var(--warning)]" },
      partial: { label: "Partial", className: "bg-[var(--accent-muted)] text-[var(--accent)]" },
      missing_2b: { label: "Missing 2B", className: "bg-[var(--danger)]/10 text-[var(--danger)]" },
      missing_pr: { label: "Missing PR", className: "bg-[var(--danger)]/10 text-[var(--danger)]" },
      duplicate: { label: "Duplicate", className: "bg-[var(--bg-tertiary)] text-[var(--text-muted)]" },
    };
    const c = config[bucket];
    return <span className={cn("px-1.5 py-0.5 rounded text-[9px] font-bold uppercase", c.className)}>{c.label}</span>;
  };

  return (
    <div className="space-y-4" data-testid="recon-workspace">
      {/* Header */}
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-3">
          <Link href="/compliance/gst" className="text-foreground-muted hover:text-foreground transition-colors">
            <ArrowLeft className="w-5 h-5" />
          </Link>
          <div>
            <h1 className="text-heading-xl text-foreground">ITC Reconciliation</h1>
            <p className="text-body-sm text-foreground-muted">
              29AABCA1234A1Z5 &middot; April 2026 &middot; Purchase Register vs GSTR-2B
            </p>
          </div>
        </div>
        <div className="flex items-center gap-2">
          <button
            data-testid="export-button"
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-foreground-muted hover:text-foreground hover:bg-[var(--bg-tertiary)] transition-colors"
          >
            <Download className="w-3.5 h-3.5" />
            Export
          </button>
          <button
            data-testid="run-recon-button"
            onClick={handleRunRecon}
            disabled={reconRunning}
            className={cn(
              "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold transition-colors",
              "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
              reconRunning && "opacity-50 cursor-not-allowed",
            )}
          >
            {reconRunning ? <Loader2 className="w-3.5 h-3.5 animate-spin" /> : <RefreshCw className="w-3.5 h-3.5" />}
            {reconRunning ? "Running..." : "Run Reconciliation"}
          </button>
        </div>
      </div>

      {/* Bucket Summary Bar */}
      <div data-testid="bucket-summary" className="bg-app-card border border-app-border rounded-card p-4">
        <div className="flex items-center justify-between mb-3">
          <div className="text-xs font-semibold text-foreground">Match Summary</div>
          <div className="text-xs text-foreground-muted">
            {totalCount} invoices &middot; <span className="text-[var(--success)] font-semibold">{matchedPct}% matched</span>
          </div>
        </div>
        <div className="flex gap-2">
          <button
            data-testid="bucket-all"
            onClick={() => setActiveBucket("all")}
            className={cn(
              "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium transition-colors border",
              activeBucket === "all"
                ? "border-[var(--accent)] bg-[var(--accent-muted)] text-[var(--accent)]"
                : "border-[var(--border-default)] text-foreground-muted hover:bg-[var(--bg-tertiary)]",
            )}
          >
            All ({totalCount})
          </button>
          {bucketCounts.map((b) => (
            <button
              key={b.bucket}
              data-testid={`bucket-${b.bucket}`}
              onClick={() => setActiveBucket(b.bucket)}
              className={cn(
                "flex items-center gap-1.5 px-3 py-2 rounded-lg text-xs font-medium transition-colors border",
                activeBucket === b.bucket
                  ? "border-[var(--accent)] bg-[var(--accent-muted)] text-[var(--accent)]"
                  : "border-[var(--border-default)] text-foreground-muted hover:bg-[var(--bg-tertiary)]",
              )}
            >
              <span style={{ color: b.color }}>{b.icon}</span>
              {b.label} ({b.count})
            </button>
          ))}
        </div>
      </div>

      {/* Filter Bar */}
      <div className="bg-app-card border border-app-border rounded-card px-4 py-3 flex items-center gap-3">
        <div className="relative flex-1 max-w-xs">
          <Search className="absolute left-2.5 top-1/2 -translate-y-1/2 w-3.5 h-3.5 text-[var(--text-muted)]" />
          <input
            data-testid="search-input"
            type="text"
            placeholder="Search invoice, supplier, GSTIN..."
            value={searchQuery}
            onChange={(e) => setSearchQuery(e.target.value)}
            className="w-full pl-8 pr-3 py-1.5 rounded-lg text-xs bg-[var(--bg-tertiary)] border border-[var(--border-default)] text-foreground placeholder:text-[var(--text-muted)] focus:outline-none focus:border-[var(--accent)]"
          />
        </div>
        <select
          data-testid="reason-filter"
          value={reasonFilter}
          onChange={(e) => setReasonFilter(e.target.value)}
          className="px-3 py-1.5 rounded-lg text-xs bg-[var(--bg-tertiary)] border border-[var(--border-default)] text-foreground focus:outline-none focus:border-[var(--accent)]"
        >
          <option value="all">All Reasons</option>
          <option value="amount">Amount Mismatch</option>
          <option value="date">Date Mismatch</option>
          <option value="not found">Missing</option>
          <option value="duplicate">Duplicate</option>
          <option value="partial">Partial</option>
        </select>
        <button
          data-testid="ai-toggle"
          onClick={() => setShowAiOnly(!showAiOnly)}
          className={cn(
            "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors",
            showAiOnly
              ? "border-[var(--accent)] bg-[var(--accent-muted)] text-[var(--accent)]"
              : "border-[var(--border-default)] text-foreground-muted hover:bg-[var(--bg-tertiary)]",
          )}
        >
          <Sparkles className="w-3.5 h-3.5" />
          AI Suggestions
        </button>
        <div className="ml-auto flex gap-2">
          <button
            data-testid="bulk-accept-button"
            onClick={handleBulkAcceptMatched}
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium bg-[var(--success-muted)] text-[var(--success)] hover:opacity-90 transition-colors"
          >
            <Check className="w-3.5 h-3.5" />
            Accept All Matched
          </button>
          <button
            data-testid="send-reminders-button"
            className="flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border border-[var(--border-default)] text-foreground-muted hover:bg-[var(--bg-tertiary)] transition-colors"
          >
            <Mail className="w-3.5 h-3.5" />
            Send Reminders
          </button>
        </div>
      </div>

      {/* Split Pane Table */}
      <div data-testid="recon-table" className="bg-app-card border border-app-border rounded-card overflow-hidden">
        {/* Table Header */}
        <div className="grid grid-cols-[40px_1fr_1fr_120px_100px_80px] items-center bg-[var(--bg-tertiary)] border-b border-[var(--border-default)] px-4 py-2.5">
          <div />
          <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">Purchase Register</div>
          <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)]">GSTR-2B / IMS</div>
          <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)] text-center">Bucket</div>
          <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)] text-center">Confidence</div>
          <div className="text-[10px] font-bold uppercase tracking-wide text-[var(--text-muted)] text-center">Actions</div>
        </div>

        {/* Table Body */}
        {loading ? (
          <div className="flex items-center justify-center py-16">
            <Loader2 className="w-6 h-6 animate-spin text-[var(--accent)]" />
            <span className="ml-2 text-xs text-foreground-muted">Running 5-stage reconciliation pipeline...</span>
          </div>
        ) : filteredMatches.length === 0 ? (
          <div className="flex flex-col items-center justify-center py-16">
            <Search className="w-8 h-8 text-[var(--text-muted)] mb-2" />
            <span className="text-xs text-foreground-muted">No matches found for current filters</span>
          </div>
        ) : (
          <div className="divide-y divide-[var(--border-default)]">
            {filteredMatches.map((match) => {
              const isExpanded = expandedRow === match.id;
              return (
                <div key={match.id} data-testid={`recon-row-${match.id}`}>
                  {/* Collapsed row */}
                  <div
                    className={cn(
                      "grid grid-cols-[40px_1fr_1fr_120px_100px_80px] items-center px-4 py-2.5 cursor-pointer hover:bg-[var(--bg-tertiary)] transition-colors",
                      match.accepted === true && "bg-[var(--success-muted)]/30",
                      match.accepted === false && "bg-[var(--danger)]/5",
                    )}
                    onClick={() => setExpandedRow(isExpanded ? null : match.id)}
                  >
                    <div className="flex items-center gap-1">
                      {isExpanded ? <ChevronDown className="w-3.5 h-3.5 text-foreground-muted" /> : <ChevronRight className="w-3.5 h-3.5 text-foreground-muted" />}
                      {match.aiSuggestion && (
                        <Sparkles className="w-3 h-3 text-[var(--accent)]" data-testid={`ai-icon-${match.id}`} />
                      )}
                    </div>
                    <div className="text-xs">
                      {match.prEntry ? (
                        <div className="flex items-center gap-3">
                          <span className="font-mono font-medium text-foreground">{match.prEntry.invoiceNumber}</span>
                          <span className="text-foreground-muted truncate max-w-[160px]">{match.prEntry.supplierName}</span>
                          <span className="font-mono text-foreground font-semibold">{formatINR(match.prEntry.totalValue)}</span>
                        </div>
                      ) : (
                        <span className="text-[var(--text-muted)] italic">—</span>
                      )}
                    </div>
                    <div className="text-xs">
                      {match.gstr2bEntry ? (
                        <div className="flex items-center gap-3">
                          <span className="font-mono font-medium text-foreground">{match.gstr2bEntry.invoiceNumber}</span>
                          <span className="text-foreground-muted truncate max-w-[160px]">{match.gstr2bEntry.supplierName}</span>
                          <span className="font-mono text-foreground font-semibold">{formatINR(match.gstr2bEntry.totalValue)}</span>
                        </div>
                      ) : (
                        <span className="text-[var(--text-muted)] italic">—</span>
                      )}
                    </div>
                    <div className="text-center">{getBucketBadge(match.bucket)}</div>
                    <div className="text-center">
                      <span className={cn(
                        "text-xs font-mono font-semibold",
                        match.confidence >= 90 ? "text-[var(--success)]" :
                        match.confidence >= 60 ? "text-[var(--warning)]" :
                        "text-[var(--danger)]",
                      )}>
                        {match.confidence > 0 ? `${match.confidence}%` : "—"}
                      </span>
                    </div>
                    <div className="flex justify-center gap-1">
                      {match.accepted === null && (
                        <>
                          <button
                            data-testid={`accept-${match.id}`}
                            onClick={(e) => { e.stopPropagation(); handleAcceptMatch(match.id); }}
                            className="p-1 rounded hover:bg-[var(--success-muted)] text-[var(--success)] transition-colors"
                            title="Accept"
                          >
                            <Check className="w-3.5 h-3.5" />
                          </button>
                          <button
                            data-testid={`reject-${match.id}`}
                            onClick={(e) => { e.stopPropagation(); handleRejectMatch(match.id); }}
                            className="p-1 rounded hover:bg-[var(--danger)]/10 text-[var(--danger)] transition-colors"
                            title="Skip"
                          >
                            <SkipForward className="w-3.5 h-3.5" />
                          </button>
                        </>
                      )}
                      {match.accepted === true && <CheckCircle2 className="w-4 h-4 text-[var(--success)]" />}
                      {match.accepted === false && <XCircle className="w-4 h-4 text-[var(--danger)]" />}
                    </div>
                  </div>

                  {/* Expanded detail */}
                  {isExpanded && (
                    <div data-testid={`detail-${match.id}`} className="px-4 py-4 bg-[var(--bg-secondary)] border-t border-[var(--border-default)]">
                      <div className="flex gap-4">
                        {getEntryDisplay(match.prEntry, "pr")}
                        <div className="flex flex-col items-center justify-center gap-2 px-2">
                          <ArrowUpDown className="w-4 h-4 text-[var(--text-muted)]" />
                          {match.mismatchReasons.length > 0 && (
                            <div className="space-y-1">
                              {match.mismatchReasons.map((reason, idx) => (
                                <div key={idx} className="text-[9px] px-1.5 py-0.5 rounded bg-[var(--warning-muted)] text-[var(--warning)] whitespace-nowrap">
                                  {reason}
                                </div>
                              ))}
                            </div>
                          )}
                        </div>
                        {getEntryDisplay(match.gstr2bEntry, "2b")}
                      </div>
                      {/* IMS Action toolbar */}
                      {match.gstr2bEntry && (
                        <div className="mt-3 flex items-center gap-2 pt-3 border-t border-[var(--border-default)]">
                          <span className="text-[10px] font-semibold text-[var(--text-muted)] uppercase tracking-wide">IMS Action:</span>
                          {(["accept", "reject", "pending"] as IMSAction[]).map((action) => (
                            <button
                              key={action}
                              data-testid={`ims-${action}-${match.id}`}
                              onClick={() => handleIMSAction(match.id, action)}
                              className={cn(
                                "px-2.5 py-1 rounded text-[10px] font-semibold uppercase transition-colors",
                                match.gstr2bEntry?.imsAction === action
                                  ? action === "accept" ? "bg-[var(--success)] text-white"
                                    : action === "reject" ? "bg-[var(--danger)] text-white"
                                    : "bg-[var(--warning)] text-white"
                                  : "bg-[var(--bg-tertiary)] text-[var(--text-muted)] hover:bg-[var(--border-default)]",
                              )}
                            >
                              {action}
                            </button>
                          ))}
                          {match.aiSuggestion && (
                            <div className="ml-auto flex items-center gap-1.5 text-[10px] text-[var(--accent)] font-medium">
                              <Sparkles className="w-3 h-3" />
                              AI suggests: Accept (high-confidence fuzzy match)
                            </div>
                          )}
                        </div>
                      )}
                    </div>
                  )}
                </div>
              );
            })}
          </div>
        )}
      </div>

      {/* Summary Footer */}
      <div className="bg-app-card border border-app-border rounded-card p-4">
        <div className="grid grid-cols-4 gap-4">
          {[
            { label: "Total ITC Available", value: formatINR(matches.filter((m) => m.gstr2bEntry).reduce((s, m) => s + (m.gstr2bEntry?.cgst ?? 0) + (m.gstr2bEntry?.sgst ?? 0) + (m.gstr2bEntry?.igst ?? 0), 0)), color: "text-foreground" },
            { label: "ITC Matched", value: formatINR(matches.filter((m) => m.bucket === "matched" && m.gstr2bEntry).reduce((s, m) => s + (m.gstr2bEntry?.cgst ?? 0) + (m.gstr2bEntry?.sgst ?? 0) + (m.gstr2bEntry?.igst ?? 0), 0)), color: "text-[var(--success)]" },
            { label: "ITC at Risk", value: formatINR(matches.filter((m) => ["mismatch", "missing_2b"].includes(m.bucket)).reduce((s, m) => s + ((m.prEntry?.cgst ?? 0) + (m.prEntry?.sgst ?? 0) + (m.prEntry?.igst ?? 0)), 0)), color: "text-[var(--danger)]" },
            { label: "Pending IMS Actions", value: String(matches.filter((m) => m.gstr2bEntry?.imsAction === "pending").length), color: "text-[var(--warning)]" },
          ].map((stat) => (
            <div key={stat.label} className="text-center">
              <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">{stat.label}</div>
              <div className={cn("text-lg font-bold font-mono", stat.color)}>{stat.value}</div>
            </div>
          ))}
        </div>
      </div>
    </div>
  );
}
