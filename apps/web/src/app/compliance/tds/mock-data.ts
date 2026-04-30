import type {
  Deductee,
  TDSEntry,
  TDSImportRow,
  ReturnStatus,
  Section2025,
  DeducteeCategory,
  ResidencyStatus,
  EntryStatus,
  FilingStatus,
} from "./types";

const NAMES: { name: string; pan: string; cat: DeducteeCategory; res: ResidencyStatus; section: Section2025; country?: string }[] = [
  { name: "Tata Consultancy Services Ltd", pan: "AAACT1234A", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
  { name: "Infosys Ltd", pan: "AABCI5678B", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
  { name: "Rajesh Kumar", pan: "ABCPK1234E", cat: "INDIVIDUAL", res: "RESIDENT", section: "392" },
  { name: "Priya Sharma", pan: "BCDPS5678F", cat: "INDIVIDUAL", res: "RESIDENT", section: "392" },
  { name: "Kumar & Associates HUF", pan: "CDEHA9012G", cat: "HUF", res: "RESIDENT", section: "393(1)" },
  { name: "Mahindra Engineering Pvt Ltd", pan: "DEFCM3456H", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
  { name: "Global Tech Solutions Inc", pan: "EFGGT7890I", cat: "COMPANY", res: "NON_RESIDENT", section: "393(2)", country: "US" },
  { name: "Singapore Consulting Pte Ltd", pan: "FGHGS1234J", cat: "COMPANY", res: "NON_RESIDENT", section: "393(2)", country: "SG" },
  { name: "Amit Jain", pan: "GHIPA5678K", cat: "INDIVIDUAL", res: "RESIDENT", section: "393(1)" },
  { name: "Reliance Industries Ltd", pan: "HIJCR9012L", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
  { name: "Patel Family Trust", pan: "IJKTP3456M", cat: "TRUST", res: "RESIDENT", section: "393(1)" },
  { name: "Sunita Devi", pan: "JKLPI7890N", cat: "INDIVIDUAL", res: "RESIDENT", section: "392" },
  { name: "London Advisory LLP", pan: "KLMFL1234O", cat: "FIRM", res: "NON_RESIDENT", section: "393(2)", country: "GB" },
  { name: "Bharat Heavy Electricals Ltd", pan: "LMNCH5678P", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
  { name: "Wipro Ltd", pan: "MNOCW9012Q", cat: "COMPANY", res: "RESIDENT", section: "393(1)" },
];

const PAYMENT_CODES_RESIDENT = ["1023", "1024", "1026", "1027", "1028", "1009", "1031"];
const PAYMENT_CODES_SALARY = ["1001", "1002", "1003"];
const STATUSES: EntryStatus[] = ["PENDING", "DEPOSITED", "FILED"];

function seededRand(seed: number): () => number {
  let s = seed;
  return () => {
    s = (s * 1103515245 + 12345) & 0x7fffffff;
    return s / 0x7fffffff;
  };
}

function formatDateDDMMYYYY(d: Date): string {
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()}`;
}

function toISODate(d: Date): string {
  return d.toISOString().split("T")[0];
}

export function generateMockDeductees(): Deductee[] {
  const rand = seededRand(42);
  return NAMES.map((n, i) => ({
    id: `ded-${String(i + 1).padStart(4, "0")}`,
    tenantId: "11111111-1111-1111-1111-111111111111",
    pan: n.pan,
    name: n.name,
    category: n.cat,
    residency: n.res,
    sectionPreference: n.section,
    countryCode: n.country,
    dtaaApplicable: n.res === "NON_RESIDENT",
    form41Filed: n.res === "NON_RESIDENT" && rand() > 0.3,
    trcAttached: n.res === "NON_RESIDENT" && rand() > 0.5,
    totalDeductedYTD: Math.round(rand() * 5000000 + 100000),
    lastTransactionDate: formatDateDDMMYYYY(
      new Date(2026, 3, 15 - Math.floor(rand() * 30))
    ),
    address: n.res === "NON_RESIDENT" ? `123 Main St, ${n.country}` : `${Math.floor(rand() * 500) + 1}, MG Road, Mumbai 400001`,
    email: `finance@${n.name.toLowerCase().replace(/[^a-z]/g, "").slice(0, 10)}.com`,
    phone: `+91 ${String(Math.floor(rand() * 9000000000 + 1000000000))}`,
  }));
}

export function generateMockEntries(deductees: Deductee[]): TDSEntry[] {
  const rand = seededRand(101);
  const entries: TDSEntry[] = [];
  const baseDate = new Date(2026, 3, 1);

  for (let i = 0; i < 60; i++) {
    const ded = deductees[i % deductees.length];
    const txDate = new Date(baseDate.getTime() - i * 86400000 * 1.5);
    let paymentCode: string;
    let subClause: string;

    if (ded.sectionPreference === "392") {
      paymentCode = PAYMENT_CODES_SALARY[i % PAYMENT_CODES_SALARY.length];
      subClause = "";
    } else if (ded.sectionPreference === "393(2)") {
      paymentCode = "1057";
      subClause = "Sl.17";
    } else {
      paymentCode = PAYMENT_CODES_RESIDENT[i % PAYMENT_CODES_RESIDENT.length];
      subClause = ["Sl.6(i).D(a)", "Sl.6(i).D(b)", "Sl.6(i).C", "Sl.6(i).E(a)", "Sl.6(i).E(b)", "Sl.6(i).A", "Sl.6(i).F"][i % 7];
    }

    const grossAmount = Math.round(rand() * 500000 + 50000);
    const baseRate = ded.sectionPreference === "392" ? 10 : ded.sectionPreference === "393(2)" ? 20 : [1, 2, 5, 10, 2, 10, 0.1][i % 7];
    const cessRate = ded.residency === "NON_RESIDENT" ? 4 : 0;
    const surchargeRate = grossAmount > 5000000 ? 10 : 0;
    const tdsAmount = Math.round(grossAmount * baseRate / 100);
    const surcharge = Math.round(tdsAmount * surchargeRate / 100);
    const cess = Math.round((tdsAmount + surcharge) * cessRate / 100);
    const status = STATUSES[Math.floor(rand() * STATUSES.length)];

    entries.push({
      id: `entry-${String(i + 1).padStart(5, "0")}`,
      tenantId: ded.tenantId,
      deducteeId: ded.id,
      deducteePan: ded.pan,
      deducteeName: ded.name,
      section2025: ded.sectionPreference,
      paymentCode,
      subClause,
      taxYear: "2026-27",
      quarter: txDate.getMonth() < 6 ? "Q1" : txDate.getMonth() < 9 ? "Q2" : "Q3",
      transactionDate: toISODate(txDate),
      paymentDate: status !== "PENDING" ? toISODate(new Date(txDate.getTime() + 7 * 86400000)) : undefined,
      grossAmount,
      baseRate,
      cessRate,
      surchargeRate,
      effectiveRate: baseRate + (baseRate * surchargeRate / 100) + ((baseRate + baseRate * surchargeRate / 100) * cessRate / 100),
      tdsAmount,
      surcharge,
      cess,
      totalTax: tdsAmount + surcharge + cess,
      challanNumber: status !== "PENDING" ? `CHL${String(i + 1000).padStart(8, "0")}` : undefined,
      challanDate: status !== "PENDING" ? toISODate(new Date(txDate.getTime() + 10 * 86400000)) : undefined,
      status,
      noPanDeduction: false,
      lowerCertApplied: rand() > 0.9,
      dtaaCountryCode: ded.countryCode,
      natureOfPayment: ded.sectionPreference === "392" ? "Salary" : ded.sectionPreference === "393(2)" ? "Technical services" : "Contractual payment",
      invoiceNumber: ded.sectionPreference !== "392" ? `INV-${String(2000 + i).padStart(6, "0")}` : undefined,
    });
  }

  return entries;
}

export function generateMockReturnStatuses(): ReturnStatus[] {
  const quarters = ["Q1", "Q2", "Q3", "Q4"];
  const forms: { type: "138" | "140" | "144"; label: string }[] = [
    { type: "138", label: "Form 138 (Salary)" },
    { type: "140", label: "Form 140 (Non-Salary Resident)" },
    { type: "144", label: "Form 144 (Non-Resident)" },
  ];
  const dueDates: Record<string, string> = {
    Q1: "31/07/2026", Q2: "31/10/2026", Q3: "31/01/2027", Q4: "31/05/2027",
  };
  const statuses: FilingStatus[] = ["FILED", "PENDING", "NOT_DUE", "NOT_DUE"];
  const results: ReturnStatus[] = [];

  for (const form of forms) {
    for (let qi = 0; qi < quarters.length; qi++) {
      const q = quarters[qi];
      const status = statuses[qi];
      results.push({
        formType: form.type,
        formLabel: form.label,
        taxYear: "2026-27",
        quarter: q,
        status,
        dueDate: dueDates[q],
        filedDate: status === "FILED" ? "28/07/2026" : undefined,
        tokenNumber: status === "FILED" ? `TKN${form.type}-${q}-001` : undefined,
        acknowledgementNumber: status === "FILED" ? `ACK${q}abcdef` : undefined,
        entryCount: status === "FILED" ? 15 : status === "PENDING" ? 12 : 0,
        totalTds: status === "FILED" ? 1250000 : status === "PENDING" ? 980000 : 0,
      });
    }
  }

  return results;
}

const ALL_DEDUCTEES = generateMockDeductees();
const ALL_ENTRIES = generateMockEntries(ALL_DEDUCTEES);
const ALL_RETURNS = generateMockReturnStatuses();

export function generateMockImportRows(): TDSImportRow[] {
  return [
    { rowNumber: 1, deducteePan: "AAACT1234A", deducteeName: "TCS Ltd", paymentCode: "1024", amount: 500000, tdsAmount: 10000, transactionDate: "15/04/2026", status: "valid", errors: [] },
    { rowNumber: 2, deducteePan: "AABCI5678B", deducteeName: "Infosys Ltd", paymentCode: "1027", amount: 300000, tdsAmount: 30000, transactionDate: "16/04/2026", status: "valid", errors: [] },
    { rowNumber: 3, deducteePan: "INVALID123", deducteeName: "Bad Co", paymentCode: "1024", amount: 100000, tdsAmount: 2000, transactionDate: "17/04/2026", status: "error", errors: ["Invalid PAN format"] },
    { rowNumber: 4, deducteePan: "ABCPK1234E", deducteeName: "Rajesh Kumar", paymentCode: "9999", amount: 200000, tdsAmount: 4000, transactionDate: "18/04/2026", status: "error", errors: ["Unknown payment code 9999"] },
    { rowNumber: 5, deducteePan: "CDEHA9012G", deducteeName: "Kumar & Associates", paymentCode: "1024", amount: 20000, tdsAmount: 400, transactionDate: "19/04/2026", status: "warning", errors: ["Amount below ₹30,000 threshold"] },
    { rowNumber: 6, deducteePan: "DEFCM3456H", deducteeName: "Mahindra Eng", paymentCode: "1009", amount: 120000, tdsAmount: 12000, transactionDate: "20/04/2026", status: "valid", errors: [] },
    { rowNumber: 7, deducteePan: "GHIPA5678K", deducteeName: "Amit Jain", paymentCode: "1026", amount: 75000, tdsAmount: 3750, transactionDate: "21/04/2026", status: "valid", errors: [] },
    { rowNumber: 8, deducteePan: "HIJCR9012L", deducteeName: "Reliance Ind", paymentCode: "1031", amount: 8000000, tdsAmount: 8000, transactionDate: "22/04/2026", status: "valid", errors: [] },
  ];
}

export { ALL_DEDUCTEES, ALL_ENTRIES, ALL_RETURNS };
