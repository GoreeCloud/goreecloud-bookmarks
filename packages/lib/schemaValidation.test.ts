import { describe, expect, it } from "vitest";
import { TokenExpiry } from "@linkwarden/types/global";
import { PostTokenSchema } from "./schemaValidation";

describe("PostTokenSchema", () => {
  it("rejects blank access-token names", () => {
    expect(
      PostTokenSchema.safeParse({
        name: "",
        expires: TokenExpiry.sevenDays,
      }).success
    ).toBe(false);

    expect(
      PostTokenSchema.safeParse({
        name: "   ",
        expires: TokenExpiry.sevenDays,
      }).success
    ).toBe(false);
  });

  it("trims a valid access-token name", () => {
    const result = PostTokenSchema.parse({
      name: "  backup automation  ",
      expires: TokenExpiry.sevenDays,
    });

    expect(result.name).toBe("backup automation");
    expect(result.expires).toBe(TokenExpiry.sevenDays);
  });

  it("keeps the existing 50-character name limit", () => {
    expect(
      PostTokenSchema.safeParse({
        name: "a".repeat(51),
        expires: TokenExpiry.sevenDays,
      }).success
    ).toBe(false);
  });
});
