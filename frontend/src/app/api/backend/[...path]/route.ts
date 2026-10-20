import { NextRequest } from "next/server";
import { backend } from "@/lib/backend";
import { jobSchema, jobsSchema, pairSchema, priceSchema, systemSchema } from "@/lib/contracts";
import { json } from "@/lib/auth";

export const runtime = "nodejs";
export const dynamic = "force-dynamic";

export async function GET(request: NextRequest, context: { params: Promise<{ path: string[] }> }) {
  const { path } = await context.params;
  if (path.length === 1 && path[0] === "system") return backend("/api/system", systemSchema);
  if (path.length === 1 && path[0] === "jobs") return backend("/api/jobs", jobsSchema);
  if (path.length === 2 && path[0] === "jobs" && /^[a-zA-Z0-9_-]{1,128}$/.test(path[1])) {
    return backend(`/api/jobs/${encodeURIComponent(path[1])}`, jobSchema);
  }
  if (path.length === 1 && path[0] === "prices") {
    const pair = pairSchema.safeParse(request.nextUrl.searchParams.get("pair"));
    if (!pair.success) return json({ error: "Choose ETH/USD or BTC/USD." }, 400);
    return backend(`/api/prices?pair=${encodeURIComponent(pair.data)}`, priceSchema);
  }
  return json({ error: "This endpoint does not exist." }, 404);
}
