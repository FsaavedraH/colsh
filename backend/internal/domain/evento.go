package domain

import (
	"time"

	"github.com/google/uuid"
)

type EventoTrazabilidad struct {
	IDEvento    uuid.UUID `json:"id_evento"`
	IDPedido    uuid.UUID `json:"id_pedido"`
	Estado      string    `json:"estado"`
	Fecha       time.Time `json:"fecha"`
	Responsable uuid.UUID `json:"responsable"`
	TxID        string    `json:"tx_id,omitempty"`
}
