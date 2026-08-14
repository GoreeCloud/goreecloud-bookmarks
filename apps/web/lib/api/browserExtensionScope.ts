import { NextApiRequest } from "next";

const STATIC_ROUTES = new Set([
  "GET /api/v1/config",
  "GET /api/v1/collections",
  "GET /api/v1/tags",
  "GET /api/v1/search",
  "POST /api/v1/links",
  "DELETE /api/v1/session",
]);

const DYNAMIC_ROUTES = [
  {
    method: "PUT",
    pathname: /^\/api\/v1\/links\/\d+$/,
  },
  {
    method: "DELETE",
    pathname: /^\/api\/v1\/links\/\d+$/,
  },
  {
    method: "POST",
    pathname: /^\/api\/v1\/archives\/\d+$/,
  },
];

export function isBrowserExtensionRequestAllowed(req: NextApiRequest) {
  const method = req.method?.toUpperCase();
  const pathname = req.url?.split("?", 1)[0];

  if (!method || !pathname) return false;

  if (STATIC_ROUTES.has(`${method} ${pathname}`)) return true;

  return DYNAMIC_ROUTES.some(
    (route) => route.method === method && route.pathname.test(pathname)
  );
}
