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

	"github.com/go-playground/validator/v10"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/redis/go-redis/v9"
	"github.com/shariski/room-booking/internal/config"
	"github.com/shariski/room-booking/internal/handler"
	"github.com/shariski/room-booking/internal/repository"
	"github.com/shariski/room-booking/internal/usecase"
)

func main() {
	config := config.Load()

	dsn := fmt.Sprintf("postgres://%s:@%s:%s/%s?sslmode=disable", config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})

	validator := validator.New()

	roomRepository := repository.NewRoomRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	roomUsecase := usecase.NewRoomUsecase(roomRepository, rdb)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepository)

	roomHandler := handler.NewRoomHandler(roomUsecase, validator)
	bookingHandler := handler.NewBookingHandler(bookingUsecase, validator)

	mux := http.NewServeMux()
	mux.HandleFunc("GET /rooms", roomHandler.List)
	mux.HandleFunc("GET /rooms/{id}", roomHandler.Get)
	mux.HandleFunc("POST /bookings", bookingHandler.Create)
	mux.HandleFunc("DELETE /bookings/{id}", bookingHandler.Delete)

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

	if err := db.Close(); err != nil {
		log.Fatal("DB closing error:", err)
	}

	if err := rdb.Close(); err != nil {
		log.Fatal("Redis closing error:", err)
	}

	fmt.Println("Server stopped")
}
