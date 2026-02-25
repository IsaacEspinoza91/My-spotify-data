package main

import (
	"context"
	"log"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/config"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/database"
)

func main() {
	// Cargar Configuración Centralizada
	cfg := config.Load()

	// Inicializar DB
	ctx := context.Background()
	dbPool, err := database.NewPostgresConnection(ctx, cfg.DBUrl)
	if err != nil {
		log.Fatalf("Error fatal conectando a la base de datos: %v", err)
	}
	defer dbPool.Close()
	log.Println("Conectado a PostgreSQL exitosamente")

}
