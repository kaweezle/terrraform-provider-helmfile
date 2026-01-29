output "ingress_release_name" {
  description = "Name of the ingress release"
  value       = helmfile_release.ingress.name
}

output "database_release_name" {
  description = "Name of the database release"
  value       = helmfile_release.database.name
}

output "backend_release_name" {
  description = "Name of the backend release"
  value       = helmfile_release.backend.name
}

output "frontend_release_name" {
  description = "Name of the frontend release"
  value       = helmfile_release.frontend.name
}

output "frontend_url" {
  description = "Frontend application URL"
  value       = "https://${var.frontend_host}"
}

output "monitoring_enabled" {
  description = "Whether monitoring stack is deployed"
  value       = var.enable_monitoring
}

output "environment" {
  description = "Deployed environment"
  value       = var.environment
}

output "cluster_name" {
  description = "Target cluster name"
  value       = var.cluster_name
}
