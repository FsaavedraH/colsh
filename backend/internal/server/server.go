package server

import (
	"context"
	"log"

	"github.com/FsaavedraH/colsh/backend/internal/handler"
	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	appmw "github.com/FsaavedraH/colsh/backend/internal/middleware"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/jackc/pgx/v5/pgxpool"
)

func NuevoRouter(pool *pgxpool.Pool, ledgerAdapter *ledger.LedgerAdapter) *chi.Mux {
	pedidoRepo := &repository.PedidoRepository{Pool: pool}
	inventarioRepo := &repository.InventarioRepository{Pool: pool}

	pedidoHandler := &handler.PedidoHandler{Repo: pedidoRepo, InventarioRepo: inventarioRepo}
	inventarioHandler := &handler.InventarioHandler{InventarioRepo: inventarioRepo, PedidoRepo: pedidoRepo}

	reporteRepo := &repository.ReporteRepository{Pool: pool}
	reporteHandler := &handler.ReporteHandler{ReporteRepo: reporteRepo}

	usuarioRepo := &repository.UsuarioRepository{Pool: pool}
	usuarioHandler := &handler.UsuarioHandler{UsuarioRepo: usuarioRepo}
	authHandler := &handler.AuthHandler{UsuarioRepo: usuarioRepo}

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

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-User-Role"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// Autenticacion - rutas publicas, sin RBAC (son el paso previo a tener un rol)
	r.Post("/api/auth/login", authHandler.Login)
	r.Post("/api/auth/registro", authHandler.Registro)

	// Pedidos
	r.With(appmw.RequireRole("Cliente")).Post("/api/pedidos", pedidoHandler.CrearPedido)
	r.With(appmw.RequireRole("Cliente", "Picking", "Empaque", "Transportista", "Administrador")).
		Get("/api/pedidos/{id}", pedidoHandler.ConsultarPedido)
	r.With(appmw.RequireRole("Cliente")).Get("/api/mis-pedidos", pedidoHandler.ListarMisPedidos)

	// Catalogo
	r.With(appmw.RequireRole("Cliente", "Picking", "Empaque", "Transportista", "Administrador")).
		Get("/api/productos", inventarioHandler.ListarCatalogo)

	// Usuarios
	r.With(appmw.RequireRole("Administrador")).Get("/api/usuarios", usuarioHandler.ListarUsuarios)
	r.With(appmw.RequireRole("Administrador")).Post("/api/usuarios", usuarioHandler.CrearUsuarioAdmin)

	// Inventario
	r.With(appmw.RequireRole("Picking", "Administrador")).Post("/api/inventario/validar", inventarioHandler.ValidarInventario)
	r.With(appmw.RequireRole("Picking", "Administrador")).Post("/api/inventario/actualizar", inventarioHandler.ActualizarInventario)

	// Picking
	r.With(appmw.RequireRole("Picking", "Administrador")).Get("/api/picking", pickingHandler.ListarOrdenes)
	r.With(appmw.RequireRole("Picking", "Administrador")).Get("/api/picking/historial", pickingHandler.ListarHistorial)
	r.With(appmw.RequireRole("Picking")).Post("/api/picking/escanear-ubicacion", pickingHandler.EscanearUbicacion)
	r.With(appmw.RequireRole("Picking")).Post("/api/picking/escanear-producto", pickingHandler.EscanearProducto)
	r.With(appmw.RequireRole("Picking")).Post("/api/recoleccion", pickingHandler.ConfirmarRecoleccion)

	// Empaque
	r.With(appmw.RequireRole("Empaque")).Get("/api/empaque", empaqueHandler.ListarOrdenes)
	r.With(appmw.RequireRole("Empaque")).Get("/api/empaque/historial", empaqueHandler.ListarHistorial)
	r.With(appmw.RequireRole("Empaque")).Post("/api/empaque/recepcion", empaqueHandler.RecepcionEmpaque)
	r.With(appmw.RequireRole("Empaque")).Post("/api/empaque/escanear", empaqueHandler.EscanearValidacion)
	r.With(appmw.RequireRole("Empaque")).Post("/api/empaque", empaqueHandler.ConfirmarEmpaque)

	// Despacho y entrega
	r.With(appmw.RequireRole("Transportista")).Get("/api/despacho", despachoHandler.ListarOrdenes)
	r.With(appmw.RequireRole("Transportista")).Get("/api/despacho/historial", despachoHandler.ListarHistorial)
	r.With(appmw.RequireRole("Transportista")).Post("/api/despacho", despachoHandler.GenerarDespacho)
	r.With(appmw.RequireRole("Transportista")).Post("/api/entrega", despachoHandler.ConfirmarEntrega)

	// Trazabilidad
	r.With(appmw.RequireRole("Cliente", "Picking", "Empaque", "Transportista", "Administrador")).
		Get("/api/trazabilidad/{id_pedido}", trazabilidadHandler.ConsultarTrazabilidad)

	// Reportes
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/pedidos", reporteHandler.ListarPedidos)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/tiempos", reporteHandler.TiemposPorEtapa)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/incidencias", reporteHandler.IndicadorIncidencias)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/conteo-estados", reporteHandler.ConteoPorEstado)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/productos-top", reporteHandler.ProductosMasVendidos)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/pedidos-por-dia", reporteHandler.PedidosPorDia)

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