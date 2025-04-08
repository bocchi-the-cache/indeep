use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};
use tokio_uring::fs::OpenOptions;
use std::sync::Arc;
use rand::Rng;

#[derive(Clone)]
pub struct Storage {
    data_path: Arc<PathBuf>,
    offset: Arc<AtomicU64>,
}


impl Storage {
    pub fn new() -> Self {
        Storage {
            data_path: Arc::new(PathBuf::new()),
            offset: Arc::new(AtomicU64::new(0)),
        }
    }

    pub async fn init(&mut self, path: impl Into<PathBuf>) -> tokio::io::Result<()> {
        self.data_path = Arc::new(path.into());
        self.offset = Arc::new(AtomicU64::new(0));
        
        let file = OpenOptions::new()
            .read(true)
            .write(true)
            .create(true)
            .open(&*self.data_path)
            .await?;
            
        drop(file);
        
        Ok(())
    }

    pub async fn write(&self) -> tokio::io::Result<()> {
        let file = OpenOptions::new()
            .read(true)
            .write(true)
            .open(&*self.data_path)
            .await?;

        // generate a random 1KiB~2KiB buffer to write
        let mut rng = rand::rng();
        let len = rng.random_range(1024..2048);
        let mut buf = vec![0; len as usize];
        rng.fill(&mut buf[..]);

        // write the buffer to the file
        let offset = self.offset.load(Ordering::SeqCst);
        self.offset.fetch_add(len as u64, Ordering::SeqCst);
        
        let (result, ret_buf) = file.write_all_at(buf, offset).await;
        match result {
            Ok(_) => {
                println!("Wrote {} bytes to file at offset {}", ret_buf.len(), offset);
                Ok(())
            }
            Err(e) => {
                println!("Error writing to file: {}", e);
                Err(e)
            }
        }
    }
}