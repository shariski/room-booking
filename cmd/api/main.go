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
	"github.com/shariski/room-booking/internal/middleware"
	"github.com/shariski/room-booking/internal/repository"
	"github.com/shariski/room-booking/internal/usecase"
)

func main() {
	config := config.Load()

	dsn := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", config.DBUser, config.DBPassword, config.DBHost, config.DBPort, config.DBName)
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatal(err)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     config.RedisAddr,
		Password: config.RedisPass,
	})

	validator := validator.New()

	userRepository := repository.NewUserRepository(db)
	roomRepository := repository.NewRoomRepository(db)
	bookingRepository := repository.NewBookingRepository(db)

	userUsecase := usecase.NewUserUsecase(config, userRepository)
	roomUsecase := usecase.NewRoomUsecase(roomRepository, rdb)
	bookingUsecase := usecase.NewBookingUsecase(bookingRepository)

	userHandler := handler.NewUserHandler(userUsecase, validator)
	roomHandler := handler.NewRoomHandler(roomUsecase, validator)
	bookingHandler := handler.NewBookingHandler(bookingUsecase, validator)

	authMiddleware := middleware.AuthMiddleware(config.JWTSecret)

	mux := http.NewServeMux()
	mux.HandleFunc("POST /users", userHandler.Create)
	mux.HandleFunc("POST /login", userHandler.Login)
	mux.Handle("GET /rooms", authMiddleware(http.HandlerFunc(roomHandler.List)))
	mux.Handle("GET /rooms/{id}", authMiddleware(http.HandlerFunc(roomHandler.Get)))
	mux.Handle("POST /bookings", authMiddleware(http.HandlerFunc(bookingHandler.Create)))
	mux.Handle("DELETE /bookings/{id}", authMiddleware(http.HandlerFunc(bookingHandler.Delete)))

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
