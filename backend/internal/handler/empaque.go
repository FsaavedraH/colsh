package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

type EmpaqueHandler struct {
	PedidoRepo  *repository.PedidoRepository
	EmpaqueRepo *repository.EmpaqueRepository
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

// POST /api/empaque/escanear - RF-17, RF-26
func (h *EmpaqueHandler) EscanearValidacion(w http.ResponseWriter, r *http.Request) {
	var req EscanearValidacionEmpaqueRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if req.IDProductoEscaneado != req.IDProductoEsperado {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Producto incorrecto. El producto escaneado no corresponde. Intenta nuevamente.",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"coincide": true,
		"mensaje":  "Producto validado correctamente para empaque",
	})
}

type ConfirmarEmpaqueRequest struct {
	IDPedido    string `json:"id_pedido"`
	Responsable string `json:"responsable"`
}

// POST /api/empaque - RF-18
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "empacado"})
}
