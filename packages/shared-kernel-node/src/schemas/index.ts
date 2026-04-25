import { z } from 'zod';

// ---------------------------------------------------------------------------
// Tenant
// ---------------------------------------------------------------------------

export const tenantSchema = z.object({
  id: z.string().min(1),
  name: z.string().min(1),
  tier: z.enum(['pooled', 'bridge', 'silo', 'on_prem']),
  status: z.enum(['active', 'suspended', 'deactivated']),
  createdAt: z.coerce.date(),
});

// ---------------------------------------------------------------------------
// User
// ---------------------------------------------------------------------------

export const userSchema = z.object({
  id: z.string().min(1),
  tenantId: z.string().min(1),
  email: z.string().email(),
  name: z.string().min(1),
  roles: z.array(z.string()),
  permissions: z.array(z.string()),
  activePan: z.string().nullable(),
  activeGstin: z.string().nullable(),
  activeTan: z.string().nullable(),
});

// ---------------------------------------------------------------------------
// Vendor
// ---------------------------------------------------------------------------

export const vendorSchema = z.object({
  id: z.string().min(1),
  tenantId: z.string().min(1),
  name: z.string().min(1),
  gstin: z.string().min(1),
  pan: z.string().min(1),
  complianceScore: z.number().min(0).max(100),
  category: z.enum(['A', 'B', 'C', 'D']),
  status: z.enum(['active', 'inactive', 'blacklisted']),
});

// ---------------------------------------------------------------------------
// Line item
// ---------------------------------------------------------------------------

const taxComponentSchema = z.object({
  rate: z.number().min(0),
  amount: z.number().min(0),
});

export const lineItemSchema = z.object({
  itemId: z.string().min(1),
  description: z.string(),
  hsn: z.string().min(1),
  unit: z.string().min(1),
  quantity: z.number().positive(),
  unitPrice: z.number().min(0),
  discount: z.number().min(0),
  taxableValue: z.number().min(0),
  cgst: taxComponentSchema,
  sgst: taxComponentSchema,
  igst: taxComponentSchema,
  cess: taxComponentSchema,
});

// ---------------------------------------------------------------------------
// Invoice (canonical — architecture §10)
// ---------------------------------------------------------------------------

const partySchema = z.object({
  gstin: z.string(),
  name: z.string().min(1),
  address: z.string(),
  stateCode: z.string().length(2),
});

const invoiceTotalsSchema = z.object({
  taxable: z.number().min(0),
  cgst: z.number().min(0),
  sgst: z.number().min(0),
  igst: z.number().min(0),
  cess: z.number().min(0),
  roundOff: z.number(),
  grandTotal: z.number(),
});

const invoicePaymentSchema = z.object({
  mode: z.string(),
  bankDetails: z.string(),
});

const invoiceReferencesSchema = z.object({
  poNumber: z.string().nullable(),
  grnNumber: z.string().nullable(),
  contractNumber: z.string().nullable(),
  irn: z.string().nullable(),
  ewbNumber: z.string().nullable(),
});

const invoiceMetadataSchema = z.object({
  sourceSystem: z.string(),
  sourceDocumentId: z.string(),
  createdBy: z.string(),
  createdAt: z.coerce.date(),
  tags: z.array(z.string()),
});

export const invoiceSchema = z.object({
  tenantId: z.string().min(1),
  pan: z.string().min(1),
  gstin: z.string().min(1),
  id: z.string().min(1),
  documentNumber: z.string().min(1),
  documentDate: z.coerce.date(),
  supplyType: z.enum(['B2B', 'B2C', 'EXPORT', 'SEZ', 'NIL', 'DEEMED_EXPORT']),
  documentType: z.enum(['INV', 'CRN', 'DBN']),
  reverseCharge: z.boolean(),
  supplier: partySchema,
  buyer: partySchema,
  lineItems: z.array(lineItemSchema).min(1),
  totals: invoiceTotalsSchema,
  payment: invoicePaymentSchema,
  references: invoiceReferencesSchema,
  metadata: invoiceMetadataSchema,
});
