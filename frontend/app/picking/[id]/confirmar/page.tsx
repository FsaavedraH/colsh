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

interface Pedido {
  productos: ItemPedido[];
}

export default function ConfirmarRecoleccionPage() {
  const params = useParams();
  const router = useRouter();
  const { usuario } = useAuth();
  const idPedido = params.id as string;

  const [producto, setProducto] = useState<ItemPedido | null>(null);
  const [cargando, setCargando] = useState(true);
  const [errorCarga, setErrorCarga] = useState("");
  const [cantidadContada, setCantidadContada] = useState(0);
  const [confirmado, setConfirmado] = useState(false);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  useEffect(() => {
    apiFetch<Pedido>(`/api/pedidos/${idPedido}`, { rol: "Picking" })
      .then((data) => {
        const item = data.productos?.[0] || null;
        setProducto(item);
        if (item) setCantidadContada(item.cantidad);
      })
      .catch((err) => setErrorCarga(err.message))
      .finally(() => setCargando(false));
  }, [idPedido]);

  const noCoincide = producto ? cantidadContada !== producto.cantidad : false;

  async function confirmarRecoleccion() {
    if (!usuario || !producto) {
      setError("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }

    setEnviando(true);
    setError("");

    try {
      await apiFetch("/api/recoleccion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: producto.id_producto,
          cantidad: cantidadContada,
          responsable: usuario.id_usuario,
        }),
      });
      setConfirmado(true);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  if (cargando) return <p className="text-gray-500">Cargando pedido...</p>;
  if (errorCarga) return <p className="text-red-600">Error: {errorCarga}</p>;
  if (!producto) return <p className="text-gray-500">Este pedido no tiene productos registrados.</p>;

  if (confirmado) {
    return (
      <div className="max-w-md">
        <div className="bg-green-100 text-green-700 rounded-xl p-6 text-center">
          <div className="text-3xl mb-2">✓</div>
          <h1 className="text-lg font-bold mb-1">¡Ítem recolectado!</h1>
          <p className="text-sm">La recolección fue registrada correctamente.</p>
        </div>
        <Button onClick={() => router.push("/picking")}>
          Volver a la lista de órdenes
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-bold mb-1">Confirmar recolección</h1>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4 space-y-3">
        <div>
          <p className="text-sm text-gray-500">Producto</p>
          <p className="font-semibold">{producto.nombre}</p>
        </div>

        <div>
          <p className="text-sm text-gray-500">Cantidad solicitada en el pedido</p>
          <p className="font-semibold">{producto.cantidad}</p>
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
          ⚠️ La cantidad contada ({cantidadContada}) no coincide con la solicitada ({producto.cantidad}).
          Verifica antes de confirmar; esto quedará registrado como incidencia de inventario.
        </div>
      )}

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      <Button onClick={confirmarRecoleccion} disabled={enviando}>
        {enviando ? "Confirmando..." : "Confirmar recolección"}
      </Button>
    </div>
  );
}