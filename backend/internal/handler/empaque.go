package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

type EmpaqueHandler struct {
	PedidoRepo  *repository.PedidoRepository
	EmpaqueRepo *repository.EmpaqueRepository
	ReporteRepo *repository.ReporteRepository
	Ledger      *ledger.LedgerAdapter
}

func (h *EmpaqueHandler) registrarEnLedgerSiDisponible(idPedido, estado, responsable string) {
	if h.Ledger == nil {
		return
	}
	idEvento := uuid.New().String()
	fecha := time.Now().Format(time.RFC3339)
	_ = h.Ledger.RegistrarEnLedger(context.Background(), idEvento, idPedido, estado, fecha, responsable)
}

type RecepcionEmpaqueRequest struct {
	IDPedido string `json:"id_pedido"`
}

// POST /api/empaque/recepcion - RF-16
func (h *EmpaqueHandler) RecepcionEmpaque(w http.ResponseWriter, r *http.Request) {
	var req RecepcionEmpaqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	pedido, err := h.PedidoRepo.ConsultarPorID(r.Context(), idPedido)
	if err != nil {
		http.Error(w, `{"error":"Pedido no encontrado"}`, http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id_pedido": pedido.IDPedido,
		"estado":    pedido.Estado,
		"mensaje":   "Pedido recibido en empaque",
	})
}

type EscanearValidacionEmpaqueRequest struct {
	IDPedido            string `json:"id_pedido"`
	IDProductoEsperado  string `json:"id_producto_esperado"`
	IDProductoEscaneado string `json:"id_producto_escaneado"`
}

// POST /api/empaque/escanear - RF-17, RF-26, RF-28
func (h *EmpaqueHandler) EscanearValidacion(w http.ResponseWriter, r *http.Request) {
	var req EscanearValidacionEmpaqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if req.IDProductoEscaneado != req.IDProductoEsperado {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "incorrecto", "empaque")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Producto incorrecto. El producto escaneado no corresponde. Intenta nuevamente.",
		})
		return
	}

	h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "correcto", "empaque")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"coincide": true,
		"mensaje":  "Producto validado correctamente para empaque",
	})
}

type ConfirmarEmpaqueRequest struct {
	IDPedido    string `json:"id_pedido"`
	Responsable string `json:"responsable"`
}

// POST /api/empaque - RF-18, RF-24
func (h *EmpaqueHandler) ConfirmarEmpaque(w http.ResponseWriter, r *http.Request) {
	var req ConfirmarEmpaqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	responsable, err := uuid.Parse(req.Responsable)
	if err != nil {
		http.Error(w, `{"error":"responsable invalido"}`, http.StatusBadRequest)
		return
	}

	if err := h.EmpaqueRepo.RegistrarEmpaque(r.Context(), idPedido, responsable); err != nil {
		http.Error(w, `{"error":"No se pudo registrar el empaque: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En despacho"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	h.registrarEnLedgerSiDisponible(req.IDPedido, "En empaque", req.Responsable)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "empacado"})
}
