package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

/* Config lee var entorno, valida que no falte nada critico (fail fast), y devolver estructura tipada */

type AppConfig struct {
	Port  string
	DBUrl string
}

// Load lee las variables de entorno y construye la configuración
func Load() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Println("Aviso: No se encontró archivo .env, leyendo variables del sistema...")
	}

	// Configurar el Puerto
	port := os.Getenv("PORT")
	if port == "" { // fallback a 8080 si no existe
		port = "8080"
	}

	// Validar Base de Datos
	dbUser := getEnvOrFatal("DB_USER")
	dbPass := getEnvOrFatal("DB_PASSWORD")
	dbHost := getEnvOrFatal("DB_HOST")
	dbPort := getEnvOrFatal("DB_PORT")
	dbName := getEnvOrFatal("DB_NAME")

	// Construir el Data Source Name
	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable&timezone=America/Santiago",
		dbUser, dbPass, dbHost, dbPort, dbName)

	return &AppConfig{
		Port:  port,
		DBUrl: dsn,
	}
}

// getEnvOrFatal asegura que la variable exista, si no, mata la aplicación
func getEnvOrFatal(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("Error Crítico: La variable de entorno %s no está definida", key)
	}
	return val
}
