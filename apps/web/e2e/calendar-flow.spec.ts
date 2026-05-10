import { test, expect } from "@playwright/test";

const MONTH_NAMES = [
  "January", "February", "March", "April", "May", "June",
  "July", "August", "September", "October", "November", "December",
];

test.describe("Compliance Calendar", () => {
  test("full calendar lifecycle: grid → events → detail panel → filters → navigation", async ({ page }) => {
    const coldStart = Date.now();

    await page.goto("/compliance/calendar");

    // Wait for skeleton to appear (or grid to load directly on fast cache)
    const skeleton = page.getByTestId("calendar-skeleton");
    const grid = page.getByTestId("calendar-grid");
    await expect(skeleton.or(grid)).toBeVisible({ timeout: 5000 });

    // Wait for the real grid to render (skeleton replaced)
    await expect(grid).toBeVisible({ timeout: 10000 });

    const coldLoadMs = Date.now() - coldStart;
    console.log(`[perf] Cold calendar load: ${coldLoadMs}ms (target <2000ms)`);
    if (coldLoadMs > 2000) {
      console.warn(`[perf] Cold load exceeded 2s target: ${coldLoadMs}ms`);
    }

    // Verify month grid structure: 7 columns of weekday headers
    const weekdayHeaders = grid.locator("div.grid-cols-7").first();
    await expect(weekdayHeaders).toBeVisible();
    const headerCells = weekdayHeaders.locator("> div");
    await expect(headerCells).toHaveCount(7);

    // Verify 42 date cells (6 rows × 7 cols)
    const dateCells = grid.locator("[data-testid^='day-cell-']");
    await expect(dateCells).toHaveCount(42);

    // Verify today's cell has the amber highlight indicator
    const todayIndicator = page.getByTestId("today-indicator");
    await expect(todayIndicator).toBeVisible();
    await expect(todayIndicator).toHaveText("Today");

    // Verify at least one event pill is rendered
    const eventPills = page.getByTestId("event-pill");
    expect(await eventPills.count()).toBeGreaterThan(0);

    // Find a day cell that has event pills and click it
    const firstPill = eventPills.first();
    const dayWithEvents = firstPill.locator("xpath=ancestor::button[@data-testid]");
    await dayWithEvents.click();

    // Verify slide-out panel opens
    const panel = page.getByTestId("day-detail-panel");
    await expect(panel).toBeVisible();

    // Verify panel shows item details
    const eventDetail = panel.getByTestId("event-detail").first();
    await expect(eventDetail).toBeVisible();

    // Verify authority text is visible
    const authorityLabel = eventDetail.getByText(/Authority/);
    await expect(authorityLabel).toBeVisible();

    // Close the slide-out via close button
    const closeBtn = page.getByTestId("close-panel");
    await closeBtn.click();

    // Verify panel is dismissed (translated off-screen)
    await expect(panel).toHaveClass(/translate-x-full/);

    // Toggle a category filter badge — hide indirect tax
    const indirectTaxBadge = page.getByTestId("category-indirect_tax");
    await expect(indirectTaxBadge).toBeVisible();
    const pillCountBefore = await eventPills.count();

    await indirectTaxBadge.click();

    // Verify opacity changed (filter toggled off)
    await expect(indirectTaxBadge).toHaveClass(/opacity-40/);

    // Some pills should have disappeared (or count changed)
    const pillCountAfterHide = await eventPills.count();
    expect(pillCountAfterHide).toBeLessThanOrEqual(pillCountBefore);

    // Toggle it back
    await indirectTaxBadge.click();
    await expect(indirectTaxBadge).toHaveClass(/opacity-100/);

    // Pills should be restored
    const pillCountAfterRestore = await eventPills.count();
    expect(pillCountAfterRestore).toBe(pillCountBefore);

    // Click "Next month" and verify month label changes
    const monthLabel = page.getByTestId("month-label");
    const currentMonthText = await monthLabel.textContent();
    const nextMonthBtn = page.getByTestId("next-month");
    await nextMonthBtn.click();

    // Wait for the grid to reload after navigation
    await expect(grid).toBeVisible({ timeout: 10000 });

    const nextMonthText = await monthLabel.textContent();
    expect(nextMonthText).not.toBe(currentMonthText);

    // Click "Today" and verify month returns to current
    const todayBtn = page.getByTestId("today-button");
    await todayBtn.click();

    await expect(grid).toBeVisible({ timeout: 10000 });

    const now = new Date();
    const expectedLabel = `${MONTH_NAMES[now.getMonth()]} ${now.getFullYear()}`;
    await expect(monthLabel).toHaveText(expectedLabel);

    // Warm cache reload timing
    const warmStart = Date.now();
    await page.reload();
    await expect(grid).toBeVisible({ timeout: 10000 });
    const warmLoadMs = Date.now() - warmStart;
    console.log(`[perf] Warm cache reload: ${warmLoadMs}ms (target <500ms)`);
    if (warmLoadMs > 500) {
      console.warn(`[perf] Warm reload exceeded 500ms target: ${warmLoadMs}ms`);
    }
  });

  test("stat cards render with correct labels", async ({ page }) => {
    await page.goto("/compliance/calendar");
    const grid = page.getByTestId("calendar-grid");
    await expect(grid).toBeVisible({ timeout: 10000 });

    await expect(page.getByText("Filed")).toBeVisible();
    await expect(page.getByText("Due in 7 days")).toBeVisible();
    await expect(page.getByText("Upcoming this month")).toBeVisible();
  });

  test("day detail panel closes on Escape key", async ({ page }) => {
    await page.goto("/compliance/calendar");
    const grid = page.getByTestId("calendar-grid");
    await expect(grid).toBeVisible({ timeout: 10000 });

    // Click a day with events to open the panel
    const firstPill = page.getByTestId("event-pill").first();
    const dayWithEvents = firstPill.locator("xpath=ancestor::button[@data-testid]");
    await dayWithEvents.click();

    const panel = page.getByTestId("day-detail-panel");
    await expect(panel).toBeVisible();

    // Press Escape
    await page.keyboard.press("Escape");

    // Panel should slide out
    await expect(panel).toHaveClass(/translate-x-full/);
  });
});
