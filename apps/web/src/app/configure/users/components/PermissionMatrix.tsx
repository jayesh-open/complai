"use client";

import { useState, useEffect, useCallback } from "react";
import { PermissionCell } from "./PermissionCell";
import type { Permission } from "../types";

const MODULE_ACTIONS: Record<string, { label: string; actions: string[] }> = {
  gst_returns: { label: "GST Returns", actions: ["view", "edit", "file"] },
  gstr_9_9c: { label: "GSTR-9 / 9C", actions: ["view", "edit", "file"] },
  e_invoicing: { label: "e-Invoicing", actions: ["view", "generate", "cancel"] },
  e_way_bill: { label: "e-Way Bill", actions: ["view", "generate", "cancel"] },
  itc_reconciliation: { label: "ITC Reconciliation", actions: ["view", "edit"] },
  vendor_compliance: { label: "Vendor Compliance", actions: ["view", "edit"] },
  tds: { label: "TDS / TCS", actions: ["view", "calculate", "file", "issue_cert"] },
  itr: { label: "ITR", actions: ["view", "calculate", "file", "approve"] },
  compliance_calendar: { label: "Compliance Calendar", actions: ["view"] },
  users_roles: { label: "Users & Roles", actions: ["view", "manage"] },
  connected_apps: { label: "Connected Apps", actions: ["view", "manage"] },
  billing: { label: "Billing", actions: ["view", "manage"] },
};

const MODULE_ORDER = Object.keys(MODULE_ACTIONS);

export type PermissionPair = { resource: string; action: string };

interface PermissionMatrixProps {
  permissions: Permission[];
  isSystemRole: boolean;
  onChange: (pairs: PermissionPair[]) => void;
  disabled?: boolean;
}

function buildGrantMap(permissions: Permission[]): Map<string, Set<string>> {
  const map = new Map<string, Set<string>>();
  for (const p of permissions) {
    if (!map.has(p.resource)) map.set(p.resource, new Set());
    map.get(p.resource)!.add(p.action);
  }
  return map;
}

function grantMapToPairs(map: Map<string, Set<string>>): PermissionPair[] {
  const pairs: PermissionPair[] = [];
  for (const [resource, actions] of map) {
    for (const action of actions) {
      pairs.push({ resource, action });
    }
  }
  return pairs;
}

export function PermissionMatrix({ permissions, isSystemRole, onChange, disabled }: PermissionMatrixProps) {
  const [grantMap, setGrantMap] = useState(() => buildGrantMap(permissions));

  useEffect(() => {
    setGrantMap(buildGrantMap(permissions));
  }, [permissions]);

  const handleToggle = useCallback(
    (module: string, action: string, granted: boolean) => {
      setGrantMap((prev) => {
        const next = new Map(prev);
        const actions = new Set(next.get(module) ?? []);
        if (granted) actions.add(action);
        else actions.delete(action);
        if (actions.size > 0) next.set(module, actions);
        else next.delete(module);
        onChange(grantMapToPairs(next));
        return next;
      });
    },
    [onChange],
  );

  const isDisabled = isSystemRole || !!disabled;

  return (
    <div className="bg-[var(--bg-secondary)] rounded-[14px] border border-[var(--border-default)] overflow-hidden">
      <table className="w-full">
        <thead>
          <tr className="border-b border-[var(--border-default)]">
            <th className="px-4 py-2.5 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left w-[200px]">
              Module
            </th>
            <th className="px-4 py-2.5 text-[10px] font-semibold uppercase tracking-wide text-[var(--text-muted)] text-left">
              Permissions
            </th>
          </tr>
        </thead>
        <tbody>
          {MODULE_ORDER.map((mod) => {
            const { label, actions } = MODULE_ACTIONS[mod];
            const granted = Array.from(grantMap.get(mod) ?? []);
            return (
              <tr
                key={mod}
                className="border-b border-[var(--border-default)] last:border-b-0 hover:bg-[var(--bg-tertiary)] transition-colors"
              >
                <td className="px-4 py-2.5">
                  <span className="block text-xs font-medium text-[var(--text-primary)]">{label}</span>
                  <span className="block text-[10px] font-mono text-[var(--text-muted)]">{mod}</span>
                </td>
                <td className="px-4 py-2.5">
                  <PermissionCell
                    module={mod}
                    actions={actions}
                    grantedActions={granted}
                    isSystemRole={isDisabled}
                    onChange={(action, g) => handleToggle(mod, action, g)}
                  />
                </td>
              </tr>
            );
          })}
        </tbody>
      </table>
    </div>
  );
}
