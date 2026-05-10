"use client";

import { useState, useCallback, useEffect, useRef } from "react";
import Link from "next/link";
import { ArrowLeft, Plus, RefreshCw } from "lucide-react";
import { RolesTable } from "../components/RolesTable";
import { AddCustomRoleModal } from "../components/AddCustomRoleModal";
import type { Role } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

export default function RolesPage() {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [showAddModal, setShowAddModal] = useState(false);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const controllerRef = useRef<AbortController | null>(null);

  const fetchRoles = useCallback(async () => {
    if (retryRef.current) clearTimeout(retryRef.current);
    controllerRef.current?.abort();
    const controller = new AbortController();
    controllerRef.current = controller;
    setLoading(true);
    try {
      const res = await fetch(
        `/api/v1/roles?tenant_id=${encodeURIComponent(DEV_TENANT)}`,
        { signal: controller.signal },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setRoles(Array.isArray(data) ? data : []);
      setError(null);
    } catch (err) {
      if (controller.signal.aborted) return;
      setError(err instanceof Error ? err.message : "Failed to fetch roles");
      retryRef.current = setTimeout(fetchRoles, 30_000);
    } finally {
      if (!controller.signal.aborted) setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchRoles();
    return () => {
      controllerRef.current?.abort();
      if (retryRef.current) clearTimeout(retryRef.current);
    };
  }, [fetchRoles]);

  const handleDelete = async (role: Role) => {
    if (!confirm(`Delete custom role "${role.display_name}"? This cannot be undone.`)) return;
    try {
      const res = await fetch(
        `/api/v1/roles/${role.id}?tenant_id=${DEV_TENANT}`,
        { method: "DELETE" },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      fetchRoles();
    } catch {
      // Refresh to show current state
      fetchRoles();
    }
  };

  return (
    <div className="p-7 max-w-[1280px]" data-testid="roles-page">
      <Link
        href="/configure/users"
        className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] mb-4"
      >
        <ArrowLeft className="w-3.5 h-3.5" />
        Back to Users
      </Link>

      <div className="flex items-start justify-between mb-6">
        <div>
          <h2 className="text-heading-lg text-[var(--text-primary)]">Roles</h2>
          <p className="text-body-sm text-[var(--text-muted)] mt-0.5">
            View and manage tenant roles. System roles cannot be deleted but can be cloned.
          </p>
        </div>
        <button
          onClick={() => setShowAddModal(true)}
          className="flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90"
        >
          <Plus className="w-3.5 h-3.5" />
          Add Custom Role
        </button>
      </div>

      <div className="flex gap-1 mb-5 border-b border-[var(--border-default)]">
        <Link
          href="/configure/users"
          className="px-4 py-2 text-xs font-medium border-b-2 border-transparent text-[var(--text-muted)] hover:text-[var(--text-primary)] -mb-px"
        >
          Users
        </Link>
        <span className="px-4 py-2 text-xs font-medium border-b-2 border-[var(--accent)] text-[var(--accent)] -mb-px">
          Roles
        </span>
      </div>

      {loading && (
        <div className="bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden">
          <div className="animate-pulse space-y-0">
            <div className="h-10 bg-[var(--bg-tertiary)] border-b border-[var(--border-default)]" />
            {Array.from({ length: 5 }).map((_, i) => (
              <div key={i} className="h-12 border-b border-[var(--border-default)] last:border-b-0 flex items-center px-[18px] gap-6">
                <div className="h-3 w-32 rounded bg-[var(--bg-tertiary)]" />
                <div className="h-3 w-44 rounded bg-[var(--bg-tertiary)]" />
                <div className="h-3 w-16 rounded bg-[var(--bg-tertiary)]" />
                <div className="h-3 w-12 rounded bg-[var(--bg-tertiary)]" />
              </div>
            ))}
          </div>
        </div>
      )}

      {error && !loading && (
        <div className="rounded-[14px] border border-[var(--danger)]/30 bg-[var(--danger)]/5 p-6 text-center">
          <p className="text-sm text-[var(--text-primary)] mb-3">Could not load roles</p>
          <button
            onClick={fetchRoles}
            className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
          >
            <RefreshCw className="w-3.5 h-3.5" />
            Retry
          </button>
        </div>
      )}

      {!loading && !error && (
        <RolesTable roles={roles} onDelete={handleDelete} />
      )}

      <AddCustomRoleModal
        open={showAddModal}
        roles={roles}
        onClose={() => setShowAddModal(false)}
        onCreated={fetchRoles}
      />
    </div>
  );
}
