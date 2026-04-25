// ---------------------------------------------------------------------------
// Branded string types
// ---------------------------------------------------------------------------

declare const __brand: unique symbol;
type Brand<T, B extends string> = T & { readonly [__brand]: B };

export type TenantId = Brand<string, 'TenantId'>;
export type UserId = Brand<string, 'UserId'>;
export type PanId = Brand<string, 'PanId'>;
export type GstinId = Brand<string, 'GstinId'>;
export type TanId = Brand<string, 'TanId'>;

// ---------------------------------------------------------------------------
// Tenant
// ---------------------------------------------------------------------------

export type TenantTier = 'pooled' | 'bridge' | 'silo' | 'on_prem';
export type TenantStatus = 'active' | 'suspended' | 'deactivated';

export interface Tenant {
  id: TenantId;
  name: string;
  tier: TenantTier;
  status: TenantStatus;
  createdAt: Date;
}

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

export interface User {
  id: UserId;
  tenantId: TenantId;
  email: string;
  name: string;
  roles: string[];
  permissions: string[];
  activePan: PanId | null;
  activeGstin: GstinId | null;
  activeTan: TanId | null;
}

// ---------------------------------------------------------------------------
// Vendor
// ---------------------------------------------------------------------------

export type VendorCategory = 'A' | 'B' | 'C' | 'D';
export type VendorStatus = 'active' | 'inactive' | 'blacklisted';

export interface Vendor {
  id: string;
  tenantId: TenantId;
  name: string;
  gstin: GstinId;
  pan: PanId;
  complianceScore: number;
  category: VendorCategory;
  status: VendorStatus;
}

// ---------------------------------------------------------------------------
// Invoice (canonical schema — architecture §10)
// ---------------------------------------------------------------------------

export interface Party {
  gstin: string;
  name: string;
  address: string;
  stateCode: string;
}

export interface TaxComponent {
  rate: number;
  amount: number;
}

export interface LineItem {
  itemId: string;
  description: string;
  hsn: string;
  unit: string;
  quantity: number;
  unitPrice: number;
  discount: number;
  taxableValue: number;
  cgst: TaxComponent;
  sgst: TaxComponent;
  igst: TaxComponent;
  cess: TaxComponent;
}

export interface InvoiceTotals {
  taxable: number;
  cgst: number;
  sgst: number;
  igst: number;
  cess: number;
  roundOff: number;
  grandTotal: number;
}

export interface InvoicePayment {
  mode: string;
  bankDetails: string;
}

export interface InvoiceReferences {
  poNumber: string | null;
  grnNumber: string | null;
  contractNumber: string | null;
  irn: string | null;
  ewbNumber: string | null;
}

export interface InvoiceMetadata {
  sourceSystem: string;
  sourceDocumentId: string;
  createdBy: UserId;
  createdAt: Date;
  tags: string[];
}

export type SupplyType = 'B2B' | 'B2C' | 'EXPORT' | 'SEZ' | 'NIL' | 'DEEMED_EXPORT';
export type DocumentType = 'INV' | 'CRN' | 'DBN';

export interface Invoice {
  tenantId: TenantId;
  pan: PanId;
  gstin: GstinId;
  id: string;
  documentNumber: string;
  documentDate: Date;
  supplyType: SupplyType;
  documentType: DocumentType;
  reverseCharge: boolean;
  supplier: Party;
  buyer: Party;
  lineItems: LineItem[];
  totals: InvoiceTotals;
  payment: InvoicePayment;
  references: InvoiceReferences;
  metadata: InvoiceMetadata;
}

// ---------------------------------------------------------------------------
// Filing status
// ---------------------------------------------------------------------------

export enum FilingStatus {
  DRAFT = 'DRAFT',
  VALIDATING = 'VALIDATING',
  VALIDATED = 'VALIDATED',
  SUBMITTING = 'SUBMITTING',
  SUBMITTED = 'SUBMITTED',
  FILED = 'FILED',
  ACKNOWLEDGED = 'ACKNOWLEDGED',
  FAILED = 'FAILED',
  AMENDMENT = 'AMENDMENT',
}

// ---------------------------------------------------------------------------
// API response wrappers
// ---------------------------------------------------------------------------

export interface ApiResponseMeta {
  requestId: string;
  latencyMs: number;
}

export interface ApiResponse<T> {
  data: T;
  meta: ApiResponseMeta;
}

export interface ApiError {
  code: string;
  message: string;
  details: Record<string, unknown>;
  retryAfterMs: number | null;
}
