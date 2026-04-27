import type {
  EwbRecord, EwbItem, SourceInvoiceForEwb, VehicleEntry, TransportMode,
} from "./types";

const SUPPLIERS = [
  { gstin: "29AABCA1234A1Z5", name: "Acme Manufacturing Ltd" },
  { gstin: "27AABCB5678B1Z3", name: "Bharat Traders Pvt Ltd" },
  { gstin: "33AABCC9012C1Z1", name: "Chennai Exports Ltd" },
];

const CONSIGNEES = [
  { gstin: "07AABCD3456D1Z9", name: "Delhi Distributors Ltd" },
  { gstin: "24AABCE7890E1Z7", name: "Gujarat Retail Corp" },
  { gstin: "36AABCF1234F1Z5", name: "Telangana Tech Solutions" },
  { gstin: "09AABCG5678G1Z3", name: "UP Logistics Pvt Ltd" },
  { gstin: "19AABCH9012H1Z1", name: "Kolkata Trading House" },
];

const VEHICLES = [
  "KA01AB1234", "MH02CD5678", "TN03EF9012", "DL04GH3456",
  "GJ05IJ7890", "UP06KL2345", "WB07MN6789",
];

const HSN_CODES = ["84713010", "85044090", "94036090", "73269099", "39269099"];
const DESCRIPTIONS = [
  "Laptop Computer (15.6 inch)", "UPS Power Supply Unit",
  "Ergonomic Office Chair", "Steel Mounting Bracket", "Plastic Storage Container",
];

function fmt(d: Date): string {
  const dd = String(d.getDate()).padStart(2, "0");
  const mm = String(d.getMonth() + 1).padStart(2, "0");
  return `${dd}/${mm}/${d.getFullYear()}`;
}

function makeItems(count: number): EwbItem[] {
  return Array.from({ length: count }, (_, i) => {
    const price = [45000, 8500, 15200, 3400, 1200][i % 5];
    const qty = [2, 5, 10, 25, 50][i % 5];
    return {
      slNo: i + 1,
      description: DESCRIPTIONS[i % 5],
      hsnCode: HSN_CODES[i % 5],
      quantity: qty,
      unit: "NOS",
      taxableValue: price * qty,
      gstRate: 18,
    };
  });
}

function validityDays(km: number): number {
  if (km <= 200) return 1;
  return Math.ceil(km / 200);
}

export function generateMockEwbRecords(): EwbRecord[] {
  const now = new Date("2026-04-27T10:00:00+05:30");
  const records: EwbRecord[] = [];

  for (let i = 0; i < 25; i++) {
    const supplier = SUPPLIERS[i % SUPPLIERS.length];
    const consignee = CONSIGNEES[i % CONSIGNEES.length];
    const genDate = new Date(now.getTime() - i * 3600000 * (1 + (i % 4)));
    const distanceKm = [150, 320, 480, 750, 1200, 60, 200][i % 7];
    const days = validityDays(distanceKm);
    const validUntil = new Date(genDate.getTime() + days * 86400000);
    const modes: TransportMode[] = ["Road", "Rail", "Air", "Ship"];
    const items = makeItems(1 + (i % 3));
    const taxableValue = items.reduce((s, it) => s + it.taxableValue, 0);
    const vehicle = VEHICLES[i % VEHICLES.length];

    let status: EwbRecord["status"] = "ACTIVE";
    if (i === 4 || i === 15) status = "CANCELLED";
    else if (i === 7 || i === 18) status = "CONSOLIDATED";
    else if (validUntil.getTime() < now.getTime()) status = "EXPIRED";

    const vehicleHistory: VehicleEntry[] = [
      { vehicleNo: vehicle, updatedAt: genDate.toISOString(), updatedBy: "finance@acme.com" },
    ];
    if (i % 5 === 0 && status === "ACTIVE") {
      vehicleHistory.push({
        vehicleNo: VEHICLES[(i + 3) % VEHICLES.length],
        updatedAt: new Date(genDate.getTime() + 3600000 * 4).toISOString(),
        updatedBy: "logistics@acme.com",
        reason: "Vehicle breakdown — replaced",
        fromPlace: "Pune",
      });
    }

    records.push({
      id: `ewb-${String(i + 1).padStart(4, "0")}`,
      ewbNumber: `${3310}${String(1000000 + i * 1117).padStart(8, "0")}`,
      invoiceNo: `INV-2026-${String(3000 + i).padStart(5, "0")}`,
      invoiceDate: fmt(new Date(genDate.getTime() - 86400000)),
      gstin: supplier.gstin,
      consigneeName: consignee.name,
      consigneeGstin: consignee.gstin,
      status,
      transportMode: modes[i % 4],
      vehicleNo: vehicleHistory[vehicleHistory.length - 1].vehicleNo,
      distanceKm,
      validFrom: genDate.toISOString(),
      validUntil: validUntil.toISOString(),
      generatedAt: genDate.toISOString(),
      cancelledAt: status === "CANCELLED"
        ? new Date(genDate.getTime() + 7200000).toISOString()
        : undefined,
      cancelReason: status === "CANCELLED" ? "Wrong consignee details" : undefined,
      totalValue: Math.round(taxableValue * 1.18),
      taxableValue,
      hsnCode: HSN_CODES[i % 5],
      consolidatedEwbNo: status === "CONSOLIDATED" ? `CEWB-${9900 + i}` : undefined,
      vehicleHistory,
      items,
    });
  }

  return records;
}

export function generateMockSourceInvoicesForEwb(): SourceInvoiceForEwb[] {
  return Array.from({ length: 15 }, (_, i) => {
    const supplier = SUPPLIERS[i % SUPPLIERS.length];
    const consignee = CONSIGNEES[i % CONSIGNEES.length];
    const items = makeItems(1 + (i % 3));
    const taxableValue = items.reduce((s, it) => s + it.taxableValue, 0);
    return {
      id: `src-ewb-${i + 1}`,
      invoiceNo: `INV-2026-${String(4000 + i).padStart(5, "0")}`,
      invoiceDate: fmt(new Date(Date.now() - i * 86400000)),
      gstin: supplier.gstin,
      consigneeName: consignee.name,
      consigneeGstin: consignee.gstin,
      totalValue: Math.round(taxableValue * 1.18),
      taxableValue,
      hsnCode: HSN_CODES[i % 5],
      items,
    };
  });
}

export function getValidityDays(km: number): number {
  return validityDays(km);
}
