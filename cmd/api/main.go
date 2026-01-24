package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omerbeden/paymentgateway/internal/adapter/handler/http/routes"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/cache"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/database"
)

func main() {

	appConfig := config.Load()
	db, err := database.NewPostgres(appConfig.DatabaseDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	redis := cache.NewRedis(appConfig.RedisAddr)
	defer redis.Close()

	router := routes.SetupRoutes(db, redis, appConfig)

	srv := &http.Server{
		Addr:    ":" + appConfig.ServerPort,
		Handler: router,
	}

	go func() {
		log.Println("Starting server on :8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Listen error: %s", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %s", err)
	}

	log.Println("Server exited gracefully")

}
