package grpccache

import (
	"api_service/internal/grpc/grpc_server"
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
)

type CacheClient struct {
	grpc_server.UnimplementedCacheServiceServer
	Client grpc_server.CacheServiceClient
	conn   *grpc.ClientConn
}

func (c *CacheClient) DeleteUser(ctx context.Context, req *grpc_server.DeleteUserRequest) (
	*grpc_server.DeleteUserResponse, error,
) {
	return c.Client.DeleteUser(ctx, req)
}

func (c *CacheClient) GetUser(ctx context.Context, req *grpc_server.GetUserRequest) (
	*grpc_server.GetUserResponse, error,
) {
	return c.Client.GetUser(ctx, req)
}

func (c *CacheClient) Write(ctx context.Context, req *grpc_server.WriteRequest) (
	*grpc_server.WriteResponse, error,
) {
	return c.Client.Write(ctx, req)
}

func (c *CacheClient) Close() error {
	return c.conn.Close()
}

func NewCacheClient(addr string) (*CacheClient, error) {
	conn, err := grpc.Dial(
		addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithTimeout(5*time.Second),
	)

	if err != nil {
		return nil, fmt.Errorf("Error to connect to server: %s", err)
	}
	return &CacheClient{
		Client: grpc_server.NewCacheServiceClient(conn),
		conn:   conn,
	}, nil
}
