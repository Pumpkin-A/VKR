package api

import (
	"context"
	"log"
	"log/slog"
	"net"
	"transaction_service/config"
	"transaction_service/internal/models"
	pb "transaction_service/pkg/pb/github.com/yourproject/pkg/pb/transaction/v1"

	"google.golang.org/grpc"
)

type PaymentManager interface {
	GetPayment(ctx context.Context, uuid string) (models.Payment, error)
}

type Server struct {
	pb.UnimplementedPaymentServiceServer
	rpcServer *grpc.Server
	pm        PaymentManager
}

func New(pm PaymentManager) *Server {
	s := &Server{
		rpcServer: grpc.NewServer(),
		pm:        pm,
	}
	pb.RegisterPaymentServiceServer(s.rpcServer, s)
	return s
}

func (s *Server) Start(cfg config.Config) error {
	lis, err := net.Listen("tcp", cfg.Server.Port)
	if err != nil {
		slog.Error("error with listen server", "err", err.Error())
		return err
	}
	slog.Info("grpc server was started")
	return s.rpcServer.Serve(lis)
}

func (s *Server) Stop() {
	s.rpcServer.GracefulStop()
}

func (s *Server) GetPayment(ctx context.Context, req *pb.PaymentRequest) (*pb.PaymentResponse, error) {
	transactionID := req.GetPaymentId()
	log.Printf("Received request for transaction ID: %s", transactionID)

	payment, err := s.pm.GetPayment(ctx, transactionID)
	if err != nil {
		return &pb.PaymentResponse{}, err
	}

	grpcPaymnet := models.ConvertPaymentToGrpc(payment)

	return grpcPaymnet, nil
}
