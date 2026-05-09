import { describe, it, expect, vi, beforeEach, afterEach } from 'vitest';
import { ConfigService } from '@nestjs/config';
import { GstServiceClient } from '../clients/gst-service-client';
import { TdsServiceClient } from '../clients/tds-service-client';
import { ItrServiceClient } from '../clients/itr-service-client';
import { Gstr9ServiceClient } from '../clients/gstr9-service-client';
import type { FilingQuery } from '../calendar.types';

const originalFetch = globalThis.fetch;

function mockConfig(overrides: Record<string, string> = {}): ConfigService {
  return { get: (key: string, def: string) => overrides[key] ?? def } as unknown as ConfigService;
}

function baseQuery(overrides: Partial<FilingQuery> = {}): FilingQuery {
  return {
    eventType: 'gstr3b_monthly',
    tenantId: 'tenant-1',
    dueDate: '2026-05-20',
    gstin: '29AABCU9603R1ZM',
    pan: 'AABCU9603R',
    ...overrides,
  };
}

afterEach(() => {
  globalThis.fetch = originalFetch;
  vi.restoreAllMocks();
});

describe('GstServiceClient', () => {
  it('makes GET to correct URL with gstin and return_period', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: { status: 'filed', filed_at: '2026-05-18T10:00:00Z' } }),
    });
    const client = new GstServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr3b_monthly' }));

    expect(globalThis.fetch).toHaveBeenCalledOnce();
    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toContain('/v1/gst/gstr3b/summary');
    expect(url).toContain('gstin=29AABCU9603R1ZM');
    expect(url).toContain('return_period=2026-05-20');
    expect(result.status).toBe('filed');
    expect(result.filedAt).toEqual(new Date('2026-05-18T10:00:00Z'));
  });

  it('returns null status on non-2xx response', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const client = new GstServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery());
    expect(result.status).toBeNull();
  });

  it('returns null status on network error', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
    const client = new GstServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery());
    expect(result.status).toBeNull();
  });

  it('returns null for unknown event types (no endpoint)', async () => {
    const client = new GstServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr7' }));
    expect(result.status).toBeNull();
  });

  it('uses base URL from config', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: { status: 'filed' } }),
    });
    const client = new GstServiceClient(mockConfig({ GST_SERVICE_URL: 'http://gst:9000' }));
    await client.getFilingStatus(baseQuery());
    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toMatch(/^http:\/\/gst:9000/);
  });

  it('returns null when response body has no data', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({}),
    });
    const client = new GstServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery());
    expect(result.status).toBeNull();
  });
});

describe('TdsServiceClient', () => {
  it('makes GET to /api/v1/tds/entries with correct params', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [{ status: 'filed', filed_at: '2026-05-06T10:00:00Z' }] }),
    });
    const client = new TdsServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'tds_deposit', dueDate: '2026-05-07' }));

    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toContain('/api/v1/tds/entries');
    expect(url).toContain('form_type=challan_281');
    expect(url).toContain('period=2026-05-07');
    expect(result.status).toBe('filed');
  });

  it('returns null for unknown TDS event type', async () => {
    const client = new TdsServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'unknown_tds' }));
    expect(result.status).toBeNull();
  });

  it('returns null on non-2xx response', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 404 });
    const client = new TdsServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'form138' }));
    expect(result.status).toBeNull();
  });

  it('returns null on network error', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('timeout'));
    const client = new TdsServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'form140' }));
    expect(result.status).toBeNull();
  });

  it('returns null when entries array is empty', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [] }),
    });
    const client = new TdsServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'form138' }));
    expect(result.status).toBeNull();
  });
});

describe('ItrServiceClient', () => {
  it('makes GET to /api/v1/itr/filings with correct params', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [{ status: 'submitted', filed_at: '2026-07-30T10:00:00Z' }] }),
    });
    const client = new ItrServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'itr_non_audit', dueDate: '2026-07-31' }));

    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toContain('/api/v1/itr/filings');
    expect(url).toContain('form_type=itr_1_2_4');
    expect(url).toContain('pan=AABCU9603R');
    expect(result.status).toBe('submitted');
  });

  it('returns null for unknown ITR event type', async () => {
    const client = new ItrServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'unknown_itr' }));
    expect(result.status).toBeNull();
  });

  it('returns null on network error', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
    const client = new ItrServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'itr_audit' }));
    expect(result.status).toBeNull();
  });

  it('returns null on non-2xx response', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 503 });
    const client = new ItrServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'tax_audit' }));
    expect(result.status).toBeNull();
  });
});

describe('Gstr9ServiceClient', () => {
  it('makes GET to /api/v1/gstr9/annual-return with correct params', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [{ status: 'filed', filed_at: '2026-12-20T10:00:00Z' }] }),
    });
    const client = new Gstr9ServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr9', dueDate: '2026-12-31' }));

    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toContain('/api/v1/gstr9/annual-return');
    expect(url).toContain('form_type=gstr9');
    expect(url).toContain('gstin=29AABCU9603R1ZM');
    expect(result.status).toBe('filed');
  });

  it('uses gstr9c form_type for 9C events', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [{ status: 'filed' }] }),
    });
    const client = new Gstr9ServiceClient(mockConfig());
    await client.getFilingStatus(baseQuery({ eventType: 'gstr9c', dueDate: '2026-12-31' }));

    const url = (globalThis.fetch as ReturnType<typeof vi.fn>).mock.calls[0]![0] as string;
    expect(url).toContain('form_type=gstr9c');
  });

  it('returns null on network error', async () => {
    globalThis.fetch = vi.fn().mockRejectedValue(new Error('ECONNREFUSED'));
    const client = new Gstr9ServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr9' }));
    expect(result.status).toBeNull();
  });

  it('returns null on non-2xx response', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({ ok: false, status: 500 });
    const client = new Gstr9ServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr9' }));
    expect(result.status).toBeNull();
  });

  it('returns null when annual returns array is empty', async () => {
    globalThis.fetch = vi.fn().mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({ data: [] }),
    });
    const client = new Gstr9ServiceClient(mockConfig());
    const result = await client.getFilingStatus(baseQuery({ eventType: 'gstr9' }));
    expect(result.status).toBeNull();
  });
});
