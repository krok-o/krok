package middleware

import (
	"context"
	"strconv"
	"strings"

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

		token, err := getToken(ctx)
		if err != nil {
			return ctx, err
		}

		claims, err := oauthProvider.Verify(token)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "failed to verify token")
		}

		userID, err := strconv.Atoi(claims.Subject)
		if err != nil {
			return ctx, status.Error(codes.Unauthenticated, "failed to convert sub")
		}
		contextWithClaims := context.WithValue(ctx, UserClaimsCtxKey{}, UserClaims{UserID: userID})

		return handler(contextWithClaims, req)
	}
}

type UserClaimsCtxKey struct{}

type UserClaims struct {
	UserID int
}

func getToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "failed to get request metadata")
	}

	header := md.Get("authorization")
	if header != nil && len(header) == 1 {
		token := strings.TrimPrefix(header[0], "Bearer ")
		return token, nil
	}

	cookie := md.Get("grpcgateway-cookie")
	if cookie == nil {
		return "", status.Error(codes.Unauthenticated, "token not present")
	}

	token := strings.TrimPrefix(cookie[0], "_token_=")
	return token, nil
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
