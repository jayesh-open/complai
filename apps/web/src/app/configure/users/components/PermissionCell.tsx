"use client";

import { PermissionIcon } from "./PermissionIcon";

interface PermissionCellProps {
  module: string;
  actions: string[];
  grantedActions: string[];
  isSystemRole: boolean;
  onChange: (action: string, granted: boolean) => void;
}

export function PermissionCell({ actions, grantedActions, isSystemRole, onChange }: PermissionCellProps) {
  return (
    <div className="flex items-center gap-1">
      {actions.map((action) => {
        const granted = grantedActions.includes(action);
        return (
          <PermissionIcon
            key={action}
            action={action}
            granted={granted}
            disabled={isSystemRole}
            onClick={() => onChange(action, !granted)}
          />
        );
      })}
    </div>
  );
}
