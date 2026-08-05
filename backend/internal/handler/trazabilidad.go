package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/go-chi/chi/v5"
)

type TrazabilidadHandler struct {
	Ledger *ledger.LedgerAdapter
}

// GET /api/trazabilidad/{id_pedido} - RF-24
func (h *TrazabilidadHandler) ConsultarTrazabilidad(w http.ResponseWriter, r *http.Request) {
	idPedido := chi.URLParam(r, "id_pedido")

	w.Header().Set("Content-Type", "application/json")

	if h.Ledger == nil {
		// RT-06, Escenario 3: si el ledger no esta disponible, se responde con un mensaje controlado
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{
			"error":   "El servicio de trazabilidad (ledger) no esta disponible en este momento.",
			"detalle": "El sistema operativo sigue funcionando con normalidad; solo la consulta de historial inmutable esta afectada.",
		})
		return
	}

	resultado, err := h.Ledger.ConsultarHistorialPedido(r.Context(), idPedido)
	if err != nil {
		http.Error(w, `{"error":"No se pudo consultar el historial en el ledger"}`, http.StatusInternalServerError)
		return
	}

	// El chaincode ya devuelve JSON valido, lo pasamos directo
	w.Write(resultado)
}
