package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
)

type ReporteHandler struct {
	ReporteRepo *repository.ReporteRepository
}

// GET /api/reportes/pedidos?estado=&fecha_desde=&fecha_hasta= - RF-28
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

// GET /api/reportes/tiempos - RF-29
func (h *ReporteHandler) TiemposPorEtapa(w http.ResponseWriter, r *http.Request) {
	tiempos, err := h.ReporteRepo.TiemposPorEtapa(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular los tiempos por etapa"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tiempos)
}

// GET /api/reportes/incidencias - indicador operativo (equivalente IIT generado por el sistema)
func (h *ReporteHandler) IndicadorIncidencias(w http.ResponseWriter, r *http.Request) {
	indicadores, err := h.ReporteRepo.IndicadorIncidenciasEscaneo(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo calcular el indicador de incidencias"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(indicadores)
}
