import { describe, it, expect, vi } from 'vitest';
import { aggregateStatuses, type StatusAggregatorDeps } from '../status-aggregator';
import type { ComplianceEvent, FilingStatusResult, FilingStatusProvider } from '../calendar.types';

function makeEvent(overrides: Partial<ComplianceEvent> = {}): ComplianceEvent {
  return {
    id: 'evt-test-001',
    title: 'Test Event',
    description: 'A test event',
    category: 'indirect_tax',
    authority: 'CBIC',
    dueDate: '2026-05-20',
    status: 'upcoming',
    eventType: 'gstr3b_monthly',
    ...overrides,
  };
}

function makeMockClient(result: FilingStatusResult = { status: null }): FilingStatusProvider & { getFilingStatus: ReturnType<typeof vi.fn> } {
  return {
    getFilingStatus: vi.fn().mockResolvedValue(result),
  };
}

function makeDeps(overrides: Partial<StatusAggregatorDeps> = {}): StatusAggregatorDeps {
  return {
    gstClient: makeMockClient(),
    gstr9Client: makeMockClient(),
    tdsClient: makeMockClient(),
    itrClient: makeMockClient(),
    ...overrides,
  };
}

describe('aggregateStatuses', () => {
  it('returns "filed" when record exists and filed before due date', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({
        status: 'filed',
        filedAt: new Date('2026-05-18'),
      }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result[0]!.status).toBe('filed');
  });

  it('returns "filed" when record status is "submitted" and filed before due', async () => {
    const event = makeEvent({
      eventType: 'gstr1_monthly',
      dueDate: '2026-05-11',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({
        status: 'submitted',
        filedAt: new Date('2026-05-10'),
      }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result[0]!.status).toBe('filed');
  });

  it('returns "filed_late" when record exists but filed after due date', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({
        status: 'filed',
        filedAt: new Date('2026-05-25'),
      }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-26'));
    expect(result[0]!.status).toBe('filed_late');
  });

  it('returns "overdue" when no filing record and past due date', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({ status: null }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-25'));
    expect(result[0]!.status).toBe('overdue');
  });

  it('returns "due_soon" when no filing record and due within 7 days', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({ status: null }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result[0]!.status).toBe('due_soon');
  });

  it('returns "upcoming" when no filing record and due > 7 days away', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({ status: null }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-01'));
    expect(result[0]!.status).toBe('upcoming');
  });

  it('falls back to "upcoming" when service client throws', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const failingClient = {
      getFilingStatus: vi.fn().mockRejectedValue(new Error('Connection refused')),
    };
    const deps = makeDeps({ gstClient: failingClient });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-25'));
    expect(result[0]!.status).toBe('upcoming');
  });

  it('routes GST event types to gstClient', async () => {
    const gstTypes = ['gstr1_monthly', 'gstr3b_monthly', 'gstr7', 'gstr8', 'pmt06', 'lut_rfd11'];
    for (const eventType of gstTypes) {
      const event = makeEvent({ eventType });
      const client = makeMockClient({ status: 'filed', filedAt: new Date('2026-05-01') });
      const deps = makeDeps({ gstClient: client });

      await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
      expect(client.getFilingStatus).toHaveBeenCalledOnce();
    }
  });

  it('routes GSTR-9 event types to gstr9Client', async () => {
    for (const eventType of ['gstr9', 'gstr9c']) {
      const event = makeEvent({ eventType, category: 'indirect_tax', authority: 'CBIC' });
      const client = makeMockClient({ status: 'filed', filedAt: new Date('2026-05-01') });
      const deps = makeDeps({ gstr9Client: client });

      await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
      expect(client.getFilingStatus).toHaveBeenCalledOnce();
    }
  });

  it('routes TDS event types to tdsClient', async () => {
    const tdsTypes = ['tds_deposit', 'form138', 'form140', 'form144', 'form27eq', 'form130', 'form131'];
    for (const eventType of tdsTypes) {
      const event = makeEvent({ eventType, category: 'direct_tax', authority: 'CBDT' });
      const client = makeMockClient({ status: 'filed', filedAt: new Date('2026-05-01') });
      const deps = makeDeps({ tdsClient: client });

      await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
      expect(client.getFilingStatus).toHaveBeenCalledOnce();
    }
  });

  it('routes ITR event types to itrClient', async () => {
    const itrTypes = ['itr_non_audit', 'itr_audit', 'itr7', 'tax_audit', 'advance_tax_q1'];
    for (const eventType of itrTypes) {
      const event = makeEvent({ eventType, category: 'direct_tax', authority: 'CBDT' });
      const client = makeMockClient({ status: 'filed', filedAt: new Date('2026-05-01') });
      const deps = makeDeps({ itrClient: client });

      await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
      expect(client.getFilingStatus).toHaveBeenCalledOnce();
    }
  });

  it('returns "upcoming" for unrecognized event types (no service)', async () => {
    const event = makeEvent({ eventType: 'pf_monthly' });
    const deps = makeDeps();

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-04-01'));
    expect(result[0]!.status).toBe('upcoming');
    expect(deps.gstClient.getFilingStatus).not.toHaveBeenCalled();
    expect(deps.tdsClient.getFilingStatus).not.toHaveBeenCalled();
    expect(deps.itrClient.getFilingStatus).not.toHaveBeenCalled();
    expect(deps.gstr9Client.getFilingStatus).not.toHaveBeenCalled();
  });

  it('processes multiple events in parallel', async () => {
    const events = [
      makeEvent({ eventType: 'gstr1_monthly', dueDate: '2026-05-11', id: 'evt-1' }),
      makeEvent({ eventType: 'tds_deposit', dueDate: '2026-05-07', id: 'evt-2', category: 'direct_tax', authority: 'CBDT' }),
      makeEvent({ eventType: 'itr_audit', dueDate: '2026-10-31', id: 'evt-3', category: 'direct_tax', authority: 'CBDT' }),
    ];
    const deps = makeDeps({
      gstClient: makeMockClient({ status: 'filed', filedAt: new Date('2026-05-10') }),
      tdsClient: makeMockClient({ status: null }),
      itrClient: makeMockClient({ status: null }),
    });

    const result = await aggregateStatuses(events, deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result).toHaveLength(3);
    expect(result[0]!.status).toBe('filed');
    expect(result[1]!.status).toBe('overdue');
    expect(result[2]!.status).toBe('upcoming');
  });

  it('returns "filed" when filedAt is exactly on due date', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({
        status: 'filed',
        filedAt: new Date('2026-05-20T15:00:00'),
      }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-21'));
    expect(result[0]!.status).toBe('filed');
  });

  it('returns "due_soon" when due date is exactly 7 days away', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      dueDate: '2026-05-20',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({ status: null }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-13'));
    expect(result[0]!.status).toBe('due_soon');
  });

  it('handles empty events array', async () => {
    const deps = makeDeps();
    const result = await aggregateStatuses([], deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result).toEqual([]);
  });

  it('preserves all original event fields in output', async () => {
    const event = makeEvent({
      eventType: 'gstr3b_monthly',
      title: 'GSTR-3B (May 2026)',
      description: 'Monthly summary return',
      sectionRef: 'GST § 39',
      formRef: 'GSTR-3B',
      penalty: '₹50/day late fee',
      linkedModule: '/compliance/gst',
    });
    const deps = makeDeps({
      gstClient: makeMockClient({ status: 'filed', filedAt: new Date('2026-05-18') }),
    });

    const result = await aggregateStatuses([event], deps, { tenantId: 'test' }, new Date('2026-05-15'));
    expect(result[0]!.title).toBe('GSTR-3B (May 2026)');
    expect(result[0]!.description).toBe('Monthly summary return');
    expect(result[0]!.sectionRef).toBe('GST § 39');
    expect(result[0]!.formRef).toBe('GSTR-3B');
    expect(result[0]!.linkedModule).toBe('/compliance/gst');
  });
});
