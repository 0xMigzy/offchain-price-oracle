import { z } from "zod";

export const pairSchema = z.enum(["ETH/USD", "BTC/USD"]);
export type Pair = z.infer<typeof pairSchema>;
const timestamp = z.string().refine((value) => Number.isFinite(Date.parse(value)), "Invalid timestamp");
const statusSchema = z.enum(["queued", "processing", "cached", "submitted", "confirmed", "failed"]);
export const jobSchema = z.object({
  id: z.string().min(1),
  pair: pairSchema,
  status: statusSchema,
  createdAt: timestamp,
  updatedAt: timestamp,
  price: z.string().regex(/^\d+$/).optional(),
  priceDecimals: z.literal(8),
  txHash: z.string().regex(/^0x[a-fA-F0-9]{64}$/).optional(),
  blockNumber: z.number().int().nonnegative().optional(),
  error: z.string().optional(),
  events: z.array(z.object({ status: z.string(), at: timestamp })),
});
export type Job = z.infer<typeof jobSchema>;
export const jobsSchema = z.object({ jobs: z.array(jobSchema) });
const healthSchema = z.enum(["healthy", "unavailable", "unknown"]);
export const systemSchema = z.object({
  chainId: z.union([z.literal(11155111), z.literal(31337)]),
  contractAddress: z.string().refine((value) => value === "" || /^0x[a-fA-F0-9]{40}$/.test(value)),
  explorerBaseUrl: z.union([z.literal("https://sepolia.etherscan.io"), z.literal("")]),
  simulated: z.literal(true),
  queueBackend: z.literal("redis-streams"),
  services: z.object({ rpc: healthSchema, redis: healthSchema, queue: healthSchema }),
  pairs: z.array(pairSchema),
});
export type System = z.infer<typeof systemSchema>;
export const priceSchema = z.object({
  pair: pairSchema,
  price: z.string().regex(/^\d+$/),
  timestamp: z.number().int().nonnegative(),
  priceDecimals: z.literal(8),
});
export type Price = z.infer<typeof priceSchema>;
export const sessionSchema = z.object({ authenticated: z.boolean(), expiresAt: z.number().nullable(), available: z.boolean() });
export type Session = z.infer<typeof sessionSchema>;
export const acceptedSchema = z.object({ id: z.string().min(1), pair: pairSchema, status: z.literal("queued") });

export function formatPrice(raw: string, decimals = 8): string {
  const padded = raw.padStart(decimals + 1, "0");
  const whole = padded.slice(0, -decimals).replace(/\B(?=(\d{3})+(?!\d))/g, ",");
  const fraction = padded.slice(-decimals).replace(/0+$/, "").padEnd(2, "0");
  return `$${whole}.${fraction}`;
}
