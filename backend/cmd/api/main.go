package main

import (
	"context"
	"errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

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

	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      router, // Middleware
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	// Canal para escuchar señales del S.O. (Ctrl+C, Docker Stop, etc)
	quit := make(chan os.Signal, 1) // Pipe
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Levantar server con goroutine
	go func() {
		log.Printf("Servidor corriendo en el puerto %s...\n", cfg.Port)
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("Error crítico en el servidor HTTP: %v", err)
		}
	}()

	<-quit
	log.Println("Señal de apagado recibida. Iniciando Graceful Shutdown...")

	// Si el servidor tarda más de 10 seg en terminar peticiones, forzamos apagado. Evitar ataque Slowloris
	ctxShutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("El servidor forzó el apagado debido a un error: %v", err)
	}

	log.Println("Servidor apagado correctamente.")

}
