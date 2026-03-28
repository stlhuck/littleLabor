package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/stlhuck/littleLabor/internal/api"
	"github.com/stlhuck/littleLabor/internal/db"

	"net/http"
)

func main() {
	fmt.Println("just doing a test")
	// Initialize database connection
	dbConfig := db.Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "littlelabor"),
		Password: getEnv("DB_PASSWORD", "littlelabor_dev_password"),
		Database: getEnv("DB_NAME", "littlelabor"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	conn, err := db.NewConnection(ctx, dbConfig)
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer conn.Close(context.Background())

	log.Println("Database connected successfully")

	// Initialize repositories
	userRepo := db.NewUserRepository(conn)
	choreRepo := db.NewChoreRepository(conn)
	rewardRepo := db.NewRewardRepository(conn)

	// Set up routes
	router := api.NewRouter(userRepo, choreRepo, rewardRepo)
	router.RegisterRoutes()

	// Start server
	port := ":8080"
	log.Printf("littleLabor API starting on %s\n", port)

	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

// getEnv retrieves an environment variable with a default fallback
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
