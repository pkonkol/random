From WSL
```s
go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
```
Generate python files
```s
python3 -m grpc_tools.protoc -I./grpc_demo_go/proto/ --python_out=. --grpc_python_out=. ./grpc_demo_go/proto/test.proto
```
Generate go files
```s
protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative src/proto/test.proto
```
And run
```s
go run src/server/server.go
go run src/client/client.go
```