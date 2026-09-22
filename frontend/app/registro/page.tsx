"use client";

import { useState } from "react";
import Link from "next/link";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Button from "@/components/ui/Button";
import IlustracionMoto from "@/components/ui/IlustracionMoto";

export default function RegistroPage() {
  const router = useRouter();
  const [nombre, setNombre] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [confirmarPassword, setConfirmarPassword] = useState("");
  const [error, setError] = useState("");
  const [cargando, setCargando] = useState(false);

  async function manejarSubmit(e: React.FormEvent) {
    e.preventDefault();
    setError("");

    if (password !== confirmarPassword) {
      setError("Las contraseñas no coinciden.");
      return;
    }

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
    <div className="min-h-screen flex">
      <div className="hidden lg:flex flex-col justify-center items-center w-1/2 bg-[#1e3a5f] text-white p-12">
        <div className="w-64 aspect-[3/2]">
          <IlustracionMoto />
        </div>
        <h2 className="text-2xl font-bold mt-6 mb-2 text-center">Únete a ColSh</h2>
        <p className="text-white/70 text-center max-w-sm">
          Crea tu cuenta para hacer pedidos y hacer seguimiento de tus repuestos.
        </p>
      </div>

      <div className="flex-1 flex items-center justify-center bg-gray-50 p-4">
        <form onSubmit={manejarSubmit} className="bg-white rounded-xl border border-gray-200 p-8 w-full max-w-sm">
          <h1 className="text-xl font-bold mb-1">Crear cuenta de cliente</h1>
          <p className="text-gray-500 text-sm mb-6">
            Regístrate para hacer pedidos y hacer seguimiento
          </p>

          <label className="text-sm text-gray-500 block mb-1">Nombre completo</label>
          <input
            type="text"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
            placeholder="Juan Pérez"
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Correo electrónico</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
            placeholder="tu@correo.com"
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Contraseña</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
            placeholder="Mínimo 6 caracteres"
            minLength={6}
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Confirmar contraseña</label>
          <input
            type="password"
            value={confirmarPassword}
            onChange={(e) => setConfirmarPassword(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
            placeholder="Repite tu contraseña"
            minLength={6}
            required
          />

          {error && (
            <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">{error}</div>
          )}

          <Button type="submit" disabled={cargando}>
            {cargando ? "Creando cuenta..." : "Crear cuenta"}
          </Button>

          <p className="text-sm text-gray-500 text-center mt-5">
            ¿Ya tienes cuenta?{" "}
            <Link href="/login" className="text-blue-600 font-medium hover:underline">
              Inicia sesión
            </Link>
          </p>
        </form>
      </div>
    </div>
  );
}