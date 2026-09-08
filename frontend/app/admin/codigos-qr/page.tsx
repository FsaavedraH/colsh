"use client";

import { useEffect, useMemo, useState } from "react";
import { QRCodeSVG } from "qrcode.react";
import { apiFetch } from "@/lib/api";
import Button from "@/components/ui/Button";

interface ProductoConToken {
  nombre: string;
  stock: number;
  ubicacion: string;
  token_producto: string;
  token_ubicacion: string;
}

export default function CodigosQRPage() {
  const [productos, setProductos] = useState<ProductoConToken[]>([]);
  const [seleccionados, setSeleccionados] = useState<Set<string>>(new Set());
  const [busqueda, setBusqueda] = useState("");
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");
  const [mostrarImpresion, setMostrarImpresion] = useState(false);

  useEffect(() => {
    apiFetch<ProductoConToken[]>("/api/inventario/qr", { rol: "Administrador" })
      .then((data) => setProductos(data || []))
      .catch((err) => setError(err.message))
      .finally(() => setCargando(false));
  }, []);

  const productosFiltrados = useMemo(() => {
    const q = busqueda.trim().toLowerCase();
    if (!q) return productos;
    return productos.filter((p) => p.nombre.toLowerCase().includes(q));
  }, [productos, busqueda]);

  function alternarSeleccion(nombre: string) {
    setSeleccionados((prev) => {
      const nuevo = new Set(prev);
      if (nuevo.has(nombre)) nuevo.delete(nombre);
      else nuevo.add(nombre);
      return nuevo;
    });
  }

  function seleccionarTodos() {
    setSeleccionados(new Set(productosFiltrados.map((p) => p.nombre)));
  }

  function limpiarSeleccion() {
    setSeleccionados(new Set());
  }

  const productosParaImprimir = productos.filter((p) => seleccionados.has(p.nombre));

  if (mostrarImpresion) {
    return (
      <div>
        <div className="mb-6 print:hidden">
          <Button onClick={() => window.print()}>Imprimir / Guardar como PDF</Button>
          <button
            onClick={() => setMostrarImpresion(false)}
            className="ml-3 text-sm text-gray-500 hover:underline"
          >
            ← Volver a la selección
          </button>
        </div>

        <div className="grid grid-cols-2 md:grid-cols-3 gap-6 print:grid-cols-3">
          {productosParaImprimir.map((p) => (
            <div key={p.nombre} className="border border-gray-300 rounded-lg p-4 break-inside-avoid">
              <p className="text-xs text-gray-400 mb-1 text-center">Producto</p>
              <div className="flex justify-center mb-2">
                <QRCodeSVG value={p.token_producto} size={110} />
              </div>
              <p className="text-sm font-semibold text-center mb-3">{p.nombre}</p>

              <p className="text-xs text-gray-400 mb-1 text-center">Ubicación</p>
              <div className="flex justify-center mb-2">
                <QRCodeSVG value={p.token_ubicacion} size={90} />
              </div>
              <p className="text-xs font-mono text-center">{p.ubicacion}</p>
            </div>
          ))}
        </div>
      </div>
    );
  }

  return (
    <div>
      <h1 className="text-2xl font-bold mb-1">Códigos QR</h1>
      <p className="text-gray-500 mb-6">
        Selecciona los productos que necesitas etiquetar (nuevos ingresos, reposición de etiquetas dañadas, etc.). Los códigos generados están firmados digitalmente: no pueden falsificarse con otro generador de QR.
      </p>

      <div className="flex items-center gap-3 mb-4">
        <input
          type="text"
          placeholder="Buscar producto..."
          value={busqueda}
          onChange={(e) => setBusqueda(e.target.value)}
          className="border border-gray-300 rounded-lg px-3 py-2 text-sm flex-1 max-w-sm"
        />
        <button onClick={seleccionarTodos} className="text-sm text-blue-600 hover:underline">
          Seleccionar todos
        </button>
        <button onClick={limpiarSeleccion} className="text-sm text-gray-500 hover:underline">
          Limpiar selección
        </button>
      </div>

      {cargando && <p className="text-gray-500">Cargando productos...</p>}
      {error && <p className="text-red-600">Error: {error}</p>}

      <div className="bg-white rounded-xl border border-gray-200 overflow-hidden mb-6">
        <table className="w-full text-sm">
          <thead className="bg-gray-50 text-gray-500 text-left">
            <tr>
              <th className="px-4 py-3 w-10"></th>
              <th className="px-4 py-3">Producto</th>
              <th className="px-4 py-3">Ubicación</th>
              <th className="px-4 py-3">Stock</th>
            </tr>
          </thead>
          <tbody>
            {productosFiltrados.map((p) => (
              <tr
                key={p.nombre}
                className="border-t border-gray-100 hover:bg-gray-50 cursor-pointer"
                onClick={() => alternarSeleccion(p.nombre)}
              >
                <td className="px-4 py-3">
                  <input
                    type="checkbox"
                    checked={seleccionados.has(p.nombre)}
                    onChange={() => alternarSeleccion(p.nombre)}
                    onClick={(e) => e.stopPropagation()}
                  />
                </td>
                <td className="px-4 py-3 font-medium">{p.nombre}</td>
                <td className="px-4 py-3 font-mono text-gray-500">{p.ubicacion}</td>
                <td className="px-4 py-3 text-gray-500">{p.stock}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      <Button onClick={() => setMostrarImpresion(true)} disabled={seleccionados.size === 0}>
        Generar códigos ({seleccionados.size} seleccionados)
      </Button>
    </div>
  );
}