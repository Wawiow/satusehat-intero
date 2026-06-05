import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "SatuSehat Intero | Dashboard Rumah Sakit",
  description: "Dashboard integrasi rumah sakit untuk operasional SatuSehat Intero.",
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="id">
      <body>{children}</body>
    </html>
  );
}
