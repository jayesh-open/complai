export type EventCategory = "direct_tax" | "indirect_tax" | "statutory";
export type EventStatus = "filed" | "filed_late" | "due_soon" | "upcoming" | "overdue";
export type Authority = "CBDT" | "CBIC" | "MCA" | "EPFO" | "ESIC" | "GSTN";

export interface ComplianceEvent {
  id: string;
  title: string;
  description: string;
  category: EventCategory;
  authority: Authority;
  sectionRef?: string;
  formRef?: string;
  dueDateOffset?: number;
  penalty?: string;
  status: EventStatus;
  linkedModule?: string;
  eventType?: string;
}

export type CalendarEvent = ComplianceEvent & { dueDate: Date };
