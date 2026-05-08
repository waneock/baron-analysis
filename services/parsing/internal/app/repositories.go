package app

import (
	"database/sql"
	"skinbaron-analyzer/services/parsing/internal/repository"
	"skinbaron-analyzer/services/parsing/internal/repository/postgres"
)

type Repositories struct {
	Offers           repository.OffersRepository
	Items            repository.ItemsRepository
	ItemWears        repository.ItemWearsRepository
	MarketSyncSource repository.MarketSyncSourceRepository
	ItemWearSale     repository.ItemWearSaleRepository
	ItemSales        repository.ItemSalesRepository
	Jobs             repository.JobRepository
}

func NewRepositories(db *sql.DB) *Repositories {
	return &Repositories{
		Offers:           postgres.NewOffersRepo(db),
		Items:            postgres.NewItemsRepo(db),
		ItemWears:        postgres.NewItemWearsRepo(db),
		MarketSyncSource: postgres.NewMarketSyncSourceRepo(db),
		ItemWearSale:     postgres.NewItemWearSaleRepo(db),
		ItemSales:        postgres.NewItemSalesRepo(db),
		Jobs:             postgres.NewJobRepo(db),
	}
}
