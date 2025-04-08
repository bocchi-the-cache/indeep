use serde::{Deserialize, Serialize};
use std::fs;
use std::path::Path;

use once_cell::sync::Lazy;
use std::sync::RwLock;


#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct GlobalConfig {
    pub server_config: ServerConfig,
    pub storage_config: StorageConfig,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct StorageConfig {
    pub data_path : String,
}

#[derive(Debug, Clone, Serialize, Deserialize)]
pub struct ServerConfig {
    pub listen_addr : String,
    pub listen_port : u16,
}

impl GlobalConfig {
    pub fn from_file<P: AsRef<Path>>(path : P) -> Result<Self, Box<dyn std::error::Error>> {
        let config_content = fs::read_to_string(path)?;
        let config = serde_yaml::from_str(&config_content)?;
    
        Ok(config)
    }

    pub fn default() -> Self {
        GlobalConfig { 
            server_config: ServerConfig {
                listen_addr: "127.0.0.1".to_string(),
                listen_port: 3000,
            },
            storage_config: StorageConfig { 
                data_path: "data/".to_string(), 
            },
        }
    }
}


pub static CONFIG: Lazy<RwLock<GlobalConfig>> = Lazy::new(|| RwLock::new(GlobalConfig::default()));

pub fn init_config<P: AsRef<Path>>(path: Option<P>) -> Result<(), Box<dyn std::error::Error>> {
    let config = if let Some(p) = path {
        GlobalConfig::from_file(p)?
    } else {
        GlobalConfig::default()
    };

    let mut global_config = CONFIG.write().unwrap();
    *global_config = config;

    Ok(())
}

pub fn get_config() -> GlobalConfig {
    let config = CONFIG.read().unwrap();
    config.clone()
}