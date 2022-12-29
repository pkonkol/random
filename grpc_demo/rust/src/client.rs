use grpc_demo::test_client::TestClient;
use grpc_demo::{TestRequest, TestReply, NumberStream};

// mod grpc_demo;
pub mod grpc_demo {
    tonic::include_proto!("grpc_demo");
}

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    let channel = tonic::transport::Channel::from_static("http://[::1]:50051")
        .connect()
        .await?;
    let mut client = TestClient::new(channel);
    let msg = "hello from rust".to_string();
    let cnt = 2137;
    println!("sending request from rust {msg}:{cnt}");
    let request = tonic::Request::new(
        TestRequest{
            message: msg,
            counter: cnt,
        }
    );
    let response = client.send_request(request).await?.into_inner();
    println!("reponse: {:?}", response);
    Ok(())
}
