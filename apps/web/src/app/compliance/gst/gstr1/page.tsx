"use client";

import { useState, useCallback } from "react";
import {
  ArrowLeft, ArrowRight, Download, CheckCircle2, FileCheck2,
  AlertTriangle, Loader2, FileSpreadsheet,
} from "lucide-react";
import Link from "next/link";
import { cn, formatINR } from "@/lib/utils";
import { FilingConfirmationModal } from "@complai/ui-components";

type WizardStep = "ingest" | "validate" | "review" | "approve" | "file" | "acknowledge";

interface StepDef {
  id: WizardStep;
  label: string;
  number: number;
}

const STEPS: StepDef[] = [
  { id: "ingest", label: "Ingest Data", number: 1 },
  { id: "validate", label: "Validate", number: 2 },
  { id: "review", label: "Review Sections", number: 3 },
  { id: "approve", label: "Approve", number: 4 },
  { id: "file", label: "File Return", number: 5 },
  { id: "acknowledge", label: "Acknowledgement", number: 6 },
];

const SECTIONS = [
  { id: "b2b", label: "B2B Invoices", description: "Registered recipients" },
  { id: "b2cl", label: "B2C Large", description: "Unregistered, > ₹2.5L inter-state" },
  { id: "b2cs", label: "B2C Small", description: "Unregistered, ≤ ₹2.5L" },
  { id: "cdnr", label: "Credit/Debit Notes (Registered)", description: "CDN to registered" },
  { id: "cdnur", label: "Credit/Debit Notes (Unregistered)", description: "CDN to unregistered" },
  { id: "exp", label: "Exports", description: "Export invoices" },
  { id: "nil", label: "NIL Rated", description: "Nil/exempt/non-GST" },
  { id: "hsn", label: "HSN Summary", description: "HSN-wise summary" },
  { id: "docs", label: "Document Issued", description: "Document serial numbers" },
];

interface SectionSummary {
  section: string;
  invoiceCount: number;
  taxableValue: number;
  cgst: number;
  sgst: number;
  igst: number;
  totalTax: number;
}

interface FilingState {
  gstin: string;
  period: string;
  filingId: string | null;
  ingestedCount: number;
  sections: SectionSummary[];
  errorCount: number;
  isApproved: boolean;
  arn: string | null;
  totalTax: number;
}

export default function GSTR1FilingWizard() {
  const [currentStep, setCurrentStep] = useState<WizardStep>("ingest");
  const [loading, setLoading] = useState(false);
  const [showConfirmModal, setShowConfirmModal] = useState(false);
  const [signMethod, setSignMethod] = useState<"dsc" | "evc">("evc");

  const [filing, setFiling] = useState<FilingState>({
    gstin: "29AABCA1234A1Z5",
    period: "042026",
    filingId: null,
    ingestedCount: 0,
    sections: [],
    errorCount: 0,
    isApproved: false,
    arn: null,
    totalTax: 0,
  });

  const stepIndex = STEPS.findIndex((s) => s.id === currentStep);

  const simulateIngest = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1200));

    const mockSections: SectionSummary[] = [
      { section: "b2b", invoiceCount: 50, taxableValue: 4500000, cgst: 202500, sgst: 202500, igst: 405000, totalTax: 810000 },
      { section: "b2cl", invoiceCount: 5, taxableValue: 1750000, cgst: 0, sgst: 0, igst: 315000, totalTax: 315000 },
      { section: "b2cs", invoiceCount: 15, taxableValue: 750000, cgst: 67500, sgst: 67500, igst: 0, totalTax: 135000 },
      { section: "cdnr", invoiceCount: 10, taxableValue: 280000, cgst: 12600, sgst: 12600, igst: 25200, totalTax: 50400 },
      { section: "cdnur", invoiceCount: 3, taxableValue: 60000, cgst: 2700, sgst: 2700, igst: 0, totalTax: 5400 },
      { section: "exp", invoiceCount: 5, taxableValue: 2500000, cgst: 0, sgst: 0, igst: 0, totalTax: 0 },
      { section: "nil", invoiceCount: 5, taxableValue: 200000, cgst: 0, sgst: 0, igst: 0, totalTax: 0 },
      { section: "hsn", invoiceCount: 0, taxableValue: 0, cgst: 0, sgst: 0, igst: 0, totalTax: 0 },
      { section: "docs", invoiceCount: 0, taxableValue: 0, cgst: 0, sgst: 0, igst: 0, totalTax: 0 },
    ];

    const totalTax = mockSections.reduce((sum, s) => sum + s.totalTax, 0);

    setFiling((prev) => ({
      ...prev,
      filingId: crypto.randomUUID(),
      ingestedCount: 93,
      sections: mockSections,
      totalTax,
    }));
    setLoading(false);
    setCurrentStep("validate");
  }, []);

  const simulateValidate = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 800));
    setFiling((prev) => ({ ...prev, errorCount: 0 }));
    setLoading(false);
    setCurrentStep("review");
  }, []);

  const handleApprove = useCallback(async () => {
    setLoading(true);
    await new Promise((r) => setTimeout(r, 600));
    setFiling((prev) => ({ ...prev, isApproved: true }));
    setLoading(false);
    setCurrentStep("file");
  }, []);

  const handleFile = useCallback(async () => {
    setShowConfirmModal(false);
    setLoading(true);
    await new Promise((r) => setTimeout(r, 1500));
    setFiling((prev) => ({ ...prev, arn: `AA290420260000${Math.floor(Math.random() * 9000 + 1000)}` }));
    setLoading(false);
    setCurrentStep("acknowledge");
  }, []);

  return (
    <div className="space-y-6">
      <div className="flex items-center gap-3">
        <Link
          href="/compliance/gst"
          className="text-foreground-muted hover:text-foreground transition-colors"
        >
          <ArrowLeft className="w-5 h-5" />
        </Link>
        <div>
          <h1 className="text-heading-xl text-foreground">GSTR-1 Filing</h1>
          <p className="text-body-sm text-foreground-muted">
            {filing.gstin} &middot; Period: Apr 2026
          </p>
        </div>
      </div>

      {/* Step Indicator */}
      <div className="bg-app-card border border-app-border rounded-card p-4">
        <div className="flex items-center gap-1">
          {STEPS.map((step, i) => {
            const isActive = step.id === currentStep;
            const isCompleted = i < stepIndex;
            return (
              <div key={step.id} className="flex items-center flex-1">
                <div className="flex items-center gap-2 flex-1">
                  <div
                    className={cn(
                      "w-7 h-7 rounded-full flex items-center justify-center text-xs font-bold flex-shrink-0",
                      isCompleted
                        ? "bg-[var(--success)] text-white"
                        : isActive
                          ? "bg-[var(--accent)] text-[var(--accent-text)]"
                          : "bg-[var(--bg-tertiary)] text-[var(--text-muted)]",
                    )}
                  >
                    {isCompleted ? <CheckCircle2 className="w-4 h-4" /> : step.number}
                  </div>
                  <span
                    className={cn(
                      "text-xs font-medium whitespace-nowrap",
                      isActive ? "text-foreground" : "text-foreground-muted",
                    )}
                  >
                    {step.label}
                  </span>
                </div>
                {i < STEPS.length - 1 && (
                  <div
                    className={cn(
                      "h-[2px] flex-1 mx-2 rounded",
                      i < stepIndex ? "bg-[var(--success)]" : "bg-[var(--border-default)]",
                    )}
                  />
                )}
              </div>
            );
          })}
        </div>
      </div>

      {/* Step Content */}
      <div className="bg-app-card border border-app-border rounded-card">
        {currentStep === "ingest" && (
          <div className="p-6 space-y-4">
            <div>
              <h2 className="text-heading-lg text-foreground">Step 1: Ingest Sales Data</h2>
              <p className="text-body-sm text-foreground-muted mt-1">
                Pull AR invoices from Aura and categorise into GSTR-1 sections.
              </p>
            </div>
            <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
              <div className="flex gap-3 text-xs">
                <span className="text-foreground-muted min-w-[100px]">GSTIN:</span>
                <span className="text-foreground font-medium font-mono">{filing.gstin}</span>
              </div>
              <div className="flex gap-3 text-xs">
                <span className="text-foreground-muted min-w-[100px]">Period:</span>
                <span className="text-foreground font-medium">April 2026</span>
              </div>
              <div className="flex gap-3 text-xs">
                <span className="text-foreground-muted min-w-[100px]">Source:</span>
                <span className="text-foreground font-medium">Aura (Order-to-Cash)</span>
              </div>
            </div>
            <button
              data-testid="ingest-button"
              onClick={simulateIngest}
              disabled={loading}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
                "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
                loading && "opacity-50 cursor-not-allowed",
              )}
            >
              {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <Download className="w-4 h-4" />}
              {loading ? "Ingesting..." : "Ingest from Aura"}
            </button>
          </div>
        )}

        {currentStep === "validate" && (
          <div className="p-6 space-y-4">
            <div>
              <h2 className="text-heading-lg text-foreground">Step 2: Validate Invoices</h2>
              <p className="text-body-sm text-foreground-muted mt-1">
                {filing.ingestedCount} invoices ingested. Run validation checks.
              </p>
            </div>
            <div className="flex items-center gap-3 p-3 bg-[var(--success-muted)] rounded-lg">
              <CheckCircle2 className="w-5 h-5 text-[var(--success)]" />
              <span className="text-xs text-[var(--success)] font-medium">
                {filing.ingestedCount} invoices categorised into {filing.sections.filter((s) => s.invoiceCount > 0).length} sections
              </span>
            </div>
            <button
              data-testid="validate-button"
              onClick={simulateValidate}
              disabled={loading}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
                "bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]",
                loading && "opacity-50 cursor-not-allowed",
              )}
            >
              {loading ? <Loader2 className="w-4 h-4 animate-spin" /> : <FileCheck2 className="w-4 h-4" />}
              {loading ? "Validating..." : "Run Validation"}
            </button>
          </div>
        )}

        {currentStep === "review" && (
          <div className="p-6 space-y-4">
            <div className="flex items-center justify-between">
              <div>
                <h2 className="text-heading-lg text-foreground">Step 3: Review Sections</h2>
                <p className="text-body-sm text-foreground-muted mt-1">
                  Review categorised data across all 9 GSTR-1 sections.
                </p>
              </div>
              {filing.errorCount === 0 && (
                <div className="flex items-center gap-1.5 text-xs text-[var(--success)] font-medium">
                  <CheckCircle2 className="w-4 h-4" />
                  All validations passed
                </div>
              )}
            </div>
            <div data-testid="section-table" className="border border-[var(--border-default)] rounded-lg overflow-hidden">
              <table className="w-full text-xs">
                <thead>
                  <tr className="bg-[var(--bg-tertiary)] border-b border-[var(--border-default)]">
                    <th className="text-left px-4 py-2.5 text-[var(--text-muted)] font-semibold">Section</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">Invoices</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">Taxable Value</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">CGST</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">SGST</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">IGST</th>
                    <th className="text-right px-4 py-2.5 text-[var(--text-muted)] font-semibold">Total Tax</th>
                  </tr>
                </thead>
                <tbody>
                  {SECTIONS.map((sec) => {
                    const data = filing.sections.find((s) => s.section === sec.id);
                    if (!data || data.invoiceCount === 0) return null;
                    return (
                      <tr key={sec.id} className="border-b border-[var(--border-default)] last:border-0 hover:bg-[var(--bg-tertiary)] transition-colors">
                        <td className="px-4 py-2.5">
                          <div className="font-medium text-foreground">{sec.label}</div>
                          <div className="text-[10px] text-foreground-muted">{sec.description}</div>
                        </td>
                        <td className="px-4 py-2.5 text-right font-mono text-foreground">{data.invoiceCount}</td>
                        <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(data.taxableValue)}</td>
                        <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(data.cgst)}</td>
                        <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(data.sgst)}</td>
                        <td className="px-4 py-2.5 text-right font-mono text-foreground">{formatINR(data.igst)}</td>
                        <td className="px-4 py-2.5 text-right font-mono font-semibold text-foreground">{formatINR(data.totalTax)}</td>
                      </tr>
                    );
                  })}
                </tbody>
                <tfoot>
                  <tr className="bg-[var(--bg-tertiary)] border-t-2 border-[var(--accent)]">
                    <td className="px-4 py-2.5 font-bold text-foreground">Total</td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {filing.sections.reduce((sum, s) => sum + s.invoiceCount, 0)}
                    </td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {formatINR(filing.sections.reduce((sum, s) => sum + s.taxableValue, 0))}
                    </td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {formatINR(filing.sections.reduce((sum, s) => sum + s.cgst, 0))}
                    </td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {formatINR(filing.sections.reduce((sum, s) => sum + s.sgst, 0))}
                    </td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {formatINR(filing.sections.reduce((sum, s) => sum + s.igst, 0))}
                    </td>
                    <td className="px-4 py-2.5 text-right font-mono font-bold text-foreground">
                      {formatINR(filing.totalTax)}
                    </td>
                  </tr>
                </tfoot>
              </table>
            </div>
            <div className="flex justify-end">
              <button
                data-testid="proceed-approve-button"
                onClick={() => setCurrentStep("approve")}
                className="flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]"
              >
                Proceed to Approval
                <ArrowRight className="w-4 h-4" />
              </button>
            </div>
          </div>
        )}

        {currentStep === "approve" && (
          <div className="p-6 space-y-4">
            <div>
              <h2 className="text-heading-lg text-foreground">Step 4: Maker-Checker Approval</h2>
              <p className="text-body-sm text-foreground-muted mt-1">
                A second user (checker) must approve before filing.
              </p>
            </div>
            <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl">
              <div className="px-5 py-3 border-b border-[var(--border-default)]">
                <div className="text-[10px] font-semibold uppercase tracking-wide text-[var(--warning)]">
                  Approval Required
                </div>
              </div>
              <div className="px-5 py-4 space-y-3">
                <div className="text-sm font-bold text-foreground">GSTR-1 Filing — Apr 2026</div>
                <div className="text-xs text-foreground-muted">Submitted by: Jayesh H &middot; {new Date().toLocaleDateString("en-IN")}</div>
                <div className="space-y-1.5">
                  {[
                    { label: "GSTIN", value: filing.gstin },
                    { label: "Period", value: "April 2026" },
                    { label: "Invoices", value: `${filing.ingestedCount}` },
                    { label: "Total Tax", value: formatINR(filing.totalTax) },
                  ].map((d) => (
                    <div key={d.label} className="flex gap-2 text-xs">
                      <span className="text-foreground-muted min-w-[80px]">{d.label}:</span>
                      <span className="text-foreground font-medium">{d.value}</span>
                    </div>
                  ))}
                </div>
                {filing.totalTax > 1000000 && (
                  <div className="flex items-center gap-1.5 text-[11px] text-[var(--warning)]">
                    <AlertTriangle className="w-3.5 h-3.5" />
                    Tax liability exceeds ₹10,00,000 — requires senior approval
                  </div>
                )}
              </div>
              <div className="px-5 py-3 border-t border-[var(--border-default)] flex justify-end">
                <button
                  data-testid="approve-button"
                  onClick={handleApprove}
                  disabled={loading}
                  className="px-4 py-1.5 text-[11px] font-bold rounded-lg bg-gradient-to-br from-[var(--accent)] to-[var(--accent-hover)] text-[var(--accent-text)] shadow-[var(--shadow-accent)]"
                >
                  {loading ? "Approving..." : "Approve & Continue"}
                </button>
              </div>
            </div>
          </div>
        )}

        {currentStep === "file" && (
          <div className="p-6 space-y-4">
            <div>
              <h2 className="text-heading-lg text-foreground">Step 5: File Return</h2>
              <p className="text-body-sm text-foreground-muted mt-1">
                Save, submit, and file GSTR-1 on the GST portal.
              </p>
            </div>
            <div className="flex items-center gap-3 p-3 bg-[var(--success-muted)] rounded-lg">
              <CheckCircle2 className="w-5 h-5 text-[var(--success)]" />
              <span className="text-xs text-[var(--success)] font-medium">
                Filing approved. Ready to submit to GST portal.
              </span>
            </div>
            <div className="bg-[var(--bg-tertiary)] rounded-lg p-4 space-y-2">
              <div className="flex gap-3 text-xs">
                <span className="text-foreground-muted min-w-[100px]">Total Tax:</span>
                <span className="text-foreground font-bold font-mono">{formatINR(filing.totalTax)}</span>
              </div>
              <div className="flex gap-3 text-xs">
                <span className="text-foreground-muted min-w-[100px]">Invoices:</span>
                <span className="text-foreground font-medium">{filing.ingestedCount}</span>
              </div>
            </div>
            <button
              data-testid="file-button"
              onClick={() => setShowConfirmModal(true)}
              className={cn(
                "flex items-center gap-2 px-4 py-2 rounded-[10px] text-xs font-semibold transition-colors",
                "bg-[var(--danger)] text-white hover:opacity-90",
              )}
            >
              <FileSpreadsheet className="w-4 h-4" />
              File GSTR-1
            </button>
          </div>
        )}

        {currentStep === "acknowledge" && (
          <div data-testid="acknowledge-step" className="p-6 space-y-4 text-center">
            <div className="flex justify-center">
              <div className="w-16 h-16 rounded-full bg-[var(--success-muted)] flex items-center justify-center">
                <CheckCircle2 className="w-8 h-8 text-[var(--success)]" />
              </div>
            </div>
            <div>
              <h2 className="text-heading-lg text-foreground">GSTR-1 Filed Successfully</h2>
              <p className="text-body-sm text-foreground-muted mt-1">
                Your return has been filed on the GST portal.
              </p>
            </div>
            <div data-testid="filing-receipt" className="bg-[var(--bg-tertiary)] rounded-lg p-4 inline-block mx-auto space-y-2 text-left">
              {[
                { label: "ARN", value: filing.arn ?? "" },
                { label: "GSTIN", value: filing.gstin },
                { label: "Period", value: "April 2026" },
                { label: "Filed At", value: new Date().toLocaleString("en-IN") },
                { label: "Total Tax", value: formatINR(filing.totalTax) },
              ].map((d) => (
                <div key={d.label} className="flex gap-3 text-xs">
                  <span className="text-foreground-muted min-w-[80px]">{d.label}:</span>
                  <span className="text-foreground font-bold font-mono">{d.value}</span>
                </div>
              ))}
            </div>
            <div className="pt-2">
              <Link
                href="/compliance/gst"
                className="text-xs text-[var(--accent)] font-medium hover:underline"
              >
                Back to GST Returns
              </Link>
            </div>
          </div>
        )}
      </div>

      <FilingConfirmationModal
        open={showConfirmModal}
        onClose={() => setShowConfirmModal(false)}
        onConfirm={handleFile}
        title="File GSTR-1 — April 2026"
        details={[
          { label: "GSTIN", value: filing.gstin },
          { label: "Period", value: "April 2026" },
          { label: "Invoices", value: `${filing.ingestedCount}` },
          { label: "Total Tax", value: formatINR(filing.totalTax) },
        ]}
        warningText="This action is irreversible. The return will be filed on the GST portal."
        confirmWord="FILE"
        signMethod={signMethod}
        onSignMethodChange={setSignMethod}
      />
    </div>
  );
}
