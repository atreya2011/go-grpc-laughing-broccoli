package server

import (
	"context"
	"log"
	"sync"

	"github.com/gofrs/uuid"
	grpc_auth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pbExample "github.com/atreya2011/go-grpc-laughing-broccoli/proto"
)

// Backend implements the protobuf interface
type Backend struct {
	mu    *sync.RWMutex
	users []*pbExample.User
}

// New initializes a new Backend struct.
func New() *Backend {
	return &Backend{
		mu: &sync.RWMutex{},
	}
}

// AddUser adds a user to the in-memory store.
func (b *Backend) AddUser(ctx context.Context, _ *pbExample.AddUserRequest) (*pbExample.User, error) {
	log.Println("user claims:", ctx.Value(userClaims{}))

	b.mu.Lock()
	defer b.mu.Unlock()

	user := &pbExample.User{
		Id: uuid.Must(uuid.NewV4()).String(),
	}
	b.users = append(b.users, user)

	return user, nil
}

// ListUsers lists all users in the store.
func (b *Backend) ListUsers(_ *pbExample.ListUsersRequest, srv pbExample.UserService_ListUsersServer) error {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, user := range b.users {
		err := srv.Send(user)
		if err != nil {
			return err
		}
	}

	return nil
}

type userClaims struct {
	id   int
	name string
}

// parse jwt to extract user claims
func parseToken(token string) (userClaims, error) {
	log.Println("parsing received token:", token)
	return userClaims{
		id:   1,
		name: "test",
	}, nil
}

// ExampleAuthFunc is used by a middleware to authenticate requests
func ExampleAuthFunc(ctx context.Context) (context.Context, error) {
	// extract bearer token from context
	token, err := grpc_auth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, err
	}

	log.Println("got token:", token)

	// parse token to extract user claims
	tokenInfo, err := parseToken(token)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %v", err)
	}

	log.Println("got token info:", tokenInfo)

	// pass it downstream handlers for retrieval
	newCtx := context.WithValue(ctx, userClaims{}, tokenInfo)

	return newCtx, nil
}
