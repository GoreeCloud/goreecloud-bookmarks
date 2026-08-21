import type { NextApiRequest, NextApiResponse } from "next";
import verifyUser from "@/lib/api/verifyUser";
import postLink from "@/lib/api/controllers/links/postLink";

type ExtensionCaptureRequest = {
  title?: unknown;
  url?: unknown;
  collectionId?: unknown;
  tags?: unknown;
  note?: unknown;
};

function cleanString(value: unknown, maxLength: number): string | null {
  if (typeof value !== "string") return null;
  const cleaned = value.trim();
  if (!cleaned) return null;
  return cleaned.slice(0, maxLength);
}

function cleanCollectionId(value: unknown): number | undefined {
  if (value === null || value === undefined || value === "") return undefined;
  const parsed = typeof value === "number" ? value : Number(value);
  if (!Number.isSafeInteger(parsed) || parsed <= 0) return undefined;
  return parsed;
}

function cleanTags(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  const normalized = value
    .filter((tag): tag is string => typeof tag === "string")
    .map((tag) => tag.trim())
    .filter(Boolean);
  return Array.from(new Set(normalized)).slice(0, 30);
}

export default async function extensionCapture(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method !== "POST") {
    res.setHeader("Allow", "POST");
    return res.status(405).json({ response: "Method not allowed." });
  }

  const user = await verifyUser({ req, res });
  if (!user) return;

  if (process.env.NEXT_PUBLIC_DEMO === "true") {
    return res.status(400).json({
      response: "Bookmark capture is disabled in the read-only demo.",
    });
  }

  const body = (req.body ?? {}) as ExtensionCaptureRequest;
  const url = cleanString(body.url, 8192);
  const title = cleanString(body.title, 1000);
  const collectionId = cleanCollectionId(body.collectionId);
  const note = cleanString(body.note, 10000);
  const tags = cleanTags(body.tags);

  if (!url) {
    return res.status(400).json({ response: "A bookmark URL is required." });
  }

  if (
    body.collectionId !== null &&
    body.collectionId !== undefined &&
    body.collectionId !== "" &&
    collectionId === undefined
  ) {
    return res.status(400).json({ response: "The collection ID is invalid." });
  }

  let parsedUrl: URL;
  try {
    parsedUrl = new URL(url);
  } catch {
    return res.status(400).json({ response: "The bookmark URL is invalid." });
  }

  if (!["http:", "https:"].includes(parsedUrl.protocol)) {
    return res.status(400).json({
      response: "Only HTTP and HTTPS bookmark URLs are supported.",
    });
  }

  const created = await postLink(
    {
      url: parsedUrl.toString(),
      name: title || parsedUrl.toString(),
      description: note || undefined,
      collection: collectionId ? { id: collectionId } : undefined,
      tags: tags.map((name) => ({ name })),
    },
    user.id
  );

  return res.status(created.status).json({ response: created.response });
}
