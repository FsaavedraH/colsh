package domain

import "github.com/google/uuid"

type Producto struct {
	IDProducto uuid.UUID `json:"id_producto"`
	Nombre     string    `json:"nombre"`
}
