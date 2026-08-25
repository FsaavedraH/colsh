"use client";

import { useEffect, useState } from "react";
import Link from "next/link";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface OrdenPicking {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  nombre_cliente: string;
  total_items: number;
}

export default function ListaPickingPage() {
  const [ordenes, setOrdenes] = useState<OrdenPicking[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    cargarOrdenes();
  }, []);

  async function cargarOrdenes() {
    setCargando(true);
    setError("");
    try {
      const data = await apiFetch<OrdenPicking[]>("/api/picking", { rol: "Picking" });
      setOrdenes(data || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

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
      <h1 className="text-2xl font-bold mb-1">Órdenes de Picking (FIFO)</h1>
      <p className="text-gray-500 mb-6">
        Mostrando {ordenes.length} orden{ordenes.length !== 1 ? "es" : ""} pendiente{ordenes.length !== 1 ? "s" : ""}
      </p>

      {cargando && <p className="text-gray-500">Cargando órdenes...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}

      {!cargando && !error && ordenes.length === 0 && (
        <p className="text-gray-500">No hay órdenes pendientes de recolección en este momento.</p>
      )}

      <div className="space-y-3">
        {ordenes.map((orden) => (
          <Link
            key={orden.id_pedido}
            href={`/picking/${orden.id_pedido}`}
            className="block bg-white rounded-xl border border-gray-200 p-4 hover:border-orange-400 hover:shadow-sm transition-all"
          >
            <div className="flex justify-between items-start">
              <div>
                <div className="font-mono text-sm text-gray-500">
                  {orden.id_pedido.slice(0, 8).toUpperCase()}
                </div>
                <div className="font-semibold text-gray-800">{orden.nombre_cliente}</div>
                <div className="text-sm text-gray-500 mt-1">
                  {formatearFecha(orden.fecha_creacion)} · {orden.total_items} ítem{orden.total_items !== 1 ? "s" : ""}
                </div>
              </div>
              <Badge color="yellow">{orden.estado}</Badge>
            </div>
          </Link>
        ))}
      </div>
    </div>
  );
}