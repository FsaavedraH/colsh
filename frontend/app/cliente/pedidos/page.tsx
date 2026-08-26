"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  nombre_cliente: string;
  total_items: number;
}

const CLIENTE_ID_TEMPORAL = "fbada7ec-ee17-4f0a-80db-b5d469adb4d4";

export default function MisPedidosPage() {
  const [pedidos, setPedidos] = useState<Pedido[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    apiFetch<Pedido[]>(`/api/mis-pedidos?cliente_id=${CLIENTE_ID_TEMPORAL}`, { rol: "Cliente" })
      .then((data) => setPedidos(data || []))
      .catch((err) => setError(err.message))
      .finally(() => setCargando(false));
  }, []);

  function formatearFecha(fecha: string) {
    return new Date(fecha).toLocaleString("es-CO", {
      day: "2-digit",
      month: "2-digit",
      hour: "2-digit",
      minute: "2-digit",
    });
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Mis pedidos</h1>
      <p className="text-gray-500 mb-6">Historial completo de tus pedidos</p>

      {cargando && <p className="text-gray-500">Cargando...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}
      {!cargando && !error && pedidos.length === 0 && (
        <p className="text-gray-500">Todavía no tienes pedidos.</p>
      )}

      <div className="space-y-3">
        {pedidos.map((p) => (
          <Link
            key={p.id_pedido}
            href={`/cliente/pedidos/${p.id_pedido}`}
            className="block bg-white rounded-xl border border-gray-200 p-4 hover:border-blue-500 hover:shadow-sm transition-all"
          >
            <div className="flex justify-between items-start">
              <div>
                <div className="font-mono text-sm text-gray-500">{p.id_pedido.slice(0, 8).toUpperCase()}</div>
                <div className="text-sm text-gray-500 mt-1">{formatearFecha(p.fecha_creacion)} · {p.total_items} ítem{p.total_items !== 1 ? "s" : ""}</div>
              </div>
              <Badge color="yellow">{p.estado}</Badge>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}