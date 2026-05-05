import { test, expect } from "@playwright/test";

test.describe("GSTR-9 + 9C Annual Return Flow", () => {
  const FY = "2025-26";
  const GSTIN = "29AABCA1234A1Z5";

  test("GSTR-9 wizard — walk through steps and file", async ({ page }) => {
    await page.goto(`/compliance/gst-returns/annual/${FY}/${GSTIN}/9`);
    const wizard = page.getByTestId("gstr9-wizard");
    await expect(wizard).toBeVisible({ timeout: 10000 });

    // Step 1: Threshold — verify turnover displayed
    await expect(wizard.getByText(/Step 1: Threshold/)).toBeVisible();
    await expect(wizard.getByText(/Aggregate Turnover/)).toBeVisible();

    // Scope Next button within wizard to avoid Next.js dev tools button
    const nextButton = wizard.getByRole("button", { name: "Next" });

    // Step 2: Review Tables
    await nextButton.click();
    await expect(wizard.getByText(/Step 2/)).toBeVisible();

    // Step 3: Late ITC
    await nextButton.click();
    await expect(wizard.getByText(/Step 3/)).toBeVisible();

    // Step 4: HSN Summary
    await nextButton.click();
    await expect(wizard.getByText(/Step 4/)).toBeVisible();

    // Step 5: Fees & Demands
    await nextButton.click();
    await expect(wizard.getByText(/Step 5/)).toBeVisible();

    // Step 6: Sign & Submit — Sign DSC first
    await nextButton.click();
    await expect(wizard.getByText(/Step 6/)).toBeVisible();

    // Sign DSC
    const signButton = wizard.getByRole("button", { name: /Sign GSTR-9 with DSC/i });
    await signButton.click();

    // Wait for signing to complete
    await expect(wizard.getByText(/Signed with DSC/i)).toBeVisible({ timeout: 5000 });

    // File GSTR-9 — triggers confirmation modal
    const fileButton = wizard.getByRole("button", { name: /File GSTR-9/i });
    await expect(fileButton).toBeEnabled();
    await fileButton.click();

    // Confirmation modal — type FILE and confirm
    const confirmInput = page.getByPlaceholder(/FILE/i);
    await expect(confirmInput).toBeVisible({ timeout: 3000 });
    await confirmInput.fill("FILE");
    const confirmBtn = page.getByRole("button", { name: /Confirm/i });
    await confirmBtn.click();

    // Verify ARN is displayed
    await expect(wizard.getByText(/ARN/)).toBeVisible({ timeout: 10000 });
  });

  test("GSTR-9C wizard — threshold check + file flow", async ({ page }) => {
    await page.goto(`/compliance/gst-returns/annual/${FY}/${GSTIN}/9c`);
    const wizard = page.getByTestId("gstr9c-wizard");
    await expect(wizard).toBeVisible({ timeout: 10000 });

    // Step 1: Threshold Check
    await expect(page.getByTestId("step-threshold-check")).toBeVisible();
    await expect(wizard.getByText("₹5 Cr", { exact: true })).toBeVisible();
    await expect(wizard.getByText(/mandatory/i)).toBeVisible();

    const nextButton = wizard.getByRole("button", { name: "Next" });
    await nextButton.click();

    // Step 2: Upload Financials
    await expect(page.getByTestId("step-upload-financials")).toBeVisible();
    const uploadButton = wizard.getByRole("button", { name: /Mock Financials/i });
    await uploadButton.click();
    await expect(wizard.getByText(/loaded/i)).toBeVisible();
    await nextButton.click();

    // Step 3: Reconciliation (split-pane sections)
    await expect(page.getByTestId("step-reconciliation")).toBeVisible();
    await expect(wizard.getByText("Reconciliation by Section")).toBeVisible();
    await expect(wizard.getByText("Books (Audited) — Part II", { exact: true })).toBeVisible();
    await expect(wizard.getByText("Books (Audited) — Part III", { exact: true })).toBeVisible();
    await nextButton.click();

    // Step 4: Resolve Mismatches
    await expect(page.getByTestId("step-resolve-mismatches")).toBeVisible();
    await expect(wizard.getByText(/unresolved/)).toBeVisible();

    // Resolve all ERROR mismatches
    const resolveButtons = wizard.getByRole("button", { name: "Resolve" });
    const count = await resolveButtons.count();
    for (let i = 0; i < count; i++) {
      const btn = wizard.getByRole("button", { name: "Resolve" }).first();
      await btn.click();
      await expect(page.getByTestId("resolve-modal")).toBeVisible();
      const textarea = page.getByTestId("resolve-modal").locator("textarea");
      await textarea.fill("Timing difference — reconciled in subsequent period with supporting documents");
      const resolveBtn = page.getByTestId("resolve-modal").getByRole("button", { name: /Resolve|Acknowledge/i });
      await resolveBtn.click();
      await expect(page.getByTestId("resolve-modal")).not.toBeVisible();
    }

    // Proceed to certification
    await nextButton.click();

    // Step 5: Self-Certification
    await expect(page.getByTestId("certification-form")).toBeVisible();
    const certInput = page.getByPlaceholder(`I CERTIFY GSTR-9C ${FY} ${GSTIN}`);
    await certInput.fill(`I CERTIFY GSTR-9C ${FY} ${GSTIN}`);
    const certButton = wizard.getByRole("button", { name: /Certify/i });
    await certButton.click();
    await expect(page.getByTestId("certification-locked")).toBeVisible();
    await nextButton.click();

    // Step 6: File with DSC
    await expect(page.getByTestId("step-file-dsc")).toBeVisible();
    const signButton = wizard.getByRole("button", { name: /Sign GSTR-9C/i });
    await signButton.click();
    await expect(wizard.getByText(/Signed/i)).toBeVisible({ timeout: 5000 });

    // File GSTR-9C
    const fileButton = wizard.getByRole("button", { name: /File GSTR-9C/i });
    await fileButton.click();

    // Verify ARN
    await expect(page.getByTestId("step-arn-9c")).toBeVisible({ timeout: 10000 });
    await expect(wizard.getByText(/ARN/)).toBeVisible();
  });

  test("GSTR-9C blocks proceed when ERROR mismatches unresolved", async ({ page }) => {
    await page.goto(`/compliance/gst-returns/annual/${FY}/${GSTIN}/9c`);
    const wizard = page.getByTestId("gstr9c-wizard");
    await expect(wizard).toBeVisible({ timeout: 10000 });

    const nextButton = wizard.getByRole("button", { name: "Next" });

    // Step 1 → Step 2
    await nextButton.click();
    await expect(page.getByTestId("step-upload-financials")).toBeVisible();

    // Upload financials
    const uploadButton = wizard.getByRole("button", { name: /Mock Financials/i });
    await uploadButton.click();

    // Step 2 → Step 3
    await nextButton.click();
    await expect(page.getByTestId("step-reconciliation")).toBeVisible();

    // Step 3 → Step 4
    await nextButton.click();
    await expect(page.getByTestId("step-resolve-mismatches")).toBeVisible();

    // Verify ERROR mismatches present (use the unresolved text) and Next is disabled
    await expect(wizard.getByText(/unresolved — blocks proceed/)).toBeVisible();
    await expect(nextButton).toBeDisabled();
  });
});
