package main

import (
	"fmt"
	"net"
	"time"

	"github.com/blockchainservice/jsonrpc"
)

func main() {

	initLogRotator("./json_rpc.log")
	setLogLevels("debug")
	// test jsonrpc
	listeners := make([]net.Listener, 0, 1)
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println(err)
	}
	listeners = append(listeners, listener)

	jsonRPC := jsonrpc.NewRPCServer(listeners)
	jsonRPCLog.Info("json rpc server start ......")
	jsonRPC.Start()
	for {
		time.Sleep(time.Second * 2)
	}
}
