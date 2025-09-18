package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/toastsandwich/kvstore/pkg/kvstore"
	"github.com/toastsandwich/kvstore/pkg/proto"
)

type grpcServer struct {
	*proto.UnimplementedServiceKVStoreServer
	s *kvstore.Store
}

func (g *grpcServer) Ping(ctx context.Context, req *proto.PingRequest) (*proto.Response, error) {
	return &proto.Response{Message: "pong", Code: 200}, nil
}

func (g *grpcServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Response, error) {
	errs := make([]string, 0)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, p := range req.Pairs {
			if err := g.s.Add(p.Key, p.Val); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			errors := strings.Join(errs, ",")
			return &proto.Response{
				Message: errors,
				Code:    fiber.StatusInternalServerError,
			}, fmt.Errorf("error while PUT: %s", errors)
		}
		return &proto.Response{Message: "Created", Code: fiber.StatusCreated}, nil
	}
}

func (g *grpcServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.Response, error) {
	errs := make([]string, 0)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, k := range req.Keys {
			if err := g.s.Remove(k); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			errors := strings.Join(errs, ",")
			return &proto.Response{
				Message: errors,
				Code:    fiber.StatusInternalServerError,
			}, fmt.Errorf("error while DELETE: %s", errors)
		}
		return &proto.Response{Message: "Deleted", Code: fiber.StatusOK}, nil
	}
}
