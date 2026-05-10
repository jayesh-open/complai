"use client";

import { create } from "zustand";
import { persist } from "zustand/middleware";
import type { Role, Permission } from "@/app/configure/users/types";

interface ViewAsState {
  viewAsRoleId: string | null;
  viewAsRole: Role | null;
  viewAsPermissions: Permission[];
  isViewAs: boolean;
  enterViewAs: (role: Role, permissions: Permission[]) => void;
  exitViewAs: () => void;
}

export const useViewAsStore = create<ViewAsState>()(
  persist(
    (set) => ({
      viewAsRoleId: null,
      viewAsRole: null,
      viewAsPermissions: [],
      isViewAs: false,
      enterViewAs: (role, permissions) =>
        set({
          viewAsRoleId: role.id,
          viewAsRole: role,
          viewAsPermissions: permissions,
          isViewAs: true,
        }),
      exitViewAs: () =>
        set({
          viewAsRoleId: null,
          viewAsRole: null,
          viewAsPermissions: [],
          isViewAs: false,
        }),
    }),
    { name: "complai-view-as" },
  ),
);
