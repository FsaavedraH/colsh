"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface EventoPendiente {
  id: string;
  id_evento_original: string;
  id_pedido: string;
  estado: string;
  fecha_evento_real: string;
  responsable: string;
  intentos: number;
  ultimo_intento: string | null;
  creado_en: string;
}

export default function LedgerPendientesPage() {
  const [pendientes, setPendientes] = useState<EventoPendiente[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    cargarDatos();
    const intervalo = setInterval(cargarDatos, 15000);
    return () => clearInterval(intervalo);
  }, []);

  async function cargarDatos() {
    try {
      const data = await apiFetch<EventoPendiente[]>("/api/ledger/pendientes", { rol: "Administrador" });
      setPendientes(data || []);
      setError("");
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
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
      second: "2-digit",
    });
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Sincronización con el Ledger</h1>
      <p className="text-gray-500 mb-6">
        Eventos que ocurrieron mientras la red Hyperledger Fabric no estaba disponible. Se
        reintenta sincronizarlos automáticamente cada 30 segundos, conservando la hora real
        en que ocurrió cada evento.
      </p>

      {cargando && <p className="text-gray-500">Cargando...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}

      {!cargando && !error && pendientes.length === 0 && (
        <div className="bg-green-50 border border-green-200 rounded-xl p-5 flex items-center gap-3">
          <span className="text-2xl">✓</span>
          <div>
            <p className="font-semibold text-green-800">Todo sincronizado</p>
            <p className="text-sm text-green-700">No hay eventos pendientes por registrar en el ledger.</p>
          </div>
        </div>
      )}

      {pendientes.length > 0 && (
        <>
          <div className="bg-amber-50 border border-amber-300 rounded-xl p-4 mb-4">
            <p className="text-sm text-amber-800">
              ⚠️ Hay <strong>{pendientes.length}</strong> evento(s) esperando sincronizar. Esto
              suele indicar que la red Hyperledger Fabric está caída en este momento.
            </p>
          </div>

          <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
            <table className="w-full text-sm">
              <thead className="bg-gray-50 text-gray-500 text-left">
                <tr>
                  <th className="px-4 py-3">Pedido</th>
                  <th className="px-4 py-3">Evento</th>
                  <th className="px-4 py-3">Ocurrió a las</th>
                  <th className="px-4 py-3">Intentos</th>
                  <th className="px-4 py-3">Último intento</th>
                </tr>
              </thead>
              <tbody>
                {pendientes.map((e) => (
                  <tr key={e.id} className="border-t border-gray-100">
                    <td className="px-4 py-3 font-mono text-gray-500">
                      {e.id_pedido.slice(0, 8).toUpperCase()}
                    </td>
                    <td className="px-4 py-3">
                      <Badge color="yellow">{e.estado}</Badge>
                    </td>
                    <td className="px-4 py-3">{formatearFecha(e.fecha_evento_real)}</td>
                    <td className="px-4 py-3">
                      {e.intentos > 0 ? (
                        <span className="text-amber-700 font-semibold">{e.intentos}</span>
                      ) : (
                        <span className="text-gray-400">0</span>
                      )}
                    </td>
                    <td className="px-4 py-3 text-gray-500">
                      {e.ultimo_intento ? formatearFecha(e.ultimo_intento) : "—"}
                    </td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        </>
      )}
    </div>
  );
}