package main

import (
	"log"
	"net/http"
	"os"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/FsaavedraH/colsh/backend/internal/server"
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

	if err := server.VerificarConexionDB(pool); err != nil {
		log.Fatal("La conexión a la base de datos falló:", err)
	}
	log.Println("Conexión exitosa a la base de datos.")

	ledgerAdapter := server.ConectarLedger()
	if ledgerAdapter != nil {
		defer ledgerAdapter.Cerrar()
	}

	router := server.NuevoRouter(pool, ledgerAdapter)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	log.Println("Servidor corriendo en el puerto", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}
