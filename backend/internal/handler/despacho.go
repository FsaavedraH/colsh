package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

type DespachoHandler struct {
	PedidoRepo   *repository.PedidoRepository
	DespachoRepo *repository.DespachoRepository
}

type GenerarDespachoRequest struct {
	IDPedido      string `json:"id_pedido"`
	Transportista string `json:"transportista"`
}

// POST /api/despacho - RF-19, RF-20, RF-21
func (h *DespachoHandler) GenerarDespacho(w http.ResponseWriter, r *http.Request) {
	var req GenerarDespachoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	transportista, err := uuid.Parse(req.Transportista)
	if err != nil {
		http.Error(w, `{"error":"transportista invalido"}`, http.StatusBadRequest)
		return
	}

	pedido, err := h.PedidoRepo.ConsultarPorID(r.Context(), idPedido)
	if err != nil {
		http.Error(w, `{"error":"Pedido no encontrado"}`, http.StatusNotFound)
		return
	}

	// RF-19: generar codigo de seguimiento (simple, basado en el id del pedido)
	codigoSeguimiento := strings.ToUpper(strings.Split(pedido.IDPedido.String(), "-")[0])

	// RF-20: registrar evento de despacho con el transportador asignado
	if err := h.DespachoRepo.RegistrarEvento(r.Context(), idPedido, "En despacho", transportista); err != nil {
		http.Error(w, `{"error":"No se pudo registrar el despacho: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// RF-21: actualizar estado del pedido
	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En despacho"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"estado":             "en_despacho",
		"id_pedido":          pedido.IDPedido,
		"codigo_seguimiento": codigoSeguimiento,
		"direccion_entrega":  pedido.DireccionEntrega,
	})
}

type ConfirmarEntregaRequest struct {
	IDPedido      string `json:"id_pedido"`
	Transportista string `json:"transportista"`
}

// POST /api/entrega - RF-22, RF-23
func (h *DespachoHandler) ConfirmarEntrega(w http.ResponseWriter, r *http.Request) {
	var req ConfirmarEntregaRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	transportista, err := uuid.Parse(req.Transportista)
	if err != nil {
		http.Error(w, `{"error":"transportista invalido"}`, http.StatusBadRequest)
		return
	}

	// RF-22: registrar fecha y hora de entrega final
	if err := h.DespachoRepo.RegistrarEvento(r.Context(), idPedido, "Entregado", transportista); err != nil {
		http.Error(w, `{"error":"No se pudo registrar la entrega: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// RF-23: cambiar estado a "Entregado"
	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "Entregado"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "entregado"})
}
