"use client";

import { useState } from "react";
import { ChevronDown, ChevronRight, Copy, Check } from "lucide-react";
import { cn } from "@/lib/utils";

interface SignedJsonViewerProps {
  json: string;
  className?: string;
}

export function SignedJsonViewer({ json, className }: SignedJsonViewerProps) {
  const [expanded, setExpanded] = useState(false);
  const [copied, setCopied] = useState(false);

  const handleCopy = async () => {
    await navigator.clipboard.writeText(json);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div
      data-testid="signed-json-viewer"
      className={cn(
        "rounded-xl border border-[var(--border-default)] overflow-hidden",
        className
      )}
    >
      <button
        onClick={() => setExpanded(!expanded)}
        className={cn(
          "flex items-center justify-between w-full px-4 py-3",
          "bg-[var(--bg-tertiary)] hover:bg-[var(--bg-secondary)]",
          "transition-colors"
        )}
      >
        <div className="flex items-center gap-2">
          {expanded ? (
            <ChevronDown className="w-4 h-4 text-[var(--text-muted)]" />
          ) : (
            <ChevronRight className="w-4 h-4 text-[var(--text-muted)]" />
          )}
          <span className="text-xs font-semibold text-[var(--text-primary)]">
            Signed Invoice JSON
          </span>
        </div>
        <button
          onClick={(e) => {
            e.stopPropagation();
            handleCopy();
          }}
          className={cn(
            "flex items-center gap-1 px-2 py-1 rounded text-[10px]",
            "border border-[var(--border-default)]",
            "text-[var(--text-muted)] hover:text-[var(--text-primary)]",
            "hover:bg-[var(--bg-secondary)] transition-colors"
          )}
        >
          {copied ? (
            <Check className="w-3 h-3 text-[var(--success)]" />
          ) : (
            <Copy className="w-3 h-3" />
          )}
          {copied ? "Copied" : "Copy"}
        </button>
      </button>
      {expanded && (
        <div className="p-4 bg-[var(--bg-primary)] overflow-x-auto max-h-[400px] overflow-y-auto">
          <pre className="text-[11px] font-mono text-[var(--text-secondary)] whitespace-pre leading-relaxed">
            {json}
          </pre>
        </div>
      )}
    </div>
  );
}
