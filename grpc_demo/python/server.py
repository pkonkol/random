from concurrent import futures
import grpc
import demo_pb2
import demo_pb2_grpc

class TestServicer(demo_pb2_grpc.TestServicer):
    def SendRequest(self, request, context):
        # return super().SendRequest(request, context)
        print(f"received request: {request.message} with {request.counter}")
        return demo_pb2.TestReply(message="reply from python", counter=request.counter+1)

    def ClientStream(self, request_iterator, context):
        return super().ClientStream(request_iterator, context)

def serve():
    server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
    demo_pb2_grpc.add_TestServicer_to_server(TestServicer(), server)
    server.add_insecure_port('[::]:50051')
    server.start()
    server.wait_for_termination()

if __name__ == '__main__':
    print("starting server")
    serve()