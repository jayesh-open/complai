import type { EInvoiceRecord, SourceInvoice, AuditEntry } from "./types";

const SUPPLIERS = [
  { gstin: "29AABCA1234A1Z5", name: "Acme Manufacturing Ltd" },
  { gstin: "27AABCB5678B1Z3", name: "Bharat Traders Pvt Ltd" },
  { gstin: "33AABCC9012C1Z1", name: "Chennai Exports Ltd" },
];

const BUYERS = [
  { gstin: "07AABCD3456D1Z9", name: "Delhi Distributors Ltd" },
  { gstin: "24AABCE7890E1Z7", name: "Gujarat Retail Corp" },
  { gstin: "36AABCF1234F1Z5", name: "Telangana Tech Solutions" },
  { gstin: "09AABCG5678G1Z3", name: "UP Logistics Pvt Ltd" },
  { gstin: "19AABCH9012H1Z1", name: "Kolkata Trading House" },
];

const HSN_CODES = ["84713010", "85044090", "94036090", "73269099", "39269099"];
const DESCRIPTIONS = [
  "Laptop Computer (15.6 inch)",
  "UPS Power Supply Unit",
  "Ergonomic Office Chair",
  "Steel Mounting Bracket",
  "Plastic Storage Container",
];

function generateIRN(index: number): string {
  const hash = Array.from({ length: 64 }, (_, i) =>
    "0123456789abcdef"[(index * 7 + i * 13) % 16]
  ).join("");
  return hash;
}

function generateItems(count: number): SourceInvoice["items"] {
  return Array.from({ length: count }, (_, i) => {
    const unitPrice = [45000, 8500, 15200, 3400, 1200][i % 5];
    const qty = [2, 5, 10, 25, 50][i % 5];
    return {
      slNo: i + 1,
      description: DESCRIPTIONS[i % 5],
      hsnCode: HSN_CODES[i % 5],
      quantity: qty,
      unit: "NOS",
      unitPrice,
      totalAmount: unitPrice * qty,
      gstRate: 18,
    };
  });
}

export function generateMockEInvoices(): EInvoiceRecord[] {
  const records: EInvoiceRecord[] = [];
  const now = new Date("2026-04-27T10:00:00+05:30");

  for (let i = 0; i < 30; i++) {
    const supplier = SUPPLIERS[i % SUPPLIERS.length];
    const buyer = BUYERS[i % BUYERS.length];
    const genDate = new Date(now.getTime() - i * 3600 * 1000 * (2 + i % 5));
    const items = generateItems(1 + (i % 3));
    const taxableValue = items.reduce((s, it) => s + it.totalAmount, 0);
    const isInterstate = supplier.gstin.slice(0, 2) !== buyer.gstin.slice(0, 2);
    const igst = isInterstate ? Math.round(taxableValue * 0.18) : 0;
    const cgst = isInterstate ? 0 : Math.round(taxableValue * 0.09);
    const sgst = isInterstate ? 0 : Math.round(taxableValue * 0.09);
    const isCancelled = i === 3 || i === 12;

    records.push({
      id: `einv-${String(i + 1).padStart(4, "0")}`,
      invoiceNo: `INV-2026-${String(1000 + i).padStart(5, "0")}`,
      invoiceDate: formatDateDDMMYYYY(new Date(genDate.getTime() - 86400000)),
      gstin: supplier.gstin,
      buyerGstin: buyer.gstin,
      buyerName: buyer.name,
      irn: generateIRN(i),
      ackNumber: `1${String(i + 1).padStart(11, "0")}`,
      ackDate: formatDateDDMMYYYY(genDate),
      status: isCancelled ? "CANCELLED" : "GENERATED",
      totalValue: taxableValue + igst + cgst + sgst,
      taxableValue,
      igstAmount: igst,
      cgstAmount: cgst,
      sgstAmount: sgst,
      generatedAt: genDate.toISOString(),
      cancelledAt: isCancelled
        ? new Date(genDate.getTime() + 3600000).toISOString()
        : undefined,
      cancelReason: isCancelled ? "Duplicate entry" : undefined,
      signedInvoice: JSON.stringify(
        {
          Version: "1.1",
          TranDtls: { TaxSch: "GST", SupTyp: "B2B", RegRev: "N" },
          DocDtls: {
            Typ: "INV",
            No: `INV-2026-${String(1000 + i).padStart(5, "0")}`,
            Dt: formatDateDDMMYYYY(new Date(genDate.getTime() - 86400000)),
          },
          SellerDtls: { Gstin: supplier.gstin, LglNm: supplier.name },
          BuyerDtls: { Gstin: buyer.gstin, LglNm: buyer.name },
          ValDtls: {
            AssVal: taxableValue,
            IgstVal: igst,
            CgstVal: cgst,
            SgstVal: sgst,
            TotInvVal: taxableValue + igst + cgst + sgst,
          },
          ItemList: items.map((it) => ({
            SlNo: String(it.slNo),
            PrdDesc: it.description,
            HsnCd: it.hsnCode,
            Qty: it.quantity,
            Unit: it.unit,
            UnitPrice: it.unitPrice,
            TotAmt: it.totalAmount,
            GstRt: it.gstRate,
          })),
        },
        null,
        2
      ),
      qrCodeData: `upi://pay?irn=${generateIRN(i)}&ack=${1e11 + i + 1}`,
      items,
    });
  }

  return records;
}

export function generateMockSourceInvoices(): SourceInvoice[] {
  const invoices: SourceInvoice[] = [];
  for (let i = 0; i < 20; i++) {
    const supplier = SUPPLIERS[i % SUPPLIERS.length];
    const buyer = BUYERS[i % BUYERS.length];
    const items = generateItems(1 + (i % 4));
    const taxableValue = items.reduce((s, it) => s + it.totalAmount, 0);
    const isInterstate = supplier.gstin.slice(0, 2) !== buyer.gstin.slice(0, 2);
    invoices.push({
      id: `src-inv-${i + 1}`,
      invoiceNo: `INV-2026-${String(2000 + i).padStart(5, "0")}`,
      invoiceDate: formatDateDDMMYYYY(
        new Date(Date.now() - i * 86400000)
      ),
      gstin: supplier.gstin,
      buyerGstin: buyer.gstin,
      buyerName: buyer.name,
      totalValue:
        taxableValue +
        (isInterstate
          ? Math.round(taxableValue * 0.18)
          : Math.round(taxableValue * 0.18)),
      taxableValue,
      igstAmount: isInterstate ? Math.round(taxableValue * 0.18) : 0,
      cgstAmount: isInterstate ? 0 : Math.round(taxableValue * 0.09),
      sgstAmount: isInterstate ? 0 : Math.round(taxableValue * 0.09),
      items,
    });
  }
  return invoices;
}

export function generateAuditTrail(record: EInvoiceRecord): AuditEntry[] {
  const entries: AuditEntry[] = [
    {
      action: "Invoice created in Aura",
      actor: "System (Aura Sync)",
      timestamp: new Date(
        new Date(record.generatedAt).getTime() - 7200000
      ).toISOString(),
      detail: `Source: ${record.invoiceNo}`,
      status: "info",
    },
    {
      action: "Payload validated",
      actor: "System",
      timestamp: new Date(
        new Date(record.generatedAt).getTime() - 300000
      ).toISOString(),
      detail: "GSTIN, HSN, totals verified",
      status: "success",
    },
    {
      action: "IRN generated",
      actor: "finance@acme.com",
      timestamp: record.generatedAt,
      detail: `ACK: ${record.ackNumber}`,
      status: "success",
    },
  ];

  if (record.status === "CANCELLED" && record.cancelledAt) {
    entries.push({
      action: "IRN cancelled",
      actor: "finance@acme.com",
      timestamp: record.cancelledAt,
      detail: record.cancelReason ?? "Cancelled",
      status: "danger",
    });
  }

  return entries;
}

function formatDateDDMMYYYY(d: Date): string {
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  const yyyy = d.getFullYear();
  return `${dd}/${mm}/${yyyy}`;
}
