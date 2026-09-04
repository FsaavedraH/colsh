"use client";

import { useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";

export default function ConfirmarEmpaquePage() {
  const params = useParams();
  const router = useRouter();
  const { usuario } = useAuth();
  const idPedido = params.id as string;

  const [confirmado, setConfirmado] = useState(false);
  const [error, setError] = useState("");
  const [enviando, setEnviando] = useState(false);

  async function confirmarEmpaque() {
    if (!usuario) {
      setError("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }

    setEnviando(true);
    setError("");

    try {
      await apiFetch("/api/empaque", {
        method: "POST",
        rol: "Empaque",
        body: JSON.stringify({
          id_pedido: idPedido,
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
        <label className="text-sm text-gray-500 block mb-1">Responsable</label>
        <p className="text-sm font-medium text-gray-700">{usuario?.nombre || "—"}</p>
      </div>

      {error && <div className="mb-4 p-3 bg-red-100 text-red-700 rounded-lg text-sm">{error}</div>}

      <Button onClick={confirmarEmpaque} disabled={enviando}>
        {enviando ? "Confirmando..." : "Confirmar empaque"}
      </Button>
    </div>
  );
}