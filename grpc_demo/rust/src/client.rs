use std::time::Duration;
use tokio::sync::mpsc;
use rand::{Rng, SeedableRng};
use rand::rngs::StdRng;
use grpc_demo::test_client::TestClient;
use grpc_demo::{TestRequest, TestReply, NumberStream};
use futures::stream;
use tokio_stream::wrappers::ReceiverStream;
use std::io;

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
    // SendMessage test
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

    // ClientStream test
    let request = tonic::Request::new(stream::iter( [ NumberStream{number: 1}, NumberStream{number: 9}, NumberStream{number: 27}, ]));
    match client.client_stream(request).await {
        Ok(response) => println!("generated res message: {:#?}", response),
        Err(e) => println!("failed to get response for stream, {e}"),
    }

    const BUFFER_SIZE: usize = 10;
    let (mut tx, rx) = mpsc::channel::<NumberStream>(BUFFER_SIZE);
    tokio::spawn(async move {
        let mut rng = StdRng::seed_from_u64(21371488);
        for i in [7, 19, 123] {
            let ms = rng.gen_range(0..10);
            tokio::time::sleep(Duration::from_millis(5*ms)).await;
            tx.send(NumberStream{number: i}).await.unwrap();
        }
        drop(tx);
    });
    match client.client_stream(ReceiverStream::new(rx)).await {
        Ok(response) => println!("generated res message: {:#?}", response),
        Err(e) => println!("failed to get response for stream, {e}"),
    }

    // Chat
    let (mut tx, rx) = mpsc::channel::<NumberStream>(BUFFER_SIZE);
    tokio::spawn(async move {
        for i in 0..3 {
            let mut input = String::new();
            io::stdin().read_line(&mut input);
            println!("input is'{}'", input.trim());
            let number = input.trim().parse::<u64>().unwrap();
            tx.send(NumberStream{number: number}).await.unwrap();
        }
        drop(tx);
    });
    match client.client_stream(ReceiverStream::new(rx)).await {
        Ok(response) => println!("generated res message: {:#?}", response),
        Err(e) => println!("failed to get response for stream, {e}"),
    }


    Ok(())
}
