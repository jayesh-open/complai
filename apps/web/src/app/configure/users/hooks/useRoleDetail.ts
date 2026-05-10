"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import type { Role, Permission } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

interface UseRoleDetailResult {
  role: Role | null;
  permissions: Permission[];
  loading: boolean;
  error: string | null;
  mutate: () => void;
}

export function useRoleDetail(roleId: string): UseRoleDetailResult {
  const [role, setRole] = useState<Role | null>(null);
  const [permissions, setPermissions] = useState<Permission[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const controllerRef = useRef<AbortController | null>(null);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);

  const doFetch = useCallback(async () => {
    if (retryRef.current) clearTimeout(retryRef.current);
    controllerRef.current?.abort();
    const controller = new AbortController();
    controllerRef.current = controller;
    setLoading(true);
    try {
      const res = await fetch(
        `/api/v1/roles/${encodeURIComponent(roleId)}?tenant_id=${encodeURIComponent(DEV_TENANT)}`,
        { signal: controller.signal },
      );
      if (!res.ok) throw new Error(`HTTP ${res.status}`);
      const data = await res.json();
      setRole(data.role ?? null);
      setPermissions(Array.isArray(data.permissions) ? data.permissions : []);
      setError(null);
    } catch (err) {
      if (controller.signal.aborted) return;
      setError(err instanceof Error ? err.message : "Failed to fetch role");
      retryRef.current = setTimeout(doFetch, 30_000);
    } finally {
      if (!controller.signal.aborted) setLoading(false);
    }
  }, [roleId]);

  useEffect(() => {
    doFetch();
    return () => {
      controllerRef.current?.abort();
      if (retryRef.current) clearTimeout(retryRef.current);
    };
  }, [doFetch]);

  return { role, permissions, loading, error, mutate: doFetch };
}
