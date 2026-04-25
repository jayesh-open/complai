import { test, expect } from "@playwright/test";
import AxeBuilder from "@axe-core/playwright";

const STORYBOOK_URL = "http://localhost:6006";

const STORIES = [
  { id: "compliance-statusbadge--all-variants", name: "StatusBadge" },
  { id: "compliance-govstatuspill--all-systems", name: "GovStatusPill" },
  { id: "compliance-kpimetriccard--card-grid", name: "KpiMetricCard" },
  { id: "compliance-filingconfirmationmodal--default", name: "FilingConfirmationModal" },
  { id: "compliance-datatable--compact", name: "DataTable" },
  { id: "compliance-periodselector--default", name: "PeriodSelector" },
  { id: "compliance-audittrailtimeline--default", name: "AuditTrailTimeline" },
  { id: "compliance-vendorcompliancescorecard--high-score", name: "VendorComplianceScoreCard" },
  { id: "compliance-bulkoperationtray--default", name: "BulkOperationTray" },
  { id: "compliance-makercheckerapprovalcard--default", name: "MakerCheckerApprovalCard" },
  { id: "compliance-reconciliationsplitpane--gstr-1-vs-2-b", name: "ReconciliationSplitPane" },
];

test.describe("Storybook a11y (axe-core)", () => {
  for (const story of STORIES) {
    test(`${story.name} has no a11y violations`, async ({ page }) => {
      await page.goto(
        `${STORYBOOK_URL}/iframe.html?id=${story.id}&viewMode=story`,
        { waitUntil: "networkidle" }
      );
      await page.waitForSelector("#storybook-root", { state: "attached" });
      await page.waitForTimeout(1000);

      const results = await new AxeBuilder({ page })
        .include("#storybook-root")
        .analyze();

      expect(results.violations).toEqual([]);
    });
  }
});
