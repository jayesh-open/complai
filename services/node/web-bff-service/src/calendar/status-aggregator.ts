import type { ComplianceEvent, FilingStatusResult, EventStatus } from './calendar.types';
import type { GstServiceClient } from './clients/gst-service-client';
import type { Gstr9ServiceClient } from './clients/gstr9-service-client';
import type { TdsServiceClient } from './clients/tds-service-client';
import type { ItrServiceClient } from './clients/itr-service-client';

const GST_EVENT_TYPES = new Set([
  'gstr1_monthly', 'gstr1_iff_m1', 'gstr1_iff_m2', 'gstr1_qrmp',
  'gstr3b_monthly', 'gstr3b_qrmp_north', 'gstr3b_qrmp_south',
  'pmt06', 'gstr7', 'gstr8', 'lut_rfd11',
]);

const GSTR9_EVENT_TYPES = new Set(['gstr9', 'gstr9c']);

const TDS_EVENT_TYPES = new Set([
  'tds_deposit', 'form138', 'form140', 'form144',
  'form27eq', 'form130', 'form131', 'property_tds_194ia',
]);

const ITR_EVENT_TYPES = new Set([
  'itr_non_audit', 'itr_audit', 'itr7', 'tax_audit',
  'advance_tax_q1', 'advance_tax_q2', 'advance_tax_q3', 'advance_tax_q4',
]);

function resolveStatus(
  filingResult: FilingStatusResult,
  dueDate: string,
  now: Date,
): EventStatus {
  const due = new Date(dueDate + 'T23:59:59');
  const sevenDaysFromNow = new Date(now);
  sevenDaysFromNow.setDate(sevenDaysFromNow.getDate() + 7);
  sevenDaysFromNow.setHours(23, 59, 59, 999);

  if (filingResult.status === 'filed' || filingResult.status === 'submitted') {
    if (filingResult.filedAt && filingResult.filedAt > due) {
      return 'filed_late';
    }
    return 'filed';
  }

  if (now > due) return 'overdue';
  if (due <= sevenDaysFromNow) return 'due_soon';
  return 'upcoming';
}

export interface StatusAggregatorDeps {
  gstClient: GstServiceClient;
  gstr9Client: Gstr9ServiceClient;
  tdsClient: TdsServiceClient;
  itrClient: ItrServiceClient;
}

export async function aggregateStatuses(
  events: ComplianceEvent[],
  deps: StatusAggregatorDeps,
  now: Date = new Date(),
): Promise<ComplianceEvent[]> {
  const enriched = await Promise.all(
    events.map(async (event) => {
      const client = getClientForEvent(event.eventType, deps);
      if (!client) {
        return { ...event, status: 'upcoming' as EventStatus };
      }

      try {
        const result = await client.getFilingStatus(
          event.eventType,
          event.id,
          event.dueDate,
        );
        return { ...event, status: resolveStatus(result, event.dueDate, now) };
      } catch {
        return { ...event, status: 'upcoming' as EventStatus };
      }
    }),
  );

  return enriched;
}

function getClientForEvent(
  eventType: string,
  deps: StatusAggregatorDeps,
): StatusAggregatorDeps[keyof StatusAggregatorDeps] | null {
  if (GST_EVENT_TYPES.has(eventType)) return deps.gstClient;
  if (GSTR9_EVENT_TYPES.has(eventType)) return deps.gstr9Client;
  if (TDS_EVENT_TYPES.has(eventType)) return deps.tdsClient;
  if (ITR_EVENT_TYPES.has(eventType)) return deps.itrClient;
  return null;
}
