import { test, expect } from "@playwright/test";

test.describe("ITR Magic Link Review (Employee-Facing)", () => {
  test("valid token renders read-only ITR display", async ({ page }) => {
    await page.goto("/itr/review/mlk-valid-001");
    await expect(page.getByTestId("magic-link-review")).toBeVisible({ timeout: 10000 });

    // Should show tax year and employer info
    await expect(page.getByText(/Tax Year/)).toBeVisible();

    // Should show ITR form recommendation
    await expect(page.getByText(/ITR-/)).toBeVisible();

    // Should show income section
    await expect(page.getByText("Income Summary")).toBeVisible();
  });

  test("approve and e-verify displays ARN", async ({ page }) => {
    await page.goto("/itr/review/mlk-valid-001");
    await expect(page.getByTestId("magic-link-review")).toBeVisible({ timeout: 10000 });

    // Click Approve & E-Verify
    const approveButton = page.getByRole("button", { name: /Approve/i });
    await approveButton.click();

    // Verify ARN displayed after approval
    await expect(page.getByText("ARN")).toBeVisible({ timeout: 10000 });
  });

  test("expired token shows expiry state", async ({ page }) => {
    await page.goto("/itr/review/mlk-expired-001");
    await expect(page.getByTestId("magic-link-review")).toBeVisible({ timeout: 10000 });

    // Should show expired heading
    await expect(page.getByRole("heading", { name: "Link Expired" })).toBeVisible();
  });

  test("unknown token shows expired/invalid state", async ({ page }) => {
    await page.goto("/itr/review/invalid-token-does-not-exist");
    await expect(page.getByTestId("magic-link-review")).toBeVisible({ timeout: 10000 });

    // Unknown tokens fall through to expired state in the mock
    await expect(page.getByRole("heading", { name: "Link Expired" })).toBeVisible();
  });
});
