import type { NextApiRequest, NextApiResponse } from "next";
import { prisma } from "@linkwarden/prisma";
import verifyByCredentials from "@/lib/api/verifyByCredentials";
import verifyToken from "@/lib/api/verifyToken";
import createSession, {
  SessionPurpose,
} from "@/lib/api/controllers/session/createSession";
import { PostSessionSchema } from "@linkwarden/lib/schemaValidation";

export default async function session(
  req: NextApiRequest,
  res: NextApiResponse
) {
  if (req.method === "DELETE") {
    const token = await verifyToken({ req });

    if (typeof token === "string") {
      return res.status(401).json({ response: token });
    }

    if (!token.jti || !token.id) {
      return res.status(400).json({ response: "Invalid session token." });
    }

    const revoked = await prisma.accessToken.updateMany({
      where: {
        token: token.jti,
        userId: Number(token.id),
        isSession: true,
        revoked: false,
      },
      data: {
        revoked: true,
      },
    });

    return res.status(200).json({
      response: {
        revoked: revoked.count > 0,
      },
    });
  }

  if (req.method !== "POST") {
    res.setHeader("Allow", ["POST", "DELETE"]);
    return res.status(405).json({ response: "Method not allowed." });
  }

  const dataValidation = PostSessionSchema.safeParse(req.body);

  if (!dataValidation.success) {
    return res.status(400).json({
      response: `Error: ${
        dataValidation.error.issues[0].message
      } [${dataValidation.error.issues[0].path.join(", ")}]`,
    });
  }

  const requestedPurpose = req.body?.purpose;
  if (
    requestedPurpose !== undefined &&
    requestedPurpose !== "browser_extension"
  ) {
    return res.status(400).json({ response: "Invalid session purpose." });
  }

  const purpose = requestedPurpose as SessionPurpose | undefined;
  const { username, password, sessionName } = dataValidation.data;

  const result = await verifyByCredentials({ username, password });

  if (result.status === "email_not_verified")
    return res.status(401).json({
      response: "Please verify your email address before logging in.",
      code: "EMAIL_NOT_VERIFIED",
      email: result.user.email,
    });

  if (result.status !== "success")
    return res.status(400).json({
      response:
        "Invalid credentials. You might need to reset your password if you're sure you already signed up with the current username/email.",
    });

  const token = await createSession(result.user.id, sessionName, purpose);
  return res.status(token.status).json({ response: token.response });
}
