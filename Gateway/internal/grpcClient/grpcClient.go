package grpcClient

import (
	"context"
	"payment_gateway/config"
	pb "payment_gateway/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.PaymentServiceClient
}

func NewPaymentClient(cfg config.Config) (*Client, error) {
	conn, err := grpc.NewClient(cfg.Client.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	service := pb.NewPaymentServiceClient(conn)

	return &Client{
		conn:    conn,
		service: service,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetPayment(ctx context.Context, paymentID string) (*pb.PaymentResponse, error) {
	req := &pb.PaymentRequest{
		PaymentId: paymentID,
	}

	return c.service.GetPayment(ctx, req)
}

func (c *Client) GetPaymentWithTimeout(paymentID string, timeout time.Duration) (*pb.PaymentResponse, error) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	return c.GetPayment(ctx, paymentID)
}
