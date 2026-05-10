import type { Permission } from "@/app/configure/users/types";

const ROUTE_RESOURCE_MAP: Record<string, { resource: string; action: string }> = {
  "/compliance/calendar": { resource: "compliance_calendar", action: "view" },
  "/compliance/gst": { resource: "gst_returns", action: "view" },
  "/compliance/gst-returns/annual": { resource: "gstr_9_9c", action: "view" },
  "/compliance/e-invoicing": { resource: "e_invoicing", action: "view" },
  "/compliance/e-way-bill": { resource: "e_way_bill", action: "view" },
  "/compliance/itc-reconciliation": { resource: "itc_reconciliation", action: "view" },
  "/compliance/vendor-compliance": { resource: "vendor_compliance", action: "view" },
  "/compliance/tds": { resource: "tds", action: "view" },
  "/compliance/itr": { resource: "itr", action: "view" },
  "/configure/users": { resource: "users_roles", action: "view" },
  "/configure/connected-apps": { resource: "connected_apps", action: "view" },
};

const ALWAYS_ACCESSIBLE = ["/dashboard", "/configure/appearance"];

export function hasPermission(
  permissions: Permission[],
  resource: string,
  action: string,
): boolean {
  return permissions.some((p) => p.resource === resource && p.action === action);
}

export function canAccessRoute(
  permissions: Permission[],
  route: string,
): boolean {
  if (ALWAYS_ACCESSIBLE.some((r) => route === r || route.startsWith(r + "/"))) {
    return true;
  }

  const sortedRoutes = Object.keys(ROUTE_RESOURCE_MAP).sort(
    (a, b) => b.length - a.length,
  );
  for (const mappedRoute of sortedRoutes) {
    if (route === mappedRoute || route.startsWith(mappedRoute + "/")) {
      const { resource, action } = ROUTE_RESOURCE_MAP[mappedRoute];
      return hasPermission(permissions, resource, action);
    }
  }

  return true;
}

export function canAccessSidebarItem(
  permissions: Permission[],
  href: string,
): boolean {
  return canAccessRoute(permissions, href);
}

export { ROUTE_RESOURCE_MAP, ALWAYS_ACCESSIBLE };
