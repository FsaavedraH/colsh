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

type PickingHandler struct {
	PedidoRepo     *repository.PedidoRepository
	InventarioRepo *repository.InventarioRepository
	ReporteRepo    *repository.ReporteRepository
	Ledger         *ledger.LedgerAdapter
}

func (h *PickingHandler) registrarEnLedgerSiDisponible(idPedido, estado, responsable string) {
	if h.Ledger == nil {
		return
	}
	idEvento := uuid.New().String()
	fecha := time.Now().Format(time.RFC3339)
	_ = h.Ledger.RegistrarEnLedger(context.Background(), idEvento, idPedido, estado, fecha, responsable)
}

type IniciarPickingRequest struct {
	IDPedido string `json:"id_pedido"`
}

// POST /api/picking/iniciar - RF-09, RF-10. El operario "toma" el pedido de la cola
// y lo pasa de "Pendiente" a "En recoleccion".
func (h *PickingHandler) IniciarPicking(w http.ResponseWriter, r *http.Request) {
	var req IniciarPickingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	err = h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En recoleccion")
	if err != nil {
		http.Error(w, `{"error":"No se pudo iniciar el picking del pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "En recoleccion"})
}

// GET /api/picking - RF-09, RF-10
func (h *PickingHandler) ListarOrdenes(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarParaPicking(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la lista de ordenes"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

// GET /api/picking/historial
func (h *PickingHandler) ListarHistorial(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarHistorialPicking(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el historial"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

type EscanearUbicacionRequest struct {
	IDPedido           string `json:"id_pedido"`
	IDProducto         string `json:"id_producto"`
	UbicacionEscaneada string `json:"ubicacion_escaneada"`
}

// POST /api/picking/escanear-ubicacion - RF-11
func (h *PickingHandler) EscanearUbicacion(w http.ResponseWriter, r *http.Request) {
	var req EscanearUbicacionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	idProducto, err := uuid.Parse(req.IDProducto)
	if err != nil {
		http.Error(w, `{"error":"id_producto invalido"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	ubicacionEsperada, err := h.InventarioRepo.ObtenerUbicacion(r.Context(), idProducto)
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la ubicacion del producto"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if req.UbicacionEscaneada != ubicacionEsperada {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "ubicacion", "incorrecto", "picking")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide":            false,
			"ubicacion_esperada":  ubicacionEsperada,
			"ubicacion_escaneada": req.UbicacionEscaneada,
			"mensaje":             "Ubicacion incorrecta. La ubicacion escaneada no coincide. Intenta nuevamente.",
		})
		return
	}

	h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "ubicacion", "correcto", "picking")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"coincide": true,
		"mensaje":  "Ubicacion correcta",
	})
}

type EscanearProductoRequest struct {
	IDPedido            string `json:"id_pedido"`
	IDProductoEsperado  string `json:"id_producto_esperado"`
	IDProductoEscaneado string `json:"id_producto_escaneado"`
}

// POST /api/picking/escanear-producto - RF-12, RF-13, RF-26
func (h *PickingHandler) EscanearProducto(w http.ResponseWriter, r *http.Request) {
	var req EscanearProductoRequest
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
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "incorrecto", "picking")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Producto incorrecto. El producto escaneado no corresponde. Intenta nuevamente.",
		})
		return
	}

	h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "correcto", "picking")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"coincide": true,
		"mensaje":  "Producto correcto, listo para recolectar",
	})
}

type ConfirmarRecoleccionRequest struct {
	IDPedido    string `json:"id_pedido"`
	IDProducto  string `json:"id_producto"`
	Cantidad    int    `json:"cantidad"`
	Responsable string `json:"responsable"`
}

// POST /api/recoleccion - RF-14, RF-15, RF-24. El stock ya fue reservado al crear
// el pedido (ver PedidoHandler.CrearPedido), asi que aqui NO se vuelve a descontar.
func (h *PickingHandler) ConfirmarRecoleccion(w http.ResponseWriter, r *http.Request) {
	var req ConfirmarRecoleccionRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Cantidad <= 0 {
		http.Error(w, `{"error":"La cantidad debe ser mayor a 0"}`, http.StatusBadRequest)
		return
	}

	idPedido, err := uuid.Parse(req.IDPedido)
	if err != nil {
		http.Error(w, `{"error":"id_pedido invalido"}`, http.StatusBadRequest)
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En empaque"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}
	h.registrarEnLedgerSiDisponible(req.IDPedido, "En recoleccion", req.Responsable)

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"estado": "recolectado"})
}