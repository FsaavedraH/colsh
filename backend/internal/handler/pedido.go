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
	Repo           *repository.PedidoRepository
	InventarioRepo *repository.InventarioRepository
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

// POST /api/pedidos - RF-01, RF-02, RF-03. Reserva (descuenta) el stock de inmediato
// para evitar sobreventa entre pedidos concurrentes (RF-05). Si no alcanza, el pedido
// queda "En espera por inventario" sin haber tocado el stock.
func (h *PedidoHandler) CrearPedido(w http.ResponseWriter, r *http.Request) {
	var req CrearPedidoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

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
		Estado:           "Pendiente",
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

	estadoFinal := pedido.Estado
	if h.InventarioRepo != nil {
		disponible, _, errReservar := h.InventarioRepo.ReservarStock(r.Context(), productos)
		if errReservar == nil && !disponible {
			estadoFinal = "En espera por inventario"
			pedido.Estado = estadoFinal
		}
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
		"estado":    estadoFinal,
	})
}

type PedidoConDetalle struct {
	IDPedido         uuid.UUID                      `json:"id_pedido"`
	FechaCreacion    time.Time                      `json:"fecha_creacion"`
	Estado           string                         `json:"estado"`
	IDCliente        uuid.UUID                      `json:"id_cliente"`
	DireccionEntrega string                         `json:"direccion_entrega"`
	Productos        []repository.ItemDetallePedido `json:"productos"`
}

// GET /api/pedidos/{id} - RF-04, incluye el detalle de productos del pedido
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

	items, err := h.Repo.ObtenerDetalleConNombres(r.Context(), id)
	if err != nil {
		items = []repository.ItemDetallePedido{}
	}

	respuesta := PedidoConDetalle{
		IDPedido:         pedido.IDPedido,
		FechaCreacion:    pedido.FechaCreacion,
		Estado:           pedido.Estado,
		IDCliente:        pedido.IDCliente,
		DireccionEntrega: pedido.DireccionEntrega,
		Productos:        items,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(respuesta)
}

// GET /api/mis-pedidos?cliente_id=... - Mis pedidos (Cliente)
func (h *PedidoHandler) ListarMisPedidos(w http.ResponseWriter, r *http.Request) {
	clienteIDParam := r.URL.Query().Get("cliente_id")
	clienteID, err := uuid.Parse(clienteIDParam)
	if err != nil {
		http.Error(w, `{"error":"cliente_id invalido o faltante"}`, http.StatusBadRequest)
		return
	}

	pedidos, err := h.Repo.ListarPorCliente(r.Context(), clienteID)
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener los pedidos"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pedidos)
}

type CancelarPedidoRequest struct {
	ClienteID string `json:"cliente_id"`
}

var estadosCancelables = map[string]bool{
	"Pendiente":                 true,
	"En espera por inventario":  true,
	"En recoleccion":            true,
	"En empaque":                true,
}

// POST /api/pedidos/{id}/cancelar - Cliente (dueno del pedido) o Administrador.
// Libera el stock reservado (si lo habia) y marca el pedido como "Cancelado".
// No se puede cancelar una vez el pedido esta "En despacho" o mas adelante.
func (h *PedidoHandler) CancelarPedido(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	id, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, `{"error":"id invalido"}`, http.StatusBadRequest)
		return
	}

	var req CancelarPedidoRequest
	json.NewDecoder(r.Body).Decode(&req) // el body es opcional (Admin puede omitirlo)

	pedido, err := h.Repo.ConsultarPorID(r.Context(), id)
	if err != nil {
		if err == pgx.ErrNoRows {
			http.Error(w, `{"error":"Pedido no encontrado"}`, http.StatusNotFound)
			return
		}
		http.Error(w, `{"error":"Error al consultar pedido"}`, http.StatusInternalServerError)
		return
	}

	rol := r.Header.Get("X-User-Role")
	if rol == "Cliente" {
		if req.ClienteID == "" || req.ClienteID != pedido.IDCliente.String() {
			http.Error(w, `{"error":"No tienes permiso para cancelar este pedido"}`, http.StatusForbidden)
			return
		}
	}

	if !estadosCancelables[pedido.Estado] {
		http.Error(w, `{"error":"El pedido ya no se puede cancelar en su estado actual: `+pedido.Estado+`"}`, http.StatusConflict)
		return
	}

	// Si el pedido nunca llego a reservar stock (estaba en espera por inventario),
	// no hay nada que liberar. En cualquier otro estado cancelable, si habia stock reservado.
	if pedido.Estado != "En espera por inventario" && h.InventarioRepo != nil {
		productos, errProd := h.Repo.ObtenerProductosDelPedido(r.Context(), id)
		if errProd == nil {
			h.InventarioRepo.LiberarStock(r.Context(), productos)
		}
	}

	if err := h.Repo.ActualizarEstado(r.Context(), id, "Cancelado"); err != nil {
		http.Error(w, `{"error":"No se pudo cancelar el pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "Cancelado"})
}