"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { MoreHorizontal } from "lucide-react";
import { cn } from "@/lib/utils";
import { RoleTypeBadge } from "./RoleTypeBadge";
import type { Role } from "../types";

type SortDir = "asc" | "desc";

interface RolesTableProps {
  roles: Role[];
  onDelete?: (role: Role) => void;
  className?: string;
}

const TH =
  "px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left select-none";

export function RolesTable({ roles, onDelete, className }: RolesTableProps) {
  const router = useRouter();
  const [sortDir, setSortDir] = useState<SortDir>("asc");
  const [menuOpen, setMenuOpen] = useState<string | null>(null);

  const sorted = [...roles].sort((a, b) => {
    const cmp = a.display_name.toLowerCase() < b.display_name.toLowerCase() ? -1 : 1;
    return sortDir === "asc" ? cmp : -cmp;
  });

  const arrow = sortDir === "asc" ? " ↑" : " ↓";

  return (
    <div
      className={cn(
        "bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden",
        className,
      )}
    >
      <table className="w-full">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            <th
              className={cn(TH, "cursor-pointer")}
              onClick={() => setSortDir((d) => (d === "asc" ? "desc" : "asc"))}
            >
              Role{arrow}
            </th>
            <th className={TH}>Description</th>
            <th className={TH}>Type</th>
            <th className={TH}>Members</th>
            <th className={cn(TH, "w-12")} />
          </tr>
        </thead>
        <tbody>
          {sorted.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center py-12 text-[var(--text-muted)] text-sm">
                No roles found. System roles should be present — try refreshing.
              </td>
            </tr>
          ) : (
            sorted.map((role) => (
              <tr
                key={role.id}
                onClick={() => router.push(`/configure/users/roles/${role.id}`)}
                className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors cursor-pointer"
              >
                <td className="px-[18px] py-3 text-xs font-medium text-[var(--text-primary)]">
                  {role.display_name}
                </td>
                <td className="px-[18px] py-3 text-xs text-[var(--text-secondary)] max-w-[300px] truncate">
                  {role.description ?? "—"}
                </td>
                <td className="px-[18px] py-3">
                  <RoleTypeBadge isSystem={role.is_system} />
                </td>
                <td className="px-[18px] py-3 text-xs text-[var(--text-muted)]">—</td>
                <td className="px-[18px] py-3 relative">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      setMenuOpen(menuOpen === role.id ? null : role.id);
                    }}
                    className="p-1 rounded-md hover:bg-[var(--bg-tertiary)] text-[var(--text-muted)]"
                    aria-label="Actions"
                  >
                    <MoreHorizontal className="w-4 h-4" />
                  </button>
                  {menuOpen === role.id && (
                    <div
                      className="absolute right-4 top-10 z-30 w-44 rounded-lg border border-[var(--border-default)] bg-[var(--bg-primary)] shadow-lg py-1"
                      onMouseLeave={() => setMenuOpen(null)}
                    >
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setMenuOpen(null);
                          router.push(`/configure/users/roles/${role.id}`);
                        }}
                        className="w-full text-left px-3 py-1.5 text-xs text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
                      >
                        View Details
                      </button>
                      {!role.is_system && onDelete && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setMenuOpen(null);
                            onDelete(role);
                          }}
                          className="w-full text-left px-3 py-1.5 text-xs text-[var(--danger)] hover:bg-[var(--bg-tertiary)]"
                        >
                          Delete
                        </button>
                      )}
                      {role.is_system && (
                        <span className="block px-3 py-1.5 text-[10px] text-[var(--text-disabled)] cursor-default">
                          System roles cannot be deleted
                        </span>
                      )}
                    </div>
                  )}
                </td>
              </tr>
            ))
          )}
        </tbody>
      </table>
    </div>
  );
}
