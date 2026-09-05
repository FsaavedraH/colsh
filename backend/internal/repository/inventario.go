package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type InventarioRepository struct {
	Pool *pgxpool.Pool
}

type FaltanteInventario struct {
	IDProducto      string `json:"id_producto"`
	StockDisponible int    `json:"stock_disponible"`
	CantidadPedida  int    `json:"cantidad_pedida"`
}

func (r *InventarioRepository) ValidarStock(ctx context.Context, productos []ProductoPedidoInput) (bool, []FaltanteInventario, error) {
	var faltantes []FaltanteInventario

	for _, p := range productos {
		idProducto, err := uuid.Parse(p.IDProducto)
		if err != nil {
			return false, nil, err
		}

		var stock int
		err = r.Pool.QueryRow(ctx,
			`SELECT COALESCE(SUM(stock), 0) FROM inventario WHERE id_producto = $1`, idProducto,
		).Scan(&stock)
		if err != nil {
			return false, nil, err
		}

		if stock < p.Cantidad {
			faltantes = append(faltantes, FaltanteInventario{
				IDProducto:      p.IDProducto,
				StockDisponible: stock,
				CantidadPedida:  p.Cantidad,
			})
		}
	}

	disponible := len(faltantes) == 0
	return disponible, faltantes, nil
}

// ReservarStock: descuenta el stock de todos los productos del pedido de forma
// atomica (todo o nada). Si algun producto no tiene stock suficiente, no se
// descuenta nada y se informa cuales faltaron. Se usa al CREAR el pedido, para
// comprometer el material desde ese momento y evitar sobreventa entre pedidos.
func (r *InventarioRepository) ReservarStock(ctx context.Context, productos []ProductoPedidoInput) (bool, []FaltanteInventario, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return false, nil, err
	}
	defer tx.Rollback(ctx)

	var faltantes []FaltanteInventario

	for _, p := range productos {
		idProducto, err := uuid.Parse(p.IDProducto)
		if err != nil {
			return false, nil, err
		}

		tag, err := tx.Exec(ctx,
			`UPDATE inventario SET stock = stock - $1 WHERE id_producto = $2 AND stock >= $1`,
			p.Cantidad, idProducto,
		)
		if err != nil {
			return false, nil, err
		}

		if tag.RowsAffected() == 0 {
			var stockActual int
			tx.QueryRow(ctx, `SELECT COALESCE(SUM(stock), 0) FROM inventario WHERE id_producto = $1`, idProducto).Scan(&stockActual)
			faltantes = append(faltantes, FaltanteInventario{
				IDProducto:      p.IDProducto,
				StockDisponible: stockActual,
				CantidadPedida:  p.Cantidad,
			})
		}
	}

	if len(faltantes) > 0 {
		// no se hace commit: al salir del bloque, el defer tx.Rollback(ctx) revierte todo
		return false, faltantes, nil
	}

	if err := tx.Commit(ctx); err != nil {
		return false, nil, err
	}

	return true, nil, nil
}

// LiberarStock: devuelve al inventario el material que se habia reservado
// para un pedido que ahora se cancela.
func (r *InventarioRepository) LiberarStock(ctx context.Context, productos []ProductoPedidoInput) error {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	for _, p := range productos {
		idProducto, err := uuid.Parse(p.IDProducto)
		if err != nil {
			return err
		}
		_, err = tx.Exec(ctx,
			`UPDATE inventario SET stock = stock + $1 WHERE id_producto = $2`,
			p.Cantidad, idProducto,
		)
		if err != nil {
			return err
		}
	}

	return tx.Commit(ctx)
}

func (r *InventarioRepository) ActualizarStock(ctx context.Context, idProducto uuid.UUID, cantidad int) error {
	_, err := r.Pool.Exec(ctx,
		`UPDATE inventario SET stock = stock - $1 WHERE id_producto = $2`,
		cantidad, idProducto,
	)
	return err
}

func (r *InventarioRepository) ObtenerUbicacion(ctx context.Context, idProducto uuid.UUID) (string, error) {
	var ubicacion string
	err := r.Pool.QueryRow(ctx,
		`SELECT ubicacion FROM inventario WHERE id_producto = $1 LIMIT 1`, idProducto,
	).Scan(&ubicacion)
	if err != nil {
		return "", err
	}
	return ubicacion, nil
}

type ProductoCatalogo struct {
	IDProducto string `json:"id_producto"`
	Nombre     string `json:"nombre"`
	Stock      int    `json:"stock"`
	Ubicacion  string `json:"ubicacion"`
}

// ListarCatalogo devuelve todos los productos con su stock disponible.
func (r *InventarioRepository) ListarCatalogo(ctx context.Context) ([]ProductoCatalogo, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT p.id_producto, p.nombre, COALESCE(SUM(i.stock), 0) as stock, MAX(i.ubicacion) as ubicacion
		FROM producto p
		LEFT JOIN inventario i ON i.id_producto = p.id_producto
		GROUP BY p.id_producto, p.nombre
		ORDER BY p.nombre ASC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []ProductoCatalogo
	for rows.Next() {
		var pc ProductoCatalogo
		if err := rows.Scan(&pc.IDProducto, &pc.Nombre, &pc.Stock, &pc.Ubicacion); err != nil {
			return nil, err
		}
		resultado = append(resultado, pc)
	}
	return resultado, nil
}

type CompraResumen struct {
	IDCompra             uuid.UUID `json:"id_compra"`
	IDProducto           uuid.UUID `json:"id_producto"`
	NombreProducto       string    `json:"nombre_producto"`
	Cantidad             int       `json:"cantidad"`
	CostoUnitarioMomento *float64  `json:"costo_unitario_momento"`
	CostoTotal           *float64  `json:"costo_total"`
	Fecha                time.Time `json:"fecha"`
	NombreResponsable    string    `json:"nombre_responsable"`
}

// RegistrarCompra: RF-07 (version simplificada). Suma stock al producto seleccionado
// y deja constancia en el historial de compras, con el costo tomado del catalogo.
func (r *InventarioRepository) RegistrarCompra(ctx context.Context, idProducto uuid.UUID, cantidad int, responsable uuid.UUID) (uuid.UUID, error) {
	tx, err := r.Pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	var costoUnitario *float64
	err = tx.QueryRow(ctx, `SELECT costo_unitario FROM producto WHERE id_producto = $1`, idProducto).Scan(&costoUnitario)
	if err != nil {
		return uuid.Nil, err
	}

	var costoTotal *float64
	if costoUnitario != nil {
		total := *costoUnitario * float64(cantidad)
		costoTotal = &total
	}

	idCompra := uuid.New()
	_, err = tx.Exec(ctx,
		`INSERT INTO compra_inventario (id_compra, id_producto, cantidad, costo_unitario_momento, costo_total, responsable)
		 VALUES ($1, $2, $3, $4, $5, $6)`,
		idCompra, idProducto, cantidad, costoUnitario, costoTotal, responsable,
	)
	if err != nil {
		return uuid.Nil, err
	}

	_, err = tx.Exec(ctx,
		`UPDATE inventario SET stock = stock + $1 WHERE id_producto = $2`,
		cantidad, idProducto,
	)
	if err != nil {
		return uuid.Nil, err
	}

	if err := tx.Commit(ctx); err != nil {
		return uuid.Nil, err
	}

	return idCompra, nil
}

// ListarCompras: historial completo de ingresos de compra, mas reciente primero.
func (r *InventarioRepository) ListarCompras(ctx context.Context) ([]CompraResumen, error) {
	rows, err := r.Pool.Query(ctx, `
		SELECT c.id_compra, c.id_producto, p.nombre, c.cantidad, c.costo_unitario_momento, c.costo_total, c.fecha, u.nombre
		FROM compra_inventario c
		JOIN producto p ON p.id_producto = c.id_producto
		JOIN usuario u ON u.id_usuario = c.responsable
		ORDER BY c.fecha DESC
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []CompraResumen
	for rows.Next() {
		var cr CompraResumen
		if err := rows.Scan(&cr.IDCompra, &cr.IDProducto, &cr.NombreProducto, &cr.Cantidad, &cr.CostoUnitarioMomento, &cr.CostoTotal, &cr.Fecha, &cr.NombreResponsable); err != nil {
			return nil, err
		}
		resultado = append(resultado, cr)
	}
	return resultado, nil
}