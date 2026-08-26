"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";
import Button from "@/components/ui/Button";

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  id_cliente: string;
  direccion_entrega: string;
}

export default function DetalleEmpaquePage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [pedido, setPedido] = useState<Pedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [recibido, setRecibido] = useState(false);

  useEffect(() => {
    cargarPedido();
  }, [idPedido]);

  async function cargarPedido() {
    setCargando(true);
    setError("");
    try {
      const data = await apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Empaque" });
      setPedido(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function confirmarRecepcion() {
    setError("");
    try {
      await apiFetch("/api/empaque/recepcion", {
        method: "POST",
        rol: "Empaque",
        body: JSON.stringify({ id_pedido: idPedido }),
      });
      setRecibido(true);
    } catch (err: any) {
      setError(err.message);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando orden...</p>;
  if (error && !pedido) return <p className="text-red-600">Error: {error}</p>;
  if (!pedido) return <p className="text-gray-500">Orden no encontrada.</p>;

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => router.push("/empaque")}
        className="text-sm text-gray-500 mb-4 hover:underline"
      >
        ← Volver a la lista
      </button>

      <div className="flex justify-between items-start mb-6">
        <div>
          <h1 className="text-xl font-bold">Orden {pedido.id_pedido.slice(0, 8).toUpperCase()}</h1>
          <p className="text-gray-500 text-sm">{pedido.direccion_entrega}</p>
        </div>
        <Badge color="green">{pedido.estado}</Badge>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <p className="text-sm text-gray-500 mb-1">Cliente</p>
        <p className="font-semibold mb-4">{pedido.id_cliente}</p>

        <p className="text-sm text-gray-500 mb-1">Fecha de creación</p>
        <p className="font-semibold">
          {new Date(pedido.fecha_creacion).toLocaleString("es-CO")}
        </p>
      </div>

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      {!recibido ? (
        <Button onClick={confirmarRecepcion}>Confirmar recepción en empaque</Button>
      ) : (
        <Button onClick={() => router.push(`/empaque/${idPedido}/escanear`)}>
          Iniciar validación →
        </Button>
      )}
    </div>
  );
}