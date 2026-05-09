import type { FilingStatusResult } from '../calendar.types';

export class ItrServiceClient {
  async getFilingStatus(
    eventType: string,
    identifier: string,
    period: string,
  ): Promise<FilingStatusResult> {
    return { status: null };
  }
}
