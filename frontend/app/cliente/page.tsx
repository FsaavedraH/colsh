"use client";

import { useEffect, useState } from "react";
import { useRouter } from "next/navigation";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";

interface Producto {
  id_producto: string;
  nombre: string;
  stock: number;
  ubicacion: string;
}

const CARRITO_PENDIENTE_KEY = "colsh_carrito_pendiente";

export default function CatalogoPage() {
  const router = useRouter();
  const { usuario } = useAuth();
  const [productos, setProductos] = useState<Producto[]>([]);
  const [cantidades, setCantidades] = useState<Record<string, number>>({});
  const [direccion, setDireccion] = useState("Cra 50 #10-25");
  const [cargando, setCargando] = useState(true);
  const [enviando, setEnviando] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    cargarCatalogo();
    restaurarCarritoPendiente();
  }, []);

  function restaurarCarritoPendiente() {
    const guardado = sessionStorage.getItem(CARRITO_PENDIENTE_KEY);
    if (guardado) {
      try {
        const { cantidades: c, direccion: d } = JSON.parse(guardado);
        setCantidades(c || {});
        setDireccion(d || "Cra 50 #10-25");
      } catch {
        // si el JSON guardado esta corrupto, simplemente lo ignoramos
      }
      sessionStorage.removeItem(CARRITO_PENDIENTE_KEY);
    }
  }

  async function cargarCatalogo() {
    setCargando(true);
    try {
      const data = await apiFetch<Producto[]>("/api/productos", { rol: "Cliente" });
      setProductos(data || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  function cambiarCantidad(idProducto: string, delta: number) {
    setCantidades((prev) => {
      const actual = prev[idProducto] || 0;
      const nueva = Math.max(0, actual + delta);
      return { ...prev, [idProducto]: nueva };
    });
  }

  const productosSeleccionados = Object.entries(cantidades).filter(([, cant]) => cant > 0);
  const totalItems = productosSeleccionados.reduce((sum, [, cant]) => sum + cant, 0);

  async function crearPedido() {
    if (productosSeleccionados.length === 0) return;

    if (!usuario) {
      sessionStorage.setItem(
        CARRITO_PENDIENTE_KEY,
        JSON.stringify({ cantidades, direccion })
      );
      router.push("/login");
      return;
    }

    setEnviando(true);
    setError("");

    try {
      const data = await apiFetch<{ id_pedido: string }>("/api/pedidos", {
        method: "POST",
        rol: "Cliente",
        body: JSON.stringify({
          cliente_id: usuario.id_usuario,
          direccion_entrega: direccion,
          productos: productosSeleccionados.map(([id_producto, cantidad]) => ({
            id_producto,
            cantidad,
          })),
        }),
      });
      router.push(`/cliente/pedidos/${data.id_pedido}`);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setEnviando(false);
    }
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Catálogo de Repuestos</h1>

      {cargando && <p className="text-gray-500">Cargando catálogo...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}

      <div className="grid md:grid-cols-3 gap-4 mb-6">
        {productos.map((p) => (
          <div key={p.id_producto} className="bg-white rounded-xl border border-gray-200 p-4">
            <div className="font-semibold text-gray-800 mb-1">{p.nombre}</div>
            <div className="text-xs text-gray-400 mb-3">Stock: {p.stock}</div>
            <div className="flex items-center gap-2">
              <button
                onClick={() => cambiarCantidad(p.id_producto, -1)}
                className="w-8 h-8 rounded-lg bg-gray-100 hover:bg-gray-200 font-bold"
              >
                −
              </button>
              <span className="w-8 text-center font-semibold">
                {cantidades[p.id_producto] || 0}
              </span>
              <button
                onClick={() => cambiarCantidad(p.id_producto, 1)}
                className="w-8 h-8 rounded-lg bg-blue-100 hover:bg-blue-200 text-blue-700 font-bold"
              >
                +
              </button>
            </div>
          </div>
        ))}
      </div>

      {totalItems > 0 && (
        <div className="bg-white rounded-xl border border-gray-200 p-5 max-w-md">
          <h2 className="font-semibold mb-3">Resumen del pedido ({totalItems} ítems)</h2>

          <label className="text-sm text-gray-500 block mb-1">Dirección de entrega</label>
          <input
            type="text"
            value={direccion}
            onChange={(e) => setDireccion(e.target.value)}
            className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
          />

          {!usuario && (
            <p className="text-xs text-amber-600 mb-3">
              Necesitas iniciar sesión para confirmar tu pedido.
            </p>
          )}

          <Button onClick={crearPedido} disabled={enviando}>
            {enviando
              ? "Creando pedido..."
              : usuario
              ? "Confirmar pedido"
              : "Iniciar sesión para continuar"}
          </Button>
        </div>
      )}
    </div>
  );
}