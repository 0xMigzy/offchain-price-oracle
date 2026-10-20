import "server-only";
import { createHash, createHmac, timingSafeEqual } from "node:crypto";
import { NextRequest, NextResponse } from "next/server";

const production = process.env.NODE_ENV === "production";
const sessionName = production ? "__Host-oracle-session" : "oracle-session";
const attemptName = production ? "__Host-oracle-attempts" : "oracle-attempts";
const lifetime = 4 * 60 * 60;
const attemptWindow = 15 * 60;
const maxAttempts = 5;
const cookieOptions = { httpOnly: true, secure: production, sameSite: "strict" as const, path: "/" };
type Attempt = { count: number; expiresAt: number };
const attempts = new Map<string, Attempt>();

type SignedValue = { purpose: "session" | "attempt"; expiresAt: number; count?: number };

export function authConfigured(): boolean {
  return Boolean(process.env.DEMO_PASSWORD && (process.env.SESSION_SECRET?.length ?? 0) >= 32);
}

function signature(value: string): Buffer {
  // Including the password digest invalidates all sessions when the demo password rotates.
  return createHmac("sha256", process.env.SESSION_SECRET!)
    .update(createHash("sha256").update(process.env.DEMO_PASSWORD!).digest())
    .update(value)
    .digest();
}

function sign(value: SignedValue): string {
  const payload = Buffer.from(JSON.stringify(value)).toString("base64url");
  return `${payload}.${signature(payload).toString("base64url")}`;
}

function verify(raw: string | undefined, purpose: SignedValue["purpose"]): SignedValue | null {
  if (!authConfigured() || !raw || raw.length > 1024) return null;
  const parts = raw.split(".");
  if (parts.length !== 2) return null;
  try {
    const supplied = Buffer.from(parts[1], "base64url");
    const expected = signature(parts[0]);
    if (supplied.length !== expected.length || !timingSafeEqual(supplied, expected)) return null;
    const value: SignedValue = JSON.parse(Buffer.from(parts[0], "base64url").toString());
    if (value.purpose !== purpose || !Number.isSafeInteger(value.expiresAt) || value.expiresAt <= Date.now()) return null;
    return value;
  } catch {
    return null;
  }
}

export function session(request: NextRequest): SignedValue | null {
  return verify(request.cookies.get(sessionName)?.value, "session");
}

export function json(body: unknown, status = 200): NextResponse {
  return NextResponse.json(body, { status, headers: { "Cache-Control": "no-store, private", Vary: "Cookie" } });
}

export function sameOrigin(request: NextRequest): boolean {
  const origin = request.headers.get("origin");
  if (!origin) return false;
  try {
    const host = request.headers.get("host");
    if (!host) return false;
    // NextURL canonicalizes loopback names; compare the browser's actual Host.
    const expected = new URL(`${request.nextUrl.protocol}//${host}`);
    return new URL(origin).origin === expected.origin;
  } catch {
    return false;
  }
}

export function passwordMatches(password: string): boolean {
  const supplied = createHash("sha256").update(password).digest();
  const expected = createHash("sha256").update(process.env.DEMO_PASSWORD!).digest();
  return timingSafeEqual(supplied, expected);
}

export function startSession(response: NextResponse, expiresAt: number): void {
  response.cookies.set(sessionName, sign({ purpose: "session", expiresAt }), { ...cookieOptions, maxAge: lifetime });
  response.cookies.set(attemptName, "", { ...cookieOptions, maxAge: 0 });
}

export function endSession(response: NextResponse): void {
  response.cookies.set(sessionName, "", { ...cookieOptions, maxAge: 0 });
}

function clientKey(request: NextRequest): string {
  // Vercel overwrites this header. Never trust arbitrary client X-Forwarded-For.
  return process.env.VERCEL === "1"
    ? (request.headers.get("x-vercel-forwarded-for")?.split(",")[0]?.trim() || "unknown")
    : "local";
}

export function loginAttempt(request: NextRequest): { allowed: boolean; attempt: Attempt; key: string } {
  const now = Date.now();
  for (const [key, value] of attempts) if (value.expiresAt <= now) attempts.delete(key);
  const key = clientKey(request);
  const stored = attempts.get(key);
  const cookie = verify(request.cookies.get(attemptName)?.value, "attempt");
  const remembered = cookie && Number.isSafeInteger(cookie.count) && cookie.count! > 0
    ? { count: cookie.count!, expiresAt: cookie.expiresAt }
    : null;
  const previous = stored && (!remembered || stored.count >= remembered.count) ? stored : remembered;
  const attempt = previous ?? { count: 0, expiresAt: now + attemptWindow * 1000 };
  if (attempt.count >= maxAttempts) return { allowed: false, attempt, key };
  const next = { ...attempt, count: attempt.count + 1 };
  if (attempts.size >= 10000 && !attempts.has(key)) {
    const oldest = attempts.keys().next().value;
    if (oldest) attempts.delete(oldest);
  }
  attempts.set(key, next);
  return { allowed: true, attempt: next, key };
}

export function rememberAttempt(response: NextResponse, attempt: Attempt): void {
  response.cookies.set(attemptName, sign({ purpose: "attempt", ...attempt }), {
    ...cookieOptions,
    maxAge: Math.max(1, Math.ceil((attempt.expiresAt - Date.now()) / 1000)),
  });
}

export function clearAttempts(key: string): void {
  attempts.delete(key);
}

export async function readJson(request: NextRequest): Promise<unknown> {
  if (!request.headers.get("content-type")?.toLowerCase().startsWith("application/json")) {
    throw new Error("Send a JSON request.");
  }
  // Bound streamed bodies too; Content-Length is not required or trusted.
  const reader = request.body?.getReader();
  if (!reader) throw new Error("A request body is required.");
  let text = "";
  let size = 0;
  const decoder = new TextDecoder();
  try {
    while (true) {
      const { value, done } = await reader.read();
      if (done) break;
      size += value.byteLength;
      if (size > 4096) {
        await reader.cancel();
        throw new Error("Request body is too large.");
      }
      text += decoder.decode(value, { stream: true });
    }
    text += decoder.decode();
    return JSON.parse(text);
  } finally {
    reader.releaseLock();
  }
}
