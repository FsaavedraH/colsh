package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type EmpaqueRepository struct {
	Pool *pgxpool.Pool
}

// RF-18: registra el evento de empaque con fecha, hora y responsable
func (r *EmpaqueRepository) RegistrarEmpaque(ctx context.Context, idPedido uuid.UUID, responsable uuid.UUID) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO evento_trazabilidad (id_evento, id_pedido, estado, fecha, responsable)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), idPedido, "En empaque", time.Now(), responsable,
	)
	return err
}
