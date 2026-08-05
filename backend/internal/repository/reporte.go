package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ReporteRepository struct {
	Pool *pgxpool.Pool
}

func (r *ReporteRepository) RegistrarIntentoEscaneo(ctx context.Context, idPedido uuid.UUID, tipo, resultado, etapa string) error {
	_, err := r.Pool.Exec(ctx,
		`INSERT INTO intento_escaneo (id_intento, id_pedido, tipo, resultado, etapa)
		 VALUES ($1, $2, $3, $4, $5)`,
		uuid.New(), idPedido, tipo, resultado, etapa,
	)
	return err
}

type PedidoReporte struct {
	IDPedido      uuid.UUID `json:"id_pedido"`
	FechaCreacion time.Time `json:"fecha_creacion"`
	Estado        string    `json:"estado"`
	NombreCliente string    `json:"nombre_cliente"`
}

// RF-28: listado de pedidos con filtros opcionales por estado y rango de fecha.
func (r *ReporteRepository) ListarPedidosFiltrado(ctx context.Context, estado, fechaDesde, fechaHasta string) ([]PedidoReporte, error) {
	query := `
		SELECT p.id_pedido, p.fecha_creacion, p.estado, u.nombre
		FROM pedido p
		JOIN usuario u ON u.id_usuario = p.id_cliente
		WHERE ($1 = '' OR p.estado = $1)
		  AND ($2 = '' OR p.fecha_creacion >= $2::timestamp)
		  AND ($3 = '' OR p.fecha_creacion <= $3::timestamp)
		ORDER BY p.fecha_creacion DESC
	`
	rows, err := r.Pool.Query(ctx, query, estado, fechaDesde, fechaHasta)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []PedidoReporte
	for rows.Next() {
		var p PedidoReporte
		if err := rows.Scan(&p.IDPedido, &p.FechaCreacion, &p.Estado, &p.NombreCliente); err != nil {
			return nil, err
		}
		resultado = append(resultado, p)
	}
	return resultado, nil
}

type TiempoPorEtapa struct {
	Etapa             string  `json:"etapa"`
	TiempoPromedioMin float64 `json:"tiempo_promedio_minutos"`
	CantidadPedidos   int     `json:"cantidad_pedidos"`
}

// RF-29: tiempo promedio (en minutos) entre cada transicion de estado, usando evento_trazabilidad.
func (r *ReporteRepository) TiemposPorEtapa(ctx context.Context) ([]TiempoPorEtapa, error) {
	query := `
		WITH eventos_ordenados AS (
			SELECT
				id_pedido,
				estado,
				fecha,
				LAG(fecha) OVER (PARTITION BY id_pedido ORDER BY fecha) AS fecha_anterior
			FROM evento_trazabilidad
		)
		SELECT
			estado AS etapa,
			AVG(EXTRACT(EPOCH FROM (fecha - fecha_anterior)) / 60) AS tiempo_promedio_minutos,
			COUNT(*) AS cantidad_pedidos
		FROM eventos_ordenados
		WHERE fecha_anterior IS NOT NULL
		GROUP BY estado
		ORDER BY tiempo_promedio_minutos ASC
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []TiempoPorEtapa
	for rows.Next() {
		var t TiempoPorEtapa
		if err := rows.Scan(&t.Etapa, &t.TiempoPromedioMin, &t.CantidadPedidos); err != nil {
			return nil, err
		}
		resultado = append(resultado, t)
	}
	return resultado, nil
}

type IndicadorIncidencias struct {
	Etapa               string  `json:"etapa"`
	TotalIntentos       int     `json:"total_intentos"`
	IntentosIncorrectos int     `json:"intentos_incorrectos"`
	TasaIncidencia      float64 `json:"tasa_incidencia_porcentaje"`
}

// Indicador operativo: tasa de escaneos incorrectos por etapa (equivalente generado
// por el sistema al IIT del planteamiento del problema, seccion 1.2).
func (r *ReporteRepository) IndicadorIncidenciasEscaneo(ctx context.Context) ([]IndicadorIncidencias, error) {
	query := `
		SELECT
			etapa,
			COUNT(*) AS total_intentos,
			COUNT(*) FILTER (WHERE resultado = 'incorrecto') AS intentos_incorrectos
		FROM intento_escaneo
		GROUP BY etapa
	`
	rows, err := r.Pool.Query(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var resultado []IndicadorIncidencias
	for rows.Next() {
		var ind IndicadorIncidencias
		if err := rows.Scan(&ind.Etapa, &ind.TotalIntentos, &ind.IntentosIncorrectos); err != nil {
			return nil, err
		}
		if ind.TotalIntentos > 0 {
			ind.TasaIncidencia = (float64(ind.IntentosIncorrectos) / float64(ind.TotalIntentos)) * 100
		}
		resultado = append(resultado, ind)
	}
	return resultado, nil
}
