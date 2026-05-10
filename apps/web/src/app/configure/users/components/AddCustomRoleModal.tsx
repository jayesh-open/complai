"use client";

import { useState, type FormEvent } from "react";
import { X } from "lucide-react";
import { cn } from "@/lib/utils";
import type { Role } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";
const SLUG_RE = /^[a-z][a-z0-9_]*$/;

interface AddCustomRoleModalProps {
  open: boolean;
  roles: Role[];
  onClose: () => void;
  onCreated: () => void;
}

export function AddCustomRoleModal({ open, roles, onClose, onCreated }: AddCustomRoleModalProps) {
  const [name, setName] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [description, setDescription] = useState("");
  const [startingPoint, setStartingPoint] = useState<"blank" | "clone">("blank");
  const [cloneFrom, setCloneFrom] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [fieldErrors, setFieldErrors] = useState<Record<string, string>>({});

  if (!open) return null;

  const systemRoles = roles.filter((r) => r.is_system);

  const validate = (): boolean => {
    const errs: Record<string, string> = {};
    if (!name.trim()) errs.name = "Name is required";
    else if (!SLUG_RE.test(name.trim())) errs.name = "Lowercase letters, digits, underscores only. Must start with a letter.";
    if (!displayName.trim()) errs.displayName = "Display name is required";
    if (startingPoint === "clone" && !cloneFrom) errs.cloneFrom = "Select a role to clone from";
    setFieldErrors(errs);
    return Object.keys(errs).length === 0;
  };

  const reset = () => {
    setName("");
    setDisplayName("");
    setDescription("");
    setStartingPoint("blank");
    setCloneFrom("");
    setError(null);
    setFieldErrors({});
  };

  const handleClose = () => { reset(); onClose(); };

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    if (!validate()) return;
    setSubmitting(true);
    setError(null);
    try {
      const body: Record<string, unknown> = {
        name: name.trim(),
        display_name: displayName.trim(),
      };
      if (description.trim()) body.description = description.trim();
      if (startingPoint === "clone") body.template = cloneFrom;

      const res = await fetch(`/api/v1/roles?tenant_id=${DEV_TENANT}`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.message ?? `HTTP ${res.status}`);
      }
      reset();
      onCreated();
      onClose();
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to create role");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <>
      <div className="fixed inset-0 z-40 bg-black/40" onClick={handleClose} aria-hidden="true" />
      <div className="fixed inset-0 z-50 flex items-center justify-center p-4">
        <div
          data-testid="add-role-modal"
          className="w-full max-w-md rounded-2xl bg-[var(--bg-primary)] border border-[var(--border-default)] shadow-xl"
          onClick={(e) => e.stopPropagation()}
        >
          <div className="flex items-center justify-between px-5 py-4 border-b border-[var(--border-default)]">
            <h3 className="text-heading-lg text-[var(--text-primary)]">Add Custom Role</h3>
            <button onClick={handleClose} className="p-1.5 rounded-lg text-[var(--text-muted)] hover:bg-[var(--bg-tertiary)] hover:text-[var(--text-primary)] transition-colors" aria-label="Close">
              <X className="w-5 h-5" />
            </button>
          </div>

          <form onSubmit={handleSubmit} className="px-5 py-4 space-y-4">
            <FormField label="Name (slug)" error={fieldErrors.name}>
              <input value={name} onChange={(e) => setName(e.target.value)} placeholder="tax_specialist" className={inputCls(!!fieldErrors.name)} />
            </FormField>
            <FormField label="Display name" error={fieldErrors.displayName}>
              <input value={displayName} onChange={(e) => setDisplayName(e.target.value)} placeholder="Tax Specialist" className={inputCls(!!fieldErrors.displayName)} />
            </FormField>
            <FormField label="Description (optional)">
              <input value={description} onChange={(e) => setDescription(e.target.value)} placeholder="Describe this role" className={inputCls(false)} />
            </FormField>

            <div>
              <span className="block text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-2">Starting point</span>
              <div className="space-y-2">
                <label className="flex items-center gap-2 cursor-pointer">
                  <input type="radio" checked={startingPoint === "blank"} onChange={() => setStartingPoint("blank")} className="accent-[var(--accent)]" />
                  <span className="text-xs text-[var(--text-primary)]">Start blank (no permissions)</span>
                </label>
                <label className="flex items-center gap-2 cursor-pointer">
                  <input type="radio" checked={startingPoint === "clone"} onChange={() => setStartingPoint("clone")} className="accent-[var(--accent)]" />
                  <span className="text-xs text-[var(--text-primary)]">Clone from existing role</span>
                </label>
              </div>
              {startingPoint === "clone" && (
                <div className="mt-2">
                  <select value={cloneFrom} onChange={(e) => setCloneFrom(e.target.value)} className={inputCls(!!fieldErrors.cloneFrom)}>
                    <option value="">Select role to clone</option>
                    {systemRoles.map((r) => <option key={r.id} value={r.name}>{r.display_name}</option>)}
                  </select>
                  {fieldErrors.cloneFrom && <p className="text-[10px] text-[var(--danger)] mt-0.5">{fieldErrors.cloneFrom}</p>}
                </div>
              )}
            </div>

            {error && <p className="text-xs text-[var(--danger)]">{error}</p>}

            <div className="flex justify-end gap-2 pt-2">
              <button type="button" onClick={handleClose} className="px-4 py-2 rounded-lg text-xs font-medium border border-[var(--border-default)] text-[var(--text-primary)] hover:bg-[var(--bg-tertiary)]">
                Cancel
              </button>
              <button type="submit" disabled={submitting} className="px-4 py-2 rounded-lg text-xs font-medium bg-[var(--accent)] text-[var(--accent-text)] hover:opacity-90 disabled:opacity-50">
                {submitting ? "Creating..." : "Create Role"}
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
    hasError ? "border-[var(--danger)] focus:border-[var(--danger)]" : "border-[var(--border-default)] focus:border-[var(--accent)]",
  );
}

function FormField({ label, error, children }: { label: string; error?: string; children: React.ReactNode }) {
  return (
    <div>
      <label className="block text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] mb-1">{label}</label>
      {children}
      {error && <p className="text-[10px] text-[var(--danger)] mt-0.5">{error}</p>}
    </div>
  );
}
