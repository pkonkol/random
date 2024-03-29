# Generated by the gRPC Python protocol compiler plugin. DO NOT EDIT!
"""Client and server classes corresponding to protobuf-defined services."""
import grpc

import demo_pb2 as demo__pb2


class TestStub(object):
    """Missing associated documentation comment in .proto file."""

    def __init__(self, channel):
        """Constructor.

        Args:
            channel: A grpc.Channel.
        """
        self.SendRequest = channel.unary_unary(
                '/grpc_demo.Test/SendRequest',
                request_serializer=demo__pb2.TestRequest.SerializeToString,
                response_deserializer=demo__pb2.TestReply.FromString,
                )
        self.ClientStream = channel.stream_unary(
                '/grpc_demo.Test/ClientStream',
                request_serializer=demo__pb2.NumberStream.SerializeToString,
                response_deserializer=demo__pb2.TestReply.FromString,
                )


class TestServicer(object):
    """Missing associated documentation comment in .proto file."""

    def SendRequest(self, request, context):
        """Missing associated documentation comment in .proto file."""
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')

    def ClientStream(self, request_iterator, context):
        """rpc ServerStream (TestRequest) returns (stream TestReply) {}
        """
        context.set_code(grpc.StatusCode.UNIMPLEMENTED)
        context.set_details('Method not implemented!')
        raise NotImplementedError('Method not implemented!')


def add_TestServicer_to_server(servicer, server):
    rpc_method_handlers = {
            'SendRequest': grpc.unary_unary_rpc_method_handler(
                    servicer.SendRequest,
                    request_deserializer=demo__pb2.TestRequest.FromString,
                    response_serializer=demo__pb2.TestReply.SerializeToString,
            ),
            'ClientStream': grpc.stream_unary_rpc_method_handler(
                    servicer.ClientStream,
                    request_deserializer=demo__pb2.NumberStream.FromString,
                    response_serializer=demo__pb2.TestReply.SerializeToString,
            ),
    }
    generic_handler = grpc.method_handlers_generic_handler(
            'grpc_demo.Test', rpc_method_handlers)
    server.add_generic_rpc_handlers((generic_handler,))


 # This class is part of an EXPERIMENTAL API.
class Test(object):
    """Missing associated documentation comment in .proto file."""

    @staticmethod
    def SendRequest(request,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.unary_unary(request, target, '/grpc_demo.Test/SendRequest',
            demo__pb2.TestRequest.SerializeToString,
            demo__pb2.TestReply.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)

    @staticmethod
    def ClientStream(request_iterator,
            target,
            options=(),
            channel_credentials=None,
            call_credentials=None,
            insecure=False,
            compression=None,
            wait_for_ready=None,
            timeout=None,
            metadata=None):
        return grpc.experimental.stream_unary(request_iterator, target, '/grpc_demo.Test/ClientStream',
            demo__pb2.NumberStream.SerializeToString,
            demo__pb2.TestReply.FromString,
            options, channel_credentials,
            insecure, call_credentials, compression, wait_for_ready, timeout, metadata)
