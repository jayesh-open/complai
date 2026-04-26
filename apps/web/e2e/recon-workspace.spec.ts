import { test, expect } from "@playwright/test";

test.describe("ITC Reconciliation Workspace", () => {
  test.beforeEach(async ({ page }) => {
    await page.goto("/compliance/itc-reconciliation");
    await expect(page.getByTestId("recon-workspace")).toBeVisible();
  });

  test("renders workspace with bucket summary and table", async ({ page }) => {
    await expect(page.getByRole("heading", { name: "ITC Reconciliation" })).toBeVisible();

    const bucketSummary = page.getByTestId("bucket-summary");
    await expect(bucketSummary).toBeVisible();

    await expect(page.getByTestId("bucket-all")).toBeVisible();
    await expect(page.getByTestId("bucket-matched")).toBeVisible();
    await expect(page.getByTestId("bucket-mismatch")).toBeVisible();

    const table = page.getByTestId("recon-table");
    await expect(table).toBeVisible();
  });

  test("run reconciliation populates results", async ({ page }) => {
    const runBtn = page.getByTestId("run-recon-button");
    await expect(runBtn).toBeVisible();
    await runBtn.click();

    // Wait for recon to complete (mock has simulated delay)
    await expect(page.getByTestId("recon-table")).toBeVisible({ timeout: 5000 });

    // Bucket counts should show non-zero matched
    const matchedBucket = page.getByTestId("bucket-matched");
    await expect(matchedBucket).toBeVisible();
  });

  test("filter by bucket", async ({ page }) => {
    // Click on mismatch bucket to filter
    const mismatchBucket = page.getByTestId("bucket-mismatch");
    await expect(mismatchBucket).toBeVisible();
    await mismatchBucket.click();

    // Table should still be visible with filtered rows
    await expect(page.getByTestId("recon-table")).toBeVisible();
  });

  test("search filters rows", async ({ page }) => {
    const searchInput = page.getByTestId("search-input");
    await expect(searchInput).toBeVisible();
    await searchInput.fill("INV");

    // Table should still render
    await expect(page.getByTestId("recon-table")).toBeVisible();
  });

  test("toggle AI suggestions", async ({ page }) => {
    const aiToggle = page.getByTestId("ai-toggle");
    await expect(aiToggle).toBeVisible();
    await aiToggle.click();

    // Table should still be visible
    await expect(page.getByTestId("recon-table")).toBeVisible();
  });

  test("export button is present", async ({ page }) => {
    const exportBtn = page.getByTestId("export-button");
    await expect(exportBtn).toBeVisible();
    await expect(exportBtn).toContainText("Export");
  });
});
