import type { FilingStatusResult } from '../calendar.types';

export class GstServiceClient {
  async getFilingStatus(
    eventType: string,
    identifier: string,
    period: string,
  ): Promise<FilingStatusResult> {
    // Mock: return null (no filing record) for all queries.
    // Real HTTP calls wired in Calendar-2b-2.
    return { status: null };
  }
}
