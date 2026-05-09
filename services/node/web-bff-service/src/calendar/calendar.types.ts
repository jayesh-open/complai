export type EventCategory = 'direct_tax' | 'indirect_tax' | 'statutory';
export type Authority = 'CBDT' | 'CBIC' | 'MCA' | 'EPFO' | 'ESIC' | 'GSTN';

export type EventStatus =
  | 'filed'
  | 'filed_late'
  | 'overdue'
  | 'due_soon'
  | 'upcoming';

export interface ComplianceEvent {
  id: string;
  title: string;
  description: string;
  category: EventCategory;
  authority: Authority;
  sectionRef?: string;
  formRef?: string;
  dueDate: string; // ISO date YYYY-MM-DD
  penalty?: string;
  status: EventStatus;
  linkedModule?: string;
  eventType: string;
}

export type FilingScheme = 'monthly' | 'qrmp';
export type BusinessType = 'company' | 'llp' | 'proprietorship' | 'trust';

export interface TenantConfig {
  tenantId: string;
  gstins: string[];
  pans: string[];
  filingScheme: FilingScheme;
  businessType: BusinessType;
  annualTurnover: number; // in INR
  requiresAudit: boolean;
  registrationDate?: string; // ISO date
}

export const DEFAULT_TENANT_CONFIG: TenantConfig = {
  tenantId: 'default',
  gstins: ['29AABCU9603R1ZM'],
  pans: ['AABCU9603R'],
  filingScheme: 'monthly',
  businessType: 'company',
  annualTurnover: 10_00_00_000, // ₹10 Cr
  requiresAudit: true,
};

export interface EventTemplate {
  eventType: string;
  title: string;
  titleSuffix?: (date: Date) => string;
  description: string;
  category: EventCategory;
  authority: Authority;
  sectionRef?: string;
  formRef?: string;
  penalty?: string;
  linkedModule?: string;
  schedule: EventSchedule;
  filter?: (config: TenantConfig) => boolean;
}

export type EventSchedule =
  | { type: 'monthly'; dayOfMonth: number }
  | { type: 'quarterly'; months: number[]; dayOfMonth: number }
  | { type: 'annual'; month: number; dayOfMonth: number }
  | { type: 'half_yearly'; months: number[]; dayOfMonth: number };
