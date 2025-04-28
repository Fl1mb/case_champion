package grpcclient

import (
	"api_service/internal/grpc/server/user_grpc"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type UserServiceClient struct {
	user_grpc.UnimplementedUserServiceServer
	Client     user_grpc.UserServiceClient
	connection *grpc.ClientConn
}

func (c *UserServiceClient) CreateUser(ctx context.Context, req *user_grpc.CreateUserRequest) (*user_grpc.UserResponse, error) {
	return c.Client.CreateUser(ctx, req)
}

func (c *UserServiceClient) GetUser(ctx context.Context, req *user_grpc.GetUserRequest) (*user_grpc.UserResponse, error) {
	return c.Client.GetUser(ctx, req)
}

func (c *UserServiceClient) Login(ctx context.Context, req *user_grpc.LoginRequest) (*user_grpc.LoginResponse, error) {
	return c.Client.Login(ctx, req)
}

func NewUserServiceClient(addr string) (*UserServiceClient, error) {
	//Make connection with user_service
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("Error to connect to server: %s", err)
	}
	return &UserServiceClient{
		Client:     user_grpc.NewUserServiceClient(conn),
		connection: conn,
	}, nil
}

func (c *UserServiceClient) Close() error {
	return c.connection.Close()
}
