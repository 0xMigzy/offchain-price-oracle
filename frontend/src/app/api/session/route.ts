import { NextRequest } from "next/server";
import { authConfigured, clearAttempts, endSession, json, loginAttempt, passwordMatches, readJson, rememberAttempt, sameOrigin, session, startSession } from "@/lib/auth";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

export async function GET(request: NextRequest) {
  const current = session(request);
  return json({ authenticated: Boolean(current), expiresAt: current?.expiresAt ?? null, available: authConfigured() });
}

export async function POST(request: NextRequest) {
  if (!sameOrigin(request)) return json({ error: "This request must come from the dashboard." }, 403);
  if (!authConfigured()) return json({ error: "Demo access is not configured. Contact the demo operator." }, 503);
  const limit = loginAttempt(request);
  if (!limit.allowed) {
    const retryAfter = Math.max(1, Math.ceil((limit.attempt.expiresAt - Date.now()) / 1000));
    const response = json({ error: "Too many unlock attempts. Please try again after the cooldown.", retryAfter }, 429);
    response.headers.set("Retry-After", String(retryAfter));
    rememberAttempt(response, limit.attempt);
    return response;
  }
  let password: unknown;
  try {
    const body = await readJson(request);
    password = body && typeof body === "object" && "password" in body ? body.password : undefined;
  } catch {
    const response = json({ error: "Enter a password using a valid JSON request." }, 400);
    rememberAttempt(response, limit.attempt);
    return response;
  }
  if (typeof password !== "string" || password.length > 1024 || !passwordMatches(password)) {
    const response = json({ error: "That demo password is incorrect." }, 401);
    rememberAttempt(response, limit.attempt);
    return response;
  }
  clearAttempts(limit.key);
  const expiresAt = Date.now() + 4 * 60 * 60 * 1000;
  const response = json({ authenticated: true, expiresAt, available: true });
  startSession(response, expiresAt);
  return response;
}

export async function DELETE(request: NextRequest) {
  if (!sameOrigin(request)) return json({ error: "This request must come from the dashboard." }, 403);
  const response = json({ authenticated: false, expiresAt: null, available: authConfigured() });
  endSession(response);
  return response;
}
