package domain

import "github.com/google/uuid"

type Usuario struct {
	IDUsuario    uuid.UUID `json:"id_usuario"`
	Nombre       string    `json:"nombre"`
	Email        string    `json:"email"`
	Rol          string    `json:"rol"`
	PasswordHash string    `json:"-"`
}
