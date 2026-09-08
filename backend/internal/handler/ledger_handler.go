package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
)

type LedgerHandler struct {
	ColaLedgerRepo *repository.ColaLedgerRepository
}

// GET /api/ledger/pendientes - Solo Administrador. Muestra los eventos que aun
// no se han logrado sincronizar con el ledger (por caidas de Hyperledger Fabric),
// con la fecha REAL en que ocurrieron y cuantas veces se ha intentado reenviarlos.
func (h *LedgerHandler) ListarPendientes(w http.ResponseWriter, r *http.Request) {
	pendientes, err := h.ColaLedgerRepo.ListarPendientes(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la cola de eventos pendientes"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(pendientes)
}