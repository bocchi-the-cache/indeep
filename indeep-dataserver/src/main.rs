pub mod config;
pub mod storage;

use std::path::Path;
use axum::{response::Html, routing::get, Router};

#[tokio::main]
async fn main() -> Result<(), Box<dyn std::error::Error>> {
    config::init_config::<&Path>(None)?;

    let app = Router::new().route("/", get(handler));

    let addr =  {
        let cfg = config::get_config();
        format!("{}:{}", cfg.server_config.listen_addr, cfg.server_config.listen_port)
    };

    let listener = tokio::net::TcpListener::bind(addr).await.unwrap();

    println!("Listening on http://{}", listener.local_addr().unwrap());

    axum::serve(listener, app).await?;

    Ok(())
}

async fn handler() -> Html<&'static str> {
    println!("Recv Request.");
    Html("<h1>Hello, World!</h1>")
}