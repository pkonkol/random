grpcurl -plaintext -import-path ./proto -proto demo.proto \
    -d '{"message": "hello from grpcurl", "counter": -2048}' '[::]:50051' \
    grpc_demo.Test/SendRequest

grpcurl -plaintext -import-path ./proto -proto demo.proto \
    -d @ '[::]:50051' \
    grpc_demo.Test/ClientStream \
    <<EOM
    {"number": 1}
    {"number": 2}
    {"number": 1}
EOM