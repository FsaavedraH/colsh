"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import { useAuth } from "@/lib/auth";
import Button from "@/components/ui/Button";

interface Producto {
  id_producto: string;
  nombre: string;
  stock: number;
  ubicacion: string;
}

interface Compra {
  id_compra: string;
  id_producto: string;
  nombre_producto: string;
  cantidad: number;
  costo_unitario_momento: number | null;
  costo_total: number | null;
  fecha: string;
  nombre_responsable: string;
}

function formatearCOP(valor: number | null) {
  if (valor === null) return "—";
  return valor.toLocaleString("es-CO", { style: "currency", currency: "COP", maximumFractionDigits: 0 });
}

export default function IngresoComprasPage() {
  const { usuario } = useAuth();
  const [productos, setProductos] = useState<Producto[]>([]);
  const [compras, setCompras] = useState<Compra[]>([]);
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  const [idProducto, setIdProducto] = useState("");
  const [cantidad, setCantidad] = useState(1);
  const [registrando, setRegistrando] = useState(false);
  const [errorFormulario, setErrorFormulario] = useState("");

  useEffect(() => {
    cargarDatos();
  }, []);

  async function cargarDatos() {
    setCargando(true);
    setError("");
    try {
      const [productosData, comprasData] = await Promise.all([
        apiFetch<Producto[]>("/api/productos", { rol: "Administrador" }),
        apiFetch<Compra[]>("/api/inventario/compras", { rol: "Administrador" }),
      ]);
      setProductos(productosData || []);
      setCompras(comprasData || []);
      if (productosData && productosData.length > 0 && !idProducto) {
        setIdProducto(productosData[0].id_producto);
      }
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function registrarCompra(e: React.FormEvent) {
    e.preventDefault();
    setErrorFormulario("");

    if (!usuario) {
      setErrorFormulario("No se pudo identificar al usuario. Inicia sesión de nuevo.");
      return;
    }
    if (!idProducto || cantidad <= 0) {
      setErrorFormulario("Selecciona un producto y una cantidad válida.");
      return;
    }

    setRegistrando(true);
    try {
      await apiFetch("/api/inventario/compras", {
        method: "POST",
        rol: "Administrador",
        body: JSON.stringify({
          id_producto: idProducto,
          cantidad: cantidad,
          responsable: usuario.id_usuario,
        }),
      });
      setCantidad(1);
      await cargarDatos();
    } catch (err: any) {
      setErrorFormulario(err.message);
    } finally {
      setRegistrando(false);
    }
  }

  const productoSeleccionado = productos.find((p) => p.id_producto === idProducto);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Ingreso de compras</h1>
      <p className="text-gray-500 mb-6">
        Registra la entrada de repuestos comprados. Solo suma stock a productos existentes.
      </p>

      <form
        onSubmit={registrarCompra}
        className="bg-white rounded-xl border border-gray-200 p-5 mb-8 max-w-lg"
      >
        <h2 className="font-semibold mb-4">Nueva compra</h2>

        <label className="text-sm text-gray-500 block mb-1">Producto</label>
        <select
          value={idProducto}
          onChange={(e) => setIdProducto(e.target.value)}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-1"
        >
          {productos.map((p) => (
            <option key={p.id_producto} value={p.id_producto}>
              {p.nombre} (stock actual: {p.stock})
            </option>
          ))}
        </select>
        {productoSeleccionado && (
          <p className="text-xs text-gray-400 mb-3">
            Ubicación: {productoSeleccionado.ubicacion}
          </p>
        )}

        <label className="text-sm text-gray-500 block mb-1 mt-3">Cantidad comprada</label>
        <input
          type="number"
          value={cantidad}
          onChange={(e) => setCantidad(Number(e.target.value))}
          min={1}
          className="w-full border border-gray-300 rounded-lg px-3 py-2 text-sm mb-4"
        />

        {errorFormulario && (
          <div className="bg-red-100 text-red-700 text-sm rounded-lg p-3 mb-4">
            {errorFormulario}
          </div>
        )}

        <Button type="submit" disabled={registrando}>
          {registrando ? "Registrando..." : "Registrar compra"}
        </Button>
      </form>

      <h2 className="font-semibold mb-3">Historial de compras</h2>

      {cargando && <p className="text-gray-500">Cargando...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}
      {!cargando && !error && compras.length === 0 && (
        <p className="text-gray-500">Todavía no se ha registrado ninguna compra.</p>
      )}

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-500 text-left">
            <tr>
              <th className="px-4 py-3">Fecha</th>
              <th className="px-4 py-3">Producto</th>
              <th className="px-4 py-3">Cantidad</th>
              <th className="px-4 py-3">Costo unitario</th>
              <th className="px-4 py-3">Costo total</th>
              <th className="px-4 py-3">Responsable</th>
            </tr>
          </thead>
          <tbody>
            {compras.map((c) => (
              <tr key={c.id_compra} className="border-t border-gray-100">
                <td className="px-4 py-3 text-gray-500">
                  {new Date(c.fecha).toLocaleString("es-CO", { day: "2-digit", month: "2-digit", year: "numeric", hour: "2-digit", minute: "2-digit" })}
                </td>
                <td className="px-4 py-3 font-medium">{c.nombre_producto}</td>
                <td className="px-4 py-3">+{c.cantidad}</td>
                <td className="px-4 py-3">{formatearCOP(c.costo_unitario_momento)}</td>
                <td className="px-4 py-3 font-medium">{formatearCOP(c.costo_total)}</td>
                <td className="px-4 py-3 text-gray-500">{c.nombre_responsable}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </div>
  );
}