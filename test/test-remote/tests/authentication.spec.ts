/**
 * Copyright 2021 Opstrace, Inc.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

import { test } from "../fixtures/authentication";
import { expect } from "@playwright/test";

import { CLUSTER_BASE_URL, CI_LOGIN_EMAIL } from "../testutils";

test.describe("auth0 authentication", () => {
  test.beforeEach(async ({ context, page, authCookies }) => {
    context.addCookies(authCookies);
    await page.goto(CLUSTER_BASE_URL);
    await page.waitForSelector("text=Getting Started");
  });

  test("should see homepage", async ({ page }) => {
    expect(await page.isVisible("text=Getting Started")).toBeTruthy();
  });

  test("should have self in user list", async ({ page }) => {
    await page.click("text=Users");
    expect(await page.isVisible(`text=${CI_LOGIN_EMAIL}`)).toBeTruthy();
  });
});
