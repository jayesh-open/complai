"use client";

import { useState } from "react";
import { Plus, RefreshCw } from "lucide-react";
import { cn } from "@/lib/utils";
import { useUsers } from "./hooks/useUsers";
import { useRoles } from "./hooks/useRoles";
import { UsersTable } from "./components/UsersTable";
import { AddUserModal } from "./components/AddUserModal";
import { UserDetailPanel } from "./components/UserDetailPanel";
import type { User } from "./types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

export default function UsersPage() {
  const { users, loading, error, mutate } = useUsers(DEV_TENANT);
  const { roles } = useRoles(DEV_TENANT);
  const [showAddModal, setShowAddModal] = useState(false);
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [activeTab, setActiveTab] = useState<"users" | "roles">("users");

  const handleDeactivate = async (user: User) => {
    if (!confirm(`Deactivate ${user.first_name} ${user.last_name}? They will lose access.`)) return;
    try {
      const res = await fetch(
        `/api/v1/users/${user.id}/deactivate?tenant_id=${DEV_TENANT}`,
        { method: "POST" },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      mutate();
    } catch {
      // Error handled visually by the user list refetching
    }
  };

  return (
    <div className="p-7 max-w-[1280px]" data-testid="users-page">
      <div className="flex items-start justify-between mb-6">
        <div>
          <h2 className="text-heading-lg text-[var(--text-primary)]">Users &amp; Roles</h2>
          <p className="text-body-sm text-[var(--text-muted)] mt-0.5">
            Manage tenant users and their role assignments
          </p>
        </div>
        {activeTab === "users" && (
          <button
            onClick={() => setShowAddModal(true)}
            className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90"
          >
            <Plus className="w-3.5 h-3.5" />
            Add User
          </button>
        )}
      </div>

      {/* Tabs */}
      <div className="flex gap-1 mb-5 border-b border-[var(--border-default)]">
        <button
          onClick={() => setActiveTab("users")}
          className={cn(
            "px-4 py-2 text-xs font-medium border-b-2 transition-colors -mb-px",
            activeTab === "users"
              ? "border-[var(--accent)] text-[var(--accent)]"
              : "border-transparent text-[var(--text-muted)] hover:text-[var(--text-primary)]",
          )}
        >
          Users
        </button>
        <button
          onClick={() => setActiveTab("roles")}
          className={cn(
            "px-4 py-2 text-xs font-medium border-b-2 transition-colors -mb-px",
            activeTab === "roles"
              ? "border-[var(--accent)] text-[var(--accent)]"
              : "border-transparent text-[var(--text-muted)] hover:text-[var(--text-primary)]",
          )}
        >
          Roles
        </button>
      </div>

      {activeTab === "users" ? (
        <>
          {loading && (
            <div className="bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden">
              <div className="animate-pulse space-y-0">
                <div className="h-10 bg-[var(--bg-tertiary)] border-b border-[var(--border-default)]" />
                {Array.from({ length: 5 }).map((_, i) => (
                  <div key={i} className="h-12 border-b border-[var(--border-default)] last:border-b-0 flex items-center px-[18px] gap-6">
                    <div className="h-3 w-32 rounded bg-[var(--bg-tertiary)]" />
                    <div className="h-3 w-44 rounded bg-[var(--bg-tertiary)]" />
                    <div className="h-3 w-20 rounded bg-[var(--bg-tertiary)]" />
                    <div className="h-3 w-16 rounded bg-[var(--bg-tertiary)]" />
                  </div>
                ))}
              </div>
            </div>
          )}

          {error && !loading && (
            <div className="rounded-[14px] border border-[var(--danger)]/30 bg-[var(--danger)]/5 p-6 text-center">
              <p className="text-sm text-[var(--text-primary)] mb-3">Could not load users</p>
              <button
                onClick={mutate}
                className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
              >
                <RefreshCw className="w-3.5 h-3.5" />
                Retry
              </button>
            </div>
          )}

          {!loading && !error && (
            <UsersTable
              users={users}
              onSelectUser={setSelectedUser}
              onDeactivate={handleDeactivate}
            />
          )}
        </>
      ) : (
        <div className="rounded-[14px] border border-[var(--border-default)] bg-[var(--bg-secondary)] p-8 text-center">
          <p className="text-sm text-[var(--text-muted)]">
            Roles management coming in Sub-part D-2.
          </p>
        </div>
      )}

      <AddUserModal
        open={showAddModal}
        roles={roles}
        onClose={() => setShowAddModal(false)}
        onCreated={mutate}
      />

      {selectedUser && (
        <UserDetailPanel
          user={selectedUser}
          roles={roles}
          open={!!selectedUser}
          onClose={() => setSelectedUser(null)}
          onUpdated={() => {
            mutate();
            setSelectedUser(null);
          }}
        />
      )}
    </div>
  );
}
