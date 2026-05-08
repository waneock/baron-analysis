package parsinggrpc

import (
	"context"
	"fmt"
	pb "skinbaron-analyzer/proto/parsing/v1"
	"skinbaron-analyzer/services/reporting/internal/domain"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn   *grpc.ClientConn
	client pb.ParsingServiceClient
}

func New(address string) (*Client, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("parsinggrpc: client: dial error %w", err)
	}

	return &Client{
		conn:   conn,
		client: pb.NewParsingServiceClient(conn),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) ListOffers(ctx context.Context, input domain.ListOffersInput) (*domain.ListOffersOutput, error) {
	pbListOffersRequest := listOffersInputToPbRequest(input)

	pbListOffersResult, err := c.client.ListOffers(ctx, pbListOffersRequest)
	if err != nil {
		return nil, err
	}

	listOffersOutput := pbResponseToListOffersOut(pbListOffersResult)

	return listOffersOutput, nil
}

func (c *Client) ListItemSales(ctx context.Context, input domain.ListItemSalesInput) (*domain.ListItemSalesOutput, error) {
	pbListItemSalesRequest := listItemSalesInputToPbRequest(input)

	pbListItemSalesResponse, err := c.client.ListItemSales(ctx, &pbListItemSalesRequest)
	if err != nil {
		return nil, err
	}

	listItemSalesOutput := pbResponseToListItemSalesOut(pbListItemSalesResponse)

	return &listItemSalesOutput, nil
}

func (c *Client) ListItemSalesStats(ctx context.Context, input domain.ListItemSalesStatInput) (*domain.ListItemSalesStatOutput, error) {
	pbListItemSalesStatRequest := listItemSalesStatsInputToPbRequest(input)

	pbListItemSalesStatResponse, err := c.client.ListItemSalesStats(ctx, &pbListItemSalesStatRequest)
	if err != nil {
		return nil, err
	}

	listItemSalesStatOutput := pbResponseToListItemSalesStatsOut(pbListItemSalesStatResponse)

	return &listItemSalesStatOutput, nil
}
