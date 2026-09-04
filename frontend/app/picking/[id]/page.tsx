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

export default function DetalleOrdenPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [pedido, setPedido] = useState<Pedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [iniciando, setIniciando] = useState(false);

  useEffect(() => {
    cargarPedido();
  }, [idPedido]);

  async function cargarPedido() {
    setCargando(true);
    setError("");
    try {
      const data = await apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Picking" });
      setPedido(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function iniciarPicking() {
    setIniciando(true);
    setError("");
    try {
      // Si el pedido esta "Pendiente", esta llamada lo pasa a "En recoleccion".
      // Si ya estaba "En recoleccion" (se reabrio la orden), simplemente lo confirma de nuevo.
      await apiFetch("/api/picking/iniciar", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({ id_pedido: idPedido }),
      });
      router.push(`/picking/${idPedido}/escanear-ubicacion`);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setIniciando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando orden...</p>;
  if (error) return <p className="text-red-600">Error: {error}</p>;
  if (!pedido) return <p className="text-gray-500">Orden no encontrada.</p>;

  return (
    <div className="max-w-2xl">
      <button
        onClick={() => router.push("/picking")}
        className="text-sm text-gray-500 mb-4 hover:underline"
      >
        ← Volver a la lista
      </button>

      <div className="flex justify-between items-start mb-6">
        <div>
          <h1 className="text-xl font-bold">Orden {pedido.id_pedido.slice(0, 8).toUpperCase()}</h1>
          <p className="text-gray-500 text-sm">{pedido.direccion_entrega}</p>
        </div>
        <Badge color="yellow">{pedido.estado}</Badge>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <p className="text-sm text-gray-500 mb-1">Cliente</p>
        <p className="font-semibold mb-4">{pedido.id_cliente}</p>

        <p className="text-sm text-gray-500 mb-1">Fecha de creación</p>
        <p className="font-semibold">
          {new Date(pedido.fecha_creacion).toLocaleString("es-CO")}
        </p>
      </div>

      <Button onClick={iniciarPicking} disabled={iniciando}>
        {iniciando ? "Iniciando..." : "Iniciar picking →"}
      </Button>
    </div>
  );
}