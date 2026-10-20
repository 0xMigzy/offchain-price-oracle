# Off-Chain Price Oracle

An educational oracle dashboard with real queue processing and on-chain receipts.
ETH/USD and BTC/USD prices are **simulated fixed values**, not market data.

## Request flow

Next.js dashboard → authenticated Go API → Redis Streams → single Go worker →
Redis price cache or `OracleConsumer` transaction. Railway hosts the Go service
and private Redis; Vercel hosts the dashboard. Hosted deployments require Sepolia
(chain ID 11155111). Local Anvil uses chain ID 31337 and is labeled separately.

Redis Streams replaces Kafka to reduce the hosted demo's idle resource costs.
The same Redis instance holds the queue, recent job records, and short-lived price
cache. The local Compose stack no longer starts a Kafka broker.

Jobs expose recorded `queued`, `processing`, `cached`, `submitted`, `confirmed`,
and `failed` events. `cached` means no new transaction. `confirmed` means a
successful mined receipt, **not Ethereum finality**. The dashboard never invents
progress or substitutes mock data for unavailable services.

## Local development

1. Configure root `.env` with the local deployer `PRIVATE_KEY` and
   `CONTRACT_ADDRESS`. Never use the publicly known Anvil key on a public network.
2. From the repository root, run `scripts/deploy-local.sh` with `PRIVATE_KEY`
   exported. It starts Anvil/Redis, deploys the contract, writes root `.env`,
   and attaches to the worker.
3. After source changes, run `docker compose up -d --build worker`.
4. In `frontend`, copy `.env.example` to `.env.local` and configure the four
   server-only variables listed below, with `BACKEND_URL=http://127.0.0.1:8080`.
5. Run `npm ci && npm run dev` in `frontend`, then open `http://localhost:3000`.

Compose explicitly enables unauthenticated **local** submissions when no
`BACKEND_API_TOKEN` is set. Set that token in the worker environment to match the
frontend to exercise authenticated submission locally. Railway forbids this
local opt-out. The frontend always requires its demo password for submissions.

`/health` remains a plain `ok` liveness endpoint, not a dependency readiness check.
Public API reads: `/api/system`, `/api/jobs`, `/api/jobs/{id}`, and
`/api/prices?pair=ETH%2FUSD`. `POST /oracle/request` accepts `{"pair":"ETH/USD"}`
and returns 202 only after Redis accepts the job. Hosted writes require
`Authorization: Bearer <BACKEND_API_TOKEN>`.

## Hosting configuration

| Component | Location |
| --- | --- |
| Dashboard (Vercel `0x-migzy/oracle-lab-demo`) | https://oracle-lab-demo.vercel.app |
| API/worker (Railway `oracle-lab-demo`, service `api`) | https://api-production-4fb0.up.railway.app |
| Redis (Railway service `redis`, volume `redis-volume` at `/data`) | Private network only: `redis.railway.internal:6379` |
| `OracleConsumer` on Sepolia | [`0x5178A4Ed5A008feE6153c85815445C1d7d76c4d7`](https://sepolia.etherscan.io/address/0x5178A4Ed5A008feE6153c85815445C1d7d76c4d7) |
| Owner and only submitter | `0x6561961228A15988d946D9017510e813377cE473` (testnet-only demo signer) |

Both Railway services run one replica with `overlapSeconds=0`, so a redeploy
never runs two workers against the same signer.

To redeploy the contract, use `forge create src/OracleConsumer.sol:OracleConsumer`
with `--gas-limit` set from the RPC provider's `eth_estimateGas` plus a margin
(about 8.7M gas). `forge script` under-estimated it at about 1.6M, and that
deployment ran out of gas on chain.

| Service | Root/config | Required variables |
| --- | --- | --- |
| Vercel dashboard | `frontend` | `BACKEND_URL`, `BACKEND_API_TOKEN`, `DEMO_PASSWORD`, `SESSION_SECRET` (at least 32 characters) |
| Railway API/worker | `worker`, `worker/railway.toml` | `RPC_URL`, `CONTRACT_ADDRESS`, `PRIVATE_KEY`, `BACKEND_API_TOKEN`, `REDIS_URL`; platform `PORT` supported |
| Railway Redis | `worker/deploy/redis` | `REDIS_PASSWORD`; optional `REDIS_MAXMEMORY` (default `64mb`) |

The signing key belongs only in Railway secret configuration. It must never be
placed in Vercel, browser code, or a `NEXT_PUBLIC_*` variable. Store local secrets
under ignored `.secrets/` with owner-only permissions. Frontend `.env.local` is
also ignored.

Redis must have a persistent volume at `/data`, AOF enabled, and no eviction.
Keep it private; only the Go HTTP service needs a public domain. Its configured
memory limit is not a limit on total Redis RSS or billed usage.

Run **one worker only**, including during deployment replacement: it owns the
consumer identity and signer nonce stream. Signed transactions are stored
privately before broadcast so interrupted work can resume using the same hash.
A stuck pending transaction blocks subsequent jobs; fee replacement is not
implemented. Losing Redis data removes this recovery protection.

Terminal job records expire after seven days; the recent index is bounded to
1,000 IDs and public history returns up to 30 jobs. Active jobs remain until
resolved. The shared-password session lasts four hours. Login throttling is
per Vercel instance, not a distributed abuse-prevention guarantee.

The Railway trial lasts at most 30 days or until its credit is spent. Measure
actual RAM/CPU/storage usage after deployment before projecting availability;
the two-service design does not guarantee that $5 covers a particular duration.

## Verification

```sh
cd worker
go vet ./...
go test ./...
go build ./...
```

The cache test uses local Redis. Durable queue tests require `TEST_REDIS_URL`
pointing to an isolated Redis database; they modify their queue/test keys:

```sh
TEST_REDIS_URL=redis://localhost:6379/15 go test -v ./internal/queue
```

Frontend: `npm run typecheck && npm run build` inside `frontend`.
Contracts: `forge test -vv` inside `contracts`.

For end-to-end proof, unlock the dashboard, submit a pair, inspect its recorded
events, and verify the successful receipt and resulting on-chain price. Repeat
within the cache TTL to observe `cached` without another transaction. Public
deployment verification must use the actual Vercel and Railway URLs and a
Sepolia receipt; a local Anvil smoke run is not equivalent.
