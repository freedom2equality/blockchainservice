package jsonrpc

import (
	"fmt"
)

type commandHandler func(*RPCServer, interface{}, <-chan struct{}) (interface{}, error)

// rpcHandlers maps RPC command strings to appropriate handler functions.
// This is set by init because help references rpcHandlers and thus causes
// a dependency loop.
var rpcHandlers = map[string]commandHandler{
	"hello_world": helloWorld,
	"echo":        echo,
	"get_data":    getData,
}

var rpcAskWallet = map[string]struct{}{
	"listtransactions": {},
}

var rpcUnimplemented = map[string]struct{}{
	"getmempoolentry": {},
	"getnetworkinfo":  {},
	"getwork":         {},
}

func helloWorld(s *RPCServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	fmt.Println("hello world")
	return "hello world", nil
}

func echo(s *RPCServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*Echo)

	return c.Content, nil
}

func getData(s *RPCServer, cmd interface{}, closeChan <-chan struct{}) (interface{}, error) {
	c := cmd.(*GetData)
	for _, data := range c.Data {
		fmt.Println(data.Content)
	}
	for key, value := range c.Arr {
		fmt.Println(key, value)
	}
	fmt.Println(c.Echo)
	if c.Pdata != nil {
		pdata := *c.Pdata
		fmt.Println(pdata)
	}

	return "ok", nil
}
