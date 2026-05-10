"use client";

import { useState } from "react";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

interface PermissionPair {
  resource: string;
  action: string;
}

interface UseUpdateRolePermissionsResult {
  updatePermissions: (roleId: string, pairs: PermissionPair[]) => Promise<boolean>;
  isUpdating: boolean;
  error: string | null;
}

export function useUpdateRolePermissions(): UseUpdateRolePermissionsResult {
  const [isUpdating, setIsUpdating] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const updatePermissions = async (roleId: string, pairs: PermissionPair[]): Promise<boolean> => {
    setIsUpdating(true);
    setError(null);
    try {
      const res = await fetch(
        `/api/v1/roles/${encodeURIComponent(roleId)}/permissions?tenant_id=${encodeURIComponent(DEV_TENANT)}`,
        {
          method: "PUT",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify({ permission_pairs: pairs }),
        },
      );
      if (!res.ok) {
        const data = await res.json().catch(() => ({}));
        throw new Error(data.message ?? `HTTP ${res.status}`);
      }
      return true;
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to update permissions");
      return false;
    } finally {
      setIsUpdating(false);
    }
  };

  return { updatePermissions, isUpdating, error };
}
