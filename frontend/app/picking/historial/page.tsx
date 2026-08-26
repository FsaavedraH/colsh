"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface Orden {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  nombre_cliente: string;
  total_items: number;
}

export default function HistorialPickingPage() {
  const [ordenes, setOrdenes] = useState<Orden[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    apiFetch<Orden[]>("/api/picking/historial", { rol: "Picking" })
      .then((data) => setOrdenes(data || []))
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
      <h1 className="text-2xl font-bold mb-1">Historial de Picking</h1>
      <p className="text-gray-500 mb-6">Pedidos ya recolectados</p>

      {cargando && <p className="text-gray-500">Cargando...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}
      {!cargando && !error && ordenes.length === 0 && (
        <p className="text-gray-500">Sin historial todavía.</p>
      )}

      <div className="space-y-3">
        {ordenes.map((o) => (
          <div key={o.id_pedido} className="bg-white rounded-xl border border-gray-200 p-4 flex justify-between items-start">
            <div>
              <div className="font-mono text-sm text-gray-500">{o.id_pedido.slice(0, 8).toUpperCase()}</div>
              <div className="font-semibold text-gray-800">{o.nombre_cliente}</div>
              <div className="text-sm text-gray-500 mt-1">{formatearFecha(o.fecha_creacion)} · {o.total_items} ítem{o.total_items !== 1 ? "s" : ""}</div>
            </div>
            <Badge color="green">{o.estado}</Badge>
          </div>
        ))}
      </div>
    </div>
  );
}