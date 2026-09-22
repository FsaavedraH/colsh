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
  costo_unitario: number;
}

const CARRITO_PENDIENTE_KEY = "colsh_carrito_pendiente";
const IVA = 0.19;

function formatearCOP(valor: number) {
  return valor.toLocaleString("es-CO", { style: "currency", currency: "COP", maximumFractionDigits: 0 });
}

function IconoRepuesto() {
  return (
    <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
      <path
        d="M14.7 6.3a4 4 0 0 0-5.4 4.6L3 17.2V21h3.8l6.3-6.3a4 4 0 0 0 4.6-5.4l-2.8 2.8-2-2 2.8-2.8Z"
        strokeLinecap="round"
        strokeLinejoin="round"
      />
    </svg>
  );
}

export default function CatalogoPage() {
  const router = useRouter();
  const { usuario } = useAuth();
  const [productos, setProductos] = useState<Producto[]>([]);
  const [cantidades, setCantidades] = useState<Record<string, number>>({});
  const [direccion, setDireccion] = useState("Cra 50 #10-25");
  const [cargando, setCargando] = useState(true);
  const [enviando, setEnviando] = useState(false);
  const [error, setError] = useState("");
  const [mostrarResumen, setMostrarResumen] = useState(false);

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
        setMostrarResumen(true);
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

  function precioDe(idProd: string) {
    return productos.find((p) => p.id_producto === idProd)?.costo_unitario ?? 0;
  }

  function nombreDe(idProd: string) {
    return productos.find((p) => p.id_producto === idProd)?.nombre ?? "";
  }

  const productosSeleccionados = Object.entries(cantidades).filter(([, cant]) => cant > 0);
  const totalItems = productosSeleccionados.reduce((sum, [, cant]) => sum + cant, 0);
  const subtotal = productosSeleccionados.reduce(
    (sum, [id, cant]) => sum + precioDe(id) * cant,
    0
  );
  const iva = subtotal * IVA;
  const total = subtotal + iva;

  const esRolNoCliente = usuario !== null && usuario.rol !== "Cliente";

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

    if (esRolNoCliente) {
      setError(
        "Tu sesión actual no es de Cliente, así que no puedes confirmar pedidos con ella. Cierra sesión e inicia con una cuenta de Cliente."
      );
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
    <div className={totalItems > 0 ? "pb-32" : ""}>
      <h1 className="text-2xl font-bold mb-1">Catálogo de Repuestos</h1>
      <p className="text-gray-500 mb-6">
        {cargando ? "Cargando catálogo..." : `${productos.length} repuestos disponibles`}
      </p>

      {error && <p className="text-red-600 mb-4">Error: {error}</p>}

      <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-4 mb-6">
        {productos.map((p) => (
          <div
            key={p.id_producto}
            className="bg-white rounded-xl border border-gray-200 p-4 hover:shadow-sm transition-shadow"
          >
            <div className="flex items-start gap-3 mb-3">
              <div className="w-11 h-11 rounded-lg bg-gray-100 text-gray-500 flex items-center justify-center shrink-0">
                <IconoRepuesto />
              </div>
              <div className="min-w-0">
                <div className="font-semibold text-gray-800 leading-snug">{p.nombre}</div>
                <div className="text-xs text-gray-400">Stock: {p.stock}</div>
              </div>
            </div>

            <div className="flex items-center justify-between">
              <span className="font-semibold text-gray-900">{formatearCOP(p.costo_unitario)}</span>
              <div className="flex items-center gap-2">
                <button
                  onClick={() => cambiarCantidad(p.id_producto, -1)}
                  className="w-8 h-8 rounded-lg bg-gray-100 hover:bg-gray-200 font-bold"
                  aria-label={`Quitar ${p.nombre}`}
                >
                  −
                </button>
                <span className="w-6 text-center font-semibold">
                  {cantidades[p.id_producto] || 0}
                </span>
                <button
                  onClick={() => cambiarCantidad(p.id_producto, 1)}
                  disabled={p.stock === 0}
                  className="w-8 h-8 rounded-lg bg-blue-100 hover:bg-blue-200 text-blue-700 font-bold disabled:opacity-40 disabled:cursor-not-allowed"
                  aria-label={`Agregar ${p.nombre}`}
                >
                  +
                </button>
              </div>
            </div>
          </div>
        ))}
      </div>

      {totalItems > 0 && (
        <div className="fixed bottom-0 left-0 right-0 md:left-56 bg-white border-t border-gray-200 shadow-lg z-40">
          {mostrarResumen && (
            <div className="border-b border-gray-100 p-4 max-w-md ml-auto mr-4 max-h-64 overflow-y-auto">
              <p className="text-xs font-semibold text-gray-500 mb-2">Resumen del pedido</p>
              <div className="space-y-1 mb-3">
                {productosSeleccionados.map(([id, cant]) => (
                  <div key={id} className="flex justify-between text-sm text-gray-600">
                    <span className="truncate pr-2">
                      {nombreDe(id)} x{cant}
                    </span>
                    <span className="shrink-0">{formatearCOP(precioDe(id) * cant)}</span>
                  </div>
                ))}
              </div>
              <div className="border-t border-gray-100 pt-2 space-y-1 text-sm">
                <div className="flex justify-between text-gray-500">
                  <span>Subtotal</span>
                  <span>{formatearCOP(subtotal)}</span>
                </div>
                <div className="flex justify-between text-gray-500">
                  <span>IVA (19%)</span>
                  <span>{formatearCOP(iva)}</span>
                </div>
                <div className="flex justify-between font-semibold text-gray-900 text-base pt-1">
                  <span>Total</span>
                  <span>{formatearCOP(total)}</span>
                </div>
              </div>

              <label className="text-sm text-gray-500 block mb-1 mt-4">Dirección de entrega</label>
              <input
                type="text"
                value={direccion}
                onChange={(e) => setDireccion(e.target.value)}
                className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm"
              />
              {!usuario && (
                <p className="text-xs text-amber-600 mt-2">
                  Necesitas iniciar sesión para confirmar tu pedido.
                </p>
              )}
              {esRolNoCliente && (
                <p className="text-xs text-red-600 mt-2">
                  Tu sesión actual ({usuario?.rol}) no puede confirmar pedidos. Cierra sesión e
                  inicia con una cuenta de Cliente.
                </p>
              )}
            </div>
          )}

          <div className="flex items-center justify-between px-4 py-3 max-w-4xl mx-auto">
            <button
              onClick={() => setMostrarResumen((v) => !v)}
              className="text-sm font-medium text-gray-700 hover:text-gray-900"
            >
              🛒 {totalItems} ítem{totalItems !== 1 ? "s" : ""} · {formatearCOP(total)}{" "}
              <span className="text-gray-400">{mostrarResumen ? "▾" : "▸"}</span>
            </button>

            <Button onClick={crearPedido} disabled={enviando || esRolNoCliente}>
              {enviando
                ? "Creando pedido..."
                : !usuario
                ? "Iniciar sesión para continuar"
                : esRolNoCliente
                ? "No disponible con este rol"
                : "Confirmar pedido"}
            </Button>
          </div>
        </div>
      )}
    </div>
  );
}