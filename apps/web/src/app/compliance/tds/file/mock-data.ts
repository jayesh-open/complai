import { ALL_DEDUCTEES, ALL_ENTRIES } from "../mock-data";
import { PAYMENT_CODE_MAP } from "../payment-codes";
import type {
  TDSFormType,
  QuarterFilingStatus,
  FilingGridCell,
  TDSFilingData,
  DeducteeSummaryRow,
  PaymentCodeDistribution,
} from "./types";
import { FORM_LABELS, QUARTERS } from "./types";

const GRID_STATUSES: Record<string, QuarterFilingStatus> = {
  "138-q1": "FILED", "138-q2": "DRAFT", "138-q3": "NOT_STARTED", "138-q4": "NOT_STARTED",
  "140-q1": "FILED", "140-q2": "SUBMITTED", "140-q3": "NOT_STARTED", "140-q4": "NOT_STARTED",
  "144-q1": "FILED", "144-q2": "DRAFT", "144-q3": "NOT_STARTED", "144-q4": "NOT_STARTED",
};

export function generateFilingGrid(): FilingGridCell[] {
  const forms: TDSFormType[] = ["138", "140", "144"];
  const cells: FilingGridCell[] = [];

  for (const form of forms) {
    for (const q of QUARTERS) {
      const key = `${form}-${q.id}`;
      const status = GRID_STATUSES[key] ?? "NOT_STARTED";
      const entryCount = status === "FILED" ? 15 : status === "DRAFT" || status === "SUBMITTED" ? 12 : 0;
      const totalTds = status === "FILED" ? 1250000 : status === "DRAFT" || status === "SUBMITTED" ? 980000 : 0;
      cells.push({ formType: form, formLabel: FORM_LABELS[form], quarter: q.id, status, dueDate: q.dueDate, entryCount, totalTds });
    }
  }
  return cells;
}

function filterEntries(formType: TDSFormType, quarter: string) {
  const sectionFilter = formType === "138" ? "392" : formType === "140" ? "393(1)" : "393(2)";
  return ALL_ENTRIES.filter((e) => e.section2025 === sectionFilter && e.quarter === quarter.toUpperCase());
}

function buildDeducteeSummaries(formType: TDSFormType, quarter: string): DeducteeSummaryRow[] {
  const entries = filterEntries(formType, quarter);
  const map = new Map<string, DeducteeSummaryRow>();

  for (const e of entries) {
    let row = map.get(e.deducteeId);
    if (!row) {
      const ded = ALL_DEDUCTEES.find((d) => d.id === e.deducteeId);
      row = {
        deducteeId: e.deducteeId, pan: e.deducteePan, name: e.deducteeName,
        category: ded?.category ?? "INDIVIDUAL", residency: ded?.residency ?? "RESIDENT",
        entryCount: 0, grossTotal: 0, tdsTotal: 0, surchargeTotal: 0, cessTotal: 0, totalTax: 0,
        paymentCodes: [], form41Filed: ded?.form41Filed, trcAttached: ded?.trcAttached,
        countryCode: ded?.countryCode,
      };
      map.set(e.deducteeId, row);
    }
    row.entryCount++;
    row.grossTotal += e.grossAmount;
    row.tdsTotal += e.tdsAmount;
    row.surchargeTotal += e.surcharge;
    row.cessTotal += e.cess;
    row.totalTax += e.totalTax;
    if (!row.paymentCodes.includes(e.paymentCode)) row.paymentCodes.push(e.paymentCode);
  }
  return Array.from(map.values());
}

function buildPaymentCodeDistribution(formType: TDSFormType, quarter: string): PaymentCodeDistribution[] {
  const entries = filterEntries(formType, quarter);
  const map = new Map<string, PaymentCodeDistribution>();

  for (const e of entries) {
    let d = map.get(e.paymentCode);
    if (!d) {
      const info = PAYMENT_CODE_MAP.get(e.paymentCode);
      d = { code: e.paymentCode, label: info?.label ?? e.paymentCode, count: 0, grossTotal: 0, tdsTotal: 0 };
      map.set(e.paymentCode, d);
    }
    d.count++;
    d.grossTotal += e.grossAmount;
    d.tdsTotal += e.tdsAmount;
  }
  return Array.from(map.values());
}

export function generateFilingData(formType: TDSFormType, taxYear: string, quarter: string): TDSFilingData {
  const entries = filterEntries(formType, quarter);
  const summaries = buildDeducteeSummaries(formType, quarter);
  const distribution = buildPaymentCodeDistribution(formType, quarter);
  const qInfo = QUARTERS.find((q) => q.id === quarter);

  const blockers: string[] = [];
  if (formType === "144") {
    for (const s of summaries) {
      if (!s.form41Filed) blockers.push(`${s.name} (${s.pan}): Form 41 not filed`);
      if (!s.trcAttached) blockers.push(`${s.name} (${s.pan}): TRC not attached`);
    }
  }

  const totalGross = entries.reduce((a, e) => a + e.grossAmount, 0);
  const totalTds = entries.reduce((a, e) => a + e.tdsAmount, 0);
  const totalSurcharge = entries.reduce((a, e) => a + e.surcharge, 0);
  const totalCess = entries.reduce((a, e) => a + e.cess, 0);

  return {
    formType, formLabel: FORM_LABELS[formType], taxYear, quarter,
    quarterLabel: qInfo ? `${qInfo.label} (${qInfo.dateRange})` : quarter,
    tan: "MUMA12345B",
    deducteeSummaries: summaries, entries, paymentCodeDistribution: distribution,
    totalGross, totalTds, totalSurcharge, totalCess,
    totalTax: totalTds + totalSurcharge + totalCess,
    deducteeCount: summaries.length, dtaaBlockers: blockers,
  };
}

export function generateFVUText(data: TDSFilingData): string {
  const lines: string[] = [
    `^FVU|${data.formType === "138" ? "FORM138" : data.formType === "140" ? "FORM140" : "FORM144"}|5.2|`,
    `^BH|${data.tan}|${data.taxYear}|${data.quarter.toUpperCase()}|ORIGINAL|`,
    `^CH|1|${data.tan}|${data.taxYear}|${data.quarter.toUpperCase()}|${data.deducteeCount}|${data.totalTds}|`,
    "",
  ];
  for (const s of data.deducteeSummaries) {
    lines.push(`^DD|${s.pan}|${s.name}|${s.entryCount}|${s.grossTotal}|${s.tdsTotal}|${s.surchargeTotal}|${s.cessTotal}|${s.totalTax}|`);
  }
  lines.push("", `^TV|TOTAL|${data.deducteeCount}|${data.totalGross}|${data.totalTds}|${data.totalSurcharge}|${data.totalCess}|${data.totalTax}|`);
  lines.push("^FV|END|");
  return lines.join("\n");
}
