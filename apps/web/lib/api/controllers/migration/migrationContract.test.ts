import { describe, expect, it } from "vitest";
import { MigrationFormat } from "@linkwarden/types/global";
import {
  DEFAULT_IMPORT_LIMIT_MB,
  getExportFilename,
  getImportLimitMb,
  isSupportedMigrationFormat,
} from "./migrationContract";

describe("migrationContract", () => {
  describe("getImportLimitMb", () => {
    it("uses the configured positive limit", () => {
      expect(getImportLimitMb("25")).toBe(25);
      expect(getImportLimitMb("2.5")).toBe(2.5);
    });

    it("fails closed to the default when configuration is absent or invalid", () => {
      expect(getImportLimitMb()).toBe(DEFAULT_IMPORT_LIMIT_MB);
      expect(getImportLimitMb("0")).toBe(DEFAULT_IMPORT_LIMIT_MB);
      expect(getImportLimitMb("-1")).toBe(DEFAULT_IMPORT_LIMIT_MB);
      expect(getImportLimitMb("not-a-number")).toBe(DEFAULT_IMPORT_LIMIT_MB);
      expect(getImportLimitMb("Infinity")).toBe(DEFAULT_IMPORT_LIMIT_MB);
    });
  });

  describe("isSupportedMigrationFormat", () => {
    it("accepts every supported import format", () => {
      expect(isSupportedMigrationFormat(MigrationFormat.htmlFile)).toBe(true);
      expect(isSupportedMigrationFormat(MigrationFormat.linkwarden)).toBe(true);
      expect(isSupportedMigrationFormat(MigrationFormat.wallabag)).toBe(true);
      expect(isSupportedMigrationFormat(MigrationFormat.omnivore)).toBe(true);
      expect(isSupportedMigrationFormat(MigrationFormat.pocket)).toBe(true);
    });

    it("rejects missing and unknown formats", () => {
      expect(isSupportedMigrationFormat(undefined)).toBe(false);
      expect(isSupportedMigrationFormat(null)).toBe(false);
      expect(isSupportedMigrationFormat("unknown-format")).toBe(false);
    });
  });

  it("creates a deterministic, product-specific export filename", () => {
    expect(getExportFilename(new Date("2026-08-22T12:00:00.000Z"))).toBe(
      "goreecloud-bookmarks-export-2026-08-22.json"
    );
  });
});
