import Link from "next/link";

export default function NotFound() {
  return <main className="fallback"><p className="eyebrow">ORACLE LAB / 404</p><h1>This page does not exist.</h1><Link className="button primary" href="/">Return to the dashboard</Link></main>;
}
