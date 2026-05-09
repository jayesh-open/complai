import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { CalendarController } from '../calendar.controller';
import { GstServiceClient } from '../clients/gst-service-client';
import { Gstr9ServiceClient } from '../clients/gstr9-service-client';
import { TdsServiceClient } from '../clients/tds-service-client';
import { ItrServiceClient } from '../clients/itr-service-client';
import { ConfigService } from '@nestjs/config';
import { BadRequestException, NotFoundException } from '@nestjs/common';

function mockConfig(): ConfigService {
  return { get: (_k: string, d: string) => d } as unknown as ConfigService;
}

function mockClient() {
  return { getFilingStatus: vi.fn().mockResolvedValue({ status: null }) };
}

function mockCache() {
  const store = new Map<string, unknown>();
  return {
    get: vi.fn(async (key: string) => store.get(key) ?? undefined),
    set: vi.fn(async (key: string, value: unknown) => { store.set(key, value); }),
    _store: store,
  };
}

function createController(cacheOverride?: ReturnType<typeof mockCache>) {
  const config = mockConfig();
  const cache = cacheOverride ?? mockCache();
  const ctrl = new CalendarController(
    new GstServiceClient(config),
    new Gstr9ServiceClient(config),
    new TdsServiceClient(config),
    new ItrServiceClient(config),
    cache as any,
  );
  return { ctrl, cache };
}

const originalFetch = globalThis.fetch;

afterEach(() => {
  globalThis.fetch = originalFetch;
  vi.restoreAllMocks();
});

beforeEach(() => {
  globalThis.fetch = vi.fn().mockResolvedValue({
    ok: false,
    status: 503,
  });
});

describe('CalendarController — getEvents', () => {
  it('returns 200 with events array for valid date range', async () => {
    const { ctrl } = createController();
    const result = await ctrl.getEvents('2026-05-01', '2026-05-31', 'acme');

    expect(result.events).toBeInstanceOf(Array);
    expect(result.events.length).toBeGreaterThan(0);
    expect(result.generated_at).toBeDefined();
    expect(typeof result.generated_at).toBe('string');
  });

  it('throws 400 when "from" is missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvents('', '2026-05-31')).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when "to" is missing', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvents('2026-05-01', '')).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when "from" is not ISO date format', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvents('05/01/2026', '2026-05-31')).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when "to" < "from"', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvents('2026-06-01', '2026-05-01')).rejects.toThrow(BadRequestException);
  });

  it('throws 400 when range exceeds 366 days', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvents('2026-01-01', '2027-06-01')).rejects.toThrow(BadRequestException);
  });

  it('all events have status field set', async () => {
    const { ctrl } = createController();
    const result = await ctrl.getEvents('2026-05-01', '2026-05-31', 'acme');
    for (const event of result.events) {
      expect(['filed', 'filed_late', 'overdue', 'due_soon', 'upcoming']).toContain(event.status);
    }
  });

  it('returns cached result on cache hit', async () => {
    const cache = mockCache();
    const cachedResult = { events: [{ id: 'cached' }], generated_at: '2026-05-01T00:00:00Z' };
    cache._store.set('calendar:events:acme:2026-05-01:2026-05-31', cachedResult);
    const { ctrl } = createController(cache);

    const result = await ctrl.getEvents('2026-05-01', '2026-05-31', 'acme');
    expect(result).toBe(cachedResult);
  });

  it('populates cache on miss', async () => {
    const cache = mockCache();
    const { ctrl } = createController(cache);
    await ctrl.getEvents('2026-05-01', '2026-05-31', 'acme');

    expect(cache.set).toHaveBeenCalledOnce();
    const [key] = cache.set.mock.calls[0]!;
    expect(key).toBe('calendar:events:acme:2026-05-01:2026-05-31');
  });

  it('defaults tenant_id to "default" when not provided', async () => {
    const cache = mockCache();
    const { ctrl } = createController(cache);
    await ctrl.getEvents('2026-05-01', '2026-05-31');

    const [key] = cache.set.mock.calls[0]!;
    expect(key).toContain(':default:');
  });
});

describe('CalendarController — getEvent', () => {
  it('returns single event for valid event ID', async () => {
    const { ctrl } = createController();
    const listResult = await ctrl.getEvents('2026-05-01', '2026-05-31', 'acme');
    const firstEvent = listResult.events[0]!;

    const single = await ctrl.getEvent(firstEvent.id, 'acme', '2026-05-01', '2026-05-31');
    expect(single!.id).toBe(firstEvent.id);
    expect(single!.eventType).toBe(firstEvent.eventType);
  });

  it('throws 404 for non-existent event ID', async () => {
    const { ctrl } = createController();
    await expect(ctrl.getEvent('evt-nonexistent', 'acme', '2026-05-01', '2026-05-31')).rejects.toThrow(NotFoundException);
  });
});
