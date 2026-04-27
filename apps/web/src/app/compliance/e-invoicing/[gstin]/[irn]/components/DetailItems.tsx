"use client";

import { formatINR } from "@complai/ui-components";
import type { EInvoiceItem } from "../../../types";

interface EInvoiceDetailItemsProps {
  items: EInvoiceItem[];
}

export function EInvoiceDetailItems({ items }: EInvoiceDetailItemsProps) {
  return (
    <div className="bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-xl overflow-hidden">
      <div className="px-5 py-3 border-b border-[var(--border-default)]">
        <h3 className="text-xs font-semibold text-[var(--text-primary)] uppercase tracking-wide">
          Line Items
        </h3>
      </div>
      <div className="overflow-x-auto">
        <table className="w-full">
          <thead>
            <tr className="border-b border-[var(--border-default)]">
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                #
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                Description
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
                HSN
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                Qty
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                Rate
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                GST %
              </th>
              <th className="px-5 py-2 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-right">
                Amount
              </th>
            </tr>
          </thead>
          <tbody>
            {items.map((item) => (
              <tr
                key={item.slNo}
                className="border-b border-[var(--border-default)] last:border-b-0"
              >
                <td className="px-5 py-2.5 text-xs text-[var(--text-muted)]">
                  {item.slNo}
                </td>
                <td className="px-5 py-2.5 text-xs text-[var(--text-primary)]">
                  {item.description}
                </td>
                <td className="px-5 py-2.5 text-xs font-mono text-[var(--text-secondary)]">
                  {item.hsnCode}
                </td>
                <td className="px-5 py-2.5 text-xs text-[var(--text-secondary)] text-right tabular-nums">
                  {item.quantity} {item.unit}
                </td>
                <td className="px-5 py-2.5 text-xs font-mono text-[var(--text-secondary)] text-right tabular-nums">
                  {formatINR(item.unitPrice)}
                </td>
                <td className="px-5 py-2.5 text-xs text-[var(--text-secondary)] text-right tabular-nums">
                  {item.gstRate}%
                </td>
                <td className="px-5 py-2.5 text-xs font-mono font-semibold text-[var(--text-primary)] text-right tabular-nums">
                  {formatINR(item.totalAmount)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}
