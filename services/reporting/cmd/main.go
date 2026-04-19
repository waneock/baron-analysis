package main

import (
	"context"
	"fmt"
	"log"
	"skinbaron-analyzer/services/reporting/internal/client/parsinggrpc"
	"skinbaron-analyzer/services/reporting/internal/usecase"
	"time"
)

func main() {
	time.Sleep(60 * time.Second)
	parsingClient, err := parsinggrpc.New("parsing:50051")
	if err != nil {
		log.Fatal("cannot create parsing client: ", err)
	}
	listOffersUC := usecase.NewListOffers(parsingClient)
	appID := 730
	limit := 50
	offset := 20
	ctx := context.Background()
	offers, err := listOffersUC.Execute(ctx, usecase.ListOffersInput{
		AppID:  &appID,
		Limit:  limit,
		Offset: offset,
	})

	if err != nil {
		log.Fatalf("error when trying to list offers: ", err)
	}

	fmt.Println("offers: ", offers)
}
