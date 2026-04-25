// ---------------------------------------------------------------------------
// Hand-written TypeScript types matching Protobuf definitions.
// These will be replaced by buf-generated code in the future.
// ---------------------------------------------------------------------------

// common.proto
export interface Timestamp {
  seconds: number;
  nanos: number;
}

export interface UUID {
  value: string;
}

export interface Money {
  amountMinor: number; // int64 — amount in paise (INR)
  currency: string; // ISO 4217, e.g. "INR"
}

// tenant.proto
export interface TenantCreated {
  tenantId: string;
  name: string;
  tier: 'pooled' | 'bridge' | 'silo' | 'on_prem';
  createdAt: Timestamp;
}

// user.proto
export interface UserCreated {
  userId: string;
  tenantId: string;
  email: string;
  roles: string[];
}

// invoice.proto
export interface InvoiceCreated {
  tenantId: string;
  invoiceId: string;
  documentNumber: string;
  gstin: string;
  supplyType: 'B2B' | 'B2C' | 'EXPORT' | 'SEZ' | 'NIL' | 'DEEMED_EXPORT';
  grandTotal: Money;
  createdAt: Timestamp;
}

// filing.proto
export interface FilingSubmitted {
  tenantId: string;
  gstin: string;
  formType: string;
  period: string;
  status: string;
  submittedAt: Timestamp;
}

export interface FilingAcknowledged {
  tenantId: string;
  gstin: string;
  formType: string;
  period: string;
  status: string;
  acknowledgedAt: Timestamp;
  arn: string;
}

// canonical-invoice.proto
export interface Party {
  gstin: string;
  name: string;
  address: string;
  stateCode: string;
}

export interface TaxComponent {
  rate: number; // basis points
  amount: Money;
}

export interface LineItem {
  itemId: string;
  description: string;
  hsn: string;
  unit: string;
  quantity: number;
  unitPrice: Money;
  discount: Money;
  taxableValue: Money;
  cgst: TaxComponent;
  sgst: TaxComponent;
  igst: TaxComponent;
  cess: TaxComponent;
}

export interface InvoiceTotals {
  taxable: Money;
  cgst: Money;
  sgst: Money;
  igst: Money;
  cess: Money;
  roundOff: Money;
  grandTotal: Money;
}

export interface Payment {
  mode: string;
  bankDetails: string;
}

export interface References {
  poNumber: string;
  grnNumber: string;
  contractNumber: string;
  irn: string;
  ewbNumber: string;
}

export interface InvoiceMetadata {
  sourceSystem: string;
  sourceDocumentId: string;
  createdBy: string;
  createdAt: Timestamp;
  tags: string[];
}

export interface CanonicalInvoice {
  tenantId: string;
  pan: string;
  gstin: string;
  id: string;
  documentNumber: string;
  documentDate: Timestamp;
  supplyType: 'B2B' | 'B2C' | 'EXPORT' | 'SEZ' | 'NIL' | 'DEEMED_EXPORT';
  documentType: 'INV' | 'CRN' | 'DBN';
  reverseCharge: boolean;
  supplier: Party;
  buyer: Party;
  lineItems: LineItem[];
  totals: InvoiceTotals;
  payment: Payment;
  references: References;
  metadata: InvoiceMetadata;
}
