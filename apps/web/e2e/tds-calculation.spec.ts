import { test, expect } from "@playwright/test";

test.describe("TDS Calculation", () => {
  test("contractor payment ₹50,000 for company deductee — verify section, code, rate, deduction", async ({ page }) => {
    await page.goto("/compliance/tds/calculate");

    // Verify page loaded
    await expect(page.getByText("Calculate TDS")).toBeVisible();

    // Step 1: Select deductee — search for Mahindra (company, resident, 393(1))
    const deducteeInput = page.locator("input[placeholder='Search by PAN or name...']");
    await deducteeInput.click();
    await deducteeInput.fill("Mahindra");

    // Click the Mahindra result in the dropdown
    await page.getByText("DEFCM3456H").click();

    // Step 2: Select payment code 1024 — Contractor Other
    const paymentSelect = page.locator("select").first();
    await paymentSelect.selectOption("1024");

    // Step 3: Enter amount ₹50,000
    const amountInput = page.locator("input[type='number']");
    await amountInput.fill("50000");

    // Step 4: Verify calculation panel
    // Section reference should show 393(1)[Sl.6(i).D(b)]
    await expect(page.getByText("393(1)[Sl.6(i).D(b)]")).toBeVisible();

    // Calculator panel should show payment code 1024 as a visible text element
    const calcPanel = page.locator("text=Live Calculation Preview").locator("..");
    await expect(calcPanel.getByText("1024")).toBeVisible();

    // Base rate 2%
    await expect(calcPanel.getByText("2%").first()).toBeVisible();

    // TDS deduction: 50000 * 2% = ₹1,000 (appears as TDS Amount and Total Deduction)
    await expect(calcPanel.getByText("₹1,000").first()).toBeVisible();

    // Threshold met indicator (30000 threshold, 50000 amount)
    await expect(page.getByText("Met")).toBeVisible();

    // Step 5: Save entry — type CONFIRM
    const confirmInput = page.locator("input[placeholder='CONFIRM']");
    await confirmInput.fill("CONFIRM");

    // Save button should now be enabled
    const saveButton = page.getByRole("button", { name: "Save Entry" });
    await expect(saveButton).toBeEnabled();
    await saveButton.click();

    // Verify success message
    await expect(page.getByText("Entry saved successfully")).toBeVisible();
  });
});
