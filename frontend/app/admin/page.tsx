"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

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

export default function UsuariosPage() {
  const [usuarios, setUsuarios] = useState<Usuario[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    apiFetch<Usuario[]>("/api/usuarios", { rol: "Administrador" })
      .then((data) => setUsuarios(data || []))
      .catch((err) => setError(err.message))
      .finally(() => setCargando(false));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Usuarios</h1>
      <p className="text-gray-500 mb-6">
        Mostrando {usuarios.length} usuario{usuarios.length !== 1 ? "s" : ""}
      </p>

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