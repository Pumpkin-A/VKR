package grpcClient

import (
	"context"
	"log/slog"
	"payment_gateway/config"
	pb "payment_gateway/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"go.opentelemetry.io/otel"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Client struct {
	conn    *grpc.ClientConn
	service pb.PaymentServiceClient
	tr      trace.Tracer
}

func NewPaymentClient(ctx context.Context, cfg config.Config) (*Client, error) {
	// conn, err := grpc.NewClient(cfg.Client.ServerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	// if err != nil {
	// 	return nil, err
	// }

	conn, err := grpc.DialContext(
		ctx,
		cfg.Client.ServerAddress,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),   // добавляем Unary intercepter
		grpc.WithStreamInterceptor(otelgrpc.StreamClientInterceptor()), // добавляем Stream intercepter
	)
	if err != nil {
		slog.Error("client connection error", "err", err.Error())
		return &Client{}, err
	}

	service := pb.NewPaymentServiceClient(conn)

	return &Client{
		conn:    conn,
		service: service,
		tr:      otel.Tracer("grpc-client"),
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetPayment(ctx context.Context, paymentID string) (*pb.PaymentResponse, error) {
	ctx, sp := c.tr.Start(ctx, "Client GetPayment")
	sp.SetAttributes(attribute.String("paymentId", paymentID))
	defer sp.End()

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
