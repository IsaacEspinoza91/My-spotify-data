package main

import (
	"context"
	"log"
	"net/http"

	"github.com/IsaacEspinoza91/My-spotify-data/internal/config"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/database"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/handler"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/repository"
	"github.com/IsaacEspinoza91/My-spotify-data/internal/service"
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

	repo := repository.NewSpotifyRepository(dbPool)
	svc := service.NewSpotifyService(repo)
	router := handler.NewRouter(svc)

	log.Println("Servidor iniciado en :8080")
	http.ListenAndServe(":8080", router)

}
