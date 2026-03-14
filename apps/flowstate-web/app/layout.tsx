import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "FlowState AI — Haftalık Optimizasyon",
  description: "MasterFabric üzerine inşa edilmiş AI destekli haftalık planlama",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="tr">
      <body className="min-h-screen antialiased">{children}</body>
    </html>
  );
}
