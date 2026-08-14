import { spawnSync } from "node:child_process";

const skipBrowserInstall =
  process.env.GOREECLOUD_SKIP_PLAYWRIGHT_BROWSER_INSTALL === "1";

if (skipBrowserInstall) {
  console.log(
    "Skipping Playwright browser installation for this GoreeCloud CI workflow."
  );
  process.exit(0);
}

const result = spawnSync(
  process.platform === "win32" ? "playwright.cmd" : "playwright",
  ["install", "--with-deps", "chromium"],
  {
    stdio: "inherit",
    shell: process.platform === "win32",
  }
);

if (result.error) {
  throw result.error;
}

process.exit(result.status ?? 1);
