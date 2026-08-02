package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/FsaavedraH/colsh/backend/internal/handler"
	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/go-chi/chi/v5"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontró archivo .env, usando variables del sistema")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no está definida")
	}

	pool, err := repository.NewPool(databaseURL)
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}
	defer pool.Close()

	var resultado int
	err = pool.QueryRow(context.Background(), "SELECT 1").Scan(&resultado)
	if err != nil {
		log.Fatal("La conexión falló al hacer una consulta:", err)
	}
	log.Println("Conexión exitosa a la base de datos.")

	pedidoRepo := &repository.PedidoRepository{Pool: pool}
	pedidoHandler := &handler.PedidoHandler{Repo: pedidoRepo}

	inventarioRepo := &repository.InventarioRepository{Pool: pool}
	inventarioHandler := &handler.InventarioHandler{InventarioRepo: inventarioRepo, PedidoRepo: pedidoRepo}

	pickingHandler := &handler.PickingHandler{PedidoRepo: pedidoRepo, InventarioRepo: inventarioRepo}

	empaqueRepo := &repository.EmpaqueRepository{Pool: pool}
	empaqueHandler := &handler.EmpaqueHandler{PedidoRepo: pedidoRepo, EmpaqueRepo: empaqueRepo}

	despachoRepo := &repository.DespachoRepository{Pool: pool}
	despachoHandler := &handler.DespachoHandler{PedidoRepo: pedidoRepo, DespachoRepo: despachoRepo}

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

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Servidor corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}
