variable "environment" {
  description = "Environment name (development, staging, production)"
  type        = string
  default     = "development"

  validation {
    condition     = contains(["development", "staging", "production"], var.environment)
    error_message = "Environment must be one of: development, staging, production"
  }
}

variable "cluster_name" {
  description = "Name of the Kubernetes cluster"
  type        = string
  default     = "my-cluster"
}

# Ingress Configuration
variable "ingress_replica_count" {
  description = "Number of ingress controller replicas"
  type        = number
  default     = 2
}

variable "ingress_service_type" {
  description = "Type of service for ingress controller"
  type        = string
  default     = "LoadBalancer"
}

# Database Configuration
variable "database_storage_size" {
  description = "Size of database persistent volume"
  type        = string
  default     = "10Gi"
}

variable "database_replication_enabled" {
  description = "Enable database replication"
  type        = bool
  default     = false
}

variable "enable_database_backup" {
  description = "Enable automated database backups"
  type        = bool
  default     = true
}

variable "database_host" {
  description = "Database host for backend connection"
  type        = string
  default     = "postgresql.database.svc.cluster.local"
}

# Backend Configuration
variable "backend_replica_count" {
  description = "Number of backend API replicas"
  type        = number
  default     = 3
}

variable "backend_image_tag" {
  description = "Docker image tag for backend"
  type        = string
  default     = "latest"
}

variable "backend_cpu" {
  description = "CPU request for backend pods"
  type        = string
  default     = "500m"
}

variable "backend_memory" {
  description = "Memory request for backend pods"
  type        = string
  default     = "512Mi"
}

variable "backend_api_endpoint" {
  description = "Backend API endpoint URL"
  type        = string
  default     = "http://backend-api.application.svc.cluster.local"
}

# Frontend Configuration
variable "frontend_replica_count" {
  description = "Number of frontend replicas"
  type        = number
  default     = 2
}

variable "frontend_image_tag" {
  description = "Docker image tag for frontend"
  type        = string
  default     = "latest"
}

variable "frontend_host" {
  description = "Hostname for frontend ingress"
  type        = string
  default     = "app.example.com"
}

# Monitoring Configuration
variable "enable_monitoring" {
  description = "Enable monitoring stack deployment"
  type        = bool
  default     = false
}

variable "metrics_retention" {
  description = "Prometheus metrics retention period"
  type        = string
  default     = "15d"
}

variable "grafana_admin_password" {
  description = "Grafana admin password"
  type        = string
  sensitive   = true
  default     = "change-me"
}

variable "monitoring_storage_class" {
  description = "Storage class for monitoring persistent volumes"
  type        = string
  default     = "standard"
}
