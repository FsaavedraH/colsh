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

  const productosSeleccionados = Object.entries(cantidades).filter(([, cant]) => cant > 0);
  const totalItems = productosSeleccionados.reduce((sum, [, cant]) => sum + cant, 0);

  // Solo un usuario con rol Cliente puede confirmar pedidos. Si hay sesion activa
  // pero con otro rol (ej. Admin navego al catalogo publico), no se le permite
  // comprar suplantando a un cliente.
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
    <div className={totalItems > 0 ? "pb-28" : ""}>
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

      {/* Barra flotante del carrito, fija abajo mientras se navega el catalogo */}
      {totalItems > 0 && (
        <div className="fixed bottom-0 left-0 right-0 md:left-56 bg-white border-t border-gray-200 shadow-lg z-40">
          {mostrarResumen && (
            <div className="border-b border-gray-100 p-4 max-w-md ml-auto mr-4">
              <label className="text-sm text-gray-500 block mb-1">Dirección de entrega</label>
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
              🛒 {totalItems} ítem{totalItems !== 1 ? "s" : ""} en tu pedido{" "}
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