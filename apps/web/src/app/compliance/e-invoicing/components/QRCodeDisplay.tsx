"use client";

import { QRCodeSVG } from "qrcode.react";
import { cn } from "@/lib/utils";

interface QRCodeDisplayProps {
  value: string;
  size?: number;
  className?: string;
  label?: string;
}

export function QRCodeDisplay({
  value,
  size = 160,
  className,
  label,
}: QRCodeDisplayProps) {
  return (
    <div
      className={cn(
        "flex flex-col items-center gap-2 p-4 rounded-xl",
        "bg-white border border-[var(--border-default)]",
        className
      )}
      data-testid="qr-code-display"
    >
      <QRCodeSVG
        value={value}
        size={size}
        level="M"
        includeMargin={false}
      />
      {label && (
        <span className="text-[10px] text-[var(--text-muted)] font-mono truncate max-w-full">
          {label}
        </span>
      )}
    </div>
  );
}
