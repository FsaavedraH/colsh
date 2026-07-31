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

// RF-11: obtiene la ubicacion esperada de un producto, para comparar contra el escaneo
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
