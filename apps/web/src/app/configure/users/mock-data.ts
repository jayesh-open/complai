import type { User, Role } from "./types";

export const MOCK_ROLES: Role[] = [
  { id: "r1", tenant_id: "t1", name: "admin", display_name: "Admin", description: "Full platform access including user management and billing", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r2", tenant_id: "t1", name: "tax_manager", display_name: "Tax Manager", description: "Manages all tax filings and compliance workflows", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r3", tenant_id: "t1", name: "ap_manager", display_name: "AP Manager", description: "Manages accounts payable compliance", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r4", tenant_id: "t1", name: "ap_executive", display_name: "AP Executive", description: "Handles day-to-day accounts payable data entry", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r5", tenant_id: "t1", name: "ar_manager", display_name: "AR Manager", description: "Manages accounts receivable compliance", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r6", tenant_id: "t1", name: "ar_executive", display_name: "AR Executive", description: "Handles day-to-day accounts receivable data entry", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
  { id: "r7", tenant_id: "t1", name: "auditor", display_name: "Auditor", description: "Read-only access for audit purposes", is_system: true, created_at: "2026-01-01T00:00:00Z", updated_at: "2026-01-01T00:00:00Z" },
];

export const MOCK_USERS: User[] = [
  { id: "u1", email: "priya.sharma@acme.in", first_name: "Priya", last_name: "Sharma", status: "active", role: { id: "r1", name: "admin", display_name: "Admin" }, created_at: "2026-01-15T10:00:00Z", updated_at: "2026-05-01T08:30:00Z" },
  { id: "u2", email: "rajeev.kumar@acme.in", first_name: "Rajeev", last_name: "Kumar", status: "active", role: { id: "r2", name: "tax_manager", display_name: "Tax Manager" }, created_at: "2026-01-20T12:00:00Z", updated_at: "2026-04-28T14:00:00Z" },
  { id: "u3", email: "ananya.patel@acme.in", first_name: "Ananya", last_name: "Patel", status: "active", role: { id: "r3", name: "ap_manager", display_name: "AP Manager" }, created_at: "2026-02-01T09:00:00Z", updated_at: "2026-05-02T11:00:00Z" },
  { id: "u4", email: "vikram.singh@acme.in", first_name: "Vikram", last_name: "Singh", status: "active", role: { id: "r4", name: "ap_executive", display_name: "AP Executive" }, created_at: "2026-02-10T14:00:00Z", updated_at: "2026-05-03T09:30:00Z" },
  { id: "u5", email: "deepa.nair@acme.in", first_name: "Deepa", last_name: "Nair", status: "active", role: { id: "r5", name: "ar_manager", display_name: "AR Manager" }, created_at: "2026-02-15T10:30:00Z", updated_at: "2026-05-04T16:00:00Z" },
  { id: "u6", email: "arjun.reddy@acme.in", first_name: "Arjun", last_name: "Reddy", status: "inactive", role: { id: "r6", name: "ar_executive", display_name: "AR Executive" }, created_at: "2026-03-01T11:00:00Z", updated_at: "2026-04-15T10:00:00Z" },
  { id: "u7", email: "meera.iyer@acme.in", first_name: "Meera", last_name: "Iyer", status: "active", role: { id: "r7", name: "auditor", display_name: "Auditor" }, created_at: "2026-03-10T08:00:00Z", updated_at: "2026-05-05T12:00:00Z" },
  { id: "u8", email: "suresh.gupta@acme.in", first_name: "Suresh", last_name: "Gupta", status: "active", role: { id: "r2", name: "tax_manager", display_name: "Tax Manager" }, created_at: "2026-03-20T15:00:00Z", updated_at: "2026-05-06T10:00:00Z" },
  { id: "u9", email: "kavita.joshi@acme.in", first_name: "Kavita", last_name: "Joshi", status: "active", role: null, created_at: "2026-04-01T09:00:00Z", updated_at: "2026-04-01T09:00:00Z" },
  { id: "u10", email: "rohit.menon@acme.in", first_name: "Rohit", last_name: "Menon", status: "inactive", role: { id: "r4", name: "ap_executive", display_name: "AP Executive" }, created_at: "2026-04-10T13:00:00Z", updated_at: "2026-04-20T09:00:00Z" },
];
