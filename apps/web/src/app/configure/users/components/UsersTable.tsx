"use client";

import { useState } from "react";
import { MoreHorizontal } from "lucide-react";
import { cn } from "@/lib/utils";
import { UserStatusPill } from "./UserStatusPill";
import type { User } from "../types";

type SortField = "name" | "email" | "role";
type SortDir = "asc" | "desc";

interface UsersTableProps {
  users: User[];
  onSelectUser: (user: User) => void;
  onDeactivate: (user: User) => void;
  className?: string;
}

function sortUsers(users: User[], field: SortField, dir: SortDir): User[] {
  const sorted = [...users].sort((a, b) => {
    let av: string;
    let bv: string;
    if (field === "name") {
      av = `${a.first_name} ${a.last_name}`.toLowerCase();
      bv = `${b.first_name} ${b.last_name}`.toLowerCase();
    } else if (field === "email") {
      av = a.email.toLowerCase();
      bv = b.email.toLowerCase();
    } else {
      av = (a.role?.display_name ?? "").toLowerCase();
      bv = (b.role?.display_name ?? "").toLowerCase();
    }
    return av < bv ? -1 : av > bv ? 1 : 0;
  });
  return dir === "desc" ? sorted.reverse() : sorted;
}

const TH =
  "px-[18px] py-[10px] text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left select-none";

export function UsersTable({ users, onSelectUser, onDeactivate, className }: UsersTableProps) {
  const [sortField, setSortField] = useState<SortField>("name");
  const [sortDir, setSortDir] = useState<SortDir>("asc");
  const [menuOpen, setMenuOpen] = useState<string | null>(null);

  const toggleSort = (field: SortField) => {
    if (sortField === field) {
      setSortDir((d) => (d === "asc" ? "desc" : "asc"));
    } else {
      setSortField(field);
      setSortDir("asc");
    }
  };

  const sorted = sortUsers(users, sortField, sortDir);
  const arrow = (f: SortField) =>
    sortField === f ? (sortDir === "asc" ? " ↑" : " ↓") : "";

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
            <th className={cn(TH, "cursor-pointer")} onClick={() => toggleSort("name")}>
              Name{arrow("name")}
            </th>
            <th className={cn(TH, "cursor-pointer")} onClick={() => toggleSort("email")}>
              Email{arrow("email")}
            </th>
            <th className={cn(TH, "cursor-pointer")} onClick={() => toggleSort("role")}>
              Role{arrow("role")}
            </th>
            <th className={TH}>Status</th>
            <th className={cn(TH, "w-12")} />
          </tr>
        </thead>
        <tbody>
          {sorted.length === 0 ? (
            <tr>
              <td colSpan={5} className="text-center py-12 text-[var(--text-muted)] text-sm">
                No users yet. Add the first user to get started.
              </td>
            </tr>
          ) : (
            sorted.map((user) => (
              <tr
                key={user.id}
                onClick={() => onSelectUser(user)}
                className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors cursor-pointer"
              >
                <td className="px-[18px] py-3 text-xs font-medium text-[var(--text-primary)]">
                  {user.first_name} {user.last_name}
                </td>
                <td className="px-[18px] py-3 text-xs text-[var(--text-secondary)]">
                  {user.email}
                </td>
                <td className="px-[18px] py-3">
                  {user.role ? (
                    <span className="inline-block text-[10px] font-medium px-2 py-0.5 rounded-full bg-[var(--accent)]/10 text-[var(--accent)]">
                      {user.role.display_name}
                    </span>
                  ) : (
                    <span className="text-[10px] text-[var(--text-muted)]">No role</span>
                  )}
                </td>
                <td className="px-[18px] py-3">
                  <UserStatusPill status={user.status} />
                </td>
                <td className="px-[18px] py-3 relative">
                  <button
                    onClick={(e) => {
                      e.stopPropagation();
                      setMenuOpen(menuOpen === user.id ? null : user.id);
                    }}
                    className="p-1 rounded-md hover:bg-[var(--bg-tertiary)] text-[var(--text-muted)]"
                    aria-label="Actions"
                  >
                    <MoreHorizontal className="w-4 h-4" />
                  </button>
                  {menuOpen === user.id && (
                    <div
                      className="absolute right-4 top-10 z-30 w-36 rounded-lg border border-[var(--border-default)] bg-[var(--bg-primary)] shadow-lg py-1"
                      onMouseLeave={() => setMenuOpen(null)}
                    >
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setMenuOpen(null);
                          onSelectUser(user);
                        }}
                        className="w-full text-left px-3 py-1.5 text-xs text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
                      >
                        Edit
                      </button>
                      {user.status === "active" && (
                        <button
                          onClick={(e) => {
                            e.stopPropagation();
                            setMenuOpen(null);
                            onDeactivate(user);
                          }}
                          className="w-full text-left px-3 py-1.5 text-xs text-[var(--danger)] hover:bg-[var(--bg-tertiary)]"
                        >
                          Deactivate
                        </button>
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
