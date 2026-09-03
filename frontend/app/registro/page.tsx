"use client";

import { useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Button from "@/components/ui/Button";

export default function RegistroPage() {
  const router = useRouter();
  const [nombre, setNombre] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  async function manejarSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");
    setCargando(true);
    try {
      await apiFetch("/api/auth/registro", {
        method: "POST",
        body: JSON.stringify({ nombre, email, password }),
      });
      router.push("/login");
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  return (
    <div className="min-h-screen flex items-center justify-center bg-gray-50">
      <form onSubmit={manejarSubmit} className="bg-white rounded-xl border border-gray-200 p-8 w-full max-w-sm">
        <h1 className="text-xl font-bold mb-1">Crear cuenta</h1>
        <p className="text-gray-500 text-sm mb-6">Registro de cliente en ColSh</p>

        <label className="text-sm text-gray-500 block mb-1">Nombre</label>
        <input
          type="text"
          value={nombre}
          onChange={(e) => setNombre(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
          required
        />

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
          minLength={6}
          required
        />

        {error && (
          <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">{error}</div>
        )}

        <Button type="submit" disabled={cargando}>
          {cargando ? "Creando cuenta..." : "Registrarme"}
        </Button>

        <p className="text-sm text-gray-500 mt-4 text-center">
          ¿Ya tienes cuenta? <a href="/login" className="text-blue-600 underline">Inicia sesión</a>
        </p>
      </form>
    </div>
  );
}