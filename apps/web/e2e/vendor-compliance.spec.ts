import { test, expect } from "@playwright/test";

test.describe("Vendor Compliance", () => {
  test("list vendors with scores + category badges → click vendor → VendorComplianceScoreCard renders", async ({ page }) => {
    await page.goto("/compliance/vendor-compliance");

    // Verify we're on the vendor compliance page
    const listView = page.getByTestId("vendor-compliance-list");
    await expect(listView).toBeVisible();

    // Verify KPI summary cards
    const kpiSummary = page.getByTestId("kpi-summary");
    await expect(kpiSummary).toBeVisible();
    await expect(kpiSummary.getByText("Total Vendors")).toBeVisible();
    await expect(kpiSummary.getByText("50")).toBeVisible();
    await expect(kpiSummary.getByText("Avg Score")).toBeVisible();

    // Verify sync button
    const syncBtn = page.getByTestId("sync-button");
    await expect(syncBtn).toBeVisible();
    await expect(syncBtn).toContainText("Sync from Apex");

    // Verify vendor table
    const vendorTable = page.getByTestId("vendor-table");
    await expect(vendorTable).toBeVisible();

    // Verify vendor rows exist (should be 50 in "all" filter)
    const vendorRows = page.getByTestId("vendor-row");
    await expect(vendorRows).toHaveCount(50);

    // Verify first vendor (Tata Steel) is visible with score and category
    await expect(vendorTable.getByText("Tata Steel Ltd")).toBeVisible();
    const categoryBadges = page.getByTestId("category-badge");
    await expect(categoryBadges.first()).toBeVisible();

    // Test category filter — click Cat A
    const filterA = page.getByTestId("filter-A");
    await filterA.click();
    await expect(page.getByTestId("vendor-row")).toHaveCount(10);

    // Test Cat D filter
    const filterD = page.getByTestId("filter-D");
    await filterD.click();
    await expect(page.getByTestId("vendor-row")).toHaveCount(10);
    await expect(vendorTable.getByText("ABC Traders")).toBeVisible();

    // Reset to all
    const filterAll = page.getByTestId("filter-all");
    await filterAll.click();
    await expect(page.getByTestId("vendor-row")).toHaveCount(50);

    // Test search
    const searchInput = page.getByTestId("vendor-search");
    await searchInput.fill("Tata");
    await expect(page.getByTestId("vendor-row")).toHaveCount(1);
    await expect(vendorTable.getByText("Tata Steel Ltd")).toBeVisible();
    await searchInput.clear();
    await expect(page.getByTestId("vendor-row")).toHaveCount(50);

    // Click on a vendor (Tata Steel) to see detail view
    await vendorTable.getByText("Tata Steel Ltd").click();

    // Verify detail view renders
    const detailView = page.getByTestId("vendor-detail-view");
    await expect(detailView).toBeVisible();

    // Verify VendorComplianceScoreCard renders with 10-dot bar
    const scoreCard = page.getByTestId("vendor-scorecard");
    await expect(scoreCard).toBeVisible();

    // Verify vendor name in scorecard
    await expect(scoreCard.getByText("Tata Steel Ltd")).toBeVisible();

    // Verify 5-dimension breakdown is visible
    await expect(scoreCard.getByText("Filing regularity")).toBeVisible();
    await expect(scoreCard.getByText("IRN compliance")).toBeVisible();
    await expect(scoreCard.getByText("Mismatch rate")).toBeVisible();
    await expect(scoreCard.getByText("Payment behaviour")).toBeVisible();
    await expect(scoreCard.getByText("Document hygiene")).toBeVisible();

    // Verify vendor details section
    await expect(detailView.getByText("Vendor Details")).toBeVisible();
    await expect(detailView.getByText("Legal Name")).toBeVisible();

    // Verify category badge on detail page
    const detailBadge = detailView.getByTestId("category-badge");
    await expect(detailBadge).toBeVisible();
    await expect(detailBadge).toContainText("A — Exemplary");

    // Go back to list
    const backBtn = page.getByTestId("back-to-list");
    await backBtn.click();

    // Verify we're back on list view
    await expect(page.getByTestId("vendor-compliance-list")).toBeVisible();

    // Test sync button interaction
    await syncBtn.click();
    await expect(syncBtn).toContainText("Syncing...");
    await expect(syncBtn).toContainText("Sync from Apex", { timeout: 5000 });
  });
});
