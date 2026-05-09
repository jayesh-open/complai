import { describe, it, expect } from 'vitest';
import { generateDueDateSeries } from '../due-date-generator';
import type { TenantConfig } from '../calendar.types';

const BASE_CONFIG: TenantConfig = {
  tenantId: 'test-tenant',
  gstins: ['29AABCU9603R1ZM'],
  pans: ['AABCU9603R'],
  filingScheme: 'monthly',
  businessType: 'company',
  annualTurnover: 10_00_00_000,
  requiresAudit: true,
};

describe('generateDueDateSeries', () => {
  it('returns events within the requested date range', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(events.length).toBeGreaterThan(0);
    for (const e of events) {
      expect(e.dueDate >= '2026-05-01').toBe(true);
      expect(e.dueDate <= '2026-05-31').toBe(true);
    }
  });

  it('generates deterministic event IDs for same tenant+type+date', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const run1 = generateDueDateSeries(BASE_CONFIG, from, to);
    const run2 = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(run1.map((e) => e.id)).toEqual(run2.map((e) => e.id));
  });

  it('generates different IDs for different tenants', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const config2 = { ...BASE_CONFIG, tenantId: 'other-tenant' };
    const run1 = generateDueDateSeries(BASE_CONFIG, from, to);
    const run2 = generateDueDateSeries(config2, from, to);
    expect(run1[0]!.id).not.toBe(run2[0]!.id);
  });

  it('returns events sorted by dueDate', () => {
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    for (let i = 1; i < events.length; i++) {
      expect(events[i]!.dueDate >= events[i - 1]!.dueDate).toBe(true);
    }
  });

  it('includes monthly GSTR-1 for monthly filers', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const gstr1 = events.filter((e) => e.eventType === 'gstr1_monthly');
    expect(gstr1.length).toBe(1);
    expect(gstr1[0]!.dueDate).toBe('2026-05-11');
  });

  it('includes monthly GSTR-3B for monthly filers', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const gstr3b = events.filter((e) => e.eventType === 'gstr3b_monthly');
    expect(gstr3b.length).toBe(1);
    expect(gstr3b[0]!.dueDate).toBe('2026-05-20');
  });

  it('includes TDS deposit on 7th', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const tds = events.filter((e) => e.eventType === 'tds_deposit');
    expect(tds.length).toBe(1);
    expect(tds[0]!.dueDate).toBe('2026-05-07');
  });

  it('includes PF and ESI on 15th', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const pf = events.filter((e) => e.eventType === 'pf_monthly');
    const esi = events.filter((e) => e.eventType === 'esi_monthly');
    expect(pf.length).toBe(1);
    expect(esi.length).toBe(1);
    expect(pf[0]!.dueDate).toBe('2026-05-15');
  });

  it('all events have status=upcoming by default', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    for (const e of events) {
      expect(e.status).toBe('upcoming');
    }
  });

  it('includes correct category and authority on events', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const tds = events.find((e) => e.eventType === 'tds_deposit')!;
    expect(tds.category).toBe('direct_tax');
    expect(tds.authority).toBe('CBDT');
    expect(tds.sectionRef).toBe('ITA 2025 § 392');

    const gstr1 = events.find((e) => e.eventType === 'gstr1_monthly')!;
    expect(gstr1.category).toBe('indirect_tax');
    expect(gstr1.authority).toBe('CBIC');
  });

  it('spans multiple months correctly', () => {
    const from = new Date(2026, 3, 1);
    const to = new Date(2026, 5, 30);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    const gstr1 = events.filter((e) => e.eventType === 'gstr1_monthly');
    expect(gstr1.length).toBe(3);
    expect(gstr1.map((e) => e.dueDate)).toEqual(['2026-04-11', '2026-05-11', '2026-06-11']);
  });
});

describe('tenant-aware filtering', () => {
  it('excludes monthly GSTR-1 for QRMP filers', () => {
    const qrmpConfig: TenantConfig = { ...BASE_CONFIG, filingScheme: 'qrmp' };
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(qrmpConfig, from, to);
    expect(events.filter((e) => e.eventType === 'gstr1_monthly').length).toBe(0);
  });

  it('excludes monthly GSTR-3B for QRMP filers', () => {
    const qrmpConfig: TenantConfig = { ...BASE_CONFIG, filingScheme: 'qrmp' };
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const events = generateDueDateSeries(qrmpConfig, from, to);
    expect(events.filter((e) => e.eventType === 'gstr3b_monthly').length).toBe(0);
  });

  it('includes PMT-06 for QRMP filers only', () => {
    const from = new Date(2026, 4, 1);
    const to = new Date(2026, 4, 31);
    const monthlyEvents = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(monthlyEvents.filter((e) => e.eventType === 'pmt06').length).toBe(0);

    const qrmpConfig: TenantConfig = { ...BASE_CONFIG, filingScheme: 'qrmp' };
    const qrmpEvents = generateDueDateSeries(qrmpConfig, from, to);
    expect(qrmpEvents.filter((e) => e.eventType === 'pmt06').length).toBe(1);
  });

  it('excludes GSTR-9 when turnover < ₹2 Cr', () => {
    const lowTurnover: TenantConfig = { ...BASE_CONFIG, annualTurnover: 1_50_00_000 };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(lowTurnover, from, to);
    expect(events.filter((e) => e.eventType === 'gstr9').length).toBe(0);
  });

  it('includes GSTR-9 when turnover ≥ ₹2 Cr', () => {
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(events.filter((e) => e.eventType === 'gstr9').length).toBe(1);
  });

  it('excludes GSTR-9C when turnover < ₹5 Cr', () => {
    const medTurnover: TenantConfig = { ...BASE_CONFIG, annualTurnover: 3_00_00_000 };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(medTurnover, from, to);
    expect(events.filter((e) => e.eventType === 'gstr9c').length).toBe(0);
  });

  it('includes GSTR-9C when turnover ≥ ₹5 Cr', () => {
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(events.filter((e) => e.eventType === 'gstr9c').length).toBe(1);
  });

  it('excludes tax audit for non-audit cases', () => {
    const noAudit: TenantConfig = { ...BASE_CONFIG, requiresAudit: false };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(noAudit, from, to);
    expect(events.filter((e) => e.eventType === 'tax_audit').length).toBe(0);
    expect(events.filter((e) => e.eventType === 'itr_audit').length).toBe(0);
  });

  it('includes non-audit ITR for non-audit cases', () => {
    const noAudit: TenantConfig = { ...BASE_CONFIG, requiresAudit: false };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(noAudit, from, to);
    expect(events.filter((e) => e.eventType === 'itr_non_audit').length).toBe(1);
  });

  it('excludes LLP Form 11 for companies', () => {
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(events.filter((e) => e.eventType === 'llp_form11').length).toBe(0);
  });

  it('includes LLP Form 11 for LLPs', () => {
    const llpConfig: TenantConfig = { ...BASE_CONFIG, businessType: 'llp' };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(llpConfig, from, to);
    expect(events.filter((e) => e.eventType === 'llp_form11').length).toBe(1);
  });

  it('excludes company-specific MCA filings for LLPs', () => {
    const llpConfig: TenantConfig = { ...BASE_CONFIG, businessType: 'llp' };
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(llpConfig, from, to);
    expect(events.filter((e) => e.eventType === 'aoc4').length).toBe(0);
    expect(events.filter((e) => e.eventType === 'mgt7').length).toBe(0);
  });

  it('includes company-specific MCA filings for companies', () => {
    const from = new Date(2026, 0, 1);
    const to = new Date(2026, 11, 31);
    const events = generateDueDateSeries(BASE_CONFIG, from, to);
    expect(events.filter((e) => e.eventType === 'aoc4').length).toBe(1);
    expect(events.filter((e) => e.eventType === 'mgt7').length).toBe(1);
  });
});
