import { NextApiRequest } from "next";
import { JWT } from "next-auth/jwt";
import { prisma } from "@linkwarden/prisma";
import getTokenFromRequest from "./getTokenFromRequest";
import { isBrowserExtensionRequestAllowed } from "./browserExtensionScope";

type Props = {
  req: NextApiRequest;
};

export default async function verifyToken({
  req,
}: Props): Promise<JWT | string> {
  const token = await getTokenFromRequest(req);
  const userId = token?.id;

  if (!userId) {
    return "You must be logged in.";
  }

  if (token.exp < Date.now() / 1000) {
    return "Your session has expired, please log in again.";
  }

  // Browser-extension sessions are intentionally narrower than normal sessions
  // and manually created API tokens. Deny everything not required by the
  // approved extension workflow.
  if (
    token.purpose === "browser_extension" &&
    !isBrowserExtensionRequestAllowed(req)
  ) {
    return "This browser extension session is not authorized for this action.";
  }

  // check if token is revoked
  const revoked = await prisma.accessToken.findFirst({
    where: {
      token: token.jti,
      revoked: true,
    },
  });

  if (revoked) {
    return "Your session has expired, please log in again.";
  }

  return token;
}
