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

interface Pedido {
  productos: ItemPedido[];
}

export default function EscanearUbicacionPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [producto, setProducto] = useState<ItemPedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [errorCarga, setErrorCarga] = useState("");
  const [resultado, setResultado] = useState<{ tipo: "ok" | "error"; mensaje: string } | null>(null);
  const [verificando, setVerificando] = useState(false);

  useEffect(() => {
    apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Picking" })
      .then((data) => {
        // Nota: por ahora se asume 1 producto por pedido. El soporte para
        // pedidos con varios productos distintos queda pendiente.
        setProducto(data.productos?.[0] || null);
      })
      .catch((err) => setErrorCarga(err.message))
      .finally(() => setCargando(false));
  }, [idPedido]);

  async function manejarEscaneo(ubicacionEscaneada: string) {
    if (verificando || !producto) return;
    setVerificando(true);
    setResultado(null);

    try {
      await apiFetch("/api/picking/escanear-ubicacion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: producto.id_producto,
          ubicacion_escaneada: ubicacionEscaneada,
        }),
      });

      setResultado({ tipo: "ok", mensaje: "Ubicación correcta" });
      setTimeout(() => {
        router.push(`/picking/${idPedido}/escanear-producto`);
      }, 1200);
    } catch (err: any) {
      setResultado({ tipo: "error", mensaje: err.message });
    } finally {
      setVerificando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (errorCarga) return <p className="text-red-600">Error: {errorCarga}</p>;
  if (!producto) return <p className="text-gray-500">Este pedido no tiene productos registrados.</p>;

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-bold mb-1">Escanear ubicación</h1>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="bg-blue-50 border border-blue-200 rounded-lg p-4 mb-4">
        <p className="text-xs text-blue-600 mb-1">Producto a recolectar</p>
        <p className="font-semibold text-blue-900">{producto.nombre}</p>
        <p className="text-sm text-blue-700">Cantidad solicitada: {producto.cantidad}</p>
        <p className="text-sm text-blue-700 mt-1">
          Ubicación esperada: <span className="font-mono font-semibold">{producto.ubicacion}</span>
        </p>
      </div>

      <p className="text-sm text-gray-500 mb-2">
        Ve a esa ubicación en la bodega y escanea el código QR pegado en el estante.
      </p>

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