use std::path::PathBuf;
use std::sync::atomic::{AtomicU64, Ordering};
use tokio_uring::fs::File;
use tokio_uring::fs::OpenOptions;

pub struct Storage {
    data_path: PathBuf,
    offset: AtomicU64,
    file: Option<File>,
}

impl Storage {
    pub async fn init(&mut self, path: impl Into<PathBuf>) -> tokio::io::Result<()> {
        self.data_path = path.into();
        self.offset.store(0, Ordering::SeqCst);
        let file = OpenOptions::new()
            .read(true)
            .write(true)
            .create(true)
            .custom_flags(libc::O_DIRECT|libc::O_SYNC)
            .open(&self.data_path)
            .await?;
        self.file = Some(file);
        Ok(())
    }

    pub async fn write(&self) -> tokio::io::Result<()> {
        // guard: return err if file is not opened
        if self.file.is_none() {
            return Err(tokio::io::Error::new(
                tokio::io::ErrorKind::NotFound,
                "File not opened",
            ));
        }

        // genearate a random 1KiB~2KiB buffer to write
        len = rand::random::<u16>() % 1024 + 1024;
        let mut buf = vec![0; len as usize];
        for i in 0..len {
            buf[i as usize] = rand::random::<u8>();
        }

        // write the buffer to the file
        let file = self.file.as_ref().unwrap();
        let offset = self.offset.load(Ordering::SeqCst);
        self.offset.fetch_add(len as u64, Ordering::SeqCst);
        
        let (result, ret_buf) = file.write_all_at(&buf, offset).await;
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

