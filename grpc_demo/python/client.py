import grpc
import demo_pb2
import demo_pb2_grpc

if __name__ == '__main__':
    channel = grpc.insecure_channel('localhost:50051')
    stub = demo_pb2_grpc.TestStub(channel)

    msg = "hello from python"
    cnt = 1
    print(f"sending request {msg}:{cnt}")
    res = stub.SendRequest(demo_pb2.TestRequest(message=msg, counter=cnt))
    print(f"test reply is {res.message} with {res.counter}")