package p2p

import (
	"net"
	"sync"
)

type Listener interface {
	Connections() <-chan net.Conn
	//	String() string
	//	Stop() error
}

// Config Server options.
type Config struct {
	ListenAddr string //fmt.Sprintf(":%d", port)
}

type temporary interface {
	Temporary() bool
}

// isTemporary returns true if err is temporary.
func isTemporary(err error) bool {
	te, ok := err.(temporary)
	return ok && te.Temporary()
}

// Server manages all peer connections. Implements Listener
type Server struct {
	Config
	listener    net.Listener
	wg          sync.WaitGroup
	connections chan net.Conn
}

// StartListening start server
func (s *Server) StartListening() error {
	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	s.listener = listener
	s.wg.Add(1)
	go s.listenLoop()
	return nil
}

func (s *Server) listenLoop() {
	s.wg.Done()
	for {
		var (
			conn net.Conn
			err  error
		)
		for {
			conn, err = s.listener.Accept()
			// 网络客户端程序代码可以使用类型断言判断网络错误是瞬时错误还是永久错误。
			// 在碰到瞬时错误的时候，等待一段时间然后重试。
			if isTemporary(err) {
				log.Debug("Temporary read error", "err", err)
				continue
			} else if err != nil {
				log.Debug("Read error", "err", err)
				close(s.connections)
				return
			}
			break
		}
		// deal conn
		//go s.inboundPeerConnected(conn)
		s.connections <- conn
	}

}


func (s *Server) Connections() <-chan net.Conn {
	return s.connections
}
