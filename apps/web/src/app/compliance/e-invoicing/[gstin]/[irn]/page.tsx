"use client";

import { useState, useMemo } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Download, XCircle } from "lucide-react";
import { cn } from "@/lib/utils";
import { formatINR, AuditTrailTimeline, GovStatusPill } from "@complai/ui-components";
import { FilingConfirmationModal } from "@complai/ui-components";
import { generateMockEInvoices, generateAuditTrail } from "../../mock-data";
import { IRNStatusPill } from "../../components/IRNStatusPill";
import { QRCodeDisplay } from "../../components/QRCodeDisplay";
import { CancellationCountdown } from "../../components/CancellationCountdown";
import { SignedJsonViewer } from "../../components/SignedJsonViewer";
import { EInvoiceDetailItems } from "./components/DetailItems";

const ALL_RECORDS = generateMockEInvoices();

export default function EInvoiceDetailPage() {
  const params = useParams<{ gstin: string; irn: string }>();
  const record = useMemo(
    () => ALL_RECORDS.find((r) => r.gstin === params.gstin && r.irn === params.irn),
    [params.gstin, params.irn]
  );

  const [cancelOpen, setCancelOpen] = useState(false);

  if (!record) {
    return (
      <div className="space-y-4">
        <Link
          href={`/compliance/e-invoicing/${params.gstin}`}
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)]"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back
        </Link>
        <div className="text-center py-20 text-[var(--text-muted)]">
          E-Invoice not found
        </div>
      </div>
    );
  }

  const auditEntries = generateAuditTrail(record);
  const canCancel =
    record.status === "GENERATED" &&
    Date.now() - new Date(record.generatedAt).getTime() < 24 * 3600 * 1000;

  return (
    <div className="space-y-6">
      <div>
        <Link
          href={`/compliance/e-invoicing/${params.gstin}`}
          className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--accent)] transition-colors mb-3"
        >
          <ArrowLeft className="w-3.5 h-3.5" />
          Back to {params.gstin}
        </Link>
        <div className="flex items-center justify-between">
          <div className="flex items-center gap-3">
            <h1 className="text-heading-lg text-[var(--text-primary)]">
              {record.invoiceNo}
            </h1>
            <IRNStatusPill status={record.status} />
            <GovStatusPill
              system="IRP"
              status={record.status === "GENERATED" ? "success" : "danger"}
              label={record.status === "GENERATED" ? "Active" : "Cancelled"}
            />
          </div>
          <div className="flex items-center gap-3">
            <button
              className={cn(
                "flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-medium",
                "border border-[var(--border-default)]",
                "text-[var(--text-secondary)] hover:bg-[var(--bg-tertiary)]",
                "transition-colors"
              )}
            >
              <Download className="w-3.5 h-3.5" />
              Download QR
            </button>
            {canCancel && (
              <button
                data-testid="cancel-irn-button"
                onClick={() => setCancelOpen(true)}
                className={cn(
                  "flex items-center gap-2 px-3 py-2 rounded-lg text-xs font-semibold",
                  "bg-[var(--danger)] text-white hover:opacity-90 transition-opacity"
                )}
              >
                <XCircle className="w-3.5 h-3.5" />
                Cancel IRN
              </button>
            )}
          </div>
        </div>
      </div>

      <div className="grid grid-cols-3 gap-6">
        <div className="col-span-2 space-y-6">
          <DetailCard title="IRN Details">
            <DetailRow label="IRN" mono>
              {record.irn}
            </DetailRow>
            <DetailRow label="Ack Number" mono>
              {record.ackNumber}
            </DetailRow>
            <DetailRow label="Ack Date">{record.ackDate}</DetailRow>
            <DetailRow label="Invoice No." mono>
              {record.invoiceNo}
            </DetailRow>
            <DetailRow label="Invoice Date">{record.invoiceDate}</DetailRow>
            <DetailRow label="Supplier GSTIN" mono>
              {record.gstin}
            </DetailRow>
            <DetailRow label="Buyer">
              {record.buyerName}
              <span className="text-[var(--text-muted)] font-mono text-[10px] ml-2">
                {record.buyerGstin}
              </span>
            </DetailRow>
          </DetailCard>

          <DetailCard title="Value Summary">
            <div className="grid grid-cols-5 gap-4">
              <ValueBox label="Taxable" value={record.taxableValue} />
              <ValueBox label="IGST" value={record.igstAmount} />
              <ValueBox label="CGST" value={record.cgstAmount} />
              <ValueBox label="SGST" value={record.sgstAmount} />
              <ValueBox label="Total" value={record.totalValue} bold />
            </div>
          </DetailCard>

          <EInvoiceDetailItems items={record.items} />

          <SignedJsonViewer json={record.signedInvoice} />
        </div>

        <div className="space-y-6">
          <QRCodeDisplay
            value={record.qrCodeData}
            size={200}
            label={`ACK: ${record.ackNumber}`}
          />

          {record.status === "GENERATED" && (
            <CancellationCountdown generatedAt={record.generatedAt} />
          )}

          <DetailCard title="Audit Trail">
            <AuditTrailTimeline entries={auditEntries} />
          </DetailCard>
        </div>
      </div>

      <FilingConfirmationModal
        open={cancelOpen}
        onClose={() => setCancelOpen(false)}
        onConfirm={() => setCancelOpen(false)}
        title="Cancel IRN"
        details={[
          { label: "Invoice", value: record.invoiceNo },
          { label: "IRN", value: `${record.irn.slice(0, 24)}...` },
          { label: "Value", value: formatINR(record.totalValue) },
        ]}
        warningText="This will permanently cancel the IRN. A new IRN must be generated for this invoice."
        confirmWord="CANCEL"
      />
    </div>
  );
}

function DetailCard({
  title,
  children,
}: {
  title: string;
  children: React.ReactNode;
}) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div className="px-5 py-3 border-b border-[var(--border-default)]">
        <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide">
          {title}
        </h3>
      </div>
      <div className="px-5 py-4">{children}</div>
    </div>
  );
}

function DetailRow({
  label,
  children,
  mono,
}: {
  label: string;
  children: React.ReactNode;
  mono?: boolean;
}) {
  return (
    <div className="flex gap-4 py-1.5 text-xs">
      <span className="text-[var(--text-muted)] min-w-[120px] shrink-0">
        {label}
      </span>
      <span
        className={cn(
          "text-[var(--text-primary)] break-all",
          mono && "font-mono text-[11px]"
        )}
      >
        {children}
      </span>
    </div>
  );
}

function ValueBox({
  label,
  value,
  bold,
}: {
  label: string;
  value: number;
  bold?: boolean;
}) {
  return (
    <div className="text-center">
      <div className="text-[10px] text-[var(--text-muted)] uppercase tracking-wide mb-1">
        {label}
      </div>
      <div
        className={cn(
          "text-xs tabular-nums font-mono",
          bold
            ? "font-bold text-[var(--text-primary)]"
            : "text-[var(--text-secondary)]"
        )}
      >
        {formatINR(value)}
      </div>
    </div>
  );
}
