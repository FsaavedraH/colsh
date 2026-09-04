"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Badge from "@/components/ui/Badge";
import Button from "@/components/ui/Button";

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  id_cliente: string;
  direccion_entrega: string;
}

export default function DetalleDespachoPage() {
  const params = useParams();
  const router = useRouter();
  const { usuario } = useAuth();
  const idPedido = params.id as string;

  const [pedido, setPedido] = useState<Pedido | null>(null);
  const [codigoSeguimiento, setCodigoSeguimiento] = useState("");
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);
  const [entregado, setEntregado] = useState(false);

  useEffect(() => {
    cargarPedido();
  }, [idPedido]);

  async function cargarPedido() {
    setCargando(true);
    setError("");
    try {
      const data = await apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Transportista" });
      setPedido(data);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function generarDespacho() {
    if (!usuario) {
      setError("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }
    setEnviando(true);
    setError("");
    try {
      const data = await apiFetch<{ codigo_seguimiento: string }>("/api/despacho", {
        method: "POST",
        rol: "Transportista",
        body: JSON.stringify({
          id_pedido: idPedido,
          transportista: usuario.id_usuario,
        }),
      });
      setCodigoSeguimiento(data.codigo_seguimiento);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  async function confirmarEntrega() {
    if (!usuario) {
      setError("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }
    setEnviando(true);
    setError("");
    try {
      await apiFetch("/api/entrega", {
        method: "POST",
        rol: "Transportista",
        body: JSON.stringify({
          id_pedido: idPedido,
          transportista: usuario.id_usuario,
        }),
      });
      setEntregado(true);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando orden...</p>;
  if (error && !pedido) return <p className="text-red-600">Error: {error}</p>;
  if (!pedido) return <p className="text-gray-500">Orden no encontrada.</p>;

  if (entregado) {
    return (
      <div className="max-w-md">
        <div className="bg-green-100 text-green-700 rounded-xl p-6 text-center">
          <div className="text-3xl mb-2">✓</div>
          <h1 className="text-lg font-bold mb-1">¡Entrega confirmada!</h1>
          <p className="text-sm">El pedido fue entregado exitosamente.</p>
        </div>
        <Button onClick={() => router.push("/transportista")}>
          Volver a mis despachos
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-xl">
      <button
        onClick={() => router.push("/transportista")}
        className="text-sm text-gray-500 mb-4 hover:underline"
      >
        ← Volver a mis despachos
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

        <p className="text-sm text-gray-500 mb-1">Destino</p>
        <p className="font-semibold">{pedido.direccion_entrega}</p>
      </div>

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      {!codigoSeguimiento ? (
        <Button onClick={generarDespacho} disabled={enviando}>
          {enviando ? "Generando..." : "Generar orden de despacho"}
        </Button>
      ) : (
        <>
          <div className="bg-amber-50 border border-amber-200 rounded-lg p-4 mb-4">
            <p className="text-sm text-gray-500">Código de seguimiento</p>
            <p className="text-lg font-bold text-amber-700">{codigoSeguimiento}</p>
          </div>
          <Button onClick={confirmarEntrega} disabled={enviando}>
            {enviando ? "Confirmando..." : "Confirmar entrega"}
          </Button>
        </>
      )}
    </div>
  );
}