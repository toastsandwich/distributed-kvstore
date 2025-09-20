package server

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

func (s *Server) startHTTPServer(app *fiber.App) {
	defer s.wg.Done()

	sigCh := make(chan os.Signal, 1)

	signal.Notify(sigCh, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)
	go func() {
		if err := app.ListenTLS(net.JoinHostPort(s.cfg.Host.Addr, s.cfg.Host.Port), s.cfg.TLS.Cert, s.cfg.TLS.Key); err != nil {
			fmt.Println(err)
		}
	}()
	<-sigCh

	if err := app.Shutdown(); err != nil {
		fmt.Println("server closed with error")
		return
	}
	fmt.Println("server closed successfully")
}

func (s *Server) startGRPCServer(g *grpc.Server) {
	defer s.wg.Done()

	ln, err := net.Listen("tcp", grpcHost)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ln.Close()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, os.Interrupt, syscall.SIGTERM)

	go func() {
		fmt.Println("started grpc server on", grpcHost)
		if err := g.Serve(ln); err != nil {
			fmt.Println(err)
			return
		}
	}()

	<-sigCh
	g.Stop()
	fmt.Println("stopped grpc server")
}
