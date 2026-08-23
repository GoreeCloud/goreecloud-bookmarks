import { beforeEach, describe, expect, it, vi } from "vitest";

const mocks = vi.hoisted(() => ({
  getTokenFromRequest: vi.fn(),
  findRevokedToken: vi.fn(),
}));

vi.mock("@linkwarden/prisma", () => ({
  prisma: {
    accessToken: {
      findFirst: mocks.findRevokedToken,
    },
  },
}));

vi.mock("./getTokenFromRequest", () => ({
  default: mocks.getTokenFromRequest,
}));

import verifyToken from "./verifyToken";

describe("verifyToken authorization", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("rejects a request without an authenticated user token", async () => {
    mocks.getTokenFromRequest.mockResolvedValue(null);

    await expect(verifyToken({ req: {} as never })).resolves.toBe(
      "You must be logged in."
    );
    expect(mocks.findRevokedToken).not.toHaveBeenCalled();
  });

  it("rejects an expired token before querying revocation state", async () => {
    mocks.getTokenFromRequest.mockResolvedValue({
      id: 7,
      jti: "expired-token",
      exp: Math.floor(Date.now() / 1000) - 60,
    });

    await expect(verifyToken({ req: {} as never })).resolves.toBe(
      "Your session has expired, please log in again."
    );
    expect(mocks.findRevokedToken).not.toHaveBeenCalled();
  });

  it("rejects a revoked access token", async () => {
    const token = {
      id: 7,
      jti: "revoked-token",
      exp: Math.floor(Date.now() / 1000) + 3600,
    };
    mocks.getTokenFromRequest.mockResolvedValue(token);
    mocks.findRevokedToken.mockResolvedValue({ token: token.jti, revoked: true });

    await expect(verifyToken({ req: {} as never })).resolves.toBe(
      "Your session has expired, please log in again."
    );
    expect(mocks.findRevokedToken).toHaveBeenCalledWith({
      where: { token: token.jti, revoked: true },
    });
  });

  it("accepts an active, unrevoked token", async () => {
    const token = {
      id: 7,
      jti: "active-token",
      exp: Math.floor(Date.now() / 1000) + 3600,
    };
    mocks.getTokenFromRequest.mockResolvedValue(token);
    mocks.findRevokedToken.mockResolvedValue(null);

    await expect(verifyToken({ req: {} as never })).resolves.toEqual(token);
  });
});
