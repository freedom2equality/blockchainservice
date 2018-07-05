package p2p

import (
	"fmt"
	"net"
	"time"
)

// Manage handles peer connections and exposes an API to receive incoming messages on `Business`
type Manage struct {
	reactors map[string]Reactor
	server   Listener
	addpeer  chan *PeerConn
	quit     chan struct{}
}

func NewManage(port uint16) {

	// NewServer create
	config := Config{ListenAddr: fmt.Sprintf(":%d", port)}
	server := Server{Config: config}
	server.StartListening()
	manage := Manage{reactors: make(map[string]Reactor), server: &server}
	fmt.Println(manage)

}

// connect is to connect to other services
func (m *Manage) connect() {

}

func (m *Manage) Start() {
	m.addpeer = make(chan *PeerConn)
	m.quit = make(chan struct{})
	//m.server.StartListening()
	go m.listenerRoutine(m.server)
	go m.run()
}

func (m *Manage) listenerRoutine(l Listener) {
	tokens := 50
	slots := make(chan struct{}, tokens)
	for i := 0; i < tokens; i++ {
		slots <- struct{}{}
	}
	for {
		// Wait for a handshake slot before accepting.
		<-slots
		inConn, ok := <-l.Connections()
		if !ok {
			break
		}
		fmt.Println(inConn)

		//deal inConn
		go func() {
			err := m.inboundPeerConnected(inConn)
			if err != nil {
				log.Error("Ignoring inbound connection: error while adding peer", "address", inConn.RemoteAddr().String(), "err", err)
			}
			slots <- struct{}{}
		}()

	}
}

func (m *Manage) inboundPeerConnected(conn net.Conn) error {
	peerConn, err := newPeerConn(conn, false, false)
	if err != nil {
		conn.Close() // peer is nil
		return err
	}
	m.addPeer(peerConn)
	return nil
}

func (m *Manage) addPeer(conn PeerConn) bool {
	// todo 检查是否存在白名单
	if m.isWhitelisted(conn.conn.RemoteAddr()) {
		log.Errorf("connection from %s dropped (banned)", conn.conn.RemoteAddr().String())
		conn.CloseConn()
		return false
	}
	// todo Timeout
	conn.HandshakeTimeout(30*time.Second, m.addpeer)
	return true
}

func (m *Manage) isWhitelisted(addr net.Addr) bool {
	return false
}

func (m *Manage) run() {
running:
	for {
		select {
		case <-m.quit:
			// The server was stopped. Run the cleanup logic.
			break running
		case c := <-m.addpeer:
			fmt.Println(c)
		}
	}
}
