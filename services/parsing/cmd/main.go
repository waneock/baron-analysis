package main

import (
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/parsing/internal/app"
	"skinbaron-analyzer/services/parsing/internal/client/baron"
	"skinbaron-analyzer/services/parsing/internal/config"
	"skinbaron-analyzer/services/parsing/internal/usecase"
	"time"

	transportgrpc "skinbaron-analyzer/services/parsing/internal/transport/grpc"
)

func main() {
	// config
	cfg := config.MustLoad()

	// logger
	log := logger.MustLoad(cfg.Env)

	// database
	dbCfg := makeDBConfigData(cfg)

	db, err := db.New(*dbCfg)
	if err != nil {
		log.Error("error when trying to create a new database",
			"error", err)
		os.Exit(1)
	}

	defer db.Close()

	// repo
	repos := app.NewRepositories(db)
	if repos == nil {
		log.Error("error when trying to create a repository")
		os.Exit(1)
	}

	log.Info("app successfully initialized")

	baronClient := baron.New("https://api.skinbaron.de", env.GetAPIKey(), 5*time.Second)
	syncOffersUC := usecase.NewSyncOffers(baronClient, repos.Offers, log)
	listOffersUC := usecase.NewListOfferService(repos.Offers, log)

	handler := transportgrpc.NewHandler(syncOffersUC, listOffersUC)
	server := transportgrpc.NewServer(cfg.GRPCConfig.Address, handler)

	log.Info("starting grpc server on: ",
		"address", cfg.GRPCConfig.Address)

	if err := server.Run(); err != nil {
		log.Error("error when running server: ",
			"error", err)
		os.Exit(1)
	}
}

func makeDBConfigData(cfg *config.Config) *db.DBConfigData {
	return &db.DBConfigData{
		ConnUrl:         env.GetDBUrl(),
		ConnTimeout:     cfg.DBConfig.ConnTimeout,
		ConnMaxIdleTime: cfg.DBConfig.ConnMaxIdleTime,
		ConnMaxLifeTime: cfg.DBConfig.ConnMaxLifeTime,
		MaxIdleConns:    cfg.DBConfig.MaxIdleConns,
		MaxOpenConns:    cfg.DBConfig.MaxOpenConns,
	}
}
