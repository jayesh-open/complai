"use client";

import { useState, useEffect, useRef } from "react";
import { Eye, ChevronDown, LogOut } from "lucide-react";
import { cn } from "@/lib/utils";
import { useViewAsStore } from "@/store/view-as-store";
import type { Role, Permission, RoleDetail } from "@/app/configure/users/types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

export function ViewAsDropdown() {
  const { isViewAs, viewAsRole, enterViewAs, exitViewAs } = useViewAsStore();
  const [open, setOpen] = useState(false);
  const [roles, setRoles] = useState<Role[]>([]);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    fetch(`/api/v1/roles?tenant_id=${DEV_TENANT}`)
      .then((r) => r.json())
      .then((data) => setRoles(Array.isArray(data) ? data.filter((r: Role) => r.name !== "admin") : []))
      .catch(() => {});
  }, []);

  useEffect(() => {
    if (!open) return;
    const onClick = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) setOpen(false);
    };
    document.addEventListener("mousedown", onClick);
    return () => document.removeEventListener("mousedown", onClick);
  }, [open]);

  const handleSelect = async (role: Role) => {
    setOpen(false);
    try {
      const res = await fetch(`/api/v1/roles/${role.id}?tenant_id=${DEV_TENANT}`);
      if (!res.ok) return;
      const data: RoleDetail = await res.json();
      enterViewAs(data.role, data.permissions);
    } catch {
      // silent fail
    }
  };

  const nonAdminRoles = roles.filter((r) => r.name !== "admin");

  return (
    <div ref={ref} className="relative">
      <button
        onClick={() => setOpen(!open)}
        className={cn(
          "flex items-center gap-1.5 px-3 py-1.5 rounded-lg text-xs font-medium border transition-colors",
          isViewAs
            ? "border-[var(--warning)]/40 bg-[var(--warning)]/10 text-[var(--warning)]"
            : "border-[var(--border-default)] text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)]",
        )}
        data-testid="view-as-trigger"
      >
        <Eye className="w-3.5 h-3.5" />
        {isViewAs ? `Viewing as ${viewAsRole?.display_name}` : "View As"}
        <ChevronDown className="w-3 h-3" />
      </button>

      {open && (
        <div className="absolute right-0 top-full mt-1 z-50 w-56 rounded-lg border border-[var(--border-default)] bg-[var(--bg-primary)] shadow-lg py-1">
          {isViewAs && (
            <>
              <button
                onClick={() => { setOpen(false); exitViewAs(); }}
                className="w-full flex items-center gap-2 px-3 py-2 text-xs text-[var(--warning)] hover:bg-[var(--bg-tertiary)]"
              >
                <LogOut className="w-3.5 h-3.5" />
                Exit View As
              </button>
              <div className="border-t border-[var(--border-default)] my-1" />
              <span className="block px-3 py-1 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">
                Switch to
              </span>
            </>
          )}
          {nonAdminRoles.map((role) => (
            <button
              key={role.id}
              data-testid={`view-as-role-${role.name}`}
              onClick={() => handleSelect(role)}
              className={cn(
                "w-full text-left px-3 py-1.5 text-xs hover:bg-[var(--bg-tertiary)]",
                viewAsRole?.id === role.id
                  ? "text-[var(--accent)] font-medium"
                  : "text-[var(--text-primary)]",
              )}
            >
              {role.display_name}
            </button>
          ))}
          {nonAdminRoles.length === 0 && (
            <span className="block px-3 py-2 text-xs text-[var(--text-muted)]">Loading roles...</span>
          )}
        </div>
      )}
    </div>
  );
}
