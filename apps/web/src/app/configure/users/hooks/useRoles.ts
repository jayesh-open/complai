"use client";

import { useState, useEffect, useRef } from "react";
import type { Role } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

interface UseRolesResult {
  roles: Role[];
  loading: boolean;
  error: string | null;
}

export function useRoles(tenantId = DEV_TENANT): UseRolesResult {
  const [roles, setRoles] = useState<Role[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  useEffect(() => {
    if (retryRef.current) clearTimeout(retryRef.current);
    const controller = new AbortController();

    const doFetch = async () => {
      try {
        const res = await fetch(
          `/api/v1/roles?tenant_id=${encodeURIComponent(tenantId)}`,
          { signal: controller.signal },
        );
        if (!res.ok) throw new Error(`HTTP ${res.status}`);
        const data = await res.json();
        setRoles(Array.isArray(data) ? data : []);
        setError(null);
      } catch (err) {
        if (controller.signal.aborted) return;
        setError(err instanceof Error ? err.message : "Failed to fetch roles");
        retryRef.current = setTimeout(doFetch, 30_000);
      } finally {
        if (!controller.signal.aborted) setLoading(false);
      }
    };

    setLoading(true);
    doFetch();

    return () => {
      controller.abort();
      if (retryRef.current) clearTimeout(retryRef.current);
    };
  }, [tenantId]);

  return { roles, loading, error };
}
