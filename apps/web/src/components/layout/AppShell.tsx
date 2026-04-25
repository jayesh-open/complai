"use client";

import { useState } from "react";
import { Sidebar } from "./Sidebar";
import { Header } from "./Header";
import { CommandPalette } from "./CommandPalette";
import { useAppStore } from "@/store/app-store";

export function AppShell({ children }: { children: React.ReactNode }) {
  const [cmdOpen, setCmdOpen] = useState(false);
  const toggleSidebar = useAppStore((s) => s.toggleSidebar);

  return (
    <div className="flex h-screen overflow-hidden bg-app-bg">
      <div className="hidden md:flex">
        <Sidebar />
      </div>
      <div className="flex flex-col flex-1 min-w-0">
        <Header
          onCommandPaletteOpen={() => setCmdOpen(true)}
          onMobileMenuToggle={toggleSidebar}
        />
        <main className="flex-1 overflow-y-auto p-7" data-testid="main-content">
          {children}
        </main>
      </div>
      <CommandPalette open={cmdOpen} onOpenChange={setCmdOpen} />
    </div>
  );
}
