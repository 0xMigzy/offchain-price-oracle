import { NextRequest } from "next/server";
import { backend } from "@/lib/backend";
import { acceptedSchema, pairSchema } from "@/lib/contracts";
import { authConfigured, endSession, json, readJson, sameOrigin, session } from "@/lib/auth";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

export async function POST(request: NextRequest) {
  if (!sameOrigin(request)) return json({ error: "This request must come from the dashboard." }, 403);
  if (!authConfigured()) return json({ error: "Demo access is not configured." }, 503);
  if (!session(request)) {
    const response = json({ error: "Your demo session has expired. Unlock submissions again." }, 401);
    endSession(response);
    return response;
  }
  try {
    const body = await readJson(request);
    const pair = pairSchema.safeParse(body && typeof body === "object" && "pair" in body ? body.pair : undefined);
    if (!pair.success) return json({ error: "Choose ETH/USD or BTC/USD." }, 400);
    return backend("/oracle/request", acceptedSchema, pair.data);
  } catch {
    return json({ error: "Send a valid JSON request with a supported pair." }, 400);
  }
}
