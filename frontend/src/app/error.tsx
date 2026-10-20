"use client";

export default function ErrorPage({ reset }: { reset: () => void }) {
  return <main className="fallback"><p className="eyebrow">ORACLE LAB / ERROR</p><h1>The dashboard could not render.</h1><p>No request is submitted by reloading this view.</p><button className="button primary" onClick={reset}>Try again</button></main>;
}
