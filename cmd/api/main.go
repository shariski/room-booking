package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/shariski/room-booking/internal/config"
)

func main() {
	config := config.Load()

	dsn := fmt.Sprintf("postgres://%s:@%s:%s/%s?sslmode=disable", config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(db)

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})
	fmt.Println(rdb)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello world!")
	})

	server := &http.Server{
		Addr:    config.Port,
		Handler: mux,
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		fmt.Println("Server running on :8080")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Fatal("Shutdown error:", err)
	}

	// TODO: close DB, redis
	fmt.Println("Server stopped")
}
