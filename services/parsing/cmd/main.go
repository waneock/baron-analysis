package main

import (
	"log"
	"skinbaron-analyzer/pkg/db"
	"skinbaron-analyzer/pkg/env"
	"skinbaron-analyzer/services/parsing/internal/config"
)

func main() {
	cfg := config.MustLoad()
	log.Println(cfg)

	dbCfg := makeDBConfigData(cfg)

	db, err := db.New(*dbCfg)
	if err != nil {
		log.Fatalf("main: error when trying to create a new db: %v", err)
	}

	defer db.Close()

	log.Println("successfully initialized")
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
