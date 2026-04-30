import { test, expect } from "@playwright/test";

test.describe("TDS Form 144 (Non-Resident) Filing", () => {
  test("DTAA validation blocks submit when evidence missing", async ({ page }) => {
    await page.goto("/compliance/tds/file/144/2026-27/q1");

    const main = page.getByTestId("main-content");

    // Verify page loaded
    await expect(page.getByText("Form 144 (Non-Resident) Filing")).toBeVisible();
    await expect(page.getByText("Tax Year 2026-27")).toBeVisible();

    // Step 1: Pull
    const pullButton = main.getByRole("button", { name: /Pull/i });
    await pullButton.click();

    // Wait for Step 2: Validate
    await expect(page.getByText("Step 2: Validate")).toBeVisible({ timeout: 5000 });

    // Verify DTAA validation flags present
    await expect(page.getByText("DTAA").first()).toBeVisible();

    // Navigate to Step 3: FVU Preview
    const nextButton = main.getByRole("button", { name: "Next" });
    await nextButton.click();
    await expect(page.getByText("Step 3: FVU Preview")).toBeVisible();

    // Navigate to Step 4: Submit to TRACES
    await nextButton.click();
    await expect(page.getByText("Step 4: Submit to TRACES")).toBeVisible();

    // Verify submit button is blocked due to DTAA evidence issues
    const blockedButton = main.getByRole("button", { name: /Blocked.*DTAA/i });
    await expect(blockedButton).toBeVisible();
    await expect(blockedButton).toBeDisabled();

    // The "File Now" button in the footer should also be disabled
    const fileNowButton = main.getByRole("button", { name: "File Now" });
    await expect(fileNowButton).toBeDisabled();
  });

  test("Form 138 (salary) wizard reaches acknowledgement with ARN", async ({ page }) => {
    await page.goto("/compliance/tds/file/138/2026-27/q1");

    const main = page.getByTestId("main-content");

    // Verify page loaded
    await expect(page.getByText("Form 138 (Salary) Filing")).toBeVisible();

    // Step 1: Pull
    const pullButton = main.getByRole("button", { name: /Pull/i });
    await pullButton.click();
    await expect(page.getByText("Step 2: Validate")).toBeVisible({ timeout: 5000 });

    // Step 2 → Step 3 → Step 4
    const nextButton = main.getByRole("button", { name: "Next" });
    await nextButton.click();
    await expect(page.getByText("Step 3: FVU Preview")).toBeVisible();

    await nextButton.click();
    await expect(page.getByText("Step 4: Submit to TRACES")).toBeVisible();

    // File button should be enabled (no DTAA blockers for salary)
    const fileNowButton = main.getByRole("button", { name: "File Now" });
    await expect(fileNowButton).toBeEnabled();
    await fileNowButton.click();

    // Confirmation modal should appear (role="alertdialog")
    const modal = page.getByRole("alertdialog");
    await expect(modal).toBeVisible({ timeout: 2000 });

    // Type confirmation word
    const confirmInput = modal.locator("input");
    await confirmInput.fill("FILE Q1 2026-27");

    // Click confirm button
    const confirmButton = modal.getByRole("button", { name: /Confirm/i });
    await expect(confirmButton).toBeEnabled();
    await confirmButton.click();

    // Step 5: Acknowledge — wait for filing to complete
    // Verify ARN is displayed (format: TDS138Q1202627XXXX)
    await expect(page.getByText(/TDS138Q1/)).toBeVisible({ timeout: 5000 });
  });
});
