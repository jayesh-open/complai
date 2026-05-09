import {
  Controller,
  Get,
  Query,
  Param,
  Inject,
  BadRequestException,
  NotFoundException,
} from '@nestjs/common';
import { CACHE_MANAGER } from '@nestjs/cache-manager';
import type { Cache } from 'cache-manager';
import { generateDueDateSeries } from './due-date-generator';
import { aggregateStatuses, type StatusAggregatorDeps } from './status-aggregator';
import { GstServiceClient } from './clients/gst-service-client';
import { Gstr9ServiceClient } from './clients/gstr9-service-client';
import { TdsServiceClient } from './clients/tds-service-client';
import { ItrServiceClient } from './clients/itr-service-client';
import type { ComplianceEvent, TenantConfig } from './calendar.types';
import { DEFAULT_TENANT_CONFIG } from './calendar.types';

const ISO_DATE = /^\d{4}-\d{2}-\d{2}$/;
const MAX_RANGE_DAYS = 366;
const CACHE_TTL_MS = 5 * 60 * 1000;

@Controller('api/v1/compliance/calendar')
export class CalendarController {
  private readonly deps: StatusAggregatorDeps;

  constructor(
    private readonly gstClient: GstServiceClient,
    private readonly gstr9Client: Gstr9ServiceClient,
    private readonly tdsClient: TdsServiceClient,
    private readonly itrClient: ItrServiceClient,
    @Inject(CACHE_MANAGER) private readonly cache: Cache,
  ) {
    this.deps = {
      gstClient: this.gstClient,
      gstr9Client: this.gstr9Client,
      tdsClient: this.tdsClient,
      itrClient: this.itrClient,
    };
  }

  @Get('events')
  async getEvents(
    @Query('from') from: string,
    @Query('to') to: string,
    @Query('tenant_id') tenantId?: string,
  ) {
    this.validateDateParams(from, to);
    const tid = tenantId ?? 'default';

    const cacheKey = `calendar:events:${tid}:${from}:${to}`;
    const cached = await this.cache.get<{ events: ComplianceEvent[]; generated_at: string }>(cacheKey);
    if (cached) return cached;

    const config: TenantConfig = { ...DEFAULT_TENANT_CONFIG, tenantId: tid };
    const fromDate = new Date(from);
    const toDate = new Date(to);

    const rawEvents = generateDueDateSeries(config, fromDate, toDate);
    const events = await aggregateStatuses(
      rawEvents,
      this.deps,
      { tenantId: tid, gstin: config.gstins[0], pan: config.pans[0] },
    );

    const result = { events, generated_at: new Date().toISOString() };
    await this.cache.set(cacheKey, result, CACHE_TTL_MS);
    return result;
  }

  @Get('event/:eventId')
  async getEvent(
    @Param('eventId') eventId: string,
    @Query('tenant_id') tenantId?: string,
    @Query('from') from?: string,
    @Query('to') to?: string,
  ) {
    const tid = tenantId ?? 'default';
    const config: TenantConfig = { ...DEFAULT_TENANT_CONFIG, tenantId: tid };

    const fromDate = from ? new Date(from) : yearAgo();
    const toDate = to ? new Date(to) : yearAhead();

    const rawEvents = generateDueDateSeries(config, fromDate, toDate);
    const match = rawEvents.find((e) => e.id === eventId);
    if (!match) throw new NotFoundException(`Event ${eventId} not found`);

    const [enriched] = await aggregateStatuses(
      [match],
      this.deps,
      { tenantId: tid, gstin: config.gstins[0], pan: config.pans[0] },
    );
    return enriched;
  }

  private validateDateParams(from: string, to: string): void {
    if (!from || !ISO_DATE.test(from)) {
      throw new BadRequestException('Missing or invalid "from" query param (expected YYYY-MM-DD)');
    }
    if (!to || !ISO_DATE.test(to)) {
      throw new BadRequestException('Missing or invalid "to" query param (expected YYYY-MM-DD)');
    }

    const fromDate = new Date(from);
    const toDate = new Date(to);

    if (isNaN(fromDate.getTime()) || isNaN(toDate.getTime())) {
      throw new BadRequestException('Invalid date value');
    }
    if (toDate < fromDate) {
      throw new BadRequestException('"to" must be >= "from"');
    }

    const diffMs = toDate.getTime() - fromDate.getTime();
    const diffDays = diffMs / (1000 * 60 * 60 * 24);
    if (diffDays > MAX_RANGE_DAYS) {
      throw new BadRequestException(`Date range exceeds ${MAX_RANGE_DAYS} days`);
    }
  }
}

function yearAgo(): Date {
  const d = new Date();
  d.setFullYear(d.getFullYear() - 1);
  return d;
}

function yearAhead(): Date {
  const d = new Date();
  d.setFullYear(d.getFullYear() + 1);
  return d;
}
