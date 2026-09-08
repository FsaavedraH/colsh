"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";

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

interface ConfirmarRecoleccionResponse {
  estado: string;
  completo: boolean;
  pendientes?: number;
}

export default function ConfirmarRecoleccionPage() {
  const params = useParams();
  const router = useRouter();
  const { usuario } = useAuth();
  const idPedido = params.id as string;

  const [item, setItem] = useState<ItemPedido | null>(null);
  const [progreso, setProgreso] = useState({ procesados: 0, total: 0 });
  const [cargando, setCargando] = useState(true);
  const [errorCarga, setErrorCarga] = useState("");
  const [cantidadContada, setCantidadContada] = useState(0);
  const [pedidoCompleto, setPedidoCompleto] = useState(false);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  useEffect(() => {
    cargarSiguienteItem();
  }, [idPedido]);

  async function cargarSiguienteItem() {
    setCargando(true);
    try {
      const data = await apiFetch<SiguienteItemResponse>(
        `/api/picking/${idPedido}/siguiente-item`,
        { rol: "Picking" }
      );
      if (data.completo || !data.item) {
        setPedidoCompleto(true);
        setItem(null);
      } else {
        setItem(data.item);
        setCantidadContada(data.item.cantidad);
        setProgreso({ procesados: data.procesados, total: data.total_items });
      }
    } catch (err: any) {
      setErrorCarga(err.message);
    } finally {
      setCargando(false);
    }
  }

  const noCoincide = item ? cantidadContada !== item.cantidad : false;

  async function confirmarRecoleccion() {
    if (!usuario || !item) {
      setError("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }

    setEnviando(true);
    setError("");

    try {
      const data = await apiFetch<ConfirmarRecoleccionResponse>("/api/recoleccion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: item.id_producto,
          cantidad: cantidadContada,
          responsable: usuario.id_usuario,
        }),
      });

      if (data.completo) {
        setPedidoCompleto(true);
        setItem(null);
      } else {
        // Quedan mas items: volvemos a escanear ubicacion para el siguiente producto
        router.push(`/picking/${idPedido}/escanear-ubicacion`);
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (errorCarga) return <p className="text-red-600">Error: {errorCarga}</p>;

  if (pedidoCompleto) {
    return (
      <div className="max-w-md">
        <div className="bg-green-100 text-green-700 rounded-xl p-6 text-center">
          <div className="text-3xl mb-2">✓</div>
          <h1 className="text-lg font-bold mb-1">¡Pedido recolectado por completo!</h1>
          <p className="text-sm">Todos los productos fueron recolectados. La orden pasó a empaque.</p>
        </div>
        <Button onClick={() => router.push("/picking")}>
          Volver a la lista de órdenes
        </Button>
      </div>
    );
  }

  if (!item) return <p className="text-gray-500">Este pedido no tiene productos pendientes.</p>;

  return (
    <div className="max-w-md">
      <div className="flex justify-between items-center mb-1">
        <h1 className="text-xl font-bold">Confirmar recolección</h1>
        {progreso.total > 1 && (
          <span className="text-xs font-semibold text-gray-500 bg-gray-100 px-2 py-1 rounded-full">
            Ítem {progreso.procesados + 1} de {progreso.total}
          </span>
        )}
      </div>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4 space-y-3">
        <div>
          <p className="text-sm text-gray-500">Producto</p>
          <p className="font-semibold">{item.nombre}</p>
        </div>

        <div>
          <p className="text-sm text-gray-500">Cantidad solicitada en el pedido</p>
          <p className="font-semibold">{item.cantidad}</p>
        </div>

        <div>
          <label className="text-sm text-gray-500 block mb-1">
            Cantidad que contaste físicamente
          </label>
          <input
            type="number"
            value={cantidadContada}
            onChange={(e) => setCantidadContada(Number(e.target.value))}
            min={0}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
          />
        </div>

        <div>
          <label className="text-sm text-gray-500 block mb-1">Responsable</label>
          <p className="text-sm font-medium text-gray-700">{usuario?.nombre || "—"}</p>
        </div>
      </div>

      {noCoincide && (
        <div className="mb-4 p-3 bg-amber-100 text-amber-800 rounded-lg text-sm">
          ⚠️ La cantidad contada ({cantidadContada}) no coincide con la solicitada ({item.cantidad}).
          Verifica antes de confirmar; esto quedará registrado como incidencia de inventario.
        </div>
      )}

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      <Button onClick={confirmarRecoleccion} disabled={enviando}>
        {enviando
          ? "Confirmando..."
          : progreso.total > 1 && progreso.procesados + 1 < progreso.total
          ? "Confirmar y continuar con el siguiente ítem"
          : "Confirmar recolección"}
      </Button>
    </div>
  );
}