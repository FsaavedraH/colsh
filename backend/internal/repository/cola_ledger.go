package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ColaLedgerRepository struct {
	Pool *pgxpool.Pool
}

type EventoLedgerPendiente struct {
	ID               uuid.UUID  `json:"id"`
	IDEventoOriginal string     `json:"id_evento_original"`
	IDPedido         uuid.UUID  `json:"id_pedido"`
	Estado           string     `json:"estado"`
	FechaEventoReal  time.Time  `json:"fecha_evento_real"`
	Responsable      string     `json:"responsable"`
	Intentos         int        `json:"intentos"`
	UltimoIntento    *time.Time `json:"ultimo_intento"`
	CreadoEn         time.Time  `json:"creado_en"`
}

// EncolarPendiente guarda un evento que no se pudo registrar en el ledger en el
// momento en que ocurrio, conservando la fecha REAL del evento (no la fecha en
// que eventualmente se logre sincronizar).
func (r *ColaLedgerRepository) EncolarPendiente(ctx context.Context, idEventoOriginal string, idPedido uuid.UUID, estado string, fechaEventoReal time.Time, responsable string) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO eventos_ledger_pendientes (id_evento_original, id_pedido, estado, fecha_evento_real, responsable)
		VALUES ($1, $2, $3, $4, $5)
	`, idEventoOriginal, idPedido, estado, fechaEventoReal, responsable)
	return err
}

// ListarPendientes devuelve todos los eventos que aun no se han logrado sincronizar.
func (r *ColaLedgerRepository) ListarPendientes(ctx context.Context) ([]EventoLedgerPendiente, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT id, id_evento_original, id_pedido, estado, fecha_evento_real, responsable, intentos, ultimo_intento, creado_en
		FROM eventos_ledger_pendientes
		ORDER BY fecha_evento_real ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []EventoLedgerPendiente
	for rows.Next() {
		var e EventoLedgerPendiente
		if err := rows.Scan(&e.ID, &e.IDEventoOriginal, &e.IDPedido, &e.Estado, &e.FechaEventoReal, &e.Responsable, &e.Intentos, &e.UltimoIntento, &e.CreadoEn); err != nil {
			return nil, err
		}
		resultado = append(resultado, e)
	}
	return resultado, nil
}

// RegistrarIntentoFallido incrementa el contador de intentos, para no perder cuenta
// de cuantas veces se ha intentado sincronizar un evento que sigue fallando.
func (r *ColaLedgerRepository) RegistrarIntentoFallido(ctx context.Context, id uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		UPDATE eventos_ledger_pendientes SET intentos = intentos + 1, ultimo_intento = now() WHERE id = $1
	`, id)
	return err
}

// EliminarPendiente quita el evento de la cola una vez que se logro sincronizar con exito.
func (r *ColaLedgerRepository) EliminarPendiente(ctx context.Context, id uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM eventos_ledger_pendientes WHERE id = $1`, id)
	return err
}