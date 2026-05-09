import { createHash } from 'crypto';
import type {
  ComplianceEvent,
  EventTemplate,
  EventSchedule,
  TenantConfig,
  DEFAULT_TENANT_CONFIG,
} from './calendar.types';
import { EVENT_TEMPLATES } from './event-templates';

function formatDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, '0');
  const day = String(d.getDate()).padStart(2, '0');
  return `${y}-${m}-${day}`;
}

function monthName(d: Date): string {
  return d.toLocaleString('en-IN', { month: 'long' });
}

function deterministicId(tenantId: string, eventType: string, dueDate: string): string {
  const hash = createHash('sha256')
    .update(`${tenantId}:${eventType}:${dueDate}`)
    .digest('hex')
    .slice(0, 12);
  return `evt-${hash}`;
}

function expandSchedule(
  schedule: EventSchedule,
  from: Date,
  to: Date,
): Date[] {
  const dates: Date[] = [];
  const fromTime = from.getTime();
  const toTime = to.getTime();

  if (schedule.type === 'monthly') {
    const start = new Date(from.getFullYear(), from.getMonth(), 1);
    for (let i = -1; i <= 14; i++) {
      const d = new Date(start.getFullYear(), start.getMonth() + i, schedule.dayOfMonth);
      if (d.getTime() >= fromTime && d.getTime() <= toTime) {
        dates.push(d);
      }
    }
  } else if (schedule.type === 'quarterly') {
    for (let y = from.getFullYear() - 1; y <= to.getFullYear() + 1; y++) {
      for (const m of schedule.months) {
        const d = new Date(y, m, schedule.dayOfMonth);
        if (d.getTime() >= fromTime && d.getTime() <= toTime) {
          dates.push(d);
        }
      }
    }
  } else if (schedule.type === 'annual') {
    for (let y = from.getFullYear() - 1; y <= to.getFullYear() + 1; y++) {
      const d = new Date(y, schedule.month, schedule.dayOfMonth);
      if (d.getTime() >= fromTime && d.getTime() <= toTime) {
        dates.push(d);
      }
    }
  } else if (schedule.type === 'half_yearly') {
    for (let y = from.getFullYear() - 1; y <= to.getFullYear() + 1; y++) {
      for (const m of schedule.months) {
        const d = new Date(y, m, schedule.dayOfMonth);
        if (d.getTime() >= fromTime && d.getTime() <= toTime) {
          dates.push(d);
        }
      }
    }
  }

  return dates;
}

export function generateDueDateSeries(
  config: TenantConfig,
  from: Date,
  to: Date,
): ComplianceEvent[] {
  const events: ComplianceEvent[] = [];

  for (const tmpl of EVENT_TEMPLATES) {
    if (tmpl.filter && !tmpl.filter(config)) continue;

    const dueDates = expandSchedule(tmpl.schedule, from, to);
    for (const dueDate of dueDates) {
      const dueDateStr = formatDate(dueDate);
      const suffix = tmpl.titleSuffix ? tmpl.titleSuffix(dueDate) : `(${monthName(dueDate)} ${dueDate.getFullYear()})`;

      events.push({
        id: deterministicId(config.tenantId, tmpl.eventType, dueDateStr),
        eventType: tmpl.eventType,
        title: `${tmpl.title} ${suffix}`,
        description: tmpl.description,
        category: tmpl.category,
        authority: tmpl.authority,
        sectionRef: tmpl.sectionRef,
        formRef: tmpl.formRef,
        dueDate: dueDateStr,
        penalty: tmpl.penalty,
        status: 'upcoming',
        linkedModule: tmpl.linkedModule,
      });
    }
  }

  events.sort((a, b) => a.dueDate.localeCompare(b.dueDate));
  return events;
}
