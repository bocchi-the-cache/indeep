pub mod config;
pub mod storage;

use std::{path::Path, sync::Arc};
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


    let storage = Arc::new(tokio::sync::RwLock::new(storage::Storage::new()));
    {
        let mut storage_lock = storage.write().await;
        storage_lock.init(Path::new("data.db")).await?;
    }

    // forever write to the storage
    loop {
        storage.read().await.write().await?;
        
        tokio::time::sleep(tokio::time::Duration::from_secs(1)).await;
    }

    
    println!("Storage initialized with data path: {}", "data.db");

    Ok(())



}

async fn handler() -> Html<&'static str> {
    println!("Recv Request.");
    Html("<h1>Hello, World!</h1>")
}