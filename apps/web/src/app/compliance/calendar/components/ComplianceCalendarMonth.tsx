"use client";

import { useState, useMemo, useCallback } from "react";
import type { ComplianceEvent, EventCategory } from "../types";
import { getComplianceEvents, resolveEventsForMonth } from "../mock-data";
import { MonthNavigation } from "./MonthNavigation";
import { StatusLegend } from "./ComplianceStatusIcon";
import { ComplianceCategoryBadge } from "./ComplianceCategoryBadge";
import { ComplianceStatCard } from "./ComplianceStatCard";
import { ComplianceDayCell } from "./ComplianceDayCell";
import { ComplianceDayDetailPanel } from "./ComplianceDayDetailPanel";

const WEEKDAY_HEADERS = ["Mon", "Tue", "Wed", "Thu", "Fri", "Sat", "Sun"];

function getMonthGrid(year: number, month: number): Date[] {
  const firstDay = new Date(year, month, 1);
  let startOffset = firstDay.getDay() - 1;
  if (startOffset < 0) startOffset = 6;

  const cells: Date[] = [];
  const gridStart = new Date(year, month, 1 - startOffset);
  for (let i = 0; i < 42; i++) {
    const d = new Date(gridStart);
    d.setDate(gridStart.getDate() + i);
    cells.push(d);
  }
  return cells;
}

function isSameDay(a: Date, b: Date): boolean {
  return a.getFullYear() === b.getFullYear() && a.getMonth() === b.getMonth() && a.getDate() === b.getDate();
}

export function ComplianceCalendarMonth() {
  const today = useMemo(() => new Date(), []);
  const [year, setYear] = useState(today.getFullYear());
  const [month, setMonth] = useState(today.getMonth());
  const [categoryFilters, setCategoryFilters] = useState<Record<EventCategory, boolean>>({
    direct_tax: true,
    indirect_tax: true,
    statutory: true,
  });
  const [selectedDate, setSelectedDate] = useState<Date | null>(null);

  const allEvents = useMemo(() => getComplianceEvents(today), [today]);

  const filteredEvents = useMemo(() => {
    return allEvents.filter((e) => categoryFilters[e.category]);
  }, [allEvents, categoryFilters]);

  const monthEvents = useMemo(
    () => resolveEventsForMonth(filteredEvents, year, month),
    [filteredEvents, year, month],
  );

  const gridCells = useMemo(() => getMonthGrid(year, month), [year, month]);

  const toggleCategory = useCallback((cat: EventCategory) => {
    setCategoryFilters((prev) => ({ ...prev, [cat]: !prev[cat] }));
  }, []);

  const goPrev = useCallback(() => {
    setMonth((m) => {
      if (m === 0) { setYear((y) => y - 1); return 11; }
      return m - 1;
    });
  }, []);

  const goNext = useCallback(() => {
    setMonth((m) => {
      if (m === 11) { setYear((y) => y + 1); return 0; }
      return m + 1;
    });
  }, []);

  const goToday = useCallback(() => {
    setYear(today.getFullYear());
    setMonth(today.getMonth());
  }, [today]);

  const selectedEvents = useMemo(() => {
    if (!selectedDate) return [];
    return filteredEvents.filter((e) => isSameDay(e.dueDate, selectedDate));
  }, [selectedDate, filteredEvents]);

  const filedCount = monthEvents.filter((e) => e.status === "filed").length;
  const dueSoonCount = filteredEvents.filter((e) => {
    const diff = Math.ceil((e.dueDate.getTime() - today.getTime()) / 86400000);
    return diff >= 0 && diff <= 7;
  }).length;
  const upcomingCount = monthEvents.filter((e) => e.status === "upcoming" || e.status === "due_soon").length;

  return (
    <div className="space-y-4">
      {/* Header row */}
      <div className="flex items-start justify-between flex-wrap gap-4">
        <div>
          <p className="text-caption text-[var(--text-muted)] mb-1">Compliance &middot; Calendar</p>
          <MonthNavigation year={year} month={month} onPrev={goPrev} onNext={goNext} onToday={goToday} />
        </div>
        <StatusLegend />
      </div>

      {/* Category filters */}
      <div className="flex items-center gap-2">
        {(["direct_tax", "indirect_tax", "statutory"] as EventCategory[]).map((cat) => (
          <ComplianceCategoryBadge
            key={cat}
            category={cat}
            active={categoryFilters[cat]}
            onClick={() => toggleCategory(cat)}
          />
        ))}
      </div>

      {/* Grid */}
      <div className="border border-[var(--border-default)] rounded-xl overflow-hidden">
        {/* Weekday headers */}
        <div className="grid grid-cols-7 border-b border-[var(--border-default)]">
          {WEEKDAY_HEADERS.map((d) => (
            <div key={d} className="text-center text-caption text-[var(--text-muted)] py-2 font-semibold bg-[var(--bg-secondary)]">
              {d}
            </div>
          ))}
        </div>

        {/* Date cells */}
        <div className="grid grid-cols-7" style={{ gridAutoRows: "110px" }}>
          {gridCells.map((cellDate, i) => {
            const isCurrentMonth = cellDate.getMonth() === month && cellDate.getFullYear() === year;
            const dayEvents = filteredEvents.filter((e) => isSameDay(e.dueDate, cellDate));
            return (
              <ComplianceDayCell
                key={i}
                date={cellDate}
                events={dayEvents}
                isToday={isSameDay(cellDate, today)}
                isCurrentMonth={isCurrentMonth}
                onClick={() => setSelectedDate(cellDate)}
              />
            );
          })}
        </div>
      </div>

      {/* Stat cards */}
      <div className="grid grid-cols-3 gap-3">
        <ComplianceStatCard label="Filed" count={filedCount} variant="success" />
        <ComplianceStatCard label="Due in 7 days" count={dueSoonCount} variant="warning" />
        <ComplianceStatCard label="Upcoming this month" count={upcomingCount} variant="neutral" />
      </div>

      {/* Day detail slide-out */}
      <ComplianceDayDetailPanel
        date={selectedDate ?? today}
        events={selectedEvents}
        open={selectedDate !== null}
        onClose={() => setSelectedDate(null)}
      />
    </div>
  );
}
