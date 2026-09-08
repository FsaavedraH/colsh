package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/FsaavedraH/colsh/backend/internal/security"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type EmpaqueHandler struct {
	PedidoRepo   *repository.PedidoRepository
	EmpaqueRepo  *repository.EmpaqueRepository
	ReporteRepo  *repository.ReporteRepository
	ProgresoRepo *repository.ProgresoItemRepository
	Ledger       *ledger.LedgerAdapter
}

func (h *EmpaqueHandler) registrarEnLedgerSiDisponible(idPedido, estado, responsable string) {
	if h.Ledger == nil {
		return
	}
	idEvento := uuid.New().String()
	fecha := time.Now().Format(time.RFC3339)
	_ = h.Ledger.RegistrarEnLedger(context.Background(), idEvento, idPedido, estado, fecha, responsable)
}

// GET /api/empaque - RF-15
func (h *EmpaqueHandler) ListarOrdenes(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarParaEmpaque(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la lista de ordenes"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
}

// GET /api/empaque/historial
func (h *EmpaqueHandler) ListarHistorial(w http.ResponseWriter, r *http.Request) {
	ordenes, err := h.PedidoRepo.ListarHistorialEmpaque(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener el historial"}`, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(ordenes)
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

// GET /api/empaque/{id}/siguiente-item - RF-17. Devuelve el proximo producto
// del pedido que aun no ha sido validado en empaque.
func (h *EmpaqueHandler) SiguienteItem(w http.ResponseWriter, r *http.Request) {
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

	pendientes, err := h.ProgresoRepo.ContarPendientes(r.Context(), idPedido, "empaque")
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

	item, err := h.ProgresoRepo.SiguienteItemPendiente(r.Context(), idPedido, "empaque")
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

type EscanearValidacionEmpaqueRequest struct {
	IDPedido            string `json:"id_pedido"`
	IDProductoEsperado  string `json:"id_producto_esperado"`
	IDProductoEscaneado string `json:"id_producto_escaneado"`
	Responsable         string `json:"responsable"`
}

// POST /api/empaque/escanear - RF-17, RF-26. El producto escaneado llega como
// token firmado (HMAC); se valida la firma antes de comparar contra el esperado.
// Si coincide, registra este item como validado en empaque (soporta pedidos
// con varios productos distintos, uno a la vez).
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

	idProductoEsperado, err := uuid.Parse(req.IDProductoEsperado)
	if err != nil {
		http.Error(w, `{"error":"id_producto_esperado invalido"}`, http.StatusBadRequest)
		return
	}

	responsable, err := uuid.Parse(req.Responsable)
	if err != nil {
		http.Error(w, `{"error":"responsable invalido"}`, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	idProductoEscaneado, err := security.ValidarValorFirmado(req.IDProductoEscaneado)
	if err != nil {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "qr_invalido", "empaque")
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Codigo QR invalido o no reconocido por el sistema.",
		})
		return
	}

	if idProductoEscaneado != req.IDProductoEsperado {
		h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "incorrecto", "empaque")

		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"coincide": false,
			"mensaje":  "Producto incorrecto. El producto escaneado no corresponde. Intenta nuevamente.",
		})
		return
	}

	h.ReporteRepo.RegistrarIntentoEscaneo(r.Context(), idPedido, "producto", "correcto", "empaque")

	if err := h.ProgresoRepo.RegistrarProgreso(r.Context(), idPedido, idProductoEsperado, "empaque", 0, responsable); err != nil {
		http.Error(w, `{"error":"No se pudo registrar el progreso de empaque"}`, http.StatusInternalServerError)
		return
	}

	pendientes, _ := h.ProgresoRepo.ContarPendientes(r.Context(), idPedido, "empaque")

	json.NewEncoder(w).Encode(map[string]interface{}{
		"coincide":   true,
		"mensaje":    "Producto validado correctamente para empaque",
		"pendientes": pendientes,
	})
}

type ConfirmarEmpaqueRequest struct {
	IDPedido    string `json:"id_pedido"`
	Responsable string `json:"responsable"`
}

// POST /api/empaque - RF-18, RF-24. Solo finaliza el empaque completo del pedido
// si TODOS sus items ya fueron validados via /api/empaque/escanear.
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

	pendientes, err := h.ProgresoRepo.ContarPendientes(r.Context(), idPedido, "empaque")
	if err != nil {
		http.Error(w, `{"error":"No se pudo verificar el progreso del pedido"}`, http.StatusInternalServerError)
		return
	}
	if pendientes > 0 {
		http.Error(w, `{"error":"Aun quedan productos por validar en empaque"}`, http.StatusConflict)
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