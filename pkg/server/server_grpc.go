package server

import "context"

// GRPCKrokServer is a server.
type GRPCKrokServer struct {
	Config
	Dependencies
}

// NewGRPCKrokServer creates a new GRPC krok server.
func NewGRPCKrokServer(cfg Config, deps Dependencies) *GRPCKrokServer {
	return &GRPCKrokServer{Config: cfg, Dependencies: deps}
}

// Run starts up listening.
func (s *GRPCKrokServer) Run(ctx context.Context) error {
	return nil
}
