package handler

import (
	"encoding/json"
	"net/http"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	UsuarioRepo *repository.UsuarioRepository
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginResponse struct {
	IDUsuario string `json:"id_usuario"`
	Nombre    string `json:"nombre"`
	Email     string `json:"email"`
	Rol       string `json:"rol"`
}

// POST /api/auth/login - ruta publica. Autentica un usuario y devuelve sus datos si las credenciales son correctas.
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"Email y contrasena son obligatorios"}`, http.StatusBadRequest)
		return
	}

	usuario, err := h.UsuarioRepo.BuscarPorEmail(r.Context(), req.Email)
	if err != nil {
		http.Error(w, `{"error":"Credenciales invalidas"}`, http.StatusUnauthorized)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(usuario.PasswordHash), []byte(req.Password))
	if err != nil {
		http.Error(w, `{"error":"Credenciales invalidas"}`, http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		IDUsuario: usuario.IDUsuario.String(),
		Nombre:    usuario.Nombre,
		Email:     usuario.Email,
		Rol:       usuario.Rol,
	})
}

type RegistroRequest struct {
	Nombre   string `json:"nombre"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// POST /api/auth/registro - ruta publica. SIEMPRE crea el usuario con rol "Cliente",
// sin importar que envie el request (evita que alguien se auto-asigne otro rol).
func (h *AuthHandler) Registro(w http.ResponseWriter, r *http.Request) {
	var req RegistroRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"error":"JSON invalido"}`, http.StatusBadRequest)
		return
	}

	if req.Nombre == "" || req.Email == "" || req.Password == "" {
		http.Error(w, `{"error":"Nombre, email y contrasena son obligatorios"}`, http.StatusBadRequest)
		return
	}

	if len(req.Password) < 6 {
		http.Error(w, `{"error":"La contrasena debe tener al menos 6 caracteres"}`, http.StatusBadRequest)
		return
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, `{"error":"No se pudo procesar la contrasena"}`, http.StatusInternalServerError)
		return
	}

	idUsuario := uuid.New()
	err = h.UsuarioRepo.CrearUsuario(r.Context(), idUsuario, req.Nombre, req.Email, "Cliente", string(hash))
	if err != nil {
		http.Error(w, `{"error":"No se pudo crear el usuario (el correo ya podria existir)"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(LoginResponse{
		IDUsuario: idUsuario.String(),
		Nombre:    req.Nombre,
		Email:     req.Email,
		Rol:       "Cliente",
	})
}