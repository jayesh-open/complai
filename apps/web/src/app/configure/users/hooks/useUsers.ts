"use client";

import { useState, useEffect, useRef, useCallback } from "react";
import type { User } from "../types";

const DEV_TENANT = "00000000-0000-0000-0000-000000000001";

interface UseUsersResult {
  users: User[];
  loading: boolean;
  error: string | null;
  mutate: () => void;
}

export function useUsers(tenantId = DEV_TENANT): UseUsersResult {
  const [users, setUsers] = useState<User[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const retryRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const versionRef = useRef(0);

  const doFetch = useCallback(
    (controller: AbortController) => {
      const run = async () => {
        try {
          const res = await fetch(
            `/api/v1/users?tenant_id=${encodeURIComponent(tenantId)}`,
            { signal: controller.signal },
          );
          if (!res.ok) throw new Error(`HTTP ${res.status}`);
          const data = await res.json();
          setUsers(data.users ?? []);
          setError(null);
        } catch (err) {
          if (controller.signal.aborted) return;
          setError(err instanceof Error ? err.message : "Failed to fetch users");
          retryRef.current = setTimeout(() => {
            if (!controller.signal.aborted) run();
          }, 30_000);
        } finally {
          if (!controller.signal.aborted) setLoading(false);
        }
      };
      setLoading(true);
      run();
    },
    [tenantId],
  );

  useEffect(() => {
    if (retryRef.current) clearTimeout(retryRef.current);
    const controller = new AbortController();
    doFetch(controller);
    return () => {
      controller.abort();
      if (retryRef.current) clearTimeout(retryRef.current);
    };
  }, [doFetch, versionRef]); // eslint-disable-line react-hooks/exhaustive-deps

  const mutate = useCallback(() => {
    versionRef.current += 1;
    if (retryRef.current) clearTimeout(retryRef.current);
    const controller = new AbortController();
    doFetch(controller);
  }, [doFetch]);

  return { users, loading, error, mutate };
}
