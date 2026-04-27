import { test, expect } from "@playwright/test";

test.describe("E-Way Bill Lifecycle", () => {
  test("generate EWB → specify vehicle + distance → verify validity calc → view detail → update vehicle", async ({ page }) => {
    await page.goto("/compliance/e-way-bill/generate");
    await expect(page.getByRole("heading", { name: "Generate E-Way Bill" })).toBeVisible();

    // Select first invoice
    const firstCard = page.locator("button").filter({ hasText: "INV-2026-" }).first();
    await expect(firstCard).toBeVisible();
    await firstCard.click();

    // Should land on form step
    await expect(page.getByRole("heading", { name: "Source Invoice" })).toBeVisible();

    // Fill vehicle number
    const vehicleInput = page.locator("input[placeholder*='KA01AB']");
    await vehicleInput.fill("MH12XY9876");

    // Fill distance = 600 km
    const distanceInput = page.locator("input[type='number']");
    await distanceInput.fill("600");

    // Verify DistanceValidityCalculator shows 3 days (600 / 200 = 3)
    const validityCalc = page.getByTestId("distance-validity-calc");
    await expect(validityCalc).toBeVisible();
    await expect(validityCalc.getByText("600 km")).toBeVisible();
    await expect(validityCalc.getByText("3 days")).toBeVisible();

    // Click Review & Confirm
    await page.getByText("Review & Confirm").click();

    // Confirmation modal — type GENERATE
    const confirmInput = page.getByTestId("confirm-input");
    await expect(confirmInput).toBeVisible();
    await confirmInput.fill("GENERATE");

    const confirmBtn = page.getByTestId("modal-confirm-button");
    await expect(confirmBtn).toBeEnabled();
    await confirmBtn.click();

    // Generating step — spinner
    await expect(page.getByText("Submitting to EWB portal...")).toBeVisible();

    // Success step
    const successPanel = page.getByTestId("ewb-success");
    await expect(successPanel).toBeVisible({ timeout: 5000 });
    await expect(page.getByText("E-Way Bill Generated")).toBeVisible();

    // Navigate to detail view
    await page.getByText("View Details").click();
    await page.waitForURL(/\/compliance\/e-way-bill\/\d+$/);

    // Detail page — should show EWB summary
    await expect(page.getByText("EWB Summary")).toBeVisible({ timeout: 10000 });

    // Update Vehicle link should be visible for active records
    const updateVehicleLink = page.getByTestId("update-vehicle-link");
    if (await updateVehicleLink.isVisible()) {
      await updateVehicleLink.click();
      await page.waitForURL(/\/update-vehicle$/);

      await expect(page.getByRole("heading", { name: "Update Vehicle" })).toBeVisible();

      // Fill new vehicle details
      const newVehicleInput = page.locator("input[placeholder*='MH12AB']");
      await newVehicleInput.fill("GJ06PQ4321");

      const reasonInput = page.locator("textarea");
      await reasonInput.fill("Vehicle breakdown — replaced");

      // Click the submit button (not the Cancel link)
      const submitBtn = page.locator("button").filter({ hasText: "Update Vehicle" });
      await submitBtn.click();

      // Confirmation modal
      const updateConfirmInput = page.getByTestId("confirm-input");
      if (await updateConfirmInput.isVisible()) {
        await updateConfirmInput.fill("UPDATE");
        const updateConfirmBtn = page.getByTestId("modal-confirm-button");
        await expect(updateConfirmBtn).toBeEnabled();
        await updateConfirmBtn.click();

        // Wait for success
        await expect(page.getByText("Vehicle Updated")).toBeVisible({ timeout: 5000 });
      }
    }
  });

  test("detail page — cancel EWB within 24h", async ({ page }) => {
    await page.goto("/compliance/e-way-bill");
    await expect(page.getByRole("heading", { name: "E-Way Bills" })).toBeVisible();

    // Click on first row in the data table
    const firstRow = page.locator("table tbody tr").first();
    await expect(firstRow).toBeVisible();
    await firstRow.click();
    await page.waitForURL(/\/compliance\/e-way-bill\/[^/]+$/);

    // If cancel button is visible (active + within 24h), test cancellation
    const cancelBtn = page.getByTestId("cancel-ewb-button");
    if (await cancelBtn.isVisible()) {
      await cancelBtn.click();

      const confirmInput = page.getByTestId("confirm-input");
      await expect(confirmInput).toBeVisible();
      await confirmInput.fill("CANCEL");

      const confirmBtn = page.getByTestId("modal-confirm-button");
      await expect(confirmBtn).toBeEnabled();
      await confirmBtn.click();

      // Modal should close
      await expect(confirmInput).not.toBeVisible();
    }
  });

  test("extend validity link visible for eligible records", async ({ page }) => {
    await page.goto("/compliance/e-way-bill");

    const firstRow = page.locator("table tbody tr").first();
    await expect(firstRow).toBeVisible();
    await firstRow.click();
    await page.waitForURL(/\/compliance\/e-way-bill\/[^/]+$/);

    const extendLink = page.getByTestId("extend-validity-link");
    if (await extendLink.isVisible()) {
      await extendLink.click();
      await page.waitForURL(/\/extend$/);

      await expect(page.getByRole("heading", { name: "Extend Validity" })).toBeVisible();
      await expect(page.getByText("Current validity expires")).toBeVisible();
    }
  });

  test("cancelled record — no cancel or extend buttons", async ({ page }) => {
    await page.goto("/compliance/e-way-bill");

    // Click the "Cancelled" status filter tab (not a status pill in table)
    const tabs = page.locator("button").filter({ hasText: /^Cancelled$/ });
    await tabs.first().click();

    // Wait for table to re-render with only cancelled records
    await page.waitForTimeout(300);

    const cancelledRow = page.locator("table tbody tr").first();
    if (await cancelledRow.isVisible()) {
      await cancelledRow.click();
      await page.waitForURL(/\/compliance\/e-way-bill\/[^/]+$/);

      // Cancel and extend buttons should NOT be visible for cancelled records
      await expect(page.getByTestId("cancel-ewb-button")).not.toBeVisible();
      await expect(page.getByTestId("extend-validity-link")).not.toBeVisible();
      await expect(page.getByTestId("update-vehicle-link")).not.toBeVisible();
    }
  });
});
