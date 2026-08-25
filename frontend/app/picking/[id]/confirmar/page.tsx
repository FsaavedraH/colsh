"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Button from "@/components/ui/Button";

export default function ConfirmarRecoleccionPage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [idProducto, setIdProducto] = useState("088b064e-fdb3-4032-b360-57174fcc60ee");
  const [cantidad, setCantidad] = useState(1);
  const [responsable, setResponsable] = useState("68215789-6d0d-441c-ac90-eedf8aca4a61");
  const [confirmado, setConfirmado] = useState(false);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  async function confirmarRecoleccion() {
    setEnviando(true);
    setError("");

    try {
      await apiFetch("/api/recoleccion", {
        method: "POST",
        rol: "Picking",
        body: JSON.stringify({
          id_pedido: idPedido,
          id_producto: idProducto,
          cantidad: cantidad,
          responsable: responsable,
        }),
      });
      setConfirmado(true);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

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
          <label className="text-sm text-gray-500 block mb-1">ID de producto (temporal)</label>
          <input
            type="text"
            value={idProducto}
            onChange={(e) => setIdProducto(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm font-mono"
          />
        </div>

        <div>
          <label className="text-sm text-gray-500 block mb-1">Cantidad</label>
          <input
            type="number"
            value={cantidad}
            onChange={(e) => setCantidad(Number(e.target.value))}
            min={1}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
          />
        </div>

        <div>
          <label className="text-sm text-gray-500 block mb-1">Responsable (temporal)</label>
          <input
            type="text"
            value={responsable}
            onChange={(e) => setResponsable(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm font-mono"
          />
        </div>
      </div>

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      <Button onClick={confirmarRecoleccion} disabled={enviando}>
        {enviando ? "Confirmando..." : "Confirmar recolección"}
      </Button>
    </div>
  );
}