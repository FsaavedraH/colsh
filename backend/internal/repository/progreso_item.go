package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProgresoItemRepository struct {
	Pool *pgxpool.Pool
}

// SiguienteItemPendiente devuelve el proximo item del pedido que AUN NO ha sido
// procesado en la etapa indicada ("picking" o "empaque"). Si devuelve nil, ya
// no quedan items pendientes en esa etapa.
func (r *ProgresoItemRepository) SiguienteItemPendiente(ctx context.Context, idPedido uuid.UUID, etapa string) (*ItemDetallePedido, error) {
	var item ItemDetallePedido
	err := r.Pool.QueryRow(ctx, `
		SELECT dp.id_producto, p.nombre, dp.cantidad, COALESCE(MAX(i.ubicacion), '')
		FROM detalle_pedido dp
		JOIN producto p ON p.id_producto = dp.id_producto
		LEFT JOIN inventario i ON i.id_producto = p.id_producto
		WHERE dp.id_pedido = $1
		  AND NOT EXISTS (
		    SELECT 1 FROM progreso_item_pedido pi
		    WHERE pi.id_pedido = dp.id_pedido
		      AND pi.id_producto = dp.id_producto
		      AND pi.etapa = $2
		  )
		GROUP BY dp.id_producto, p.nombre, dp.cantidad
		ORDER BY p.nombre
		LIMIT 1
	`, idPedido, etapa).Scan(&item.IDProducto, &item.Nombre, &item.Cantidad, &item.Ubicacion)

	if err != nil {
		return nil, nil // no hay mas items pendientes (o no se encontro ninguno)
	}
	return &item, nil
}

// ContarPendientes devuelve cuantos items del pedido faltan por procesar en la etapa dada.
func (r *ProgresoItemRepository) ContarPendientes(ctx context.Context, idPedido uuid.UUID, etapa string) (int, error) {
	var total int
	err := r.Pool.QueryRow(ctx, `
		SELECT COUNT(*)
		FROM detalle_pedido dp
		WHERE dp.id_pedido = $1
		  AND NOT EXISTS (
		    SELECT 1 FROM progreso_item_pedido pi
		    WHERE pi.id_pedido = dp.id_pedido
		      AND pi.id_producto = dp.id_producto
		      AND pi.etapa = $2
		  )
	`, idPedido, etapa).Scan(&total)
	return total, err
}

// RegistrarProgreso marca un item como procesado en la etapa indicada.
func (r *ProgresoItemRepository) RegistrarProgreso(ctx context.Context, idPedido, idProducto uuid.UUID, etapa string, cantidad int, responsable uuid.UUID) error {
	_, err := r.Pool.Exec(ctx, `
		INSERT INTO progreso_item_pedido (id_pedido, id_producto, etapa, cantidad_procesada, responsable)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (id_pedido, id_producto, etapa) DO UPDATE
		SET cantidad_procesada = EXCLUDED.cantidad_procesada, fecha = now(), responsable = EXCLUDED.responsable
	`, idPedido, idProducto, etapa, cantidad, responsable)
	return err
}

// TotalItems devuelve cuantos productos distintos tiene el pedido en total.
func (r *ProgresoItemRepository) TotalItems(ctx context.Context, idPedido uuid.UUID) (int, error) {
	var total int
	err := r.Pool.QueryRow(ctx, `SELECT COUNT(*) FROM detalle_pedido WHERE id_pedido = $1`, idPedido).Scan(&total)
	return total, err
}

// LimpiarProgreso borra el progreso de una etapa para un pedido (por si se necesita reiniciar).
func (r *ProgresoItemRepository) LimpiarProgreso(ctx context.Context, idPedido uuid.UUID, etapa string) error {
	_, err := r.Pool.Exec(ctx, `DELETE FROM progreso_item_pedido WHERE id_pedido = $1 AND etapa = $2`, idPedido, etapa)
	return err
}