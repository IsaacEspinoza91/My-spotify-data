package database

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Abre conexion, hace ping y retorna pool de conexiones
func NewPostgresConnection(ctx context.Context, dsn string) (*pgxpool.Pool, error) {
	// pgxpool maneja el conjunto de conexiones abiertas
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, fmt.Errorf("error al crear el pool de conexiones: %v", err)
	}

	// Verificar conexion ping
	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("error al hacer ping a la base de datos: %v", err)
	}

	return pool, nil
}
