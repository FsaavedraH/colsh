package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

type DespachoHandler struct {
	PedidoRepo   *repository.PedidoRepository
	DespachoRepo *repository.DespachoRepository
	Ledger       *ledger.LedgerAdapter
}

func (h *DespachoHandler) registrarEnLedgerSiDisponible(idPedido, estado, responsable string) {
	if h.Ledger == nil {
		return
	}
	idEvento := uuid.New().String()
	fecha := time.Now().Format(time.RFC3339)
	_ = h.Ledger.RegistrarEnLedger(context.Background(), idEvento, idPedido, estado, fecha, responsable)
}

// GET /api/despacho - RF-20
func (h *DespachoHandler) ListarOrdenes(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarParaDespacho(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la lista de ordenes"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

// GET /api/despacho/historial
func (h *DespachoHandler) ListarHistorial(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarHistorialDespacho(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el historial"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

type GenerarDespachoRequest struct {
	IDPedido      string `json:"id_pedido"`
	Transportista string `json:"transportista"`
}

// POST /api/despacho - RF-19, RF-20, RF-21, RF-24
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

	codigoSeguimiento := strings.ToUpper(strings.Split(pedido.IDPedido.String(), "-")[0])

	if err := h.DespachoRepo.RegistrarEvento(r.Context(), idPedido, "En despacho", transportista); err != nil {
		http.Error(w, `{"error":"No se pudo registrar el despacho: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En despacho"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	h.registrarEnLedgerSiDisponible(req.IDPedido, "En despacho", req.Transportista)

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

// POST /api/entrega - RF-22, RF-23, RF-24
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

	if err := h.DespachoRepo.RegistrarEvento(r.Context(), idPedido, "Entregado", transportista); err != nil {
		http.Error(w, `{"error":"No se pudo registrar la entrega: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "Entregado"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	h.registrarEnLedgerSiDisponible(req.IDPedido, "Entregado", req.Transportista)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "entregado"})
}