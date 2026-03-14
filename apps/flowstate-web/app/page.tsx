"use client";

import { useState } from "react";
import { Calendar, ListTodo, Sparkles, LogIn, UserPlus } from "lucide-react";
import Link from "next/link";
import { gql } from "@/lib/graphql";

export default function Home() {
  const [token, setToken] = useState<string | null>(
    typeof window !== "undefined" ? localStorage.getItem("accessToken") : null
  );
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [displayName, setDisplayName] = useState("");
  const [mode, setMode] = useState<"login" | "register">("login");
  const [error, setError] = useState("");
  const [loading, setLoading] = useState(false);

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await gql<{ login: { accessToken: string } }>(
        `mutation Login($input: LoginInput!) { login(input: $input) { accessToken } }`,
        { input: { email, password } }
      );
      setToken(data.login.accessToken);
      localStorage.setItem("accessToken", data.login.accessToken);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Giriş başarısız");
    } finally {
      setLoading(false);
    }
  };

  const handleRegister = async (e: React.FormEvent) => {
    e.preventDefault();
    setError("");
    setLoading(true);
    try {
      const data = await gql<{ register: { accessToken: string } }>(
        `mutation Register($input: RegisterInput!) { register(input: $input) { accessToken } }`,
        { input: { email, password, displayName } }
      );
      setToken(data.register.accessToken);
      localStorage.setItem("accessToken", data.register.accessToken);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Kayıt başarısız");
    } finally {
      setLoading(false);
    }
  };

  const handleLogout = () => {
    setToken(null);
    localStorage.removeItem("accessToken");
  };

  if (token) {
    return (
      <div className="min-h-screen bg-slate-50">
        <nav className="border-b bg-white px-6 py-4">
          <div className="mx-auto flex max-w-6xl items-center justify-between">
            <div className="flex items-center gap-2">
              <Sparkles className="h-6 w-6 text-flowstate-600" />
              <span className="text-xl font-semibold">FlowState AI</span>
            </div>
            <button onClick={handleLogout} className="rounded-lg bg-slate-200 px-4 py-2 text-sm font-medium hover:bg-slate-300">
              Çıkış
            </button>
          </div>
        </nav>
        <main className="mx-auto max-w-6xl px-6 py-12">
          <p className="mb-8 text-slate-600">MasterFabric GraphQL API üzerinde çalışıyor.</p>
          <div className="grid gap-8 md:grid-cols-2">
            <Link
              href="/fixed-events"
              className="flex flex-col items-center gap-4 rounded-xl border bg-white p-8 shadow-sm transition hover:shadow-md"
            >
              <div className="rounded-full bg-flowstate-100 p-4">
                <Calendar className="h-8 w-8 text-flowstate-600" />
              </div>
              <h2 className="text-lg font-semibold">Sabit Program</h2>
              <p className="text-center text-sm text-slate-600">İş, okul, toplantı gibi değişmeyen etkinlikler</p>
            </Link>
            <Link
              href="/flexible-tasks"
              className="flex flex-col items-center gap-4 rounded-xl border bg-white p-8 shadow-sm transition hover:shadow-md"
            >
              <div className="rounded-full bg-amber-100 p-4">
                <ListTodo className="h-8 w-8 text-amber-600" />
              </div>
              <h2 className="text-lg font-semibold">Esnek Görevler</h2>
              <p className="text-center text-sm text-slate-600">Spor, ev işi gibi haftalık görevler</p>
            </Link>
          </div>
          <div className="mt-12 text-center">
            <Link
              href="/schedule"
              className="inline-flex items-center gap-2 rounded-lg bg-flowstate-600 px-6 py-3 font-medium text-white hover:bg-flowstate-700"
            >
              <Sparkles className="h-5 w-5" />
              Haftalık Takvim Oluştur
            </Link>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="flex min-h-screen items-center justify-center bg-gradient-to-br from-slate-50 to-flowstate-50">
      <div className="w-full max-w-md rounded-2xl border bg-white p-8 shadow-lg">
        <div className="mb-8 flex items-center justify-center gap-2">
          <Sparkles className="h-8 w-8 text-flowstate-600" />
          <h1 className="text-2xl font-bold">FlowState AI</h1>
        </div>
        <p className="mb-6 text-center text-slate-600">MasterFabric ile giriş yapın</p>
        <div className="mb-4 flex gap-2">
          <button
            onClick={() => setMode("login")}
            className={`flex-1 rounded-lg py-2 text-sm font-medium ${mode === "login" ? "bg-flowstate-600 text-white" : "bg-slate-100 text-slate-600"}`}
          >
            <LogIn className="mr-2 inline h-4 w-4" /> Giriş
          </button>
          <button
            onClick={() => setMode("register")}
            className={`flex-1 rounded-lg py-2 text-sm font-medium ${mode === "register" ? "bg-flowstate-600 text-white" : "bg-slate-100 text-slate-600"}`}
          >
            <UserPlus className="mr-2 inline h-4 w-4" /> Kayıt
          </button>
        </div>
        {mode === "login" ? (
          <form onSubmit={handleLogin} className="space-y-4">
            <input type="email" placeholder="E-posta" value={email} onChange={(e) => setEmail(e.target.value)} className="w-full rounded-lg border px-4 py-2" required />
            <input type="password" placeholder="Şifre" value={password} onChange={(e) => setPassword(e.target.value)} className="w-full rounded-lg border px-4 py-2" required />
            <button type="submit" disabled={loading} className="w-full rounded-lg bg-flowstate-600 py-2 font-medium text-white hover:bg-flowstate-700 disabled:opacity-50">
              {loading ? "Giriş yapılıyor..." : "Giriş Yap"}
            </button>
          </form>
        ) : (
          <form onSubmit={handleRegister} className="space-y-4">
            <input type="text" placeholder="Görünen ad" value={displayName} onChange={(e) => setDisplayName(e.target.value)} className="w-full rounded-lg border px-4 py-2" required />
            <input type="email" placeholder="E-posta" value={email} onChange={(e) => setEmail(e.target.value)} className="w-full rounded-lg border px-4 py-2" required />
            <input type="password" placeholder="Şifre (min 8 karakter)" value={password} onChange={(e) => setPassword(e.target.value)} className="w-full rounded-lg border px-4 py-2" required minLength={8} />
            <button type="submit" disabled={loading} className="w-full rounded-lg bg-flowstate-600 py-2 font-medium text-white hover:bg-flowstate-700 disabled:opacity-50">
              {loading ? "Kayıt yapılıyor..." : "Kayıt Ol"}
            </button>
          </form>
        )}
        {error && <p className="mt-4 text-center text-sm text-red-600">{error}</p>}
      </div>
    </div>
  );
}
