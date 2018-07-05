package p2p

import (
	"fmt"
	"net"
	"sync"

	"github.com/blockchainservice/p2p/nat"
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
	natSpec     string
	//	extIP       net.IP
}

// StartListening start server
func (s *Server) StartListening() error {

	listener, err := net.Listen("tcp", s.ListenAddr)
	if err != nil {
		return err
	}
	if err = s.mappingExternalNetwork(); err != nil {
		log.Error(err)
		return err
	}
	// add nat
	s.listener = listener
	s.wg.Add(1)
	go s.listenLoop()
	return nil
}

func (s *Server) mappingExternalNetwork() error {
	natm, err := nat.Parse(s.natSpec)
	if err != nil {
	}
	s.listener.Addr()
	realaddr := s.listener.Addr().(*net.TCPAddr)
	if natm != nil {
		go nat.Map(natm, nil, "tcp", realaddr.Port, realaddr.Port, "ethereum discovery")

		// TODO: react to external IP changes over time.
		if ext, err := natm.GetExternalAddress(); err == nil {
			fmt.Println(ext)
			realaddr = &net.TCPAddr{IP: ext, Port: realaddr.Port}
		}
	}
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
		s.connections <- conn
	}

}

func (s *Server) Connections() <-chan net.Conn {
	return s.connections
}
