import type { NextApiRequest, NextApiResponse } from "next";
import exportData from "@/lib/api/controllers/migration/exportData";
import importFromHTMLFile from "@/lib/api/controllers/migration/importFromHTMLFile";
import importFromLinkwarden from "@/lib/api/controllers/migration/importFromLinkwarden";
import { MigrationFormat, MigrationRequest } from "@linkwarden/types/global";
import verifyUser from "@/lib/api/verifyUser";
import importFromWallabag from "@/lib/api/controllers/migration/importFromWallabag";
import importFromOmnivore from "@/lib/api/controllers/migration/importFromOmnivore";
import importFromPocket from "@/lib/api/controllers/migration/importFromPocket";
import {
  getExportFilename,
  getImportLimitMb,
  isSupportedMigrationFormat,
} from "@/lib/api/controllers/migration/migrationContract";

export const config = {
  api: {
    bodyParser: false,
  },
};

const parseJsonStream = (
  req: NextApiRequest,
  limitMb: number
): Promise<unknown> => {
  return new Promise((resolve, reject) => {
    let body = "";
    let totalLength = 0;
    let rejected = false;
    const limitBytes = limitMb * 1024 * 1024;

    req.setEncoding("utf8");

    req.on("data", (chunk: string) => {
      if (rejected) return;

      totalLength += Buffer.byteLength(chunk, "utf8");
      if (totalLength > limitBytes) {
        rejected = true;
        reject(new Error("Payload Too Large"));
        return;
      }

      body += chunk;
    });

    req.on("end", () => {
      if (rejected) return;

      try {
        resolve(body ? JSON.parse(body) : {});
      } catch {
        reject(new Error("Invalid JSON"));
      }
    });

    req.on("error", (error) => {
      if (!rejected) reject(error);
    });
  });
};

export default async function migrationHandler(
  req: NextApiRequest,
  res: NextApiResponse
) {
  const user = await verifyUser({ req, res });
  if (!user) return;

  if (req.method === "GET") {
    try {
      const data = await exportData(user.id);

      if (data.status !== 200) {
        return res.status(data.status).json({ response: data.response });
      }

      return res
        .setHeader("Content-Type", "application/json; charset=utf-8")
        .setHeader(
          "Content-Disposition",
          `attachment; filename="${getExportFilename()}"`
        )
        .setHeader("Cache-Control", "private, no-store, max-age=0")
        .setHeader("Pragma", "no-cache")
        .setHeader("X-Content-Type-Options", "nosniff")
        .status(200)
        .json(data.response);
    } catch {
      return res
        .status(500)
        .json({ response: "Unable to export bookmark data." });
    }
  }

  if (req.method === "POST") {
    if (process.env.NEXT_PUBLIC_DEMO === "true") {
      return res.status(400).json({
        response:
          "This action is disabled because this GoreeCloud Bookmarks demo is read-only.",
      });
    }

    const contentType = req.headers["content-type"] || "";
    if (!contentType.toLowerCase().startsWith("application/json")) {
      return res.status(415).json({
        response: "Import requests must use application/json.",
      });
    }

    let parsedRequest: unknown;
    const limitMb = getImportLimitMb(process.env.IMPORT_LIMIT);

    try {
      parsedRequest = await parseJsonStream(req, limitMb);
    } catch (error: any) {
      if (error.message === "Payload Too Large") {
        return res.status(413).json({
          response: `Import file exceeds the ${limitMb}MB size limit.`,
        });
      }

      return res
        .status(400)
        .json({ response: "Invalid request body provided." });
    }

    if (
      !parsedRequest ||
      typeof parsedRequest !== "object" ||
      !isSupportedMigrationFormat((parsedRequest as MigrationRequest).format)
    ) {
      return res.status(400).json({
        response: "Unsupported or missing migration format.",
      });
    }

    const request = parsedRequest as MigrationRequest;

    try {
      let data;
      if (request.format === MigrationFormat.htmlFile)
        data = await importFromHTMLFile(user.id, request.data);
      else if (request.format === MigrationFormat.linkwarden)
        data = await importFromLinkwarden(user.id, request.data);
      else if (request.format === MigrationFormat.wallabag)
        data = await importFromWallabag(user.id, request.data);
      else if (request.format === MigrationFormat.omnivore)
        data = await importFromOmnivore(user.id, request.data);
      else if (request.format === MigrationFormat.pocket)
        data = await importFromPocket(user.id, request.data);

      if (!data) {
        return res.status(400).json({
          response: "Unsupported migration format.",
        });
      }

      return res.status(data.status).json({ response: data.response });
    } catch {
      return res
        .status(500)
        .json({ response: "Unable to import bookmark data." });
    }
  }

  res.setHeader("Allow", ["GET", "POST"]);
  return res.status(405).json({ response: "Method not allowed." });
}
