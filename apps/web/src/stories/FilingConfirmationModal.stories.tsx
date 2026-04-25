import type { Meta, StoryObj } from "@storybook/react";
import { useState } from "react";

function FilingConfirmationModal({ open, onClose, onConfirm, title, details, warningText = "This action is irreversible.", confirmWord = "FILE", signMethod, onSignMethodChange }: {
  open: boolean; onClose: () => void; onConfirm: () => void; title: string;
  details: { label: string; value: string }[]; warningText?: string; confirmWord?: string;
  signMethod?: "dsc" | "evc" | null; onSignMethodChange?: (m: "dsc" | "evc") => void;
}) {
  const [typed, setTyped] = useState("");
  const confirmed = typed.toUpperCase() === confirmWord.toUpperCase();
  if (!open) return null;
  return (
    <div className="relative bg-[var(--bg-secondary)] border border-[var(--border-default)] rounded-2xl shadow-[var(--shadow-lg)] w-full max-w-lg" role="alertdialog">
      <div className="px-6 py-4 border-b border-[var(--border-default)] flex items-center gap-2">
        <span className="text-[var(--warning)] text-lg">⚠</span>
        <h2 className="text-sm font-bold text-[var(--text-primary)]">{title}</h2>
      </div>
      <div className="px-6 py-5 space-y-4">
        <div className="space-y-2">
          {details.map((d) => (
            <div key={d.label} className="flex gap-3 text-[13px]">
              <span className="text-[var(--text-muted)] min-w-[100px]">{d.label}:</span>
              <span className="text-[var(--text-primary)] font-medium">{d.value}</span>
            </div>
          ))}
        </div>
        <p className="text-[13px] text-[var(--danger)] font-medium">{warningText}</p>
        {onSignMethodChange && (
          <div className="space-y-2">
            <p className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">Sign with:</p>
            <div className="flex gap-4">
              {(["dsc", "evc"] as const).map((m) => (
                <label key={m} className="flex items-center gap-2 text-[13px] text-[var(--text-secondary)] cursor-pointer">
                  <input type="radio" name="sign" checked={signMethod === m} onChange={() => onSignMethodChange(m)} className="accent-[var(--accent)]" />
                  {m === "dsc" ? "DSC token (USB)" : "EVC OTP"}
                </label>
              ))}
            </div>
          </div>
        )}
        <div className="space-y-1.5">
          <label className="text-[11px] font-semibold uppercase tracking-wide text-[var(--text-muted)]">
            Type &quot;{confirmWord}&quot; to confirm:
          </label>
          <input type="text" value={typed} onChange={(e) => setTyped(e.target.value)}
            className="w-full px-3 py-2 rounded-lg border text-sm font-mono bg-[var(--bg-tertiary)] border-[var(--border-default)] text-[var(--text-primary)] focus:outline-none focus:border-[var(--accent)]"
            placeholder={confirmWord} />
        </div>
      </div>
      <div className="px-6 py-4 border-t border-[var(--border-default)] flex justify-end gap-3">
        <button onClick={onClose} className="px-4 py-2 text-xs font-medium rounded-lg border border-[var(--border-default)] text-[var(--text-secondary)]">Cancel</button>
        <button onClick={() => { if (confirmed) { onConfirm(); setTyped(""); } }} disabled={!confirmed}
          className={`px-4 py-2 text-xs font-semibold rounded-lg ${confirmed ? "bg-[var(--danger)] text-white" : "bg-[var(--danger-muted)] text-[var(--text-disabled)] cursor-not-allowed"}`}>
          Confirm &amp; File
        </button>
      </div>
    </div>
  );
}

const meta: Meta<typeof FilingConfirmationModal> = {
  title: "Compliance/FilingConfirmationModal",
  component: FilingConfirmationModal,
  tags: ["autodocs"],
};
export default meta;
type Story = StoryObj<typeof FilingConfirmationModal>;

export const Default: Story = {
  render: () => {
    const [sign, setSign] = useState<"dsc" | "evc" | null>(null);
    return (
      <FilingConfirmationModal
        open={true} onClose={() => {}} onConfirm={() => alert("Filed!")}
        title="File GSTR-3B for April 2026"
        details={[
          { label: "Period", value: "April 2026 (FY 2026-27)" },
          { label: "GSTIN", value: "29AABCA1234A1Z5 (Karnataka)" },
          { label: "Tax payable", value: "₹12,45,678" },
        ]}
        warningText="This action is irreversible. Once filed, you cannot revise this return."
        signMethod={sign} onSignMethodChange={setSign}
      />
    );
  },
};
