import type { PaymentCodeInfo, Section2025 } from "./types";

export const PAYMENT_CODES: PaymentCodeInfo[] = [
  { code: "1001", label: "Salary — Government", section: "392", subClause: "", baseRate: 0, description: "Salary to government employees (slab rates)" },
  { code: "1002", label: "Salary — Private", section: "392", subClause: "", baseRate: 0, description: "Salary to non-government employees (slab rates)" },
  { code: "1003", label: "Salary — PSU", section: "392", subClause: "", baseRate: 0, description: "Salary to PSU employees (slab rates)" },
  { code: "1004", label: "Salary — Others", section: "392", subClause: "", baseRate: 0, description: "Salary to other employees (slab rates)" },
  { code: "1008", label: "Insurance Commission", section: "393(1)", subClause: "Sl.1(a)", baseRate: 5, description: "Insurance commission to resident" },
  { code: "1009", label: "Rent — Land/Building", section: "393(1)", subClause: "Sl.6(i).A", baseRate: 10, description: "Rent on land, building, or furniture" },
  { code: "1010", label: "Rent — Plant/Machinery", section: "393(1)", subClause: "Sl.6(i).B", baseRate: 2, description: "Rent on plant, machinery, or equipment" },
  { code: "1023", label: "Contractor — HUF/Individual", section: "393(1)", subClause: "Sl.6(i).D(a)", baseRate: 1, description: "Payment to individual/HUF contractor" },
  { code: "1024", label: "Contractor — Other", section: "393(1)", subClause: "Sl.6(i).D(b)", baseRate: 2, description: "Payment to non-individual contractor" },
  { code: "1026", label: "Commission/Brokerage", section: "393(1)", subClause: "Sl.6(i).C", baseRate: 5, description: "Commission or brokerage payments" },
  { code: "1027", label: "Professional Fees", section: "393(1)", subClause: "Sl.6(i).E(a)", baseRate: 10, description: "Fees for professional services" },
  { code: "1028", label: "Technical Fees", section: "393(1)", subClause: "Sl.6(i).E(b)", baseRate: 2, description: "Fees for technical services" },
  { code: "1031", label: "Purchase of Goods", section: "393(1)", subClause: "Sl.6(i).F", baseRate: 0.1, description: "Purchase of goods exceeding threshold" },
  { code: "1057", label: "Non-Resident Payment", section: "393(2)", subClause: "Sl.17", baseRate: 20, description: "Payment to non-resident (rates per DTAA or Act)" },
];

export const PAYMENT_CODE_MAP = new Map(PAYMENT_CODES.map((pc) => [pc.code, pc]));

export function getPaymentCodeInfo(code: string): PaymentCodeInfo | undefined {
  return PAYMENT_CODE_MAP.get(code);
}

export function getPaymentCodesForSection(section: Section2025): PaymentCodeInfo[] {
  return PAYMENT_CODES.filter((pc) => pc.section === section);
}

export const SECTION_LABELS: Record<Section2025, string> = {
  "392": "Section 392 — Salary",
  "393(1)": "Section 393(1) — Resident Non-Salary",
  "393(2)": "Section 393(2) — Non-Resident",
  "393(3)": "Section 393(3) — TCS",
};
