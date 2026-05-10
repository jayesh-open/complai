"use client";

import type { EventCategory } from "../types";
import { cn } from "@/lib/utils";

const CATEGORY_STYLES: Record<EventCategory, { bg: string; text: string; label: string }> = {
  direct_tax: { bg: "#E6F1FB", text: "#0C447C", label: "Direct Tax (IT)" },
  indirect_tax: { bg: "#EEEDFE", text: "#3C3489", label: "Indirect Tax (GST)" },
  statutory: { bg: "#F1EFE8", text: "#444441", label: "Statutory (PF/ESI/ROC)" },
};

export function ComplianceCategoryBadge({
  category,
  active = true,
  onClick,
  size = "default",
}: {
  category: EventCategory;
  active?: boolean;
  onClick?: () => void;
  size?: "default" | "small";
}) {
  const style = CATEGORY_STYLES[category];

  return (
    <button
      type="button"
      onClick={onClick}
      data-testid={`category-${category}`}
      className={cn(
        "inline-flex items-center rounded-md font-medium transition-all duration-150",
        size === "small" ? "px-2 py-0.5 text-[10px]" : "px-3 py-1 text-xs",
        active ? "opacity-100" : "opacity-40",
        onClick && "cursor-pointer hover:opacity-80",
        !onClick && "cursor-default",
      )}
      style={{ backgroundColor: style.bg, color: style.text }}
    >
      {style.label}
    </button>
  );
}

export function getCategoryStyle(category: EventCategory) {
  return CATEGORY_STYLES[category];
}
