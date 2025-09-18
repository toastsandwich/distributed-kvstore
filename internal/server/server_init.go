package server

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/toastsandwich/kvstore/internal/bootstrap"
	"github.com/toastsandwich/kvstore/internal/handler"
	"github.com/toastsandwich/kvstore/pkg/kvstore"
	"github.com/toastsandwich/kvstore/pkg/proto"
	"google.golang.org/grpc"
)

var (
	parentCtx context.Context

	store *kvstore.Store
)

func initKVStore() {
	store = kvstore.New()
}

func registerEvents(nodes bootstrap.Nodes) error {
	if nodes == nil {
		return fmt.Errorf("nodes cannot be nill")
	}

	kvstore.AddEvent = func(s string, a any) error {
		wg := &sync.WaitGroup{}
		for _, n := range nodes {
			wg.Add(1)
			go func(n proto.ServiceKVStoreClient) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
				defer cancel()

				req := &proto.PutRequest{
					Pairs: []*proto.KVPair{{Key: s, Val: fmt.Sprint(a)}},
				}
				resp, err := n.Put(ctx, req)
				if err != nil {
					fmt.Println("add event error", err)
					return
				}
				fmt.Printf("AddEvent: Code:%d Message:%s\n", resp.GetCode(), resp.GetMessage())
			}(n)
		}
		wg.Wait()
		return nil
	}

	kvstore.UpdateEvent = func(s string, a any) error {
		wg := &sync.WaitGroup{}
		for _, n := range nodes {
			wg.Add(1)
			go func(n proto.ServiceKVStoreClient) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
				defer cancel()

				req := &proto.PutRequest{
					Pairs: []*proto.KVPair{{Key: s, Val: fmt.Sprint(a)}},
				}
				resp, err := n.Update(ctx, req)
				if err != nil {
					fmt.Println("update event error", err)
					return
				}
				fmt.Printf("UpdateEvent: Code:%d Message:%s\n", resp.GetCode(), resp.GetMessage())
			}(n)
		}
		wg.Wait()
		return nil
	}

	kvstore.RemoveEvent = func(s string, _ any) error {
		wg := &sync.WaitGroup{}
		for _, n := range nodes {
			wg.Add(1)
			go func(n proto.ServiceKVStoreClient) {
				defer wg.Done()

				ctx, cancel := context.WithTimeout(parentCtx, 2*time.Second)
				defer cancel()

				req := &proto.DeleteRequest{
					Keys: []string{s},
				}
				resp, err := n.Delete(ctx, req)
				if err != nil {
					fmt.Println("delete event err", err)
					return
				}
				fmt.Printf("DeleteEvent: Code:%d Message:%s\n", resp.GetCode(), resp.GetMessage())
			}(n)
		}
		wg.Wait()
		return nil
	}
	return nil
}

func initHTTPServer(name string) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: name,

		// Prefork: true,

		CaseSensitive: true,

		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  3 * time.Minute,
	})

	h := handler.New(store).SetGrpcEndpoint(grpcHost)
	app.Get("/endpoint", h.GetGrpcEndpointHandler)
	app.Get("/get", h.GetHandler)
	app.Post("/put", h.PutHandler)
	app.Put("/update", h.UpdateHandler)
	app.Delete("/remove", h.RemoveHandler)

	return app
}

func initGRPCServer() *grpc.Server {
	gs := &grpcServer{
		UnimplementedServiceKVStoreServer: &proto.UnimplementedServiceKVStoreServer{},
		s:                                 store,
	}
	s := grpc.NewServer()

	proto.RegisterServiceKVStoreServer(s, gs)
	return s
}
