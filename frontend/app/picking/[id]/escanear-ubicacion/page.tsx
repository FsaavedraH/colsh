"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import ScanBox from "@/components/ui/ScanBox";
import { apiFetch } from "@/lib/api";

export default function EscanearUbicacionPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [idProducto, setIdProducto] = useState("088b064e-fdb3-4032-b360-57174fcc60ee");
  const [resultado, setResultado] = useState<{ tipo: "ok" | "error"; mensaje: string } | null>(null);
  const [verificando, setVerificando] = useState(false);

  async function manejarEscaneo(ubicacionEscaneada: string) {
    if (verificando) return;
    setVerificando(true);
    setResultado(null);

    try {
      await apiFetch("/api/picking/escanear-ubicacion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: idProducto,
          ubicacion_escaneada: ubicacionEscaneada,
        }),
      });

      setResultado({ tipo: "ok", mensaje: `Ubicación correcta: ${ubicacionEscaneada}` });
      setTimeout(() => {
        router.push(`/picking/${idPedido}/escanear-producto`);
      }, 1200);
    } catch (err: any) {
      setResultado({ tipo: "error", mensaje: err.message });
    } finally {
      setVerificando(false);
    }
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-bold mb-1">Escanear ubicación</h1>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="mb-4">
        <label className="text-sm text-gray-500 block mb-1">
          ID de producto a validar (temporal, para pruebas)
        </label>
        <input
          type="text"
          value={idProducto}
          onChange={(e) => setIdProducto(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm font-mono"
        />
      </div>

      <ScanBox onScan={manejarEscaneo} />

      {resultado && (
        <div
          className={`mt-4 p-4 rounded-lg ${
            resultado.tipo === "ok" ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
          }`}
        >
          {resultado.mensaje}
        </div>
      )}
    </div>
  );
}