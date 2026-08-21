import type { NextApiRequest, NextApiResponse } from "next";
import verifyUser from "@/lib/api/verifyUser";
import postLink from "@/lib/api/controllers/links/postLink";

type ExtensionCaptureRequest = {
  title?: unknown;
  url?: unknown;
  collection?: unknown;
  tags?: unknown;
  note?: unknown;
};

function cleanString(value: unknown, maxLength: number): string | null {
  if (typeof value !== "string") return null;
  const cleaned = value.trim();
  if (!cleaned) return null;
  return cleaned.slice(0, maxLength);
}

function cleanTags(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  return [
    ...new Set(
      value
        .filter((tag): tag is string => typeof tag === "string")
        .map((tag) => tag.trim())
        .filter(Boolean)
    ),
  ].slice(0, 30);
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
  const collection = cleanString(body.collection, 255);
  const note = cleanString(body.note, 10000);
  const tags = cleanTags(body.tags);

  if (!url) {
    return res.status(400).json({ response: "A bookmark URL is required." });
  }

  let parsedUrl: URL;
  try {
    parsedUrl = new URL(url);
  } catch {
    return res.status(400).json({ response: "The bookmark URL is invalid." });
  }

  if (!['http:', 'https:'].includes(parsedUrl.protocol)) {
    return res.status(400).json({
      response: "Only HTTP and HTTPS bookmark URLs are supported.",
    });
  }

  const created = await postLink(
    {
      url: parsedUrl.toString(),
      name: title || parsedUrl.toString(),
      description: note || undefined,
      collection: collection ? { name: collection } : undefined,
      tags: tags.map((name) => ({ name })),
    },
    user.id
  );

  return res.status(created.status).json({ response: created.response });
}
