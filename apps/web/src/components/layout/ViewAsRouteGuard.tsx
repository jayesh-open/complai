"use client";

import { useEffect } from "react";
import { usePathname, useRouter } from "next/navigation";
import { useViewAsStore } from "@/store/view-as-store";
import { canAccessRoute } from "@/lib/permissions";

export function ViewAsRouteGuard({ children }: { children: React.ReactNode }) {
  const pathname = usePathname();
  const router = useRouter();
  const { isViewAs, viewAsPermissions } = useViewAsStore();

  useEffect(() => {
    if (!isViewAs) return;
    if (!canAccessRoute(viewAsPermissions, pathname)) {
      router.replace("/dashboard");
    }
  }, [isViewAs, viewAsPermissions, pathname, router]);

  return <>{children}</>;
}
