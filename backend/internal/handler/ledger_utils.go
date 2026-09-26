package handler

import (
	"context"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/google/uuid"
)

// registrarEnLedgerOEncolar intenta registrar el evento en el ledger inmediatamente.
// Si el ledger no esta disponible o la escritura falla, el evento se guarda en la
// cola de pendientes CON LA FECHA REAL en que ocurrio (no la fecha del reintento),
// para que un proceso en segundo plano lo sincronice automaticamente mas adelante
// sin perder informacion ni alterar el momento real del evento.
func registrarEnLedgerOEncolar(ctx context.Context, ledgerAdapter *ledger.LedgerAdapter, colaRepo *repository.ColaLedgerRepository, idPedido uuid.UUID, estado, responsable string) {
	idEvento := uuid.New().String()
	fechaEvento := time.Now().UTC()
	fechaStr := fechaEvento.Format(time.RFC3339)

	if ledgerAdapter != nil {
		if err := ledgerAdapter.RegistrarEnLedger(ctx, idEvento, idPedido.String(), estado, fechaStr, responsable); err == nil {
			return // registrado con exito, no hace falta encolar
		}
	}

	if colaRepo != nil {
		_ = colaRepo.EncolarPendiente(ctx, idEvento, idPedido, estado, fechaEvento, responsable)
	}
}