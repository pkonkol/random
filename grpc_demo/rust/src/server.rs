use tonic::{transport::Server, Request, Response, Status, Streaming};
use grpc_demo::test_server::{Test, TestServer};
use grpc_demo::{TestRequest, TestReply, NumberStream};
use futures::{Stream, StreamExt};

pub mod grpc_demo {
    tonic::include_proto!("grpc_demo");
}

#[derive(Debug,Default)]
pub struct MyTest {}

#[tonic::async_trait]
impl Test for MyTest {
    async fn send_request (
        &self,
        request: Request<TestRequest>,
    ) -> Result<Response<TestReply>, Status> {
        println!("received request: {:?}", request);
        let req_inner = request.into_inner();
        let reply = TestReply {
            message: format!("reply from rust {}", req_inner.message).into(),
            counter: req_inner.counter * 10,
        };

        Ok(Response::new(reply))
    }

    async fn client_stream(
        &self,
        request: Request<Streaming<NumberStream>>,
    ) -> Result<Response<TestReply>, Status> {
        println!("received client stream: {:?}", request);
        let mut stream = request.into_inner();
        let mut sum = 1u64;
        while let Some(req) = stream.next().await {
            let req = req?;
            println!("\treceived stream number {}", req.number);
            sum *= req.number;
        }
        Ok(Response::new(TestReply{
            message: format!("Received number stream, product of numbers is {sum}\n"),
            counter: sum as i32,
        }))
    }
}


#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let addr = "[::1]:50051".parse().unwrap();
    let test = MyTest::default();
    println!("server listening on {addr}");
    Server::builder()
        .add_service(TestServer::new(test))
        .serve(addr)
        .await?;
    Ok(())
}
