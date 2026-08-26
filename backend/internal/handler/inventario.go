package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

type InventarioHandler struct {
	InventarioRepo *repository.InventarioRepository
	PedidoRepo     *repository.PedidoRepository
}

type ValidarInventarioRequest struct {
	IDPedido string `json:"id_pedido"`
}

func (h *InventarioHandler) ValidarInventario(w http.ResponseWriter, r *http.Request) {
	var req ValidarInventarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	productos, err := h.PedidoRepo.ObtenerProductosDelPedido(r.Context(), idPedido)
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el detalle del pedido"}`, http.StatusInternalServerError)
		return
	}

	disponible, faltantes, err := h.InventarioRepo.ValidarStock(r.Context(), productos)
	if err != nil {
		http.Error(w, `{"error":"Error al validar inventario"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if !disponible {
		if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En espera por inventario"); err != nil {
			http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"disponibilidad": false,
			"faltantes":      faltantes,
			"estado_pedido":  "En espera por inventario",
		})
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En recoleccion"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(map[string]interface{}{
		"disponibilidad": true,
		"estado_pedido":  "En recoleccion",
	})
}

type ActualizarInventarioRequest struct {
	IDProducto string `json:"id_producto"`
	Cantidad   int    `json:"cantidad"`
}

func (h *InventarioHandler) ActualizarInventario(w http.ResponseWriter, r *http.Request) {
	var req ActualizarInventarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idProducto, err := uuid.Parse(req.IDProducto)
	if err != nil {
		http.Error(w, `{"error":"id_producto invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Cantidad <= 0 {
		http.Error(w, `{"error":"La cantidad debe ser mayor a 0"}`, http.StatusBadRequest)
		return
	}

	err = h.InventarioRepo.ActualizarStock(r.Context(), idProducto, req.Cantidad)
	if err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el inventario: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "actualizado"})
}

// GET /api/productos - catalogo para clientes
func (h *InventarioHandler) ListarCatalogo(w http.ResponseWriter, r *http.Request) {
	productos, err := h.InventarioRepo.ListarCatalogo(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el catalogo"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(productos)
}