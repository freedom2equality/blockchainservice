package jsonrpc

import "github.com/blockchainservice/common"

type HelloWorld struct {
}

type Echo struct {
	Content string
}

type Data struct {
	Content string `json:"content"`
}

type GetData struct {
	Data  []Data
	Arr   map[string]int64 `jsonrpcusage:"{\"a\":1,...}"`
	Echo  int64
	Pdata *int64
}

func init() {
	// No special flags for commands in this file.
	flags := common.UsageFlag(0)

	common.MustRegisterCmd("hello_world", (*HelloWorld)(nil), flags)
	common.MustRegisterCmd("echo", (*Echo)(nil), flags)
	common.MustRegisterCmd("get_data", (*GetData)(nil), flags)
}
