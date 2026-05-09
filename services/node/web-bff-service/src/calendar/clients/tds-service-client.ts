import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { FilingStatusResult, FilingQuery } from '../calendar.types';

function extractTaxYear(dueDate: string): string {
  const d = new Date(dueDate);
  const month = d.getMonth();
  const year = d.getFullYear();
  const fy = month >= 3 ? year : year - 1;
  return `${fy}-${String(fy + 1).slice(2)}`;
}

function mapFormType(eventType: string): string | null {
  const map: Record<string, string> = {
    tds_deposit: 'challan_281',
    form138: '138',
    form140: '140',
    form144: '144',
    form27eq: '27eq',
    form130: '130',
    form131: '131',
    property_tds_194ia: '26qb',
  };
  return map[eventType] ?? null;
}

@Injectable()
export class TdsServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(TdsServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('TDS_SERVICE_URL', 'http://localhost:8099');
  }

  async getFilingStatus(query: FilingQuery): Promise<FilingStatusResult> {
    const formType = mapFormType(query.eventType);
    if (!formType) return { status: null };

    const taxYear = extractTaxYear(query.dueDate);
    // TODO: actual tds-service list endpoint may use different query param names
    const params = new URLSearchParams({
      form_type: formType,
      tax_year: taxYear,
      period: query.dueDate,
    });
    const url = `${this.baseUrl}/api/v1/tds/entries?${params}`;

    try {
      const controller = new AbortController();
      const timeout = setTimeout(() => controller.abort(), 5000);
      const res = await fetch(url, {
        headers: {
          'X-Tenant-Id': query.tenantId,
          'Content-Type': 'application/json',
        },
        signal: controller.signal,
      });
      clearTimeout(timeout);

      if (!res.ok) return { status: null };

      const body = await res.json().catch(() => null);
      const entries = body?.data;
      if (!Array.isArray(entries) || entries.length === 0) return { status: null };

      const entry = entries[0];
      return {
        status: entry.status === 'filed' || entry.status === 'submitted' ? entry.status : null,
        filedAt: entry.filed_at ? new Date(entry.filed_at) : undefined,
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      this.logger.warn(`TDS filing status query failed for ${query.eventType}: ${message}`);
      return { status: null };
    }
  }
}
