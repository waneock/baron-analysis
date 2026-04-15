package main

import (
	"context"
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/parsing/internal/client/baron"
	"skinbaron-analyzer/services/parsing/internal/config"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"skinbaron-analyzer/services/parsing/internal/usecase"
	"time"
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
	repo := repository.New(db)
	if repo == nil {
		log.Error("error when trying to create a repository")
		os.Exit(1)
	}

	log.Info("app successfully initialized")

	baronClient := baron.New("https://api.skinbaron.de", env.GetAPIKey(), 5*time.Second)
	sales := usecase.New(baronClient, repo.OffersRepository, log)
	ctx := context.Background()
	sales.SyncOffers(ctx)
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
