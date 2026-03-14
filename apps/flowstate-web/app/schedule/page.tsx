"use client";

import { useState } from "react";
import { Sparkles, ArrowLeft, Loader2 } from "lucide-react";
import Link from "next/link";
import { gql } from "@/lib/graphql";

const DAY_NAMES: Record<string, string> = { monday: "Pzt", tuesday: "Sal", wednesday: "Çar", thursday: "Per", friday: "Cum", saturday: "Cmt", sunday: "Paz" };

export default function SchedulePage() {
  const [schedule, setSchedule] = useState<Record<string, unknown> | null>(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState("");
  const token = typeof window !== "undefined" ? localStorage.getItem("accessToken") : null;

  const handleGenerate = async () => {
    setLoading(true);
    setError("");
    if (!token) {
      setError("Giriş yapmanız gerekiyor");
      setLoading(false);
      return;
    }
    try {
      const data = await gql<{ flowstateGenerateSchedule: { weekIdentifier: string; scheduleData: string } }>(
        `mutation { flowstateGenerateSchedule { weekIdentifier scheduleData } }`,
        undefined,
        token
      );
      const parsed = JSON.parse(data.flowstateGenerateSchedule.scheduleData || "{}");
      setSchedule(parsed);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Oluşturulamadı");
    } finally {
      setLoading(false);
    }
  };

  const days = schedule?.days as Record<string, unknown[]> | undefined;

  if (!token) {
    return (
      <div className="p-8">
        Giriş yapmanız gerekiyor. <Link href="/" className="text-flowstate-600">Ana sayfa</Link>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-50">
      <nav className="border-b bg-white px-6 py-4">
        <Link href="/" className="flex items-center gap-2 text-slate-600 hover:text-slate-900"><ArrowLeft className="h-5 w-5" /> Geri</Link>
      </nav>
      <main className="mx-auto max-w-6xl px-6 py-8">
        <h1 className="mb-6 flex items-center gap-2 text-2xl font-bold"><Sparkles className="h-7 w-7 text-flowstate-600" /> Haftalık Takvim</h1>
        <button onClick={handleGenerate} disabled={loading} className="mb-8 flex items-center gap-2 rounded-lg bg-flowstate-600 px-6 py-3 font-medium text-white hover:bg-flowstate-700 disabled:opacity-50">
          {loading ? <><Loader2 className="h-5 w-5 animate-spin" /> Oluşturuluyor...</> : <><Sparkles className="h-5 w-5" /> AI ile Takvim Oluştur</>}
        </button>
        {error && <p className="mb-4 text-red-600">{error}</p>}
        {schedule && (
          <div className="space-y-6">
            <p className="text-slate-600">Hafta: {String(schedule.week_identifier || "")}</p>
            <div className="grid gap-4 md:grid-cols-2 lg:grid-cols-3">
              {days && Object.entries(days).map(([dayKey, slots]) => (
                <div key={dayKey} className="rounded-xl border bg-white p-4 shadow-sm">
                  <h3 className="mb-3 font-semibold">{DAY_NAMES[dayKey] || dayKey}</h3>
                  <div className="space-y-2">
                    {Array.isArray(slots) && slots.map((slot: Record<string, unknown>, i: number) => (
                      <div key={i} className={`rounded-lg px-3 py-2 text-sm ${slot.type === "fixed" ? "bg-slate-100" : "bg-flowstate-50"}`}>
                        <span className="font-medium">{String(slot.task)}</span>
                        {slot.time && <span className="ml-2 text-slate-600">{String(slot.time)}{slot.duration ? ` (${slot.duration} dk)` : ""}</span>}
                      </div>
                    ))}
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
