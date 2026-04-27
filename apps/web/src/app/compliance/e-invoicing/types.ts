export type IRNStatus = "GENERATED" | "CANCELLED";

export interface EInvoiceRecord {
  id: string;
  invoiceNo: string;
  invoiceDate: string;
  gstin: string;
  buyerGstin: string;
  buyerName: string;
  irn: string;
  ackNumber: string;
  ackDate: string;
  status: IRNStatus;
  totalValue: number;
  taxableValue: number;
  igstAmount: number;
  cgstAmount: number;
  sgstAmount: number;
  generatedAt: string;
  cancelledAt?: string;
  cancelReason?: string;
  signedInvoice: string;
  qrCodeData: string;
  items: EInvoiceItem[];
}

export interface EInvoiceItem {
  slNo: number;
  description: string;
  hsnCode: string;
  quantity: number;
  unit: string;
  unitPrice: number;
  totalAmount: number;
  gstRate: number;
}

export interface AuditEntry {
  action: string;
  actor: string;
  timestamp: string;
  detail?: string;
  status: "success" | "warning" | "info" | "danger" | "default";
}

export interface SourceInvoice {
  id: string;
  invoiceNo: string;
  invoiceDate: string;
  gstin: string;
  buyerGstin: string;
  buyerName: string;
  totalValue: number;
  taxableValue: number;
  igstAmount: number;
  cgstAmount: number;
  sgstAmount: number;
  items: EInvoiceItem[];
  selected?: boolean;
}
