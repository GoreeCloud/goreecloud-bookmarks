import type { NextApiRequest } from "next";
import { describe, expect, it } from "vitest";
import { isBrowserExtensionRequestAllowed } from "./browserExtensionScope";

const request = (method: string, url: string) =>
  ({ method, url } as NextApiRequest);

describe("browser extension API scope", () => {
  it.each([
    ["GET", "/api/v1/config"],
    ["GET", "/api/v1/collections"],
    ["GET", "/api/v1/tags?sort=2&cursor=0"],
    ["GET", "/api/v1/search?sort=0&searchQueryString=url%3Ahttps%3A%2F%2Fexample.com"],
    ["POST", "/api/v1/links"],
    ["PUT", "/api/v1/links/42"],
    ["DELETE", "/api/v1/links/42"],
    ["POST", "/api/v1/archives/42?format=0"],
    ["DELETE", "/api/v1/session"],
  ])("allows %s %s", (method, url) => {
    expect(isBrowserExtensionRequestAllowed(request(method, url))).toBe(true);
  });

  it.each([
    ["GET", "/api/v1/tokens"],
    ["DELETE", "/api/v1/tokens/1"],
    ["GET", "/api/v1/users"],
    ["POST", "/api/v1/collections"],
    ["DELETE", "/api/v1/collections/1"],
    ["GET", "/api/v1/links/42"],
    ["POST", "/api/v1/archives/not-a-number"],
    ["GET", "/api/v1/session"],
  ])("denies %s %s", (method, url) => {
    expect(isBrowserExtensionRequestAllowed(request(method, url))).toBe(false);
  });
});
