package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"time"

	"github.com/mcfe19/redis-clone/client"
)

const defaultListenAddr = ":5001"

type Config struct {
	ListenAddress string
}
type Server struct {
	Config
	peers map[*Peer]bool
	ln net.Listener
	addPeerCh chan *Peer
	quitCh chan struct{}
	msgCh chan []byte
  kv *KV

}

func NewServer(cfg Config) *Server {
	if len(cfg.ListenAddress) == 0 {
		cfg.ListenAddress = defaultListenAddr
	}
	return &Server{
		Config: cfg,
		peers: make(map[*Peer]bool), 
		addPeerCh: make(chan *Peer),
		quitCh: make(chan struct{}),
		msgCh: make(chan []byte),
    kv: NewKV(),
	}
}

func (s *Server) Start() error {
	ln, err := net.Listen("tcp", s.ListenAddress)
	if err != nil {
		return err
	}
	s.ln = ln
	
	go s.loop()

	slog.Info("server running", "listenAddr", s.ListenAddress)

	return s.acceptLoop()
}

func (s *Server) handleRawMessage(rawMsg []byte) error {
	cmd, err := parseCommand(string(rawMsg))
	if err != nil {
		return err
	}
	switch v := cmd.(type) {
	case SetCommand:
    s.kv.Set(v.key, v.val)
	}
	return nil
}

func (s *Server) loop() {
	for {
		select {
		case rawMsg := <- s.msgCh:
			if err := s.handleRawMessage(rawMsg); err != nil {
				slog.Error("raw message error", "err", err)
			}
		case peer := <- s.addPeerCh:
			s.peers[peer] = true
		case <- s.quitCh:
			return
		}
	}
}

func (s *Server) acceptLoop() error {
	for {
		conn, err := s.ln.Accept()
		if err != nil {
			slog.Error("accept error", "err", err)
			continue
		}
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	peer := NewPeer(conn, s.msgCh)
	s.addPeerCh <- peer

	slog.Info("new peer connected", "remoteAddr", conn.RemoteAddr())

	if err := peer.readLoop(); err != nil {
		slog.Error("peer read error", err, "remoteAddr", conn.RemoteAddr())
	}
}

func main() {
  server := NewServer(Config{})
	go func() {
		log.Fatal(server.Start())
	}()
	time.Sleep(time.Second)

  client := client.New("localhost:5001")
	for i := 0; i < 10; i++ {
		if err := client.Set(context.Background(), fmt.Sprintf("foo %d", i), fmt.Sprintf("bar %d", i)); err != nil {
			log.Fatal(err)
		}
	}

	time.Sleep(time.Second)
  fmt.Println(server.kv.data)
	
}