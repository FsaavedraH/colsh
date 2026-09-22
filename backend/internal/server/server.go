package server

import (
	"context"
	"log"
	"time"

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
	progresoRepo := &repository.ProgresoItemRepository{Pool: pool}
	colaLedgerRepo := &repository.ColaLedgerRepository{Pool: pool}

	pedidoHandler := &handler.PedidoHandler{
		Repo:           pedidoRepo,
		InventarioRepo: inventarioRepo,
		ColaLedgerRepo: colaLedgerRepo,
		Ledger:         ledgerAdapter,
	}
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
		ProgresoRepo:   progresoRepo,
		ColaLedgerRepo: colaLedgerRepo,
		Ledger:         ledgerAdapter,
	}

	empaqueRepo := &repository.EmpaqueRepository{Pool: pool}
	empaqueHandler := &handler.EmpaqueHandler{
		PedidoRepo:     pedidoRepo,
		EmpaqueRepo:    empaqueRepo,
		ReporteRepo:    reporteRepo,
		ProgresoRepo:   progresoRepo,
		ColaLedgerRepo: colaLedgerRepo,
		Ledger:         ledgerAdapter,
	}

	despachoRepo := &repository.DespachoRepository{Pool: pool}
	despachoHandler := &handler.DespachoHandler{
		PedidoRepo:     pedidoRepo,
		DespachoRepo:   despachoRepo,
		ColaLedgerRepo: colaLedgerRepo,
		Ledger:         ledgerAdapter,
	}

	trazabilidadHandler := &handler.TrazabilidadHandler{Ledger: ledgerAdapter}
	ledgerHandler := &handler.LedgerHandler{ColaLedgerRepo: colaLedgerRepo}

	iniciarReintentoLedger(ledgerAdapter, colaLedgerRepo)

	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "https://colsh.vercel.app", "https://*.vercel.app"},
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
	r.With(appmw.RequireRole("Cliente", "Administrador")).Post("/api/pedidos/{id}/cancelar", pedidoHandler.CancelarPedido)

	// Catalogo
	r.With(appmw.RequireRole("Cliente", "Picking", "Empaque", "Transportista", "Administrador")).
		Get("/api/productos", inventarioHandler.ListarCatalogo)

	// Usuarios
	r.With(appmw.RequireRole("Administrador")).Get("/api/usuarios", usuarioHandler.ListarUsuarios)
	r.With(appmw.RequireRole("Administrador")).Post("/api/usuarios", usuarioHandler.CrearUsuarioAdmin)

	// Inventario
	r.With(appmw.RequireRole("Picking", "Administrador")).Post("/api/inventario/validar", inventarioHandler.ValidarInventario)
	r.With(appmw.RequireRole("Picking", "Administrador")).Post("/api/inventario/actualizar", inventarioHandler.ActualizarInventario)
	r.With(appmw.RequireRole("Administrador")).Post("/api/inventario/compras", inventarioHandler.RegistrarCompra)
	r.With(appmw.RequireRole("Administrador")).Get("/api/inventario/compras", inventarioHandler.ListarCompras)
	r.With(appmw.RequireRole("Administrador")).Get("/api/inventario/qr", inventarioHandler.ListarParaCodigosQR)
	r.With(appmw.RequireRole("Administrador")).Get("/api/inventario/demanda-espera", inventarioHandler.ListarDemandaEnEspera)

	// Picking
	r.With(appmw.RequireRole("Picking", "Administrador")).Get("/api/picking", pickingHandler.ListarOrdenes)
	r.With(appmw.RequireRole("Picking", "Administrador")).Get("/api/picking/historial", pickingHandler.ListarHistorial)
	r.With(appmw.RequireRole("Picking")).Post("/api/picking/iniciar", pickingHandler.IniciarPicking)
	r.With(appmw.RequireRole("Picking")).Get("/api/picking/{id}/siguiente-item", pickingHandler.SiguienteItem)
	r.With(appmw.RequireRole("Picking")).Post("/api/picking/escanear-ubicacion", pickingHandler.EscanearUbicacion)
	r.With(appmw.RequireRole("Picking")).Post("/api/picking/escanear-producto", pickingHandler.EscanearProducto)
	r.With(appmw.RequireRole("Picking")).Post("/api/recoleccion", pickingHandler.ConfirmarRecoleccion)

	// Empaque
	r.With(appmw.RequireRole("Empaque")).Get("/api/empaque", empaqueHandler.ListarOrdenes)
	r.With(appmw.RequireRole("Empaque")).Get("/api/empaque/historial", empaqueHandler.ListarHistorial)
	r.With(appmw.RequireRole("Empaque")).Post("/api/empaque/recepcion", empaqueHandler.RecepcionEmpaque)
	r.With(appmw.RequireRole("Empaque")).Get("/api/empaque/{id}/siguiente-item", empaqueHandler.SiguienteItem)
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
	r.With(appmw.RequireRole("Administrador")).Get("/api/ledger/pendientes", ledgerHandler.ListarPendientes)

	// Reportes
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/pedidos", reporteHandler.ListarPedidos)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/tiempos", reporteHandler.TiemposPorEtapa)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/incidencias", reporteHandler.IndicadorIncidencias)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/conteo-estados", reporteHandler.ConteoPorEstado)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/productos-top", reporteHandler.ProductosMasVendidos)
	r.With(appmw.RequireRole("Administrador")).Get("/api/reportes/pedidos-por-dia", reporteHandler.PedidosPorDia)

	return r
}

// iniciarReintentoLedger lanza un proceso en segundo plano que revisa cada 30
// segundos si hay eventos pendientes de sincronizar con el ledger (por ejemplo,
// porque Hyperledger Fabric estaba caido en el momento del evento original), e
// intenta reenviarlos SIN alterar la fecha real en que ocurrieron.
func iniciarReintentoLedger(ledgerAdapter *ledger.LedgerAdapter, colaRepo *repository.ColaLedgerRepository) {
	go func() {
		ticker := time.NewTicker(30 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			if ledgerAdapter == nil {
				continue // sin conexion al ledger todavia; se reintentara en el siguiente ciclo
			}

			pendientes, err := colaRepo.ListarPendientes(context.Background())
			if err != nil || len(pendientes) == 0 {
				continue
			}

			log.Printf("[ledger] Reintentando sincronizar %d evento(s) pendiente(s)...\n", len(pendientes))

			for _, evento := range pendientes {
				fechaStr := evento.FechaEventoReal.Format(time.RFC3339)
				err := ledgerAdapter.RegistrarEnLedger(
					context.Background(),
					evento.IDEventoOriginal,
					evento.IDPedido.String(),
					evento.Estado,
					fechaStr,
					evento.Responsable,
				)
				if err != nil {
					colaRepo.RegistrarIntentoFallido(context.Background(), evento.ID)
					continue
				}
				colaRepo.EliminarPendiente(context.Background(), evento.ID)
				log.Printf("[ledger] Sincronizado: pedido %s, estado '%s' (ocurrido %s)\n", evento.IDPedido, evento.Estado, fechaStr)
			}
		}
	}()
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