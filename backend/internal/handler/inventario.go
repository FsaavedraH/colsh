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

type RegistrarCompraRequest struct {
	IDProducto  string `json:"id_producto"`
	Cantidad    int    `json:"cantidad"`
	Responsable string `json:"responsable"`
}

// POST /api/inventario/compras - RF-07 (version simplificada, solo Admin).
// Registra un ingreso de compra: suma stock al producto seleccionado (nunca crea
// productos nuevos) y, si con ese stock ya alcanza para pedidos que estaban
// "En espera por inventario", los reactiva automaticamente a "Pendiente".
func (h *InventarioHandler) RegistrarCompra(w http.ResponseWriter, r *http.Request) {
	var req RegistrarCompraRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idProducto, err := uuid.Parse(req.IDProducto)
	if err != nil {
		http.Error(w, `{"error":"id_producto invalido"}`, http.StatusBadRequest)
		return
	}

	responsable, err := uuid.Parse(req.Responsable)
	if err != nil {
		http.Error(w, `{"error":"responsable invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Cantidad <= 0 {
		http.Error(w, `{"error":"La cantidad debe ser mayor a 0"}`, http.StatusBadRequest)
		return
	}

	idCompra, err := h.InventarioRepo.RegistrarCompra(r.Context(), idProducto, req.Cantidad, responsable)
	if err != nil {
		http.Error(w, `{"error":"No se pudo registrar la compra: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	// Reactivar pedidos que estaban "En espera por inventario" y que ahora
	// ya tengan stock suficiente para todos sus productos.
	pedidosEnEspera, err := h.PedidoRepo.ListarPorEstado(r.Context(), "En espera por inventario")
	if err == nil {
		for _, pedido := range pedidosEnEspera {
			productosDelPedido, errProd := h.PedidoRepo.ObtenerProductosDelPedido(r.Context(), pedido.IDPedido)
			if errProd != nil {
				continue
			}
			disponible, _, errValidar := h.InventarioRepo.ValidarStock(r.Context(), productosDelPedido)
			if errValidar == nil && disponible {
				h.PedidoRepo.ActualizarEstado(r.Context(), pedido.IDPedido, "Pendiente")
			}
		}
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id_compra": idCompra.String(), "estado": "registrado"})
}

// GET /api/inventario/compras - historial de compras (Admin)
func (h *InventarioHandler) ListarCompras(w http.ResponseWriter, r *http.Request) {
	compras, err := h.InventarioRepo.ListarCompras(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el historial de compras"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(compras)
}