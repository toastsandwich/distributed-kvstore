package server

import (
	"context"
	"fmt"
	"strings"

	"github.com/toastsandwich/kvstore/pkg/kvstore"
	"github.com/toastsandwich/kvstore/pkg/proto"
	"google.golang.org/grpc/codes"
)

type grpcServer struct {
	*proto.UnimplementedServiceKVStoreServer
	s *kvstore.Store
}

func (g *grpcServer) Ping(ctx context.Context, req *proto.PingRequest) (*proto.Response, error) {
	return &proto.Response{Message: "pong", Code: int32(codes.OK)}, nil
}

func (g *grpcServer) Put(ctx context.Context, req *proto.PutRequest) (*proto.Response, error) {
	errs := make([]string, 0)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, p := range req.Pairs {
			if err := g.s.Add(p.Key, p.Val, false); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			errors := strings.Join(errs, ",")
			return &proto.Response{
				Message: errors,
				Code:    int32(codes.Unknown),
			}, fmt.Errorf("error while PUT: %s", errors)
		}
		return &proto.Response{Message: "Created", Code: int32(codes.OK)}, nil
	}
}

func (g *grpcServer) Update(ctx context.Context, req *proto.PutRequest) (*proto.Response, error) {
	errs := make([]string, 0)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, p := range req.Pairs {
			if err := g.s.Update(p.Key, p.Val, false); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			errors := strings.Join(errs, ",")
			return &proto.Response{
				Message: errors,
				Code:    int32(codes.Unknown),
			}, fmt.Errorf("error while PUT: %s", errors)
		}
		return &proto.Response{Message: "Created", Code: int32(codes.OK)}, nil
	}
}

func (g *grpcServer) Delete(ctx context.Context, req *proto.DeleteRequest) (*proto.Response, error) {
	errs := make([]string, 0)
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		for _, k := range req.Keys {
			if err := g.s.Remove(k, false); err != nil {
				errs = append(errs, err.Error())
			}
		}
		if len(errs) > 0 {
			errors := strings.Join(errs, ",")
			return &proto.Response{
				Message: errors,
				Code:    int32(codes.Unknown),
			}, fmt.Errorf("error while DELETE: %s", errors)
		}
		return &proto.Response{Message: "Deleted", Code: int32(codes.OK)}, nil
	}
}
