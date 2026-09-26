package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/FsaavedraH/colsh/backend/internal/security"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type PickingHandler struct {
	PedidoRepo     *repository.PedidoRepository
	InventarioRepo *repository.InventarioRepository
	ReporteRepo    *repository.ReporteRepository
	ProgresoRepo   *repository.ProgresoItemRepository
	ColaLedgerRepo *repository.ColaLedgerRepository
	Ledger         *ledger.LedgerAdapter
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

type SiguienteItemResponse struct {
	Completo   bool                           `json:"completo"`
	Item       *repository.ItemDetallePedido `json:"item,omitempty"`
	Procesados int                            `json:"procesados"`
	TotalItems int                            `json:"total_items"`
}

// GET /api/picking/{id}/siguiente-item - RF-09, RF-10, RF-14. Devuelve el proximo
// producto del pedido que aun no ha sido recolectado. Si "completo" es true, ya
// no quedan items pendientes en Picking.
func (h *PickingHandler) SiguienteItem(w http.ResponseWriter, r *http.Request) {
	idParam := chi.URLParam(r, "id")
	idPedido, err := uuid.Parse(idParam)
	if err != nil {
		http.Error(w, `{"error":"id invalido"}`, http.StatusBadRequest)
		return
	}

	total, err := h.ProgresoRepo.TotalItems(r.Context(), idPedido)
	if err != nil {
		http.Error(w, `{"error":"No se pudo consultar el pedido"}`, http.StatusInternalServerError)
		return
	}

	pendientes, err := h.ProgresoRepo.ContarPendientes(r.Context(), idPedido, "picking")
	if err != nil {
		http.Error(w, `{"error":"No se pudo consultar el progreso"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if pendientes == 0 {
		json.NewEncoder(w).Encode(SiguienteItemResponse{
			Completo:   true,
			Procesados: total,
			TotalItems: total,
		})
		return
	}

	item, err := h.ProgresoRepo.SiguienteItemPendiente(r.Context(), idPedido, "picking")
	if err != nil || item == nil {
		http.Error(w, `{"error":"No se pudo obtener el siguiente item"}`, http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(SiguienteItemResponse{
		Completo:   false,
		Item:       item,
		Procesados: total - pendientes,
		TotalItems: total,
	})
}

type EscanearUbicacionRequest struct {
	IDPedido           string `json:"id_pedido"`
	IDProducto         string `json:"id_producto"`
	UbicacionEscaneada string `json:"ubicacion_escaneada"`
}

// POST /api/picking/escanear-ubicacion - RF-11. La ubicacion escaneada llega como
// un token firmado (HMAC); si la firma no es valida, se rechaza como QR falsificado
// o ajeno al sistema, sin llegar siquiera a comparar contra la ubicacion esperada.
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

	w.Header().Set("Content-Type", "application/json")

	ubicacionEscaneada, err := security.ValidarValorFirmado(req.UbicacionEscaneada)
	if err != nil {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "ubicacion", "qr_invalido", "picking")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Codigo QR invalido o no reconocido por el sistema.",
		})
		return
	}

	ubicacionEsperada, err := h.InventarioRepo.ObtenerUbicacion(r.Context(), idProducto)
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la ubicacion del producto"}`, http.StatusInternalServerError)
		return
	}

	if ubicacionEscaneada != ubicacionEsperada {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "ubicacion", "incorrecto", "picking")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide":            false,
			"ubicacion_esperada":  ubicacionEsperada,
			"ubicacion_escaneada": ubicacionEscaneada,
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

// POST /api/picking/escanear-producto - RF-12, RF-13, RF-26. El producto escaneado
// llega como token firmado; se valida la firma antes de comparar contra el esperado.
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

	idProductoEscaneado, err := security.ValidarValorFirmado(req.IDProductoEscaneado)
	if err != nil {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "qr_invalido", "picking")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Codigo QR invalido o no reconocido por el sistema.",
		})
		return
	}

	if idProductoEscaneado != req.IDProductoEsperado {
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
// el pedido, asi que aqui NO se vuelve a descontar. Registra este item como
// procesado; el pedido solo pasa a "En empaque" cuando TODOS sus items ya
// fueron recolectados (soporta pedidos con varios productos distintos).
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

	if err := h.ProgresoRepo.RegistrarProgreso(r.Context(), idPedido, idProducto, "picking", req.Cantidad, responsable); err != nil {
		http.Error(w, `{"error":"No se pudo registrar la recoleccion: `+err.Error()+`"}`, http.StatusInternalServerError)
		return
	}

	pendientes, err := h.ProgresoRepo.ContarPendientes(r.Context(), idPedido, "picking")
	if err != nil {
		http.Error(w, `{"error":"No se pudo verificar el progreso del pedido"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	if pendientes > 0 {
		json.NewEncoder(w).Encode(map[string]interface{}{
			"estado":     "item_recolectado",
			"completo":   false,
			"pendientes": pendientes,
		})
		return
	}

	if err := h.PedidoRepo.ActualizarEstado(r.Context(), idPedido, "En empaque"); err != nil {
		http.Error(w, `{"error":"No se pudo actualizar el estado del pedido"}`, http.StatusInternalServerError)
		return
	}
	_ = h.PedidoRepo.RegistrarEventoTrazabilidad(r.Context(), idPedido, "En recoleccion", responsable)
	registrarEnLedgerOEncolar(r.Context(), h.Ledger, h.ColaLedgerRepo, idPedido, "En recoleccion", req.Responsable)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"estado":   "En empaque",
		"completo": true,
	})
}