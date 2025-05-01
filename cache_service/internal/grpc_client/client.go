package grpcclient

import (
	"cache_service/internal/cache"
	"cache_service/internal/grpc/grpc_server"
	"context"
	"fmt"
	"time"
)

type CacheServiceServer struct {
	grpc_server.UnimplementedCacheServiceServer
	Cch *cache.Cache
}

func (c *CacheServiceServer) DeleteUser(ctx context.Context, req *grpc_server.DeleteUserRequest) (
	*grpc_server.DeleteUserResponse, error,
) {
	if c.Cch.Exists(ctx, req.JwtKey) {
		err := c.Cch.Delete(ctx, req.JwtKey)
		if err != nil {
			return &grpc_server.DeleteUserResponse{
				Success: false,
			}, err
		}
		return &grpc_server.DeleteUserResponse{
			Success: true,
		}, nil
	} else {
		return &grpc_server.DeleteUserResponse{
			Success: false,
		}, fmt.Errorf("Cache doesn't exists")
	}
}

func (c *CacheServiceServer) GetUser(ctx context.Context, req *grpc_server.GetUserRequest) (
	*grpc_server.GetUserResponse, error,
) {
	if c.Cch.Exists(ctx, req.JwtKey) {
		id, login, err := c.Cch.GetData(ctx, req.JwtKey)
		if err != nil {
			return &grpc_server.GetUserResponse{
				Success:   false,
				UserId:    0,
				UserLogin: "",
			}, err
		}
		return &grpc_server.GetUserResponse{
			Success:   true,
			UserLogin: login,
			UserId:    int32(id),
		}, nil
	}
	return &grpc_server.GetUserResponse{
		Success:   false,
		UserId:    0,
		UserLogin: "",
	}, fmt.Errorf("Doesn't exist")
}

func (c *CacheServiceServer) Write(ctx context.Context, req *grpc_server.WriteRequest) (
	*grpc_server.WriteResponse, error,
) {
	if c.Cch.Exists(ctx, req.JwtKey) {
		return &grpc_server.WriteResponse{
			Success: false,
		}, fmt.Errorf("Is already in redis")
	}
	err := c.Cch.SetData(ctx, req.JwtKey, req.UserLogin, req.UserId)
	if err != nil {
		return &grpc_server.WriteResponse{
			Success: false,
		}, fmt.Errorf("Error in redis set")
	}
	_ = c.Cch.SetDeadTime(ctx, req.JwtKey, time.Hour*24)
	return &grpc_server.WriteResponse{
		Success: true,
	}, nil
}
