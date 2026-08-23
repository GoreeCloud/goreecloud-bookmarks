import { describe, expect, it } from "vitest";
import browserCapture from "./capture";
import extensionCapture from "./extension-capture";

describe("native browser bookmark capture contract", () => {
  it("reuses the hardened authenticated capture handler", () => {
    expect(browserCapture).toBe(extensionCapture);
  });
});
