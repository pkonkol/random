export PATH="$PATH:$(go env GOPATH)/bin"
protoc -I../proto/ \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    ../proto/demo.proto

mv demo_grpc.pb.go src/proto                                                                                                      master
mv demo.pb.go src/proto     