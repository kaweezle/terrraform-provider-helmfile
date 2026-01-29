terraform {
  required_version = ">= 1.0"

  required_providers {
    helmfile = {
      source = "kaweezle/helmfile"
    }
  }
}

provider "helmfile" {
  perform_init = true

  additional_plugins = [
    "https://github.com/databus23/helm-diff",
    "https://github.com/jkroepke/helm-secrets"
  ]

  strip_args_values_on_exit_error = true
}

# Ingress Controller
resource "helmfile_release" "ingress" {
  name         = "ingress-nginx"
  file_or_path = "${path.module}/helmfile.yaml"
  environment  = var.environment
  namespace    = "ingress-system"

  selector = ["component=ingress"]

  values = [
    "${path.module}/values/${var.environment}/ingress.yaml"
  ]

  set = [
    "controller.replicaCount=${var.ingress_replica_count}",
    "controller.service.type=${var.ingress_service_type}"
  ]

  wait          = true
  wait_for_jobs = true

  state_values_set = [
    "environment=${var.environment}",
    "cluster=${var.cluster_name}"
  ]
}

# Database
resource "helmfile_release" "database" {
  name         = "postgresql"
  file_or_path = "${path.module}/helmfile.yaml"
  environment  = var.environment
  namespace    = "database"

  selector = ["component=database"]

  values = [
    "${path.module}/values/${var.environment}/database.yaml"
  ]

  set = [
    "persistence.size=${var.database_storage_size}",
    "replication.enabled=${var.database_replication_enabled}"
  ]

  suppress_secrets = true
  wait             = true
  wait_for_jobs    = true

  state_values_set = [
    "environment=${var.environment}",
    "backup_enabled=${var.enable_database_backup}"
  ]
}

# Backend Application
resource "helmfile_release" "backend" {
  name         = "backend-api"
  file_or_path = "${path.module}/helmfile.yaml"
  environment  = var.environment
  namespace    = "application"

  selector = ["component=backend"]

  values = [
    "${path.module}/values/${var.environment}/backend.yaml"
  ]

  set = [
    "replicaCount=${var.backend_replica_count}",
    "image.tag=${var.backend_image_tag}",
    "resources.requests.cpu=${var.backend_cpu}",
    "resources.requests.memory=${var.backend_memory}"
  ]

  wait          = true
  wait_for_jobs = true

  include_needs = true

  state_values_set = [
    "environment=${var.environment}",
    "database_host=${var.database_host}",
    "monitoring_enabled=${var.enable_monitoring}"
  ]

  depends_on = [
    helmfile_release.database
  ]
}

# Frontend Application
resource "helmfile_release" "frontend" {
  name         = "frontend-web"
  file_or_path = "${path.module}/helmfile.yaml"
  environment  = var.environment
  namespace    = "application"

  selector = ["component=frontend"]

  values = [
    "${path.module}/values/${var.environment}/frontend.yaml"
  ]

  set = [
    "replicaCount=${var.frontend_replica_count}",
    "image.tag=${var.frontend_image_tag}",
    "ingress.enabled=true",
    "ingress.host=${var.frontend_host}"
  ]

  wait          = true
  wait_for_jobs = true

  include_needs = true

  state_values_set = [
    "environment=${var.environment}",
    "api_endpoint=${var.backend_api_endpoint}"
  ]

  depends_on = [
    helmfile_release.backend,
    helmfile_release.ingress
  ]
}

# Monitoring Stack (optional)
resource "helmfile_release" "monitoring" {
  count = var.enable_monitoring ? 1 : 0

  name         = "monitoring"
  file_or_path = "${path.module}/helmfile.yaml"
  environment  = var.environment
  namespace    = "monitoring"

  selector = ["component=monitoring"]

  values = [
    "${path.module}/values/${var.environment}/monitoring.yaml"
  ]

  set = [
    "prometheus.retention=${var.metrics_retention}",
    "grafana.adminPassword=${var.grafana_admin_password}"
  ]

  suppress_secrets = true
  wait             = true
  wait_for_jobs    = true

  state_values_set = [
    "environment=${var.environment}",
    "storage_class=${var.monitoring_storage_class}"
  ]
}
