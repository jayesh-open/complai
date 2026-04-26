import { test, expect } from "@playwright/test";

test.describe("GSTR-3B Filing Wizard", () => {
  test("full wizard lifecycle: Auto-populate → Review → Pay → Sign → File → Acknowledge", async ({ page }) => {
    await page.goto("/compliance/gst-returns/29AABCA1234A1Z5/2026-04/gstr-3b");

    const wizard = page.getByTestId("gstr3b-wizard");
    await expect(wizard).toBeVisible();

    // Step 1: Auto-populate
    const populateBtn = page.getByTestId("populate-button");
    await expect(populateBtn).toBeVisible();
    await expect(populateBtn).toContainText("Auto-Populate");
    await populateBtn.click();

    // Step 2: Review — wait for auto-populate to finish and advance
    const reviewStep = page.getByTestId("step-review");
    await expect(reviewStep).toBeVisible({ timeout: 5000 });

    // Verify tabs are present
    await expect(page.getByTestId("tab-liability")).toBeVisible();
    await expect(page.getByTestId("tab-itc")).toBeVisible();

    // Source badges should be visible (from SourceBadge component)
    await expect(reviewStep.getByText("From GSTR-1").first()).toBeVisible();

    // Click Next to go to Pay
    const nextBtn = page.getByTestId("next-button");
    await nextBtn.click();

    // Step 3: Pay
    const payStep = page.getByTestId("step-pay");
    await expect(payStep).toBeVisible();

    // Click Next to go to Sign
    await nextBtn.click();

    // Step 4: Sign
    const signStep = page.getByTestId("step-sign");
    await expect(signStep).toBeVisible();

    // EVC should be selected by default
    const evcButton = page.getByTestId("sign-evc");
    await expect(evcButton).toBeVisible();

    // Switch to DSC and back
    const dscButton = page.getByTestId("sign-dsc");
    await dscButton.click();
    await evcButton.click();

    // Click Next to go to File
    await nextBtn.click();

    // Step 5: File
    const fileStep = page.getByTestId("step-file");
    await expect(fileStep).toBeVisible();

    // Total payable should be visible
    const totalPayable = page.getByTestId("total-payable");
    await expect(totalPayable).toBeVisible();

    // File button in the step
    const fileButton = page.getByTestId("file-button");
    await expect(fileButton).toBeVisible();
    await fileButton.click();

    // Filing confirmation modal should appear
    const modal = page.getByRole("alertdialog");
    await expect(modal).toBeVisible();

    // Test cancel first
    const cancelBtn = page.getByTestId("modal-cancel-button");
    await cancelBtn.click();
    await expect(modal).not.toBeVisible();

    // Re-open modal via the Next button (which triggers confirm on file step)
    await nextBtn.click();
    await expect(page.getByRole("alertdialog")).toBeVisible();

    // Select EVC in modal
    const evcRadio = page.getByRole("alertdialog").getByLabel("EVC OTP");
    await evcRadio.check();

    // Confirm button should be disabled before typing
    const confirmBtn = page.getByTestId("modal-confirm-button");
    await expect(confirmBtn).toBeDisabled();

    // Type "FILE" to confirm
    const confirmInput = page.getByTestId("confirm-input");
    await confirmInput.fill("FILE");
    await expect(confirmBtn).toBeEnabled();
    await confirmBtn.click();

    // Step 6: Acknowledge — wait for filing to complete
    const ackStep = page.getByTestId("step-acknowledge");
    await expect(ackStep).toBeVisible({ timeout: 5000 });

    // Verify success message
    await expect(ackStep.getByText("GSTR-3B Filed Successfully")).toBeVisible();

    // Verify filing receipt with ARN
    const receipt = page.getByTestId("filing-receipt");
    await expect(receipt).toBeVisible();
    await expect(receipt.getByText(/AA29042026/)).toBeVisible();
    await expect(receipt.getByText("GSTR-3B")).toBeVisible();

    // Download button should exist
    const downloadBtn = page.getByTestId("download-ack-button");
    await expect(downloadBtn).toBeVisible();

    // Back to GST Returns link
    await expect(ackStep.getByText("Back to GST Returns")).toBeVisible();
  });

  test("step indicator shows all 6 steps", async ({ page }) => {
    await page.goto("/compliance/gst-returns/29AABCA1234A1Z5/2026-04/gstr-3b");

    const wizard = page.getByTestId("gstr3b-wizard");
    await expect(wizard).toBeVisible();

    // Step labels should be visible in the step indicator
    const indicator = page.getByTestId("step-indicator");
    await expect(indicator).toBeVisible();
    await expect(indicator.getByText("Auto-Populate", { exact: true })).toBeVisible();
    await expect(indicator.getByText("Review", { exact: true })).toBeVisible();
    await expect(indicator.getByText("Pay", { exact: true })).toBeVisible();
    await expect(indicator.getByText("Sign", { exact: true })).toBeVisible();
    await expect(indicator.getByText("File Return", { exact: true })).toBeVisible();
    await expect(indicator.getByText("Acknowledgement", { exact: true })).toBeVisible();
  });

  test("previous button navigates back", async ({ page }) => {
    await page.goto("/compliance/gst-returns/29AABCA1234A1Z5/2026-04/gstr-3b");

    // Auto-populate first
    const populateBtn = page.getByTestId("populate-button");
    await populateBtn.click();
    await expect(page.getByTestId("step-review")).toBeVisible({ timeout: 5000 });

    // Go to Pay
    await page.getByTestId("next-button").click();
    await expect(page.getByTestId("step-pay")).toBeVisible();

    // Go back to Review
    const prevBtn = page.getByTestId("previous-button");
    await prevBtn.click();
    await expect(page.getByTestId("step-review")).toBeVisible();
  });
});
