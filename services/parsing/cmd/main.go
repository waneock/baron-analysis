package main

import (
	"context"
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/parsing/internal/app"
	"skinbaron-analyzer/services/parsing/internal/client/baron"
	"skinbaron-analyzer/services/parsing/internal/config"
	"skinbaron-analyzer/services/parsing/internal/consumer/kafka"
	transportgrpc "skinbaron-analyzer/services/parsing/internal/transport/grpc"
	"skinbaron-analyzer/services/parsing/internal/usecase"
	"time"

	kafkago "github.com/segmentio/kafka-go"
)

func main() {
	// add startup delay, give kafka time to start
	// TODO: find another solution
	time.Sleep(60 * time.Second)
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

	baronClient := baron.New("https://api.skinbaron.de", env.GetAPIKey(), 30*time.Second)

	// syncItemPrices := usecase.NewSyncItemPrices(
	// 	repos.Items,
	// 	repos.MarketSyncSource,
	// 	repos.ItemWearSale,
	// 	baronClient,
	// 	log,
	// )

	// ctx := context.Background()
	// err = syncItemPrices.Execute(ctx)
	// if err != nil {
	// 	log.Error("error happens during sync items",
	// 		"err", err)
	// }

	syncOffersUC := usecase.NewSyncOffers(baronClient, repos.Offers, repos.Jobs, log)
	// ctx := context.Background()
	// syncOffersUC.Execute(ctx)

	reader := kafkago.NewReader(kafkago.ReaderConfig{
		Brokers: []string{"kafka:29092"},
		Topic:   "sync.jobs.requested",
		GroupID: "parsing-sync-workers",
	})

	jobsHandler := kafka.NewJobsEventHandler(syncOffersUC)
	consumer := kafka.NewConsumer(reader, jobsHandler, log)

	go func() {
		if err := consumer.Run(context.Background()); err != nil {
			log.Error("kafka consumer stopped", "error", err)
		}
	}()

	listOffersUC := usecase.NewListOfferService(repos.Offers, log)

	handler := transportgrpc.NewHandler(listOffersUC)
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
