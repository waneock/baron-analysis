package app

import (
	"database/sql"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"skinbaron-analyzer/services/parsing/internal/repository/postgres"
)

type Repositories struct {
	Offers repository.OffersRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Offers: postgres.NewOffersRepo(db),
	}
}
