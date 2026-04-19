package parsinggrpc

import (
	"context"
	"fmt"
	pb "skinbaron-analyzer/proto/parsing/v1"

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

func (c *Client) SyncOffers(ctx context.Context) (*pb.SyncOffersResponse, error) {
	return c.client.SyncOffers(ctx, &pb.SyncOffersRequest{})
}

func (c *Client) ListOffers(ctx context.Context, req *pb.ListOffersRequest) (*pb.ListOffersResponse, error) {
	return c.client.ListOffers(ctx, req)
}
