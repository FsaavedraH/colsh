package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/domain"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type PedidoHandler struct {
	Repo *repository.PedidoRepository
}

type CrearPedidoRequest struct {
	ClienteID        string           `json:"cliente_id"`
	DireccionEntrega string           `json:"direccion_entrega"`
	Productos        []ProductoPedido `json:"productos"`
}

type ProductoPedido struct {
	IDProducto string `json:"id_producto"`
	Cantidad   int    `json:"cantidad"`
}

// POST /api/pedidos - RF-01
func (h *PedidoHandler) CrearPedido(w http.ResponseWriter, r *http.Request) {
	var req CrearPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	// RF-02: Validar datos obligatorios completos
	if req.ClienteID == "" || req.DireccionEntrega == "" || len(req.Productos) == 0 {
		http.Error(w, `{"error":"Faltan datos obligatorios: cliente_id, direccion_entrega o productos"}`, http.StatusBadRequest)
		return
	}

	clienteID, err := uuid.Parse(req.ClienteID)
	if err != nil {
		http.Error(w, `{"error":"cliente_id invalido"}`, http.StatusBadRequest)
		return
	}

	pedido := domain.Pedido{
		IDPedido:         uuid.New(),
		FechaCreacion:    time.Now(),
		Estado:           "Pendiente", // RF-03
		IDCliente:        clienteID,
		DireccionEntrega: req.DireccionEntrega,
	}

	var productos []repository.ProductoPedidoInput
	for _, p := range req.Productos {
		productos = append(productos, repository.ProductoPedidoInput{
			IDProducto: p.IDProducto,
			Cantidad:   p.Cantidad,
		})
	}

	err = h.Repo.Crear(r.Context(), &pedido, productos)
	if err != nil {
		http.Error(w, `{"error":"No se pudo crear el pedido: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"id_pedido": pedido.IDPedido.String(),
		"estado":    pedido.Estado,
	})
}

// GET /api/pedidos/{id} - RF-04
func (h *PedidoHandler) ConsultarPedido(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, `{"error":"id invalido"}`, http.StatusBadRequest)
		return
	}

	pedido, err := h.Repo.ConsultarPorID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, `{"error":"Pedido no encontrado"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Error al consultar pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedido)
}
