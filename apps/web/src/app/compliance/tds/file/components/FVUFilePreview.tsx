import { FileText, Copy } from "lucide-react";
import { cn } from "@complai/ui-components";

interface FVUFilePreviewProps {
  content: string;
  formLabel: string;
}

export function FVUFilePreview({ content, formLabel }: FVUFilePreviewProps) {
  const handleCopy = () => {
    navigator.clipboard.writeText(content);
  };

  return (
    <div className="space-y-3">
      <div className="flex items-center justify-between">
        <div className="flex items-center gap-2">
          <FileText className="w-4 h-4 text-[var(--text-muted)]" />
          <span className="text-xs font-semibold text-[var(--text-primary)]">
            FVU Preview — {formLabel}
          </span>
        </div>
        <button
          onClick={handleCopy}
          className={cn(
            "flex items-center gap-1 px-2 py-1 rounded text-[10px] font-medium",
            "border border-[var(--border-default)] text-[var(--text-muted)]",
            "hover:text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)] transition-colors",
          )}
        >
          <Copy className="w-3 h-3" />
          Copy
        </button>
      </div>
      <pre className={cn(
        "p-4 rounded-lg font-mono text-[11px] leading-relaxed overflow-x-auto",
        "bg-[var(--bg-tertiary)] border border-[var(--border-default)]",
        "text-[var(--text-secondary)] select-all",
      )}>
        {content}
      </pre>
    </div>
  );
}
