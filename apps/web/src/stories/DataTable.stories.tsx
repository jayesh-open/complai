import type { Meta, StoryObj } from "@storybook/react";

interface Column<T> {
  key: string;
  header: string;
  render?: (row: T) => React.ReactNode;
  className?: string;
  align?: "left" | "right" | "center";
}

function DataTable<T extends Record<string, unknown>>({
  columns, data, density = "compact", onRowClick, emptyMessage = "No data found",
}: {
  columns: Column<T>[]; data: T[]; density?: "compact" | "comfortable" | "spacious";
  onRowClick?: (row: T) => void; emptyMessage?: string;
}) {
  const rowHeight = { compact: 40, comfortable: 52, spacious: 64 }[density];
  return (
    <div className="bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden">
      <table className="w-full" data-density={density}>
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            {columns.map((col) => (
              <th key={col.key} className={`px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] ${col.align === "right" ? "text-right" : ""}`}>
                {col.header}
              </th>
            ))}
          </tr>
        </thead>
        <tbody>
          {data.length === 0 ? (
            <tr><td colSpan={columns.length} className="text-center py-12 text-[var(--text-muted)] text-sm">{emptyMessage}</td></tr>
          ) : (
            data.map((row, i) => (
              <tr key={i} onClick={() => onRowClick?.(row)}
                className={`border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors ${onRowClick ? "cursor-pointer" : ""}`}
                style={{ height: `${rowHeight}px` }}>
                {columns.map((col) => (
                  <td key={col.key} className={`px-[18px] text-xs text-[var(--text-secondary)] ${col.align === "right" ? "text-right tabular-nums font-semibold" : ""}`}>
                    {col.render ? col.render(row) : String(row[col.key] ?? "")}
                  </td>
                ))}
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}

type InvoiceRow = Record<string, unknown> & { invoiceNo: string; vendor: string; date: string; amount: string; status: string };

const SAMPLE_DATA: InvoiceRow[] = [
  { invoiceNo: "INV-2026-0042", vendor: "Tata Steel Ltd", date: "15/04/2026", amount: "₹4,52,800", status: "Approved" },
  { invoiceNo: "INV-2026-0041", vendor: "Infosys BPO", date: "14/04/2026", amount: "₹1,28,000", status: "Pending" },
  { invoiceNo: "INV-2026-0040", vendor: "Reliance Jio", date: "12/04/2026", amount: "₹86,400", status: "Rejected" },
  { invoiceNo: "INV-2026-0039", vendor: "Wipro Technologies", date: "11/04/2026", amount: "₹2,15,600", status: "Approved" },
  { invoiceNo: "INV-2026-0038", vendor: "Mahindra Logistics", date: "10/04/2026", amount: "₹3,72,000", status: "Pending" },
];

const statusColor: Record<string, string> = {
  Approved: "text-[var(--success)]",
  Pending: "text-[var(--warning)]",
  Rejected: "text-[var(--danger)]",
};

const COLUMNS: Column<InvoiceRow>[] = [
  { key: "invoiceNo", header: "Invoice No", render: (r) => <span className="font-mono text-[var(--text-primary)]">{r.invoiceNo}</span> },
  { key: "vendor", header: "Vendor" },
  { key: "date", header: "Date" },
  { key: "amount", header: "Amount", align: "right" },
  { key: "status", header: "Status", render: (r) => <span className={`font-semibold ${statusColor[r.status] ?? ""}`}>{r.status}</span> },
];

const meta: Meta<typeof DataTable> = {
  title: "Compliance/DataTable",
  component: DataTable as Meta<typeof DataTable>["component"],
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DataTable>;

export const Compact: Story = {
  render: () => <DataTable columns={COLUMNS} data={SAMPLE_DATA} density="compact" onRowClick={(r) => alert(r.invoiceNo)} />,
};

export const Comfortable: Story = {
  render: () => <DataTable columns={COLUMNS} data={SAMPLE_DATA} density="comfortable" />,
};

export const Spacious: Story = {
  render: () => <DataTable columns={COLUMNS} data={SAMPLE_DATA} density="spacious" />,
};

export const Empty: Story = {
  render: () => <DataTable columns={COLUMNS} data={[]} emptyMessage="No invoices match your filter criteria" />,
};
