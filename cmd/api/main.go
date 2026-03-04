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
	"github.com/omerbeden/paymentgateway/internal/adapter/handler/messaging"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/cache"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/config"
	"github.com/omerbeden/paymentgateway/internal/infrastructure/database"
	infrakafka "github.com/omerbeden/paymentgateway/internal/infrastructure/queue/kafka"
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

	admin, err := infrakafka.NewAdminClient(appConfig.Kafka.Brokers)
	if err != nil {
		log.Fatal("kafka admin: %v", err)
	}

	if err := admin.EnsureTopics(context.Background(), infrakafka.DefaultTopics()); err != nil {
		log.Fatal("kafka admin: ensure topics: %v", err)
	}

	producer, err := infrakafka.NewKafkaProducer(*appConfig.Kafka)
	if err != nil {
		log.Fatal("kafka producer: %v", err)
	}

	defer producer.Close()
	publisher := messaging.NewKafkaPublisher(producer)

	router := routes.SetupRoutes(db, redis, appConfig, publisher)

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
