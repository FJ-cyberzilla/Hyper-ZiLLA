# infrastructure/terraform/main.tf
"""
Main Terraform configuration for enterprise OSINT platform
"""

# Kubernetes Cluster
resource "google_container_cluster" "osint_platform" {
  name               = "osint-platform-${var.environment}"
  location           = var.region
  initial_node_count = 3

  node_config {
    machine_type = "e2-standard-4"
    disk_size_gb = 100
    
    oauth_scopes = [
      "https://www.googleapis.com/auth/cloud-platform",
      "https://www.googleapis.com/auth/logging.write",
      "https://www.googleapis.com/auth/monitoring"
    ]

    labels = {
      environment = var.environment
      workload    = "osint-platform"
    }

    tags = ["osint-platform", var.environment]
  }

  network_policy {
    enabled = true
  }

  private_cluster_config {
    enable_private_nodes    = true
    enable_private_endpoint = false
    master_ipv4_cidr_block  = "172.16.0.0/28"
  }

  monitoring_service = "monitoring.googleapis.com/kubernetes"
  logging_service    = "logging.googleapis.com/kubernetes"

  maintenance_policy {
    daily_maintenance_window {
      start_time = "02:00"
    }
  }
}

# Cloud SQL Database
resource "google_sql_database_instance" "osint_database" {
  name             = "osint-db-${var.environment}"
  database_version = "POSTGRES_14"
  region           = var.region

  settings {
    tier = "db-custom-4-16384"
    
    ip_configuration {
      ipv4_enabled    = false
      private_network = google_compute_network.osint_network.id
    }

    backup_configuration {
      enabled    = true
      start_time = "23:00"
    }

    maintenance_window {
      day  = 7
      hour = 3
    }
  }

  deletion_protection = false
}

# Storage Buckets
resource "google_storage_bucket" "data_lake" {
  name          = "osint-data-lake-${var.environment}"
  location      = var.region
  storage_class = "STANDARD"

  versioning {
    enabled = true
  }

  lifecycle_rule {
    condition {
      age = 30
    }
    action {
      type = "SetStorageClass"
      storage_class = "ARCHIVE"
    }
  }

  encryption {
    default_kms_key_name = google_kms_crypto_key.data_encryption.id
  }
}

# Networking
resource "google_compute_network" "osint_network" {
  name                    = "osint-network-${var.environment}"
  auto_create_subnetworks = false
}

resource "google_compute_subnetwork" "osint_subnet" {
  name          = "osint-subnet-${var.environment}"
  ip_cidr_range = "10.0.0.0/16"
  region        = var.region
  network       = google_compute_network.osint_network.id
}

# Security
resource "google_kms_key_ring" "encryption_keys" {
  name     = "osint-encryption-${var.environment}"
  location = "global"
}

resource "google_kms_crypto_key" "data_encryption" {
  name            = "osint-data-encryption"
  key_ring        = google_kms_key_ring.encryption_keys.id
  rotation_period = "7776000s" # 90 days

  version_template {
    algorithm        = "GOOGLE_SYMMETRIC_ENCRYPTION"
    protection_level = "HSM"
  }
}
