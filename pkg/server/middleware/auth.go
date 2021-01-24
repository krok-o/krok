package middleware

import (
	"context"
	"strings"

	"github.com/gobwas/glob"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/krok-o/krok/pkg/krok/providers"
)

var allowList = []string{
	"/auth.v1.AuthService/*",
}

// JwtAuthInterceptor is a grpc.UnaryServerInterceptor that enforcing the JWT authentication from the token provider.
func JwtAuthInterceptor(oauthProvider providers.OAuthProvider) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		allowed := false
		for _, p := range allowList {
			if checkAllowed(p, info.FullMethod) {
				allowed = true
				break
			}
		}
		if allowed {
			return handler(ctx, req)
		}

		token, err := getHeader(ctx, "authorization")
		if err != nil {
			return ctx, err
		}
		token = strings.TrimPrefix(token, "Bearer ")

		_, err = oauthProvider.Verify(token)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "failed to verify token")
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

func checkAllowed(pattern, input string) bool {
	if pattern == input {
		return true
	}

	g, err := glob.Compile(pattern, '/')
	if err != nil {
		return false
	}

	return g.Match(input)
}
