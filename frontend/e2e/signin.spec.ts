import { expect, test } from "@playwright/test";

test("signs in to an existing account with a passkey and stays signed in after reload", async ({ page, context }) => {
  await context.credentials.install();

  await page.goto("/sign-up");
  await page.getByRole("button", { name: "Create account with passkey" }).click();
  await expect(page).toHaveURL(/\/user$/);
  const account = page.getByText(/^Signed in as Laga /);
  const name = await account.textContent();

  // End the browser session without removing the account's passkey.
  await context.clearCookies();
  await page.reload();
  await expect(page).toHaveURL(/\/sign-up$/);
  await page.getByRole("link", { name: "Sign in", exact: true }).click();
  await expect(page).toHaveURL(/\/sign-in$/);

  await page.getByRole("button", { name: "Sign in with passkey" }).click();
  await expect(page).toHaveURL(/\/user$/);
  await expect(account).toHaveText(name ?? "");
  expect(await context.credentials.get({ rpId: "localhost" })).toHaveLength(1);

  await page.reload();
  await expect(page).toHaveURL(/\/user$/);
  await expect(account).toHaveText(name ?? "");
});
