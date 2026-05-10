"use client";

import {
  Eye, Pencil, Upload, PlusCircle, XCircle,
  Calculator, CheckCircle, FileText, Settings,
} from "lucide-react";
import { cn } from "@/lib/utils";

const ACTION_CONFIG: Record<string, { icon: typeof Eye; label: string; activeColor: string }> = {
  view: { icon: Eye, label: "View", activeColor: "text-[var(--info)]" },
  edit: { icon: Pencil, label: "Edit", activeColor: "text-[var(--warning)]" },
  file: { icon: Upload, label: "File / Submit", activeColor: "text-[var(--success)]" },
  generate: { icon: PlusCircle, label: "Generate", activeColor: "text-[var(--success)]" },
  cancel: { icon: XCircle, label: "Cancel", activeColor: "text-[var(--danger)]" },
  calculate: { icon: Calculator, label: "Calculate", activeColor: "text-[var(--info)]" },
  approve: { icon: CheckCircle, label: "Approve", activeColor: "text-[var(--success)]" },
  issue_cert: { icon: FileText, label: "Issue Certificate", activeColor: "text-[var(--info)]" },
  manage: { icon: Settings, label: "Manage", activeColor: "text-[var(--warning)]" },
};

interface PermissionIconProps {
  action: string;
  granted: boolean;
  disabled: boolean;
  onClick: () => void;
}

export function PermissionIcon({ action, granted, disabled, onClick }: PermissionIconProps) {
  const config = ACTION_CONFIG[action];
  if (!config) return null;

  const Icon = config.icon;

  const tooltip = disabled
    ? "System role permissions cannot be modified"
    : granted
      ? `${config.label} access granted`
      : `${config.label} access not granted`;

  return (
    <button
      type="button"
      onClick={disabled ? undefined : onClick}
      title={tooltip}
      className={cn(
        "p-1 rounded transition-colors",
        granted ? config.activeColor : "text-[var(--text-muted)] opacity-40",
        disabled
          ? "cursor-not-allowed"
          : "cursor-pointer hover:bg-[var(--bg-tertiary)]",
      )}
      disabled={disabled}
      aria-label={tooltip}
    >
      <Icon className="w-4 h-4" />
    </button>
  );
}
