"use client";

import { useState } from "react";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";

export default function LoginPage() {
  const { login } = useAuth();
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  async function manejarSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setCargando(true);
    try {
      await login(email, password);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <form onSubmit={manejarSubmit} className="bg-white rounded-xl border border-gray-200 p-8 w-full max-w-sm">
        <h1 className="text-xl font-bold mb-1">ColSh</h1>
        <p className="text-gray-500 text-sm mb-6">Sistema de Trazabilidad Logística</p>

        <label className="text-sm text-gray-500 block mb-1">Correo</label>
        <input
          type="email"
          value={email}
          onChange={(e) => setEmail(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
          required
        />

        <label className="text-sm text-gray-500 block mb-1">Contraseña</label>
        <input
          type="password"
          value={password}
          onChange={(e) => setPassword(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
          required
        />

        {error && (
          <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">{error}</div>
        )}

        <Button type="submit" disabled={cargando}>
          {cargando ? "Ingresando..." : "Iniciar sesión"}
        </Button>
      </form>
    </div>
  );
}