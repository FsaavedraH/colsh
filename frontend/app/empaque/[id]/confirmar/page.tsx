"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import Button from "@/components/ui/Button";

export default function ConfirmarEmpaquePage() {
  const params = useParams();
  const router = useRouter();
  const idPedido = params.id as string;

  const [responsable, setResponsable] = useState("e83f286d-7f7c-4f6d-ae41-b3c9bf63df14");
  const [confirmado, setConfirmado] = useState(false);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  async function confirmarEmpaque() {
    setEnviando(true);
    setError("");

    try {
      await apiFetch("/api/empaque", {
        method: "POST",
        rol: "Empaque",
        body: JSON.stringify({
          id_pedido: idPedido,
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
          <h1 className="text-lg font-bold mb-1">¡Orden empaquetada!</h1>
          <p className="text-sm">La orden está lista para despacho.</p>
        </div>
        <Button onClick={() => router.push("/empaque")}>
          Volver a la lista de órdenes
        </Button>
      </div>
    );
  }

  return (
    <div className="max-w-md">
      <h1 className="text-xl font-bold mb-1">Registro de empaque</h1>
      <p className="text-gray-500 text-sm mb-4">
        Orden {idPedido.slice(0, 8).toUpperCase()}
      </p>

      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-4">
        <label className="text-sm text-gray-500 block mb-1">Responsable (temporal)</label>
        <input
          type="text"
          value={responsable}
          onChange={(e) => setResponsable(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm font-mono"
        />
      </div>

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      <Button onClick={confirmarEmpaque} disabled={enviando}>
        {enviando ? "Confirmando..." : "Confirmar empaque"}
      </Button>
    </div>
  );
}