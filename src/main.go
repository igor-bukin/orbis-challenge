package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/orbis-challenge/src/config"
	"github.com/orbis-challenge/src/handlers"
	"github.com/orbis-challenge/src/persistence/postgres"
	"github.com/orbis-challenge/src/persistence/redis"
	"github.com/orbis-challenge/src/services"
	"github.com/orbis-challenge/src/validator"
	"github.com/rs/cors"
	"github.com/sirupsen/logrus"
)

// defaultConfigPath defines a path to JSON-config file
const defaultConfigPath = "config.json"

func main() {
	fmt.Println("Start API")

	err := config.Load(defaultConfigPath)
	if err != nil {
		log.Fatalf("Failed to initialize Config: %s", err.Error())
	}

	// setup log-level
	logLevel, err := logrus.ParseLevel(config.Config.LogPreset)
	if err != nil {
		logrus.Fatal("cannot parse log level", err)
	}
	logrus.SetLevel(logLevel)

	err = postgres.Load(&config.Config.Postgres, logrus.StandardLogger())
	if err != nil {
		logrus.Fatalf("cannot connect to the Postgres server with config [%+v]: %v",
			config.Config.Postgres, err)
	}

	err = redis.Load(config.Config.Redis)
	if err != nil {
		logrus.Fatalf("cannot connect to the Redis server with config [%+v]: %v",
			config.Config.Redis, err)
	}

	if err = validator.Load(); err != nil {
		logrus.Fatalf("cannot initialize validator: %v", err)
	}

	if err = services.Load(redis.GetRedis(), postgres.GetDB(), &config.Config); err != nil {
		logrus.Fatalf("cannot initialize services: %v", err)
	}

	corsInstance := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:*", "https://localhost:*",
			"http://127.0.0.1:*", "https://127.0.0.1:*"},
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodDelete, http.MethodOptions, http.MethodPut, http.MethodPatch},
		AllowCredentials: true,
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"*"},
	})

	server := &http.Server{
		Addr:    config.Config.ListenURL,
		Handler: corsInstance.Handler(handlers.NewRouter()),
	}

	err = server.ListenAndServe()
	if err != nil {
		logrus.Error("Failed to initialize HTTP server", "error", err)
		os.Exit(1)
	}
}
