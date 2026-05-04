package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/pkg/logger"
	"skinbaron-analyzer/services/reporting/internal/client/parsinggrpc"
	"skinbaron-analyzer/services/reporting/internal/config"
	"skinbaron-analyzer/services/reporting/internal/producer/kafka"
	"skinbaron-analyzer/services/reporting/internal/repository"
	httphndl "skinbaron-analyzer/services/reporting/internal/transport/http"
	"skinbaron-analyzer/services/reporting/internal/usecase"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	kafkago "github.com/segmentio/kafka-go"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
	// add startup delay, give kafka time to start
	// TODO: find another solution
	time.Sleep(60 * time.Second)
	cfg := config.MustLoad()

	if cfg == nil {
		log.Fatal("config is nil")
	}

	// logger
	log := logger.MustLoad(cfg.Env)

	// database
	dbCfg := makeDBConfigData(cfg)

	db, err := db.New(*dbCfg)
	if err != nil {
		log.Error("error when trying to connect to a database",
			"error", err)
		os.Exit(1)
	}

	defer db.Close()

	repo := repository.New(db)
	if repo == nil {
		log.Error("cannot create a repository")
		os.Exit(1)
	}

	writer := &kafkago.Writer{
		Addr:     kafkago.TCP("kafka:29092"),
		Balancer: &kafkago.LeastBytes{},
	}

	jobsProducer := kafka.NewProducer(writer)
	defer writer.Close()

	parsingClient, err := parsinggrpc.New("parsing:50051")
	if err != nil {
		log.Error("cannot create parsing client",
			"error", err)
		os.Exit(1)
	}

	listOffersUC := usecase.NewListOffers(parsingClient)
	syncOffersUC := usecase.NewSyncOffers(repo.JobsRepo, jobsProducer)

	httpHandler := httphndl.NewOffersHandler(syncOffersUC, listOffersUC, log)

	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Route("/v1", func(r chi.Router) {
		r.Route("/offers", func(r chi.Router) {
			r.Get("/list", httpHandler.ListOffers)
			r.Post("/sync", httpHandler.SyncOffers)
		})

		r.Route("/items", func(r chi.Router) {
			r.Post("/sync", httpHandler.SyncItems)
		})

		r.Route("/sales", func(r chi.Router) {
			r.Post("/sync", httpHandler.SyncItemSales)
		})
	})

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./web/index.html")
	})

	server := http.Server{
		Addr:    "0.0.0.0:8080",
		Handler: r,
	}

	fmt.Println("run server on: ", server.Addr)

	if err = server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Error("server stopped",
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
