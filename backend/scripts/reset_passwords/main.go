package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

// Este script asigna una contraseña real (hasheada con bcrypt) a cada usuario
// existente, reemplazando los valores "placeholder" que tenian desde el seed.
// Contraseña para TODOS los usuarios en esta fase de prototipo: colsh2026
func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontro .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("No se pudo conectar:", err)
	}
	defer pool.Close()

	const passwordDemo = "colsh2026"
	hash, err := bcrypt.GenerateFromPassword([]byte(passwordDemo), bcrypt.DefaultCost)
	if err != nil {
		log.Fatal("Error al generar hash:", err)
	}

	ctx := context.Background()
	tag, err := pool.Exec(ctx, "UPDATE usuario SET password_hash = $1", string(hash))
	if err != nil {
		log.Fatal("Error al actualizar contraseñas:", err)
	}

	fmt.Printf("Contraseñas actualizadas para %d usuarios.\n", tag.RowsAffected())
	fmt.Println("Contraseña de demo para todos:", passwordDemo)
}