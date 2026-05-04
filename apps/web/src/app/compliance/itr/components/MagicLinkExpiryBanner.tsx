"use client";

import { cn } from "@/lib/utils";
import { Clock, AlertTriangle } from "lucide-react";

interface MagicLinkExpiryBannerProps {
  expiresAt: string;
  className?: string;
}

function formatRelativeExpiry(expiresAt: string): { text: string; urgent: boolean } {
  const expiry = new Date(expiresAt);
  const now = new Date();
  const diffMs = expiry.getTime() - now.getTime();

  if (diffMs <= 0) return { text: "Expired", urgent: true };

  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  const diffHours = Math.floor((diffMs % (1000 * 60 * 60 * 24)) / (1000 * 60 * 60));

  if (diffDays > 1) return { text: `Expires in ${diffDays} days`, urgent: false };
  if (diffDays === 1) return { text: `Expires tomorrow`, urgent: false };
  if (diffHours > 1) return { text: `Expires in ${diffHours} hours`, urgent: true };
  return { text: "Expires soon", urgent: true };
}

function formatDate(iso: string): string {
  return new Date(iso).toLocaleDateString("en-IN", {
    day: "2-digit",
    month: "short",
    year: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function MagicLinkExpiryBanner({ expiresAt, className }: MagicLinkExpiryBannerProps) {
  const { text, urgent } = formatRelativeExpiry(expiresAt);

  return (
    <div
      className={cn(
        "flex items-center gap-2 px-4 py-2.5 rounded-lg border text-xs font-medium",
        urgent
          ? "bg-[var(--warning-muted)] border-[var(--warning)] text-[var(--warning)]"
          : "bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-secondary)]",
        className
      )}
    >
      {urgent ? (
        <AlertTriangle className="w-3.5 h-3.5 flex-shrink-0" />
      ) : (
        <Clock className="w-3.5 h-3.5 flex-shrink-0" />
      )}
      <span>{text}</span>
      <span className="text-[var(--text-muted)] ml-auto text-[10px]">
        {formatDate(expiresAt)}
      </span>
    </div>
  );
}
