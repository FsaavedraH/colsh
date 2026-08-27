package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
)

type ReporteHandler struct {
	ReporteRepo *repository.ReporteRepository
}

func (h *ReporteHandler) ListarPedidos(w http.ResponseWriter, r *http.Request) {
	estado := r.URL.Query().Get("estado")
	fechaDesde := r.URL.Query().Get("fecha_desde")
	fechaHasta := r.URL.Query().Get("fecha_hasta")

	pedidos, err := h.ReporteRepo.ListarPedidosFiltrado(r.Context(), estado, fechaDesde, fechaHasta)
	if err != nil {
		http.Error(w, `{"error":"No se pudo generar el reporte"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

func (h *ReporteHandler) TiemposPorEtapa(w http.ResponseWriter, r *http.Request) {
	tiempos, err := h.ReporteRepo.TiemposPorEtapa(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular los tiempos por etapa"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tiempos)
}

func (h *ReporteHandler) IndicadorIncidencias(w http.ResponseWriter, r *http.Request) {
	indicadores, err := h.ReporteRepo.IndicadorIncidenciasEscaneo(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular el indicador de incidencias"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indicadores)
}

func (h *ReporteHandler) ConteoPorEstado(w http.ResponseWriter, r *http.Request) {
	conteo, err := h.ReporteRepo.ConteoPedidosPorEstado(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular el conteo por estado"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(conteo)
}

func (h *ReporteHandler) ProductosMasVendidos(w http.ResponseWriter, r *http.Request) {
	productos, err := h.ReporteRepo.ProductosMasVendidos(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular el ranking de productos"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}

// GET /api/reportes/pedidos-por-dia
func (h *ReporteHandler) PedidosPorDia(w http.ResponseWriter, r *http.Request) {
	datos, err := h.ReporteRepo.PedidosPorDia(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular la tendencia de pedidos"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(datos)
}