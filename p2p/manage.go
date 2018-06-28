package p2p

import (
	"fmt"
)

// Manage handles peer connections and exposes an API to receive incoming messages on `Business`
type Manage struct {
	reactors map[string]Reactor
	server   Listener
}

func NewManage(port uint16) {

	// NewServer create
	config := Config{ListenAddr: fmt.Sprintf(":%d", port)}
	server := Server{Config: config}
	server.StartListening()
	manage := Manage{reactors: make(map[string]Reactor), server: &server}
	fmt.Println(manage)

}

func (m *Manage) Start() {
	//m.server.StartListening()
	go m.listenerRoutine(m.server)
}

func (m *Manage) listenerRoutine(l Listener) {
	for {
		inConn, ok := <-l.Connections()
		if !ok {
			break
		}
		fmt.Println(inConn)
	}

	// cleanup
}