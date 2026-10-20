import "server-only";
import { json } from "@/lib/auth";
import type { ZodType } from "zod";

export async function backend<T>(path: string, schema: ZodType<T>, pair?: string) {
  const configured = process.env.BACKEND_URL;
  if (!configured) return json({ error: "The backend connection is not configured." }, 503);
  let base: URL;
  try {
    base = new URL(configured);
    if (!["http:", "https:"].includes(base.protocol) || base.username || base.password || base.search || base.hash) throw new Error();
  } catch {
    return json({ error: "The backend connection is not configured correctly." }, 503);
  }
  const headers: Record<string, string> = { Accept: "application/json" };
  if (pair) {
    if (!process.env.BACKEND_API_TOKEN) return json({ error: "Demo submissions are not configured." }, 503);
    headers.Authorization = `Bearer ${process.env.BACKEND_API_TOKEN}`;
    headers["Content-Type"] = "application/json";
  }
  try {
    const response = await fetch(new URL(path, base), {
      method: pair ? "POST" : "GET",
      headers,
      ...(pair ? { body: JSON.stringify({ pair }) } : {}),
      cache: "no-store",
      signal: AbortSignal.timeout(15000),
      redirect: "error",
    });
    if (!response.ok) {
      const messages: Record<number, string> = {
        400: "The backend rejected this request.",
        401: "Backend authorization is not configured correctly. Contact the demo operator.",
        403: "The backend did not permit this request.",
        404: "This job is no longer available in the recent history.",
        429: "The submission limit has been reached. Please wait before trying again.",
        503: "The backend dependency is unavailable. Live data could not be retrieved.",
      };
      const status = response.status === 401 || response.status === 403 ? 502 : response.status;
      const result = json({ error: messages[response.status] ?? "The backend could not complete this request." }, status);
      const retryAfter = response.headers.get("Retry-After");
      if (retryAfter && /^\d+$/.test(retryAfter)) result.headers.set("Retry-After", retryAfter);
      return result;
    }
    const parsed = schema.safeParse(await response.json());
    if (!parsed.success) return json({ error: "The backend returned an unexpected response. No substitute data is displayed." }, 502);
    return json(parsed.data, pair ? 202 : 200);
  } catch (error) {
    const timedOut = error instanceof Error && (error.name === "TimeoutError" || error.name === "AbortError");
    return json({ error: timedOut
      ? (pair ? "The submission timed out; it may have been queued. Check recent jobs before submitting again." : "The backend timed out. Please try again.")
      : (pair ? "The backend connection was interrupted; the request may have been queued. Check recent jobs before trying again." : "Cannot reach the backend. Live data is unavailable.") }, timedOut ? 504 : 502);
  }
}
