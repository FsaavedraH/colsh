package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/FsaavedraH/colsh/backend/internal/repository"
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

	fmt.Println("Conexión exitosa a la base de datos. Resultado de prueba:", resultado)
}
