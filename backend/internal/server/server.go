package server

import (
	"context"
	"log"

	"github.com/FsaavedraH/colsh/backend/internal/handler"
	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NuevoRouter(pool *pgxpool.Pool, ledgerAdapter *ledger.LedgerAdapter) *chi.Mux {
	pedidoRepo := &repository.PedidoRepository{Pool: pool}
	pedidoHandler := &handler.PedidoHandler{Repo: pedidoRepo}

	inventarioRepo := &repository.InventarioRepository{Pool: pool}
	inventarioHandler := &handler.InventarioHandler{InventarioRepo: inventarioRepo, PedidoRepo: pedidoRepo}

	reporteRepo := &repository.ReporteRepository{Pool: pool}
	reporteHandler := &handler.ReporteHandler{ReporteRepo: reporteRepo}

	pickingHandler := &handler.PickingHandler{
		PedidoRepo:     pedidoRepo,
		InventarioRepo: inventarioRepo,
		ReporteRepo:    reporteRepo,
		Ledger:         ledgerAdapter,
	}

	empaqueRepo := &repository.EmpaqueRepository{Pool: pool}
	empaqueHandler := &handler.EmpaqueHandler{
		PedidoRepo:  pedidoRepo,
		EmpaqueRepo: empaqueRepo,
		ReporteRepo: reporteRepo,
		Ledger:      ledgerAdapter,
	}

	despachoRepo := &repository.DespachoRepository{Pool: pool}
	despachoHandler := &handler.DespachoHandler{PedidoRepo: pedidoRepo, DespachoRepo: despachoRepo, Ledger: ledgerAdapter}

	trazabilidadHandler := &handler.TrazabilidadHandler{Ledger: ledgerAdapter}

	r := chi.NewRouter()

	r.Post("/api/pedidos", pedidoHandler.CrearPedido)
	r.Get("/api/pedidos/{id}", pedidoHandler.ConsultarPedido)

	r.Post("/api/inventario/validar", inventarioHandler.ValidarInventario)
	r.Post("/api/inventario/actualizar", inventarioHandler.ActualizarInventario)

	r.Get("/api/picking", pickingHandler.ListarOrdenes)
	r.Post("/api/picking/escanear-ubicacion", pickingHandler.EscanearUbicacion)
	r.Post("/api/picking/escanear-producto", pickingHandler.EscanearProducto)
	r.Post("/api/recoleccion", pickingHandler.ConfirmarRecoleccion)

	r.Post("/api/empaque/recepcion", empaqueHandler.RecepcionEmpaque)
	r.Post("/api/empaque/escanear", empaqueHandler.EscanearValidacion)
	r.Post("/api/empaque", empaqueHandler.ConfirmarEmpaque)

	r.Post("/api/despacho", despachoHandler.GenerarDespacho)
	r.Post("/api/entrega", despachoHandler.ConfirmarEntrega)

	r.Get("/api/trazabilidad/{id_pedido}", trazabilidadHandler.ConsultarTrazabilidad)

	r.Get("/api/reportes/pedidos", reporteHandler.ListarPedidos)
	r.Get("/api/reportes/tiempos", reporteHandler.TiemposPorEtapa)
	r.Get("/api/reportes/incidencias", reporteHandler.IndicadorIncidencias)

	return r
}

func ConectarLedger() *ledger.LedgerAdapter {
	adapter, err := ledger.NuevoLedgerAdapter()
	if err != nil {
		log.Println("Advertencia: no se pudo conectar con el ledger (Hyperledger Fabric). El sistema seguira funcionando, pero las operaciones de trazabilidad inmutable no estaran disponibles. Detalle:", err)
		return nil
	}
	log.Println("Conexion exitosa con el ledger (Hyperledger Fabric).")
	return adapter
}

func VerificarConexionDB(pool *pgxpool.Pool) error {
	var resultado int
	err := pool.QueryRow(context.Background(), "SELECT 1").Scan(&resultado)
	return err
}
