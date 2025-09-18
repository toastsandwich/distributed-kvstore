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
	ctx    context.Context
	cancel context.CancelFunc

	store *kvstore.Store

	addEvent kvstore.EventTrigger
	delEvent kvstore.EventTrigger
)

func initKVStore() {
	store = kvstore.New(ctx)
	kvstore.AddEvent = addEvent
	kvstore.UpdateEvent = addEvent
	kvstore.RemoveEvent = delEvent
}

func registerEvents(nodes bootstrap.Nodes) error {
	if nodes == nil {
		return fmt.Errorf("nodes cannot be nill")
	}
	addEvent = func(ctx context.Context, s string, a any) error {
		wg := &sync.WaitGroup{}
		for _, n := range nodes {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				req := &proto.PutRequest{
					Pairs: []*proto.KVPair{{Key: s, Val: fmt.Sprint(a)}},
				}
				resp, err := n.Put(ctx, req)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("AddEvent: Code:%d Message:%s\n", resp.GetCode(), resp.GetMessage())
			}(wg)
		}
		wg.Wait()
		return nil
	}

	delEvent = func(ctx context.Context, s string, _ any) error {
		wg := &sync.WaitGroup{}
		for _, n := range nodes {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				req := &proto.DeleteRequest{
					Keys: []string{s},
				}
				resp, err := n.Delete(ctx, req)
				if err != nil {
					fmt.Println(err)
					return
				}
				fmt.Printf("DeleteEvent: Code:%d Message:%s\n", resp.GetCode(), resp.GetMessage())
			}(wg)
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
