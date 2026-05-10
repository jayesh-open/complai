"use client";

import { useState, type FormEvent } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";
import type { Role } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";
const EMAIL_RE = /^[^\s@]+@[^\s@]+\.[^\s@]+$/;

interface AddUserModalProps {
  open: boolean;
  roles: Role[];
  onClose: () => void;
  onCreated: () => void;
}

export function AddUserModal({ open, roles, onClose, onCreated }: AddUserModalProps) {
  const [email, setEmail] = useState("");
  const [firstName, setFirstName] = useState("");
  const [lastName, setLastName] = useState("");
  const [roleId, setRoleId] = useState(roles[0]?.id ?? "");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  if (!open) return null;

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!email.trim()) errs.email = "Email is required";
    else if (!EMAIL_RE.test(email.trim())) errs.email = "Invalid email format";
    if (!firstName.trim()) errs.firstName = "First name is required";
    if (!lastName.trim()) errs.lastName = "Last name is required";
    if (!roleId) errs.roleId = "Role is required";
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const reset = () => {
    setEmail("");
    setFirstName("");
    setLastName("");
    setRoleId(roles[0]?.id ?? "");
    setError(null);
    setFieldErrors({});
  };

  const handleClose = () => {
    reset();
    onClose();
  };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setError(null);
    try {
      const res = await fetch(
        `/api/v1/users?tenant_id=${DEV_TENANT}`,
        {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({
            email: email.trim(),
            first_name: firstName.trim(),
            last_name: lastName.trim(),
            role_id: roleId,
          }),
        },
      );
      if (!res.ok) {
        const body = await res.json().catch(() => ({}));
        throw new Error(body.message ?? `HTTP ${res.status}`);
      }
      reset();
      onCreated();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create user");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className="fixed inset-0 z-40 bg-black/40" onClick={handleClose} aria-hidden="true" />
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div
          data-testid="add-user-modal"
          className="w-full max-w-md rounded-2xl bg-[var(--bg-primary)] border border-[var(--border-default)] shadow-xl"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--border-default)]">
            <h3 className="text-heading-lg text-[var(--text-primary)]">Add User</h3>
            <button
              onClick={handleClose}
              className="p-1.5 rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-colors"
              aria-label="Close"
            >
              <X className="w-5 h-5" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="px-5 py-4 space-y-4">
            <FormField label="Email" error={fieldErrors.email}>
              <input
                type="email"
                value={email}
                onChange={(e) => setEmail(e.target.value)}
                placeholder="user@company.com"
                className={inputCls(!!fieldErrors.email)}
              />
            </FormField>

            <div className="grid grid-cols-2 gap-3">
              <FormField label="First name" error={fieldErrors.firstName}>
                <input
                  value={firstName}
                  onChange={(e) => setFirstName(e.target.value)}
                  placeholder="First name"
                  className={inputCls(!!fieldErrors.firstName)}
                />
              </FormField>
              <FormField label="Last name" error={fieldErrors.lastName}>
                <input
                  value={lastName}
                  onChange={(e) => setLastName(e.target.value)}
                  placeholder="Last name"
                  className={inputCls(!!fieldErrors.lastName)}
                />
              </FormField>
            </div>

            <FormField label="Role" error={fieldErrors.roleId}>
              <select
                value={roleId}
                onChange={(e) => setRoleId(e.target.value)}
                className={inputCls(!!fieldErrors.roleId)}
              >
                <option value="">Select role</option>
                {roles.map((r) => (
                  <option key={r.id} value={r.id}>{r.display_name}</option>
                ))}
              </select>
            </FormField>

            {error && (
              <p className="text-xs text-[var(--danger)]">{error}</p>
            )}

            <div className="flex justify-end gap-2 pt-2">
              <button
                type="button"
                onClick={handleClose}
                className="px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]"
              >
                Cancel
              </button>
              <button
                type="submit"
                disabled={submitting}
                className="px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90 disabled:opacity-50"
              >
                {submitting ? "Adding..." : "Add User"}
              </button>
            </div>
          </form>
        </div>
      </div>
    </>
  );
}

function inputCls(hasError: boolean) {
  return cn(
    "w-full rounded-lg border bg-[var(--bg-secondary)] px-3 py-2 text-xs text-[var(--text-primary)] outline-none",
    hasError
      ? "border-[var(--danger)] focus:border-[var(--danger)]"
      : "border-[var(--border-default)] focus:border-[var(--accent)]",
  );
}

function FormField({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="block text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">
        {label}
      </label>
      {children}
      {error && <p className="text-[10px] text-[var(--danger)] mt-0.5">{error}</p>}
    </div>
  );
}
