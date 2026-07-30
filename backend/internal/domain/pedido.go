package domain

import (
	"time"

	"github.com/google/uuid"
)

type Pedido struct {
	IDPedido         uuid.UUID `json:"id_pedido"`
	FechaCreacion    time.Time `json:"fecha_creacion"`
	Estado           string    `json:"estado"`
	IDCliente        uuid.UUID `json:"id_cliente"`
	DireccionEntrega string    `json:"direccion_entrega"`
}
