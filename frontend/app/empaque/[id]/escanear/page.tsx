"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import ScanBox from "@/components/ui/ScanBox";
import { apiFetch } from "@/lib/api";

export default function EscanearEmpaquePage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [idProductoEsperado, setIdProductoEsperado] = useState("088b064e-fdb3-4032-b360-57174fcc60ee");
  const [resultado, setResultado] = useState<{ tipo: "ok" | "error"; mensaje: string } | null>(null);
  const [verificando, setVerificando] = useState(false);

  async function manejarEscaneo(idProductoEscaneado: string) {
    if (verificando) return;
    setVerificando(true);
    setResultado(null);

    try {
      await apiFetch("/api/empaque/escanear", {
        method: "POST",
        rol: "Empaque",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto_esperado: idProductoEsperado,
          id_producto_escaneado: idProductoEscaneado,
        }),
      });

      setResultado({ tipo: "ok", mensaje: "Producto validado correctamente" });
      setTimeout(() => {
        router.push(`/empaque/${idPedido}/confirmar`);
      }, 1200);
    } catch (err: any) {
      setResultado({ tipo: "error", mensaje: err.message });
    } finally {
      setVerificando(false);
    }
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-bold mb-1">Validar producto para empaque</h1>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="mb-4">
        <label className="text-sm text-gray-500 block mb-1">
          ID de producto esperado (temporal, para pruebas)
        </label>
        <input
          type="text"
          value={idProductoEsperado}
          onChange={(e) => setIdProductoEsperado(e.target.value)}
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