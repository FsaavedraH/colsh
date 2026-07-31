package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Advertencia: no se encontró .env")
	}

	databaseURL := os.Getenv("DATABASE_URL")
	pool, err := pgxpool.New(context.Background(), databaseURL)
	if err != nil {
		log.Fatal("No se pudo conectar:", err)
	}
	defer pool.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Nombre del producto: ")
	nombre, _ := reader.ReadString('\n')
	nombre = strings.TrimSpace(nombre)

	fmt.Print("Stock inicial: ")
	stockStr, _ := reader.ReadString('\n')
	stock, err := strconv.Atoi(strings.TrimSpace(stockStr))
	if err != nil {
		log.Fatal("Stock invalido:", err)
	}

	fmt.Print("Ubicacion (ej. B-12-03-02): ")
	ubicacion, _ := reader.ReadString('\n')
	ubicacion = strings.TrimSpace(ubicacion)

	ctx := context.Background()
	idProducto := uuid.New()

	_, err = pool.Exec(ctx, `INSERT INTO producto (id_producto, nombre) VALUES ($1, $2)`, idProducto, nombre)
	if err != nil {
		log.Fatal("Error al crear producto:", err)
	}

	idInventario := uuid.New()
	_, err = pool.Exec(ctx,
		`INSERT INTO inventario (id_inventario, id_producto, stock, ubicacion) VALUES ($1, $2, $3, $4)`,
		idInventario, idProducto, stock, ubicacion,
	)
	if err != nil {
		log.Fatal("Error al crear inventario:", err)
	}

	fmt.Println("\nProducto creado con exito.")
	fmt.Println("ID Producto:", idProducto)
	fmt.Println("Nombre:", nombre)
	fmt.Println("Stock:", stock)
	fmt.Println("Ubicacion:", ubicacion)
}