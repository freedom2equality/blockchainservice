#!/bin/bash

curl -s -X POST -d '{"jsonrpc":"1.0","method":"hello_world","params":[],"id":1}' http://127.0.0.1:8080
curl -s -X POST -d '{"jsonrpc":"1.0","method":"echo","params":["123456"],"id":1}' http://127.0.0.1:8080
curl -s -X POST -d '{"jsonrpc":"1.0","method":"get_data","params":[[{"content": "987654321"}], {"a":1}, 87978],"id":1}' http://127.0.0.1:8080

curl -s -X POST -d '{"jsonrpc":"1.0","method":"get_data","params":[[{"content": "987654321"}], {"a":1}, 87978, 7654321],"id":1}' http://127.0.0.1:8080