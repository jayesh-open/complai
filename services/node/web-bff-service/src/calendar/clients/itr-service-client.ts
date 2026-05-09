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
    itr_non_audit: 'itr_1_2_4',
    itr_audit: 'itr_3_5_6',
    itr7: 'itr_7',
    tax_audit: 'tax_audit',
    advance_tax_q1: 'advance_tax',
    advance_tax_q2: 'advance_tax',
    advance_tax_q3: 'advance_tax',
    advance_tax_q4: 'advance_tax',
  };
  return map[eventType] ?? null;
}

@Injectable()
export class ItrServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(ItrServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('ITR_SERVICE_URL', 'http://localhost:8100');
  }

  async getFilingStatus(query: FilingQuery): Promise<FilingStatusResult> {
    const formType = mapFormType(query.eventType);
    if (!formType) return { status: null };

    const taxYear = extractTaxYear(query.dueDate);
    // TODO: actual itr-service /api/v1/itr/filings may use different filter params
    const params = new URLSearchParams({
      form_type: formType,
      tax_year: taxYear,
    });
    if (query.pan) params.set('pan', query.pan);
    const url = `${this.baseUrl}/api/v1/itr/filings?${params}`;

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
      const filings = body?.data;
      if (!Array.isArray(filings) || filings.length === 0) return { status: null };

      const filing = filings[0];
      return {
        status: filing.status === 'filed' || filing.status === 'submitted' ? filing.status : null,
        filedAt: filing.filed_at ? new Date(filing.filed_at) : undefined,
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      this.logger.warn(`ITR filing status query failed for ${query.eventType}: ${message}`);
      return { status: null };
    }
  }
}
