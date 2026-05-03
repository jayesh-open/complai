import type {
  EmployeeITRDetail,
  IncomeHeadDetail,
  DeductionItem,
  TaxComputation,
  AISMismatch,
  TaxSlab,
} from "./types";
import { ALL_EMPLOYEES } from "./mock-data";

const NEW_REGIME_SLABS: Omit<TaxSlab, "tax">[] = [
  { from: 0, to: 400000, rate: 0 },
  { from: 400000, to: 800000, rate: 5 },
  { from: 800000, to: 1200000, rate: 10 },
  { from: 1200000, to: 1600000, rate: 15 },
  { from: 1600000, to: 2000000, rate: 20 },
  { from: 2000000, to: 2400000, rate: 25 },
  { from: 2400000, to: null, rate: 30 },
];

function computeSlabs(taxableIncome: number): TaxSlab[] {
  return NEW_REGIME_SLABS.map((slab) => {
    const upper = slab.to ?? taxableIncome;
    if (taxableIncome <= slab.from) return { ...slab, tax: 0 };
    const taxable = Math.min(taxableIncome, upper) - slab.from;
    return { ...slab, tax: Math.round(taxable * slab.rate / 100) };
  });
}

function buildIncomeHeads(grossIncome: number, hasBusinessIncome: boolean): IncomeHeadDetail[] {
  const salary = Math.round(grossIncome * 0.7);
  const salaryExemptions = Math.round(salary * 0.08);

  return [
    {
      head: "SALARY",
      gross: salary,
      deductions: salaryExemptions,
      net: salary - salaryExemptions,
      subItems: [
        { label: "Basic Pay", amount: Math.round(salary * 0.5) },
        { label: "HRA", amount: Math.round(salary * 0.2) },
        { label: "Special Allowance", amount: Math.round(salary * 0.2) },
        { label: "Other Allowances", amount: Math.round(salary * 0.1) },
        { label: "HRA Exemption", amount: -Math.round(salaryExemptions * 0.6), section: "§392" },
        { label: "LTA Exemption", amount: -Math.round(salaryExemptions * 0.4) },
      ],
      visible: true,
    },
    {
      head: "HOUSE_PROPERTY",
      gross: Math.round(grossIncome * 0.05),
      deductions: Math.round(grossIncome * 0.015),
      net: Math.round(grossIncome * 0.035),
      subItems: [
        { label: "Self-Occupied — NAV", amount: 0 },
        { label: "Let-Out — Gross Rental", amount: Math.round(grossIncome * 0.05) },
        { label: "Municipal Tax", amount: -Math.round(grossIncome * 0.005) },
        { label: "Std Deduction (30%)", amount: -Math.round(grossIncome * 0.01) },
      ],
      visible: true,
    },
    {
      head: "CAPITAL_GAINS",
      gross: Math.round(grossIncome * 0.08),
      deductions: 0,
      net: Math.round(grossIncome * 0.08),
      subItems: [
        { label: "STCG (Listed Equity)", amount: Math.round(grossIncome * 0.03) },
        { label: "LTCG (§112A — up to ₹1.25L exempt)", amount: Math.round(grossIncome * 0.04), section: "§112A" },
        { label: "Schedule VDA (Crypto)", amount: Math.round(grossIncome * 0.01) },
      ],
      visible: true,
    },
    {
      head: "BUSINESS_PROFESSION",
      gross: hasBusinessIncome ? Math.round(grossIncome * 0.12) : 0,
      deductions: hasBusinessIncome ? Math.round(grossIncome * 0.04) : 0,
      net: hasBusinessIncome ? Math.round(grossIncome * 0.08) : 0,
      subItems: hasBusinessIncome
        ? [
            { label: "Gross Receipts", amount: Math.round(grossIncome * 0.12) },
            { label: "Presumptive (§44AD)", amount: -Math.round(grossIncome * 0.04), section: "§44AD" },
          ]
        : [],
      visible: hasBusinessIncome,
    },
    {
      head: "OTHER_SOURCES",
      gross: Math.round(grossIncome * 0.05),
      deductions: 0,
      net: Math.round(grossIncome * 0.05),
      subItems: [
        { label: "Savings Interest", amount: Math.round(grossIncome * 0.02) },
        { label: "FD Interest", amount: Math.round(grossIncome * 0.015) },
        { label: "Dividend Income", amount: Math.round(grossIncome * 0.01) },
        { label: "Gift Income", amount: Math.round(grossIncome * 0.005) },
      ],
      visible: true,
    },
  ];
}

function buildDeductions(grossIncome: number): DeductionItem[] {
  return [
    { section: "80C", label: "PPF, ELSS, LIC, Tuition", declared: Math.min(150000, Math.round(grossIncome * 0.06)), limit: 150000 },
    { section: "80D", label: "Health Insurance Premium", declared: Math.min(25000, Math.round(grossIncome * 0.01)), limit: 25000 },
    { section: "80E", label: "Education Loan Interest", declared: Math.round(grossIncome * 0.005), limit: 0 },
    { section: "24(b)", label: "Home Loan Interest", declared: Math.min(200000, Math.round(grossIncome * 0.05)), limit: 200000 },
    { section: "80CCD(1B)", label: "NPS Contribution", declared: Math.min(50000, Math.round(grossIncome * 0.02)), limit: 50000 },
    { section: "80TTA", label: "Savings Interest Deduction", declared: Math.min(10000, Math.round(grossIncome * 0.005)), limit: 10000 },
  ];
}

function buildComputation(incomeHeads: IncomeHeadDetail[], regime: "NEW" | "OLD", deductions: DeductionItem[]): TaxComputation {
  const totalIncome = incomeHeads.reduce((sum, h) => sum + h.net, 0);
  const standardDeduction = 75000;
  const deductionTotal = regime === "OLD" ? deductions.reduce((sum, d) => sum + d.declared, 0) : 0;
  const taxableIncome = Math.max(0, totalIncome - standardDeduction - deductionTotal);

  const slabs = computeSlabs(taxableIncome);
  const slabTax = slabs.reduce((sum, s) => sum + s.tax, 0);

  let surchargeRate = 0;
  let surchargeThreshold = "N/A";
  if (taxableIncome > 50000000) { surchargeRate = 25; surchargeThreshold = "> ₹5 Cr"; }
  else if (taxableIncome > 20000000) { surchargeRate = 25; surchargeThreshold = "> ₹2 Cr"; }
  else if (taxableIncome > 10000000) { surchargeRate = 15; surchargeThreshold = "> ₹1 Cr"; }
  else if (taxableIncome > 5000000) { surchargeRate = 10; surchargeThreshold = "> ₹50 L"; }

  const surchargeAmount = Math.round(slabTax * surchargeRate / 100);
  const taxPlusSurcharge = slabTax + surchargeAmount;
  const healthEducationCess = Math.round(taxPlusSurcharge * 4 / 100);
  const grossTax = taxPlusSurcharge + healthEducationCess;

  const rebate87A = regime === "NEW" && taxableIncome <= 700000 ? Math.min(grossTax, 25000) : 0;
  const totalLiability = Math.max(0, grossTax - rebate87A);

  const tdsCredit = Math.round(totalIncome * 0.08);
  const advanceTax = Math.round(totalIncome * 0.01);
  const selfAssessmentTax = 0;
  const refundOrPayable = totalLiability - tdsCredit - advanceTax - selfAssessmentTax;

  return {
    totalIncome, standardDeduction, taxableIncome, slabs, slabTax,
    surchargeRate, surchargeAmount, surchargeThreshold,
    healthEducationCess, grossTax, rebate87A, totalLiability,
    tdsCredit, advanceTax, selfAssessmentTax, refundOrPayable,
  };
}

function buildMismatches(employeeIndex: number): AISMismatch[] {
  const base: AISMismatch[] = [
    { id: "mis-1", category: "SALARY", field: "Gross Salary", itrValue: 1850000, aisValue: 1845000, severity: "info", resolved: true, resolution: "Rounding difference — accepted" },
    { id: "mis-2", category: "TDS", field: "TDS on Salary", itrValue: 185000, aisValue: 192000, severity: "warn", resolved: false },
    { id: "mis-3", category: "INTEREST", field: "FD Interest (HDFC)", itrValue: 45000, aisValue: 52000, severity: "error", resolved: false },
    { id: "mis-4", category: "DIVIDEND", field: "Dividend — Reliance", itrValue: 12000, aisValue: 12000, severity: "info", resolved: true, resolution: "Exact match" },
    { id: "mis-5", category: "SECURITIES", field: "Equity Sale Proceeds", itrValue: 320000, aisValue: 318000, severity: "warn", resolved: false },
  ];
  if (employeeIndex % 3 === 0) {
    base.push({ id: "mis-6", category: "PROPERTY", field: "Rental Income", itrValue: 180000, aisValue: 240000, severity: "error", resolved: false });
  }
  return base;
}

const FORM_REASONS: Record<string, string> = {
  "ITR-1": "Salary + one house property + other sources (≤ ₹50L total)",
  "ITR-2": "Capital gains / multiple house properties / foreign assets",
  "ITR-3": "Business or professional income present",
  "ITR-4": "Presumptive income under §44AD / §44ADA",
  "ITR-5": "Partnership firm / LLP",
  "ITR-6": "Company (other than §11 exempt)",
  "ITR-7": "Trust / political party / institution",
};

export function getEmployeeDetail(employeeId: string): EmployeeITRDetail | null {
  const idx = ALL_EMPLOYEES.findIndex((e) => e.id === employeeId);
  if (idx === -1) return null;
  const emp = ALL_EMPLOYEES[idx];

  const hasBusinessIncome = emp.recommendedForm === "ITR-3" || emp.recommendedForm === "ITR-4";
  const incomeHeads = buildIncomeHeads(emp.grossIncome, hasBusinessIncome);
  const deductions = buildDeductions(emp.grossIncome);
  const computation = buildComputation(incomeHeads, emp.regime, deductions);
  const mismatches = buildMismatches(idx);

  return {
    employee: emp,
    formReason: FORM_REASONS[emp.recommendedForm] ?? "",
    incomeHeads,
    deductions,
    computation,
    mismatches,
    form10IEAFiled: emp.regime === "OLD" ? idx % 2 === 0 : undefined,
    auditTrail: [
      { action: "AIS Data Fetched", actor: "System", timestamp: "2026-06-15T10:30:00Z", status: "success" as const },
      { action: "Form Generated", actor: "System", timestamp: "2026-06-16T09:00:00Z", detail: `${emp.recommendedForm} auto-generated`, status: "info" as const },
      { action: "Review Assigned", actor: "Shruti Kapoor", timestamp: "2026-06-18T14:00:00Z", status: "info" as const },
      ...(emp.filingStatus === "FILED" || emp.filingStatus === "ACKNOWLEDGED"
        ? [{ action: "ITR Filed", actor: "System", timestamp: "2026-07-10T16:00:00Z", detail: `ARN: CPC/2627/${String(idx + 1).padStart(6, "0")}`, status: "success" as const }]
        : []),
    ],
  };
}
