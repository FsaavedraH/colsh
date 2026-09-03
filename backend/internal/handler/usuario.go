package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
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

type CrearUsuarioRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
	Rol      string `json:"rol"`
}

type CrearUsuarioResponse struct {
	IDUsuario string `json:"id_usuario"`
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Rol       string `json:"rol"`
}

var rolesPermitidosPorAdmin = map[string]bool{
	"Picking":       true,
	"Empaque":       true,
	"Transportista": true,
	"Administrador": true,
}

// POST /api/usuarios - RF-27. Solo Administrador. Crea usuarios operativos con el rol indicado.
// El rol "Cliente" no se permite aqui: el cliente se crea via /api/auth/registro (autorregistro).
func (h *UsuarioHandler) CrearUsuarioAdmin(w http.ResponseWriter, r *http.Request) {
	var req CrearUsuarioRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Email == "" || req.Password == "" || req.Rol == "" {
		http.Error(w, `{"error":"Nombre, email, contrasena y rol son obligatorios"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"error":"La contrasena debe tener al menos 6 caracteres"}`, http.StatusBadRequest)
		return
	}

	if !rolesPermitidosPorAdmin[req.Rol] {
		http.Error(w, `{"error":"Rol invalido. Debe ser Picking, Empaque, Transportista o Administrador"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"No se pudo procesar la contrasena"}`, http.StatusInternalServerError)
		return
	}

	idUsuario := uuid.New()
	err = h.UsuarioRepo.CrearUsuario(r.Context(), idUsuario, req.Nombre, req.Email, req.Rol, string(hash))
	if err != nil {
		http.Error(w, `{"error":"No se pudo crear el usuario (el correo ya podria existir)"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(CrearUsuarioResponse{
		IDUsuario: idUsuario.String(),
		Nombre:    req.Nombre,
		Email:     req.Email,
		Rol:       req.Rol,
	})
}