"use client";

import { useEffect, useState } from "react";
import { ListTodo, Plus, Trash2, ArrowLeft } from "lucide-react";
import Link from "next/link";
import { gql } from "@/lib/graphql";

type FlexibleTask = { id: string; title: string; durationMinutes: number; frequencyPerWeek: number; priority: string; preferredContext: string };

export default function FlexibleTasksPage() {
  const [tasks, setTasks] = useState<FlexibleTask[]>([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState({ title: "", durationMinutes: 60, frequencyPerWeek: 3, priority: "Medium", preferredContext: "Morning" });
  const token = typeof window !== "undefined" ? localStorage.getItem("accessToken") : null;

  const fetchTasks = async () => {
    if (!token) return;
    try {
      const data = await gql<{ flowstateFlexibleTasks: FlexibleTask[] }>(
        `query { flowstateFlexibleTasks { id title durationMinutes frequencyPerWeek priority preferredContext } }`,
        undefined,
        token
      );
      setTasks(data.flowstateFlexibleTasks || []);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchTasks();
  }, [token]);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    await gql(
      `mutation Create($input: FlowstateFlexibleTaskInput!) { flowstateCreateFlexibleTask(input: $input) { id } }`,
      { input: form },
      token
    );
    setForm({ title: "", durationMinutes: 60, frequencyPerWeek: 3, priority: "Medium", preferredContext: "Morning" });
    setShowForm(false);
    fetchTasks();
  };

  const handleDelete = async (id: string) => {
    if (!token) return;
    await gql(`mutation Delete($id: UUID!) { flowstateDeleteFlexibleTask(id: $id) }`, { id }, token);
    fetchTasks();
  };

  if (loading) return <div className="p-8">Yükleniyor...</div>;
  if (!token) return <div className="p-8">Giriş yapmanız gerekiyor. <Link href="/" className="text-flowstate-600">Ana sayfa</Link></div>;

  return (
    <div className="min-h-screen bg-slate-50">
      <nav className="border-b bg-white px-6 py-4">
        <Link href="/" className="flex items-center gap-2 text-slate-600 hover:text-slate-900"><ArrowLeft className="h-5 w-5" /> Geri</Link>
      </nav>
      <main className="mx-auto max-w-4xl px-6 py-8">
        <h1 className="mb-6 flex items-center gap-2 text-2xl font-bold"><ListTodo className="h-7 w-7 text-amber-600" /> Esnek Görevler</h1>
        {showForm ? (
          <form onSubmit={handleSubmit} className="mb-8 rounded-xl border bg-white p-6 shadow-sm">
            <div className="space-y-4">
              <input type="text" placeholder="Başlık" value={form.title} onChange={(e) => setForm({ ...form, title: e.target.value })} className="w-full rounded-lg border px-4 py-2" required />
              <input type="number" min={1} value={form.durationMinutes} onChange={(e) => setForm({ ...form, durationMinutes: parseInt(e.target.value) || 60 })} className="w-full rounded-lg border px-4 py-2" />
              <input type="number" min={1} max={7} value={form.frequencyPerWeek} onChange={(e) => setForm({ ...form, frequencyPerWeek: parseInt(e.target.value) || 1 })} className="w-full rounded-lg border px-4 py-2" />
              <select value={form.priority} onChange={(e) => setForm({ ...form, priority: e.target.value })} className="w-full rounded-lg border px-4 py-2">
                <option value="High">Yüksek</option><option value="Medium">Orta</option><option value="Low">Düşük</option>
              </select>
              <select value={form.preferredContext} onChange={(e) => setForm({ ...form, preferredContext: e.target.value })} className="w-full rounded-lg border px-4 py-2">
                <option value="Morning">Sabah</option><option value="Evening">Akşam</option><option value="WorkBreak">İş molası</option><option value="Weekend">Hafta sonu</option>
              </select>
            </div>
            <div className="mt-4 flex gap-2">
              <button type="submit" className="rounded-lg bg-flowstate-600 px-4 py-2 text-white">Kaydet</button>
              <button type="button" onClick={() => setShowForm(false)} className="rounded-lg bg-slate-200 px-4 py-2">İptal</button>
            </div>
          </form>
        ) : (
          <button onClick={() => setShowForm(true)} className="mb-8 flex items-center gap-2 rounded-lg bg-flowstate-600 px-4 py-2 text-white"><Plus className="h-4 w-4" /> Yeni Görev</button>
        )}
        <div className="space-y-3">
          {tasks.map((t) => (
            <div key={t.id} className="flex items-center justify-between rounded-xl border bg-white p-4 shadow-sm">
              <div>
                <p className="font-medium">{t.title}</p>
                <p className="text-sm text-slate-600">{t.durationMinutes} dk • Haftada {t.frequencyPerWeek} kez • {t.priority}</p>
              </div>
              <button onClick={() => handleDelete(t.id)} className="rounded p-2 text-red-600 hover:bg-red-50"><Trash2 className="h-4 w-4" /></button>
            </div>
          ))}
        </div>
      </main>
    </div>
  );
}
