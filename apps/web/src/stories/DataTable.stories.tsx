import type { Meta, StoryObj } from "@storybook/react";
import { DataTable } from "@complai/ui-components";

type InvoiceRow = Record<string, unknown> & {
  invoiceNo: string; vendor: string; date: string; amount: string; status: string;
};

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

const COLUMNS = [
  { key: "invoiceNo", header: "Invoice No", render: (r: InvoiceRow) => <span className="font-mono text-[var(--text-primary)]">{r.invoiceNo}</span> },
  { key: "vendor", header: "Vendor" },
  { key: "date", header: "Date" },
  { key: "amount", header: "Amount", align: "right" as const },
  { key: "status", header: "Status", render: (r: InvoiceRow) => <span className={`font-semibold ${statusColor[r.status] ?? ""}`}>{r.status}</span> },
];

const meta: Meta<typeof DataTable> = {
  title: "Compliance/DataTable",
  component: DataTable,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof DataTable>;

export const Compact: Story = {
  render: () => <DataTable columns={COLUMNS} data={SAMPLE_DATA} density="compact" onRowClick={(r) => alert(String(r.invoiceNo))} />,
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
