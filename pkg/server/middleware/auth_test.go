package middleware

import (
	"context"
	"errors"
	"log"
	"net"
	"testing"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/krok-o/krok/pkg/krok/providers/mocks"
	repov1 "github.com/krok-o/krok/proto/repository/v1"
)

type mockRepoService struct {
	repov1.RepositoryServiceServer
}

func (s *mockRepoService) GetRepository(context.Context, *repov1.GetRepositoryRequest) (*repov1.Repository, error) {
	return &repov1.Repository{}, nil
}

// TODO: Maybe a test against the actual server (e2e) would be better for this? Could even valid REST too then.
func TestJwtAuthInterceptor(t *testing.T) {
	t.Run("valid token success", func(t *testing.T) {
		ctx := context.Background()

		mockOp := &mocks.OAuthProvider{}
		mockOp.On("Verify", "test").Return(jwt.StandardClaims{}, nil)

		conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(mockOp)))
		require.NoError(t, err)
		defer conn.Close()

		client := repov1.NewRepositoryServiceClient(conn)

		md := metadata.MD{
			"authorization": []string{"test"},
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		_, err = client.GetRepository(ctx, &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1)})
		mockOp.AssertExpectations(t)
		assert.NoError(t, err)
	})

	t.Run("no auth header", func(t *testing.T) {
		ctx := context.Background()

		mockOp := &mocks.OAuthProvider{}

		conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(mockOp)))
		require.NoError(t, err)
		defer conn.Close()

		client := repov1.NewRepositoryServiceClient(conn)

		repo, err := client.GetRepository(ctx, &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1)})
		mockOp.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = failed to get header")
		assert.Nil(t, repo)
	})

	t.Run("empty auth header", func(t *testing.T) {
		ctx := context.Background()

		mockOp := &mocks.OAuthProvider{}

		conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(mockOp)))
		require.NoError(t, err)
		defer conn.Close()

		client := repov1.NewRepositoryServiceClient(conn)

		md := metadata.MD{
			"authorization": []string{},
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		repo, err := client.GetRepository(ctx, &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1)})
		mockOp.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = failed to get header")
		assert.Nil(t, repo)
	})

	t.Run("token provider error", func(t *testing.T) {
		ctx := context.Background()

		mockOp := &mocks.OAuthProvider{}
		mockOp.On("Verify", "test").Return(jwt.StandardClaims{}, errors.New("token err"))

		conn, err := grpc.DialContext(ctx, "", grpc.WithInsecure(), grpc.WithContextDialer(dialer(mockOp)))
		require.NoError(t, err)
		defer conn.Close()

		client := repov1.NewRepositoryServiceClient(conn)

		md := metadata.MD{
			"authorization": []string{"test"},
		}
		ctx = metadata.NewOutgoingContext(ctx, md)

		repo, err := client.GetRepository(ctx, &repov1.GetRepositoryRequest{Id: wrapperspb.Int32(1)})
		mockOp.AssertExpectations(t)
		assert.EqualError(t, err, "rpc error: code = Unauthenticated desc = failed to verify token")
		assert.Nil(t, repo)
	})
}

func dialer(op *mocks.OAuthProvider) func(context.Context, string) (net.Conn, error) {
	listener := bufconn.Listen(1024 * 1024)
	server := grpc.NewServer(
		grpc.UnaryInterceptor(JwtAuthInterceptor(op)),
	)

	repov1.RegisterRepositoryServiceServer(server, &mockRepoService{})

	go func() {
		if err := server.Serve(listener); err != nil {
			log.Fatal(err)
		}
	}()

	return func(context.Context, string) (net.Conn, error) {
		return listener.Dial()
	}
}
