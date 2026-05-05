import { test, expect } from "@playwright/test";

test.describe("ITR Bulk Filing Flow", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/compliance/itr");
    await expect(page.getByTestId("itr-landing")).toBeVisible({ timeout: 10000 });
  });

  test("landing page shows employees and batches tabs", async ({ page }) => {
    // Verify heading
    await expect(page.getByText("Income Tax Returns")).toBeVisible();

    // Switch to Bulk Batches tab
    await page.getByRole("button", { name: /Bulk Batches/i }).click();

    // Verify batch table renders
    await expect(page.getByTestId("batch-table")).toBeVisible();
  });

  test("navigates to bulk batch wizard and shows select step", async ({ page }) => {
    // Click "New Batch" link
    await page.getByRole("link", { name: /New Batch/i }).click();

    // Verify wizard opened
    await expect(page.getByTestId("bulk-wizard")).toBeVisible({ timeout: 5000 });
    await expect(page.getByText("Create Bulk Filing Batch")).toBeVisible();

    // Step 1: Select Employees — verify employee rows present
    await expect(page.getByText("Select Employees", { exact: true })).toBeVisible();
    const checkboxes = page.locator("input[type='checkbox']");
    const count = await checkboxes.count();
    expect(count).toBeGreaterThan(0);
  });

  test("navigates to batch detail and employee detail", async ({ page }) => {
    // Switch to Bulk Batches tab
    await page.getByRole("button", { name: /Bulk Batches/i }).click();
    await expect(page.getByTestId("batch-table")).toBeVisible();

    // Click first batch link
    const firstBatchLink = page.getByTestId("batch-table").locator("a").first();
    await firstBatchLink.click();

    // Verify batch detail
    await expect(page.getByTestId("batch-detail")).toBeVisible({ timeout: 5000 });
    await expect(page.getByText("Employees in Batch")).toBeVisible();

    // Click first employee link in the batch
    const employeeLink = page.getByTestId("batch-detail").locator("table tbody tr a").first();
    await employeeLink.click();

    // Verify employee detail view
    await expect(page.getByTestId("employee-detail")).toBeVisible({ timeout: 5000 });
    await expect(page.getByText(/PAN/)).toBeVisible();
    await expect(page.getByText("Income Summary")).toBeVisible();
  });
});
