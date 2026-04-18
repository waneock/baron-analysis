package main

import (
	"context"
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/parsing/internal/app"
	"skinbaron-analyzer/services/parsing/internal/config"
	"skinbaron-analyzer/services/parsing/internal/usecase"
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

	// baronClient := baron.New("https://api.skinbaron.de", env.GetAPIKey(), 5*time.Second)
	// sales := usecase.NewGetSalesService(baronClient, repos.Offers, log)
	// ctx := context.Background()
	// sales.SyncOffers(ctx)

	ctx := context.Background()
	listOffersSvc := usecase.NewListOfferService(repos.Offers, log)
	ctx = context.Background()
	appId := 730
	state := 2
	filter := usecase.ListOffersInput{
		Limit:  100,
		Offset: 0,
		AppID:  &appId,

		State: &state,
	}
	listOffers, err := listOffersSvc.GetOffers(ctx, filter)
	if err != nil {
		log.Info("get offers",
			"error", err)
	}

	log.Info("list offers: ", "offers", listOffers)
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
