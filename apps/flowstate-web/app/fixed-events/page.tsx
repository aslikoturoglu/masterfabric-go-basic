"use client";

import { useEffect, useState } from "react";
import { Calendar, Plus, Trash2, ArrowLeft } from "lucide-react";
import Link from "next/link";
import { gql } from "@/lib/graphql";

const DAY_NAMES: Record<number, string> = { 1: "Pzt", 2: "Sal", 3: "Çar", 4: "Per", 5: "Cum", 6: "Cmt", 7: "Paz" };

type FixedEvent = { id: string; title: string; startTime: string; endTime: string; daysOfWeek: number[]; category: string };

export default function FixedEventsPage() {
  const [events, setEvents] = useState<FixedEvent[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ title: "", startTime: "09:00", endTime: "18:00", daysOfWeek: [1, 2, 3, 4, 5], category: "Work" });
  const token = typeof window !== "undefined" ? localStorage.getItem("accessToken") : null;

  const fetchEvents = async () => {
    if (!token) return;
    try {
      const data = await gql<{ flowstateFixedEvents: FixedEvent[] }>(
        `query { flowstateFixedEvents { id title startTime endTime daysOfWeek category } }`,
        undefined,
        token
      );
      setEvents(data.flowstateFixedEvents || []);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchEvents();
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    await gql(
      `mutation Create($input: FlowstateFixedEventInput!) { flowstateCreateFixedEvent(input: $input) { id } }`,
      { input: { ...form, category: form.category as "Work" | "School" | "Meeting" } },
      token
    );
    setForm({ title: "", startTime: "09:00", endTime: "18:00", daysOfWeek: [1, 2, 3, 4, 5], category: "Work" });
    setShowForm(false);
    fetchEvents();
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    await gql(`mutation Delete($id: UUID!) { flowstateDeleteFixedEvent(id: $id) }`, { id }, token);
    fetchEvents();
  };

  const toggleDay = (d: number) => {
    setForm((f) => ({
      ...f,
      daysOfWeek: f.daysOfWeek.includes(d) ? f.daysOfWeek.filter((x) => x !== d) : [...f.daysOfWeek, d].sort(),
    }));
  };

  if (loading) return <div className="p-8">Yükleniyor...</div>;
  if (!token) return <div className="p-8">Giriş yapmanız gerekiyor. <Link href="/" className="text-flowstate-600">Ana sayfa</Link></div>;

  return (
    <div className="min-h-screen bg-slate-50">
      <nav className="border-b bg-white px-6 py-4">
        <Link href="/" className="flex items-center gap-2 text-slate-600 hover:text-slate-900"><ArrowLeft className="h-5 w-5" /> Geri</Link>
      </nav>
      <main className="mx-auto max-w-4xl px-6 py-8">
        <h1 className="mb-6 flex items-center gap-2 text-2xl font-bold"><Calendar className="h-7 w-7 text-flowstate-600" /> Sabit Program</h1>
        {showForm ? (
          <form onSubmit={handleSubmit} className="mb-8 rounded-xl border bg-white p-6 shadow-sm">
            <h2 className="mb-4 font-semibold">Yeni Etkinlik</h2>
            <div className="space-y-4">
              <input type="text" placeholder="Başlık" value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} className="w-full rounded-lg border px-4 py-2" required />
              <div className="flex gap-4">
                <input type="time" value={form.startTime} onChange={(e) => setForm({ ...form, startTime: e.target.value })} className="rounded-lg border px-4 py-2" />
                <input type="time" value={form.endTime} onChange={(e) => setForm({ ...form, endTime: e.target.value })} className="rounded-lg border px-4 py-2" />
              </div>
              <div>
                <p className="mb-2 text-sm text-slate-600">Günler</p>
                <div className="flex gap-2">
                  {[1, 2, 3, 4, 5, 6, 7].map((d) => (
                    <button key={d} type="button" onClick={() => toggleDay(d)} className={`rounded-lg px-3 py-1 text-sm ${form.daysOfWeek.includes(d) ? "bg-flowstate-600 text-white" : "bg-slate-200 text-slate-600"}`}>
                      {DAY_NAMES[d]}
                    </button>
                  ))}
                </div>
              </div>
              <select value={form.category} onChange={(e) => setForm({ ...form, category: e.target.value })} className="rounded-lg border px-4 py-2">
                <option value="Work">İş</option><option value="School">Okul</option><option value="Meeting">Toplantı</option>
              </select>
            </div>
            <div className="mt-4 flex gap-2">
              <button type="submit" className="rounded-lg bg-flowstate-600 px-4 py-2 text-white">Kaydet</button>
              <button type="button" onClick={() => setShowForm(false)} className="rounded-lg bg-slate-200 px-4 py-2">İptal</button>
            </div>
          </form>
        ) : (
          <button onClick={() => setShowForm(true)} className="mb-8 flex items-center gap-2 rounded-lg bg-flowstate-600 px-4 py-2 text-white"><Plus className="h-4 w-4" /> Yeni Etkinlik</button>
        )}
        <div className="space-y-3">
          {events.map((e) => (
            <div key={e.id} className="flex items-center justify-between rounded-xl border bg-white p-4 shadow-sm">
              <div>
                <p className="font-medium">{e.title}</p>
                <p className="text-sm text-slate-600">{e.startTime} - {e.endTime} • {e.daysOfWeek.map((d) => DAY_NAMES[d]).join(", ")} • {e.category}</p>
              </div>
              <button onClick={() => handleDelete(e.id)} className="rounded p-2 text-red-600 hover:bg-red-50"><Trash2 className="h-4 w-4" /></button>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
