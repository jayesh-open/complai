"use client";

import { useState, useEffect, useRef } from "react";
import type { CalendarEvent } from "../types";
import { getComplianceEvents } from "../mock-data";

interface UseCalendarEventsResult {
  events: CalendarEvent[];
  loading: boolean;
  error: string | null;
  usingMockData: boolean;
}

function formatDate(d: Date): string {
  const y = d.getFullYear();
  const m = String(d.getMonth() + 1).padStart(2, "0");
  const day = String(d.getDate()).padStart(2, "0");
  return `${y}-${m}-${day}`;
}

function getGridDateRange(year: number, month: number): { from: string; to: string } {
  const firstDay = new Date(year, month, 1);
  let startOffset = firstDay.getDay() - 1;
  if (startOffset < 0) startOffset = 6;
  const gridStart = new Date(year, month, 1 - startOffset);
  const gridEnd = new Date(gridStart);
  gridEnd.setDate(gridStart.getDate() + 41);
  return { from: formatDate(gridStart), to: formatDate(gridEnd) };
}

export function useCalendarEvents(
  year: number,
  month: number,
  tenantId = "default",
): UseCalendarEventsResult {
  const [events, setEvents] = useState<CalendarEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [usingMockData, setUsingMockData] = useState(false);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (retryRef.current) clearTimeout(retryRef.current);

    const controller = new AbortController();

    const doFetch = async () => {
      const { from, to } = getGridDateRange(year, month);
      try {
        const res = await fetch(
          `/api/v1/compliance/calendar/events?from=${from}&to=${to}&tenant_id=${encodeURIComponent(tenantId)}`,
          { signal: controller.signal },
        );
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        const parsed: CalendarEvent[] = (data.events ?? []).map((e: any) => ({
          ...e,
          dueDate: new Date(e.dueDate + "T00:00:00"),
        }));
        setEvents(parsed);
        setError(null);
        setUsingMockData(false);
      } catch (err) {
        if (controller.signal.aborted) return;
        setEvents(getComplianceEvents(new Date()));
        setError(err instanceof Error ? err.message : "Failed to fetch events");
        setUsingMockData(true);
        retryRef.current = setTimeout(doFetch, 30_000);
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    };

    setLoading(true);
    doFetch();

    return () => {
      controller.abort();
      if (retryRef.current) clearTimeout(retryRef.current);
    };
  }, [year, month, tenantId]);

  return { events, loading, error, usingMockData };
}
