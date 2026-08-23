import { MigrationFormat } from "@linkwarden/types/global";

export const DEFAULT_IMPORT_LIMIT_MB = 10;

const SUPPORTED_MIGRATION_FORMATS = new Set<unknown>([
  MigrationFormat.htmlFile,
  MigrationFormat.linkwarden,
  MigrationFormat.wallabag,
  MigrationFormat.omnivore,
  MigrationFormat.pocket,
]);

export const getImportLimitMb = (value?: string) => {
  if (!value) return DEFAULT_IMPORT_LIMIT_MB;

  const parsedValue = Number(value);
  return Number.isFinite(parsedValue) && parsedValue > 0
    ? parsedValue
    : DEFAULT_IMPORT_LIMIT_MB;
};

export const isSupportedMigrationFormat = (
  format: unknown
): format is MigrationFormat =>
  typeof format === "string" && SUPPORTED_MIGRATION_FORMATS.has(format);

export const getExportFilename = (date = new Date()) =>
  `goreecloud-bookmarks-export-${date.toISOString().slice(0, 10)}.json`;
