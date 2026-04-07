package main

import (
	"log"
	"skinbaron-analyzer/pkg/config"
)

func main() {
	cfg := config.MustLoad()

	log.Println(cfg)
}
