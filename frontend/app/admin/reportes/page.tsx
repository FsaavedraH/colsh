"use client";

import { useEffect, useState } from "react";
import { apiFetch } from "@/lib/api";
import Badge from "@/components/ui/Badge";

interface Pedido {
  id_pedido: string;
  fecha_creacion: string;
  estado: string;
  nombre_cliente: string;
}

interface TiempoEtapa {
  etapa: string;
  tiempo_promedio_minutos: number;
  cantidad_pedidos: number;
}

interface Incidencia {
  etapa: string;
  total_intentos: number;
  intentos_incorrectos: number;
  tasa_incidencia_porcentaje: number;
}

interface ConteoEstado {
  estado: string;
  cantidad: number;
}

interface ProductoTop {
  nombre_producto: string;
  cantidad_vendida: number;
  total_pedidos: number;
}

interface PedidoPorDia {
  fecha: string;
  cantidad: number;
}

interface EventoTrazabilidad {
  id_evento: string;
  id_pedido: string;
  estado: string;
  fecha: string;
  responsable: string;
}

const colorTarjeta: Record<string, string> = {
  "Pendiente": "bg-gray-100 text-gray-700",
  "En espera por inventario": "bg-red-100 text-red-700",
  "En recoleccion": "bg-orange-100 text-orange-700",
  "En empaque": "bg-teal-100 text-teal-700",
  "En despacho": "bg-amber-100 text-amber-700",
  "Entregado": "bg-green-100 text-green-700",
};

// Grafica de linea simple con SVG, sin dependencias externas.
function GraficaTendencia({ datos }: { datos: PedidoPorDia[] }) {
  if (datos.length === 0) {
    return <p className="text-gray-400 text-sm py-8 text-center">Sin pedidos en los últimos 14 días.</p>;
  }

  const ancho = 600;
  const alto = 160;
  const padding = 30;
  const maxCantidad = Math.max(...datos.map((d) => d.cantidad), 1);

  const puntos = datos.map((d, i) => {
    const x = padding + (i / Math.max(datos.length - 1, 1)) * (ancho - padding * 2);
    const y = alto - padding - (d.cantidad / maxCantidad) * (alto - padding * 2);
    return { x, y, ...d };
  });

  const lineaPath = puntos.map((p, i) => `${i === 0 ? "M" : "L"} ${p.x} ${p.y}`).join(" ");
  const areaPath = `${lineaPath} L ${puntos[puntos.length - 1].x} ${alto - padding} L ${puntos[0].x} ${alto - padding} Z`;

  return (
    <svg viewBox={`0 0 ${ancho} ${alto}`} className="w-full h-40">
      <line x1={padding} y1={alto - padding} x2={ancho - padding} y2={alto - padding} stroke="#e5e7eb" strokeWidth="1" />
      <path d={areaPath} fill="#3b82f6" fillOpacity="0.08" />
      <path d={lineaPath} fill="none" stroke="#3b82f6" strokeWidth="2" />
      {puntos.map((p, i) => (
        <g key={i}>
          <circle cx={p.x} cy={p.y} r="3.5" fill="#3b82f6" />
          <text x={p.x} y={alto - padding + 16} textAnchor="middle" fontSize="9" fill="#9ca3af">
            {new Date(p.fecha).toLocaleDateString("es-CO", { day: "2-digit", month: "2-digit" })}
          </text>
          {p.cantidad > 0 && (
            <text x={p.x} y={p.y - 8} textAnchor="middle" fontSize="10" fill="#3b82f6" fontWeight="bold">
              {p.cantidad}
            </text>
          )}
        </g>
      ))}
    </svg>
  );
}

export default function ReportesPage() {
  const [pedidos, setPedidos] = useState<Pedido[]>([]);
  const [tiempos, setTiempos] = useState<TiempoEtapa[]>([]);
  const [incidencias, setIncidencias] = useState<Incidencia[]>([]);
  const [conteos, setConteos] = useState<ConteoEstado[]>([]);
  const [productosTop, setProductosTop] = useState<ProductoTop[]>([]);
  const [pedidosPorDia, setPedidosPorDia] = useState<PedidoPorDia[]>([]);
  const [filtroEstado, setFiltroEstado] = useState("");
  const [cargando, setCargando] = useState(true);
  const [error, setError] = useState("");

  const [pedidoSeleccionado, setPedidoSeleccionado] = useState<string | null>(null);
  const [eventos, setEventos] = useState<EventoTrazabilidad[]>([]);
  const [cargandoEventos, setCargandoEventos] = useState(false);
  const [errorEventos, setErrorEventos] = useState("");

  useEffect(() => {
    cargarDatos();
  }, [filtroEstado]);

  async function cargarDatos() {
    setCargando(true);
    setError("");
    try {
      const query = filtroEstado ? `?estado=${encodeURIComponent(filtroEstado)}` : "";
      const [pedidosData, tiemposData, incidenciasData, conteosData, productosTopData, pedidosPorDiaData] =
        await Promise.all([
          apiFetch<Pedido[]>(`/api/reportes/pedidos${query}`, { rol: "Administrador" }),
          apiFetch<TiempoEtapa[]>("/api/reportes/tiempos", { rol: "Administrador" }),
          apiFetch<Incidencia[]>("/api/reportes/incidencias", { rol: "Administrador" }),
          apiFetch<ConteoEstado[]>("/api/reportes/conteo-estados", { rol: "Administrador" }),
          apiFetch<ProductoTop[]>("/api/reportes/productos-top", { rol: "Administrador" }),
          apiFetch<PedidoPorDia[]>("/api/reportes/pedidos-por-dia", { rol: "Administrador" }),
        ]);
      setPedidos(pedidosData || []);
      setTiempos(tiemposData || []);
      setIncidencias(incidenciasData || []);
      setConteos(conteosData || []);
      setProductosTop(productosTopData || []);
      setPedidosPorDia(pedidosPorDiaData || []);
    } catch (err: any) {
      setError(err.message);
    } finally {
      setCargando(false);
    }
  }

  async function verTrazabilidad(idPedido: string) {
    setPedidoSeleccionado(idPedido);
    setCargandoEventos(true);
    setErrorEventos("");
    setEventos([]);
    try {
      const data = await apiFetch<EventoTrazabilidad[]>(`/api/trazabilidad/${idPedido}`, {
        rol: "Administrador",
      });
      const ordenados = (data || []).sort(
        (a, b) => new Date(a.fecha).getTime() - new Date(b.fecha).getTime()
      );
      setEventos(ordenados);
    } catch (err: any) {
      setErrorEventos(err.message);
    } finally {
      setCargandoEventos(false);
    }
  }

  const estados = ["", "Pendiente", "En espera por inventario", "En recoleccion", "En empaque", "En despacho", "Entregado"];
  const maxTiempo = Math.max(...tiempos.map((t) => t.tiempo_promedio_minutos), 1);
  const maxCantidadVendida = Math.max(...productosTop.map((p) => p.cantidad_vendida), 1);
  const totalPedidos = conteos.reduce((sum, c) => sum + c.cantidad, 0);

  return (
    <div>
      <h1 className="text-2xl font-bold mb-6">Reportes y Trazabilidad</h1>

      {error && <p className="text-red-600 mb-4">Error: {error}</p>}

      {/* SECCION 1: Tendencia principal, arriba y destacada */}
      <div className="bg-white rounded-xl border border-gray-200 p-5 mb-6">
        <div className="flex justify-between items-baseline mb-4">
          <h2 className="font-semibold">Tendencia de pedidos (últimos 14 días)</h2>
          <span className="text-sm text-gray-400">{totalPedidos} pedidos en total</span>
        </div>
        <GraficaTendencia datos={pedidosPorDia} />
      </div>

      {/* SECCION 2: Tarjetas de conteo por estado */}
      <div className="grid grid-cols-2 md:grid-cols-6 gap-3 mb-6">
        {estados.slice(1).map((estado) => {
          const conteo = conteos.find((c) => c.estado === estado);
          return (
            <div key={estado} className={`rounded-xl p-4 ${colorTarjeta[estado] || "bg-gray-100"}`}>
              <div className="text-2xl font-bold">{conteo?.cantidad || 0}</div>
              <div className="text-xs mt-1">{estado}</div>
            </div>
          );
        })}
      </div>

      {/* SECCION 3: Indicadores de apoyo, en grid de 3 columnas */}
      <div className="grid md:grid-cols-3 gap-4 mb-6">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h2 className="font-semibold mb-3 text-sm">Tiempo promedio por etapa</h2>
          {tiempos.length === 0 && !cargando && (
            <p className="text-gray-400 text-xs">Sin datos suficientes.</p>
          )}
          <div className="space-y-3">
            {tiempos.map((t) => (
              <div key={t.etapa}>
                <div className="flex justify-between text-xs text-gray-500 mb-1">
                  <span>{t.etapa}</span>
                  <span>{t.tiempo_promedio_minutos.toFixed(1)} min</span>
                </div>
                <div className="w-full bg-gray-100 rounded-full h-1.5">
                  <div
                    className="bg-blue-500 h-1.5 rounded-full"
                    style={{ width: `${Math.max((t.tiempo_promedio_minutos / maxTiempo) * 100, 4)}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h2 className="font-semibold mb-3 text-sm">Incidencias en escaneo</h2>
          {incidencias.length === 0 && !cargando && (
            <p className="text-gray-400 text-xs">Sin datos suficientes.</p>
          )}
          <div className="space-y-2">
            {incidencias.map((i) => (
              <div key={i.etapa} className="flex justify-between items-center text-xs">
                <span className="text-gray-600 capitalize">{i.etapa}</span>
                <div className="flex items-center gap-2">
                  <span className="text-gray-400">
                    {i.intentos_incorrectos}/{i.total_intentos}
                  </span>
                  <Badge color={i.tasa_incidencia_porcentaje > 20 ? "red" : "green"}>
                    {i.tasa_incidencia_porcentaje.toFixed(1)}%
                  </Badge>
                </div>
              </div>
            ))}
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h2 className="font-semibold mb-3 text-sm">Productos más vendidos</h2>
          {productosTop.length === 0 && !cargando && (
            <p className="text-gray-400 text-xs">Sin datos suficientes.</p>
          )}
          <div className="space-y-3">
            {productosTop.slice(0, 4).map((p) => (
              <div key={p.nombre_producto}>
                <div className="flex justify-between text-xs text-gray-500 mb-1">
                  <span className="truncate">{p.nombre_producto}</span>
                  <span>{p.cantidad_vendida}u</span>
                </div>
                <div className="w-full bg-gray-100 rounded-full h-1.5">
                  <div
                    className="bg-teal-500 h-1.5 rounded-full"
                    style={{ width: `${Math.max((p.cantidad_vendida / maxCantidadVendida) * 100, 4)}%` }}
                  />
                </div>
              </div>
            ))}
          </div>
        </div>
      </div>

      {/* SECCION 4: Tabla de pedidos + Trazabilidad, lado a lado */}
      <div className="grid md:grid-cols-2 gap-4">
        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <div className="flex justify-between items-center mb-4">
            <h2 className="font-semibold">Listado de pedidos</h2>
            <select
              value={filtroEstado}
              onChange={(e) => setFiltroEstado(e.target.value)}
              className="border border-gray-300 rounded-lg px-3 py-1.5 text-sm"
            >
              {estados.map((e) => (
                <option key={e} value={e}>
                  {e === "" ? "Todos los estados" : e}
                </option>
              ))}
            </select>
          </div>

          {cargando && <p className="text-gray-500 text-sm">Cargando...</p>}

          <table className="w-full text-sm">
            <thead className="text-gray-500 text-left border-b border-gray-100">
              <tr>
                <th className="py-2">Pedido</th>
                <th className="py-2">Cliente</th>
                <th className="py-2">Estado</th>
              </tr>
            </thead>
            <tbody>
              {pedidos.map((p) => (
                <tr
                  key={p.id_pedido}
                  onClick={() => verTrazabilidad(p.id_pedido)}
                  className={`border-b border-gray-50 cursor-pointer hover:bg-gray-50 ${
                    pedidoSeleccionado === p.id_pedido ? "bg-blue-50" : ""
                  }`}
                >
                  <td className="py-2 font-mono text-xs">{p.id_pedido.slice(0, 8).toUpperCase()}</td>
                  <td className="py-2">{p.nombre_cliente}</td>
                  <td className="py-2">
                    <Badge color="yellow">{p.estado}</Badge>
                  </td>
                </tr>
              ))}
            </tbody>
          </table>

          {!cargando && pedidos.length === 0 && (
            <p className="text-gray-500 text-sm py-4">No hay pedidos con ese filtro.</p>
          )}
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-5">
          <h2 className="font-semibold mb-4">Trazabilidad (Hyperledger Fabric)</h2>

          {!pedidoSeleccionado && (
            <p className="text-gray-400 text-sm">Selecciona un pedido de la tabla para ver su historial en el ledger.</p>
          )}

          {pedidoSeleccionado && (
            <>
              <p className="text-xs text-gray-500 font-mono mb-3">
                {pedidoSeleccionado.slice(0, 8).toUpperCase()}
              </p>

              {cargandoEventos && <p className="text-gray-500 text-sm">Consultando ledger...</p>}

              {errorEventos && (
                <div className="bg-amber-50 text-amber-700 text-sm rounded-lg p-3">
                  {errorEventos}
                </div>
              )}

              {!cargandoEventos && !errorEventos && eventos.length === 0 && (
                <p className="text-gray-400 text-sm">Sin eventos registrados en el ledger.</p>
              )}

              <div className="space-y-3">
                {eventos.map((ev, i) => (
                  <div key={ev.id_evento} className="flex gap-3">
                    <div className="flex flex-col items-center">
                      <div className="w-2.5 h-2.5 rounded-full bg-blue-600" />
                      {i < eventos.length - 1 && <div className="w-px flex-1 bg-gray-200" />}
                    </div>
                    <div className="pb-3">
                      <div className="text-sm font-semibold">{ev.estado}</div>
                      <div className="text-xs text-gray-400">
                        {new Date(ev.fecha).toLocaleString("es-CO")}
                      </div>
                    </div>
                  </div>
                ))}
              </div>
            </>
          )}
        </div>
      </div>
    </div>
  );
}