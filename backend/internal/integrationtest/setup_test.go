package integrationtest

import (
	"log"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
	"github.com/FsaavedraH/colsh/backend/internal/server"
	"github.com/joho/godotenv"
)

var testServer *httptest.Server

// TestMain levanta el backend real (sin Fabric) una sola vez, para todas las pruebas.
func TestMain(m *testing.M) {
	_ = godotenv.Load("../../.env")

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL no esta definida para las pruebas")
	}

	pool, err := repository.NewPool(databaseURL)
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos de pruebas:", err)
	}

	router := server.NuevoRouter(pool, nil) // ledger nil: las pruebas no dependen de Fabric
	testServer = httptest.NewServer(router)
	defer testServer.Close()

	code := m.Run()
	pool.Close()
	os.Exit(code)
}
