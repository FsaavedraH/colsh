"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";
import Button from "@/components/ui/Button";

interface Usuario {
  id_usuario: string;
  nombre: string;
  email: string;
  rol: string;
}

const colorPorRol: Record<string, "green" | "red" | "yellow" | "gray"> = {
  Cliente: "gray",
  Picking: "yellow",
  Empaque: "green",
  Transportista: "yellow",
  Administrador: "red",
};

const rolesCreables = ["Picking", "Empaque", "Transportista", "Administrador"];

export default function UsuariosPage() {
  const [usuarios, setUsuarios] = useState<Usuario[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  const [mostrarFormulario, setMostrarFormulario] = useState(false);
  const [nombre, setNombre] = useState("");
  const [email, setEmail] = useState("");
  const [password, setPassword] = useState("");
  const [rol, setRol] = useState(rolesCreables[0]);
  const [errorFormulario, setErrorFormulario] = useState("");
  const [creando, setCreando] = useState(false);

  useEffect(() => {
    cargarUsuarios();
  }, []);

  function cargarUsuarios() {
    setCargando(true);
    apiFetch<Usuario[]>("/api/usuarios", { rol: "Administrador" })
      .then((data) => setUsuarios(data || []))
      .catch((err) => setError(err.message))
      .finally(() => setCargando(false));
  }

  function limpiarFormulario() {
    setNombre("");
    setEmail("");
    setPassword("");
    setRol(rolesCreables[0]);
    setErrorFormulario("");
  }

  async function manejarCrearUsuario(e: React.FormEvent) {
    e.preventDefault();
    setErrorFormulario("");
    setCreando(true);
    try {
      await apiFetch("/api/usuarios", {
        method: "POST",
        rol: "Administrador",
        body: JSON.stringify({ nombre, email, password, rol }),
      });
      limpiarFormulario();
      setMostrarFormulario(false);
      cargarUsuarios();
    } catch (err: any) {
      setErrorFormulario(err.message);
    } finally {
      setCreando(false);
    }
  }

  return (
    <div>
      <div className="flex items-center justify-between mb-1">
        <h1 className="text-2xl font-bold">Usuarios</h1>
        <Button onClick={() => setMostrarFormulario((v) => !v)}>
          {mostrarFormulario ? "Cancelar" : "+ Nuevo usuario"}
        </Button>
      </div>
      <p className="text-gray-500 mb-6">
        Mostrando {usuarios.length} usuario{usuarios.length !== 1 ? "s" : ""}
      </p>

      {mostrarFormulario && (
        <form
          onSubmit={manejarCrearUsuario}
          className="bg-white rounded-xl border border-gray-200 p-5 mb-6 max-w-lg"
        >
          <h2 className="font-semibold mb-4">Crear usuario operativo</h2>

          <label className="text-sm text-gray-500 block mb-1">Nombre</label>
          <input
            type="text"
            value={nombre}
            onChange={(e) => setNombre(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-3"
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Correo</label>
          <input
            type="email"
            value={email}
            onChange={(e) => setEmail(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-3"
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Contraseña</label>
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-3"
            minLength={6}
            required
          />

          <label className="text-sm text-gray-500 block mb-1">Rol</label>
          <select
            value={rol}
            onChange={(e) => setRol(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
          >
            {rolesCreables.map((r) => (
              <option key={r} value={r}>
                {r}
              </option>
            ))}
          </select>

          {errorFormulario && (
            <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">
              {errorFormulario}
            </div>
          )}

          <Button type="submit" disabled={creando}>
            {creando ? "Creando..." : "Crear usuario"}
          </Button>
        </form>
      )}

      {cargando && <p className="text-gray-500">Cargando usuarios...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-500 text-left">
            <tr>
              <th className="px-4 py-3">Nombre</th>
              <th className="px-4 py-3">Correo</th>
              <th className="px-4 py-3">Rol</th>
            </tr>
          </thead>
          <tbody>
            {usuarios.map((u) => (
              <tr key={u.id_usuario} className="border-t border-gray-100">
                <td className="px-4 py-3 font-medium">{u.nombre}</td>
                <td className="px-4 py-3 text-gray-500">{u.email}</td>
                <td className="px-4 py-3">
                  <Badge color={colorPorRol[u.rol] || "gray"}>{u.rol}</Badge>
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}