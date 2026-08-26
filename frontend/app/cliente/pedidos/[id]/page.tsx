"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  id_cliente: string;
  direccion_entrega: string;
}

const ETAPAS = ["Pendiente", "En recoleccion", "En empaque", "En despacho", "Entregado"];

export default function ConsultaPedidoPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [pedido, setPedido] = useState<Pedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  useEffect(() => {
    cargarPedido();
  }, [idPedido]);

  async function cargarPedido() {
    setCargando(true);
    setError("");
    try {
      const data = await apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Cliente" });
      setPedido(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (error) return <p className="text-red-600">Error: {error}</p>;
  if (!pedido) return <p className="text-gray-500">Pedido no encontrado.</p>;

  const indiceEtapaActual = ETAPAS.indexOf(pedido.estado);

  return (
    <div className="max-w-xl">
      <button
        onClick={() => router.push("/cliente")}
        className="text-sm text-gray-500 mb-4 hover:underline"
      >
        ← Volver al catálogo
      </button>

      <div className="bg-white rounded-xl border border-gray-200 p-6 mb-6">
        <div className="text-center mb-6">
          <div className="text-3xl mb-2">✓</div>
          <h1 className="text-lg font-bold">¡Pedido creado con éxito!</h1>
          <p className="text-sm text-gray-500">
            Número de pedido: {pedido.id_pedido.slice(0, 8).toUpperCase()}
          </p>
        </div>

        <div className="flex justify-between items-center">
          {ETAPAS.map((etapa, i) => (
            <div key={etapa} className="flex-1 flex flex-col items-center">
              <div
                className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${
                  i <= indiceEtapaActual ? "bg-blue-600 text-white" : "bg-gray-200 text-gray-400"
                }`}
              >
                {i + 1}
              </div>
              <span className="text-[10px] text-gray-500 mt-1 text-center">{etapa}</span>
            </div>
          ))}
        </div>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5">
        <p className="text-sm text-gray-500 mb-1">Estado actual</p>
        <Badge color="yellow">{pedido.estado}</Badge>

        <p className="text-sm text-gray-500 mt-4 mb-1">Dirección de entrega</p>
        <p className="font-semibold">{pedido.direccion_entrega}</p>

        <p className="text-sm text-gray-500 mt-4 mb-1">Fecha de creación</p>
        <p className="font-semibold">{new Date(pedido.fecha_creacion).toLocaleString("es-CO")}</p>
      </div>
    </div>
  );
}