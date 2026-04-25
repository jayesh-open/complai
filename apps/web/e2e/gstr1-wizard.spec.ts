import { test, expect } from "@playwright/test";

test.describe("GSTR-1 Filing Wizard", () => {
  test("full wizard lifecycle: Ingest → Validate → Review → Approve → File → Acknowledge", async ({ page }) => {
    await page.goto("/compliance/gst/gstr1");

    // Step 1: Ingest
    const ingestBtn = page.getByTestId("ingest-button");
    await expect(ingestBtn).toBeVisible();
    await expect(ingestBtn).toContainText("Ingest from Aura");
    await ingestBtn.click();

    // Step 2: Validate — wait for ingest to finish and auto-advance
    const validateBtn = page.getByTestId("validate-button");
    await expect(validateBtn).toBeVisible({ timeout: 5000 });
    await expect(validateBtn).toContainText("Run Validation");
    await validateBtn.click();

    // Step 3: Review — wait for validation to complete
    const sectionTable = page.getByTestId("section-table");
    await expect(sectionTable).toBeVisible({ timeout: 5000 });

    // Verify section rows are rendered (B2B should be present)
    await expect(sectionTable.getByText("B2B Invoices")).toBeVisible();
    await expect(sectionTable.getByText("B2C Large")).toBeVisible();

    const proceedBtn = page.getByTestId("proceed-approve-button");
    await expect(proceedBtn).toBeVisible();
    await proceedBtn.click();

    // Step 4: Approve
    const approveBtn = page.getByTestId("approve-button");
    await expect(approveBtn).toBeVisible();
    await expect(approveBtn).toContainText("Approve & Continue");
    await approveBtn.click();

    // Step 5: File — wait for approval
    const fileBtn = page.getByTestId("file-button");
    await expect(fileBtn).toBeVisible({ timeout: 5000 });
    await expect(fileBtn).toContainText("File GSTR-1");

    // Open confirmation modal
    await fileBtn.click();

    // Modal should be visible
    const modal = page.getByRole("alertdialog");
    await expect(modal).toBeVisible();

    // Test Cancel path first
    const cancelBtn = page.getByTestId("modal-cancel-button");
    await expect(cancelBtn).toBeVisible();
    await cancelBtn.click();

    // Modal should close, file button should still be visible
    await expect(modal).not.toBeVisible();
    await expect(fileBtn).toBeVisible();

    // Open modal again for the actual filing
    await fileBtn.click();
    await expect(page.getByRole("alertdialog")).toBeVisible();

    // Select EVC radio
    const evcRadio = page.getByRole("alertdialog").getByLabel("EVC OTP");
    await evcRadio.check();

    // Confirm button should be disabled before typing
    const confirmBtn = page.getByTestId("modal-confirm-button");
    await expect(confirmBtn).toBeDisabled();

    // Type "FILE" to confirm
    const confirmInput = page.getByTestId("confirm-input");
    await confirmInput.fill("FILE");

    // Confirm button should now be enabled
    await expect(confirmBtn).toBeEnabled();
    await confirmBtn.click();

    // Step 6: Acknowledge — wait for filing to complete
    const ackStep = page.getByTestId("acknowledge-step");
    await expect(ackStep).toBeVisible({ timeout: 5000 });

    // Verify success message
    await expect(ackStep.getByText("GSTR-1 Filed Successfully")).toBeVisible();

    // Verify ARN is displayed (format: AA290420260000XXXX)
    const receipt = page.getByTestId("filing-receipt");
    await expect(receipt).toBeVisible();
    await expect(receipt.getByText(/AA29042026/)).toBeVisible();

    // Verify "Back to GST Returns" link exists
    const backLink = ackStep.getByText("Back to GST Returns");
    await expect(backLink).toBeVisible();
  });
});
