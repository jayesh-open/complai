"use client";

import { useState, useMemo, useRef, useEffect } from "react";
import { Search } from "lucide-react";
import { cn } from "@complai/ui-components";
import type { Deductee } from "../types";

interface DeducteeSearchComboboxProps {
  deductees: Deductee[];
  value: Deductee | null;
  onChange: (deductee: Deductee | null) => void;
  className?: string;
}

export function DeducteeSearchCombobox({
  deductees,
  value,
  onChange,
  className,
}: DeducteeSearchComboboxProps) {
  const [query, setQuery] = useState("");
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  const filtered = useMemo(() => {
    if (!query) return deductees.slice(0, 10);
    const q = query.toLowerCase();
    return deductees
      .filter((d) => d.pan.toLowerCase().includes(q) || d.name.toLowerCase().includes(q))
      .slice(0, 10);
  }, [query, deductees]);

  useEffect(() => {
    function handleClick(e: MouseEvent) {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    }
    document.addEventListener("mousedown", handleClick);
    return () => document.removeEventListener("mousedown", handleClick);
  }, []);

  return (
    <div ref={ref} className={cn("relative", className)}>
      <div className="relative">
        <Search className="absolute left-3 top-1/2 -translate-y-1/2 w-4 h-4 text-[var(--text-muted)]" />
        <input
          type="text"
          value={value ? `${value.pan} — ${value.name}` : query}
          onChange={(e) => {
            setQuery(e.target.value);
            if (value) onChange(null);
            setOpen(true);
          }}
          onFocus={() => setOpen(true)}
          placeholder="Search by PAN or name..."
          className={cn(
            "w-full pl-9 pr-3 py-2 rounded-lg text-xs",
            "bg-[var(--bg-secondary)] border border-[var(--border-default)]",
            "text-[var(--text-primary)] placeholder:text-[var(--text-muted)]",
            "focus:outline-none focus:border-[var(--accent)]",
            "focus:ring-2 focus:ring-[var(--accent-muted)]"
          )}
        />
      </div>
      {open && !value && (
        <div
          className={cn(
            "absolute z-50 mt-1 w-full rounded-lg border border-[var(--border-default)]",
            "bg-[var(--bg-secondary)] shadow-[var(--shadow-lg)]",
            "max-h-64 overflow-y-auto"
          )}
        >
          {filtered.length === 0 ? (
            <div className="px-3 py-2 text-xs text-[var(--text-muted)]">
              No deductees found
            </div>
          ) : (
            filtered.map((d) => (
              <button
                key={d.id}
                onClick={() => {
                  onChange(d);
                  setQuery("");
                  setOpen(false);
                }}
                className={cn(
                  "w-full text-left px-3 py-2 text-xs",
                  "hover:bg-[var(--bg-tertiary)] transition-colors",
                  "border-b border-[var(--border-default)] last:border-b-0"
                )}
              >
                <div className="flex items-center justify-between">
                  <div>
                    <span className="font-mono text-[var(--text-primary)] font-medium">
                      {d.pan}
                    </span>
                    <span className="text-[var(--text-muted)] ml-2">{d.name}</span>
                  </div>
                  <span className="text-[10px] text-[var(--text-muted)] uppercase">
                    {d.category}
                  </span>
                </div>
              </button>
            ))
          )}
        </div>
      )}
    </div>
  );
}
