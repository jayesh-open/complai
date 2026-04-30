import { cn } from "@complai/ui-components";

export function StatChip({
  icon: Icon,
  color,
  label,
  count,
}: {
  icon: React.ElementType;
  color: string;
  label: string;
  count: number;
}) {
  return (
    <div className="flex items-center gap-1.5">
      <Icon className={cn("w-4 h-4", color)} />
      <span className="text-xs font-medium text-[var(--text-primary)]">{count}</span>
      <span className="text-[10px] text-[var(--text-muted)] uppercase">{label}</span>
    </div>
  );
}
