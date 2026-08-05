package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/FsaavedraH/colsh/backend/internal/ledger"
	"github.com/google/uuid"
)

func main() {
	adapter, err := ledger.NuevoLedgerAdapter()
	if err != nil {
		log.Fatal("Error al conectar con Fabric:", err)
	}
	defer adapter.Cerrar()

	fmt.Println("Conexion con Fabric establecida correctamente.")

	ctx := context.Background()
	idEvento := uuid.New().String()
	idPedido := "65d8d459-ef57-4b70-b67c-4c45624962ec"

	err = adapter.RegistrarEnLedger(ctx, idEvento, idPedido, "Prueba desde backend Go", time.Now().Format(time.RFC3339), "test-responsable")
	if err != nil {
		log.Fatal("Error al registrar evento:", err)
	}

	fmt.Println("Evento registrado con exito. ID:", idEvento)

	resultado, err := adapter.ConsultarEvento(ctx, idEvento)
	if err != nil {
		log.Fatal("Error al consultar evento:", err)
	}

	fmt.Println("Evento consultado desde el ledger:", string(resultado))
}