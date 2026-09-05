"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Badge from "@/components/ui/Badge";

interface ItemPedido {
  nombre: string;
  cantidad: number;
}

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  id_cliente: string;
  direccion_entrega: string;
  productos: ItemPedido[];
}

const ETAPAS = ["Pendiente", "En recoleccion", "En empaque", "En despacho", "Entregado"];
const ESTADOS_CANCELABLES = ["Pendiente", "En espera por inventario", "En recoleccion", "En empaque"];

export default function ConsultaPedidoPage() {
  const params = useParams();
  const router = useRouter();
  const { usuario } = useAuth();
  const idPedido = params.id as string;

  const [pedido, setPedido] = useState<Pedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [cancelando, setCancelando] = useState(false);
  const [errorCancelar, setErrorCancelar] = useState("");
  const [confirmandoCancelacion, setConfirmandoCancelacion] = useState(false);

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

  async function cancelarPedido() {
    if (!usuario) return;
    setCancelando(true);
    setErrorCancelar("");
    try {
      await apiFetch(`/api/pedidos/${idPedido}/cancelar`, {
        method: "POST",
        rol: "Cliente",
        body: JSON.stringify({ cliente_id: usuario.id_usuario }),
      });
      setConfirmandoCancelacion(false);
      await cargarPedido();
    } catch (err: any) {
      setErrorCancelar(err.message);
    } finally {
      setCancelando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (error) return <p className="text-red-600">Error: {error}</p>;
  if (!pedido) return <p className="text-gray-500">Pedido no encontrado.</p>;

  const indiceEtapaActual = ETAPAS.indexOf(pedido.estado);
  const totalItems = (pedido.productos || []).reduce((sum, p) => sum + p.cantidad, 0);
  const enEsperaPorInventario = pedido.estado === "En espera por inventario";
  const esCancelable = ESTADOS_CANCELABLES.includes(pedido.estado);
  const estaCancelado = pedido.estado === "Cancelado";

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
          <div className="text-3xl mb-2">{estaCancelado ? "✕" : "✓"}</div>
          <h1 className="text-lg font-bold">
            {estaCancelado ? "Pedido cancelado" : "¡Pedido creado con éxito!"}
          </h1>
          <p className="text-sm text-gray-500">
            Número de pedido: {pedido.id_pedido.slice(0, 8).toUpperCase()}
          </p>
        </div>

        {!estaCancelado && (
          <>
            <div className="flex justify-between items-center">
              {ETAPAS.map((etapa, i) => {
                const esPasoUno = i === 0;
                const activo = enEsperaPorInventario ? esPasoUno : i <= indiceEtapaActual;
                const colorCirculo =
                  enEsperaPorInventario && esPasoUno
                    ? "bg-amber-500 text-white"
                    : activo
                    ? "bg-blue-600 text-white"
                    : "bg-gray-200 text-gray-400";

                return (
                  <div key={etapa} className="flex-1 flex flex-col items-center">
                    <div className={`w-6 h-6 rounded-full flex items-center justify-center text-xs font-bold ${colorCirculo}`}>
                      {enEsperaPorInventario && esPasoUno ? "!" : i + 1}
                    </div>
                    <span className="text-[10px] text-gray-500 mt-1 text-center">{etapa}</span>
                  </div>
                );
              })}
            </div>
            {enEsperaPorInventario && (
              <p className="text-xs text-amber-600 text-center mt-3">
                Esperando disponibilidad de stock para continuar con la recolección.
              </p>
            )}
          </>
        )}
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <p className="text-sm text-gray-500 mb-1">Estado actual</p>
        <Badge color={estaCancelado ? "red" : "yellow"}>{pedido.estado}</Badge>

        <p className="text-sm text-gray-500 mt-4 mb-1">Dirección de entrega</p>
        <p className="font-semibold">{pedido.direccion_entrega}</p>

        <p className="text-sm text-gray-500 mt-4 mb-1">Fecha de creación</p>
        <p className="font-semibold">{new Date(pedido.fecha_creacion).toLocaleString("es-CO")}</p>
      </div>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <p className="text-sm text-gray-500 mb-3">
          Detalle del pedido ({totalItems} ítem{totalItems !== 1 ? "s" : ""})
        </p>
        <div className="divide-y divide-gray-100">
          {(pedido.productos || []).map((item, i) => (
            <div key={i} className="flex justify-between py-2 text-sm">
              <span>{item.nombre}</span>
              <span className="font-semibold">x{item.cantidad}</span>
            </div>
          ))}
        </div>
      </div>

      {esCancelable && !confirmandoCancelacion && (
        <button
          onClick={() => setConfirmandoCancelacion(true)}
          className="text-sm text-red-600 hover:underline"
        >
          Cancelar pedido
        </button>
      )}

      {confirmandoCancelacion && (
        <div className="bg-red-50 border border-red-200 rounded-xl p-4">
          <p className="text-sm text-red-800 mb-3">
            ¿Seguro que quieres cancelar este pedido? Esta acción no se puede deshacer.
          </p>
          {errorCancelar && (
            <p className="text-sm text-red-600 mb-3">{errorCancelar}</p>
          )}
          <div className="flex gap-3">
            <button
              onClick={cancelarPedido}
              disabled={cancelando}
              className="text-sm font-semibold text-white bg-red-600 hover:bg-red-700 rounded-lg px-4 py-2"
            >
              {cancelando ? "Cancelando..." : "Sí, cancelar"}
            </button>
            <button
              onClick={() => setConfirmandoCancelacion(false)}
              className="text-sm text-gray-600 hover:underline"
            >
              No, mantener pedido
            </button>
          </div>
        </div>
      )}
    </div>
  );
}