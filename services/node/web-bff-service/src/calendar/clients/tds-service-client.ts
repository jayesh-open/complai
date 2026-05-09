import type { FilingStatusResult } from '../calendar.types';

export class TdsServiceClient {
  async getFilingStatus(
    eventType: string,
    identifier: string,
    period: string,
  ): Promise<FilingStatusResult> {
    return { status: null };
  }
}
