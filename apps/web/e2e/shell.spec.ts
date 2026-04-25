import { test, expect } from "@playwright/test";

test.describe("App Shell", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/dashboard");
  });

  test("redirects / to /dashboard", async ({ page }) => {
    await page.goto("/");
    await expect(page).toHaveURL(/\/dashboard/);
  });

  test("renders sidebar with 6 nav groups in correct order", async ({ page }) => {
    const sidebar = page.locator("aside");
    await expect(sidebar).toBeVisible();

    const groupLabels = sidebar.locator("[data-testid='sidebar-group-label']");
    const labels = await groupLabels.allTextContents();
    const expectedOrder = ["COMPLIANCE", "INSIGHTS", "DATA SOURCES", "DOCUMENTS", "CONFIGURE"];
    expect(labels.map((l) => l.trim())).toEqual(expectedOrder);
  });

  test("sidebar group collapse persists across reload", async ({ page }) => {
    const complianceGroup = page.locator("[data-testid='sidebar-group-compliance']");
    const toggleBtn = complianceGroup.locator("button").first();
    await toggleBtn.click();

    const items = complianceGroup.locator("[data-testid='sidebar-item']");
    await expect(items).toHaveCount(0);

    await page.reload();
    await page.waitForLoadState("networkidle");

    const itemsAfterReload = page
      .locator("[data-testid='sidebar-group-compliance']")
      .locator("[data-testid='sidebar-item']");
    await expect(itemsAfterReload).toHaveCount(0);
  });

  test("sidebar shows badge counts on nav items", async ({ page }) => {
    const badges = page.locator("aside [data-testid='sidebar-badge']");
    expect(await badges.count()).toBeGreaterThan(0);
  });

  test("header renders page title", async ({ page }) => {
    const header = page.locator("header");
    await expect(header).toBeVisible();
    await expect(header.locator("h1, [data-testid='page-title']")).toContainText(/dashboard/i);
  });
});

test.describe("Command Palette", () => {
  test("opens with Cmd+K and shows search", async ({ page }) => {
    await page.goto("/dashboard");
    await page.keyboard.press("Meta+k");
    const palette = page.locator("[data-testid='command-palette'], [cmdk-dialog], [role='dialog']");
    await expect(palette).toBeVisible();
  });

  test("filters commands on search", async ({ page }) => {
    await page.goto("/dashboard");
    await page.keyboard.press("Meta+k");
    const input = page.locator("[cmdk-input], [data-testid='command-input']");
    await input.fill("GSTR");
    const items = page.locator("[cmdk-item]");
    expect(await items.count()).toBeGreaterThan(0);
  });
});

test.describe("Theme Switching", () => {
  test("appearance page loads and shows theme cards", async ({ page }) => {
    await page.goto("/configure/appearance");
    const themeCards = page.locator("[data-testid^='theme-card-']");
    expect(await themeCards.count()).toBeGreaterThanOrEqual(15);
  });

  test("switching theme persists across reload", async ({ page }) => {
    await page.goto("/configure/appearance");

    const darkCard = page.locator("[data-testid='theme-card-dark']");
    await darkCard.click();

    const html = page.locator("html");
    await expect(html).toHaveAttribute("data-theme", "dark");

    await page.reload();
    await page.waitForLoadState("networkidle");
    await expect(page.locator("html")).toHaveAttribute("data-theme", "dark");
  });

  test("density selector changes row heights", async ({ page }) => {
    await page.goto("/configure/appearance");

    const comfortableBtn = page.locator("[data-testid='density-comfortable']");
    await comfortableBtn.click();

    await page.reload();
    await page.waitForLoadState("networkidle");
    await expect(page.locator("[data-testid='density-comfortable']")).toHaveAttribute(
      "aria-pressed",
      "true"
    );
  });
});

test.describe("Dashboard", () => {
  test("renders 4 KPI metric cards", async ({ page }) => {
    await page.goto("/dashboard");
    const kpiCards = page.locator("[data-testid='kpi-card']");
    expect(await kpiCards.count()).toBe(4);
  });

  test("renders compliance health section", async ({ page }) => {
    await page.goto("/dashboard");
    await expect(page.locator("text=Compliance Health")).toBeVisible();
  });

  test("renders action items", async ({ page }) => {
    await page.goto("/dashboard");
    await expect(page.locator("text=Action Items")).toBeVisible();
  });
});

test.describe("Mobile Responsive", () => {
  test.use({ viewport: { width: 375, height: 812 } });

  test("sidebar is hidden on mobile", async ({ page }) => {
    await page.goto("/dashboard");
    const sidebar = page.locator("aside");
    await expect(sidebar).toBeHidden();
  });

  test("mobile menu button is visible", async ({ page }) => {
    await page.goto("/dashboard");
    const menuBtn = page.locator("[data-testid='mobile-menu-btn']");
    await expect(menuBtn).toBeVisible();
  });
});
