#![allow(unused_imports, dead_code, unused_variables)]
use std::fs;
use std::path::Path;

#[derive(Debug, Clone)]
pub struct NetworkConfig {
    pub default_interface: String,
    pub scan_threads: u32,
    pub timeout_seconds: u64,
    pub monitor_mode: bool,
}

#[derive(Debug, Clone)]
pub struct SecurityConfig {
    pub vulnerability_db_path: String,
    pub risk_threshold: String,
    pub enable_metasploit: bool,
    pub enable_nessus: bool,
}

#[derive(Debug, Clone)]
pub struct ReportingConfig {
    pub company_name: String,
    pub report_template: String,
    pub output_directory: String,
    pub enable_pdf: bool,
    pub enable_html: bool,
}

#[derive(Debug, Clone)]
pub struct ApiConfig {
    pub metasploit_host: String,
    pub metasploit_port: u16,
    pub nessus_host: String,
    pub nessus_port: u16,
}

#[derive(Debug, Clone)]
pub struct AppConfig {
    pub network: NetworkConfig,
    pub security: SecurityConfig,
    pub reporting: ReportingConfig,
    pub api: ApiConfig,
}

impl Default for AppConfig {
    fn default() -> Self {
        Self {
            network: NetworkConfig {
                default_interface: "eth0".to_string(),
                scan_threads: 100,
                timeout_seconds: 5,
                monitor_mode: false,
            },
            security: SecurityConfig {
                vulnerability_db_path: "/var/lib/entynet/vulndb".to_string(),
                risk_threshold: "medium".to_string(),
                enable_metasploit: true,
                enable_nessus: false,
            },
            reporting: ReportingConfig {
                company_name: "Security Team".to_string(),
                report_template: "enterprise".to_string(),
                output_directory: "/var/reports/entynet".to_string(),
                enable_pdf: true,
                enable_html: true,
            },
            api: ApiConfig {
                metasploit_host: "127.0.0.1".to_string(),
                metasploit_port: 55553,
                nessus_host: "127.0.0.1".to_string(),
                nessus_port: 8834,
            },
        }
    }
}

impl AppConfig {
    pub fn load() -> Result<Self, Box<dyn std::error::Error>> {
        // For now, just return default config
        Ok(AppConfig::default())
    }

    pub fn save(&self, _path: &str) -> Result<(), Box<dyn std::error::Error>> {
        println!("Saving config");
        Ok(())
    }

    pub fn validate(&self) -> Result<(), Box<dyn std::error::Error>> {
        if self.network.scan_threads == 0 {
            return Err("Scan threads must be greater than 0".into());
        }
        Ok(())
    }
}
