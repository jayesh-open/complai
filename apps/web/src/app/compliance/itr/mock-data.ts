import type {
  ITREmployee,
  BulkBatch,
  ITRForm,
  TaxRegime,
  FilingStatus,
  BatchStatus,
} from "./types";

function seededRand(seed: number): () => number {
  let s = seed;
  return () => {
    s = (s * 1103515245 + 12345) & 0x7fffffff;
    return s / 0x7fffffff;
  };
}

const EMPLOYEES: { name: string; pan: string; dept: string; desig: string }[] = [
  { name: "Rajesh Kumar", pan: "ABCPK1234E", dept: "Engineering", desig: "Senior Engineer" },
  { name: "Priya Sharma", pan: "BCDPS5678F", dept: "Engineering", desig: "Staff Engineer" },
  { name: "Amit Jain", pan: "GHIPA5678K", dept: "Finance", desig: "Finance Manager" },
  { name: "Sunita Devi", pan: "JKLPI7890N", dept: "HR", desig: "HR Lead" },
  { name: "Vikram Reddy", pan: "MNOPV2345R", dept: "Engineering", desig: "Engineering Manager" },
  { name: "Neha Gupta", pan: "QRSPN6789S", dept: "Product", desig: "Product Manager" },
  { name: "Ravi Patel", pan: "TUVPR0123T", dept: "Engineering", desig: "Software Engineer" },
  { name: "Ananya Singh", pan: "WXYAS4567U", dept: "Design", desig: "Senior Designer" },
  { name: "Mohit Verma", pan: "BCDEM8901V", dept: "Engineering", desig: "Principal Engineer" },
  { name: "Kavya Nair", pan: "FGHPK2345W", dept: "Legal", desig: "Legal Counsel" },
  { name: "Deepak Mehta", pan: "IJKLD6789X", dept: "Sales", desig: "Sales Director" },
  { name: "Pooja Iyer", pan: "MNOPB0123Y", dept: "Marketing", desig: "Marketing Lead" },
  { name: "Arjun Saxena", pan: "QRSPA4567Z", dept: "Engineering", desig: "Tech Lead" },
  { name: "Shruti Kapoor", pan: "TUVPS8901A", dept: "Finance", desig: "CFO" },
  { name: "Rahul Joshi", pan: "WXYZR2345B", dept: "Engineering", desig: "Junior Engineer" },
  { name: "Meera Rao", pan: "BCDEM6789C", dept: "Engineering", desig: "QA Engineer" },
  { name: "Karthik Iyer", pan: "FGHPK0123D", dept: "Product", desig: "APM" },
  { name: "Divya Tiwari", pan: "IJKLD4567E", dept: "HR", desig: "Recruiter" },
  { name: "Sanjay Kumar", pan: "MNOPS8901F", dept: "Engineering", desig: "DevOps Engineer" },
  { name: "Rekha Menon", pan: "QRSPR2345G", dept: "Finance", desig: "Accounts Manager" },
];

const FORMS: ITRForm[] = ["ITR-1", "ITR-2", "ITR-1", "ITR-1", "ITR-2", "ITR-1", "ITR-1", "ITR-2", "ITR-3", "ITR-2", "ITR-3", "ITR-1", "ITR-2", "ITR-3", "ITR-1", "ITR-1", "ITR-1", "ITR-1", "ITR-2", "ITR-1"];
const STATUSES: FilingStatus[] = ["NOT_STARTED", "AIS_FETCHED", "FORM_GENERATED", "REVIEW_PENDING", "EMPLOYEE_APPROVED", "FILED", "ACKNOWLEDGED", "FILED"];

export function generateMockEmployees(): ITREmployee[] {
  const rand = seededRand(99);
  return EMPLOYEES.map((emp, i) => {
    const status = STATUSES[i % STATUSES.length];
    const grossIncome = 600000 + Math.floor(rand() * 3000000);
    const taxPayable = Math.floor(grossIncome * (0.05 + rand() * 0.25));
    const taxPaid = status === "FILED" || status === "ACKNOWLEDGED"
      ? taxPayable
      : Math.floor(taxPayable * (0.6 + rand() * 0.4));
    const refund = taxPaid > taxPayable ? taxPaid - taxPayable : 0;
    const regime: TaxRegime = rand() > 0.3 ? "NEW" : "OLD";

    return {
      id: `emp-itr-${String(i + 1).padStart(4, "0")}`,
      tenantId: "11111111-1111-1111-1111-111111111111",
      pan: emp.pan,
      name: emp.name,
      email: `${emp.name.toLowerCase().replace(/ /g, ".")}@company.com`,
      designation: emp.desig,
      department: emp.dept,
      taxYear: "2026-27",
      regime,
      recommendedForm: FORMS[i],
      filingStatus: status,
      grossIncome,
      taxPayable,
      taxPaid,
      refundDue: refund,
      aisReconciled: status !== "NOT_STARTED" && rand() > 0.2,
      lastUpdated: `2026-${String(3 + (i % 5)).padStart(2, "0")}-${String(10 + (i % 20)).padStart(2, "0")}`,
      acknowledgementNumber: status === "ACKNOWLEDGED" ? `CPC/2627/${String(i + 1).padStart(6, "0")}` : undefined,
      filedDate: status === "FILED" || status === "ACKNOWLEDGED" ? `2026-07-${String(10 + (i % 20)).padStart(2, "0")}` : undefined,
    };
  });
}

export function generateMockBatches(): BulkBatch[] {
  return [
    {
      id: "batch-001",
      tenantId: "11111111-1111-1111-1111-111111111111",
      name: "Tax Year 2026-27 Salary Employees — Batch 1",
      taxYear: "2026-27",
      status: "COMPLETED",
      totalEmployees: 15,
      filed: 12,
      pending: 0,
      failed: 3,
      createdAt: "2026-07-01",
      createdBy: "Shruti Kapoor",
      completedAt: "2026-07-15",
    },
    {
      id: "batch-002",
      tenantId: "11111111-1111-1111-1111-111111111111",
      name: "Tax Year 2026-27 — Senior Leadership",
      taxYear: "2026-27",
      status: "IN_PROGRESS",
      totalEmployees: 5,
      filed: 2,
      pending: 2,
      failed: 1,
      createdAt: "2026-07-10",
      createdBy: "Shruti Kapoor",
    },
    {
      id: "batch-003",
      tenantId: "11111111-1111-1111-1111-111111111111",
      name: "Tax Year 2026-27 — New Joiners",
      taxYear: "2026-27",
      status: "DRAFT",
      totalEmployees: 8,
      filed: 0,
      pending: 8,
      failed: 0,
      createdAt: "2026-07-20",
      createdBy: "Sunita Devi",
    },
  ];
}

export const ALL_EMPLOYEES = generateMockEmployees();
export const ALL_BATCHES = generateMockBatches();
