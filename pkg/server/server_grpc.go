package server

import (
	"context"
	"fmt"
	"net"

	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/krok-o/krok/pkg/protos"
)

// GRPCKrokServer is a server.
type GRPCKrokServer struct {
	Config
	Dependencies
}

// HandleHooks handles hooks.
func (s *GRPCKrokServer) HandleHooks(ctx context.Context, request *protos.HandleHooksRequest) (*protos.HandleHooksResponse, error) {
	panic("implement me")
}

// NewGRPCKrokServer creates a new GRPC krok server.
func NewGRPCKrokServer(cfg Config, deps Dependencies) *GRPCKrokServer {
	return &GRPCKrokServer{Config: cfg, Dependencies: deps}
}

// Run starts up listening.
func (s *GRPCKrokServer) Run(ctx context.Context) error {
	// setup grpc server details
	grpcLis, err := net.Listen("tcp:", ":50051")
	if err != nil {
		return fmt.Errorf("failed to listen on address %s: %v", ":50051", err)
	}
	grpcSrv := grpc.NewServer(
		grpc.StreamInterceptor(grpc_prometheus.StreamServerInterceptor),
		grpc.UnaryInterceptor(grpc_prometheus.UnaryServerInterceptor),
	)
	reflection.Register(grpcSrv)

	// setup krok server details
	protos.RegisterKrokServiceServer(grpcSrv, &GRPCKrokServer{})
	g, _ := errgroup.WithContext(ctx)
	g.Go(func() error {
		if err := grpcSrv.Serve(grpcLis); err != nil {
			s.Logger.Error().Err(err).Msg("unable to start grpc server")
			return err
		}
		return nil
	})
	return g.Wait()
}
