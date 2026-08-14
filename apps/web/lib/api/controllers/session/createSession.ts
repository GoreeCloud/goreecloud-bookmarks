import { prisma } from "@linkwarden/prisma";
import crypto from "crypto";
import { decode, encode } from "next-auth/jwt";

export type SessionPurpose = "browser_extension";

const ONE_DAY_IN_SECONDS = 86400;
const DEFAULT_SESSION_DAYS = 73000;
const DEFAULT_SESSION_MAX_AGE_DAYS = 73050;
const BROWSER_EXTENSION_SESSION_DAYS = 30;

export default async function createSession(
  userId: number,
  sessionName?: string,
  purpose?: SessionPurpose
) {
  const now = Date.now();
  const expiryDate = new Date();
  const isBrowserExtensionSession = purpose === "browser_extension";
  const sessionDays = isBrowserExtensionSession
    ? BROWSER_EXTENSION_SESSION_DAYS
    : DEFAULT_SESSION_DAYS;
  const maxAgeDays = isBrowserExtensionSession
    ? BROWSER_EXTENSION_SESSION_DAYS
    : DEFAULT_SESSION_MAX_AGE_DAYS;

  expiryDate.setDate(expiryDate.getDate() + sessionDays);

  // A browser-extension installation has one stable, named session. Reconnecting
  // replaces older credentials for that installation before issuing a new one.
  if (isBrowserExtensionSession && sessionName) {
    await prisma.accessToken.updateMany({
      where: {
        userId,
        name: sessionName,
        isSession: true,
        revoked: false,
      },
      data: {
        revoked: true,
      },
    });
  }

  const token = await encode({
    token: {
      id: userId,
      iat: now / 1000,
      exp: expiryDate.getTime() / 1000,
      jti: crypto.randomUUID(),
      ...(purpose ? { purpose } : {}),
    },
    maxAge: maxAgeDays * ONE_DAY_IN_SECONDS,
    secret: process.env.NEXTAUTH_SECRET as string,
  });

  const tokenBody = await decode({
    token,
    secret: process.env.NEXTAUTH_SECRET as string,
  });

  await prisma.accessToken.create({
    data: {
      name: sessionName || "Unknown Device",
      userId,
      token: tokenBody?.jti as string,
      isSession: true,
      expires: expiryDate,
    },
  });

  return {
    response: {
      token,
      expires: expiryDate.toISOString(),
      purpose: purpose || "general",
    },
    status: 200,
  };
}
