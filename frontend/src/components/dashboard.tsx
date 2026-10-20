"use client";

import { useCallback, useEffect, useRef, useState, type FormEvent } from "react";
import { acceptedSchema, formatPrice, jobSchema, jobsSchema, priceSchema, sessionSchema, systemSchema, type Job, type Pair, type Price, type Session, type System } from "@/lib/contracts";
import type { ZodType } from "zod";

type Resource<T> = { data: T | null; error: string | null; loading: boolean; updatedAt: number | null };
const empty = { data: null, error: null, loading: true, updatedAt: null };
const descriptions: Record<Job["status"], string> = {
  queued: "The API accepted this request and placed it on the Redis Streams queue.",
  processing: "The Go worker is processing the request and checking the price cache.",
  cached: "A cached price was reused. This request did not create a new transaction.",
  submitted: "A transaction was broadcast to the connected chain. This is not yet a successful confirmation.",
  confirmed: "The worker observed a successful transaction receipt on the connected chain.",
  failed: "The request failed. See the recorded error below; no success is assumed.",
};

class ApiError extends Error {
  constructor(message: string, readonly status: number, readonly retryAfter?: number) {
    super(message);
  }
}

async function request<T>(path: string, schema: ZodType<T>, init?: RequestInit): Promise<T> {
  let response: Response;
  try {
    response = await fetch(path, { ...init, cache: "no-store", credentials: "same-origin", signal: init?.signal ?? AbortSignal.timeout(20000) });
  } catch (error) {
    if (init?.signal?.aborted) throw error;
    throw new ApiError(init?.method === "POST" && path === "/api/request"
      ? "The connection was interrupted. Your request may have reached the queue; check recent jobs before resubmitting."
      : "The dashboard connection failed. Check your connection and try again.", 0);
  }
  const body: unknown = await response.json().catch(() => null);
  if (!response.ok) {
    const payload = body && typeof body === "object" ? body as { error?: unknown; retryAfter?: unknown } : {};
    throw new ApiError(typeof payload.error === "string" ? payload.error : `The request failed (HTTP ${response.status}).`, response.status,
      typeof payload.retryAfter === "number" ? payload.retryAfter : undefined);
  }
  const parsed = schema.safeParse(body);
  if (!parsed.success) throw new ApiError("Unexpected server response. No substitute data is displayed.", 502);
  return parsed.data;
}

function localTime(value: number | string): string {
  return new Date(value).toLocaleTimeString([], { hour: "2-digit", minute: "2-digit", second: "2-digit" });
}

function Icon({ name, size = 18 }: { name: "arrow" | "lock" | "refresh" | "external" | "bolt"; size?: number }) {
  const paths: Record<typeof name, React.ReactNode> = {
    arrow: <><path d="M4 12h16m-6-6 6 6-6 6" /></>,
    lock: <><rect x="5" y="10" width="14" height="11" rx="2" /><path d="M8 10V7a4 4 0 0 1 8 0v3m-4 5v2" /></>,
    refresh: <><path d="M20 7v5h-5M4 17v-5h5" /><path d="M6 7a7 7 0 0 1 12-1l2 3M4 15l2 3a7 7 0 0 0 12-1" /></>,
    external: <><path d="M14 3h7v7m0-7L10 14M10 3H5a2 2 0 0 0-2 2v14a2 2 0 0 0 2 2h14a2 2 0 0 0 2-2v-5" /></>,
    bolt: <path d="m13 2-9 12h7l-1 8 10-13h-7l1-7Z" />,
  };
  return <svg width={size} height={size} viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.6" strokeLinecap="round" strokeLinejoin="round" aria-hidden="true">{paths[name]}</svg>;
}

function Status({ value }: { value: string }) {
  return <span className={`status status-${value}`}><span className="status-dot" />{value}</span>;
}

function PriceCard({ pair, resource }: { pair: Pair; resource: Resource<Price> }) {
  const value = resource.data;
  const initialized = Boolean(value && value.timestamp > 0);
  return <article className="price-card">
    <div className="price-top"><div className={`coin coin-${pair === "ETH/USD" ? "eth" : "btc"}`}>{pair === "ETH/USD" ? "Ξ" : "₿"}</div><div><h3>{pair}</h3><p>{pair === "ETH/USD" ? "Ethereum" : "Bitcoin"} / US dollar</p></div><span className="tag">ON-CHAIN</span></div>
    <div className={`price-value ${!initialized ? "price-empty" : ""}`}>{resource.loading && !value ? "Reading contract…" : initialized ? formatPrice(value!.price) : resource.error ? "Unavailable" : "Not published yet"}</div>
    <div className="price-caption"><span className="mini-dot" />Simulated price · 8 decimal precision</div>
    {resource.error && <p className="inline-error" role="status">{value ? "Last successful read shown. " : ""}{resource.error}</p>}
    <div className="price-footer"><span>{initialized ? "Contract updated" : "Waiting for an on-chain update"}</span>{initialized && <time dateTime={new Date(value!.timestamp * 1000).toISOString()} title={new Date(value!.timestamp * 1000).toLocaleString()}>{localTime(value!.timestamp * 1000)}</time>}</div>
  </article>;
}

export default function Dashboard() {
  const [system, setSystem] = useState<Resource<System>>(empty);
  const [jobs, setJobs] = useState<Resource<Job[]>>(empty);
  const [eth, setEth] = useState<Resource<Price>>(empty);
  const [btc, setBtc] = useState<Resource<Price>>(empty);
  const [access, setAccess] = useState<Resource<Session>>(empty);
  const [refreshKey, setRefreshKey] = useState(0);
  const [refreshing, setRefreshing] = useState(false);
  const [lastChecked, setLastChecked] = useState<number | null>(null);
  const [selectedId, setSelectedId] = useState<string | null>(null);
  const [detail, setDetail] = useState<Resource<Job>>(empty);
  const [pair, setPair] = useState<Pair>("ETH/USD");
  const [password, setPassword] = useState("");
  const [authBusy, setAuthBusy] = useState(false);
  const [authError, setAuthError] = useState<string | null>(null);
  const [cooldown, setCooldown] = useState(0);
  const [submitting, setSubmitting] = useState(false);
  const [submitError, setSubmitError] = useState<string | null>(null);
  const [notice, setNotice] = useState<string | null>(null);
  const [clock, setClock] = useState(0);
  const mounted = useRef(true);
  const authRevision = useRef(0);
  const passwordInput = useRef<HTMLInputElement>(null);
  const refresh = useCallback(() => setRefreshKey((value) => value + 1), []);

  useEffect(() => {
    mounted.current = true;
    setClock(Date.now());
    const interval = window.setInterval(() => setClock(Date.now()), 1000);
    return () => { mounted.current = false; window.clearInterval(interval); };
  }, []);

  useEffect(() => {
    const controller = new AbortController();
    let timer: number | undefined;
    let running = false;
    async function load() {
      if (running || controller.signal.aborted) return;
      running = true;
      setRefreshing(true);
      const revision = authRevision.current;
      async function update<T>(path: string, schema: ZodType<T>, setter: React.Dispatch<React.SetStateAction<Resource<T>>>, isSession = false) {
        try {
          const data = await request(path, schema, { signal: AbortSignal.any([controller.signal, AbortSignal.timeout(20000)]) });
          if (!controller.signal.aborted && (!isSession || revision === authRevision.current)) setter({ data, error: null, loading: false, updatedAt: Date.now() });
        } catch (error) {
          if (!controller.signal.aborted && (!isSession || revision === authRevision.current)) setter((previous) => ({ ...previous, error: error instanceof Error ? error.message : "Unable to load live data.", loading: false }));
        }
      }
      await Promise.all([
        update("/api/backend/system", systemSchema, setSystem),
        update("/api/backend/jobs", jobsSchema.transform((value) => value.jobs), setJobs),
        update("/api/backend/prices?pair=ETH%2FUSD", priceSchema, setEth),
        update("/api/backend/prices?pair=BTC%2FUSD", priceSchema, setBtc),
        update("/api/session", sessionSchema, setAccess, true),
      ]);
      running = false;
      if (!controller.signal.aborted) {
        setRefreshing(false);
        setLastChecked(Date.now());
        timer = window.setTimeout(() => { if (!document.hidden) void load(); }, 5000);
      }
    }
    const onVisible = () => { if (!document.hidden) { clearTimeout(timer); void load(); } };
    void load();
    document.addEventListener("visibilitychange", onVisible);
    window.addEventListener("focus", onVisible);
    return () => { controller.abort(); clearTimeout(timer); document.removeEventListener("visibilitychange", onVisible); window.removeEventListener("focus", onVisible); };
  }, [refreshKey]);

  useEffect(() => {
    if (!selectedId) return;
    const controller = new AbortController();
    let timer: number | undefined;
    setDetail({ ...empty });
    async function load() {
      try {
        const data = await request(`/api/backend/jobs/${encodeURIComponent(selectedId!)}`, jobSchema, { signal: AbortSignal.any([controller.signal, AbortSignal.timeout(20000)]) });
        if (!controller.signal.aborted) setDetail({ data, loading: false, error: null, updatedAt: Date.now() });
      } catch (error) {
        if (!controller.signal.aborted) setDetail((previous) => ({ ...previous, loading: false, error: error instanceof Error ? error.message : "Cannot load this job." }));
      }
      if (!controller.signal.aborted) timer = window.setTimeout(load, 5000);
    }
    void load();
    return () => { controller.abort(); clearTimeout(timer); };
  }, [selectedId, refreshKey]);

  const unlocked = Boolean(access.data?.authenticated && access.data.expiresAt && access.data.expiresAt > clock);
  const sessionExpired = Boolean(access.data?.authenticated && !unlocked);
  const cooldownSeconds = Math.max(0, Math.ceil((cooldown - clock) / 1000));
  const activeJob = selectedId
    ? (detail.data?.id === selectedId ? detail.data : jobs.data?.find((job) => job.id === selectedId))
    : jobs.data?.[0];
  const pending = jobs.data?.filter((job) => ["queued", "processing", "submitted"].includes(job.status)).length ?? 0;
  const contractAddress = system.data?.contractAddress;
  const systemLive = Boolean(system.data && !system.error);
  const chainId = system.data?.chainId;
  const networkName = chainId === 31337 ? "Local Anvil" : chainId === 11155111 ? "Sepolia testnet" : "Network unverified";
  const explorer = chainId === 11155111 ? system.data!.explorerBaseUrl : "";

  async function unlock(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (authBusy || cooldownSeconds > 0) return;
    setAuthBusy(true);
    setAuthError(null);
    authRevision.current += 1;
    try {
      const data = await request("/api/session", sessionSchema, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ password }) });
      if (!mounted.current) return;
      authRevision.current += 1;
      setAccess({ data, error: null, loading: false, updatedAt: Date.now() });
      setPassword("");
      setNotice("Submissions unlocked for this browser for four hours.");
      refresh();
    } catch (error) {
      if (!mounted.current) return;
      setAuthError(error instanceof Error ? error.message : "Unable to unlock submissions.");
      if (error instanceof ApiError && error.retryAfter) setCooldown(Date.now() + error.retryAfter * 1000);
      setPassword("");
      passwordInput.current?.focus();
    } finally {
      if (mounted.current) setAuthBusy(false);
    }
  }

  async function lock() {
    if (authBusy) return;
    setAuthBusy(true);
    setAuthError(null);
    authRevision.current += 1;
    try {
      const data = await request("/api/session", sessionSchema, { method: "DELETE" });
      if (!mounted.current) return;
      authRevision.current += 1;
      setAccess({ data, error: null, loading: false, updatedAt: Date.now() });
      setNotice("Submissions locked. Public reads remain available.");
      refresh();
    } catch (error) {
      if (mounted.current) setAuthError(error instanceof Error ? error.message : "Could not end the session. Try again.");
    } finally {
      if (mounted.current) setAuthBusy(false);
    }
  }

  async function submit(event: FormEvent<HTMLFormElement>) {
    event.preventDefault();
    if (!unlocked || submitting) return;
    setSubmitting(true);
    setSubmitError(null);
    setNotice(null);
    try {
      const accepted = await request("/api/request", acceptedSchema, { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ pair }) });
      if (!mounted.current) return;
      setSelectedId(accepted.id);
      setNotice(`Request accepted for ${accepted.pair}. Follow its recorded events below.`);
      refresh();
    } catch (error) {
      if (!mounted.current) return;
      setSubmitError(error instanceof Error ? error.message : "The request could not be submitted.");
      if (error instanceof ApiError && error.status === 401) {
        authRevision.current += 1;
        setAccess({ data: { authenticated: false, expiresAt: null, available: true }, error: null, loading: false, updatedAt: Date.now() });
      }
      refresh();
    } finally {
      if (mounted.current) setSubmitting(false);
    }
  }

  return <>
    <a className="skip-link" href="#main">Skip to dashboard</a>
    <header className="site-header"><div className="header-inner"><a className="brand" href="/" aria-label="Oracle Lab home"><span className="brand-mark"><Icon name="bolt" size={21} /></span><span>ORACLE<span className="brand-light">LAB</span><small>DECENTRALIZED DATA PIPELINE</small></span></a><div className="header-right"><span className="network-pill"><span className="status-dot" />{networkName}<span className="network-id"> / {chainId ?? "—"}</span></span><a className="text-link desktop-link" href="#how-it-works">How it works <Icon name="arrow" size={14} /></a></div></div></header>
    <main id="main" className="dashboard">
      <section className="hero"><div><p className="eyebrow"><span className="accent-line" /> AN EDUCATIONAL ORACLE, END TO END</p><h1>Off-chain data.<br /><span>On-chain proof.</span></h1><p className="hero-copy">Follow a price request through a real event-driven pipeline.<br className="desktop-break" /> Simulated market data. Real infrastructure. Verifiable testnet transactions.</p></div><div className="hero-aside"><span className="hero-index">01 — LIVE OBSERVATORY</span><div className="hero-orbit" aria-hidden="true"><span className="orbit-inner"><Icon name="bolt" size={37} /></span><span className="orbit-label orbit-label-one">OFF-CHAIN</span><span className="orbit-label orbit-label-two">ON-CHAIN</span></div><span className="hero-footnote">NO WALLET REQUIRED</span></div></section>
      <div className="demo-banner"><span className="banner-icon">i</span><p><strong>Learning environment, not a market feed.</strong> Prices are simulated fixed values. Transactions use test ETH on {networkName}, not real funds. Cached requests do not send a new transaction.</p><span className="tag amber">SIMULATED DATA</span></div>

      <section className="panel pipeline-panel" id="how-it-works" aria-labelledby="pipeline-heading"><div className="section-header"><div><p className="eyebrow">THE ARCHITECTURE</p><h2 id="pipeline-heading">One request. Five stages.</h2></div><span className="muted mono small">ASYNCHRONOUS BY DESIGN</span></div><div className="pipeline">
        {[
          { number: "01", name: "API", tech: "Go · HTTP", text: "Validates & enqueues", health: systemLive ? "healthy" : "unknown" },
          { number: "02", name: "Redis Streams", tech: "Durable queue", text: "Buffers the request", health: systemLive ? system.data!.services.queue : "unknown" },
          { number: "03", name: "Worker", tech: "Go · Consumer", text: "Processes each job", health: "unobserved" },
          { number: "04", name: "Redis", tech: "Price cache", text: "Reuses fresh prices", health: systemLive ? system.data!.services.redis : "unknown" },
          { number: "05", name: chainId === 31337 ? "Anvil" : "Sepolia", tech: "Ethereum · RPC", text: "Records the update", health: systemLive ? system.data!.services.rpc : "unknown" },
        ].map((stage, index) => <div className="pipeline-stage" key={stage.name}><div className="stage-top"><span className="stage-number">{stage.number}</span><span className={`health-dot health-${stage.health}`} title={stage.health === "unobserved" ? "Worker has no independent health probe; observe job events." : `${stage.name}: ${stage.health}`} aria-label={`${stage.name}: ${stage.health}`} /></div><h3>{stage.name}</h3><span className="stage-tech">{stage.tech}</span><p>{stage.text}</p>{index < 4 && <span className="stage-arrow"><Icon name="arrow" size={17} /></span>}</div>)}
      </div><div className="pipeline-footer"><span><span className="health-dot health-healthy" /> Reachable</span><span><span className="health-dot health-unavailable" /> Unavailable</span><span><span className="health-dot health-unknown" /> Unknown / unprobed</span><p>On a cache hit, the worker returns the cached value and skips the chain write.</p></div>{system.error && <p className="inline-error" role="status">{system.error} Service indicators are unknown until the next successful check.</p>}</section>

      <div className="content-grid"><div className="main-column"><section aria-labelledby="prices-heading"><div className="section-header compact"><div><p className="eyebrow">DIRECT FROM THE CONTRACT</p><h2 id="prices-heading">On-chain prices</h2></div><span className="tag">PUBLIC READ</span></div><div className="price-grid"><PriceCard pair="ETH/USD" resource={eth} /><PriceCard pair="BTC/USD" resource={btc} /></div></section>
      <section className="panel jobs-panel" aria-labelledby="jobs-heading"><div className="section-header"><div><p className="eyebrow">THE EVENT STREAM</p><h2 id="jobs-heading">Recent requests <span className="count-badge">{jobs.data?.length ?? "—"}</span></h2></div><span className="muted small">{pending > 0 ? `${pending} in flight` : "Latest 30 jobs"}</span></div>{jobs.error && <p className="inline-error padded" role="status">{jobs.data ? "Last successful history shown. " : ""}{jobs.error}</p>}
        {jobs.loading && !jobs.data ? <div className="empty-state"><span className="loader" /><h3>Reading the job history</h3><p>Waiting for the real backend response.</p></div> : !jobs.data?.length ? <div className="empty-state"><span className="empty-symbol"><Icon name="bolt" size={24} /></span><h3>{jobs.error ? "History unavailable" : "The next request starts here"}</h3><p>{jobs.error ? "Use refresh to try the backend again." : "No retained jobs yet. Unlock the demo and submit a price request to follow its journey."}</p></div> : <div className="jobs-list" role="list" aria-label="Recent oracle requests">{jobs.data.map((job) => <button className={`job-row ${activeJob?.id === job.id ? "selected" : ""}`} key={job.id} onClick={() => setSelectedId(job.id)} aria-pressed={activeJob?.id === job.id}><span className="job-pair"><span>{job.pair}</span><code title={job.id}>{job.id.slice(0, 8)}</code></span><Status value={job.status} /><span className="job-time"><time dateTime={job.createdAt} title={new Date(job.createdAt).toLocaleString()}>{localTime(job.createdAt)}</time><span>View events <Icon name="arrow" size={12} /></span></span></button>)}</div>}
        <div className="panel-footnote">History is retained for a limited time. Only recorded backend events are shown.</div>
      </section></div>

      <aside className="side-column"><section className="panel request-panel" aria-labelledby="request-heading"><div className="request-heading"><span className="square-icon"><Icon name="bolt" size={20} /></span><span className={`tag ${unlocked ? "green" : ""}`}>{unlocked ? "UNLOCKED" : "PASSWORD GATED"}</span></div><p className="eyebrow">TRY THE PIPELINE</p><h2 id="request-heading">Request a price update</h2><p className="muted request-copy">Queue a simulated price for the worker to publish. No wallet or browser signing required.</p>
        <form onSubmit={submit}><label htmlFor="pair">Trading pair</label><div className="select-wrap"><select id="pair" value={pair} onChange={(event) => setPair(event.target.value as Pair)} disabled={submitting}><option value="ETH/USD">ETH / USD — Ethereum</option><option value="BTC/USD">BTC / USD — Bitcoin</option></select></div><button type="submit" className="button primary full" disabled={!unlocked || submitting || authBusy}>{submitting ? <><span className="loader" />Submitting request…</> : unlocked ? <>Submit request <Icon name="arrow" size={17} /></> : <><Icon name="lock" size={16} />Unlock to submit</>}</button></form>
        {submitError && <p className="inline-error" role="alert">{submitError}</p>}
        <div className="access-area">{access.loading && !access.data ? <p className="muted small">Checking browser session…</p> : unlocked ? <div className="session-active"><div><span className="session-label"><span className="health-dot health-healthy" />Demo session active</span><p>Expires at {localTime(access.data!.expiresAt!)}. Public reads stay open.</p></div><button className="button subtle" onClick={lock} disabled={authBusy || submitting}>{authBusy ? "Locking…" : "Lock"}</button></div> : <form onSubmit={unlock}><label htmlFor="demo-password">Shared demo password</label><p className="input-hint">{sessionExpired ? "Your session expired. Enter the password to unlock again." : "Ask the demo operator for access. Viewing is always public."}</p><div className="password-row"><input ref={passwordInput} id="demo-password" name="password" type="password" value={password} onChange={(event) => setPassword(event.target.value)} autoComplete="current-password" placeholder="Enter demo password" required maxLength={1024} disabled={authBusy || cooldownSeconds > 0 || access.data?.available === false} aria-describedby={authError ? "auth-error" : undefined} /><button className="button secondary" disabled={authBusy || !password || cooldownSeconds > 0 || access.data?.available === false}>{authBusy ? "Unlocking…" : "Unlock"}</button></div>{cooldownSeconds > 0 && <p className="input-hint" role="status">Try again in {Math.ceil(cooldownSeconds / 60)} minute{cooldownSeconds > 60 ? "s" : ""}.</p>}{access.data?.available === false && <p className="inline-error">The operator has not configured demo access.</p>}</form>}{authError && <p id="auth-error" className="inline-error" role="alert">{authError}</p>}{access.error && <p className="inline-error" role="status">Session check failed: {access.error}</p>}</div>
        <p className="security-note"><Icon name="lock" size={12} />The signer and API token stay on the server.</p>
      </section>
      <section className="contract-card" aria-labelledby="contract-heading"><div className="contract-heading"><span className="health-dot health-unknown" /><h3 id="contract-heading">OracleConsumer</h3><span className="muted small">{networkName}</span></div><p>{explorer ? "Confirmed updates are independently verifiable on the testnet explorer." : "Local development transactions are not published to a public explorer."}</p>{contractAddress ? explorer ? <a className="contract-link" href={`${explorer}/address/${contractAddress}`} target="_blank" rel="noopener noreferrer"><code>{contractAddress}</code><Icon name="external" size={14} /></a> : <code className="break-code">{contractAddress}</code> : <span className="muted small">{system.loading ? "Loading contract address…" : "Contract address unavailable"}</span>}<span className="contract-chain">CHAIN ID <code>{chainId ?? "—"}</code></span></section></aside></div>

      {notice && <div className="notice" role="status"><span className="health-dot health-healthy" /><p>{notice}</p><button aria-label="Dismiss notification" onClick={() => setNotice(null)}>×</button></div>}
      <section className="panel timeline-panel" aria-labelledby="timeline-heading"><div className="section-header"><div><p className="eyebrow">FOLLOW THE EXECUTION</p><h2 id="timeline-heading">Request timeline</h2></div>{activeJob && <Status value={activeJob.status} />}</div>{selectedId && detail.error && <p className="inline-error padded" role="status">{detail.data ? "Last successful job detail shown. " : ""}{detail.error}</p>}
        {!activeJob ? <div className="timeline-empty"><Icon name="arrow" size={22} /><div><h3>{selectedId && detail.loading ? "Loading this request…" : "Select a request to inspect its journey"}</h3><p>Real status transitions will appear here, with timestamps and transaction proof when available.</p>{selectedId && <code className="break-code">{selectedId}</code>}</div></div> : <div className="timeline-content"><div className="job-summary"><div className="summary-pair">{activeJob.pair}</div><code className="break-code">{activeJob.id}</code><p>{descriptions[activeJob.status]}</p>{activeJob.price !== undefined && <div className="summary-value"><span>Recorded price</span><strong>{formatPrice(activeJob.price)}</strong></div>}{activeJob.error && <p className="inline-error" role="alert">{activeJob.error}</p>}{activeJob.txHash && (explorer ? <a className="transaction-link" href={`${explorer}/tx/${activeJob.txHash}`} target="_blank" rel="noopener noreferrer"><span>View transaction<code>{activeJob.txHash.slice(0, 12)}…{activeJob.txHash.slice(-8)}</code></span><Icon name="external" size={16} /></a> : <code className="break-code">{activeJob.txHash}</code>)}{activeJob.blockNumber !== undefined && <p className="small muted">Included in block {activeJob.blockNumber.toLocaleString()}</p>}<p className="small muted">Created {new Date(activeJob.createdAt).toLocaleString()}</p></div><ol className="timeline">{activeJob.events.length === 0 ? <li className="no-events">No recorded transitions are available yet.</li> : activeJob.events.map((event, index) => <li key={`${event.status}-${event.at}-${index}`} className={index === activeJob.events.length - 1 ? "latest-event" : ""}><span className={`event-node event-${event.status}`} /><div><span className="event-title">{event.status}</span><time dateTime={event.at} title={new Date(event.at).toLocaleString()}>{localTime(event.at)}</time></div><p>{descriptions[event.status as Job["status"]] ?? "Status recorded by the backend."}</p></li>)}</ol></div>}
      </section>
      <footer className="dashboard-footer"><span>ORACLE LAB <span className="footer-divider">/</span> BUILT TO MAKE THE PIPELINE VISIBLE</span><div><span className="poll-note">{refreshing ? "Checking live data…" : lastChecked ? `Last checked ${localTime(lastChecked)} · polls every 5s` : "Connecting…"}</span><button className={`refresh-button ${refreshing ? "refreshing" : ""}`} onClick={refresh} disabled={refreshing} aria-label="Refresh dashboard"><Icon name="refresh" size={15} /></button></div></footer>
    </main>
  </>;
}
