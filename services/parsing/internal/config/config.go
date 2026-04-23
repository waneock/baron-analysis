package config

import (
	"log"
	"os"
	"skinbaron-analyzer/pkg/env"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Config struct {
	Env        string     `yaml:"env" env-required:"true"`
	DBConfig   DBConfig   `yaml:"db_config"`
	GRPCConfig GRPCConfig `yaml:"grpc_config"`
	JSONData   JSONData   `yaml:"json_data"`
}

type DBConfig struct {
	ConnTimeout     time.Duration `yaml:"conn_timeout" env-default:"5s"`
	ConnMaxIdleTime time.Duration `yaml:"conn_max_idle_time" env-default:"5m"`
	ConnMaxLifeTime time.Duration `yaml:"conn_max_life_time" env-default:"30m"`
	MaxIdleConns    int           `yaml:"max_idle_conns" env-default:"10"`
	MaxOpenConns    int           `yaml:"max_open_conns" env-default:"25"`
}

type GRPCConfig struct {
	Address string `yaml:"address" env-default:"0.0.0.0:50051"`
}

type JSONData struct {
	ItemPath string `yaml:"items_path"`
}

func MustLoad() *Config {
	configPath := env.GetConfigPath()

	if configPath == "" {
		log.Fatal("config path not is set")
	}

	_, err := os.Stat(configPath)

	if os.IsNotExist(err) {
		log.Fatal("config file does not exists at path: ", configPath)
	}

	var cfg Config
	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		log.Fatal("cannot read config file at path: ", configPath, err)
	}

	return &cfg
}
