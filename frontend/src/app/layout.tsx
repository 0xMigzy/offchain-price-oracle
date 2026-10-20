import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Oracle Lab · Sepolia Price Pipeline",
  description: "Follow a real asynchronous oracle request from an API through a Redis Streams queue, a Go worker, and a Redis cache to an Ethereum Sepolia contract. Educational simulated prices, real testnet transactions.",
  robots: { index: false, follow: false },
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return <html lang="en"><body>{children}</body></html>;
}
