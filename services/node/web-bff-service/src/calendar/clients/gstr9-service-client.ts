import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { FilingStatusResult, FilingQuery } from '../calendar.types';

function extractFy(dueDate: string): string {
  const d = new Date(dueDate);
  const month = d.getMonth();
  const year = d.getFullYear();
  const fy = month >= 3 ? year : year - 1;
  return `${fy}-${String(fy + 1).slice(2)}`;
}

@Injectable()
export class Gstr9ServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(Gstr9ServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('GSTR9_SERVICE_URL', 'http://localhost:8102');
  }

  async getFilingStatus(query: FilingQuery): Promise<FilingStatusResult> {
    const fy = extractFy(query.dueDate);
    const formType = query.eventType === 'gstr9c' ? 'gstr9c' : 'gstr9';
    // TODO: actual gstr9-service list endpoint may use different query param names
    const params = new URLSearchParams({
      fy,
      form_type: formType,
    });
    if (query.gstin) params.set('gstin', query.gstin);
    const url = `${this.baseUrl}/api/v1/gstr9/annual-return?${params}`;

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
      const returns = body?.data;
      if (!Array.isArray(returns) || returns.length === 0) return { status: null };

      const annual = returns[0];
      return {
        status: annual.status === 'filed' || annual.status === 'submitted' ? annual.status : null,
        filedAt: annual.filed_at ? new Date(annual.filed_at) : undefined,
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      this.logger.warn(`GSTR-9 filing status query failed for ${query.eventType}: ${message}`);
      return { status: null };
    }
  }
}
