import { expect, test } from "@playwright/test";

test("signs up with a passkey and stays signed in after reload", async ({ page, context }) => {
  await context.credentials.install();

  await page.goto("/user");
  await expect(page).toHaveURL(/\/sign-in$/);
  await page.getByRole("button", { name: "Sign up", exact: true }).click();

  await expect(page).toHaveURL(/\/sign-up$/);

  await page.getByRole("button", { name: "Create account with passkey" }).click();
  await expect(page).toHaveURL(/\/user$/);

  const account = page.getByText(/^You are signed in as: /);
  await expect(account).toBeVisible();
  const name = await account.textContent();
  expect(await context.credentials.get({ rpId: "localhost" })).toHaveLength(1);

  const cookies = await context.cookies();
  const session = cookies.find((cookie) => cookie.name === "laga_session");
  expect(session).toBeDefined();
  expect(session?.secure).toBe(true);
  expect(session?.httpOnly).toBe(true);
  expect(session?.sameSite).toBe("Strict");

  await page.reload();
  await expect(page).toHaveURL(/\/user$/);
  await expect(account).toHaveText(name ?? "");
});
