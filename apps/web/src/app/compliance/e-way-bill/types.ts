export type EwbStatus = "ACTIVE" | "EXPIRED" | "CANCELLED" | "CONSOLIDATED";

export type TransportMode = "Road" | "Rail" | "Air" | "Ship";

export interface VehicleEntry {
  vehicleNo: string;
  updatedAt: string;
  updatedBy: string;
  reason?: string;
  fromPlace?: string;
}

export interface EwbRecord {
  id: string;
  ewbNumber: string;
  invoiceNo: string;
  invoiceDate: string;
  gstin: string;
  consigneeName: string;
  consigneeGstin: string;
  status: EwbStatus;
  transportMode: TransportMode;
  vehicleNo: string;
  distanceKm: number;
  validFrom: string;
  validUntil: string;
  generatedAt: string;
  cancelledAt?: string;
  cancelReason?: string;
  totalValue: number;
  taxableValue: number;
  hsnCode: string;
  consolidatedEwbNo?: string;
  vehicleHistory: VehicleEntry[];
  items: EwbItem[];
}

export interface EwbItem {
  slNo: number;
  description: string;
  hsnCode: string;
  quantity: number;
  unit: string;
  taxableValue: number;
  gstRate: number;
}

export interface SourceInvoiceForEwb {
  id: string;
  invoiceNo: string;
  invoiceDate: string;
  gstin: string;
  consigneeName: string;
  consigneeGstin: string;
  totalValue: number;
  taxableValue: number;
  hsnCode: string;
  items: EwbItem[];
}
