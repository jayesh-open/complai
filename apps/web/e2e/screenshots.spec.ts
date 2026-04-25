import { test } from "@playwright/test";

const STORYBOOK = "http://localhost:6006";
const DIR = "/tmp/storybook-screenshots";

const SHOTS = [
  { id: "compliance-govstatuspill--all-systems", name: "GovStatusPill-AllSystems" },
  { id: "compliance-govstatuspill--success", name: "GovStatusPill-Success" },
  { id: "compliance-govstatuspill--warning", name: "GovStatusPill-Warning" },
  { id: "compliance-govstatuspill--danger", name: "GovStatusPill-Danger" },
  { id: "compliance-filingconfirmationmodal--default", name: "FilingConfirmationModal" },
  { id: "compliance-datatable--compact", name: "DataTable-Compact" },
  { id: "compliance-audittrailtimeline--default", name: "AuditTrailTimeline" },
  { id: "compliance-vendorcompliancescorecard--high-score", name: "VendorComplianceScoreCard-High" },
  { id: "compliance-vendorcompliancescorecard--low-score", name: "VendorComplianceScoreCard-Low" },
];

for (const shot of SHOTS) {
  test(`screenshot: ${shot.name}`, async ({ page }) => {
    await page.goto(`${STORYBOOK}/iframe.html?id=${shot.id}&viewMode=story`, { waitUntil: "networkidle" });
    await page.waitForSelector("#storybook-root", { state: "attached" });
    await page.waitForTimeout(1000);
    await page.screenshot({ path: `${DIR}/${shot.name}.png`, fullPage: true });
  });
}
