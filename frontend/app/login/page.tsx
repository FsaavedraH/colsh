"use client";

import { Suspense, useState } from "react";
import Link from "next/link";
import { useSearchParams } from "next/navigation";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";
import IlustracionMoto from "@/components/ui/IlustracionMoto";

function FormularioLogin() {
  const { login } = useAuth();
  const searchParams = useSearchParams();
  const expiroPorInactividad = searchParams.get("motivo") === "inactividad";

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
    <form onSubmit={manejarSubmit} className="bg-white rounded-xl border border-gray-200 p-8 w-full max-w-sm">
      <h1 className="text-xl font-bold mb-1">Iniciar sesión</h1>
      <p className="text-gray-500 text-sm mb-6">Ingresa tus credenciales para continuar</p>

      {expiroPorInactividad && (
        <div className="bg-amber-100 text-amber-800 text-sm rounded-lg p-3 mb-4">
          Tu sesión se cerró por inactividad (15 minutos sin uso). Inicia sesión de nuevo.
        </div>
      )}

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
        placeholder="Ingresa tu contraseña"
        required
      />

      {error && (
        <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">{error}</div>
      )}

      <Button type="submit" disabled={cargando}>
        {cargando ? "Ingresando..." : "Iniciar sesión"}
      </Button>

      <p className="text-sm text-gray-500 text-center mt-5">
        ¿No tienes cuenta?{" "}
        <Link href="/registro" className="text-blue-600 font-medium hover:underline">
          Regístrate
        </Link>
      </p>
    </form>
  );
}

export default function LoginPage() {
  return (
    <div className="min-h-screen flex">
      <div className="hidden lg:flex flex-col justify-center items-center w-1/2 bg-[#1e3a5f] text-white p-12">
        <div className="w-64 aspect-[3/2]">
          <IlustracionMoto />
        </div>
        <h2 className="text-2xl font-bold mt-6 mb-2 text-center">Trazabilidad en cada kilómetro</h2>
        <p className="text-white/70 text-center max-w-sm">
          Controla tus pedidos de repuestos de motocicletas, desde la compra hasta la entrega.
        </p>
      </div>

      <div className="flex-1 flex items-center justify-center bg-gray-50 p-4">
        <Suspense fallback={<div className="text-gray-400 text-sm">Cargando...</div>}>
          <FormularioLogin />
        </Suspense>
      </div>
    </div>
  );
}