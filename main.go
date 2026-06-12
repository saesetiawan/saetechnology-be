package main

import (
	"context"
	"saetechnology-be/internal/config"
	"saetechnology-be/internal/infrastructure/database"
	"log"
	"os"

	"saetechnology-be/internal/di"
	"saetechnology-be/internal/pkg/tracing"
)

func main() {
	ctx := context.Background()
	cfg := config.Load()

	_, shutdownTracer, err := tracing.InitTracer(
		ctx,
		"sae-technology-solution-backend",
		cfg.TracerEndPoint,
	)
	if err != nil {
		log.Fatal(err)
	}

	defer func() {
		if err := shutdownTracer(ctx); err != nil {
			log.Println("failed to shutdown tracer:", err)
		}
	}()

	db, err := di.InitPostgresql()
	if err != nil {
		panic(err)
	}
	sql, err := db.DB()
	if err != nil {
		panic(err)
	}
	err = database.RunMigrationPostgresql(sql, "migrations", os.Getenv("DB_NAME"))
	if err != nil {
		panic(err)
	}

	server := di.InitServer()
	if cfg.QueueRegisterEmail != "" {
		consumer := di.InitConsumer()
		registerEmailConsumer := di.InitRegisterEmailConsumer()
		go func() {
			if err := consumer.Start(ctx, cfg.QueueRegisterEmail, registerEmailConsumer.Handle); err != nil {
				log.Println("failed to start register email consumer:", err)
			}
		}()
	}

	log.Println("server running on", server.Addr)

	if err := server.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
