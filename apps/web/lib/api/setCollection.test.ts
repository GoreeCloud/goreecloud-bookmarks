import { beforeEach, describe, expect, it, vi } from "vitest";

const mocks = vi.hoisted(() => ({
  findUnique: vi.fn(),
  findFirst: vi.fn(),
  create: vi.fn(),
  userUpdate: vi.fn(),
  getPermission: vi.fn(),
}));

vi.mock("@linkwarden/prisma", () => ({
  prisma: {
    collection: {
      findUnique: mocks.findUnique,
      findFirst: mocks.findFirst,
      create: mocks.create,
    },
    user: {
      update: mocks.userUpdate,
    },
  },
}));

vi.mock("./getPermission", () => ({
  default: mocks.getPermission,
}));

import setCollection from "./setCollection";

describe("setCollection authorization", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  it("allows the collection owner to capture into the collection", async () => {
    const collection = { id: 42, ownerId: 7, name: "Owner collection" };
    mocks.findUnique.mockResolvedValue(collection);
    mocks.getPermission.mockResolvedValue({ ownerId: 7, members: [] });

    await expect(
      setCollection({ userId: 7, collectionId: 42 })
    ).resolves.toEqual(collection);
  });

  it("allows a collection member with create permission", async () => {
    const collection = { id: 42, ownerId: 8, name: "Shared collection" };
    mocks.findUnique.mockResolvedValue(collection);
    mocks.getPermission.mockResolvedValue({
      ownerId: 8,
      members: [{ userId: 7, canCreate: true }],
    });

    await expect(
      setCollection({ userId: 7, collectionId: 42 })
    ).resolves.toEqual(collection);
  });

  it("rejects a different user without create permission", async () => {
    const collection = { id: 42, ownerId: 8, name: "Private collection" };
    mocks.findUnique.mockResolvedValue(collection);
    mocks.getPermission.mockResolvedValue({
      ownerId: 8,
      members: [{ userId: 7, canCreate: false }],
    });

    await expect(
      setCollection({ userId: 7, collectionId: 42 })
    ).resolves.toBeNull();
  });

  it("rejects a collection that does not exist", async () => {
    mocks.findUnique.mockResolvedValue(null);

    await expect(
      setCollection({ userId: 7, collectionId: 999 })
    ).resolves.toBeNull();
    expect(mocks.getPermission).not.toHaveBeenCalled();
  });
});
