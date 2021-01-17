package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/pkg/krok/providers"
)

// JwtAuthInterceptor is a grpc.UnaryServerInterceptor that enforcing the JWT authentication from the token provider.
func JwtAuthInterceptor(tokenProvider providers.TokenProvider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		token, err := getHeader(ctx, "authorization")
		if err != nil {
			return ctx, err
		}

		// Validate the token against the TokenProvider.
		// TODO: Do we want to push the claims object into the context?
		_, err = tokenProvider.GetTokenRaw(token)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "failed to get token")
		}

		return handler(ctx, req)
	}
}

func getHeader(ctx context.Context, key string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Internal, "failed to get request metadata")
	}

	header := md.Get(key)
	if header == nil {
		return "", status.Error(codes.Unauthenticated, "failed to get header")
	}

	if len(header) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than one header")
	}

	return header[0], nil
}
