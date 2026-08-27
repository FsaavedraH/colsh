package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
)

type UsuarioHandler struct {
	UsuarioRepo *repository.UsuarioRepository
}

// GET /api/usuarios - RF-27
func (h *UsuarioHandler) ListarUsuarios(w http.ResponseWriter, r *http.Request) {
	usuarios, err := h.UsuarioRepo.ListarUsuarios(r.Context())
	if err != nil {
		http.Error(w, `{"error":"No se pudo obtener la lista de usuarios"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(usuarios)
}