package repository

import (
	"context"

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