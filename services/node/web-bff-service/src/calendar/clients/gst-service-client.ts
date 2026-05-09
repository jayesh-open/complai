import { Injectable, Logger } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import type { FilingStatusResult, FilingQuery } from '../calendar.types';

const SUMMARY_ENDPOINTS: Record<string, string> = {
  gstr1_monthly: '/v1/gst/gstr1/summary',
  gstr1_iff_m1: '/v1/gst/gstr1/summary',
  gstr1_iff_m2: '/v1/gst/gstr1/summary',
  gstr1_qrmp: '/v1/gst/gstr1/summary',
  gstr3b_monthly: '/v1/gst/gstr3b/summary',
  gstr3b_qrmp_north: '/v1/gst/gstr3b/summary',
  gstr3b_qrmp_south: '/v1/gst/gstr3b/summary',
  pmt06: '/v1/gst/gstr3b/summary',
  // TODO: gstr7, gstr8, lut_rfd11 endpoints not yet available in gst-service
};

@Injectable()
export class GstServiceClient {
  private readonly baseUrl: string;
  private readonly logger = new Logger(GstServiceClient.name);

  constructor(config: ConfigService) {
    this.baseUrl = config.get('GST_SERVICE_URL', 'http://localhost:8093');
  }

  async getFilingStatus(query: FilingQuery): Promise<FilingStatusResult> {
    const endpoint = SUMMARY_ENDPOINTS[query.eventType];
    if (!endpoint) return { status: null };

    const params = new URLSearchParams({
      gstin: query.gstin ?? '',
      return_period: query.dueDate,
    });
    const url = `${this.baseUrl}${endpoint}?${params}`;

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
      if (!body?.data) return { status: null };

      const filing = body.data;
      return {
        status: filing.status === 'filed' || filing.status === 'submitted' ? filing.status : null,
        filedAt: filing.filed_at ? new Date(filing.filed_at) : undefined,
      };
    } catch (err: unknown) {
      const message = err instanceof Error ? err.message : String(err);
      this.logger.warn(`GST filing status query failed for ${query.eventType}: ${message}`);
      return { status: null };
    }
  }
}
