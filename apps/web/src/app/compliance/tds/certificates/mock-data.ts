import { ALL_DEDUCTEES, ALL_ENTRIES } from "../mock-data";
import { PAYMENT_CODE_MAP } from "../payment-codes";
import type {
  Form130Row, Form130Detail, Form131Row, Form131Detail,
  CertificateStatus, TaxRegime, SalaryBreakup, Deduction80C,
  ChallanRef, SectionBreakdownRow,
} from "./types";

const REGIMES: TaxRegime[] = ["NEW", "OLD", "NEW"];
const CERT_STATUSES: CertificateStatus[] = ["GENERATED", "PENDING", "ISSUED"];

function bsr(i: number): string {
  return `000${(1200 + i).toString()}`.slice(-7);
}

export function generateForm130Rows(): Form130Row[] {
  const salaryDeductees = ALL_DEDUCTEES.filter((d) => d.sectionPreference === "392");
  return salaryDeductees.map((d, i) => {
    const entries = ALL_ENTRIES.filter((e) => e.deducteeId === d.id);
    const totalTds = entries.reduce((a, e) => a + e.totalTax, 0);
    const grossSalary = entries.reduce((a, e) => a + e.grossAmount, 0) || 1200000;
    return {
      employeeId: d.id, employeeName: d.name, pan: d.pan,
      grossSalary, totalTds: totalTds || 120000,
      regime: REGIMES[i % REGIMES.length],
      status: CERT_STATUSES[i % CERT_STATUSES.length],
      taxYear: "2026-27",
    };
  });
}

export function generateForm130Detail(employeeId: string, taxYear: string): Form130Detail {
  const ded = ALL_DEDUCTEES.find((d) => d.id === employeeId);
  const entries = ALL_ENTRIES.filter((e) => e.deducteeId === employeeId);
  const name = ded?.name ?? "Employee";
  const pan = ded?.pan ?? "ABCDE1234F";
  const totalTds = entries.reduce((a, e) => a + e.totalTax, 0) || 120000;
  const breakup: SalaryBreakup[] = [
    { component: "Basic Salary", amount: 720000 },
    { component: "HRA", amount: 288000 },
    { component: "Special Allowance", amount: 192000 },
    { component: "Bonus", amount: 100000 },
  ];
  const totalGross = breakup.reduce((a, b) => a + b.amount, 0);
  const regime: TaxRegime = employeeId === "ded-0003" ? "OLD" : "NEW";
  const deductions: Deduction80C[] = regime === "OLD" ? [
    { section: "80C", description: "PPF + ELSS", declared: 150000, verified: 150000 },
    { section: "80D", description: "Medical Insurance", declared: 25000, verified: 25000 },
    { section: "80TTA", description: "Savings Interest", declared: 10000, verified: 10000 },
  ] : [];
  const totalDeductions = deductions.reduce((a, d) => a + d.verified, 0);
  const challans: ChallanRef[] = [
    { bsrCode: bsr(0), challanSerial: "CHL00001000", depositDate: "07/05/2026", amount: Math.round(totalTds * 0.25), quarter: "Q1" },
    { bsrCode: bsr(1), challanSerial: "CHL00001001", depositDate: "07/08/2026", amount: Math.round(totalTds * 0.25), quarter: "Q2" },
    { bsrCode: bsr(2), challanSerial: "CHL00001002", depositDate: "07/11/2026", amount: Math.round(totalTds * 0.25), quarter: "Q3" },
    { bsrCode: bsr(3), challanSerial: "CHL00001003", depositDate: "07/02/2027", amount: Math.round(totalTds * 0.25), quarter: "Q4" },
  ];

  return {
    employeeId, employeeName: name, pan, designation: "Senior Engineer",
    taxYear, regime, employerName: "Complai Technologies Pvt Ltd", employerTan: "MUMA12345B",
    salaryBreakup: breakup, totalGross, standardDeduction: 75000,
    deductions80C: deductions, totalDeductions,
    taxableIncome: totalGross - 75000 - totalDeductions,
    taxComputed: totalTds,
    tdsChallanRefs: challans, totalTdsDeducted: totalTds,
    status: "GENERATED",
  };
}

export function generateForm131Rows(): Form131Row[] {
  const nonSalary = ALL_DEDUCTEES.filter((d) => d.sectionPreference !== "392");
  const quarters = ["Q1", "Q2"];
  const rows: Form131Row[] = [];
  for (const d of nonSalary) {
    for (const q of quarters) {
      const entries = ALL_ENTRIES.filter((e) => e.deducteeId === d.id && e.quarter === q);
      if (entries.length === 0) continue;
      const codes = [...new Set(entries.map((e) => e.paymentCode))];
      rows.push({
        deducteeId: d.id, deducteeName: d.name, pan: d.pan,
        sectionCodes: codes,
        totalAmount: entries.reduce((a, e) => a + e.grossAmount, 0),
        totalTds: entries.reduce((a, e) => a + e.totalTax, 0),
        status: q === "Q1" ? "ISSUED" : "PENDING",
        taxYear: "2026-27", quarter: q,
      });
    }
  }
  return rows;
}

export function generateForm131Detail(deducteeId: string, taxYear: string, quarter: string): Form131Detail {
  const ded = ALL_DEDUCTEES.find((d) => d.id === deducteeId);
  const entries = ALL_ENTRIES.filter((e) => e.deducteeId === deducteeId && e.quarter === quarter.toUpperCase());
  const qLabel = quarter.toUpperCase() === "Q1" ? "Q1 (Apr – Jun)" : quarter.toUpperCase() === "Q2" ? "Q2 (Jul – Sep)" : quarter.toUpperCase();

  const breakdown: SectionBreakdownRow[] = entries.map((e) => ({
    paymentCode: e.paymentCode,
    paymentLabel: PAYMENT_CODE_MAP.get(e.paymentCode)?.label ?? e.paymentCode,
    amount: e.grossAmount, tdsRate: e.baseRate, tdsAmount: e.tdsAmount,
    surcharge: e.surcharge, cess: e.cess,
    dateOfPayment: e.transactionDate,
  }));

  const challans: ChallanRef[] = entries
    .filter((e) => e.challanNumber)
    .map((e, i) => ({
      bsrCode: bsr(10 + i), challanSerial: e.challanNumber!, depositDate: e.challanDate!,
      amount: e.totalTax, quarter: e.quarter,
    }));

  return {
    deducteeId, deducteeName: ded?.name ?? "Deductee", pan: ded?.pan ?? "AAACT1234A",
    category: ded?.category ?? "COMPANY",
    taxYear, quarter: quarter.toUpperCase(), quarterLabel: qLabel,
    deductorName: "Complai Technologies Pvt Ltd", deductorTan: "MUMA12345B",
    sectionBreakdown: breakdown,
    totalAmount: entries.reduce((a, e) => a + e.grossAmount, 0),
    totalTds: entries.reduce((a, e) => a + e.tdsAmount, 0),
    totalSurcharge: entries.reduce((a, e) => a + e.surcharge, 0),
    totalCess: entries.reduce((a, e) => a + e.cess, 0),
    totalTax: entries.reduce((a, e) => a + e.totalTax, 0),
    challanRefs: challans,
    status: quarter.toUpperCase() === "Q1" ? "ISSUED" : "PENDING",
  };
}
