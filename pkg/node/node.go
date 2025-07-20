package node

import (
	"net"

	"go.etcd.io/bbolt"

	"github.com/toastsandwich/kvstore/internal"
)

type Node interface {
	listen()      // listen for connections
	Start()       // start the node
	Stop()        // gracefully stop the node
	watchErrors() // watch errors
}

type ServerNodeConfig struct {
	Addr string
	Path string
}

type ServerNode struct {
	ln             net.Listener
	connectionPool map[net.Conn]struct{}
	store          *bbolt.DB

	quitch  chan struct{}
	errorCh chan error
}

// Creates a  new server node
func NewServerNode(config ServerNodeConfig) (*ServerNode, error) {
	store, err := bbolt.Open(config.Path, 0600, nil)
	if err != nil {
		return nil, err
	}

	ln, err := net.Listen("tcp", config.Addr)
	if err != nil {
		return nil, err
	}

	server := &ServerNode{
		ln:             ln,
		connectionPool: make(map[net.Conn]struct{}),
		store:          store,

		errorCh: make(chan error),
		quitch:  make(chan struct{}),
	}
	return server, nil
}

func (s *ServerNode) listen() {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			s.errorCh <- internal.NewNodeEror(err, false)
		}
		s.connectionPool[conn] = struct{}{}
		go s.readConn(conn)
	}
}

func (s *ServerNode) readConn(c net.Conn) {
	defer func() {
		delete(s.connectionPool, c) // remove connection from connectionPool
		c.Close()
	}()
	loop := true
	// TODO: implement a reader for connection with our protocol

	// for loop {}
}

func (s *ServerNode) watchErrors() {
	for e := range s.errorCh {
		if e != nil {
			if nodeErr, ok := e.(internal.NodeError); ok && nodeErr.IsFatal {
				close(s.quitch)
			}
		}
	}
}

func (s *ServerNode) Start() {
	go s.listen()      // start listening for connections
	go s.watchErrors() // start watching for errors happening

	<-s.quitch
	s.Stop()
}

func (s *ServerNode) Stop() {
}
