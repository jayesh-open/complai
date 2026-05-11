import { test, expect } from "@playwright/test";

test.describe("Users & Roles", () => {
  test("users list page loads with table and add-user modal", async ({ page }) => {
    const t0 = Date.now();
    await page.goto("/configure/users");
    await expect(page.getByTestId("users-page")).toBeVisible({ timeout: 10000 });
    console.log(`[perf] Users list load: ${Date.now() - t0}ms`);

    await expect(page.getByTestId("page-title")).toHaveText("Users & Roles");

    const addBtn = page.getByRole("button", { name: /add user/i });
    await expect(addBtn).toBeVisible();
    await addBtn.click();
    await expect(page.getByTestId("add-user-modal")).toBeVisible();

    const cancelBtn = page.getByTestId("add-user-modal").getByRole("button", { name: /cancel/i });
    await cancelBtn.click();
    await expect(page.getByTestId("add-user-modal")).not.toBeVisible();
  });

  test("roles list page shows system roles with badges", async ({ page }) => {
    await page.goto("/configure/users/roles");
    await expect(page.getByTestId("roles-page")).toBeVisible({ timeout: 10000 });

    const systemBadges = page.getByText("System", { exact: true });
    await expect(systemBadges.first()).toBeVisible({ timeout: 10000 });
    expect(await systemBadges.count()).toBeGreaterThanOrEqual(7);

    const adminRow = page.getByRole("cell", { name: "Admin" });
    await expect(adminRow).toBeVisible();
    await adminRow.click();

    await expect(page.getByTestId("role-detail-page")).toBeVisible({ timeout: 10000 });
  });

  test("role detail page shows permission matrix for system role", async ({ page }) => {
    const t0 = Date.now();
    await page.goto("/configure/users/roles");
    await expect(page.getByTestId("roles-page")).toBeVisible({ timeout: 10000 });
    await page.getByRole("cell", { name: "Admin" }).click();

    await expect(page.getByTestId("role-detail-page")).toBeVisible({ timeout: 15000 });
    console.log(`[perf] Role detail load: ${Date.now() - t0}ms`);

    await expect(page.getByTestId("system-role-banner")).toBeVisible();
    await expect(page.getByTestId("permission-matrix")).toBeVisible();

    const moduleRows = page.locator("[data-testid^='module-row-']");
    expect(await moduleRows.count()).toBe(12);

    const saveBtn = page.getByRole("button", { name: /save changes/i });
    await expect(saveBtn).toBeDisabled();
  });

  test("view-as entry: select Auditor and verify banner + sidebar filtering", async ({ page }) => {
    await page.goto("/dashboard");
    await expect(page.getByTestId("main-content")).toBeVisible({ timeout: 10000 });

    const t0 = Date.now();
    const trigger = page.getByTestId("view-as-trigger");
    await expect(trigger).toBeVisible();
    await trigger.click();

    const auditorBtn = page.getByTestId("view-as-role-auditor");
    await expect(auditorBtn).toBeVisible({ timeout: 5000 });
    await auditorBtn.click();

    await expect(page.getByTestId("view-as-banner")).toBeVisible({ timeout: 10000 });
    console.log(`[perf] View As switch: ${Date.now() - t0}ms`);

    await expect(page.getByTestId("view-as-banner")).toContainText("Auditor");

    const sidebar = page.getByTestId("sidebar-nav");
    await expect(sidebar.getByText("Users & Roles")).not.toBeVisible();
    await expect(sidebar.getByText("Dashboard")).toBeVisible();
  });

  test("view-as route guard redirects to dashboard", async ({ page }) => {
    await page.goto("/dashboard");
    await expect(page.getByTestId("main-content")).toBeVisible({ timeout: 10000 });

    const trigger = page.getByTestId("view-as-trigger");
    await trigger.click();
    const auditorBtn = page.getByTestId("view-as-role-auditor");
    await expect(auditorBtn).toBeVisible({ timeout: 5000 });
    await auditorBtn.click();
    await expect(page.getByTestId("view-as-banner")).toBeVisible({ timeout: 10000 });

    await page.goto("/configure/users");
    await page.waitForURL("**/dashboard", { timeout: 10000 });
    expect(page.url()).toContain("/dashboard");
  });

  test("exit view-as restores admin view", async ({ page }) => {
    await page.goto("/dashboard");
    await expect(page.getByTestId("main-content")).toBeVisible({ timeout: 10000 });

    const trigger = page.getByTestId("view-as-trigger");
    await trigger.click();
    const auditorBtn = page.getByTestId("view-as-role-auditor");
    await expect(auditorBtn).toBeVisible({ timeout: 5000 });
    await auditorBtn.click();
    await expect(page.getByTestId("view-as-banner")).toBeVisible({ timeout: 10000 });

    const exitBtn = page.getByTestId("exit-view-as");
    await exitBtn.click();

    await expect(page.getByTestId("view-as-banner")).not.toBeVisible();

    const sidebar = page.getByTestId("sidebar-nav");
    await expect(sidebar.getByText("Users & Roles")).toBeVisible();
    await expect(sidebar.getByText("Dashboard")).toBeVisible();
  });
});
