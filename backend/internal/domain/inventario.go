package domain

import "github.com/google/uuid"

type Inventario struct {
	IDInventario uuid.UUID `json:"id_inventario"`
	IDProducto   uuid.UUID `json:"id_producto"`
	Stock        int       `json:"stock"`
	Ubicacion    string    `json:"ubicacion"`
}
