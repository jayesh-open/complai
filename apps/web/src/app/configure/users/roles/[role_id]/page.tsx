"use client";

import { useState, useRef, useCallback } from "react";
import { useParams } from "next/navigation";
import Link from "next/link";
import { ArrowLeft, RefreshCw, Info } from "lucide-react";
import { RoleTypeBadge } from "../../components/RoleTypeBadge";
import { PermissionMatrix, type PermissionPair } from "../../components/PermissionMatrix";
import { useRoleDetail } from "../../hooks/useRoleDetail";
import { useUpdateRolePermissions } from "../../hooks/useUpdateRolePermissions";

export default function RoleDetailPage() {
  const { role_id } = useParams<{ role_id: string }>();
  const { role, permissions, loading, error, mutate } = useRoleDetail(role_id);
  const { updatePermissions, isUpdating, error: updateError } = useUpdateRolePermissions();
  const [dirty, setDirty] = useState(false);
  const pendingPairsRef = useRef<PermissionPair[]>([]);
  const [saveSuccess, setSaveSuccess] = useState(false);

  const handleMatrixChange = useCallback((pairs: PermissionPair[]) => {
    pendingPairsRef.current = pairs;
    setDirty(true);
    setSaveSuccess(false);
  }, []);

  const handleSave = async () => {
    const ok = await updatePermissions(role_id, pendingPairsRef.current);
    if (ok) {
      setDirty(false);
      setSaveSuccess(true);
      mutate();
    }
  };

  const handleDiscard = () => {
    setDirty(false);
    setSaveSuccess(false);
    mutate();
  };

  if (loading) {
    return (
      <div className="p-7 max-w-[1280px]" data-testid="role-detail-page">
        <div className="animate-pulse space-y-4">
          <div className="h-4 w-24 rounded bg-[var(--bg-tertiary)]" />
          <div className="h-6 w-48 rounded bg-[var(--bg-tertiary)]" />
          <div className="h-3 w-64 rounded bg-[var(--bg-tertiary)]" />
          <div className="h-[500px] rounded-[14px] bg-[var(--bg-tertiary)]" />
        </div>
      </div>
    );
  }

  if (error || !role) {
    return (
      <div className="p-7 max-w-[1280px]" data-testid="role-detail-page">
        <Link href="/configure/users/roles" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] mb-4">
          <ArrowLeft className="w-3.5 h-3.5" /> Back to Roles
        </Link>
        <div className="rounded-[14px] border border-[var(--danger)]/30 bg-[var(--danger)]/5 p-6 text-center">
          <p className="text-sm text-[var(--text-primary)] mb-3">{error ?? "Role not found"}</p>
          <button onClick={mutate} className="inline-flex items-center gap-1.5 px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]">
            <RefreshCw className="w-3.5 h-3.5" /> Retry
          </button>
        </div>
      </div>
    );
  }

  return (
    <div className="p-7 max-w-[1280px]" data-testid="role-detail-page">
      <Link href="/configure/users/roles" className="inline-flex items-center gap-1 text-xs text-[var(--text-muted)] hover:text-[var(--text-primary)] mb-4">
        <ArrowLeft className="w-3.5 h-3.5" /> Back to Roles
      </Link>

      <div className="flex items-start justify-between mb-6">
        <div>
          <div className="flex items-center gap-3 mb-1">
            <h2 className="text-heading-lg text-[var(--text-primary)]">{role.display_name}</h2>
            <RoleTypeBadge isSystem={role.is_system} />
          </div>
          <p className="text-[11px] font-mono text-[var(--text-muted)]">{role.name}</p>
          {role.description && (
            <p className="text-body-sm text-[var(--text-secondary)] mt-1">{role.description}</p>
          )}
        </div>
        <span className="text-xs text-[var(--text-muted)]">{permissions.length} permissions</span>
      </div>

      {role.is_system && (
        <div data-testid="system-role-banner" className="flex items-start gap-2.5 px-4 py-3 mb-5 rounded-lg bg-[var(--warning)]/10 border border-[var(--warning)]/30">
          <Info className="w-4 h-4 text-[var(--warning)] mt-0.5 shrink-0" />
          <p className="text-xs text-[var(--text-primary)]">
            This is a system role. Permissions cannot be modified directly. Clone this role to create a customizable copy.
          </p>
        </div>
      )}

      <PermissionMatrix
        permissions={permissions}
        isSystemRole={role.is_system}
        onChange={handleMatrixChange}
      />

      {(updateError || saveSuccess) && (
        <p className={`text-xs mt-3 ${updateError ? "text-[var(--danger)]" : "text-[var(--success)]"}`}>
          {updateError ?? "Permissions saved successfully."}
        </p>
      )}

      <div className="flex items-center gap-2 mt-5">
        <button
          onClick={handleSave}
          disabled={role.is_system || !dirty || isUpdating}
          className="px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90 disabled:opacity-50"
        >
          {isUpdating ? "Saving..." : "Save Changes"}
        </button>
        {dirty && !role.is_system && (
          <button
            onClick={handleDiscard}
            className="px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
          >
            Discard Changes
          </button>
        )}
      </div>
    </div>
  );
}
