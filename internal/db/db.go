package db

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

func Connect() {
	connStr := "postgres://postgres:postgres@localhost:5432/job_dispatcher?sslmode=disable"
	
	// Configure the connection pool
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("Unable to parse connection string: %v", err)
	}
	
	// Set pool configuration
	config.MaxConns = 10
	config.MinConns = 2
	config.MaxConnLifetime = time.Hour
	config.MaxConnIdleTime = time.Minute * 30
	
	// Create the connection pool
	DB, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Fatalf("Unable to create connection pool: %v", err)
	}

	fmt.Println("Connected to the database successfully!")

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	
	err = DB.Ping(ctx)
	if err != nil {
		log.Fatalf("Unable to ping database: %v", err)
	}
	fmt.Println("Pinged the database successfully!")
}

func Close() {
	if DB != nil {
		DB.Close()
	}
}