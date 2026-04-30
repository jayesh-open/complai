"use client";

import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, Download } from "lucide-react";
import { cn, formatINR } from "@complai/ui-components";
import { generateChallanDetail } from "../mock-data";
import { ChallanStatusPill } from "../../components/ChallanStatusPill";
import { ChallanAllocationCard } from "../../components/ChallanAllocationCard";

export default function ChallanDetailPage() {
  const params = useParams<{ challanId: string }>();
  const challanId = params.challanId;

  const data = generateChallanDetail(challanId);

  if (!data) {
    return (
      <div className="space-y-4">
        <Link href="/compliance/tds/challans" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]">
          <ArrowLeft className="w-3.5 h-3.5" /> Back to Challans
        </Link>
        <div className="text-center py-12 text-[var(--text-muted)]">Challan not found</div>
      </div>
    );
  }

  return (
    <div className="space-y-6">
      <Link href="/compliance/tds/challans" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)]">
        <ArrowLeft className="w-3.5 h-3.5" /> Back to Challans
      </Link>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6">
        <div className="flex items-start justify-between">
          <div>
            <h1 className="text-heading-lg text-[var(--text-primary)]">Challan {data.challanSerial}</h1>
            <div className="flex items-center gap-3 mt-1">
              <span className="font-mono text-xs text-[var(--text-muted)]">BSR: {data.bsrCode}</span>
              <ChallanStatusPill status={data.status} />
              <span className="text-xs text-[var(--text-muted)]">{data.quarter}</span>
            </div>
          </div>
          <button className={cn("flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-semibold", "bg-[var(--accent)] text-[var(--accent-text)] hover:bg-[var(--accent-hover)] transition-colors")}>
            <Download className="w-3.5 h-3.5" /> Download
          </button>
        </div>

        <div className="grid grid-cols-4 gap-4 mt-6 pt-4 border-t border-[var(--border-default)]">
          <MiniStat label="Deposit Date" value={data.depositDate} />
          <MiniStat label="Bank" value={data.bankName} />
          <MiniStat label="Branch" value={data.branchName} />
          <MiniStat label="TAN" value={data.tan} />
        </div>
        <div className="grid grid-cols-3 gap-4 mt-4">
          <MiniStat label="Total Amount" value={formatINR(data.amount)} />
          <MiniStat label="Allocated" value={formatINR(data.allocatedAmount)} accent />
          <MiniStat label="Unallocated" value={formatINR(data.unallocatedAmount)} warn={data.unallocatedAmount > 0} />
        </div>
      </div>

      <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl p-6">
        <ChallanAllocationCard allocations={data.allocations} totalAmount={data.amount} allocatedAmount={data.allocatedAmount} />
      </div>
    </div>
  );
}

function MiniStat({ label, value, accent, warn }: { label: string; value: string; accent?: boolean; warn?: boolean }) {
  return (
    <div>
      <div className="text-[10px] text-[var(--text-muted)] uppercase font-medium">{label}</div>
      <div className={cn("text-sm font-semibold mt-0.5 tabular-nums", accent && "text-[var(--accent)]", warn && "text-[var(--warning)]", !accent && !warn && "text-[var(--text-primary)]")}>{value}</div>
    </div>
  );
}
