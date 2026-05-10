"use client";

import { useState, useEffect, useRef } from "react";
import { X, Copy, Check } from "lucide-react";
import { cn } from "@/lib/utils";
import { UserStatusPill } from "./UserStatusPill";
import type { User, Role } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";
const INPUT_CLS = "w-full rounded-lg border border-[var(--border-default)] bg-[var(--bg-secondary)] px-3 py-2 text-xs text-[var(--text-primary)] outline-none focus:border-[var(--accent)]";

interface UserDetailPanelProps {
  user: User;
  roles: Role[];
  open: boolean;
  onClose: () => void;
  onUpdated: () => void;
}

export function UserDetailPanel({ user, roles, open, onClose, onUpdated }: UserDetailPanelProps) {
  const panelRef = useRef<HTMLDivElement>(null);
  const [editing, setEditing] = useState(false);
  const [firstName, setFirstName] = useState(user.first_name);
  const [lastName, setLastName] = useState(user.last_name);
  const [roleId, setRoleId] = useState(user.role?.id ?? "");
  const [saving, setSaving] = useState(false);
  const [saveError, setSaveError] = useState<string | null>(null);
  const [copied, setCopied] = useState(false);
  const [deactivating, setDeactivating] = useState(false);

  useEffect(() => {
    setEditing(false);
    setFirstName(user.first_name);
    setLastName(user.last_name);
    setRoleId(user.role?.id ?? "");
    setSaveError(null);
  }, [user]);

  useEffect(() => {
    if (!open) return;
    const onKey = (e: KeyboardEvent) => {
      if (e.key === "Escape") onClose();
    };
    document.addEventListener("keydown", onKey);
    return () => document.removeEventListener("keydown", onKey);
  }, [open, onClose]);

  const handleSave = async () => {
    if (!roleId || roleId === user.role?.id) {
      setEditing(false);
      return;
    }
    setSaving(true);
    setSaveError(null);
    try {
      const res = await fetch(
        `/api/v1/users/${user.id}/role?tenant_id=${DEV_TENANT}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ role_id: roleId }),
        },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      setEditing(false);
      onUpdated();
    } catch (err) {
      setSaveError(err instanceof Error ? err.message : "Save failed");
    } finally {
      setSaving(false);
    }
  };

  const handleDeactivate = async () => {
    if (!confirm("Deactivate this user? They will lose access.")) return;
    setDeactivating(true);
    try {
      const res = await fetch(
        `/api/v1/users/${user.id}/deactivate?tenant_id=${DEV_TENANT}`,
        { method: "POST" },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      onUpdated();
      onClose();
    } catch {
      setSaveError("Deactivation failed");
    } finally {
      setDeactivating(false);
    }
  };

  const copyEmail = () => {
    navigator.clipboard.writeText(user.email);
    setCopied(true);
    setTimeout(() => setCopied(false), 1500);
  };

  const cancelEdit = () => {
    setEditing(false);
    setFirstName(user.first_name);
    setLastName(user.last_name);
    setRoleId(user.role?.id ?? "");
    setSaveError(null);
  };

  const currentRole = roles.find((r) => r.id === roleId);
  const created = new Date(user.created_at).toLocaleDateString("en-IN", {
    day: "2-digit", month: "short", year: "numeric",
  });

  return (
    <>
      <div
        className={cn(
          "fixed inset-0 z-40 transition-opacity duration-200",
          open ? "bg-black/40 pointer-events-auto" : "bg-transparent pointer-events-none",
        )}
        onClick={onClose}
        aria-hidden="true"
      />

      <div
        ref={panelRef}
        data-testid="user-detail-panel"
        className={cn(
          "fixed top-0 right-0 z-50 h-screen w-[480px] max-w-[90vw] bg-[var(--bg-primary)]",
          "border-l border-[var(--border-default)] shadow-xl flex flex-col",
          "transition-transform duration-200 ease-out",
          open ? "translate-x-0" : "translate-x-full",
        )}
      >
        <div className="flex items-start justify-between p-5 border-b border-[var(--border-default)]">
          <div>
            <h3 className="text-heading-lg text-[var(--text-primary)]">
              {user.first_name} {user.last_name}
            </h3>
            <p className="text-body-sm text-[var(--text-muted)] mt-0.5">User details</p>
          </div>
          <button
            onClick={onClose}
            className="p-1.5 rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-colors"
            aria-label="Close panel"
          >
            <X className="w-5 h-5" />
          </button>
        </div>

        <div className="flex-1 overflow-y-auto p-5 space-y-5">
          {editing ? (
            <>
              <Field label="First name">
                <input value={firstName} onChange={(e) => setFirstName(e.target.value)} className={INPUT_CLS} />
              </Field>
              <Field label="Last name">
                <input value={lastName} onChange={(e) => setLastName(e.target.value)} className={INPUT_CLS} />
              </Field>
              <Field label="Email"><p className="text-xs text-[var(--text-muted)]">{user.email}</p></Field>
              <Field label="Role">
                <select value={roleId} onChange={(e) => setRoleId(e.target.value)} className={INPUT_CLS}>
                  <option value="">Select role</option>
                  {roles.map((r) => <option key={r.id} value={r.id}>{r.display_name}</option>)}
                </select>
              </Field>
              {saveError && <p className="text-xs text-[var(--danger)]">{saveError}</p>}
              <div className="flex gap-2 pt-2">
                <button onClick={handleSave} disabled={saving} className="px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90 disabled:opacity-50">
                  {saving ? "Saving..." : "Save Changes"}
                </button>
                <button onClick={cancelEdit} className="px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]">
                  Cancel
                </button>
              </div>
            </>
          ) : (
            <>
              <Field label="Email">
                <div className="flex items-center gap-2">
                  <span className="text-xs text-[var(--text-primary)]">{user.email}</span>
                  <button onClick={copyEmail} className="text-[var(--text-muted)] hover:text-[var(--text-primary)]" aria-label="Copy email">
                    {copied ? <Check className="w-3.5 h-3.5 text-[var(--success)]" /> : <Copy className="w-3.5 h-3.5" />}
                  </button>
                </div>
              </Field>
              <Field label="Role">
                {currentRole ? (
                  <div>
                    <span className="inline-block text-[10px] font-medium px-2 py-0.5 rounded-full bg-[var(--accent)]/10 text-[var(--accent)]">
                      {currentRole.display_name}
                    </span>
                    {currentRole.description && (
                      <p className="text-[10px] text-[var(--text-muted)] mt-1">{currentRole.description}</p>
                    )}
                  </div>
                ) : (
                  <span className="text-xs text-[var(--text-muted)]">No role assigned</span>
                )}
              </Field>
              <Field label="Status">
                <UserStatusPill status={user.status} />
              </Field>
              <Field label="Created">
                <span className="text-xs text-[var(--text-secondary)]">{created}</span>
              </Field>
              <button
                onClick={() => setEditing(true)}
                className="px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
              >
                Edit
              </button>
            </>
          )}
        </div>

        {user.status === "active" && !editing && (
          <div className="border-t border-[var(--border-default)] px-5 py-3">
            <button
              onClick={handleDeactivate}
              disabled={deactivating}
              className="text-xs font-medium text-[var(--danger)] hover:underline disabled:opacity-50"
            >
              {deactivating ? "Deactivating..." : "Deactivate user"}
            </button>
          </div>
        )}
      </div>
    </>
  );
}

function Field({ label, children }: { label: string; children: React.ReactNode }) {
  return (
    <div>
      <dt className="text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">{label}</dt>
      <dd>{children}</dd>
    </div>
  );
}
