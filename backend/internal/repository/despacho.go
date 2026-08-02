package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type DespachoRepository struct {
	Pool *pgxpool.Pool
}

// RF-20, RF-22: registra un evento de despacho o entrega en la trazabilidad
func (r *DespachoRepository) RegistrarEvento(ctx context.Context, idPedido uuid.UUID, estado string, responsable uuid.UUID) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO evento_trazabilidad (id_evento, id_pedido, estado, fecha, responsable)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), idPedido, estado, time.Now(), responsable,
	)
	return err
}
