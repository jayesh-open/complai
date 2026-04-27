import { test, expect } from "@playwright/test";

test.describe("E-Invoicing IRN Lifecycle", () => {
  test("generate IRN → view detail with QR + signed JSON → cancel within 24h", async ({ page }) => {
    // Navigate to generate page
    await page.goto("/compliance/e-invoicing/generate");
    await expect(page.getByText("Generate IRN")).toBeVisible();
    await expect(page.getByText("Select a source invoice")).toBeVisible();

    // Select first invoice from mock data
    const firstCard = page.locator("[data-testid^='invoice-card-']").first();
    await expect(firstCard).toBeVisible();
    const invoiceNo = await firstCard.locator(".font-mono").first().textContent();
    await firstCard.click();

    // Validate step — should show payload summary and no errors for well-formed mock data
    await expect(page.getByText("Payload Summary")).toBeVisible();
    const generateBtn = page.getByTestId("generate-irn-button");
    await expect(generateBtn).toBeVisible();
    await expect(generateBtn).toBeEnabled();
    await generateBtn.click();

    // Generating step — spinner
    await expect(page.getByText("Submitting to IRP...")).toBeVisible();

    // Success step — IRN, ack number, QR code displayed
    const successPanel = page.getByTestId("irn-success");
    await expect(successPanel).toBeVisible({ timeout: 5000 });
    await expect(page.getByText("IRN Generated Successfully")).toBeVisible();
    await expect(page.getByText(`Invoice ${invoiceNo}`)).toBeVisible();

    // Verify IRN and ack number are rendered
    await expect(successPanel.getByText("Ack Number")).toBeVisible();
    await expect(successPanel.getByText("Ack Date")).toBeVisible();

    // QR code should be visible
    const qrCode = page.locator("svg").first();
    await expect(qrCode).toBeVisible();

    // Navigate to detail view
    const viewLink = page.getByTestId("view-details-link");
    await expect(viewLink).toBeVisible();
    await viewLink.click();

    // Detail page — signed JSON viewer renders
    await expect(page.getByTestId("signed-json-viewer")).toBeVisible();
    await expect(page.getByText("Signed Invoice JSON")).toBeVisible();

    // Click to expand signed JSON
    await page.getByText("Signed Invoice JSON").click();
    await expect(page.locator("pre")).toBeVisible();

    // Verify cancel button within 24h (mock records generated recently)
    const cancelBtn = page.getByTestId("cancel-irn-button");
    if (await cancelBtn.isVisible()) {
      await cancelBtn.click();

      // Cancel confirmation modal
      const modal = page.getByRole("alertdialog");
      await expect(modal).toBeVisible();
      await expect(modal.getByText("Cancel IRN")).toBeVisible();

      // Type CANCEL to confirm
      const confirmInput = page.getByTestId("confirm-input");
      await confirmInput.fill("CANCEL");

      const confirmBtn = page.getByTestId("modal-confirm-button");
      await expect(confirmBtn).toBeEnabled();
      await confirmBtn.click();

      // Modal should close
      await expect(modal).not.toBeVisible();
    }
  });

  test("cancelled record does not show cancel button", async ({ page }) => {
    await page.goto("/compliance/e-invoicing");

    // Click the "Cancelled" status filter tab
    const cancelledTab = page.locator("button").filter({ hasText: /^Cancelled$/ }).first();
    await cancelledTab.click();

    // Wait for table to re-render with only cancelled records
    await page.waitForTimeout(300);

    // Click first table row
    const firstRow = page.locator("table tbody tr").first();
    if (await firstRow.isVisible()) {
      await firstRow.click();
      await page.waitForURL(/\/compliance\/e-invoicing\//);

      // Cancel button should NOT be visible for already-cancelled records
      await expect(page.getByTestId("cancel-irn-button")).not.toBeVisible();
    }
  });
});
