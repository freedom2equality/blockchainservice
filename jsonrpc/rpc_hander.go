package jsonrpc

type commandHandler func(*RpcServer, interface{}, <-chan struct{}) (interface{}, error)

// rpcHandlers maps RPC command strings to appropriate handler functions.
// This is set by init because help references rpcHandlers and thus causes
// a dependency loop.
var rpcHandlers map[string]commandHandler

var rpcAskWallet = map[string]struct{}{
	"listtransactions": {},
}

var rpcUnimplemented = map[string]struct{}{
	"getmempoolentry": {},
	"getnetworkinfo":  {},
	"getwork":         {},
}
