package server

import (
	"context"
	"fmt"
	"net"
	"os"
	"sync"
	"time"

	"github.com/toastsandwich/kvstore/internal/bootstrap"
	"github.com/toastsandwich/kvstore/internal/config"
)

var (
	grpcHost string
)

func (s *Server) findFreePort() (int, error) {
	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return 0, err
	}
	defer ln.Close()

	return ln.Addr().(*net.TCPAddr).Port, nil

}

type Server struct {
	cfg config.Config

	wg *sync.WaitGroup
}

func New(cfg config.Config) *Server {
	return &Server{
		cfg: cfg,
		wg:  &sync.WaitGroup{},
	}
}

func (s *Server) Start() {
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	initKVStore()

	app := initHTTPServer(s.cfg.Name)

	port, err := s.findFreePort()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	grpcHost = net.JoinHostPort(s.cfg.Host.Addr, fmt.Sprint(port))
	grpc := initGRPCServer()

	s.wg.Add(1)
	go s.startHTTPServer(app)

	s.wg.Add(1)
	go s.startGRPCServer(grpc)

	nodes := bootstrap.Init(s.cfg.Peer)
	defer bootstrap.Close()

	if err := registerEvents(nodes); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	s.wg.Wait()
	fmt.Println("kvstore stopped successfully")
}
