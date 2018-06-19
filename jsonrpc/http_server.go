package jsonrpc

import (
	"net"
	"net/http"
	"sync"
	"time"
)

// todo Reading and writing separation
const (
	// rpcAuthTimeoutSeconds is the number of seconds a connection to the
	// RPC server is allowed to stay open without authenticating before it
	// is closed.
	rpcAuthTimeoutSeconds = 10
)

type RpcServer struct {
	Listeners []net.Listener
	wg        sync.WaitGroup
}

func (s *RpcServer) Start() {
	rpcServeMux := http.NewServeMux()
	httpServer := &http.Server{
		Handler:     rpcServeMux,
		ReadTimeout: time.Second * rpcAuthTimeoutSeconds,
	}

	rpcServeMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Connection", "close")
		w.Header().Set("Content-Type", "application/json")
		r.Close = true

		// todo: Limit the number of connections to max allowed.

		// todo:User Authentication
		var isAdmin = true

		// Read and respond to the request.
		s.jsonRPCRead(w, r, isAdmin)
	})
	// todo: Websocket endpoint.
	rpcServeMux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
	})
	for _, listener := range s.Listeners {
		s.wg.Add(1)
		go func(listener net.Listener) {

			httpServer.Serve(listener)
			s.wg.Done()
		}(listener)
	}
}

func (s *RpcServer) jsonRPCRead(w http.ResponseWriter, r *http.Request, isAdmin bool) {

}
