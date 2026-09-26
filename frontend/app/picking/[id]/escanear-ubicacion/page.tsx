"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import ScanBox from "@/components/ui/ScanBox";
import { apiFetch } from "@/lib/api";

interface ItemPedido {
  id_producto: string;
  nombre: string;
  cantidad: number;
  ubicacion: string;
}

interface SiguienteItemResponse {
  completo: boolean;
  item?: ItemPedido;
  procesados: number;
  total_items: number;
}

export default function EscanearUbicacionPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [item, setItem] = useState<ItemPedido | null>(null);
  const [progreso, setProgreso] = useState({ procesados: 0, total: 0 });
  const [cargando, setCargando] = useState(true);
  const [errorCarga, setErrorCarga] = useState("");
  const [resultado, setResultado] = useState<{ tipo: "ok" | "error"; mensaje: string } | null>(null);
  const [verificando, setVerificando] = useState(false);

  useEffect(() => {
    cargarSiguienteItem();
  }, [idPedido]);

  async function cargarSiguienteItem() {
    setCargando(true);
    setResultado(null);
    try {
      const data = await apiFetch<SiguienteItemResponse>(
        `/api/picking/${idPedido}/siguiente-item`,
        { rol: "Picking" }
      );
      if (data.completo || !data.item) {
        router.push(`/picking/${idPedido}/confirmar`);
        return;
      }
      setItem(data.item);
      setProgreso({ procesados: data.procesados, total: data.total_items });
    } catch (err: any) {
      setErrorCarga(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function manejarEscaneo(ubicacionEscaneada: string) {
    if (verificando || !item) return;
    setVerificando(true);
    setResultado(null);

    try {
      await apiFetch("/api/picking/escanear-ubicacion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: item.id_producto,
          ubicacion_escaneada: ubicacionEscaneada,
        }),
      });

      setResultado({ tipo: "ok", mensaje: "Ubicación correcta" });
      setTimeout(() => {
        router.push(`/picking/${idPedido}/escanear-producto`);
      }, 400);
    } catch (err: any) {
      // Se deja el escaneo en pausa (ver "pausado" en ScanBox de abajo)
      // hasta que el operario retire el codigo y pulse "Escanear de nuevo",
      // para que el mensaje de error no se borre solo mientras el QR
      // incorrecto sigue frente a la camara.
      setResultado({ tipo: "error", mensaje: err.message });
    } finally {
      setVerificando(false);
    }
  }

  function reintentar() {
    setResultado(null);
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (errorCarga) return <p className="text-red-600">Error: {errorCarga}</p>;
  if (!item) return <p className="text-gray-500">Este pedido no tiene productos pendientes.</p>;

  const pausado = verificando || resultado !== null;

  return (
    <div className="max-w-md">
      <div className="flex justify-between items-center mb-1">
        <h1 className="text-xl font-bold">Escanear ubicación</h1>
        {progreso.total > 1 && (
          <span className="text-xs font-semibold text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
            Ítem {progreso.procesados + 1} de {progreso.total}
          </span>
        )}
      </div>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
        <p className="text-xs text-blue-600 mb-1">Producto a recolectar</p>
        <p className="font-semibold text-blue-900">{item.nombre}</p>
        <p className="text-sm text-blue-700">Cantidad solicitada: {item.cantidad}</p>
        <p className="text-sm text-blue-700 mt-1">
          Ubicación esperada: <span className="font-mono font-semibold">{item.ubicacion}</span>
        </p>
      </div>

      <p className="text-sm text-gray-500 mb-2">
        Ve a esa ubicación en la bodega y escanea el código QR pegado en el estante.
      </p>

      <ScanBox onScan={manejarEscaneo} pausado={pausado} />

      {resultado && (
        <div
          className={`mt-4 p-4 rounded-lg ${
            resultado.tipo === "ok" ? "bg-green-100 text-green-700" : "bg-red-100 text-red-700"
          }`}
        >
          <p>{resultado.mensaje}</p>
          {resultado.tipo === "error" && (
            <button
              onClick={reintentar}
              className="mt-2 text-sm font-semibold text-red-700 underline"
            >
              Escanear de nuevo
            </button>
          )}
        </div>
      )}
    </div>
  );
}