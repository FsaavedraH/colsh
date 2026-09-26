package repository

import (
	"context"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/domain"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PedidoRepository struct {
	Pool *pgxpool.Pool
}

type ProductoPedidoInput struct {
	IDProducto string
	Cantidad   int
}

type PedidoPickingResumen struct {
	IDPedido      uuid.UUID `json:"id_pedido"`
	FechaCreacion time.Time `json:"fecha_creacion"`
	Estado        string    `json:"estado"`
	NombreCliente string    `json:"nombre_cliente"`
	TotalItems    int       `json:"total_items"`
}

type ItemDetallePedido struct {
	IDProducto string `json:"id_producto"`
	Nombre     string `json:"nombre"`
	Cantidad   int    `json:"cantidad"`
	Ubicacion  string `json:"ubicacion"`
}

func (r *PedidoRepository) Crear(ctx context.Context, pedido *domain.Pedido, productos []ProductoPedidoInput) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	_, err = tx.Exec(ctx,
		`INSERT INTO pedido (id_pedido, fecha_creacion, estado, id_cliente, direccion_entrega)
		 VALUES ($1, $2, $3, $4, $5)`,
		pedido.IDPedido, pedido.FechaCreacion, pedido.Estado, pedido.IDCliente, pedido.DireccionEntrega,
	)
	if err != nil {
		return err
	}

	for _, p := range productos {
		idProducto, err := uuid.Parse(p.IDProducto)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx,
			`INSERT INTO detalle_pedido (id_detalle, id_pedido, id_producto, cantidad)
			 VALUES ($1, $2, $3, $4)`,
			uuid.New(), pedido.IDPedido, idProducto, p.Cantidad,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *PedidoRepository) ConsultarPorID(ctx context.Context, id uuid.UUID) (*domain.Pedido, error) {
	var p domain.Pedido
	err := r.Pool.QueryRow(ctx,
		`SELECT p.id_pedido, p.fecha_creacion, p.estado, p.id_cliente, u.nombre, p.direccion_entrega
		 FROM pedido p
		 JOIN usuario u ON u.id_usuario = p.id_cliente
		 WHERE p.id_pedido = $1`, id,
	).Scan(&p.IDPedido, &p.FechaCreacion, &p.Estado, &p.IDCliente, &p.NombreCliente, &p.DireccionEntrega)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *PedidoRepository) ActualizarEstado(ctx context.Context, id uuid.UUID, nuevoEstado string) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE pedido SET estado = $1 WHERE id_pedido = $2`,
		nuevoEstado, id,
	)
	return err
}

// RegistrarEventoTrazabilidad inserta un evento en la tabla local evento_trazabilidad,
// independiente del registro en el ledger de Hyperledger Fabric (que se maneja aparte
// via registrarEnLedgerOEncolar). Sin esto, el historial completo del pedido (RF-24)
// y los reportes de tiempo por etapa quedan incompletos para los estados que solo
// pasaban por el ledger: creacion, cancelacion y confirmacion de recoleccion.
func (r *PedidoRepository) RegistrarEventoTrazabilidad(ctx context.Context, idPedido uuid.UUID, estado string, responsable uuid.UUID) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO evento_trazabilidad (id_evento, id_pedido, estado, fecha, responsable)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), idPedido, estado, time.Now().UTC(), responsable,
	)
	return err
}

func (r *PedidoRepository) ObtenerProductosDelPedido(ctx context.Context, idPedido uuid.UUID) ([]ProductoPedidoInput, error) {
	rows, err := r.Pool.Query(ctx,
		`SELECT id_producto, cantidad FROM detalle_pedido WHERE id_pedido = $1`, idPedido,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var productos []ProductoPedidoInput
	for rows.Next() {
		var idProducto uuid.UUID
		var cantidad int
		if err := rows.Scan(&idProducto, &cantidad); err != nil {
			return nil, err
		}
		productos = append(productos, ProductoPedidoInput{
			IDProducto: idProducto.String(),
			Cantidad:   cantidad,
		})
	}
	return productos, nil
}

func (r *PedidoRepository) ObtenerDetalleConNombres(ctx context.Context, idPedido uuid.UUID) ([]ItemDetallePedido, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id_producto, p.nombre, dp.cantidad, COALESCE(MAX(i.ubicacion), '')
		FROM detalle_pedido dp
		JOIN producto p ON p.id_producto = dp.id_producto
		LEFT JOIN inventario i ON i.id_producto = p.id_producto
		WHERE dp.id_pedido = $1
		GROUP BY p.id_producto, p.nombre, dp.cantidad
		ORDER BY p.nombre
	`, idPedido)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var items []ItemDetallePedido
	for rows.Next() {
		var item ItemDetallePedido
		if err := rows.Scan(&item.IDProducto, &item.Nombre, &item.Cantidad, &item.Ubicacion); err != nil {
			return nil, err
		}
		items = append(items, item)
	}
	return items, nil
}

func (r *PedidoRepository) listarPorEstados(ctx context.Context, estados []string) ([]PedidoPickingResumen, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id_pedido, p.fecha_creacion, p.estado, u.nombre,
		       COALESCE(SUM(dp.cantidad), 0) as total_items
		FROM pedido p
		JOIN usuario u ON u.id_usuario = p.id_cliente
		LEFT JOIN detalle_pedido dp ON dp.id_pedido = p.id_pedido
		WHERE p.estado = ANY($1)
		GROUP BY p.id_pedido, p.fecha_creacion, p.estado, u.nombre
		ORDER BY p.fecha_creacion ASC
	`, estados)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []PedidoPickingResumen
	for rows.Next() {
		var pr PedidoPickingResumen
		if err := rows.Scan(&pr.IDPedido, &pr.FechaCreacion, &pr.Estado, &pr.NombreCliente, &pr.TotalItems); err != nil {
			return nil, err
		}
		resultado = append(resultado, pr)
	}
	return resultado, nil
}

func (r *PedidoRepository) ListarPorEstado(ctx context.Context, estado string) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{estado})
}

func (r *PedidoRepository) ListarParaPicking(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"Pendiente", "En recoleccion"})
}

func (r *PedidoRepository) ListarHistorialPicking(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"En empaque", "En despacho", "Entregado"})
}

func (r *PedidoRepository) ListarParaEmpaque(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"En empaque"})
}

func (r *PedidoRepository) ListarHistorialEmpaque(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"En despacho", "Entregado"})
}

func (r *PedidoRepository) ListarParaDespacho(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"En despacho"})
}

func (r *PedidoRepository) ListarHistorialDespacho(ctx context.Context) ([]PedidoPickingResumen, error) {
	return r.listarPorEstados(ctx, []string{"Entregado"})
}

func (r *PedidoRepository) ListarPorCliente(ctx context.Context, idCliente uuid.UUID) ([]PedidoPickingResumen, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id_pedido, p.fecha_creacion, p.estado, u.nombre,
		       COALESCE(SUM(dp.cantidad), 0) as total_items
		FROM pedido p
		JOIN usuario u ON u.id_usuario = p.id_cliente
		LEFT JOIN detalle_pedido dp ON dp.id_pedido = p.id_pedido
		WHERE p.id_cliente = $1
		GROUP BY p.id_pedido, p.fecha_creacion, p.estado, u.nombre
		ORDER BY p.fecha_creacion DESC
	`, idCliente)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []PedidoPickingResumen
	for rows.Next() {
		var pr PedidoPickingResumen
		if err := rows.Scan(&pr.IDPedido, &pr.FechaCreacion, &pr.Estado, &pr.NombreCliente, &pr.TotalItems); err != nil {
			return nil, err
		}
		resultado = append(resultado, pr)
	}
	return resultado, nil
}